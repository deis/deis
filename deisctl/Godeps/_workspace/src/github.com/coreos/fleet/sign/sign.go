package sign

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	gossh "github.com/coreos/fleet/third_party/code.google.com/p/gosshnew/ssh"
	gosshagent "github.com/coreos/fleet/third_party/code.google.com/p/gosshnew/ssh/agent"
	log "github.com/coreos/fleet/third_party/github.com/golang/glog"

	"github.com/coreos/fleet/pkg"
	"github.com/coreos/fleet/ssh"
)

const DefaultAuthorizedKeysFile = "~/.ssh/authorized_keys"

// SignatureSet contains a set of SSH signatures for a blob of data, and is named by a Tag.
type SignatureSet struct {
	Tag        string
	Signatures []*gossh.Signature
}

// SignatureCreator provides the ability to sign a blob of data with multiple SSH public keys, contained in the keyring
type SignatureCreator struct {
	keyring gosshagent.Agent
}

// NewSignatureCreator instantiates a SignatureCreator with the given keyring
func NewSignatureCreator(keyring gosshagent.Agent) *SignatureCreator {
	return &SignatureCreator{keyring}
}

// NewSignatureCreatorFromSSHAgent return a SignatureCreator which uses the local ssh-agent as its keyring
func NewSignatureCreatorFromSSHAgent() (*SignatureCreator, error) {
	keyring, err := ssh.SSHAgentClient()
	if err != nil {
		return nil, err
	}
	return &SignatureCreator{keyring}, nil
}

// Sign generates a SignatureSet for the given data, labelled by the supplied tag. It returns a *SignatureSet and any error encountere
func (sc *SignatureCreator) Sign(tag string, data []byte) (*SignatureSet, error) {
	if sc.keyring == nil {
		return nil, errors.New("signature creator is uninitialized")
	}

	sigs := make([]*gossh.Signature, 0)

	keys, err := sc.keyring.List()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, errors.New("signature creator keyring is empty")
	}

	for _, k := range keys {
		sig, err := sc.keyring.Sign(k, data)
		if err != nil {
			return nil, err
		}

		sigs = append(sigs, sig)
	}

	return &SignatureSet{tag, sigs}, nil
}

type SignatureVerifier struct {
	// keys is used to verify signing, created when needed
	pubkeys []gossh.PublicKey
}

func NewSignatureVerifier() *SignatureVerifier {
	return &SignatureVerifier{}
}

// NewSignatureVerifierFromSSHAgent return SignatureVerifier which uses public keys in the local ssh-agent to verify signatures
func NewSignatureVerifierFromSSHAgent() (*SignatureVerifier, error) {
	keyring, err := ssh.SSHAgentClient()
	if err != nil {
		return nil, err
	}
	return NewSignatureVerifierFromKeyring(keyring)
}

// NewSignatureVerifierFromKeyring creates a SignatureVerifier
// which uses public keys from the given keyring to verify signatures
func NewSignatureVerifierFromKeyring(keyring gosshagent.Agent) (*SignatureVerifier, error) {
	keys, err := keyring.List()
	if err != nil {
		return nil, err
	}

	pubkeys := make([]gossh.PublicKey, len(keys))

	for i, k := range keys {
		pubkeys[i] = k
	}

	return &SignatureVerifier{pubkeys}, nil
}

// NewSignatureVerifierFromAuthorizedKeysFile creates a
// SignatureVerifier which uses public keys from the specified
// authorized_keys file to verify signatures
func NewSignatureVerifierFromAuthorizedKeysFile(filepath string) (*SignatureVerifier, error) {
	out, err := ioutil.ReadFile(pkg.ParseFilepath(filepath))
	if err != nil {
		return nil, err
	}

	pubkeys, err := parseAuthorizedKeys(out)
	if err != nil {
		return nil, err
	}
	if len(pubkeys) == 0 {
		return nil, errors.New("no authorized keys found in file")
	}

	return &SignatureVerifier{pubkeys}, nil
}

// Verify verifies that at least one of the signatures in the provided
// SignatureSet is a valid signature of the given data blob.
func (sv *SignatureVerifier) Verify(data []byte, s *SignatureSet) (bool, error) {
	if sv.pubkeys == nil {
		return false, errors.New("signature verifier is uninitialized")
	}

	for _, authKey := range sv.pubkeys {
		// Manually coerce the keys provided by the agent to PublicKeys,
		// since gosshnew does not want to do it for us and agent.Keys are unable to Verify
		key, err := gossh.ParsePublicKey(authKey.Marshal())
		if err != nil {
			log.Errorf("Unable to use SSH key: %v", err)
			continue
		}

		for _, sign := range s.Signatures {
			if err := key.Verify(data, sign); err == nil {
				return true, nil
			}
		}
	}

	return false, nil
}

func parseAuthorizedKeys(in []byte) ([]gossh.PublicKey, error) {
	pubkeys := make([]gossh.PublicKey, 0)
	for len(in) > 0 {
		pubkey, _, _, rest, err := gossh.ParseAuthorizedKey(in)
		if err != nil {
			return nil, err
		}
		in = rest

		pubkeys = append(pubkeys, pubkey)
	}

	return pubkeys, nil
}

func marshal(obj interface{}) ([]byte, error) {
	encoded, err := json.Marshal(obj)
	if err == nil {
		return encoded, nil
	}
	return nil, fmt.Errorf("unable to JSON-serialize object: %s", err)
}
