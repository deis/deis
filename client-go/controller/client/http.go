package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/deis/deis/version"
)

// CreateHTTPClient creates a HTTP Client with proper SSL options.
func CreateHTTPClient(sslVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !sslVerify},
	}
	return &http.Client{Transport: tr}
}

func rawRequest(client *http.Client, method string, url string, body io.Reader, headers http.Header,
	expectedStatusCode int) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)
	req.Header = headers

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != expectedStatusCode {
		defer res.Body.Close()

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, checkForErrors(res, string(resBody))
	}

	return res, nil
}

// Request makes a HTTP request on the controller.
func (c Client) Request(method string, path string, body []byte) (*http.Response, error) {
	url := c.ControllerURL

	if strings.Contains(path, "?") {
		parts := strings.Split(path, "?")
		url.Path = parts[0]
		url.RawQuery = parts[1]
	} else {
		url.Path = path
	}

	req, err := http.NewRequest(method, url.String(), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.Token)
	addUserAgent(&req.Header)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	checkAPICompatability(res.Header.Get("DEIS_API_VERSION"))

	return res, nil
}

// BasicRequest makes a simple http request on the controller.
func (c Client) BasicRequest(method string, path string, body []byte) (string, error) {
	res, err := c.Request(method, path, body)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}
	return string(resBody), checkForErrors(res, string(resBody))
}

func checkForErrors(res *http.Response, body string) error {
	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil
	}

	bodyMap := make(map[string]interface{})

	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		return err
	}

	errorMessage := "\n"
	for key, value := range bodyMap {
		switch v := value.(type) {
		case string:
			errorMessage += key + ": " + v + "\n"
		case []interface{}:
			for _, subValue := range v {
				switch sv := subValue.(type) {
				case string:
					errorMessage += key + ": " + sv + "\n"
				default:
					fmt.Printf("Unexpected type in %s error message array. Contents: %v",
						reflect.TypeOf(value), sv)
				}
			}
		default:
			fmt.Printf("Cannot handle key %s in error message, type %s. Contents: %v",
				key, reflect.TypeOf(value), bodyMap[key])
		}
	}

	errorMessage += res.Status + "\n"
	return errors.New(errorMessage)
}

// CheckConection checks that the user is connected to a network and the URL points to a valid controller.
func CheckConection(client *http.Client, controllerURL url.URL) error {
	errorMessage := `%s does not appear to be a valid Deis controller.
Make sure that the Controller URI is correct and the server is running.`

	baseURL := controllerURL.String()

	controllerURL.Path = "/v1/"

	req, err := http.NewRequest("GET", controllerURL.String(), bytes.NewBuffer(nil))
	addUserAgent(&req.Header)

	if err != nil {
		return err
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf(errorMessage+"\n", baseURL)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 401 {
		return fmt.Errorf(errorMessage, baseURL)
	}

	checkAPICompatability(res.Header.Get("DEIS_API_VERSION"))

	return nil
}

func addUserAgent(headers *http.Header) {
	headers.Add("User-Agent", "Deis Client v"+version.Version)
}
