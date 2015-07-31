package fleet

import (
	"fmt"
	"io"
	"sync"
)

// RollingRestart for instance units
func (c *FleetClient) RollingRestart(component string, wg *sync.WaitGroup, out, ew io.Writer) {
	if component != "router" {
		fmt.Fprint(ew, "invalid component. supported for: router")
		return
	}

	components, err := c.Units(component)
	if err != nil {
		io.WriteString(ew, err.Error())
		return
	}
	if len(components) < 1 {
		fmt.Fprint(ew, "rolling restart requires at least 1 component")
		return
	}
	for num := range components {
		unitName := fmt.Sprintf("%s@%v", component, num+1)

		c.Stop([]string{unitName}, wg, out, ew)
		wg.Wait()
		c.Destroy([]string{unitName}, wg, out, ew)
		wg.Wait()
		c.Create([]string{unitName}, wg, out, ew)
		wg.Wait()
		c.Start([]string{unitName}, wg, out, ew)
		wg.Wait()
	}
}
