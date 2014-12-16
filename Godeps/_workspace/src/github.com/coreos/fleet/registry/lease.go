/*
   Copyright 2014 CoreOS, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package registry

import (
	"encoding/json"
	"path"
	"time"

	"github.com/coreos/fleet/etcd"
)

const (
	leasePrefix = "lease"
)

func (r *EtcdRegistry) leasePath(name string) string {
	return path.Join(r.keyPrefix, leasePrefix, name)
}

func (r *EtcdRegistry) GetLease(name string) (Lease, error) {
	key := r.leasePath(name)
	req := etcd.Get{
		Key: key,
	}

	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isKeyNotFound(err) {
			err = nil
		}
		return nil, err
	}

	l := leaseFromResult(resp, r.etcd)
	return l, nil
}

func (r *EtcdRegistry) StealLease(name, machID string, ver int, period time.Duration, idx uint64) (Lease, error) {
	val, err := serializeLeaseMetadata(machID, ver)
	if err != nil {
		return nil, err
	}

	req := etcd.Set{
		Key:           r.leasePath(name),
		Value:         val,
		PreviousIndex: idx,
		TTL:           period,
	}

	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isNodeExist(err) {
			err = nil
		}
		return nil, err
	}

	l := leaseFromResult(resp, r.etcd)
	return l, nil
}

func (r *EtcdRegistry) AcquireLease(name string, machID string, ver int, period time.Duration) (Lease, error) {
	val, err := serializeLeaseMetadata(machID, ver)
	if err != nil {
		return nil, err
	}

	req := etcd.Create{
		Key:   r.leasePath(name),
		Value: val,
		TTL:   period,
	}

	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isNodeExist(err) {
			err = nil
		}
		return nil, err
	}

	l := leaseFromResult(resp, r.etcd)
	return l, nil
}

type etcdLeaseMetadata struct {
	MachineID string
	Version   int
}

// etcdLease implements the Lease interface
type etcdLease struct {
	key  string
	meta etcdLeaseMetadata
	idx  uint64
	ttl  time.Duration
	etcd etcd.Client
}

func (l *etcdLease) Release() error {
	req := etcd.Delete{
		Key:           l.key,
		PreviousIndex: l.idx,
	}
	_, err := l.etcd.Do(&req)
	return err
}

func (l *etcdLease) Renew(period time.Duration) error {
	val, err := serializeLeaseMetadata(l.meta.MachineID, l.meta.Version)
	req := etcd.Set{
		Key:           l.key,
		Value:         val,
		PreviousIndex: l.idx,
		TTL:           period,
	}

	resp, err := l.etcd.Do(&req)
	if err != nil {
		return err
	}

	renewed := leaseFromResult(resp, l.etcd)
	*l = *renewed

	return nil
}

func (l *etcdLease) MachineID() string {
	return l.meta.MachineID
}

func (l *etcdLease) Version() int {
	return l.meta.Version
}

func (l *etcdLease) Index() uint64 {
	return l.idx
}

func (l *etcdLease) TimeRemaining() time.Duration {
	return l.ttl
}

func leaseFromResult(res *etcd.Result, ec etcd.Client) *etcdLease {
	lease := &etcdLease{
		key:  res.Node.Key,
		idx:  res.Node.ModifiedIndex,
		ttl:  res.Node.TTLDuration(),
		etcd: ec,
	}

	err := json.Unmarshal([]byte(res.Node.Value), &lease.meta)

	// fall back to using the entire value as the MachineID for
	// backwards-compatibility with engines that are not aware
	// of this versioning mechanism
	if err != nil {
		lease.meta = etcdLeaseMetadata{
			MachineID: res.Node.Value,
			Version:   0,
		}
	}

	return lease
}

func serializeLeaseMetadata(machID string, ver int) (string, error) {
	meta := etcdLeaseMetadata{
		MachineID: machID,
		Version:   ver,
	}

	b, err := json.Marshal(meta)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
