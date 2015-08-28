package cmd

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/deis/deis/client/controller/models/builds"
)

// BuildsList lists an app's builds.
func BuildsList(appID string, results int) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = c.ResponseLimit
	}

	builds, count, err := builds.List(c, appID, results)

	if err != nil {
		return err
	}

	fmt.Printf("=== %s Builds%s", appID, limitCount(len(builds), count))

	for _, build := range builds {
		fmt.Println(build.UUID, build.Created)
	}
	return nil
}

// BuildsCreate creates a build for an app.
func BuildsCreate(appID, image, procfile string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	var procfileMap map[string]string

	if procfile != "" {
		err = yaml.Unmarshal([]byte(procfile), &procfileMap)

		if err != nil {
			return err
		}
	}

	fmt.Print("Creating build... ")
	quit := progress()
	_, err = builds.New(c, appID, image, procfileMap)
	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Println("done")

	return nil
}
