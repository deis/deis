package update

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/go-omaha/omaha"
	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/constant"
	"github.com/deis/deis/deisctl/lock"
	"github.com/deis/deis/deisctl/utils"
)

type Client struct {
	ID             string
	SessionID      string
	Version        string
	AppID          string
	Track          string
	config         *serverConfig
	errorRate      int
	pingsRemaining int
	lock           *lock.Lock
}

func (c *Client) Logf(format string, args ...interface{}) {
	format = c.ID + ": " + format
	fmt.Printf(format, args...)
}

func (c *Client) failed(tag string, err error) {
	c.Logf("%s %v\n", tag, err)
	c.MakeRequest("3", "0", false, false)
}

func (c *Client) getCodebaseURL(uc *omaha.UpdateCheck) string {
	return uc.Urls.Urls[0].CodeBase + uc.Manifest.Packages.Packages[0].Name
}

func (c *Client) RequestLock() {
	elc, err := lock.NewEtcdLockClient(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing etcd client:", err)
	}
	var mID string
	mID = utils.GetMachineID("/")
	if mID == "" {
		fmt.Fprintln(os.Stderr, "Cannot read machine-id")
	}
	c.lock = lock.New(mID, elc)
}

func (c *Client) OmahaRequest(otype, result string, updateCheck, isPing bool) *omaha.Request {
	req := omaha.NewRequest("lsb", "CoreOS", "", "")
	app := req.AddApp(c.AppID, c.Version)
	app.MachineID = c.ID
	app.BootId = c.SessionID
	app.Track = c.Track
	app.OEM = Flags.OEM

	if updateCheck {
		app.AddUpdateCheck()
	}

	if isPing {
		app.AddPing()
		app.Ping.LastReportDays = "1"
		app.Ping.Status = "1"
	}

	if otype != "" {
		event := app.AddEvent()
		event.Type = otype
		event.Result = result
		if result == "0" {
			event.ErrorCode = "2000"
		} else {
			event.ErrorCode = ""
		}
	}

	return req
}

func (c *Client) MakeRequest(otype, result string, updateCheck, isPing bool) (*omaha.Response, error) {
	client := &http.Client{}
	req := c.OmahaRequest(otype, result, updateCheck, isPing)
	raw, err := xml.MarshalIndent(req, "", " ")
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(c.config.server+"/v1/update/", "text/xml", bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	oresp := new(omaha.Response)
	err = xml.NewDecoder(resp.Body).Decode(oresp)
	if err != nil {
		return nil, err
	}

	if Flags.verbose {
		raw, _ := xml.MarshalIndent(req, "", " ")
		c.Logf("request: %s\n", string(raw))
		raw, _ = xml.MarshalIndent(oresp, "", " ")
		c.Logf("response: %s\n", string(raw))
	}

	return oresp, nil
}

// Loop between n and m seconds
func (c *Client) Loop(n, m int) {
	interval := constant.InitialInterval
	for {
		randSleep(n, m)
		resp, err := c.MakeRequest("3", "2", true, false)
		if err != nil {
			log.Println(err)
			continue
		}
		uc := resp.Apps[0].UpdateCheck
		if uc.Status != "ok" {
			c.Logf("update check status: %s\n", uc.Status)
		} else {
			url := c.getCodebaseURL(uc)
			if !strings.Contains(url, "deis") {
				c.failed("Wrong Url", err)
				continue
			}
			c.MakeRequest("13", "1", false, false)
			err = c.downloadFromURL(url, "/tmp/deis.tar.gz")
			if err != nil {
				c.failed("Download failed", err)
				continue
			}
			err = utils.Extract("/tmp/deis.tar.gz", "/")
			if err != nil {
				c.failed("Extract failed", err)
				continue
			}
			c.MakeRequest("14", "1", false, false)
			c.RequestLock()
			for {
				err = c.lock.Lock()
				if err != nil && err != lock.ErrExist {
					interval = expBackoff(interval)
					fmt.Printf("Retrying in %v. Error locking: %v\n", interval, err)
					time.Sleep(interval)
					continue
				} else {
					break
				}
			}
			c.SetVersion(resp)
			err = c.lock.Unlock()
			if err == lock.ErrNotExist {
				fmt.Println("no lock found")
			} else if err == nil {
				fmt.Println("Unlocked existing lock for this machine")
			} else {
				fmt.Fprintln(os.Stderr, "Error unlocking:", err)
			}
		}
	}
}

func (c *Client) SetVersion(resp *omaha.Response) {
	// A field can potentially be nil.
	defer func() {
		if err := recover(); err != nil {
			c.Logf("%s: error setting version: %v", c.ID, err)
		}
	}()
	uc := resp.Apps[0].UpdateCheck

	err := c.Update()
	if err != nil {
		c.failed("update failed", err)
		return
	}
	c.Logf("Installation done ")
	c.MakeRequest("2", "1", false, false)
	c.Logf("Update done ")
	c.MakeRequest("3", "1", false, false)
	// installed

	// simulate reboot lock for a while
	for c.pingsRemaining > 0 {
		c.MakeRequest("", "", false, true)
		c.pingsRemaining--
		time.Sleep(1 * time.Second)
	}

	c.Logf("updated from %s to %s\n", c.Version, uc.Manifest.Version)

	c.Version = uc.Manifest.Version
	utils.PutVersion(c.Version)

	_, err = c.MakeRequest("3", "2", false, false) // Send complete with new version.
	if err != nil {
		log.Println(err)
	}

	c.SessionID = uuid.New()
}

func (c *Client) Update() (err error) {
	deis, _ := fleet.NewClient()
	localServices := deis.GetLocaljobs()
	fmt.Printf("local services: %v\n", localServices)

	if localServices.Len() == 0 {
		fmt.Println("no local services")
	}

	for _, service := range localServices {
		if strings.HasSuffix(service, "-data.service") {
			continue
		}
		localService := strings.Split(strings.Split(service, "-")[1], ".service")[0]
		fmt.Printf("destroying %v\n", localService)
		err := deis.Destroy([]string{localService})
		if err != nil {
			return err
		}
		fmt.Printf("re-creating %v\n", localService)
		err = deis.Create([]string{localService})
		if err != nil {
			return err
		}
		fmt.Printf("starting %v\n", localService)
		err = deis.Start([]string{localService})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) downloadFromURL(url, filePath string) (err error) {
	fmt.Printf("Downloading %s to %s\n", url, filePath)

	output, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error while creating", filePath, "-", err)
		return
	}
	defer output.Close()
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
	return
}

func expBackoff(interval time.Duration) time.Duration {
	interval = interval * 2
	if interval > constant.MaxInterval {
		interval = constant.MaxInterval
	}
	return interval
}

// Sleeps randomly between n and m seconds.
func randSleep(n, m int) {
	r := m
	if m-n > 0 {
		r = rand.Intn(m-n) + n
	}
	time.Sleep(time.Duration(r) * time.Second)
}
