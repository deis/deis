package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/deis/deis/client-go/controller/api"
)

// Register creates a account on a Deis controller.
func Register(controllerURL url.URL, username string, password string, email string, sslVerify bool,
	loginAfter bool) error {
	client := CreateHTTPClient(sslVerify)

	user := api.AuthRegisterRequest{Username: username, Password: password, Email: email}
	body, err := json.Marshal(user)

	if err != nil {
		return err
	}

	controllerURL.Path = "/v1/auth/register/"

	headers := http.Header{}

	controllerClient, err := New()

	if err == nil {
		headers.Add("Authorization", "token "+controllerClient.Token)
	}

	headers.Add("Content-Type", "application/json")
	addUserAgent(&headers)

	res, err := rawRequest(client, "POST", controllerURL.String(), bytes.NewBuffer(body), headers, 201)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	fmt.Printf("Registered %s\n", username)

	// Remove the path of the URL.
	controllerURL.Path = ""
	if loginAfter {
		return Login(controllerURL, username, password, sslVerify)
	}

	return nil
}

// Login logs a user into a Deis controller.
func Login(controllerURL url.URL, username string, password string, sslVerify bool) error {
	client := CreateHTTPClient(sslVerify)

	user := api.AuthLoginRequest{Username: username, Password: password}
	body, err := json.Marshal(user)

	if err != nil {
		return err
	}

	controllerURL.Path = "/v1/auth/login/"

	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	addUserAgent(&headers)

	res, err := rawRequest(client, "POST", controllerURL.String(), bytes.NewBuffer(body), headers, 200)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)

	token := api.AuthLoginResponse{}

	if err = json.Unmarshal([]byte(resBody), &token); err != nil {
		return err
	}

	// Remove the path of the URL.
	controllerURL.Path = ""
	controllerClient := Client{Username: username, SSLVerify: sslVerify,
		ControllerURL: controllerURL, Token: token.Token}

	if err = controllerClient.Save(); err != nil {
		return err
	}

	fmt.Printf("Logged in as %s\n", username)
	return nil
}

// Logout from a Deis controller by deleting config file.
func Logout() error {
	if err := deleteSettings(); err != nil {
		return err
	}

	fmt.Println("Logged out")
	return nil
}

// Passwd changes a user's password.
func Passwd(username string, password string, newPassword string) error {
	client, err := New()

	if err != nil {
		return err
	}

	req := api.AuthPasswdRequest{Password: password, NewPassword: newPassword}

	if username != "" {
		req.Username = username
	}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	resBody, status, err := client.BasicRequest("POST", "/v1/auth/passwd/", body)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("Password change failed: %s", resBody)
	}

	fmt.Println("Password change succeeded.")
	return nil
}

// Cancel deletes a user's account.
func Cancel() error {
	client, err := New()

	if err != nil {
		return err
	}

	body, status, err := client.BasicRequest("DELETE", "/v1/auth/cancel/", nil)

	if status != 204 {
		return fmt.Errorf("Cancellation failed: %s", body)
	}

	if err = deleteSettings(); err != nil {
		return err
	}

	fmt.Println("Account cancelled")
	return nil
}

// Regenerate regenenerates a user's token.
func Regenerate(username string, all bool) error {
	client, err := New()

	if err != nil {
		return err
	}

	var body []byte

	if all == true {
		body, err = json.Marshal(api.AuthRegenerateRequest{All: all})
	} else if username != "" {
		body, err = json.Marshal(api.AuthRegenerateRequest{Name: username})
	} else {
		body = []byte{}
	}

	if err != nil {
		return err
	}

	resBody, status, err := client.BasicRequest("POST", "/v1/auth/tokens/", body)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("Token regeneration failed: %s", resBody)
	}

	// If the token regenerated is the current user's, update it.
	if username == "" && all == false {
		token := api.AuthRegenerateResponse{}
		if err = json.Unmarshal([]byte(resBody), &token); err != nil {
			return nil
		}

		client.Token = token.Token

		if err = client.Save(); err != nil {
			return nil
		}
	}

	fmt.Println("Token Regenerated")
	return nil
}
