package etcd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

var (
	validCA = []byte(`-----BEGIN CERTIFICATE-----
MIIFNDCCAx6gAwIBAgIBATALBgkqhkiG9w0BAQUwLTEMMAoGA1UEBhMDVVNBMRAw
DgYDVQQKEwdldGNkLWNhMQswCQYDVQQLEwJDQTAeFw0xNDAzMTMwMjA5MDlaFw0y
NDAzMTMwMjA5MDlaMC0xDDAKBgNVBAYTA1VTQTEQMA4GA1UEChMHZXRjZC1jYTEL
MAkGA1UECxMCQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDdlBlw
Jiakc4C1UpMUvQ+2fttyBMfMLivQgj51atpKd8qIBvpZwz1wtpzdRG0hSYMF0IUk
MfBqyg+T5tt2Lfs3Gx3cYKS7G0HTfmABC7GdG8gNvEVNl/efxqvhis7p7hur765e
J+N2GR4oOOP5Wa8O5flv10cp3ZJLhAguc2CONLzfh/iAYAItFgktGHXJ/AnUhhaj
KWdKlK9Cv71YsRPOiB1hCV+LKfNSqrXPMvQ4sarz3yECIBhpV/KfskJoDyeNMaJd
gabX/S7gUCd2FvuOpGWdSIsDwyJf0tnYmQX5XIQwBZJib/IFMmmoVNYc1bFtYvRH
j0g0Ax4tHeXU/0mglqEcaTuMejnx8jlxZAM8Z94wHLfKbtaP0zFwMXkaM4nmfZqh
vLZwowDGMv9M0VRFEhLGYIc3xQ8G2u8cFAGw1UqTxKhwAdRmrcFaQ38sk4kziy0u
AkpGavS7PKcFjjB/fdDFO/kwGQOthX/oTn9nP3BT+IK2h1A6ATMPI4lVnhb5/KBt
9M/fGgbiU+I9QT0Ilz/LlrcCuzyRXREvIZvoUL77Id+JT3qQxqPn/XMKLN4WEFII
112MFGqCD85JZzNoC4RkZd8kFlR4YJWsS4WqJlWprESr5cCDuLviK+31cnIRF4fJ
mz0gPsVgY7GFEan3JJnL8oRUVzdTPKfPt0atsQIDAQABo2MwYTAOBgNVHQ8BAf8E
BAMCAAQwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUnVlVvktY+zlLpG43nTpG
AWmUkrYwHwYDVR0jBBgwFoAUnVlVvktY+zlLpG43nTpGAWmUkrYwCwYJKoZIhvcN
AQEFA4ICAQAqIcPFux3V4h1N0aGM4fCS/iT50TzDnRb5hwILKbmyA6LFnH4YF7PZ
aA0utDNo1XSRDMpR38HWk0weh5Sfx6f2danaKZHAsea8oVEtdrz16ZMOvoh0CPIM
/hn0CGQOoXDADDNFASuExhhpoyYkDqTVTCQ/zbhZg1mjBljJ+BBzlSgeoE4rUDpn
nuDcmD9LtjpsVQL+J662rd51xV4Z6a7aZLvN9GfO8tYkfCGCD9+fGh1Cpz0IL7qw
VRie+p/XpjoHemswnRhYJ4wn10a1UkVSR++wld6Gvjb9ikyr9xVyU5yrRM55pP2J
VguhzjhTIDE1eDfIMMxv3Qj8+BdVQwtKFD+zQYQcbcjsvjTErlS7oCbM2DVlPnRT
QaCM0q0yorfzc4hmml5P95ngz2xlohavgNMhsYIlcWyq3NVbm7mIXz2pjqa16Iit
vL7WX6OVupv/EOMRx5cVcLqqEaYJmAlNd/CCD8ihDQCwoJ6DJhczPRexrVp+iZHK
SnIUONdXb/g8ungXUGL1jGNQrWuq49clpI5sLWNjMDMFAQo0qu5bLkOIMlK/evCt
gctOjXDvGXCk5h6Adf14q9zDGFdLoxw0/aciUSn9IekdzYPmkYUTifuzkVRsPKzS
nmI4dQvz0rHIh4FBUKWWrJhRWhrv9ty/YFuJXVUHeAwr5nz6RFZ4wQ==
-----END CERTIFICATE-----`)
	validKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIJKAIBAAKCAgEAyNxL6iay1rJz24wE/BDYjEcgSDYYWn7m4uTW/oJRM5GwtpL9
s15FZKZAbmw0cMod3qJkm3cCmJN8s/iKKU++d7XibnkaTD6vQMq//j2ZeGNbRtOC
nI3zrzpbOsz7A3x85bkfExO9OSH+cMGbtwXcMc3bcfU9ETsyBIEbdAMbnHuapIPd
yFjcTqyK/uCwsWH06b6U1zttJc9CLkDZtTqaPT1aFp+z13Tprgs0htoVtQ3Cqksk
D+yJKZQSUtBIaKLyLF2r0pDyibLL0I+92RSAVYCoV7h5jzXa8qWkJArcbKm1XTjp
aIyLamE0wwImncEUFpGIAzkkAhiYj6mFScfqx+DJc8UOp/cdqiHJ3pXzK/lRQxHN
WLx7tVyzIOW9SJg+gobrWFtEYRSdwkFXUEdouJCfE9Q0iWCyEjDg2bsdXGWlKEi/
xJKwuf/DzlmZj/JyVzugOMK2Qxxd9P6lqaPk+T77AOnAAX19Y5HE8TwVxitajmfK
06E8aayds3N87mTcUoDN9p843D1IJ+efTIHZdB0eHOCXk2RrHm1psTFppM//wVeH
lGhh6gqc0UB392CMcrLwwtl3+M9gJZPAJS0V6e/5LGrXcQLcnPsvPjFgnOjdGGyP
c47/nswgakfprtT+U29B3mzxc93TnSKYgt5FPEMjBGoMPLucZYmbOAMcHTcCAwEA
AQKCAgBS1vCESKOXgo/f61ae8v+skyUQQyc2I4Jr739wBiUhRKQCGIuDr4ylHyAR
qpTSM7mv+X/O0n2CmcljnEy3Dwl568zQTSf4bB3xde1LGPKzwR6DDnaexLjM+x9n
F+UqoewM/pV/U7PF3WxH6sGi8UrIS6OG02L1OVm+m9TLuwBnQF8eHLiaiXOLCwRk
bBzTe5f70zslrX+tiVY9J0fiw6GbQjNmg0UzxicePcbTGxy6yEsR2t2rp51GRahs
+TPz28hPXe6gcGFnQxNmF/JvllH7cY18aDvSQZ7kVkZlCwmv0ypWoUM6eESDgkW1
a6yrgVccm7bhxW5BYw2AqqSrMkV0oMcCUjh2rYvex7w6dM374Ok3DD/dXjTHLNV5
+0tHMxXUiCKwe7hVEg+iGD4E1jap5n5c4RzpEtAXsGEK5WUBksHi9qOBv+lubjZn
Kcfbos+BbnmUCU3MmU48EZwyFQIu9djkLXfJV2Cbbg9HmkrIOYgi4tFjoBKeQLE4
6GCucMWnNfMO7Kq/z7c+7sfWOAA55pu0Ojel8VH6US+Y/1mEuSUhQudrJn8GxAmc
4t+C2Ie1Q1bK3iJbd0NUqtlwd9xI9wQgCbaxfQceUmBBjuTUu3YFctZ7Jia7h18I
gZ3wsKfySDhW29XTFvnT3FUpc+AN9Pv4sB7uobm6qOBV8/AdKQKCAQEA1zwIuJki
bSgXxsD4cfKgQsyIk0eMj8bDOlf/A8AFursXliH3rRASoixXNgzWrMhaEIE2BeeT
InE13YCUjNCKoz8oZJqKYpjh3o/diZf1vCo6m/YUSR+4amynWE4FEAa58Og2WCJ3
Nx8/IMpmch2VZ+hSQuNr5uvpH84+eZADQ1GB6ypzqxb5HjIEeryLJecDQGe4ophd
JCo3loezq/K0XJQI8GTBe2GQPjXSmLMZKksyZoWEXAaC1Q+sdJWZvBpm3GfVQbXu
q7wyqTMknVIlEOy0sHxstsbayysSFFQ/fcgKjyQb8f4efOkyQg8mH5vQOZghbHJ+
7I8wVSSBt+bE2wKCAQEA7udRoo2NIoIpJH+2+SPqJJVq1gw/FHMM4oXNZp+AAjR1
hTWcIzIXleMyDATl5ZFzZIY1U2JMifS5u2R7fDZEu9vfZk4e6BJUJn+5/ahjYFU8
m8WV4rFWR6XN0SZxPb43Mn6OO7EoMqr8InRufiN4LwIqnPqDm2D9Fdijb9QFJ2UG
QLKNnIkLTcUfx1RYP4T48CHkeZdxV8Cp49SzSSV8PbhIVBx32bm/yO6nLHoro7Wl
YqXGW0wItf2BUA5a5eYNO0ezVkOkTp2aj/p9i+0rqbsYa480hzlnOzYI5F72Z8V2
iPltUAeQn53Vg1azySa1x8/0Xp5nVsgQSh18CH3p1QKCAQBxZv4pVPXgkXlFjTLZ
xr5Ns7pZ7x7OOiluuiJw9WGPazgYMDlxA8DtlXM11Tneu4lInOu73LGXOhLpa+/Y
6Z/CN2qu5wX2wRpwy1gsQNaGl7FdryAtDvt5h1n8ms7sDL83gQHxGee6MUpvmnSz
t4aawrtk5rJZbv7bdS1Rm2E8vNs47psXD/mdwTi++kxOYhNCgeO0N5cLkPrM4x71
f+ErzguPrWaL/XGkdXNKZULjF8+sWLjOS9fvLlzs6E2h4D9F7addAeCIt5XxtDKc
eUVyT2U8f7I/8zIgTccu0tzJBvcZSCs5K20g3zVNvPGXQd9KGS+zFfht51vN4HhA
TuR1AoIBAGuQBKZeexP1bJa9VeF4dRxBldeHrgMEBeIbgi5ZU+YqPltaltEV5Z6b
q1XUArpIsZ6p+mpvkKxwXgtsI1j6ihnW1g+Wzr2IOxEWYuQ9I3klB2PPIzvswj8B
/NfVKhk1gl6esmVXzxR4/Yp5x6HNUHhBznPdKtITaf+jCXr5B9UD3DvW6IF5Bnje
bv9tD0qSEQ71A4xnTiXHXfZxNsOROA4F4bLVGnUR97J9GRGic/GCgFMY9mT2p9lg
qQ8lV3G5EW4GS01kqR6oQQXgLxSIFSeXUFhlIq5bfwoeuwQvaVuxgTwMqVXmAgyL
oK1ApTPE1QWAsLLFORvOed8UxVqBbn0CggEBALfr/wheXCKLdzFzm03sO1i9qVz2
vnpxzexXW3V/TtM6Dff2ojgkDC+CVximtAiLA/Wj60hXnQxw53g5VVT5rESx0J3c
pq+azbi1eWzFeOrqJvKQhMfYc0nli7YuGnPkKzeepJJtWZHYkAjL4QZAn1jt0RqV
DQmlGPGiOuGP8uh59c23pbjgh4eSJnvhOT2BFKhKZpBdTBYeiQiZBqIyme8rNTFr
NmpBxtUr77tccVTrcWWhhViG36UNpetAP7b5QCHScIXZJXrEqyK5HaePqi5UMH8o
alSz6s2REG/xP7x54574TvRG/3cIamv1AfZAOjin7BwhlSLhPl2eeh4Cgas=
-----END RSA PRIVATE KEY-----`)
	validCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFWzCCA0WgAwIBAgIBAjALBgkqhkiG9w0BAQUwLTEMMAoGA1UEBhMDVVNBMRAw
DgYDVQQKEwdldGNkLWNhMQswCQYDVQQLEwJDQTAeFw0xNDAzMTMwMjA5MjJaFw0y
NDAzMTMwMjA5MjJaMEUxDDAKBgNVBAYTA1VTQTEQMA4GA1UEChMHZXRjZC1jYTEP
MA0GA1UECxMGc2VydmVyMRIwEAYDVQQDEwkxMjcuMC4wLjEwggIiMA0GCSqGSIb3
DQEBAQUAA4ICDwAwggIKAoICAQDI3EvqJrLWsnPbjAT8ENiMRyBINhhafubi5Nb+
glEzkbC2kv2zXkVkpkBubDRwyh3eomSbdwKYk3yz+IopT753teJueRpMPq9Ayr/+
PZl4Y1tG04KcjfOvOls6zPsDfHzluR8TE705If5wwZu3Bdwxzdtx9T0ROzIEgRt0
Axuce5qkg93IWNxOrIr+4LCxYfTpvpTXO20lz0IuQNm1Opo9PVoWn7PXdOmuCzSG
2hW1DcKqSyQP7IkplBJS0EhoovIsXavSkPKJssvQj73ZFIBVgKhXuHmPNdrypaQk
CtxsqbVdOOlojItqYTTDAiadwRQWkYgDOSQCGJiPqYVJx+rH4MlzxQ6n9x2qIcne
lfMr+VFDEc1YvHu1XLMg5b1ImD6ChutYW0RhFJ3CQVdQR2i4kJ8T1DSJYLISMODZ
ux1cZaUoSL/EkrC5/8POWZmP8nJXO6A4wrZDHF30/qWpo+T5PvsA6cABfX1jkcTx
PBXGK1qOZ8rToTxprJ2zc3zuZNxSgM32nzjcPUgn559Mgdl0HR4c4JeTZGsebWmx
MWmkz//BV4eUaGHqCpzRQHf3YIxysvDC2Xf4z2Alk8AlLRXp7/ksatdxAtyc+y8+
MWCc6N0YbI9zjv+ezCBqR+mu1P5Tb0HebPFz3dOdIpiC3kU8QyMEagw8u5xliZs4
AxwdNwIDAQABo3IwcDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYD
VR0OBBYEFD6UrVN8uolWz6et79jVeZetjd4XMB8GA1UdIwQYMBaAFJ1ZVb5LWPs5
S6RuN506RgFplJK2MA8GA1UdEQQIMAaHBH8AAAEwCwYJKoZIhvcNAQEFA4ICAQCo
sKn1Rjx0tIVWAZAZB4lCWvkQDp/txnb5zzQUlKhIW2o98IklASmOYYyZbE2PXlda
/n8TwKIzWgIoNh5AcgLWhtASrnZdGFXY88n5jGk6CVZ1+Dl+IX99h+r+YHQzf1jU
BjGrZHGv3pPjwhFGDS99lM/TEBk/eLI2Kx5laL+nWMTwa8M1OwSIh6ZxYPVlWUqb
rurk5l/YqW+UkYIXIQhe6LwtB7tBjr6nDIWBfHQ7uN8IdB8VIAF6lejr22VmERTW
j+zJ5eTzuQN1f0s930mEm8pW7KgGxlEqrUlSJtxlMFCv6ZHZk1Y4yEiOCBKlPNme
X3B+lhj//PH3gLNm3+ZRr5ena3k+wL9Dd3d3GDCIx0ERQyrGS/rJpqNPI+8ZQlG0
nrFlm7aP6UznESQnJoSFbydiD0EZ4hXSdmDdXQkTklRpeXfMcrYBGN7JrGZOZ2T2
WtXBMx2bgPeEH50KRrwUMFe122bchh0Fr+hGvNK2Q9/gRyQPiYHq6vSF4GzorzLb
aDuWA9JRH8/c0z8tMvJ7KjmmmIxd39WWGZqiBrGQR7utOJjpQl+HCsDIQM6yZ/Bu
RpwKj2yBz0OQg4tWbtqUuFkRMTkCR6vo3PadgO1VWokM7UFUXlScnYswcM5EwnzJ
/IsYJ2s1V706QVUzAGIbi3+wYi3enk7JfYoGIqa2oA==
-----END CERTIFICATE-----`)
	corruptedKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
	Corrupted
-----END RSA PRIVATE KEY-----`)
	corruptedCert = []byte(`-----BEGIN CERTIFICATE-----
	Corrupted
-----END CERTIFICATE-----`)
)

// Spot-check NewClient can identify good and bad endpoints
func TestNewClient(t *testing.T) {
	tests := []struct {
		endpoints []string
		pass      bool
	}{
		// these should result in the default endpoint being used
		{[]string{}, true},
		{nil, true},

		// simplest good endpoint, just a scheme and IP
		{[]string{"http://192.0.2.3"}, true},

		// multiple valid values
		{[]string{"http://192.0.2.3", "http://192.0.2.4"}, true},

		// completely invalid URL
		{[]string{"://"}, false},

		// bogus endpoint filtered by our own logic
		{[]string{"boots://pants"}, false},

		// good endpoint followed by a bogus endpoint
		{[]string{"http://192.0.2.3", "boots://pants"}, false},
	}

	for i, tt := range tests {
		_, err := NewClient(tt.endpoints, http.Transport{}, time.Second)
		if tt.pass != (err == nil) {
			t.Errorf("case %d %v: expected to pass=%t, err=%v", i, tt.endpoints, tt.pass, err)
		}
	}
}

// client.SetDefaultPath should only overwrite the path if it is unset
func TestSetDefaultPath(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"http://example.com", "http://example.com/"},
		{"http://example.com/", "http://example.com/"},
		{"http://example.com/foo", "http://example.com/foo"},
	}

	for i, tt := range tests {
		u, _ := url.Parse(tt.in)
		if tt.in != u.String() {
			t.Errorf("case %d: url.Parse modified the URL before we could test it", i)
			continue
		}

		setDefaultPath(u)
		if tt.out != u.String() {
			t.Errorf("case %d: expected output of %s did not match actual value %s", i, tt.out, u.String())
		}
	}
}

// Enumerate the many permutations of an endpoint, asserting whether or
// not they should be acceptable
func TestFilterURL(t *testing.T) {
	tests := []struct {
		endpoint string
		pass     bool
	}{
		// IP & port
		{"http://192.0.2.3:4001/", true},

		// trailing slash
		{"http://192.0.2.3/", true},

		// hostname
		{"http://example.com/", true},

		// https scheme
		{"https://192.0.2.3:4002/", true},

		// no host info
		{"http:///foo/bar", false},

		// empty path
		{"http://192.0.2.3", false},

		// custom path
		{"http://192.0.2.3/foo/bar", false},

		// custom query params
		{"http://192.0.2.3/?foo=bar", false},

		// no scheme
		{"192.0.2.3:4002/", false},

		// non-http scheme
		{"boots://192.0.2.3:4002/", false},

		// no slash after scheme (url.URL.Opaque)
		{"http:192.0.2.3/", false},

		// user info
		{"http://elroy@192.0.2.3/", false},

		// fragment
		{"http://192.0.2.3/#foo", false},
	}

	for i, tt := range tests {
		u, _ := url.Parse(tt.endpoint)
		if tt.endpoint != u.String() {
			t.Errorf("case %d: url.Parse modified the URL before we could test it", i)
			continue
		}

		err := filterURL(u)

		if tt.pass != (err == nil) {
			t.Errorf("case %d %v: expected to pass=%t, err=%v", i, tt.endpoint, tt.pass, err)
		}
	}
}

// Ensure the channel passed into c.resolve is actually wired up
func TestClientCancel(t *testing.T) {
	act := Get{Key: "/foo"}
	c, err := NewClient(nil, http.Transport{}, time.Second)
	if err != nil {
		t.Fatalf("Failed building Client: %v", err)
	}

	cancel := make(chan struct{})
	sentinel := make(chan struct{}, 2)

	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		<-cancel
		sentinel <- struct{}{}
		return nil, nil, errors.New("Cancelled")
	}

	go func() {
		c.resolve(&act, rf, cancel)
		sentinel <- struct{}{}
	}()

	select {
	case <-sentinel:
		t.Fatalf("sentinel should not be ready")
	default:
	}

	close(cancel)

	for i := 0; i < 2; i++ {
		select {
		case <-sentinel:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("timed out waiting for sentinel value")
		}
	}
}

type clientStep struct {
	method string
	url    string

	resp http.Response
}

func assertClientSteps(t *testing.T, c *client, act Action, steps []clientStep, expectSuccess bool) {
	idx := 0
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		if idx >= len(steps) {
			t.Fatalf("Received too many requests")
		}
		step := steps[idx]
		idx = idx + 1

		if step.method != req.Method {
			t.Fatalf("step %d: request method is %s, expected %s", idx, req.Method, step.method)
		}

		if step.url != req.URL.String() {
			t.Fatalf("step %d: request URL is %s, expected %s", idx, req.URL, step.url)
		}

		var body []byte
		if step.resp.Body != nil {
			var err error
			body, err = ioutil.ReadAll(step.resp.Body)
			if err != nil {
				t.Fatalf("step %d: failed preparing body: %v", idx, err)
			}
		}

		return &step.resp, body, nil
	}

	_, err := c.resolve(act, rf, make(chan struct{}))
	if expectSuccess != (err == nil) {
		t.Fatalf("expected to pass=%t, err=%v", expectSuccess, err)
	}
}

// Follow all redirects, using the full Location header regardless of how crazy it seems
func TestClientRedirectsFollowed(t *testing.T) {
	steps := []clientStep{
		{
			"GET", "http://192.0.2.1:4001/v2/keys/foo?consistent=true&recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://192.0.2.2:4001/v2/keys/foo?recursive=false&sorted=false"},
				},
			},
		},
		{
			"GET", "http://192.0.2.2:4001/v2/keys/foo?recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://192.0.2.3:4002/pants?recursive=true"},
				},
			},
		},
		{
			"GET", "http://192.0.2.3:4002/pants?recursive=true",
			http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"X-Etcd-Index": {"123"}},
				Body:       ioutil.NopCloser(strings.NewReader("{}")),
			},
		},
	}

	c, err := NewClient([]string{"http://192.0.2.1:4001"}, http.Transport{}, time.Second)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	act := &Get{Key: "/foo"}
	assertClientSteps(t, c, act, steps, true)
}

// Follow a redirect to a failing node, then fall back to the healthy second endpoint
func TestClientRedirectsAndAlternateEndpoints(t *testing.T) {
	steps := []clientStep{
		{
			"GET", "http://192.0.2.1:4001/v2/keys/foo?consistent=true&recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://192.0.2.5:4001/v2/keys/foo?recursive=true"},
				},
			},
		},
		{
			"GET", "http://192.0.2.5:4001/v2/keys/foo?recursive=true",
			http.Response{
				StatusCode: http.StatusGatewayTimeout,
			},
		},
		{
			"GET", "http://192.0.2.2:4002/v2/keys/foo?consistent=true&recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"X-Etcd-Index": {"123"}},
				Body:       ioutil.NopCloser(strings.NewReader("{}")),
			},
		},
	}

	c, err := NewClient([]string{"http://192.0.2.1:4001", "http://192.0.2.2:4002"}, http.Transport{}, time.Second)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	act := &Get{Key: "/foo"}
	assertClientSteps(t, c, act, steps, true)
}

func TestClientRedirectOverLimit(t *testing.T) {
	reqCount := 0
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		reqCount = reqCount + 1

		if reqCount > 10 {
			t.Fatalf("c.resolve made %d requests, expected max of 10", reqCount)
		}

		resp := http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Header: http.Header{
				"Location": {"http://127.0.0.1:4001/"},
			},
		}

		return &resp, []byte{}, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req != nil || err != nil {
		t.Errorf("Expected nil resp and nil err, got resp=%v and err=%v", req, err)
	}

	if reqCount != 10 {
		t.Fatalf("c.resolve should have made 10 responses, got %d", reqCount)
	}
}

func TestClientRedirectMax(t *testing.T) {
	count := 0
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		var resp http.Response
		var body []byte

		count = count + 1

		if count == 10 {
			resp = http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Etcd-Index": {"123"},
				},
			}
			body = []byte("{}")
		} else {
			resp = http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://127.0.0.1:4001/"},
				},
			}
		}

		return &resp, body, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req == nil || err != nil {
		t.Errorf("Expected non-nil resp and nil err, got resp=%v and err=%v", req, err)
	}
}

func TestClientRequestFuncError(t *testing.T) {
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		return nil, nil, errors.New("bogus error")
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req != nil {
		t.Errorf("Expected req=nil, got %v", nil)
	}
	if err != nil {
		t.Errorf("Expected err=nil, got %v", err)
	}
}

func TestClientRedirectNowhere(t *testing.T) {
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		resp := http.Response{StatusCode: http.StatusTemporaryRedirect}
		return &resp, []byte{}, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req != nil {
		t.Errorf("Expected req=nil, got %v", nil)
	}
	if err != nil {
		t.Errorf("Expected err=nil, got %v", err)
	}
}

func newTestingRequestAndClient(t *testing.T, handler http.Handler) (*client, *http.Request) {
	ts := httptest.NewServer(handler)
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	c, err := NewClient(nil, http.Transport{}, time.Second)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	return c, req
}

func TestGoodRequestHTTP(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing")
	})
	c, req := newTestingRequestAndClient(t, h)

	cancel := make(chan struct{})
	resp, body, err := c.requestHTTP(req, cancel)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	} else {
		// ensure the body was closed
		var b []byte
		if n, err := resp.Body.Read(b); n != 0 || err == nil {
			t.Errorf("resp.Body.Read() returned unexpectedly: want (0, err), got (%d, %v)", n, err)
		}
	}
	if string(body) != "testing" {
		t.Errorf("unexpected body: got %q, want %q", body, "testing")
	}
}

// transport that returns a nil Response and nil error
type nilNilTransport struct{}

func (n *nilNilTransport) RoundTrip(req *http.Request) (*http.Response,
	error) {
	return nil, nil
}
func (n *nilNilTransport) CancelRequest(req *http.Request) {}

// Ensure that any request that somehow returns (nil, nil) propagates an actual error
func TestNilNilRequestHTTP(t *testing.T) {
	c := &client{[]url.URL{}, &nilNilTransport{}, time.Second}
	cancel := make(chan struct{})
	resp, body, err := c.requestHTTP(nil, cancel)
	if err == nil {
		t.Error("unexpected nil error")
	} else if err.Error() != "nil error and nil response" {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("unexpected response: got %q, want %q", resp, nil)
	}
	if body != nil {
		t.Errorf("unexpected body: got %q, want %q", body, nil)
	}
}

// Simple implementation of ReadCloser to serve as response.Body
type rc struct{}

func (r *rc) Read(p []byte) (n int, err error) { return 0, nil }
func (r *rc) Close() error                     { return nil }

// transport that returns a non-nil Response and non-nil error
type respAndErrTransport struct{}

func (r *respAndErrTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: &rc{},
	}, errors.New("some error")
}
func (r *respAndErrTransport) CancelRequest(req *http.Request) {}

// Ensure that the body of a response is closed even when an error is returned
func TestRespAndErrRequestHTTP(t *testing.T) {
	c := &client{[]url.URL{}, &respAndErrTransport{}, time.Second}
	cancel := make(chan struct{})
	resp, body, err := c.requestHTTP(nil, cancel)
	if err == nil {
		t.Error("unexpected nil error")
	} else if err.Error() == "cancelled" {
		t.Error("unexpected error, should not be cancelled")
	}
	if resp != nil {
		t.Errorf("unexpected response: got %q, want %q", resp, nil)
	}
	if body != nil {
		t.Errorf("unexpected body: got %q, want %q", body, nil)
	}
}

func TestCancelledRequestHTTP(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Hour)
	})
	c, req := newTestingRequestAndClient(t, h)

	cancel := make(chan struct{})
	close(cancel)
	resp, body, err := c.requestHTTP(req, cancel)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if err.Error() != "cancelled" {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("unexpected response: got %q, want %q", resp, nil)
	}
	if body != nil {
		t.Errorf("unexpected body: got %q, want %q", body, nil)
	}
}

func TestBuildTLSClientConfigNoCertificate(t *testing.T) {
	config, err := buildTLSClientConfig([]byte{}, []byte{}, []byte{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !config.InsecureSkipVerify {
		t.Errorf("insecureSkipVerify not set")
	}
	if len(config.Certificates) != 0 {
		t.Errorf("unexpected certificates")
	}
	if config.RootCAs != nil {
		t.Errorf("unexpected root CA")
	}
}

func TestBuildTLSClientConfigWithValidCertificateAndWithCA(t *testing.T) {
	config, err := buildTLSClientConfig(validCA, validCert, validKey)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if config.InsecureSkipVerify {
		t.Errorf("insecureSkipVerify should not be set")
	}
	if len(config.Certificates) == 0 {
		t.Errorf("missing certificates")
	}
	if config.RootCAs == nil {
		t.Errorf("missing root CA")
	}
}

func TestBuildTLSClientConfigWithValidCertificateAndWithoutCA(t *testing.T) {
	config, err := buildTLSClientConfig([]byte{}, validCert, validKey)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if config.InsecureSkipVerify {
		t.Errorf("insecureSkipVerify should not be set")
	}
	if len(config.Certificates) == 0 {
		t.Errorf("missing certificates")
	}
	if config.RootCAs != nil {
		t.Errorf("unexpected root CA")
	}
}

func TestBuildTLSClientConfigWithOnlyKeyfileIsAnError(t *testing.T) {
	config, err := buildTLSClientConfig([]byte{}, []byte{}, corruptedKey)
	if err == nil {
		t.Errorf("error expected")
	}
	if config != nil {
		t.Errorf("config should be nil")
	}
}

func TestBuildTLSClientConfigWithOnlyCertfileIsAnError(t *testing.T) {
	config, err := buildTLSClientConfig([]byte{}, validCert, []byte{})
	if err == nil {
		t.Errorf("error expected")
	}
	if config != nil {
		t.Errorf("config should be nil")
	}
}

func TestBuildTLSClientConfigWithCorruptedCertificate(t *testing.T) {
	config, err := buildTLSClientConfig([]byte{}, corruptedCert, corruptedKey)
	if err == nil {
		t.Errorf("error expected")
	}
	if config != nil {
		t.Errorf("config should be nil")
	}
}
