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

package pkg

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

type keypairFunc func(certPEMBlock, keyPEMBlock []byte) (cert tls.Certificate, err error)

func buildTLSClientConfig(ca, cert, key []byte, parseKeyPair keypairFunc) (*tls.Config, error) {
	if len(cert) == 0 && len(key) == 0 {
		return &tls.Config{InsecureSkipVerify: true}, nil
	}

	tlsCert, err := parseKeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	cfg := tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS10,
	}

	if len(ca) != 0 {
		cp, err := newCertPool(ca)
		if err != nil {
			return nil, err
		}
		cfg.RootCAs = cp
	}

	return &cfg, nil
}

func newCertPool(ca []byte) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	for {
		var block *pem.Block
		block, ca = pem.Decode(ca)
		if block == nil {
			break
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		certPool.AddCert(cert)
	}
	return certPool, nil
}

func ReadTLSConfigFiles(cafile, certfile, keyfile string) (cfg *tls.Config, err error) {
	var ca, cert, key []byte

	if certfile != "" {
		cert, err = ioutil.ReadFile(certfile)
		if err != nil {
			return
		}
	}

	if keyfile != "" {
		key, err = ioutil.ReadFile(keyfile)
		if err != nil {
			return
		}
	}

	if cafile != "" {
		ca, err = ioutil.ReadFile(cafile)
		if err != nil {
			return
		}
	}

	cfg, err = buildTLSClientConfig(ca, cert, key, tls.X509KeyPair)

	return
}
