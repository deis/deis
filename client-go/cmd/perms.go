package cmd

import (
	"fmt"

	"github.com/deis/deis/client-go/controller/client"
	"github.com/deis/deis/client-go/controller/models/perms"
)

// PermsList prints which users have permissions.
func PermsList(appID string, admin bool) error {
	c, appID, err := permsLoad(appID, admin)

	if err != nil {
		return err
	}

	var users []string

	if admin {
		users, err = perms.ListAdmins(c)
	} else {
		users, err = perms.List(c, appID)
	}

	if err != nil {
		return err
	}

	if admin {
		fmt.Println("=== Administrators")
	} else {
		fmt.Printf("=== %s's Users\n", appID)
	}

	for _, user := range users {
		fmt.Println(user)
	}

	return nil
}

// PermCreate adds a user to an app or makes them an administrator.
func PermCreate(appID string, username string, admin bool) error {

	c, appID, err := permsLoad(appID, admin)

	if err != nil {
		return err
	}

	if admin {
		fmt.Printf("Adding %s to system administrators... ", username)
		perms.NewAdmin(c, username)
	} else {
		fmt.Printf("Adding %s to %s collaborators... ", username, appID)
		perms.New(c, appID, username)
	}

	if err != nil {
		return err
	}

	fmt.Println("done")

	return nil
}

// PermDelete removes a user from an app or revokes admin privilages.
func PermDelete(appID string, username string, admin bool) error {

	c, appID, err := permsLoad(appID, admin)

	if err != nil {
		return err
	}

	if admin {
		fmt.Printf("Removing %s from system administrators... ", username)
		perms.DeleteAdmin(c, username)
	} else {
		fmt.Printf("Removing %s from %s collaborators... ", username, appID)
		perms.Delete(c, appID, username)
	}

	if err != nil {
		return err
	}

	fmt.Println("done")

	return nil
}

func permsLoad(appID string, admin bool) (*client.Client, string, error) {
	c, err := client.New()

	if err != nil {
		return nil, "", err
	}

	if !admin && appID == "" {
		appID, err = c.DetectApp()

		if err != nil {
			return nil, "", err
		}
	}

	return c, appID, err
}
