package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/deis/deis/client-go/controller/client"
	"golang.org/x/crypto/ssh/terminal"
)

// Register creates a account on a Deis controller.
func Register(controller string, username string, password string, email string,
	sslVerify bool) error {

	u, err := url.Parse(controller)

	if err != nil {
		return err
	}

	controllerURL, err := chooseScheme(*u)

	if err != nil {
		return err
	}

	if err = client.CheckConection(client.CreateHTTPClient(sslVerify), controllerURL); err != nil {
		return err
	}

	if username == "" {
		fmt.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Printf("\npassword (confirm): ")
		passwordConfirm, err := readPassword()
		fmt.Println()

		if err != nil {
			return err
		}

		if password != passwordConfirm {
			return errors.New("Password mismatch, aborting registration.")
		}
	}

	if email == "" {
		fmt.Print("email: ")
		fmt.Scanln(&email)
	}

	return client.Register(controllerURL, username, password, email, sslVerify, true)
}

// Login to a Deis controller.
func Login(controller string, username string, password string, sslVerify bool) error {
	u, err := url.Parse(controller)

	if err != nil {
		return err
	}

	controllerURL, err := chooseScheme(*u)

	if err != nil {
		return err
	}

	if err = client.CheckConection(client.CreateHTTPClient(sslVerify), controllerURL); err != nil {
		return err
	}

	if username == "" {
		fmt.Print("username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("password: ")
		password, err = readPassword()
		fmt.Println()

		if err != nil {
			return err
		}
	}

	return client.Login(controllerURL, username, password, sslVerify)
}

// Logout from a Deis controller.
func Logout() error {
	return client.Logout()
}

// Passwd changes a user's password.
func Passwd(username string, password string, newPassword string) error {
	var err error

	if password == "" {
		fmt.Print("current password: ")
		password, err = readPassword()
		fmt.Println()

		if err != nil {
			return err
		}
	}

	if newPassword == "" {
		fmt.Print("new password: ")
		newPassword, err = readPassword()
		fmt.Printf("\nnew password (confirm): ")
		passwordConfirm, err := readPassword()

		fmt.Println()

		if err != nil {
			return err
		}

		if newPassword != passwordConfirm {
			return errors.New("Password mismatch, not changing.")
		}
	}

	return client.Passwd(username, password, newPassword)
}

// Cancel deletes a user's account.
func Cancel(username string, password string, yes bool) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	fmt.Println("Please log in again in order to cancel this account")

	if err = Login(c.ControllerURL.String(), username, password, c.SSLVerify); err != nil {
		return err
	}

	if yes == false {
		confirm := ""

		c, err = client.New()

		if err != nil {
			return err
		}

		fmt.Printf("cancel account %s at %s? (y/N): ", c.Username, c.ControllerURL.String())
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) == "y" {
			yes = true
		}
	}

	if yes == false {
		fmt.Println("Account not changed")
		return nil
	}

	return client.Cancel()
}

// Whoami prints the logged in user.
func Whoami() error {
	c, err := client.New()

	if err != nil {
		return err
	}

	fmt.Printf("You are %s at %s\n", c.Username, c.ControllerURL.String())
	return nil
}

// Regenerate regenenerates a user's token.
func Regenerate(username string, all bool) error {
	return client.Regenerate(username, all)
}

func readPassword() (string, error) {
	password, err := terminal.ReadPassword(0)

	return string(password), err
}

func chooseScheme(u url.URL) (url.URL, error) {
	if u.Scheme == "" {
		u.Scheme = "http"
		u, err := url.Parse(u.String())

		if err != nil {
			return url.URL{}, err
		}

		return *u, nil
	}

	return u, nil
}
