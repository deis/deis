// Package hawk implements the Hawk HTTP authentication scheme.
package hawk

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Now is a func() time.Time that is used by the package to get the current time.
var Now = time.Now

// MaxTimestampSkew is the maximum ±skew that a request timestamp can have without returning ErrTimestampSkew.
var MaxTimestampSkew = time.Minute

var (
	ErrNoAuth             = AuthError("no Authorization header or bewit parameter found")
	ErrReplay             = AuthError("request nonce is being replayed")
	ErrInvalidMAC         = AuthError("invalid MAC")
	ErrBewitExpired       = AuthError("bewit expired")
	ErrTimestampSkew      = AuthError("timestamp skew too high")
	ErrMissingServerAuth  = AuthError("missing Server-Authentication header")
	ErrInvalidBewitMethod = AuthError("bewit only allows HEAD and GET requests")
)

type AuthError string

func (e AuthError) Error() string { return "hawk: " + string(e) }

type CredentialErrorType int

const (
	UnknownID CredentialErrorType = iota
	UnknownApp
	IDAppMismatch
)

func (t CredentialErrorType) String() string {
	switch t {
	case UnknownApp:
		return "unknown app"
	case IDAppMismatch:
		return "id/app mismatch"
	}
	return "unknown id"
}

// CredentialError is returned by a CredentialsLookupFunc when the provided credentials
// ID is invalid.
type CredentialError struct {
	Type        CredentialErrorType
	Credentials *Credentials
}

func (e *CredentialError) Error() string {
	return fmt.Sprintf("hawk: credential error with id %s and app %s: %s", e.Credentials.ID, e.Credentials.App, e.Type)
}

type Credentials struct {
	ID   string
	Key  string
	Hash func() hash.Hash

	// Data may be set in a CredentialsLookupFunc to correlate the credentials
	// with an internal data record.
	Data interface{}

	App      string
	Delegate string
}

func (creds *Credentials) MAC() hash.Hash { return hmac.New(creds.Hash, []byte(creds.Key)) }

type AuthType int

const (
	AuthHeader AuthType = iota
	AuthResponse
	AuthBewit
)

func (a AuthType) String() string {
	switch a {
	case AuthResponse:
		return "response"
	case AuthBewit:
		return "bewit"
	default:
		return "header"
	}
}

// A CredentialsLookupFunc is called by NewAuthFromRequest after parsing the
// request auth. The Credentials will never be nil and ID will always be set.
// App and Delegate will be set if provided in the request. This function must
// look up the corresponding Key and Hash and set them on the provided
// Credentials. If the Key/Hash are found and the App/Delegate are valid (if
// provided) the error should be nil. If the Key or App could not be found or
// the App does not match the ID, then a CredentialError must be returned.
// Errors will propagate to the caller of NewAuthFromRequest, so internal errors
// may be returned.
type CredentialsLookupFunc func(*Credentials) error

// A NonceCheckFunc is called by NewAuthFromRequest and should make sure that
// the provided nonce is unique within the context of the provided time.Time and
// Credentials. It should return false if the nonce is being replayed.
type NonceCheckFunc func(string, time.Time, *Credentials) bool

type AuthFormatError struct {
	Field string
	Err   string
}

func (e AuthFormatError) Error() string { return "hawk: invalid " + e.Field + ", " + e.Err }

// ParseRequestHeader parses a Hawk header (provided in the Authorization
// HTTP header) and populates an Auth. If an error is returned it will always be
// of type AuthFormatError.
func ParseRequestHeader(header string) (*Auth, error) {
	auth := &Auth{ActualTimestamp: Now()}
	err := auth.ParseHeader(header, AuthHeader)
	if err != nil {
		return nil, err
	}

	if auth.Credentials.ID == "" {
		return nil, AuthFormatError{"id", "missing or empty"}
	}
	if auth.Timestamp.IsZero() {
		return nil, AuthFormatError{"ts", "missing, empty, or zero"}
	}
	if auth.Nonce == "" {
		return nil, AuthFormatError{"nonce", "missing or empty"}
	}
	auth.ReqHash = true

	return auth, nil
}

// ParseBewit parses a bewit token provided in a URL parameter and populates an
// Auth. If an error is returned it will always be of type AuthFormatError.
func ParseBewit(bewit string) (*Auth, error) {
	if len(bewit)%4 != 0 {
		bewit += strings.Repeat("=", 4-len(bewit)%4)
	}
	decoded, err := base64.URLEncoding.DecodeString(bewit)
	if err != nil {
		return nil, AuthFormatError{"bewit", "malformed base64 encoding"}
	}
	components := bytes.SplitN(decoded, []byte(`\`), 4)
	if len(components) != 4 {
		return nil, AuthFormatError{"bewit", "missing components"}
	}

	auth := &Auth{
		Credentials:     Credentials{ID: string(components[0])},
		Ext:             string(components[3]),
		Method:          "GET",
		ActualTimestamp: Now(),
		IsBewit:         true,
	}

	ts, err := strconv.ParseInt(string(components[1]), 10, 64)
	if err != nil {
		return nil, AuthFormatError{"ts", "not an integer"}
	}
	auth.Timestamp = time.Unix(ts, 0)

	auth.MAC = make([]byte, base64.StdEncoding.DecodedLen(len(components[2])))
	n, err := base64.StdEncoding.Decode(auth.MAC, components[2])
	if err != nil {
		return nil, AuthFormatError{"mac", "malformed base64 encoding"}
	}
	auth.MAC = auth.MAC[:n]

	return auth, nil
}

// NewAuthFromRequest parses a request containing an Authorization header or
// bewit parameter and populates an Auth. If creds is not nil it will be called
// to look up the associated credentials. If nonce is not nil it will be called
// to make sure the nonce is not replayed.
//
// If the request does not contain a bewit or Authorization header, ErrNoAuth is
// returned. If the request contains a bewit and it is not a GET or HEAD
// request, ErrInvalidBewitMethod is returned. If there is an error parsing the
// provided auth details, an AuthFormatError will be returned. If creds returns
// an error, it will be returned. If nonce returns false, ErrReplay will be
// returned.
func NewAuthFromRequest(req *http.Request, creds CredentialsLookupFunc, nonce NonceCheckFunc) (*Auth, error) {
	header := req.Header.Get("Authorization")
	bewit := req.URL.Query().Get("bewit")

	var auth *Auth
	var err error
	if header != "" {
		auth, err = ParseRequestHeader(header)
		if err != nil {
			return nil, err
		}
	}
	if auth == nil && bewit != "" {
		if req.Method != "GET" && req.Method != "HEAD" {
			return nil, ErrInvalidBewitMethod
		}
		auth, err = ParseBewit(bewit)
		if err != nil {
			return nil, err
		}
	}
	if auth == nil {
		return nil, ErrNoAuth
	}

	auth.Method = req.Method
	auth.RequestURI = req.URL.Path
	if req.URL.RawQuery != "" {
		auth.RequestURI += "?" + req.URL.RawQuery
	}
	if bewit != "" {
		auth.Method = "GET"
		bewitPattern, _ := regexp.Compile(`\?bewit=` + bewit + `\z|bewit=` + bewit + `&|&bewit=` + bewit + `\z`)
		auth.RequestURI = bewitPattern.ReplaceAllString(auth.RequestURI, "")
	}
	auth.Host, auth.Port = extractReqHostPort(req)
	if creds != nil {
		err = creds(&auth.Credentials)
		if err != nil {
			return nil, err
		}
	}
	if nonce != nil && !auth.IsBewit && !nonce(auth.Nonce, auth.Timestamp, &auth.Credentials) {
		return nil, ErrReplay
	}
	return auth, nil
}

func extractReqHostPort(req *http.Request) (host string, port string) {
	if idx := strings.Index(req.Host, ":"); idx != -1 {
		host, port, _ = net.SplitHostPort(req.Host)
	} else {
		host = req.Host
	}
	if req.URL.Host != "" {
		if idx := strings.Index(req.Host, ":"); idx != -1 {
			host, port, _ = net.SplitHostPort(req.Host)
		} else {
			host = req.URL.Host
		}
	}
	if port == "" {
		if req.URL.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}
	return
}

// NewRequestAuth builds a client Auth based on req and creds. tsOffset will be
// applied to Now when setting the timestamp.
func NewRequestAuth(req *http.Request, creds *Credentials, tsOffset time.Duration) *Auth {
	auth := &Auth{
		Method:      req.Method,
		Credentials: *creds,
		Timestamp:   Now().Add(tsOffset),
		Nonce:       nonce(),
		RequestURI:  req.URL.RequestURI(),
	}
	auth.Host, auth.Port = extractReqHostPort(req)
	return auth
}

// NewRequestAuth builds a client Auth based on uri and creds. tsOffset will be
// applied to Now when setting the timestamp.
func NewURLAuth(uri string, creds *Credentials, tsOffset time.Duration) (*Auth, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	auth := &Auth{
		Method:      "GET",
		Credentials: *creds,
		Timestamp:   Now().Add(tsOffset),
	}
	if u.Path != "" {
		// url.Parse unescapes the path, which is unexpected
		auth.RequestURI = "/" + strings.SplitN(uri[8:], "/", 2)[1]
	} else {
		auth.RequestURI = "/"
	}
	auth.Host, auth.Port = extractURLHostPort(u)
	return auth, nil
}

func extractURLHostPort(u *url.URL) (host string, port string) {
	if idx := strings.Index(u.Host, ":"); idx != -1 {
		host, port, _ = net.SplitHostPort(u.Host)
	} else {
		host = u.Host
	}
	if port == "" {
		if u.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}
	return
}

func nonce() string {
	b := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)[:8]
}

const headerVersion = "1"

type Auth struct {
	Credentials Credentials

	Method     string
	RequestURI string
	Host       string
	Port       string

	MAC   []byte
	Nonce string
	Ext   string
	Hash  []byte

	// ReqHash is true if the request contained a hash
	ReqHash   bool
	IsBewit   bool
	Timestamp time.Time

	// ActualTimestamp is when the request was received
	ActualTimestamp time.Time
}

var headerRegex = regexp.MustCompile(`(id|ts|nonce|hash|ext|mac|app|dlg)="([ !#-\[\]-~]+)"`) // character class is ASCII printable [\x20-\x7E] without \ and "

// ParseHeader parses a Hawk request or response header and populates auth.
// t must be AuthHeader if the header is an Authorization header from a request
// or AuthResponse if the header is a Server-Authorization header from
// a response.
func (auth *Auth) ParseHeader(header string, t AuthType) error {
	if len(header) < 4 || strings.ToLower(header[:4]) != "hawk" {
		return AuthFormatError{"scheme", "must be Hawk"}
	}

	matches := headerRegex.FindAllStringSubmatch(header, 8)

	var err error
	for _, match := range matches {
		switch match[1] {
		case "hash":
			auth.Hash, err = base64.StdEncoding.DecodeString(match[2])
			if err != nil {
				return AuthFormatError{"hash", "malformed base64 encoding"}
			}
		case "ext":
			auth.Ext = match[2]
		case "mac":
			auth.MAC, err = base64.StdEncoding.DecodeString(match[2])
			if err != nil {
				return AuthFormatError{"mac", "malformed base64 encoding"}
			}
		default:
			if t == AuthHeader {
				switch match[1] {
				case "app":
					auth.Credentials.App = match[2]
				case "dlg":
					auth.Credentials.Delegate = match[2]
				case "id":
					auth.Credentials.ID = match[2]
				case "ts":
					ts, err := strconv.ParseInt(match[2], 10, 64)
					if err != nil {
						return AuthFormatError{"ts", "not an integer"}
					}
					auth.Timestamp = time.Unix(ts, 0)
				case "nonce":
					auth.Nonce = match[2]

				}
			}
		}

	}

	if len(auth.MAC) == 0 {
		return AuthFormatError{"mac", "missing or empty"}
	}

	return nil
}

// Valid confirms that the timestamp is within skew and verifies the MAC.
//
// If the request is valid, nil will be returned. If auth is a bewit and the
// method is not GET or HEAD, ErrInvalidBewitMethod will be returned. If auth is
// a bewit and the timestamp is after the the specified expiry, ErrBewitExpired
// will be returned. If auth is from a request header and the timestamp is
// outside the maximum skew, ErrTimestampSkew will be returned. If the MAC is
// not the expected value, ErrInvalidMAC will be returned.
func (auth *Auth) Valid() error {
	t := AuthHeader
	if auth.IsBewit {
		t = AuthBewit
		if auth.Method != "GET" && auth.Method != "HEAD" {
			return ErrInvalidBewitMethod
		}
		if auth.ActualTimestamp.After(auth.Timestamp) {
			return ErrBewitExpired
		}
	} else {
		skew := auth.ActualTimestamp.Sub(auth.Timestamp)
		if abs(skew) > MaxTimestampSkew {
			return ErrTimestampSkew
		}
	}
	if !hmac.Equal(auth.mac(t), auth.MAC) {
		if auth.IsBewit && strings.HasPrefix(auth.RequestURI, "http") && len(auth.RequestURI) > 9 {
			// try just the path
			uri := auth.RequestURI
			auth.RequestURI = "/" + strings.SplitN(auth.RequestURI[8:], "/", 2)[1]
			if auth.Valid() == nil {
				return nil
			}
			auth.RequestURI = uri
		}
		return ErrInvalidMAC
	}
	return nil
}

func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// ValidResponse checks that a response Server-Authorization header is correct.
//
// ErrMissingServerAuth is returned if header is an empty string. ErrInvalidMAC
// is returned if the MAC is not the expected value.
func (auth *Auth) ValidResponse(header string) error {
	if header == "" {
		return ErrMissingServerAuth
	}
	err := auth.ParseHeader(header, AuthResponse)
	if err != nil {
		return err
	}
	if !hmac.Equal(auth.mac(AuthResponse), auth.MAC) {
		return ErrInvalidMAC
	}
	return nil
}

// PayloadHash initializes a hash for body validation. To validate a request or
// response body, call PayloadHash with contentType set to the body Content-Type
// with all parameters and prefix/suffix whitespace stripped, write the entire
// body to the returned hash, and then validate the hash with ValidHash.
func (auth *Auth) PayloadHash(contentType string) hash.Hash {
	h := auth.Credentials.Hash()
	h.Write([]byte("hawk." + headerVersion + ".payload\n" + contentType + "\n"))
	return h
}

// ValidHash writes the final newline to h and checks if it matches auth.Hash.
func (auth *Auth) ValidHash(h hash.Hash) bool {
	h.Write([]byte("\n"))
	return bytes.Equal(h.Sum(nil), auth.Hash)
}

// SetHash writes the final newline to h and sets auth.Hash to the sum. This is
// used to specify a response payload hash.
func (auth *Auth) SetHash(h hash.Hash) {
	h.Write([]byte("\n"))
	auth.Hash = h.Sum(nil)
	auth.ReqHash = false
}

// ResponseHeader builds a response header based on the auth and provided ext,
// which may be an empty string. Use PayloadHash and SetHash before
// ResponseHeader to include a hash of the response payload.
func (auth *Auth) ResponseHeader(ext string) string {
	auth.Ext = ext
	if auth.ReqHash {
		auth.Hash = nil
	}

	h := `Hawk mac="` + base64.StdEncoding.EncodeToString(auth.mac(AuthResponse)) + `"`
	if auth.Ext != "" {
		h += `, ext="` + auth.Ext + `"`
	}
	if auth.Hash != nil {
		h += `, hash="` + base64.StdEncoding.EncodeToString(auth.Hash) + `"`
	}

	return h
}

// RequestHeader builds a request header based on the auth.
func (auth *Auth) RequestHeader() string {
	auth.MAC = auth.mac(AuthHeader)

	h := `Hawk id="` + auth.Credentials.ID +
		`", mac="` + base64.StdEncoding.EncodeToString(auth.MAC) +
		`", ts="` + strconv.FormatInt(auth.Timestamp.Unix(), 10) +
		`", nonce="` + auth.Nonce + `"`

	if len(auth.Hash) > 0 {
		h += `, hash="` + base64.StdEncoding.EncodeToString(auth.Hash) + `"`
	}
	if auth.Ext != "" {
		h += `, ext="` + auth.Ext + `"`
	}
	if auth.Credentials.App != "" {
		h += `, app="` + auth.Credentials.App + `"`
	}
	if auth.Credentials.Delegate != "" {
		h += `, dlg="` + auth.Credentials.Delegate + `"`
	}

	return h
}

// Bewit creates and encoded request bewit parameter based on the auth.
func (auth *Auth) Bewit() string {
	auth.Method = "GET"
	auth.Nonce = ""
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(auth.Credentials.ID+`\`+
		strconv.FormatInt(auth.Timestamp.Unix(), 10)+`\`+
		base64.StdEncoding.EncodeToString(auth.mac(AuthBewit))+`\`+
		auth.Ext)), "=")
}

// NormalizedString builds the string that will be HMACed to create a request
// MAC.
func (auth *Auth) NormalizedString(t AuthType) string {
	str := "hawk." + headerVersion + "." + t.String() + "\n" +
		strconv.FormatInt(auth.Timestamp.Unix(), 10) + "\n" +
		auth.Nonce + "\n" +
		auth.Method + "\n" +
		auth.RequestURI + "\n" +
		auth.Host + "\n" +
		auth.Port + "\n" +
		base64.StdEncoding.EncodeToString(auth.Hash) + "\n" +
		auth.Ext + "\n"

	if auth.Credentials.App != "" {
		str += auth.Credentials.App + "\n"
		str += auth.Credentials.Delegate + "\n"
	}

	return str
}

func (auth *Auth) mac(t AuthType) []byte {
	mac := auth.Credentials.MAC()
	mac.Write([]byte(auth.NormalizedString(t)))
	return mac.Sum(nil)
}

func (auth *Auth) tsMac(ts string) []byte {
	mac := auth.Credentials.MAC()
	mac.Write([]byte("hawk." + headerVersion + ".ts\n" + ts + "\n"))
	return mac.Sum(nil)
}

// StaleTimestampHeader builds a signed WWW-Authenticate response header for use
// when Valid returns ErrTimestampSkew.
func (auth *Auth) StaleTimestampHeader() string {
	ts := strconv.FormatInt(Now().Unix(), 10)
	return `Hawk ts="` + ts +
		`", tsm="` + base64.StdEncoding.EncodeToString(auth.tsMac(ts)) +
		`", error="Stale timestamp"`
}

var tsHeaderRegex = regexp.MustCompile(`(ts|tsm|error)="([ !#-\[\]-~]+)"`) // character class is ASCII printable [\x20-\x7E] without \ and "

// UpdateOffset parses a signed WWW-Authenticate response header containing
// a stale timestamp error and updates auth.Timestamp with an adjusted
// timestamp.
func (auth *Auth) UpdateOffset(header string) (time.Duration, error) {
	if len(header) < 4 || strings.ToLower(header[:4]) != "hawk" {
		return 0, AuthFormatError{"scheme", "must be Hawk"}
	}

	matches := tsHeaderRegex.FindAllStringSubmatch(header, 3)

	var err error
	var ts time.Time
	var tsm []byte
	var errMsg string

	for _, match := range matches {
		switch match[1] {
		case "ts":
			t, err := strconv.ParseInt(match[2], 10, 64)
			if err != nil {
				return 0, AuthFormatError{"ts", "not an integer"}
			}
			ts = time.Unix(t, 0)
		case "tsm":
			tsm, err = base64.StdEncoding.DecodeString(match[2])
			if err != nil {
				return 0, AuthFormatError{"tsm", "malformed base64 encoding"}
			}
		case "error":
			errMsg = match[2]
		}
	}

	if errMsg != "Stale timestamp" {
		return 0, AuthFormatError{"error", "missing or unknown"}
	}

	if !hmac.Equal(tsm, auth.tsMac(strconv.FormatInt(ts.Unix(), 10))) {
		return 0, ErrInvalidMAC
	}

	offset := ts.Sub(Now())
	auth.Timestamp = ts
	auth.Nonce = nonce()
	return offset, nil
}
