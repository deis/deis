package sshd

import (
	"bufio"
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/log"
)

// ParseAuthorizedKeys reads and process an authorized_keys file.
//
// The file is merely parsed into lines, which are then returned in an array.
//
// Params:
// 	- path (string): The path to the authorized_keys file.
//
// Returns:
//  []string of keys.
//
func ParseAuthorizedKeys(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	path := p.Get("path", "~/.ssh/authorized_keys").(string)

	file, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	reader := bufio.NewScanner(file)
	buf := []string{}

	for reader.Scan() {
		data := reader.Text()
		if len(data) > 0 {
			log.Infof(c, "Adding key '%s'", data)
			buf = append(buf, strings.TrimSpace(data))
		}
	}

	return buf, nil

}

// ParseHostKeys parses the host key files.
//
// By default it looks in /etc/ssh for host keys of the patterh ssh_host_{{TYPE}}_key.
//
// Params:
// 	- keytypes ([]string): Key types to parse. Defaults to []string{rsa, dsa, ecdsa}
// 	- enabeV1 (bool): Allow V1 keys. By default this is disabled.
// 	- path (string): Override the lookup pattern. If %s, it will be replaced with the keytype.
//
// Returns:
// 	[]ssh.Signer
func ParseHostKeys(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	hostKeyTypes := p.Get("keytypes", []string{"rsa", "dsa", "ecdsa"}).([]string)
	pathTpl := p.Get("path", "/etc/ssh/ssh_host_%s_key").(string)
	hostKeys := make([]ssh.Signer, 0, len(hostKeyTypes))
	for _, t := range hostKeyTypes {
		path := fmt.Sprintf(pathTpl, t)

		if key, err := ioutil.ReadFile(path); err == nil {
			if hk, err := ssh.ParsePrivateKey(key); err == nil {
				log.Infof(c, "Parsed host key %s.", path)
				hostKeys = append(hostKeys, hk)
			} else {
				log.Errf(c, "Failed to parse host key %s (skipping): %s", path, err)
			}
		}
	}
	if c.Get("enableV1", false).(bool) {
		path := "/etc/ssh/ssh_host_key"
		if key, err := ioutil.ReadFile(path); err != nil {
			log.Errf(c, "Failed to read ssh_host_key")
		} else if hk, err := ssh.ParsePrivateKey(key); err == nil {
			log.Infof(c, "Parsed host key %s.", path)
			hostKeys = append(hostKeys, hk)
		} else {
			log.Errf(c, "Failed to parse host key %s: %s", path, err)
		}
	}
	return hostKeys, nil
}

// AuthKey authenticates based on a public key.
//
// Params:
// 	- metadata (ssh.ConnMetadata)
// 	- key (ssh.PublicKey)
// 	- authorizedKeys ([]string): List of lines from an authorized keys file.
//
// Returns:
// 	*ssh.Permissions
//
func AuthKey(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	meta := p.Get("metadata", nil).(ssh.ConnMetadata)
	key := p.Get("key", nil).(ssh.PublicKey)
	authorized := p.Get("authorizedKeys", []string{}).([]string)

	auth := new(ssh.CertChecker)
	auth.UserKeyFallback = func(meta ssh.ConnMetadata, pk ssh.PublicKey) (*ssh.Permissions, error) {

		// This gives us a string in the form "ssh-rsa LONG_KEY"
		//supplied := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pk)))
		suppliedType := key.Type()
		supplied := key.Marshal()

		for _, allowedKey := range authorized {
			allowed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(allowedKey))
			if err != nil {
				log.Infof(c, "Could not parse authorized key '%q': %s", allowedKey, err)
				continue
			}

			// We use a contstant time compare more as a precaution than anything
			// else. A timing attack here would be very difficult, but... better
			// safe than sorry.
			if allowed.Type() == suppliedType && subtle.ConstantTimeCompare(allowed.Marshal(), supplied) == 1 {
				log.Infof(c, "Key accepted for user %s.", meta.User())
				perm := &ssh.Permissions{
					Extensions: map[string]string{
						//"key":  ssh.MarshalAuthorizedKey(supplied.Marshal()),
						"user": meta.User(),
					},
				}
				return perm, nil
			}
		}

		return nil, fmt.Errorf("No matching keys found.")
	}

	return auth.Authenticate(meta, key)
}

// compareKeys compares to key files and returns true of they match.
func compareKeys(a, b ssh.PublicKey) bool {
	if a.Type() != b.Type() {
		return false
	}
	// The best way to compare just the key seems to be to marshal both and
	// then compare the output byte sequence.
	return subtle.ConstantTimeCompare(a.Marshal(), b.Marshal()) == 1
}

// Start starts an instance of /usr/sbin/sshd.
//
//
func Start(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dargs := []string{"-e", "-D"}

	sshd := exec.Command("/usr/sbin/sshd", dargs...)
	sshd.Stdout = os.Stdout
	sshd.Stderr = os.Stderr

	if err := sshd.Start(); err != nil {
		return 0, err
	}

	return sshd.Process.Pid, nil
}

// Configure creates a new SSH configuration object.
//
// Config sets a PublicKeyCallback handler that forwards public key auth
// requests to the route named "pubkeyAuth".
//
// This assumes certain details about our environment, like the location of the
// host keys. It also provides only key-based authentication.
// ConfigureServerSshConfig
//
// Params:
//
// Returns:
//  An *ssh.ServerConfig
func Configure(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	router := c.Get("cookoo.Router", nil).(*cookoo.Router)

	//hostKeyTypes := []string{"rsa", "dsa", "ecdsa"}
	cfg := &ssh.ServerConfig{
		PublicKeyCallback: func(m ssh.ConnMetadata, k ssh.PublicKey) (*ssh.Permissions, error) {
			c.Put("metadata", m)
			c.Put("key", k)

			pubkeyAuth := c.Get("route.sshd.pubkeyAuth", "pubkeyAuth").(string)
			err := router.HandleRequest(pubkeyAuth, c, true)
			return c.Get("pubkeyAuth", &ssh.Permissions{}).(*ssh.Permissions), err
		},
	}

	return cfg, nil
}

// FingerprintKey fingerprints a key and returns the colon-formatted version
//
// Params:
// 	- key (ssh.PublicKey): The key to fingerprint.
//
// Returns:
// 	- A string representation of the key fingerprint.
func FingerprintKey(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	key := p.Get("key", nil).(ssh.PublicKey)
	return Fingerprint(key), nil
}

// Fingerprint generates a colon-separated fingerprint string from a public key.
func Fingerprint(key ssh.PublicKey) string {
	hash := md5.Sum(key.Marshal())
	buf := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(buf, hash[:])
	// We need this in colon notation:
	fp := make([]byte, len(buf)+15)

	i, j := 0, 0
	for ; i < len(buf); i++ {
		if i > 0 && i%2 == 0 {
			fp[j] = ':'
			j++
		}
		fp[j] = buf[i]
		j++
	}

	return string(fp)
}
