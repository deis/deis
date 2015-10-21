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

package ssh

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	gossh "golang.org/x/crypto/ssh"

	"github.com/coreos/fleet/log"
	"github.com/coreos/fleet/pkg"
)

const (
	DefaultKnownHostsFile = "~/.fleetctl/known_hosts"

	sshDefaultPort = 22  // ssh.h
	sshHashDelim   = "|" // hostfile.h

	warningRemoteHostChanged = `@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The fingerprint for the %v key sent by the remote host is
%v.
Please contact your system administrator.
Add correct host key in %v to get rid of this message.
Host key verification failed.
`
	promptToTrustHost = `The authenticity of host '%v' can't be established.
%v key fingerprint is %v.
Are you sure you want to continue connecting (yes/no)? `
)

// askToTrustHost prompts the user to trust a new key fingerprint while connecting to a host
func askToTrustHost(addr, algo, fingerprint string) bool {
	var ans string

	fmt.Fprintf(os.Stderr, promptToTrustHost, addr, algo, fingerprint)
	fmt.Scanf("%s\n", &ans)

	ans = strings.ToLower(ans)
	if ans != "yes" && ans != "y" {
		return false
	}

	return true
}

var (
	ErrUntrustHost = errors.New("unauthorized host")
	ErrUnmatchKey  = errors.New("host key mismatch")
)

// HostKeyChecker implements the gossh.HostKeyChecker interface
// It is used for key validation during the cryptographic handshake
type HostKeyChecker struct {
	m         HostKeyManager
	trustHost func(addr, algo, fingerprint string) bool
}

// NewHostKeyChecker returns a new HostKeyChecker
func NewHostKeyChecker(m HostKeyManager) *HostKeyChecker {
	return &HostKeyChecker{m, askToTrustHost}
}

// Check is called during the handshake to check the server's public key for
// unexpected changes. The key argument is in SSH wire format. It can be parsed
// using ssh.ParsePublicKey. The address before DNS resolution is passed in the
// addr argument, so the key can also be checked against the hostname.
// It returns any error encountered while checking the public key. A nil return
// value indicates that the key was either successfully verified (against an
// existing known_hosts entry), or accepted by the user as a new key.
func (kc *HostKeyChecker) Check(addr string, remote net.Addr, key gossh.PublicKey) error {
	remoteAddr, err := kc.addrToHostPort(remote.String())
	if err != nil {
		return err
	}

	algoStr := algoString(key.Type())
	keyFingerprintStr := md5String(md5.Sum(key.Marshal()))

	hostKeys, err := kc.m.GetHostKeys()
	_, ok := err.(*os.PathError)
	if err != nil && !ok {
		log.Errorf("Failed to read known_hosts file %v: %v", kc.m.String(), err)
	}

	mismatched := false
	for pattern, keys := range hostKeys {
		if !matchHost(remoteAddr, pattern) {
			continue
		}
		for _, hostKey := range keys {
			// Any matching key is considered a success, irrespective of previous failures
			if hostKey.Type() == key.Type() && bytes.Compare(hostKey.Marshal(), key.Marshal()) == 0 {
				return nil
			}
			// TODO(jonboulle): could be super friendly like the OpenSSH client
			// and note exactly which key failed (file + line number)
			mismatched = true
		}
	}

	if mismatched {
		fmt.Fprintf(os.Stderr, warningRemoteHostChanged, algoStr, keyFingerprintStr, kc.m.String())
		return ErrUnmatchKey
	}

	// If we get this far, we haven't matched on any of the hostname patterns,
	// so it's considered a new key. Prompt the user to trust it.
	if !kc.trustHost(remoteAddr, algoStr, keyFingerprintStr) {
		fmt.Fprintln(os.Stderr, "Host key verification failed.")
		return ErrUntrustHost
	}

	if err := kc.m.PutHostKey(remoteAddr, key); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add the host to the list of known hosts (%v).\n", kc.m)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Warning: Permanently added '%v' (%v) to the list of known hosts.\n", remoteAddr, algoStr)
	return nil
}

// addrToHostPort takes the given address and parses it into a string suitable
// for use in the 'hostnames' field in a known_hosts file.  For more details,
// see the `SSH_KNOWN_HOSTS FILE FORMAT` section of `man 8 sshd`
func (kc *HostKeyChecker) addrToHostPort(a string) (string, error) {
	if !strings.Contains(a, ":") {
		// No port, so return unadulterated
		return a, nil
	}
	host, p, err := net.SplitHostPort(a)
	if err != nil {
		log.Debugf("Unable to parse addr %s: %v", a, err)
		return "", err
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		log.Debugf("Error parsing port %s: %v", p, err)
		return "", err
	}

	// Default port should be omitted from the entry.
	// (see `put_host_port` in openssh/misc.c)
	if port == 0 || port == sshDefaultPort {
		// IPv6 addresses must be enclosed in square brackets
		if strings.Contains(host, ":") {
			host = fmt.Sprintf("[%s]", host)
		}
		return host, nil
	}

	return fmt.Sprintf("[%s]:%d", host, port), nil
}

// HostKeyManager defines an interface for managing "known hosts" keys
type HostKeyManager interface {
	String() string
	// GetHostKeys returns a map from host patterns to a list of PublicKeys
	GetHostKeys() (map[string][]gossh.PublicKey, error)
	// put new host key under management
	PutHostKey(addr string, hostKey gossh.PublicKey) error
}

// HostKeyFile is an implementation of HostKeyManager that saves and loads
// "known hosts" keys from a file
type HostKeyFile struct {
	path string
}

// NewHostKeyFile returns a new HostKeyFile using the given file path
func NewHostKeyFile(path string) *HostKeyFile {
	return &HostKeyFile{pkg.ParseFilepath(path)}
}

func (f *HostKeyFile) String() string {
	return f.path
}

func (f *HostKeyFile) GetHostKeys() (map[string][]gossh.PublicKey, error) {
	in, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	hostKeys := make(map[string][]gossh.PublicKey)
	n := 0
	s := bufio.NewScanner(in)
	for s.Scan() {
		n++
		line := s.Bytes()

		hosts, key, err := parseKnownHostsLine(line)

		if err != nil {
			log.Warningf("%v:%d - %v\n", f.path, n, err)
			continue
		}

		if hosts == "" {
			// Comment/empty line
			continue
		}

		// It is permissible to have several lines for the same host name(s)
		hostKeys[hosts] = append(hostKeys[hosts], key)
	}

	return hostKeys, nil
}

// parseKnownHostsLine parses a line from a known hosts file.  It returns a
// string containing the hosts section of the line, a gossh.PublicKey parsed
// from the line, and any error encountered during the parsing.
func parseKnownHostsLine(line []byte) (string, gossh.PublicKey, error) {

	// Skip any leading whitespace.
	line = bytes.TrimLeft(line, "\t ")

	// Skip comments and empty lines.
	if bytes.HasPrefix(line, []byte("#")) || len(line) == 0 {
		return "", nil, nil
	}

	// Skip markers.
	if bytes.HasPrefix(line, []byte("@")) {
		return "", nil, errors.New("marker functionality not implemented")
	}

	// Find the end of the host name(s) portion.
	end := bytes.IndexAny(line, "\t ")
	if end <= 0 {
		return "", nil, errors.New("bad format (insufficient fields)")
	}
	hosts := string(line[:end])
	keyBytes := line[end+1:]

	// Check for hashed host names.
	if strings.HasPrefix(hosts, sshHashDelim) {
		return "", nil, errors.New("hashed hosts not implemented")
	}

	// Finally, actually try to extract the key.
	key, _, _, _, err := gossh.ParseAuthorizedKey(keyBytes)
	if err != nil {
		return "", nil, fmt.Errorf("error parsing key: %v", err)
	}

	return hosts, key, nil
}

func (f *HostKeyFile) PutHostKey(addr string, hostKey gossh.PublicKey) error {
	// Make necessary directories if needed
	err := os.MkdirAll(path.Dir(f.path), 0700)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(renderHostLine(addr, hostKey))
	if err != nil {
		return err
	}
	return nil
}

func renderHostLine(addr string, key gossh.PublicKey) []byte {
	keyByte := gossh.MarshalAuthorizedKey(key)
	// allocate line space in advance
	length := len(addr) + 1 + len(keyByte)
	line := make([]byte, 0, length)

	w := bytes.NewBuffer(line)
	w.Write([]byte(addr))
	w.WriteByte(' ')
	w.Write(keyByte)
	return w.Bytes()
}

// algoString returns a short-name representation of an algorithm type
func algoString(algo string) string {
	switch algo {
	case gossh.KeyAlgoRSA:
		return "RSA"
	case gossh.KeyAlgoDSA:
		return "DSA"
	case gossh.KeyAlgoECDSA256, gossh.KeyAlgoECDSA384, gossh.KeyAlgoECDSA521:
		return "ECDSA"
	}
	return algo
}

// md5String returns a formatted string representing the given md5Sum in hex
func md5String(md5Sum [16]byte) string {
	md5Str := fmt.Sprintf("% x", md5Sum)
	md5Str = strings.Replace(md5Str, " ", ":", -1)
	return md5Str
}
