package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/deis/deis/pkg/prettyprint"

	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/certs"
)

// CertsList lists certs registered with the controller.
func CertsList(results int) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = c.ResponseLimit
	}

	certList, _, err := certs.List(c, results)

	if err != nil {
		return err
	}

	if len(certList) == 0 {
		fmt.Println("No certs")
		return nil
	}

	certMap := make(map[string]string)
	nameMax := 0
	expiresMax := 0
	for _, cert := range certList {
		certMap[cert.Name] = cert.Expires

		if len(cert.Name) > nameMax {
			nameMax = len(cert.Name)
		}
		if len(cert.Expires) > nameMax {
			expiresMax = len(cert.Expires)
		}
	}

	nameHeader := "Common Name"
	expiresHeader := "Expires"
	tabSpaces := 5
	bufferSpaces := tabSpaces

	if nameMax < len(nameHeader) {
		tabSpaces += len(nameHeader) - nameMax
		nameMax = len(nameHeader)
	} else {
		bufferSpaces += nameMax - len(nameHeader)
	}

	if expiresMax < len(expiresHeader) {
		expiresMax = len(expiresHeader)
	}

	fmt.Printf("%s%s%s\n", nameHeader, strings.Repeat(" ", bufferSpaces), expiresHeader)
	fmt.Printf("%s%s%s\n", strings.Repeat("-", nameMax), strings.Repeat(" ", 5),
		strings.Repeat("-", expiresMax))
	fmt.Print(prettyprint.PrettyTabs(certMap, tabSpaces))
	return nil
}

// CertAdd adds a cert to the controller.
func CertAdd(cert, key, commonName, sans string) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	fmt.Print("Adding SSL endpoint... ")
	quit := progress()
	err = processCertsAdd(c, cert, key, commonName, sans)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}

func processCertsAdd(c *client.Client, cert, key, commonName, sans string) error {
	if sans != "" {
		for _, san := range strings.Split(sans, ",") {
			if err := doCertAdd(c, cert, key, san); err != nil {
				return err
			}
		}
		return nil
	}

	return doCertAdd(c, cert, key, commonName)
}

func doCertAdd(c *client.Client, cert string, key string, commonName string) error {
	certFile, err := ioutil.ReadFile(cert)

	if err != nil {
		return err
	}

	keyFile, err := ioutil.ReadFile(key)

	if err != nil {
		return err
	}

	_, err = certs.New(c, string(certFile), string(keyFile), commonName)
	return err
}

// CertRemove deletes a cert from the controller.
func CertRemove(commonName string) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	fmt.Printf("Removing %s... ", commonName)
	quit := progress()

	certs.Delete(c, commonName)

	quit <- true
	<-quit

	if err == nil {
		fmt.Println("done")
	}

	return err
}
