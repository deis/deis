// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// authenticate authenticates with the remote server. See RFC 4252.
func (c *connection) clientAuthenticate(config *ClientConfig) error {
	// initiate user auth session
	if err := c.transport.writePacket(Marshal(&serviceRequestMsg{serviceUserAuth})); err != nil {
		return err
	}
	packet, err := c.transport.readPacket()
	if err != nil {
		return err
	}
	var serviceAccept serviceAcceptMsg
	if err := Unmarshal(packet, &serviceAccept); err != nil {
		return err
	}

	// during the authentication phase the client first attempts the "none" method
	// then any untried methods suggested by the server.
	tried, remain := make(map[string]bool), make(map[string]bool)
	for auth := AuthMethod(new(noneAuth)); auth != nil; {
		ok, methods, err := auth.auth(c.transport.getSessionID(), config.User, c.transport, config.Rand)
		if err != nil {
			return err
		}
		if ok {
			// success
			return nil
		}
		tried[auth.method()] = true
		delete(remain, auth.method())
		for _, meth := range methods {
			if tried[meth] {
				// if we've tried meth already, skip it.
				continue
			}
			remain[meth] = true
		}
		auth = nil
		for _, a := range config.Auth {
			if remain[a.method()] {
				auth = a
				break
			}
		}
	}
	return fmt.Errorf("ssh: unable to authenticate, attempted methods %v, no supported methods remain", keys(tried))
}

func keys(m map[string]bool) (s []string) {
	for k := range m {
		s = append(s, k)
	}
	return
}

// An AuthMethod represents an instance of an RFC 4252 authentication method.
type AuthMethod interface {
	// auth authenticates user over transport t.
	// Returns true if authentication is successful.
	// If authentication is not successful, a []string of alternative
	// method names is returned.
	auth(session []byte, user string, p packetConn, rand io.Reader) (bool, []string, error)

	// method returns the RFC 4252 method name.
	method() string
}

// "none" authentication, RFC 4252 section 5.2.
type noneAuth int

func (n *noneAuth) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) {
	if err := c.writePacket(Marshal(&userAuthRequestMsg{
		User:    user,
		Service: serviceSSH,
		Method:  "none",
	})); err != nil {
		return false, nil, err
	}

	return handleAuthResponse(c)
}

func (n *noneAuth) method() string {
	return "none"
}

// passwordCallback is an AuthMethod that fetches the password through
// a function call, e.g. by prompting the user.
type passwordCallback func() (password string, err error)

func (cb passwordCallback) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) {
	type passwordAuthMsg struct {
		User     string `sshtype:"50"`
		Service  string
		Method   string
		Reply    bool
		Password string
	}

	pw, err := cb()
	if err != nil {
		return false, nil, err
	}

	if err := c.writePacket(Marshal(&passwordAuthMsg{
		User:     user,
		Service:  serviceSSH,
		Method:   cb.method(),
		Reply:    false,
		Password: pw,
	})); err != nil {
		return false, nil, err
	}

	return handleAuthResponse(c)
}

func (cb passwordCallback) method() string {
	return "password"
}

// Password returns an AuthMethod using the given password.
func Password(secret string) AuthMethod {
	return passwordCallback(func() (string, error) { return secret, nil })
}

// PasswordCallback returns an AuthMethod that uses a callback for
// fetching a password.
func PasswordCallback(prompt func() (secret string, err error)) AuthMethod {
	return passwordCallback(prompt)
}

type publickeyAuthMsg struct {
	User    string `sshtype:"50"`
	Service string
	Method  string
	// HasSig indicates to the receiver packet that the auth request is signed and
	// should be used for authentication of the request.
	HasSig   bool
	Algoname string
	Pubkey   string
	// Sig is defined as []byte so Marshal will exclude it during validateKey
	Sig []byte `ssh:"rest"`
}

// publicKeyCallback is an AuthMethod that uses a set of key
// pairs fors authentication.
type publicKeyCallback func() ([]Signer, error)

func (cb publicKeyCallback) method() string {
	return "publickey"
}

func (cb publicKeyCallback) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) {
	// Authentication is performed in two stages. The first stage sends an
	// enquiry to test if each key is acceptable to the remote. The second
	// stage attempts to authenticate with the valid keys obtained in the
	// first stage.

	// a map of public keys to their index in the keyring
	signers, err := cb()
	if err != nil {
		return false, nil, err
	}
	var validKeys []Signer
	for _, signer := range signers {
		if ok, err := validateKey(signer.PublicKey(), user, c); ok {
			validKeys = append(validKeys, signer)
		} else {
			if err != nil {
				return false, nil, err
			}
		}
	}

	// methods that may continue if this auth is not successful.
	var methods []string
	for _, signer := range validKeys {
		pub := signer.PublicKey()

		pubkey := pub.Marshal()
		sign, err := signer.Sign(rand, buildDataSignedForAuth(session, userAuthRequestMsg{
			User:    user,
			Service: serviceSSH,
			Method:  cb.method(),
		}, []byte(pub.Type()), pubkey))
		if err != nil {
			return false, nil, err
		}

		// manually wrap the serialized signature in a string
		s := Marshal(sign)
		sig := make([]byte, stringLength(len(s)))
		marshalString(sig, s)
		msg := publickeyAuthMsg{
			User:     user,
			Service:  serviceSSH,
			Method:   cb.method(),
			HasSig:   true,
			Algoname: pub.Type(),
			Pubkey:   string(pubkey),
			Sig:      sig,
		}
		p := Marshal(&msg)
		if err := c.writePacket(p); err != nil {
			return false, nil, err
		}
		success, methods, err := handleAuthResponse(c)
		if err != nil {
			return false, nil, err
		}
		if success {
			return success, methods, err
		}
	}
	return false, methods, nil
}

// validateKey validates the key provided is acceptable to the server.
func validateKey(key PublicKey, user string, c packetConn) (bool, error) {
	pubkey := key.Marshal()
	msg := publickeyAuthMsg{
		User:     user,
		Service:  serviceSSH,
		Method:   "publickey",
		HasSig:   false,
		Algoname: key.Type(),
		Pubkey:   string(pubkey),
	}
	if err := c.writePacket(Marshal(&msg)); err != nil {
		return false, err
	}

	return confirmKeyAck(key, c)
}

func confirmKeyAck(key PublicKey, c packetConn) (bool, error) {
	pubkey := key.Marshal()
	algoname := key.Type()

	for {
		packet, err := c.readPacket()
		if err != nil {
			return false, err
		}
		switch packet[0] {
		case msgUserAuthBanner:
			// TODO(gpaul): add callback to present the banner to the user
		case msgUserAuthPubKeyOk:
			msg := userAuthPubKeyOkMsg{}
			if err := Unmarshal(packet, &msg); err != nil {
				return false, err
			}
			if msg.Algo != algoname || !bytes.Equal(msg.PubKey, pubkey) {
				return false, nil
			}
			return true, nil
		case msgUserAuthFailure:
			return false, nil
		default:
			return false, unexpectedMessageError(msgUserAuthSuccess, packet[0])
		}
	}
	panic("unreachable")
}

// PublicKeys returns an AuthMethod that uses the given key
// pairs.
func PublicKeys(signers ...Signer) AuthMethod {
	return publicKeyCallback(func() ([]Signer, error) { return signers, nil })
}

// PublicKeysCallback returns an AuthMethod that runs the given
// function to obtain a list of key pairs.
func PublicKeysCallback(getSigners func() (signers []Signer, err error)) AuthMethod {
	return publicKeyCallback(getSigners)
}

// handleAuthResponse returns whether the preceding authentication request succeeded
// along with a list of remaining authentication methods to try next and
// an error if an unexpected response was received.
func handleAuthResponse(c packetConn) (bool, []string, error) {
	for {
		packet, err := c.readPacket()
		if err != nil {
			return false, nil, err
		}

		switch packet[0] {
		case msgUserAuthBanner:
			// TODO: add callback to present the banner to the user
		case msgUserAuthFailure:
			msg := userAuthFailureMsg{}
			if err := Unmarshal(packet, &msg); err != nil {
				return false, nil, err
			}
			return false, msg.Methods, nil
		case msgUserAuthSuccess:
			return true, nil, nil
		case msgDisconnect:
			return false, nil, io.EOF
		default:
			return false, nil, unexpectedMessageError(msgUserAuthSuccess, packet[0])
		}
	}
	panic("unreachable")
}

// KeyboardInteractiveChallenge should print questions, optionally
// disabling echoing (e.g. for passwords), and return all the answers.
// Challenge may be called multiple times in a single session. After
// successful authentication, the server may send a challenge with no
// questions, for which the user and instruction messages should be
// printed.  RFC 4256 section 3.3 details how the UI should behave for
// both CLI and GUI environments.
type KeyboardInteractiveChallenge func(user, instruction string, questions []string, echos []bool) (answers []string, err error)

// KeyboardInteractive returns a AuthMethod using a prompt/response
// sequence controlled by the server.
func KeyboardInteractive(challenge KeyboardInteractiveChallenge) AuthMethod {
	return challenge
}

func (cb KeyboardInteractiveChallenge) method() string {
	return "keyboard-interactive"
}

func (cb KeyboardInteractiveChallenge) auth(session []byte, user string, c packetConn, rand io.Reader) (bool, []string, error) {
	type initiateMsg struct {
		User       string `sshtype:"50"`
		Service    string
		Method     string
		Language   string
		Submethods string
	}

	if err := c.writePacket(Marshal(&initiateMsg{
		User:    user,
		Service: serviceSSH,
		Method:  "keyboard-interactive",
	})); err != nil {
		return false, nil, err
	}

	for {
		packet, err := c.readPacket()
		if err != nil {
			return false, nil, err
		}

		// like handleAuthResponse, but with less options.
		switch packet[0] {
		case msgUserAuthBanner:
			// TODO: Print banners during userauth.
			continue
		case msgUserAuthInfoRequest:
			// OK
		case msgUserAuthFailure:
			var msg userAuthFailureMsg
			if err := Unmarshal(packet, &msg); err != nil {
				return false, nil, err
			}
			return false, msg.Methods, nil
		case msgUserAuthSuccess:
			return true, nil, nil
		default:
			return false, nil, unexpectedMessageError(msgUserAuthInfoRequest, packet[0])
		}

		var msg userAuthInfoRequestMsg
		if err := Unmarshal(packet, &msg); err != nil {
			return false, nil, err
		}

		// Manually unpack the prompt/echo pairs.
		rest := msg.Prompts
		var prompts []string
		var echos []bool
		for i := 0; i < int(msg.NumPrompts); i++ {
			prompt, r, ok := parseString(rest)
			if !ok || len(r) == 0 {
				return false, nil, errors.New("ssh: prompt format error")
			}
			prompts = append(prompts, string(prompt))
			echos = append(echos, r[0] != 0)
			rest = r[1:]
		}

		if len(rest) != 0 {
			return false, nil, fmt.Errorf("ssh: junk following message %q", rest)
		}

		answers, err := cb(msg.User, msg.Instruction, prompts, echos)
		if err != nil {
			return false, nil, err
		}

		if len(answers) != len(prompts) {
			return false, nil, errors.New("ssh: not enough answers from keyboard-interactive callback")
		}
		responseLength := 1 + 4
		for _, a := range answers {
			responseLength += stringLength(len(a))
		}
		serialized := make([]byte, responseLength)
		p := serialized
		p[0] = msgUserAuthInfoResponse
		p = p[1:]
		p = marshalUint32(p, uint32(len(answers)))
		for _, a := range answers {
			p = marshalString(p, []byte(a))
		}

		if err := c.writePacket(serialized); err != nil {
			return false, nil, err
		}
	}
}
