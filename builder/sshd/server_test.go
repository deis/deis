package sshd

import (
	"net"
	"testing"
	"time"

	"github.com/Masterminds/cookoo"
	"golang.org/x/crypto/ssh"
)

var testingServerAddr = "127.0.0.1:2244"

// TestServer tests the SSH server.
//
// This listens on the non-standard port 2244 of localhost. This will generate
// an entry in your known_hosts file, and will tie that to the testing key
// used here. It's not recommended that you try to start another SSH server on
// the same port (at a later time) or else you will have key issues that you
// must manually resolve.
func TestServer(t *testing.T) {
	key, err := sshTestingHostKey()
	if err != nil {
		t.Fatal(err)
	}

	cfg := ssh.ServerConfig{
		NoClientAuth: true,
	}
	cfg.AddHostKey(key)

	cxt := runServer(&cfg, t)

	// Give server time to initialize.
	time.Sleep(200 * time.Millisecond)

	// Connect to the server and issue env var set. This should return true.
	client, err := ssh.Dial("tcp", testingServerAddr, &ssh.ClientConfig{})
	if err != nil {
		t.Fatalf("Failed to connect client to local server: %s", err)
	}
	sess, err := client.NewSession()
	if err != nil {
		t.Fatalf("Failed to create client session: %s", err)
	}
	defer sess.Close()

	if err := sess.Setenv("HELLO", "world"); err != nil {
		t.Fatal(err)
	}

	if out, err := sess.Output("ping"); err != nil {
		t.Errorf("Output '%s' Error %s", out, err)
	} else if string(out) != "pong" {
		t.Errorf("Expected 'pong', got '%s'", out)
	}

	// Create a new session because the success of the last one closed the
	// connection.
	sess, err = client.NewSession()
	if err != nil {
		t.Fatalf("Failed to create client session: %s", err)
	}
	if err := sess.Run("illegal"); err == nil {
		t.Fatalf("expected a failed run with command 'illegal'")
	}
	if err := sess.Run("illegal command"); err == nil {
		t.Fatalf("expected a failed run with command 'illegal command'")
	}

	closer := cxt.Get("sshd.Closer", nil).(chan interface{})
	closer <- true
}

// sshTestingHostKey loads the testing key.
func sshTestingHostKey() (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(testingHostKey))
}

func sshTestingClientKey() (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(testingClientKey))
}

func runServer(config *ssh.ServerConfig, t *testing.T) cookoo.Context {
	reg, router, cxt := cookoo.Cookoo()
	cxt.Put(ServerConfig, config)
	cxt.Put(Address, testingServerAddr)
	cxt.Put("cookoo.Router", router)

	reg.AddRoute(cookoo.Route{
		Name: "sshPing",
		Help: "Handles an ssh exec ping.",
		Does: cookoo.Tasks{
			cookoo.Cmd{
				Name: "ping",
				Fn:   Ping,
				Using: []cookoo.Param{
					{Name: "request", From: "cxt:request"},
					{Name: "channel", From: "cxt:channel"},
				},
			},
		},
	})

	go func() {
		if err := Serve(reg, router, cxt); err != nil {
			t.Fatalf("Failed serving with %s", err)
		}
	}()

	return cxt

}

// connMetadata mocks ssh.ConnMetadata for authentication.
type connMetadata struct{}

func (cm *connMetadata) User() string          { return "deis" }
func (cm *connMetadata) SessionID() []byte     { return []byte("1") }
func (cm *connMetadata) ClientVersion() []byte { return []byte("2.3.4") }
func (cm *connMetadata) ServerVersion() []byte { return []byte("2.3.4") }
func (cm *connMetadata) RemoteAddr() net.Addr  { return cm.localhost() }
func (cm *connMetadata) LocalAddr() net.Addr   { return cm.localhost() }
func (cm *connMetadata) localhost() net.Addr {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	return addrs[0]
}

var (
	testingHostKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0xOK/wubqj+e4HNp+yAdK4WJnLZCvcjS2DwaxwF+E968kSeU
27SOqiol7Y0UwLGLpB6rpIBnSqXo70xiMUSrnteKmMejddzfbGkvnyvo0dwE4nDd
vnbz64I25xfjTldb4RtNvpk6ymr0soq0EEYssLmdnt7pIgHT71n9RNtu+RPpRe5n
B2ImVeeEsQBhxFsIkkT21JqBhZQRVpeAAOHwainWpkP2MF2ajYUoirs5qOkPxxaw
Mc4i5CSvmFDkWjqkNt84QH9M9M/ws8qX76nImYOPHiF0KRbxamWsYjvdHJCSckdC
mOM7UtsQs8wC3E0xpuPEI0pNRTHCsgH7+KGxmwIDAQABAoIBAAOQufFS7d8zUeiy
qmCeiz+X8todzgTMppsWcNFZuhp10bOV+pK3ew1uxtM7ZdVXamdsSTPvI0+Ee+nG
3YW9hjSZqXKpNJ6iC3gWUsKaiEU7NS3qACTed4JL4ceHhMRm/1tPDcIhbnfK1LVL
WH1J4ileCUaMt11msIDDgV6vYjF81733O+8kPnh5BaFLIOuPdmAPfsZC2WQfBTka
6F5bhe9mcraQohWOGC/NKBbV9o6Ua2GT5ZJILtyPwfx8ctnQHLfmlTOI7qpRyMCU
1hGwlWxyvZRyY4loZehy0c7DaEWJqWS1AST9AbUcNXciYSt/5pUP76W0L6NzwJdh
C1jIY2ECgYEA+JwlIzhsZRsN0jA3A2qWRt3WGdliujAqDvVj4e8E+QnlTh/MDVKF
x3F+w58DHRKJrH7d1nD1fq2id6vh3Sl7xGHZiztOpolY0xlOt71X+2anX+QTEX5Q
d1jB/zQliUsxzIjqn31dKUlAfoI5XiWrxuP1Py8gZSTnnBl8bkdKZysCgYEA2VnG
+bhBdw/0RJVsleyHBrq0+MnQ80dxj6XatKvniVDqjHQefq088W2ULeI5wVjdMy59
CVnDVS6759pLkWu5br7Agb+NGyVKd3o0CT0Jn6JJj9kq1Wq7iOedJF+GtabVp4gk
efIYECkS7BKe1GFH5vRM8FbyyepRFBCgrH1ep1ECgYEAiRojaO7+6CspThcE379y
LJa+MfcueRuCtkkh0kFsbqLEcHccouQ1nq26iMsyfl/wyM4WLOKSoE/FX1XM85ij
BsQnop8MWs83ywMT5ERpNt1/xGQVF/qfCZJLOiBZ6wMq7W88ZMRQEiqxhJLwbDk+
KCsi3rtwlBbsG6v6cR6jq40CgYAzH4nMvQkw7yC+bQMgdIUCETJ1/kpWnqxYZGN/
8ZtBUjYJGVr+4tKd2u9qp3Z8QuGsozen1mQ6igaKr27s4pC4Osfe/OY8x1Wvqp/I
uIGl+a8h1avcjQFVX1036/wsh/RjNoOV51q/mlmoC20ueT9HVJkwQtNSqPmvJYYV
bFuyMQKBgQCsRVEJ6eqai+Pz4bY2UfBnkU6ZHdySI+fQB/T770p0/SbrYMBxNrPQ
v3+ZZfZMlci4pxBtXqrnoyj4uUoqZtR3ENLz53SN1i0vpT7DtC6gMnEF1UWiaoJ6
6mGH5/bxCg9wpV7qpqR0EbFM/dhQFZmmnirOS8x+00hJvc1HFiuN/A==
-----END RSA PRIVATE KEY-----
`
	testingClientKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAvziJnoiaaVyUGPnyqVC49XLzNRS+TPW63Nw4qovCG8lVbxKG
DIHC64tJrCDiZd0ppEhY+RQDGaPwrMInHnV8IwdS1wX22UTRuXA/oXmHcIxO2zmU
nrjFDlpKm2o+2Xd167ifdV9AiqNBtquO0M882RaGy99LbNPcl9ugAnxo5DVI1jES
l5vYqtiOAnRSvmJn2c+hkJfKXryH7hU4y+blDK5Vz44eSsC7bgG3ZbKfKGR9mlf2
ozVlzMNi2ACZ58vDBxn5WVLb1bPV1LHpicJ00fU3TDRnK3MkwvvAnqp78bNzi+ou
YIAwYSZ41iHNd596LQJchr1vs3Fo8qbgYaLY8wIDAQABAoIBAEJgQL0ME/Vw0mOd
F5OYVqu0vCF30trqDXQu6Wih3L5Cc+p7Vpau0Fds4STjwVK0o4jIKEJFpRHYa2m8
d1HGXFHYb/P9uQMQNXCWOzA0/EOgIJtOcH1sC9MAmpc6GRjps8AgNRHL/55gLyZW
hNuMpEWC4UWRfCAJpq/7554VS1+zWK0vy1GszikROjsZnopLTshMV+/7217tSk4O
1GY9ucNJX5iX3M83pmBOJX0ce8fqxeNnAdQIaAtp+ytm5TRzyaQtTjMlq0oqP8+7
Zx9aZKT11IpbOKBSIc6twRArlV1dT9kEI15zS9hfbWuvguB0zuhbhejS4wmZb9Tt
X8rGL4ECgYEA+MZcRzxpBKL+VNuQ4iSwF3RUYL1FIglJV7AM8UdM2hiNeiKidhD5
kmNXVf9C6XWg3OIHCno7HetBo0WZIPmOQMy4CDGC2bWEnQN+/bf3xsKzbECCLtH+
DALXSztihGGiY2zSoOCwTe7WZjGaF9s4C2rVkhsU/9di4qbapGTaWMECgYEAxMZD
c/sVTTT+/thdcLbBDhAfy6RMQwAy/1IPxNVR4C4O+l/rspbKxvV7JyaErP66g871
dBwrOGMfEsYoOOsUBFaj2/jJZdHvQj9jY/kdsfMBivHzkWEFte09NROOThbq+sgX
5bIPwS+IcVCgcA4We+aBv+rYKdvk05RJ8owPSrMCgYBjz4H6erxPxe1wsl8gvEOC
RYQNBCMWks9ARTwMGeU1o6AvnnG8GPdoyj6iHDYGYNFXjb/xbjUFvfupvCTB3B48
1WYIs4SiQHeiX2K1/PeGYVuHVSJmEo5w1zr1zi+qmVmDtoeTUFKsEeUnP0NpyuRj
gEuLwR3dv9bGxNb4GhaYgQKBgDNQCFL8TMe/ZCeMwIEeByXlqoTuKTznlmTiP15y
ylENcbZ0wP/nNqW/aggBkWOTYYvxsiw/FD42CupYZjDBjIy9EynPrKUyo5PA9+gg
FFBNMD/NbFii1lxkqytmGBvg+hG/kAvD7TvRa2ExR0UxR0e0Cm3Dje8MepV5+/aV
837lAoGBAPcvnrDFWKUy8dlrw05+9esiuZgCrCzZPw5xIxhrnRPcBOBl+QdpMscP
eWVutcVy5Frxl5tTf71WK/YhGPgWBt/CQz73Bf1+CX80CeApWWAqiAr240NED5a0
dBAFNBWp8IdHnQmdp9HKvxEXSK+RgOzPNLrpaRv+FPuiD6OtvhmD
-----END RSA PRIVATE KEY-----`

	testingClientPubKey      = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC/OImeiJppXJQY+fKpULj1cvM1FL5M9brc3Diqi8IbyVVvEoYMgcLri0msIOJl3SmkSFj5FAMZo/CswicedXwjB1LXBfbZRNG5cD+heYdwjE7bOZSeuMUOWkqbaj7Zd3XruJ91X0CKo0G2q47QzzzZFobL30ts09yX26ACfGjkNUjWMRKXm9iq2I4CdFK+YmfZz6GQl8pevIfuFTjL5uUMrlXPjh5KwLtuAbdlsp8oZH2aV/ajNWXMw2LYAJnny8MHGflZUtvVs9XUsemJwnTR9TdMNGcrcyTC+8Ceqnvxs3OL6i5ggDBhJnjWIc13n3otAlyGvW+zcWjypuBhotjz donotuse`
	testingClientFingerprint = `78:b9:21:20:1a:ed:e6:10:05:35:47:da:d4:1f:b6:73`
)
