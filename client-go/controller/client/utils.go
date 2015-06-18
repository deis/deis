package client

import (
	"fmt"
	"os"
	"path"

	"github.com/deis/deis/version"
)

func locateSettingsFile() string {
	filename := os.Getenv("DEIS_PROFILE")

	if filename == "" {
		filename = "client"
	}

	return path.Join(os.Getenv("HOME"), ".deis", filename+".json")
}

func deleteSettings() error {
	filename := locateSettingsFile()

	_, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if err = os.Remove(filename); err != nil {
		return err
	}

	return nil
}

func checkAPICompatability(serverAPIVersion string) {
	if serverAPIVersion != version.APIVersion {
		fmt.Printf(`!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: %s
!    Server version: %s
`, version.APIVersion, serverAPIVersion)
	}
}
