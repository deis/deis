// +build integration

package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	gitCloneCmd  = "if [ ! -d {{.ExampleApp}} ] ; then git clone https://github.com/deis/{{.ExampleApp}}.git ; fi"
	gitRemoveCmd = "git remote remove deis"
	gitPushCmd   = "git push deis master"
)

// Client represents the client data structure in ~/.deis/client.json
type Client struct {
	Controller string `json:"controller"`
	Username   string `json:"username"`
	Token      string `json:"token"`
}

func TestGlobal(t *testing.T) {
	params := utils.GetGlobalConfig()
	utils.Execute(t, authRegisterCmd, params, false, "")
	clientTest(t, params)
	utils.Execute(t, keysAddCmd, params, false, "")
}

func clientTest(t *testing.T, params *utils.DeisTestConfig) {
	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	profile := os.Getenv("DEIS_PROFILE")
	if profile == "" {
		profile = "client"
	}
	clientJsonFilePath := ".deis/" + profile + ".json"
	data, err := ioutil.ReadFile(path.Join(user.HomeDir, clientJsonFilePath))
	if err != nil {
		t.Fatal(err)
	}
	client := &Client{}
	json.Unmarshal(data, &client)
	if client.Token == "" {
		t.Error("token not present in client.json")
	}
	if client.Controller == "" {
		t.Error("controller endpoint not present in client.json")
	}
	if client.Username == "" {
		t.Error("username not present in client.json")
	}
}
