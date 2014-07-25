package registry

import (
	"path"
	"strings"

	log "github.com/coreos/fleet/third_party/github.com/golang/glog"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/event"
	"github.com/coreos/fleet/job"
)

const (
	offerPrefix = "offer"
)

// jobOfferModel is used for serializing and deserializing JobOffers stored in the Registry
type jobOfferModel struct {
	Job        jobModel
	MachineIDs []string
}

// CreateJobOffer attempts to store a JobOffer and a reference to its associated Job in the repository
func (r *EtcdRegistry) CreateJobOffer(jo *job.JobOffer) error {
	jom := jobOfferModel{
		Job: jobModel{
			Name:     jo.Job.Name,
			UnitHash: jo.Job.UnitHash,
		},
		MachineIDs: jo.MachineIDs,
	}

	json, err := marshal(jom)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	req := etcd.Set{
		Key:   path.Join(r.keyPrefix, offerPrefix, jo.Job.Name, "object"),
		Value: json,
	}
	_, err = r.etcd.Do(&req)
	return err
}

// getJobOfferFromJSON hydrates a JobOffer from a JSON-encoded jobOfferModel
func (r *EtcdRegistry) getJobOfferFromJSON(val string) *job.JobOffer {
	var jom jobOfferModel
	if err := unmarshal(val, &jom); err != nil {
		return nil
	}

	j := r.getJobFromModel(jom.Job)
	if j == nil {
		return nil
	}

	jo := job.JobOffer{
		Job:        *j,
		MachineIDs: jom.MachineIDs,
	}

	return &jo
}

// Bids returns a list of JobBids that have been submitted for the given JobOffer
func (r *EtcdRegistry) Bids(jo *job.JobOffer) ([]job.JobBid, error) {
	var bids []job.JobBid

	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, offerPrefix, jo.Job.Name, "bids"),
		Recursive: true,
	}
	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isKeyNotFound(err) {
			return bids, nil
		}
		return nil, err
	}

	for _, node := range resp.Node.Nodes {
		machID := path.Base(node.Key)
		jb := job.NewBid(jo.Job.Name, machID)
		bids = append(bids, *jb)
	}

	return bids, nil
}

// UnresolvedJobOffers returns a list of hydrated JobOffers from the Registry
func (r *EtcdRegistry) UnresolvedJobOffers() []job.JobOffer {
	var offers []job.JobOffer

	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, offerPrefix),
		Sorted:    true,
		Recursive: true,
	}
	resp, err := r.etcd.Do(&req)

	if err != nil {
		return offers
	}

	for _, node := range resp.Node.Nodes {
		req := etcd.Get{
			Key: path.Join(node.Key, "object"),

			//TODO(bcwaldon): This request should not need to be sorted/recursive
			Sorted:    true,
			Recursive: true,
		}
		resp, err := r.etcd.Do(&req)

		// The object was probably handled between when we attempted to
		// start resolving offers and when we actually tried to get it
		if err != nil {
			continue
		}

		jo := r.getJobOfferFromJSON(resp.Node.Value)
		if jo == nil {
			continue
		}

		offers = append(offers, *jo)
	}

	return offers
}

func (r *EtcdRegistry) LockJobOffer(jobName, context string) *TimedResourceMutex {
	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, offerPrefix, jobName),
		Recursive: true,
	}
	_, err := r.etcd.Do(&req)
	if err != nil {
		return nil
	}

	return r.lockResource("offer", jobName, context)
}

func (r *EtcdRegistry) ResolveJobOffer(jobName string) error {
	req := etcd.Delete{
		Key: path.Join(r.keyPrefix, offerPrefix, jobName, "object"),
	}
	if _, err := r.etcd.Do(&req); err != nil {
		if !isKeyNotFound(err) {
			return err
		}
	}

	req = etcd.Delete{
		Key:       path.Join(r.keyPrefix, offerPrefix, jobName),
		Recursive: true,
	}

	r.etcd.Do(&req)
	return nil
}

func (r *EtcdRegistry) SubmitJobBid(jb *job.JobBid) {
	req := etcd.Set{
		Key: path.Join(r.keyPrefix, offerPrefix, jb.JobName, "bids", jb.MachineID),
	}
	r.etcd.Do(&req)
}

func (es *EventStream) filterEventJobOffered(resp *etcd.Result) *event.Event {
	if resp.Action != "set" {
		return nil
	}

	dir, base := path.Split(resp.Node.Key)

	if base != "object" {
		return nil
	}

	dir = path.Dir(strings.TrimSuffix(dir, "/"))
	prefix := path.Base(strings.TrimSuffix(dir, "/"))

	if prefix != offerPrefix {
		return nil
	}

	jo := es.registry.getJobOfferFromJSON(resp.Node.Value)
	if jo == nil {
		return nil
	}

	return &event.Event{"EventJobOffered", *jo, nil}
}

func filterEventJobBidSubmitted(resp *etcd.Result) *event.Event {
	if resp.Action != "set" {
		return nil
	}

	dir, machID := path.Split(resp.Node.Key)
	dir, prefix := path.Split(strings.TrimSuffix(dir, "/"))

	if prefix != "bids" {
		return nil
	}

	dir, jobName := path.Split(strings.TrimSuffix(dir, "/"))
	prefix = path.Base(strings.TrimSuffix(dir, "/"))

	if prefix != offerPrefix {
		return nil
	}

	jb := job.NewBid(jobName, machID)
	return &event.Event{"EventJobBidSubmitted", *jb, nil}
}
