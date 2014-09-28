package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/deis/deisctl/backend"
)

func Scale(b backend.Backend, targets []string) error {
	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		// the router is the only component that can scale past 1 at the moment
		if num > 1 && !strings.Contains(component, "router") {
			return fmt.Errorf("cannot scale %s past 1", component)
		}
		if err := b.Scale(component, num); err != nil {
			return err
		}
	}
	return nil
}

func splitScaleTarget(target string) (c string, num int, err error) {
	r := regexp.MustCompile(`([a-z-]+)=([\d]+)`)
	match := r.FindStringSubmatch(target)
	if len(match) == 0 {
		err = fmt.Errorf("Could not parse: %v", target)
		return
	}
	c = match[1]
	num, err = strconv.Atoi(match[2])
	if err != nil {
		return
	}
	return
}
