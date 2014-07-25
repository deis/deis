package sign

import (
	"fmt"

	gossh "github.com/coreos/fleet/third_party/code.google.com/p/gosshnew/ssh"
	"github.com/coreos/fleet/third_party/code.google.com/p/gosshnew/ssh/testdata"
)

var (
	testPrivateKeys map[string]interface{}
	testSigners     map[string]gossh.Signer
	testPublicKeys  map[string]gossh.PublicKey
)

func init() {
	var err error

	n := len(testdata.PEMBytes)
	testPrivateKeys = make(map[string]interface{}, n)
	testSigners = make(map[string]gossh.Signer, n)
	testPublicKeys = make(map[string]gossh.PublicKey, n)
	for t, k := range testdata.PEMBytes {
		testPrivateKeys[t], err = gossh.ParseRawPrivateKey(k)
		if err != nil {
			panic(fmt.Sprintf("Unable to parse test key %s: %v", t, err))
		}
		testSigners[t], err = gossh.NewSignerFromKey(testPrivateKeys[t])
		if err != nil {
			panic(fmt.Sprintf("Unable to create signer for test key %s: %v", t, err))
		}
		testPublicKeys[t] = testSigners[t].PublicKey()
	}

	// Create a cert and sign it for use in tests.
	testCert := &gossh.Certificate{
		Nonce:           []byte{},                       // To pass reflect.DeepEqual after marshal & parse, this must be non-nil
		ValidPrincipals: []string{"gopher1", "gopher2"}, // increases test coverage
		ValidAfter:      0,                              // unix epoch
		ValidBefore:     gossh.CertTimeInfinity,         // The end of currently representable time.
		Reserved:        []byte{},                       // To pass reflect.DeepEqual after marshal & parse, this must be non-nil
		Key:             testPublicKeys["ecdsa"],
		SignatureKey:    testPublicKeys["rsa"],
		Permissions: gossh.Permissions{
			CriticalOptions: map[string]string{},
			Extensions:      map[string]string{},
		},
	}
	testCert.SignCert(testSigners["rsa"])
	testPrivateKeys["cert"] = testPrivateKeys["ecdsa"]
	testSigners["cert"], err = gossh.NewCertSigner(testCert, testSigners["ecdsa"])
	if err != nil {
		panic(fmt.Sprintf("Unable to create certificate signer: %v", err))
	}
}
