package resource

// ResourceTuple groups together cpu, memory and disk space. This could be
// total, available or consumed. It could also be used by job resource requirements.
type ResourceTuple struct {
	// in hundreds, ie 100=1core, 50=0.5core, 200=2cores, etc
	Cores int
	// in MB
	Memory int
	// in MB
	Disk int
}
