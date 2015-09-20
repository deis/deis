package fleet

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func nextUnitNum(units []string) (num int, err error) {
	count, err := countUnits(units)
	if err != nil {
		return
	}
	sort.Ints(count)
	num = 1
	for _, i := range count {
		if num < i {
			return num, nil
		}
		num++
	}
	return num, nil
}

func lastUnitNum(units []string) (num int, err error) {
	count, err := countUnits(units)
	if err != nil {
		return
	}
	num = 1
	sort.Sort(sort.Reverse(sort.IntSlice(count)))
	if len(count) == 0 {
		return num, fmt.Errorf("Component not found")
	}
	return count[0], nil
}

func countUnits(units []string) (count []int, err error) {
	for _, unit := range units {
		_, n, err := splitJobName(unit)
		if err != nil {
			return []int{}, err
		}
		count = append(count, n)
	}
	return
}

func splitJobName(component string) (c string, num int, err error) {
	r := regexp.MustCompile(`deis\-([a-z-]+)\@([\d]+)\.service`)
	match := r.FindStringSubmatch(component)
	if len(match) == 0 {
		c, err = "", fmt.Errorf("Could not parse component: %v", component)
		return
	}
	c = match[1]
	num, err = strconv.Atoi(match[2])
	if err != nil {
		return
	}
	return
}

func splitTarget(target string) (component string, num int, err error) {
	// see if we were provided a specific target
	r := regexp.MustCompile(`^([a-z-]+)(@\d+)?(\.service)?$`)
	match := r.FindStringSubmatch(target)
	// check for failed match
	if len(match) < 3 {
		err = fmt.Errorf("Could not parse target: %v", target)
		return
	}
	if match[2] == "" {
		component = match[1]
		return component, 0, nil
	}
	num, err = strconv.Atoi(match[2][1:])
	if err != nil {
		return "", 0, err
	}
	return match[1], num, err
}

// expand a target to all installed units
func (c *FleetClient) expandTargets(targets []string) (expandedTargets []string, err error) {
	for _, t := range targets {
		// ensure unit name starts with "deis-"
		if !strings.HasPrefix(t, "deis-") {
			t = "deis-" + t
		}
		if strings.HasSuffix(t, "@*") {
			var targets []string
			targets, err = c.Units(strings.TrimSuffix(t, "@*"))
			if err != nil {
				return
			}
			expandedTargets = append(expandedTargets, targets...)
		} else {
			if !strings.HasSuffix(t, ".service") {
				t += ".service"
			}

			expandedTargets = append(expandedTargets, t)
		}

	}
	return
}

// randomValue returns a random string from a slice of string
func randomValue(src []string) string {
	s := rand.NewSource(int64(time.Now().Unix()))
	r := rand.New(s)
	idx := r.Intn(len(src))
	return src[idx]
}
