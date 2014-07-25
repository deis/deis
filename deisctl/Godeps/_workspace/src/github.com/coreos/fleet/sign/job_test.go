package sign

import (
	"bytes"
	"testing"

	gossh "github.com/coreos/fleet/third_party/code.google.com/p/gosshnew/ssh"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/unit"
)

func TestSignJob(t *testing.T) {
	c, _ := initSign(t)

	u, err := unit.NewUnit("Echo")
	if err != nil {
		t.Fatalf("unexpected error creating new unit: %v", err)
	}
	j := job.NewJob("echo.service", *u)

	data, err := marshal(u)
	if err != nil {
		t.Fatal("marshal error:", err)
	}

	expectedSig, err := c.keyring.Sign(testPublicKeys["rsa"], data)
	if err != nil {
		t.Fatal("sign error:", err)
	}

	s, err := c.SignJob(j)
	if err != nil {
		t.Fatal("sign payload error:", err)
	}
	if s.Tag != TagForJob("echo.service") {
		t.Fatal("sign tag error:", err)
	}

	if len(s.Signatures) != 1 {
		t.Fatal("expect 1 signature instead of", len(s.Signatures))
	}
	if bytes.Compare(s.Signatures[0].Blob, expectedSig.Blob) != 0 {
		t.Fatal("wrong signature")
	}
}

func TestVerifyJob(t *testing.T) {
	c, v := initSign(t)

	u, err := unit.NewUnit("Echo")
	if err != nil {
		t.Fatalf("unexpected error creating new unit: %v", err)
	}
	j := job.NewJob("echo.service", *u)

	data, err := marshal(u)
	if err != nil {
		t.Fatal("marshal error:", err)
	}

	v.pubkeys = append(v.pubkeys, testPublicKeys["rsa"])
	signature, err := c.keyring.Sign(testPublicKeys["rsa"], data)
	if err != nil {
		t.Fatal("sign error:", err)
	}

	ss := &SignatureSet{TagForJob("echo.service"), []*gossh.Signature{signature}}

	ok, err := v.VerifyJob(j, ss)
	if err != nil {
		t.Fatal("error verifying job:", err)
	}
	if !ok {
		t.Fatal("job verification failed")
	}

	ss.Tag = ""
	ok, err = v.VerifyJob(j, ss)
	if err == nil || ok == true {
		t.Fatal("should fail on job verification")
	}

	ok, err = v.VerifyJob(j, nil)
	if err == nil || ok == true {
		t.Fatal("should fail on job verification")
	}
}
