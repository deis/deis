package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/deis/deis/pkg/prettyprint"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/models/config"
)

// ConfigList lists an app's config.
func ConfigList(appID string, oneLine bool) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	config, err := config.List(c, appID)

	if err != nil {
		return err
	}

	var keys []string
	for k := range config.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if oneLine {
		for _, key := range keys {
			fmt.Printf("%s=%s ", key, config.Values[key])
		}
		fmt.Println()
	} else {
		fmt.Printf("=== %s Config\n", appID)

		configMap := make(map[string]string)

		// config.Values is type interface, so it needs to be converted to a string
		for _, key := range keys {
			configMap[key] = fmt.Sprintf("%v", config.Values[key])
		}

		fmt.Print(prettyprint.PrettyTabs(configMap, 6))
	}

	return nil
}

// ConfigSet sets an app's config variables.
func ConfigSet(appID string, configVars []string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	configMap := parseConfig(configVars)

	value, ok := configMap["SSH_KEY"]

	if ok {
		sshKey := value.(string)

		if _, err := os.Stat(value.(string)); err == nil {
			contents, err := ioutil.ReadFile(value.(string))

			if err != nil {
				return err
			}

			sshKey = string(contents)
		}

		sshRegex := regexp.MustCompile("^-.+ .SA PRIVATE KEY-*")

		if !sshRegex.MatchString(sshKey) {
			return fmt.Errorf("Could not parse SSH private key:\n %s", sshKey)
		}

		configMap["SSH_KEY"] = base64.StdEncoding.EncodeToString([]byte(sshKey))
	}

	fmt.Print("Creating config... ")

	quit := progress()
	configObj := api.Config{Values: configMap}
	configObj, err = config.Set(c, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	if release, ok := configObj.Values["DEIS_RELEASE"]; ok {
		fmt.Printf("done, %s\n\n", release)
	} else {
		fmt.Print("done\n\n")
	}

	return ConfigList(appID, false)
}

// ConfigUnset removes a config variable from an app.
func ConfigUnset(appID string, configVars []string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Print("Removing config... ")

	quit := progress()

	configObj := api.Config{}

	valuesMap := make(map[string]interface{})

	for _, configVar := range configVars {
		valuesMap[configVar] = nil
	}

	configObj.Values = valuesMap

	_, err = config.Set(c, appID, configObj)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Print("done\n\n")

	return ConfigList(appID, false)
}

// ConfigPull pulls an app's config to a file.
func ConfigPull(appID string, interactive bool, overwrite bool) error {
	filename := ".env"

	if !overwrite {
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("%s already exists, pass -o to overwrite", filename)
		}
	}

	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	configVars, err := config.List(c, appID)

	if interactive {
		contents, err := ioutil.ReadFile(filename)

		if err != nil {
			return err
		}
		localConfigVars := strings.Split(string(contents), "\n")

		configMap := parseConfig(localConfigVars[:len(localConfigVars)-1])

		for key, value := range configVars.Values {
			localValue, ok := configMap[key]

			if ok {
				if value != localValue {
					var confirm string
					fmt.Printf("%s: overwrite %s with %s? (y/N) ", key, localValue, value)

					fmt.Scanln(&confirm)

					if strings.ToLower(confirm) == "y" {
						configMap[key] = value
					}
				}
			} else {
				configMap[key] = value
			}
		}

		return ioutil.WriteFile(filename, []byte(formatConfig(configMap)), 0755)
	}

	return ioutil.WriteFile(filename, []byte(formatConfig(configVars.Values)), 0755)
}

// ConfigPush pushes an app's config from a file.
func ConfigPush(appID string, fileName string) error {
	contents, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	config := strings.Split(string(contents), "\n")
	return ConfigSet(appID, config[:len(config)-1])
}

func parseConfig(configVars []string) map[string]interface{} {
	configMap := make(map[string]interface{})

	regex := regexp.MustCompile(`^([A-z_]+[A-z0-9_]*)=([\s\S]+)$`)
	for _, config := range configVars {
		if regex.MatchString(config) {
			captures := regex.FindStringSubmatch(config)
			configMap[captures[1]] = captures[2]
		} else {
			fmt.Printf("'%s' does not match the pattern 'key=var', ex: MODE=test\n", config)
		}
	}

	return configMap
}

func formatConfig(configVars map[string]interface{}) string {
	var formattedConfig string

	for key, value := range configVars {
		formattedConfig += fmt.Sprintf("%s=%s\n", key, value)
	}

	return formattedConfig
}
