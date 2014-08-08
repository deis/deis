package deisctl

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
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
			return count, err
		}
		count = append(count, n)
	}
	return
}

func splitJobName(component string) (c string, num int, err error) {
	r := regexp.MustCompile(`deis\-([a-z-]+)\.([\d]+)\.service`)
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

func splitComponentTarget(target string) (c string, num int, err error) {
	r := regexp.MustCompile(`([a-z-]+)\.?([\d]+)?`)
	match := r.FindStringSubmatch(target)
	if len(match) == 0 {
		err = fmt.Errorf("Could not parse: %v", target)
		return
	}
	if match[2] == "" {
		return match[1], 0, nil
	}
	c = match[1]
	num, err = strconv.Atoi(match[2])
	if err != nil {
		return
	}
	return
}
