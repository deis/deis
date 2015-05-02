package wizard

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/utils"
)

// Run executes a wizard to interactively collect install options and
// configuration..
func Run() error {

	fmt.Println(utils.DeisIfy("Deis Install Wizard"))

	err := collectDomainInfo()
	if err != nil {
		return err
	}

	incStore, err := collectStorageSubsystemOptions()
	if err != nil {
		return err
	}

	err = collectDbOptions(incStore)
	if err != nil {
		return err
	}

	err = collectObjectStoreOptions(incStore)
	if err != nil {
		return err
	}

	err = collectLogSubsystemOptions(incStore)
	if err != nil {
		return err
	}

	fmt.Println("\nThat's all the information needed.\n")
	return nil

}

func collectDomainInfo() error {
	domain := promptForResponse("Enter the domain on which you wish to host Deis: ", false)
	err := config.SetConfig("platform", "domain", domain)
	if err != nil {
		return err
	}
	return nil
}

func collectStorageSubsystemOptions() (bool, error) {
	// Find out if the operator wants to install the storage subsystem
	fmt.Println("\nDeis' storage subsystem is used for storing platform state, but it is optional.")
	fmt.Println("It is possible to, instead, bind Deis components to externally-managed services for storing platform state.")
	incStore := promptForYesNo("Do you want to install the storage subsystem?")
	if !incStore {
		// If they do not, make doubly sure they know what they are getting
		// themselves into
		fmt.Println("\nThis will require you to manage your own external PostgreSQL server, S3 bucket or other compatible object store, and log server.")
		fmt.Println("This is an advanced option.")
		incStore = !promptForYesNo("Is this really what you want?")
	}
	err := config.SetConfig("installer", "includeStorageSubsystem", strconv.FormatBool(incStore))
	if err != nil {
		return incStore, err
	}
	return incStore, nil
}

func collectDbOptions(incStore bool) error {
	var incDb bool
	// The option to install the deis-database component is only applicable if
	// the storage subsystem is being installed
	if incStore {
		fmt.Println("\nNumerous Deis components store platform state in a PostgreSQL database.")
		fmt.Println("You have opted to install Deis' storage subsystem.")
		fmt.Println("It is possible for Deis to manage a PostgreSQL server on top of that.")
		fmt.Println("Alternatively, Deis components can bind to an externally managed PostgreSQL server.")
		// Find out if the operator wishes to install the deis-database component
		incDb = promptForYesNo("Do you want Deis to install its own PostgreSQL server?")
		if !incDb {
			// If they do not, make doubly sure they know what they are getting
			// themselves into
			fmt.Println("\nThis will require you to manage your own external PostgreSQL server.")
			fmt.Println("This is an advanced option.")
			incDb = !promptForYesNo("Is this really what you want?")
		}
	} else {
		// Storage subsystem is not being installed, so we definitely cannot
		// install deis-database
		incDb = false
	}

	err := config.SetConfig("installer", "includeDatabase", strconv.FormatBool(incDb))
	if err != nil {
		return err
	}

	// If not installing the deis-database component, collect connection details
	// for an exteral PostgreSQL server
	if !incDb {
		if incStore {
			// Operator opted into the storage subsystem, but opted out of the
			// deis-database component
			fmt.Println("\nYou have opted not to install Deis' own, internal PostgreSQL server.")
		} else {
			// Operator implicitly opted out of the deis-database component by opting
			// out of storage subsystem
			fmt.Println("\nNumerous Deis components store platform state in a PostgreSQL database.")
			fmt.Println("You have opted not to install Deis' storage subsystem.")
		}
		fmt.Println("This requires you to manage your own external PostgreSQL server.")
		fmt.Println("Let's collect database connection details.")
		err = config.SetConfig("database", "engine", "postgresql_psycopg2")
		if err != nil {
			return err
		}
		dbHost := promptForResponse("Enter the PostgreSQL server's host address: ", false)
		err = config.SetConfig("database", "host", dbHost)
		if err != nil {
			return err
		}
		dbPort := promptForResponse("Enter the port number that PostgreSQL server listens on (usually 5432): ", false)
		err = config.SetConfig("database", "port", dbPort)
		if err != nil {
			return err
		}
		dbName := promptForResponse("Enter the database name: ", false)
		err = config.SetConfig("database", "name", dbName)
		if err != nil {
			return err
		}
		dbUser := promptForResponse("Enter the database username: ", false)
		err = config.SetConfig("database", "user", dbUser)
		if err != nil {
			return err
		}
		dbPass := promptForResponse("Enter the database user's password: ", false)
		err = config.SetConfig("database", "password", dbPass)
		if err != nil {
			return err
		}
	}
	return nil
}

func collectObjectStoreOptions(incStore bool) error {
	// If not installing the storage subsystem, collect connection details for an
	// S3 bucket or other external, S3-compatible object store
	if !incStore {
		fmt.Println("\nDeis' internal Docker registry must be backed by an S3 bucket or other compatible object store.")
		fmt.Println("You have opted not to install Deis' storage subsystem.")
		fmt.Println("This requires you to manage your own S3 bucket or other compatible object store.")
		fmt.Println("Let's collect connection details.")
		objStoreRegistryRegion := promptForResponse("Enter the region if using S3, otherwise leave blank): ", true)
		err := config.SetConfig("registry", "s3region", objStoreRegistryRegion)
		if err != nil {
			return err
		}
		objStoreRegistryHost := promptForResponse("Enter host for the object store gateway (s3.amazonaws.com if using S3): ", false)
		err = config.SetConfig("store", "gateway/host", objStoreRegistryHost)
		if err != nil {
			return err
		}
		objStoreRegistryPort := promptForResponse("Enter port number for the object store gateway (usually 80 or 443): ", false)
		err = config.SetConfig("store", "gateway/port", objStoreRegistryPort)
		if err != nil {
			return err
		}
		objStoreRegistryAccessKey := promptForResponse("Enter the access key: ", false)
		err = config.SetConfig("registry", "s3accessKey", objStoreRegistryAccessKey)
		if err != nil {
			return err
		}
		err = config.SetConfig("store", "gateway/accessKey", objStoreRegistryAccessKey)
		if err != nil {
			return err
		}
		objStoreRegistrySecretKey := promptForResponse("Enter the secret key: ", false)
		err = config.SetConfig("registry", "s3secretKey", objStoreRegistrySecretKey)
		if err != nil {
			return err
		}
		err = config.SetConfig("store", "gateway/secretKey", objStoreRegistrySecretKey)
		if err != nil {
			return err
		}
		objStoreRegistryBucket := promptForResponse("Enter the bucket name: ", false)
		err = config.SetConfig("registry", "s3bucket", objStoreRegistryBucket)
		if err != nil {
			return err
		}
		err = config.SetConfig("registry", "bucketName", objStoreRegistryBucket)
		if err != nil {
			return err
		}
		objStoreRegistryPath := promptForResponse("Enter the path to where Docker images may be stored (/registry is recommended): ", false)
		err = config.SetConfig("registry", "s3path", objStoreRegistryPath)
		if err != nil {
			return err
		}
		useSecure := promptForYesNo("Should HTTPS be used when connecting?")
		err = config.SetConfig("registry", "s3secure", strconv.FormatBool(useSecure))
		if err != nil {
			return err
		}
		useEncryption := promptForYesNo("Should server-size encryption be used?")
		err = config.SetConfig("registry", "s3encrypt", strconv.FormatBool(useEncryption))
		if err != nil {
			return err
		}
	}
	return nil
}

func collectLogSubsystemOptions(incStore bool) error {
	var drainLogs bool
	if incStore {
		// Operator opted into the storage subsystem, deis-logger and
		// deis-controller will, jointly, do a fine job of aggregating logs and
		// providing developers with access to those logs. Draining to an external
		// log server is therefore optional and is achieved by setting a value
		// for /deis/logs/drain in etcd.
		drainLogs = promptForYesNo("\nDo you wish to drain application logs to an external log server?")
		if drainLogs {
			logDrainUrl := promptForResponse("Enter URL for external log service, including protocol, host, and port number (e.g. syslog://log.example.com:514): ", false)
			err := config.SetConfig("logs", "drain", logDrainUrl)
			if err != nil {
				return err
			}
		}
	} else {
		// Operator opted out of the storage subsystem. This requires logs to be
		// drained to an external log server. Since there is no point in even
		// running the deis-logger component without the storage subsystem, this
		// will be achieved, instead, by configuring deis-logspout to forward logs
		// directly to the external log server (instead of forwarding them first to
		// deis-logger).
		fmt.Println("\nYou have opted not to install Deis' storage subsystem.")
		fmt.Println("Without this, logs MUST be forwarded to an external log server.")
		logHost := promptForResponse("Enter the log server's host address: ", false)
		err := config.SetConfig("logs", "host", logHost)
		if err != nil {
			return err
		}
		logPort := promptForResponse("Enter the port number that the log server server listens on: ", false)
		err = config.SetConfig("logs", "port", logPort)
		if err != nil {
			return err
		}
	}
	return nil
}

func promptForResponse(prompt string, allowEmpty bool) string {
	fmt.Print(prompt)
	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(response)
	if !allowEmpty && response == "" {
		fmt.Println("Invalid response")
		return promptForResponse(prompt, allowEmpty)
	}
	return response
}

func promptForYesNo(prompt string) bool {
	fmt.Printf("%s [Yn]: ", prompt)
	var response string
	fmt.Scanln(&response)
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Invalid response")
		return promptForYesNo(prompt)
	}
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}
