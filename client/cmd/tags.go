package cmd

import (
	"fmt"
	"strings"

	"github.com/deis/deis/pkg/prettyprint"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/models/config"
)

// TagsList lists an app's tags.
func TagsList(appID string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	config, err := config.List(c, appID)

	fmt.Printf("=== %s Tags\n", appID)

	tagMap := make(map[string]string)

	for key, value := range config.Tags {
		tagMap[key] = fmt.Sprintf("%v", value)
	}

	fmt.Print(prettyprint.PrettyTabs(tagMap, 5))

	return nil
}

// TagsSet sets an app's tags.
func TagsSet(appID string, tags []string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	tagsMap := parseTags(tags)

	fmt.Print("Applying tags... ")

	quit := progress()
	configObj := api.Config{}
	configObj.Tags = tagsMap

	_, err = config.Set(c, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")

	return TagsList(appID)
}

// TagsUnset removes an app's tags.
func TagsUnset(appID string, tags []string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Print("Applying tags... ")

	quit := progress()

	configObj := api.Config{}

	tagsMap := make(map[string]interface{})

	for _, tag := range tags {
		tagsMap[tag] = nil
	}

	configObj.Tags = tagsMap

	_, err = config.Set(c, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")

	return TagsList(appID)
}

func parseTags(tags []string) map[string]interface{} {
	tagMap := make(map[string]interface{})

	for _, tag := range tags {
		key, value, err := parseTag(tag)

		if err != nil {
			fmt.Println(err)
			continue
		}

		tagMap[key] = value
	}

	return tagMap
}

func parseTag(tag string) (string, string, error) {
	parts := strings.Split(tag, "=")

	if len(parts) != 2 {
		return "", "", fmt.Errorf(`%s is invalid, Must be in format key=value
Examples: rack=1 evironment=production`, tag)
	}

	return parts[0], parts[1], nil
}
