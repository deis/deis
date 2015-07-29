package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

// Client oversees the interaction between the client and controller
type Client struct {
	// HTTP client used to communicate with the API.
	HTTPClient *http.Client

	// SSLVerify determines whether or not to verify SSL conections.
	SSLVerify bool

	// URL used to communicate with the controller.
	ControllerURL url.URL

	// Token is used to authenticate the request against the API.
	Token string

	// Username is the name of the user performing requests against the API.
	Username string

	// ResponseLimit is the number of results to return on requests that can be limited.
	ResponseLimit int
}

// DefaultResponseLimit is the default number of responses to return on requests that can
// be limited.
var DefaultResponseLimit = 100

type settingsFile struct {
	Username   string `json:"username"`
	SslVerify  bool   `json:"ssl_verify"`
	Controller string `json:"controller"`
	Token      string `json:"token"`
	Limit      int    `json:"response_limit"`
}

// New creates a new client from a settings file.
func New() (*Client, error) {
	filename := locateSettingsFile()

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("Not logged in. Use 'deis login' or 'deis register' to get started.")
		}

		return nil, err
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	settings := settingsFile{}
	if err = json.Unmarshal(contents, &settings); err != nil {
		return nil, err
	}

	u, err := url.Parse(settings.Controller)
	if err != nil {
		return nil, err
	}

	if settings.Limit <= 0 {
		settings.Limit = DefaultResponseLimit
	}

	return &Client{HTTPClient: CreateHTTPClient(settings.SslVerify), SSLVerify: settings.SslVerify,
		ControllerURL: *u, Token: settings.Token, Username: settings.Username,
		ResponseLimit: settings.Limit}, nil
}

// Save settings to a file
func (c Client) Save() error {
	settings := settingsFile{Username: c.Username, SslVerify: c.SSLVerify,
		Controller: c.ControllerURL.String(), Token: c.Token, Limit: c.ResponseLimit}

	if settings.Limit <= 0 {
		settings.Limit = DefaultResponseLimit
	}

	settingsContents, err := json.Marshal(settings)

	if err != nil {
		return err
	}

	if err = os.MkdirAll(path.Join(FindHome(), "/.deis/"), 0775); err != nil {
		return err
	}

	return ioutil.WriteFile(locateSettingsFile(), settingsContents, 0775)
}

// Delete user's settings file.
func Delete() error {
	filename := locateSettingsFile()

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}
