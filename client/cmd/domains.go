package cmd

import (
	"fmt"

	"github.com/deis/deis/client/controller/models/domains"
)

// DomainsList lists domains registered with an app.
func DomainsList(appID string, results int) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = c.ResponseLimit
	}

	domains, count, err := domains.List(c, appID, results)

	if err != nil {
		return err
	}

	fmt.Printf("=== %s Domains%s", appID, limitCount(len(domains), count))

	for _, domain := range domains {
		fmt.Println(domain.Domain)
	}
	return nil
}

// DomainsAdd adds a domain to an app.
func DomainsAdd(appID, domain string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Adding %s to %s... ", domain, appID)

	quit := progress()
	_, err = domains.New(c, appID, domain)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}

// DomainsRemove removes a domain registered with an app.
func DomainsRemove(appID, domain string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Removing %s from %s... ", domain, appID)

	quit := progress()
	err = domains.Delete(c, appID, domain)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Println("done")
	return nil
}
