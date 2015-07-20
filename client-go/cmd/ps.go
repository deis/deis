package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/models/ps"
)

// PsList lists an app's processes.
func PsList(appID string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	processes, err := ps.List(c, appID)

	if err != nil {
		return err
	}

	printProcesses(appID, processes)

	return nil
}

// PsScale scales an app's processes.
func PsScale(appID string, targets []string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	targetMap := make(map[string]int)
	regex := regexp.MustCompile("^([A-z]+)=([0-9]+)$")

	for _, target := range targets {
		if regex.MatchString(target) {
			captures := regex.FindStringSubmatch(target)
			targetMap[captures[1]], err = strconv.Atoi(captures[2])

			if err != nil {
				return err
			}
		} else {
			fmt.Printf("'%s' does not match the pattern 'type=num', ex: web=2\n", target)
		}
	}

	fmt.Printf("Scaling processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress()

	err = ps.Scale(c, appID, targetMap)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	processes, err := ps.List(c, appID)

	if err != nil {
		return err
	}

	printProcesses(appID, processes)
	return nil
}

// PsRestart restarts an app's processes.
func PsRestart(appID, target string) error {
	c, appID, err := load(appID)

	if err != nil {
		return err
	}

	psType := ""
	psNum := -1

	if target != "" {
		if strings.Contains(target, ".") {
			parts := strings.Split(target, ".")
			psType = parts[0]
			psNum, err = strconv.Atoi(parts[1])

			if err != nil {
				return err
			}
		} else {
			psType = target
		}
	}

	fmt.Printf("Restarting processes... but first, %s!\n", drinkOfChoice())
	startTime := time.Now()
	quit := progress()

	_, err = ps.Restart(c, appID, psType, psNum)

	quit <- true
	<-quit

	if err != nil {
		return err
	}

	fmt.Printf("done in %ds\n", int(time.Since(startTime).Seconds()))

	processes, err := ps.List(c, appID)

	if err != nil {
		return err
	}

	printProcesses(appID, processes)
	return nil
}

func printProcesses(appID string, processes []api.Process) {
	psMap := ps.ByType(processes)

	fmt.Printf("=== %s Processes\n", appID)

	for psType, procs := range psMap {
		fmt.Printf("--- %s:\n", psType)

		for _, proc := range procs {
			fmt.Printf("%s.%d %s (%s)\n", proc.Type, proc.Num, proc.State, proc.Release)
		}
	}
}
