package registry

import "github.com/coreos/fleet/Godeps/_workspace/src/github.com/coreos/go-semver/semver"

// LatestVersion attempts to retrieve the latest version of fleet that has
// been registered in the Registry. It returns the version if it can be
// determined (or nil otherwise), and any error encountered.
func (r *EtcdRegistry) LatestVersion() (*semver.Version, error) {
	machs, err := r.Machines()
	if err != nil {
		if isKeyNotFound(err) {
			err = nil
		}
		return nil, err
	}
	var lv *semver.Version
	for _, m := range machs {
		v, err := semver.NewVersion(m.Version)
		if err != nil {
			continue
		} else if lv == nil || lv.LessThan(*v) {
			lv = v
		}
	}
	return lv, nil
}
