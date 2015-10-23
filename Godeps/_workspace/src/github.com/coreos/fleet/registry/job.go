// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"errors"
	"fmt"
	"path"
	"sort"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/log"
	"github.com/coreos/fleet/unit"
)

const (
	jobPrefix = "job"
)

// Schedule returns all ScheduledUnits known by fleet, ordered by name
func (r *EtcdRegistry) Schedule() ([]job.ScheduledUnit, error) {
	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, jobPrefix),
		Sorted:    true,
		Recursive: true,
	}

	res, err := r.etcd.Do(&req)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			err = nil
		}
		return nil, err
	}

	heartbeats := make(map[string]string)
	uMap := make(map[string]*job.ScheduledUnit)

	for _, dir := range res.Node.Nodes {
		_, name := path.Split(dir.Key)
		u := &job.ScheduledUnit{
			Name:            name,
			TargetMachineID: dirToTargetMachineID(&dir),
		}
		heartbeats[name] = dirToHeartbeat(&dir)
		uMap[name] = u
	}

	states, err := r.statesByMUSKey()
	if err != nil {
		return nil, err
	}

	var sortable sort.StringSlice

	// Determine the JobState of each ScheduledUnit
	for name, su := range uMap {
		sortable = append(sortable, name)
		key := MUSKey{
			machID: su.TargetMachineID,
			name:   name,
		}
		us := states[key]
		js := determineJobState(heartbeats[name], su.TargetMachineID, us)
		su.State = &js
	}
	sortable.Sort()

	units := make([]job.ScheduledUnit, 0, len(sortable))
	for _, name := range sortable {
		units = append(units, *uMap[name])
	}
	return units, nil
}

// Units lists all Units stored in the Registry, ordered by name. This includes both global and non-global units.
func (r *EtcdRegistry) Units() ([]job.Unit, error) {
	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, jobPrefix),
		Sorted:    true,
		Recursive: true,
	}

	res, err := r.etcd.Do(&req)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			err = nil
		}
		return nil, err
	}

	uMap := make(map[string]*job.Unit)
	for _, dir := range res.Node.Nodes {
		u, err := r.dirToUnit(&dir)
		if err != nil {
			log.Errorf("Failed to parse Unit from etcd: %v", err)
			continue
		}
		if u == nil {
			continue
		}
		uMap[u.Name] = u
	}

	var sortable sort.StringSlice
	for name, _ := range uMap {
		sortable = append(sortable, name)
	}
	sortable.Sort()

	units := make([]job.Unit, 0, len(sortable))
	for _, name := range sortable {
		units = append(units, *uMap[name])
	}

	return units, nil
}

// Unit retrieves the Unit by the given name from the Registry. Returns nil if
// no such Unit exists, and any error encountered.
func (r *EtcdRegistry) Unit(name string) (*job.Unit, error) {
	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, jobPrefix, name),
		Recursive: true,
	}

	res, err := r.etcd.Do(&req)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			err = nil
		}
		return nil, err
	}

	return r.dirToUnit(res.Node)
}

// dirToUnit takes a Node containing a Job's constituent objects (in child
// nodes) and returns a *job.Unit, or any error encountered
func (r *EtcdRegistry) dirToUnit(dir *etcd.Node) (*job.Unit, error) {
	objKey := path.Join(dir.Key, "object")
	var objNode *etcd.Node
	for _, node := range dir.Nodes {
		node := node
		if node.Key == objKey {
			objNode = &node
		}
	}
	if objNode == nil {
		return nil, nil
	}
	u, err := r.getUnitFromObjectNode(objNode)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("unable to parse Unit in Registry at key %s", objKey)
	}
	if tgtstate := dirToTargetState(dir); tgtstate != "" {
		ts, err := job.ParseJobState(tgtstate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Unit(%s) target-state: %v", u.Name, err)
		}
		u.TargetState = ts
	}

	return u, nil
}

// ScheduledUnit retrieves the ScheduledUnit by the given name from the Registry.
// Returns nil if no such ScheduledUnit exists, and any error encountered.
func (r *EtcdRegistry) ScheduledUnit(name string) (*job.ScheduledUnit, error) {
	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, jobPrefix, name),
		Recursive: true,
	}

	res, err := r.etcd.Do(&req)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			err = nil
		}
		return nil, err
	}

	su := job.ScheduledUnit{
		Name:            name,
		TargetMachineID: dirToTargetMachineID(res.Node),
	}

	var us *unit.UnitState
	if len(su.TargetMachineID) > 0 {
		us, err = r.getUnitState(name, su.TargetMachineID)
		if err != nil {
			return nil, err
		}
	}

	js := determineJobState(dirToHeartbeat(res.Node), su.TargetMachineID, us)
	su.State = &js

	return &su, nil
}

func (r *EtcdRegistry) UnscheduleUnit(name, machID string) error {
	req := etcd.Delete{
		Key:           r.jobTargetAgentPath(name),
		PreviousValue: machID,
	}

	_, err := r.etcd.Do(&req)
	if etcd.IsKeyNotFound(err) {
		err = nil
	}

	return err
}

// getValueInDir takes a *etcd.Node containing a job, and returns the value of
// the given key within that directory (i.e. child node) as a string, or an
// empty string if the child node does not exist
func getValueInDir(dir *etcd.Node, key string) (value string) {
	valPath := path.Join(dir.Key, key)
	for _, node := range dir.Nodes {
		if node.Key == valPath {
			value = node.Value
			break
		}
	}
	return
}

func dirToTargetMachineID(dir *etcd.Node) (tgtMID string) {
	return getValueInDir(dir, "target")
}

func dirToTargetState(dir *etcd.Node) (tgtState string) {
	return getValueInDir(dir, "target-state")
}

func dirToHeartbeat(dir *etcd.Node) (heartbeat string) {
	return getValueInDir(dir, "job-state")
}

// getUnitFromObject takes a *etcd.Node containing a Unit's jobModel, and
// instantiates and returns a representative *job.Unit, transitively fetching the
// associated UnitFile as necessary
func (r *EtcdRegistry) getUnitFromObjectNode(node *etcd.Node) (*job.Unit, error) {
	var err error
	var jm jobModel
	if err = unmarshal(node.Value, &jm); err != nil {
		return nil, err
	}

	var unit *unit.UnitFile

	unit = r.getUnitByHash(jm.UnitHash)
	if unit == nil {
		log.Warningf("No Unit found in Registry for Job(%s)", jm.Name)
		return nil, nil
	}

	ju := &job.Unit{
		Name: jm.Name,
		Unit: *unit,
	}
	return ju, nil

}

// jobModel is used for serializing and deserializing Jobs stored in the Registry
type jobModel struct {
	Name     string
	UnitHash unit.Hash
}

// DestroyUnit removes a Job object from the repository. It does not yet remove underlying
// UnitFiles from the repository.
func (r *EtcdRegistry) DestroyUnit(name string) error {
	req := etcd.Delete{
		Key:       path.Join(r.keyPrefix, jobPrefix, name),
		Recursive: true,
	}

	_, err := r.etcd.Do(&req)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			err = errors.New("job does not exist")
		}

		return err
	}

	// TODO(jonboulle): add unit reference counting and actually destroying Units
	return nil
}

// CreateUnit attempts to store a Unit and its associated unit file in the registry
func (r *EtcdRegistry) CreateUnit(u *job.Unit) (err error) {
	if err := r.storeOrGetUnitFile(u.Unit); err != nil {
		return err
	}

	jm := jobModel{
		Name:     u.Name,
		UnitHash: u.Unit.Hash(),
	}
	json, err := marshal(jm)
	if err != nil {
		return
	}

	req := etcd.Create{
		Key:   path.Join(r.keyPrefix, jobPrefix, u.Name, "object"),
		Value: json,
	}

	_, err = r.etcd.Do(&req)
	if err != nil {
		if etcd.IsNodeExist(err) {
			err = errors.New("job already exists")
		}
		return
	}

	return r.SetUnitTargetState(u.Name, u.TargetState)
}

func (r *EtcdRegistry) SetUnitTargetState(name string, state job.JobState) error {
	req := etcd.Set{
		Key:   r.jobTargetStatePath(name),
		Value: string(state),
	}
	_, err := r.etcd.Do(&req)
	return err
}

func (r *EtcdRegistry) ScheduleUnit(name string, machID string) error {
	req := etcd.Create{
		Key:   r.jobTargetAgentPath(name),
		Value: machID,
	}
	_, err := r.etcd.Do(&req)
	return err
}

func (r *EtcdRegistry) jobTargetAgentPath(jobName string) string {
	return path.Join(r.keyPrefix, jobPrefix, jobName, "target")
}

func (r *EtcdRegistry) jobTargetStatePath(jobName string) string {
	return path.Join(r.keyPrefix, jobPrefix, jobName, "target-state")
}
