package registry

import (
	"path"
	"strconv"

	"github.com/coreos/fleet/Godeps/_workspace/src/github.com/coreos/go-semver/semver"

	"github.com/coreos/fleet/etcd"
)

// LatestDaemonVersion attempts to retrieve the latest version of fleetd
// that has been registered in the Registry. It returns the version if
// it can be determined (or nil otherwise), and any error encountered.
func (r *EtcdRegistry) LatestDaemonVersion() (*semver.Version, error) {
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

// EngineVersion implements the ClusterRegistry interface
func (r *EtcdRegistry) EngineVersion() (int, error) {
	req := etcd.Get{
		Key: r.engineVersionPath(),
	}

	res, err := r.etcd.Do(&req)
	if err != nil {
		// no big deal, either the cluster is new or is just
		// upgrading from old unversioned code
		if isKeyNotFound(err) {
			err = nil
		}
		return 0, err
	}

	return strconv.Atoi(res.Node.Value)
}

// UpdateEngineVersion implements the ClusterRegistry interface
func (r *EtcdRegistry) UpdateEngineVersion(from, to int) error {
	key := r.engineVersionPath()

	strTo := strconv.Itoa(to)
	strFrom := strconv.Itoa(from)

	var req etcd.Action
	req = &etcd.Set{
		Key:           key,
		Value:         strTo,
		PreviousValue: strFrom,
	}

	_, err := r.etcd.Do(req)
	if err == nil {
		return nil
	} else if !isKeyNotFound(err) {
		return err
	}

	req = &etcd.Create{
		Key:   key,
		Value: strTo,
	}

	_, err = r.etcd.Do(req)
	return err
}

func (r *EtcdRegistry) engineVersionPath() string {
	return path.Join(r.keyPrefix, "/engine/version")
}
