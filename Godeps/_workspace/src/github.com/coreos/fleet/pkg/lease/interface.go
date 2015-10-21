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

package lease

import "time"

// Lease proxies to an auto-expiring lease stored in a LeaseRegistry.
// The creator of a Lease must repeatedly call Renew to keep their lease
// from expiring.
type Lease interface {
	// Renew attempts to extend the Lease TTL to the provided duration.
	// The operation will succeed only if the Lease has not changed in
	// the LeaseRegistry since it was last renewed or first acquired.
	// An error is returned if the Lease has already expired, or if the
	// operation fails for any other reason.
	Renew(time.Duration) error

	// Release relinquishes the ownership of a Lease back to the Registry.
	// After calling Release, the Lease object should be discarded. An
	// error is returned if the Lease has already expired, or if the
	// operation fails for any other reason.
	Release() error

	// MachineID returns the ID of the Machine that holds this Lease. This
	// value must be considered a cached value as it is not guaranteed to
	// be correct.
	MachineID() string

	// Version returns the current version at which the lessee is operating.
	// This value has the same correctness guarantees as MachineID.
	// It is up to the caller to determine what this Version means.
	Version() int

	// Index exposes the relative time at which the Lease was created or
	// renewed. For example, this could be implemented as the ModifiedIndex
	// field of a node in etcd.
	Index() uint64

	// TimeRemaining represents the amount of time left on the Lease when
	// it was fetched from the LeaseRegistry.
	TimeRemaining() time.Duration
}

type Manager interface {
	// GetLease fetches a Lease only if it exists. If it does not
	// exist, a nil Lease will be returned. Any other failures
	// result in non-nil error and nil Lease objects.
	GetLease(name string) (Lease, error)

	// AcquireLease acquires a named lease only if the lease is not
	// currently held. If a Lease cannot be acquired, a nil Lease
	// object is returned. An error is returned only if there is a
	// failure communicating with the Registry.
	AcquireLease(name, machID string, ver int, period time.Duration) (Lease, error)

	// StealLease attempts to replace the lessee of the Lease identified
	// by the provided name and index with a new lessee. This function
	// will fail if the named Lease has progressed past the given index.
	StealLease(name, machID string, ver int, period time.Duration, idx uint64) (Lease, error)
}
