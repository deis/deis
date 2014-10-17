package machine

import (
	"github.com/coreos/fleet/log"
	"github.com/coreos/fleet/pkg"
)

type Machine interface {
	State() MachineState
}

// HasMetadata determine if the Metadata of a given MachineState
// matches the indicated values.
func HasMetadata(state *MachineState, metadata map[string]pkg.Set) bool {
	for key, values := range metadata {
		local, ok := state.Metadata[key]
		if !ok {
			log.V(1).Infof("No local values found for Metadata(%s)", key)
			return false
		}

		log.V(1).Infof("Asserting local Metadata(%s) meets requirements", key)

		if values.Contains(local) {
			log.V(1).Infof("Local Metadata(%s) meets requirement", key)
		} else {
			log.V(1).Infof("Local Metadata(%s) does not match requirement", key)
			return false
		}
	}

	return true
}
