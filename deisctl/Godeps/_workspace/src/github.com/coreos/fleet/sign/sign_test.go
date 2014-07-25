package sign

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	gossh "github.com/coreos/fleet/third_party/code.google.com/p/gosshnew/ssh"
	gosshagent "github.com/coreos/fleet/third_party/code.google.com/p/gosshnew/ssh/agent"
)

const (
	privateKeyFile                = "../fixtures/insecure_private_key"
	authorizedKeysFile            = "../fixtures/authorized_keys"
	nonexistentAuthorizedKeysFile = "../fixtures/nonexistent_authorized_keys"
	signatureForGood              = "903a528536371744b4f7f3720e531321f128e164254600dba3753e26aa0bd4d6f86cd9da2f4463aca90549427f26547df821ff403722825651abbdb5a674b9bab07ed89a0b89e249cf93341325dd243300dc72a168b0faf06d3e154efd7e42f24aba312407b658634cb89e2f37d19cb7341feba9aca591d09894da4b9d5fe2713f69408a8d7c3fe28fbe07e80b2b1617158b510aadb487e37baf33a2497ffeb2e2e4091ec1a025adc59acae1ea28ee41632806389ffefc47272ef37405cf1c30933e427b8996106df6ca4cd4e5fa8c8f27d4ef74b8a4632d687ef2ccb520015034c72573ed4c95d927b53732bde72641eeb438c8e8f9374d091ba8deb2bfd929"
	firstAuthorizedKey            = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDMmq9hew++XUyYKDRIuq4K3VaVJUaE76LscrJ4Ov+UPJ0nTm0/VH2z0eOX9fQijmZ3c0/uMN03bdTfZG2w4TLYwxkgtIgch6nYG540oeKGHfcx3D/LYQ1isTwlLFelSAnDjaIsiLmxv0XHc4lojhLEtjf1OyHMf06snQscizYTmin29/7qSehf9WBEAPxdMuBGWYMi4yHnDn0cT4b7iowzZ+LQFjhZDthz2WTdSqofHbjPQSLGm65IotCJh8WRROKYPVLqnlZtQV7ntkzxsDSVpv5gsGMfZpuF1LGkQ89p/dCvpShoiklORMDA0Stm0wSemoKkwWvaYTbiyj1ZreXl\n"
	secondAuthorizedKey           = "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ==\n"
)

type badStubKeyring struct {
}

func newBadStubKeyring() *badStubKeyring {
	return &badStubKeyring{}
}

func (k *badStubKeyring) List() ([]*gosshagent.Key, error) {
	return nil, errors.New("")
}

func (k *badStubKeyring) Sign(key gossh.PublicKey, data []byte) (*gossh.Signature, error) {
	return nil, errors.New("")
}

func (k *badStubKeyring) Add(s interface{}, cert *gossh.Certificate, comment string) error {
	return nil
}

func (k *badStubKeyring) Lock(passphrase []byte) error {
	return nil
}

func (k *badStubKeyring) Unlock(passphrase []byte) error {
	return nil
}

func (k *badStubKeyring) Remove(key gossh.PublicKey) error {
	return nil
}

func (k *badStubKeyring) RemoveAll() error {
	return nil
}

func (k *badStubKeyring) Signers() ([]gossh.Signer, error) {
	return nil, errors.New("")
}

func initSign(t *testing.T) (*SignatureCreator, *SignatureVerifier) {
	keyring := gosshagent.NewKeyring()

	keyring.Add(testPrivateKeys["rsa"], nil, "")

	c := NewSignatureCreator(keyring)
	v, err := NewSignatureVerifierFromKeyring(keyring)
	if err != nil {
		t.Fatal("Failed to initialize signature verifier:", err)
	}

	return c, v
}

// TestNewSignatureVerifierFromFile tests initializing SignatureVerifier from file
func TestNewSignatureVerifierFromFile(t *testing.T) {
	v, err := NewSignatureVerifierFromAuthorizedKeysFile(authorizedKeysFile)
	if err != nil {
		t.Error("failed to read from authorized keys file:", err)
	}

	keys := v.pubkeys
	if bytes.Compare(gossh.MarshalAuthorizedKey(keys[0]), []byte(firstAuthorizedKey)) != 0 {
		t.Error("wrong first authorized key")
	}
	if bytes.Compare(gossh.MarshalAuthorizedKey(keys[1]), []byte(secondAuthorizedKey)) != 0 {
		t.Error("wrong second authorized key")
	}
}

// TestBadNewSignatureVerifierFromFile tests initializing SignatureVerifier from file incorrectly
func TestBadNewSignatureVerifierFromFile(t *testing.T) {
	_, err := NewSignatureVerifierFromAuthorizedKeysFile(nonexistentAuthorizedKeysFile)
	if err == nil {
		t.Error("new signature verifier succeeded from nonexistent keyfile")
	}

	_, err = NewSignatureVerifierFromAuthorizedKeysFile(privateKeyFile)
	if err == nil {
		t.Error("new signature verifier succeeded from private key file")
	}

	emptyFile, err := ioutil.TempFile(os.TempDir(), "empty_authorized_keys")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(emptyFile.Name())

	_, err = NewSignatureVerifierFromAuthorizedKeysFile(emptyFile.Name())
	if err == nil {
		t.Error("new signature verifier succeeded from empty key file")
	}
}

// TestNewSignatureVerifierFromKeyring tests initializing SignatureVerifier from keyring
func TestNewSignatureVerifierFromKeyring(t *testing.T) {
	c, _ := initSign(t)
	v, err := NewSignatureVerifierFromKeyring(c.keyring)
	if err != nil {
		t.Error("fail to new signature verifier")
	}
	keys := len(v.pubkeys)
	// there should be at least one added in the tests.
	if keys == 0 {
		t.Error("expected pubkeys but found none")
	}
}

// TestBadNewSignatureVerifierFromKeyring tests initializing SignatureVerifier from keyring incorrectly
func TestBadNewSignatureVerifierFromKeyring(t *testing.T) {
	_, err := NewSignatureVerifierFromKeyring(newBadStubKeyring())
	if err == nil {
		t.Error("new signature verifier succeeded from bad keyring")
	}
}

// TestSign tests the creation of correct signature
func TestSign(t *testing.T) {
	c, _ := initSign(t)
	tag := "1"
	data := []byte("Good")

	expectedSig, err := c.keyring.Sign(testPublicKeys["rsa"], data)
	if err != nil {
		t.Error("Sign:", err)
	}

	sig, err := c.Sign(tag, data)
	if err != nil {
		t.Error("sign failed:", err)
	}

	if sig.Tag != tag {
		t.Errorf("expected tag %v instead of %v", tag, sig.Tag)
	}

	if len(sig.Signatures) == 0 {
		t.Fatal("Expected signatures for data 'Good'")
	}

	if bytes.Compare(sig.Signatures[0].Blob, expectedSig.Blob) != 0 {
		t.Fatal("wrong signature for data 'Good'")
	}
}

// TestBadSign tests the incorrect creation of correct signature
func TestBadSign(t *testing.T) {
	c, _ := initSign(t)
	tag := "1"
	data := []byte("Good")

	c.keyring.RemoveAll()
	sig, err := c.Sign(tag, data)
	if sig != nil || err == nil {
		t.Error("sign succeeded with empty keyring")
	}

	c.keyring = newBadStubKeyring()
	sig, err = c.Sign(tag, nil)
	if sig != nil || err == nil {
		t.Error("sign succeeded with bad keyring")
	}

	c.keyring = nil
	sig, err = c.Sign(tag, data)
	if sig != nil || err == nil {
		t.Error("sign succeeded with nil keyring")
	}
}

// TestVerify tests the verification of correct signature
func TestVerify(t *testing.T) {
	c, v := initSign(t)
	badData := []byte("Bad")
	data := []byte("Good")

	v.pubkeys = append(v.pubkeys, testPublicKeys["rsa"])
	sig, err := c.keyring.Sign(testPublicKeys["rsa"], data)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	set := &SignatureSet{"1", []*gossh.Signature{sig}}

	ok, err := v.Verify(data, set)
	if err != nil {
		t.Error("failed to verify signature:", err)
	}
	if !ok {
		t.Error("failed to verify signature is correct")
	}

	ok, err = v.Verify(badData, set)
	if err != nil {
		t.Error("failed to verify signature:", err)
	}
	if ok {
		t.Error("failed to verify signature is incorrect")
	}
}

// TestBadVerify tests the incorrect verification of correct signature
func TestBadVerify(t *testing.T) {
	_, v := initSign(t)
	data := []byte("Good")

	v.pubkeys = nil
	ok, err := v.Verify(data, nil)
	if ok != false || err == nil {
		t.Error("verify with nil signature set succeeded")
	}
}

// TestBadMarshal tests incorrect marshal
func TestBadMarshal(t *testing.T) {
	c := make(chan bool)
	_, err := marshal(c)
	if err == nil {
		t.Error("marshaling empty channel succeeded")
	}
}
