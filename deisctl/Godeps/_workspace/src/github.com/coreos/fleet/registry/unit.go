package registry

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/log"
	"github.com/coreos/fleet/unit"
)

const (
	unitPrefix = "/unit/"
	// Legacy versions of fleet stored payloads instead of units
	payloadPrefix = "/payload/"
)

func (r *EtcdRegistry) storeOrGetUnitFile(u unit.UnitFile) (err error) {
	um := unitModel{
		Raw: u.String(),
	}

	json, err := marshal(um)
	if err != nil {
		return err
	}

	req := etcd.Create{
		Key:   r.hashedUnitPath(u.Hash()),
		Value: json,
	}
	_, err = r.etcd.Do(&req)
	// unit is already stored
	if err != nil && isNodeExist(err) {
		// TODO(jonboulle): verify more here?
		err = nil
	}
	return
}

// getUnitFromLegacyPayload tries to extract a Unit from a legacy JobPayload of the given name
func (r *EtcdRegistry) getUnitFromLegacyPayload(name string) (*unit.UnitFile, error) {
	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, payloadPrefix, name),
		Recursive: true,
		Sorted:    true,
	}
	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isKeyNotFound(err) {
			err = nil
		}
		return nil, err
	}

	var ljp LegacyJobPayload
	if err := unmarshal(resp.Node.Value, &ljp); err != nil {
		return nil, fmt.Errorf("error unmarshaling LegacyJobPayload(%s): %v", name, err)
	}
	if ljp.Name != name {
		return nil, fmt.Errorf("payload name in Registry (%s) does not match expected name (%s)", ljp.Name, name)
	}
	// After the unmarshaling, the LegacyPayload should contain a fully hydrated Unit
	return &ljp.Unit, nil
}

// getUnitByHash retrieves from the Registry the Unit associated with the given Hash
func (r *EtcdRegistry) getUnitByHash(hash unit.Hash) *unit.UnitFile {
	req := etcd.Get{
		Key:       r.hashedUnitPath(hash),
		Recursive: true,
	}
	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isKeyNotFound(err) {
			err = nil
		}
		return nil
	}
	var um unitModel
	if err := unmarshal(resp.Node.Value, &um); err != nil {
		log.Errorf("error unmarshaling Unit(%s): %v", hash, err)
		return nil
	}

	u, err := unit.NewUnitFile(um.Raw)
	if err != nil {
		log.Errorf("error parsing Unit(%s): %v", hash, err)
		return nil
	}

	return u
}

func (r *EtcdRegistry) hashedUnitPath(hash unit.Hash) string {
	return path.Join(r.keyPrefix, unitPrefix, hash.String())
}

// LegacyJobPayload deals with the legacy concept of a "JobPayload" (deprecated by UnitFiles).
// The associated marshaling/unmarshaling methods deal with Payloads encoded in this legacy format.
type LegacyJobPayload struct {
	Name string
	Unit unit.UnitFile
}

func (ljp *LegacyJobPayload) UnmarshalJSON(data []byte) error {
	var ljpm legacyJobPayloadModel
	err := json.Unmarshal(data, &ljpm)
	if err != nil {
		return fmt.Errorf("unable to JSON-deserialize object: %s", err)
	}

	var u *unit.UnitFile
	if len(ljpm.Unit.Raw) > 0 {
		u, err = unit.NewUnitFile(ljpm.Unit.Raw)
	} else {
		u, err = unit.NewUnitFromLegacyContents(ljpm.Unit.Contents)
	}
	if err != nil {
		return err
	}

	ljp.Unit = *u
	ljp.Name = ljpm.Name

	return nil
}

// legacyJobPayloadModel is an abstraction to deal with serialized LegacyJobPayloads
type legacyJobPayloadModel struct {
	Name string
	Unit struct {
		Contents map[string]map[string]string
		Raw      string
	}
}

type unitModel struct {
	Raw string
}
