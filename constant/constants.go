package constant

import "time"

const (
	UnitsDir        = "/var/lib/deis/units/"
	HooksDir        = "/var/lib/deis/hooks/"
	Version         = "/etc/deis-version"
	MachineID       = "/etc/machine-id"
	UpdatekeyDir    = "/deis/update/"
	InitialInterval = time.Second * 10
	MaxInterval     = time.Minute * 7
)
