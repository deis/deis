package constant

import "time"

const (
	// UnitsDir is the default directory for unit files.
	UnitsDir = "/var/lib/deis/units/"

	// HooksDir is the default directory for hook scripts.
	HooksDir = "/var/lib/deis/hooks/"

	// Version is the location of the deis-version text file.
	Version = "/etc/deis-version"

	// MachineID is the location of the machine-id file.
	MachineID = "/etc/machine-id"

	// UpdatekeyDir is the etcd directory for update data.
	UpdatekeyDir = "/deis/update/"

	// InitialInterval specifies how long to wait at first between update loops.
	InitialInterval = time.Second * 10

	// MaxInterval is the longest time allowed to wait between loops.
	MaxInterval = time.Minute * 7
)
