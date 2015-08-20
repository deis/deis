package sshd

import (
	"testing"

	"github.com/Masterminds/cookoo"
	"golang.org/x/crypto/ssh"
)

// TestAuthKey tests the AuthKey command for authorizing public keys.
func TestAuthKey(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	key, err := sshTestingClientKey()
	if err != nil {
		t.Fatal(err)
	}

	cxt.Put("cookoo.Router", router)
	cxt.Put("authorizedKeys", []string{testingClientPubKey})
	cxt.Put("key", key.PublicKey())
	cxt.Put("metadata", &connMetadata{})

	reg.AddRoute(cookoo.Route{
		Name: "auth", Help: "Authenticate a route.",
		Does: cookoo.Tasks{
			cookoo.Cmd{
				Name: "auth",
				Fn:   AuthKey,
				Using: []cookoo.Param{
					{Name: "key", From: "cxt:key"},
					{Name: "authorizedKeys", From: "cxt:authorizedKeys"},
					{Name: "metadata", From: "cxt:metadata"},
				},
			},
		},
	})
	if err := router.HandleRequest("auth", cxt, true); err != nil {
		t.Fatalf("Failed auth run with %s", err)
	}

	auth := cxt.Get("auth", &ssh.Permissions{}).(*ssh.Permissions)
	if user, ok := auth.Extensions["user"]; !ok {
		t.Errorf("Expected a user, but got none.")
	} else if user != "deis" {
		t.Errorf("Expected user to be 'deis', got '%s'", user)
	}
}

func TestFingerprint(t *testing.T) {
	key, _ := sshTestingClientKey()
	fp := Fingerprint(key.PublicKey())
	if fp != testingClientFingerprint {
		t.Errorf("Expected fingerprint %s to match %s.", fp, testingClientFingerprint)
	}
}
