package update

import (
	"fmt"
	"net/http"
	"strconv"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/updateservicectl/auth"
	update "github.com/coreos/updateservicectl/client/update/v1"
	"github.com/deis/deis/deisctl/constant"
	"github.com/deis/deis/deisctl/utils"
	docopt "github.com/docopt/docopt-go"
)

const (
	// DefaultOmahaServer to communicate with
	DefaultOmahaServer = "https://opdemand.update.core-os.net"
	// DefaultOEM string to report to Omaha Server
	DefaultOEM = "deisctl"
	// DefaultAppID used for Omaha protocol
	DefaultAppID = "0ccac0df-ca24-4f2b-bb7b-4a265bd0eb33"
	// DefaultGroupID used for Omaha protocol
	DefaultGroupID = "2e87b742-68c9-4d08-8f37-5cb7bb2c9d3a"
)

// Flags for update package
var Flags struct {
	Server        string
	groupID       string
	appID         string
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

func parseInt(arg string) (i int, err error) {
	i, err = strconv.Atoi(arg)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func setUpdateFlags(args map[string]interface{}) error {

	appID := utils.GetKey(constant.UpdatekeyDir, "app-id", "DEISCTL_APP_ID")
	if args["--app-id"] != nil {
		Flags.appID = args["--app-id"].(string)
	} else if appID != "" {
		Flags.appID = appID
	} else {
		Flags.appID = DefaultAppID
	}

	groupID := utils.GetKey(constant.UpdatekeyDir, "group-id", "DEISCTL_GROUP_ID")
	if args["--group-id"] != nil {
		Flags.groupID = args["--group-id"].(string)
	} else if groupID != "" {
		Flags.groupID = groupID
	} else {
		Flags.groupID = DefaultGroupID
	}

	// read version from /etc/deis-version
	if args["--version"] == nil {
		Flags.version = utils.GetVersion()
	} else {
		Flags.version = args["--version"].(string)
	}

	// read update server
	if args["--server"] == nil {
		Flags.Server = DefaultOmahaServer
	} else {
		Flags.Server = args["--server"].(string)
	}

	minSleep, err := parseInt(args["--min-sleep"].(string))
	if err != nil {
		return err
	}
	Flags.minSleep = minSleep

	maxSleep, err := parseInt(args["--max-sleep"].(string))
	if err != nil {
		return err
	}
	Flags.maxSleep = maxSleep

	Flags.verbose = args["--verbose"].(bool)
	Flags.OEM = DefaultOEM

	return nil
}

// Update runs the Deis update engine daemon
func Update() error {
	usage := `Deis Update Daemon

	Usage:
	deisctl update [options]

	Options:
	--verbose                   print out the request bodies [default: false]
	--min-sleep=<sec>           minimum time between update checks [default: 10]
	--max-sleep=<sec>           maximum time between update checks [default: 30]
	--server=<server>           alternate update server URL (optional)
	`
	// parse command-line arguments
	args, err := docopt.Parse(usage, nil, true, "", true)
	if err != nil {
		return err
	}
	fmt.Printf("args: %v\n", args)
	err = setUpdateFlags(args)
	if err != nil {
		return err
	}
	fmt.Printf("flags: %v\n", Flags)
	return doUpdate()
}

type serverConfig struct {
	server string
}

func doUpdate() error {

	// construct hawk http client
	user, key, skipSSLVerify := "", "", true
	client := getHawkClient(user, key, skipSSLVerify)

	// use http client to construct update service
	service, err := update.New(client)
	if err != nil {
		return err
	}
	service.BasePath = Flags.Server + "/_ah/api/update/v1/"

	// create update client
	conf := &serverConfig{
		server: Flags.Server,
	}
	c := &Client{
		ID:        utils.GetClientID(),
		SessionID: uuid.New(),
		Version:   Flags.version,
		AppID:     Flags.appID,
		Track:     Flags.groupID,
		config:    conf,
	}
	go c.Loop(Flags.minSleep, Flags.maxSleep)

	// run forever
	wait := make(chan bool)
	<-wait
	return nil
}

func getHawkClient(user string, key string, skipSSLVerify bool) *http.Client {
	return &http.Client{
		Transport: &auth.HawkRoundTripper{
			User:          user,
			Token:         key,
			SkipSSLVerify: skipSSLVerify,
		},
	}
}
