package cmd

import (
	"fmt"

	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/users"
)

// UsersList lists users registered with the controller.
func UsersList(results int) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = c.ResponseLimit
	}

	users, count, err := users.List(c, results)

	if err != nil {
		return err
	}

	fmt.Printf("=== Users%s", limitCount(len(users), count))

	for _, user := range users {
		fmt.Println(user.Username)
	}
	return nil
}
