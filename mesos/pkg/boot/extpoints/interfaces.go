package extpoints

import (
	"github.com/deis/deis/mesos/pkg/types"
)

// BootComponent interface that defines the steps
// required to initialize a component
type BootComponent interface {
	// EtcdDefaults required initial values in etcd
	EtcdDefaults() map[string]string

	// MkdirsEtcd required directories in etcd
	MkdirsEtcd() []string

	// PreBoot custom pre-boot task (custom go code)
	PreBoot(currentBoot *types.CurrentBoot)

	// PreBootScripts scripts to execute before the component starts
	PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script

	// UseConfd is required the use of confd?
	UseConfd() (bool, bool)

	// BootDaemons required commands to start the component
	BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon

	// WaitForPorts ports that must be open to indicate that the component is running
	WaitForPorts() []int

	// PostBootScripts scripts to execute after the component starts
	PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script

	// PostBoot custom post-boot task (custom go code)
	PostBoot(currentBoot *types.CurrentBoot)

	// ScheduleTasks tasks that must run during the lifecycle of the component
	ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron

	// PreShutdownScripts scripts to execute before the component execution ends
	PreShutdownScripts(currentBoot *types.CurrentBoot) []*types.Script
}
