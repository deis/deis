// Package net contains commonly useful network functions
package net

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Download downloads a resource from a specified
// source (URL) to the specified destination
func Download(src string, dest string) error {
	res, err := http.Get(src)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(dest, data, 0644); err != nil {
		return err
	}
	return nil
}
