package updatectl

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"encoding/xml"
	"fmt"
	"github.com/coreos/go-omaha/omaha"
	update "github.com/coreos/updatectl/client/update/v1"
	"github.com/deis/deisctl/client"
	"github.com/deis/deisctl/cmd"
	"github.com/deis/deisctl/utils"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	initialInterval = time.Second * 10
	maxInterval     = time.Minute * 7
	downloadDir     = "/home/core/deis/systemd/"
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
	cmdInstanceDeis.Flags.StringVar(&instanceFlags.appId, "app-id", os.Getenv("DEISCTL_APP_ID"), "Application ID to update.")
	//instanceFlags.appId.required = true
	cmdInstanceDeis.Flags.StringVar(&instanceFlags.groupId, "group-id", os.Getenv("DEISCTL_GROUP_ID"), "Group ID to update.")
	//instanceFlags.groupId.required = true
	cmdInstanceDeis.Flags.StringVar(&instanceFlags.version, "version", os.Getenv("DEISCTL_APP_VERSION"), "Version to report.")
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

func (c *Client) getCodebaseUrl(uc *omaha.UpdateCheck) string {
	return uc.Urls.Urls[0].CodeBase
}

func (c *Client) updateservice() {
	fmt.Println("starting systemd units")
	files, _ := utils.ListFiles(downloadDir + "*.service")
	fmt.Println(files)
	deis, _ := client.NewClient()
	localServices := deis.GetLocaljobs()
	Services := utils.GetServices()
	if localServices.Len() == 0 {
		fmt.Println("no local services")
		return
	}
	for _, service := range localServices {
		cmd.Uninstall(deis, []string{strings.Split(strings.Split(service, "-")[1], ".")[0]})
		cmd.Install(deis, []string{strings.Split(strings.Split(service, "-")[1], ".")[0]})
	}
	var count int
	for _, service := range Services {
		count = 0
		for _, lserv := range localServices {
			if strings.Contains(lserv, strings.Split(strings.Split(service, "-")[1], ".")[0]) {
				count = count + 1
			}
		}
		if count == 0 {
			go func() {
				_ = cmd.PullImage(service)
			}()
		}
	}

	// pre-install hook (download all new docker images)
	// use systemd to list local deis-* units, ignore -data units
	// cmd.Unistall([]string{"router.3"})
	// cmd.Install([]string{"router.3"})
	// post-install hook (make sure upgrade was successful)
}

func (c *Client) downloadFromUrl(url, fileName string) (err error) {
	url = url + "deis.tar.gz"
	fmt.Printf("Downloading %s to %s", url, fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(downloadDir + fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
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
	url := c.getCodebaseUrl(uc)
	c.MakeRequest("13", "1", false, false)
	c.downloadFromUrl(url, "deis.tar.gz")
	utils.Extract(downloadDir+"deis.tar.gz", downloadDir)
	c.MakeRequest("14", "1", false, false)
	c.updateservice()
	fmt.Println("updated done")
	c.MakeRequest("3", "1", false, false)
	// installed
	fmt.Println("updated done")
	// simulate reboot lock for a while
	for c.pingsRemaining > 0 {
		c.MakeRequest("", "", false, true)
		c.pingsRemaining--
		time.Sleep(1 * time.Second)
	}

	c.Log("updated from %s to %s\n", c.Version, uc.Manifest.Version)

	c.Version = uc.Manifest.Version

	_, err := c.MakeRequest("3", "2", false, false) // Send complete with new version.
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
		Id:             fmt.Sprintf("{update-client-" + utils.NewUuid()),
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
