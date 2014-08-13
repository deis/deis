package updatectl

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
	"text/tabwriter"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/go-omaha/omaha"
	update "github.com/coreos/updatectl/client/update/v1"
	"github.com/deis/deisctl/client"
	"github.com/deis/deisctl/cmd"
	"github.com/deis/deisctl/constant"
	"github.com/deis/deisctl/utils"
)

var (
	instanceFlags struct {
		groupId       string
		appId         string
		start         int64
		end           int64
		verbose       bool
		clientsPerApp int
		minSleep      int
		maxSleep      int
		errorRate     int
		OEM           string
		pingOnly      int
		version       string
	}

	cmdInstance = &Command{
		Name:    "instance",
		Usage:   "[OPTION]...",
		Summary: "Operations to view instances.",
		Subcommands: []*Command{
			cmdInstanceDeis,
		},
	}

	cmdInstanceDeis = &Command{
		Name:        "instance deis",
		Usage:       "[OPTION]...",
		Description: "Simulate single deis to update instances.",
		Run:         instanceDeis,
	}
)

func init() {
	cmdInstanceDeis.Flags.BoolVar(&instanceFlags.verbose, "verbose", false, "Print out the request bodies")
	cmdInstanceDeis.Flags.IntVar(&instanceFlags.clientsPerApp, "clients-per-app", 1, "Number of fake fents per appid.")
	cmdInstanceDeis.Flags.IntVar(&instanceFlags.minSleep, "min-sleep", 5, "Minimum time between update checks.")
	cmdInstanceDeis.Flags.IntVar(&instanceFlags.maxSleep, "max-sleep", 10, "Maximum time between update checks.")
	cmdInstanceDeis.Flags.IntVar(&instanceFlags.errorRate, "errorrate", 1, "Chance of error (0-100)%.")
	cmdInstanceDeis.Flags.StringVar(&instanceFlags.OEM, "oem", "deisclient", "oem to report")
	// simulate reboot lock.
	cmdInstanceDeis.Flags.IntVar(&instanceFlags.pingOnly, "ping-only", 0, "halt update and just send ping requests this many times.")
	cmdInstanceDeis.Flags.StringVar(&instanceFlags.appId, "app-id", utils.GetKey(constant.UpdatekeyDir, "app-id", "DEISCTL_APP_ID"), "Application ID to update.")
	//instanceFlags.appId.required = true
	cmdInstanceDeis.Flags.StringVar(&instanceFlags.groupId, "group-id", utils.GetKey(constant.UpdatekeyDir, "group-id", "DEISCTL_GROUP_ID"), "Group ID to update.")
	//instanceFlags.groupId.required = true
	cmdInstanceDeis.Flags.StringVar(&instanceFlags.version, "version", utils.GetVersion(), "Version to report.")
}

//+ downloadDir + "deis.tar.gz"

type serverConfig struct {
	server string
}

type Client struct {
	Id             string
	SessionId      string
	Version        string
	AppId          string
	Track          string
	config         *serverConfig
	errorRate      int
	pingsRemaining int
}

func (c *Client) Log(format string, v ...interface{}) {
	format = c.Id + ": " + format
	fmt.Printf(format, v...)
}

func (c *Client) failed(tag string, err error) {
	c.Log("%s %v\n", tag, err)
	c.MakeRequest("3", "0", false, false)
}

func (c *Client) getCodebaseUrl(uc *omaha.UpdateCheck) string {
	return uc.Urls.Urls[0].CodeBase + uc.Manifest.Packages.Packages[0].Name
}

func (c *Client) updateservice() (err error) {
	fmt.Println("starting systemd units")
	// files, _ := utils.ListFiles(constant.UnitsDir + "*.service")
	deis, _ := client.NewClient()
	localServices := deis.GetLocaljobs()
	fmt.Printf("local services: %v\n", localServices)
	Services := utils.GetServices()
	if localServices.Len() == 0 {
		fmt.Println("no local services")
		return
	}
	for _, service := range localServices {
		if strings.HasSuffix(service, "-data.service") {
			continue
		}
		localService := strings.Split(strings.Split(service, "-")[1], ".service")[0]
		err := cmd.Uninstall(deis, []string{localService})
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		err = cmd.Install(deis, []string{localService})
		if err != nil {
			return err
		}
	}
	var count int
	for _, service := range Services {
		count = 0
		if strings.HasSuffix(service, "-data.service") {
			continue
		}
		for _, lserv := range localServices {
			if strings.Contains(lserv, strings.Split(strings.Split(service, "-")[1], ".")[0]) {
				count = count + 1
			}
		}
		if count == 0 {
			if err = cmd.PullImage(service); err != nil {
				fmt.Println("failed pulling image", service)
				return err
			}
		}
	}
	return nil
	// pre-install hook (download all new docker images)
	// use systemd to list local deis-* units, ignore -data units
	// cmd.Unistall([]string{"router.3"})
	// cmd.Install([]string{"router.3"})
	// post-install hook (make sure upgrade was successful)
}

func (c *Client) downloadFromUrl(url, filePath string) (err error) {
	fmt.Printf("Downloading %s to %s", url, filePath)

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

func (c *Client) OmahaRequest(otype, result string, updateCheck, isPing bool) *omaha.Request {
	req := omaha.NewRequest("lsb", "CoreOS", "", "")
	app := req.AddApp(c.AppId, c.Version)
	app.MachineID = c.Id
	app.BootId = c.SessionId
	app.Track = c.Track
	app.OEM = instanceFlags.OEM

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

	if instanceFlags.verbose {
		raw, _ := xml.MarshalIndent(req, "", " ")
		c.Log("request: %s\n", string(raw))
		raw, _ = xml.MarshalIndent(oresp, "", " ")
		c.Log("response: %s\n", string(raw))
	}

	return oresp, nil
}

func (c *Client) SetVersion(resp *omaha.Response) {
	// A field can potentially be nil.
	defer func() {
		if err := recover(); err != nil {
			c.Log("%s: error setting version: %v", c.Id, err)
		}
	}()
	uc := resp.Apps[0].UpdateCheck

	err := c.updateservice()
	if err != nil {
		c.failed("update failed", err)
		return
	}
	c.Log("Installation done ")
	c.MakeRequest("2", "1", false, false)
	c.Log("Update done ")
	c.MakeRequest("3", "1", false, false)
	// installed

	// simulate reboot lock for a while
	for c.pingsRemaining > 0 {
		c.MakeRequest("", "", false, true)
		c.pingsRemaining--
		time.Sleep(1 * time.Second)
	}

	c.Log("updated from %s to %s\n", c.Version, uc.Manifest.Version)

	c.Version = uc.Manifest.Version
	utils.PutVersion(c.Version)

	_, err = c.MakeRequest("3", "2", false, false) // Send complete with new version.
	if err != nil {
		log.Println(err)
	}

	c.SessionId = uuid.New()
}

// Sleep between n and m seconds
func (c *Client) Loop(n, m int) {
	for {
		randSleep(n, m)
		resp, err := c.MakeRequest("3", "2", true, false)
		if err != nil {
			log.Println(err)
			continue
		}
		uc := resp.Apps[0].UpdateCheck
		if uc.Status != "ok" {
			c.Log("update check status: %s\n", uc.Status)
		} else {
			url := c.getCodebaseUrl(uc)
			if !strings.Contains(url, "deis") {
				c.failed("Wrong Url", err)
				continue
			}
			c.MakeRequest("13", "1", false, false)
			err = c.downloadFromUrl(url, "/tmp/deis.tar.gz")
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
			c.SetVersion(resp)
		}
	}
}

// Sleeps randomly between n and m seconds.
func randSleep(n, m int) {
	r := m
	if m-n > 0 {
		r = rand.Intn(m-n) + n
	}
	time.Sleep(time.Duration(r) * time.Second)
}

func instanceDeis(args []string, service *update.Service, out *tabwriter.Writer) int {
	if instanceFlags.appId == "" || instanceFlags.groupId == "" {
		return ERROR_USAGE
	}

	conf := &serverConfig{
		server: globalFlags.Server,
	}

	c := &Client{
		Id:             utils.GetClientID(),
		SessionId:      uuid.New(),
		Version:        instanceFlags.version,
		AppId:          instanceFlags.appId,
		Track:          instanceFlags.groupId,
		config:         conf,
		errorRate:      instanceFlags.errorRate,
		pingsRemaining: instanceFlags.pingOnly,
	}
	go c.Loop(instanceFlags.minSleep, instanceFlags.maxSleep)

	// run forever
	wait := make(chan bool)
	<-wait
	return OK
}
