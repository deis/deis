package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/deis/deis/pkg/prettyprint"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/client/controller/models/apps"
	"github.com/deis/deis/client/controller/models/config"
	"github.com/deis/deis/client/pkg/git"
	"github.com/deis/deis/client/pkg/webbrowser"
)

// AppCreate creates an app.
func AppCreate(id string, buildpack string, remote string, noRemote bool) error {
	c, err := client.New()
	if err != nil {
		return err
	}

	fmt.Print("Creating Application... ")
	quit := progress()
	app, err := apps.New(c, id)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Printf("done, created %s\n", app.ID)

	if buildpack != "" {
		configValues := api.Config{
			Values: map[string]interface{}{
				"BUILDPACK_URL": buildpack,
			},
		}
		if _, err = config.Set(c, app.ID, configValues); err != nil {
			return err
		}
	}

	if !noRemote {
		if err = git.CreateRemote(c.ControllerURL.Host, remote, app.ID); err != nil {
			if err.Error() == "exit status 128" {
				fmt.Println("To replace the existing git remote entry, run:")
				fmt.Printf("  git remote rename deis deis.old && deis git:remote -a %s\n", app.ID)
			}
			return err
		}
	}

	fmt.Println("remote available at", git.RemoteURL(c.ControllerURL.Host, app.ID))

	return nil
}

// AppsList lists apps on the Deis controller.
func AppsList(results int) error {
	c, err := client.New()

	if err != nil {
		return err
	}

	if results == defaultLimit {
		results = c.ResponseLimit
	}

	apps, count, err := apps.List(c, results)

	if err != nil {
		return err
	}

	fmt.Printf("=== Apps%s", limitCount(len(apps), count))

	for _, app := range apps {
		fmt.Println(app.ID)
	}
	return nil
}

// AppInfo prints info about app.
func AppInfo(appID string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	app, err := apps.Get(c, appID)

	if err != nil {
		return err
	}

	fmt.Printf("=== %s Application\n", app.ID)
	fmt.Println("updated: ", app.Updated)
	fmt.Println("uuid:    ", app.UUID)
	fmt.Println("created: ", app.Created)
	fmt.Println("url:     ", app.URL)
	fmt.Println("owner:   ", app.Owner)
	fmt.Println("id:      ", app.ID)

	fmt.Println()
	// print the app processes
	if err = PsList(app.ID, defaultLimit); err != nil {
		return err
	}

	fmt.Println()
	// print the app domains
	if err = DomainsList(app.ID, defaultLimit); err != nil {
		return err
	}

	fmt.Println()

	return nil
}

// AppOpen opens an app in the default webbrowser.
func AppOpen(appID string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	app, err := apps.Get(c, appID)

	if err != nil {
		return err
	}

	u := app.URL
	if !(strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
		u = "http://" + u
	}

	return webbrowser.Webbrowser(u)
}

// AppLogs returns the logs from an app.
func AppLogs(appID string, lines int) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	logs, err := apps.Logs(c, appID, lines)

	if err != nil {
		return err
	}

	return printLogs(logs)
}

// printLogs prints each log line with a color matched to its category.
func printLogs(logs string) error {
	for _, log := range strings.Split(logs, "\n") {
		category := "unknown"
		parts := strings.Split(strings.Split(log, ": ")[0], " ")
		if len(parts) >= 2 {
			category = parts[1]
		}
		colorVars := map[string]string{
			"Color": chooseColor(category),
			"Log":   log,
		}
		fmt.Println(prettyprint.ColorizeVars("{{.V.Color}}{{.V.Log}}{{.C.Default}}", colorVars))
	}

	return nil
}

// AppRun runs a one time command in the app.
func AppRun(appID, command string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Running '%s'...\n", command)

	out, err := apps.Run(c, appID, command)

	if err != nil {
		return err
	}

	fmt.Print(out.Output)
	os.Exit(out.ReturnCode)
	return nil
}

// AppDestroy destroys an app.
func AppDestroy(appID, confirm string) error {
	gitSession := false

	c, err := client.New()

	if err != nil {
		return err
	}

	if appID == "" {
		appID, err = git.DetectAppName(c.ControllerURL.Host)

		if err != nil {
			return err
		}

		gitSession = true
	}

	if confirm == "" {
		fmt.Printf(` !    WARNING: Potentially Destructive Action
 !    This command will destroy the application: %s
 !    To proceed, type "%s" or re-run this command with --confirm=%s

> `, appID, appID, appID)

		fmt.Scanln(&confirm)
	}

	if confirm != appID {
		return fmt.Errorf("App %s does not match confirm %s, aborting.", appID, confirm)
	}

	startTime := time.Now()
	fmt.Printf("Destroying %s...\n", appID)

	if err = apps.Delete(c, appID); err != nil {
		return err
	}

	fmt.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	if gitSession {
		return git.DeleteRemote(appID)
	}

	return nil
}

// AppTransfer transfers app ownership to another user.
func AppTransfer(appID, username string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	fmt.Printf("Transferring %s to %s... ", appID, username)

	err = apps.Transfer(c, appID, username)

	if err != nil {
		return err
	}

	fmt.Println("done")

	return nil
}
