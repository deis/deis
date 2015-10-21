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
	"io"
	"net"

	gossh "golang.org/x/crypto/ssh"
)

func DialCommand(client *SSHForwardingClient, cmd string) (net.Conn, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}

	err = client.ForwardAgentAuthentication(session)
	if err != nil {
		return nil, err
	}

	err = session.Start(cmd)
	if err != nil {
		return nil, err
	}

	pc := &proxyConn{
		session: session,
		writer:  stdin,
		reader:  stdout,
		errchan: make(chan error),
	}

	go func() {
		if err := session.Wait(); err != nil {
			pc.errchan <- err
		}
		close(pc.errchan)
	}()

	return pc, nil
}

type proxyConn struct {
	session *gossh.Session
	writer  io.WriteCloser
	reader  io.Reader
	errchan chan error

	// proxyConn does not fully implement the net.Conn
	// interface, so we have to embed it here.
	net.Conn
}

func (pc *proxyConn) Read(b []byte) (int, error) {
	n, err := pc.reader.Read(b)
	if err == nil {
		return n, err
	}

	perr := <-pc.errchan
	if perr != nil {
		err = perr
	}

	return n, err
}

func (pc *proxyConn) Write(b []byte) (int, error) {
	return pc.writer.Write(b)
}

func (pc *proxyConn) Close() error {
	pc.session.Signal(gossh.SIGTERM)
	pc.session.Close()
	pc.writer.Close()
	return nil
}
