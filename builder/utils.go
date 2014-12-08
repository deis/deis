package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// YamlToJson takes an input yaml string, parses it and returns a string formatted as json.
func YamlToJson(bytes []byte) (string, error) {
	var anomaly map[string]string

	if err := yaml.Unmarshal(bytes, &anomaly); err != nil {
		return "", err
	}

	retVal, err := json.Marshal(&anomaly)

	if err != nil {
		return "", err
	}

	return string(retVal), nil
}

// ParseConfig takes a response body from the controller and returns a Config object.
func ParseConfig(res *http.Response) (*Config, error) {
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(body, &config)
	return &config, err
}

func ParseDomain(bytes []byte) (string, error) {
	var hook BuildHookResponse
	if err := json.Unmarshal(bytes, &hook); err != nil {
		return "", err
	}

	if hook.Domains == nil {
		return "", fmt.Errorf("invalid application domain")
	}

	if len(hook.Domains) < 1 {
		return "", fmt.Errorf("invalid application domain")
	}

	return hook.Domains[0], nil
}

func ParseReleaseVersion(bytes []byte) (int, error) {
	var hook BuildHookResponse
	if err := json.Unmarshal(bytes, &hook); err != nil {
		return 0, fmt.Errorf("invalid application json configuration")
	}

	if hook.Release == nil {
		return 0, fmt.Errorf("invalid application version")
	}

	return hook.Release["version"], nil
}

func GetDefaultType(bytes []byte) (string, error) {
	type YamlTypeMap struct {
		DefaultProcessTypes ProcessType `default_process_types`
	}

	var p YamlTypeMap

	if err := yaml.Unmarshal(bytes, &p); err != nil {
		return "", err
	}

	retVal, err := json.Marshal(&p.DefaultProcessTypes)

	if err != nil {
		return "", err
	}

	if len(p.DefaultProcessTypes) == 0 {
		return "{}", nil
	}

	return string(retVal), nil
}

func ParseControllerConfig(bytes []byte) ([]string, error) {
	var controllerConfig Config
	if err := json.Unmarshal(bytes, &controllerConfig); err != nil {
		return []string{}, err
	}

	if controllerConfig.Values == nil {
		return []string{""}, nil
	}

	retVal := []string{}
	for k, v := range controllerConfig.Values {
		retVal = append(retVal, fmt.Sprintf(" -e %s=\"%v\"", k, v))
	}
	return retVal, nil
}
