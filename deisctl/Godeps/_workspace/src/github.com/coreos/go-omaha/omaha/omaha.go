/*
   Implements the Google omaha protocol.

   Omaha is a request/response protocol using XML. Requests are made by
   clients and responses are given by the Omaha server.
   http://code.google.com/p/omaha/wiki/ServerProtocol
*/
package omaha

import (
	"encoding/xml"
)

type Request struct {
	XMLName        xml.Name `xml:"request" datastore:"-"`
	Os             Os       `xml:"os"`
	Apps           []*App   `xml:"app"`
	Protocol       string   `xml:"protocol,attr"`
	Version        string   `xml:"version,attr,omitempty"`
	IsMachine      string   `xml:"ismachine,attr,omitempty"`
	SessionId      string   `xml:"sessionid,attr,omitempty"`
	UserId         string   `xml:"userid,attr,omitempty"`
	InstallSource  string   `xml:"installsource,attr,omitempty"`
	TestSource     string   `xml:"testsource,attr,omitempty"`
	RequestId      string   `xml:"requestid,attr,omitempty"`
	UpdaterVersion string   `xml:"updaterversion,attr,omitempty"`
}

func NewRequest(version string, platform string, sp string, arch string) *Request {
	r := new(Request)
	r.Protocol = "3.0"
	r.Os = Os{Version: version, Platform: platform, Sp: sp, Arch: arch}
	return r
}

func (r *Request) AddApp(id string, version string) *App {
	a := NewApp(id)
	a.Version = version
	r.Apps = append(r.Apps, a)
	return a
}

/* Response
 */
type Response struct {
	XMLName  xml.Name `xml:"response" datastore:"-" json:"-"`
	DayStart DayStart `xml:"daystart"`
	Apps     []*App   `xml:"app"`
	Protocol string   `xml:"protocol,attr"`
	Server   string   `xml:"server,attr"`
}

func NewResponse(server string) *Response {
	r := &Response{Server: server, Protocol: "3.0"}
	r.DayStart.ElapsedSeconds = "0"
	return r
}

type DayStart struct {
	ElapsedSeconds string `xml:"elapsed_seconds,attr"`
}

func (r *Response) AddApp(id string) *App {
	a := NewApp(id)
	r.Apps = append(r.Apps, a)
	return a
}

type App struct {
	XMLName     xml.Name     `xml:"app" datastore"-" json:"-"`
	Ping        *Ping        `xml:"ping"`
	UpdateCheck *UpdateCheck `xml:"updatecheck"`
	Events      []*Event     `xml:"event" json:",omitempty"`
	Id          string       `xml:"appid,attr,omitempty"`
	Version     string       `xml:"version,attr,omitempty"`
	NextVersion string       `xml:"nextversion,attr,omitempty"`
	Lang        string       `xml:"lang,attr,omitempty"`
	Client      string       `xml:"client,attr,omitempty"`
	InstallAge  string       `xml:"installage,attr,omitempty"`
	Status      string       `xml:"status,attr,omitempty"`

	// update engine extensions
	Track     string `xml:"track,attr,omitempty"`
	FromTrack string `xml:"from_track,attr,omitempty"`

	// coreos update engine extensions
	BootId    string `xml:"bootid,attr,omitempty"`
	MachineID string `xml:"machineid,attr,omitempty"`
	OEM       string `xml:"oem,attr,omitempty"`
}

func NewApp(id string) *App {
	a := &App{Id: id}
	return a
}

func (a *App) AddUpdateCheck() *UpdateCheck {
	a.UpdateCheck = new(UpdateCheck)
	return a.UpdateCheck
}

func (a *App) AddPing() *Ping {
	a.Ping = new(Ping)
	return a.Ping
}

func (a *App) AddEvent() *Event {
	event := new(Event)
	a.Events = append(a.Events, event)
	return event
}

type UpdateCheck struct {
	XMLName             xml.Name  `xml:"updatecheck" datastore:"-" json:"-"`
	Urls                *Urls     `xml:"urls"`
	Manifest            *Manifest `xml:"manifest"`
	TargetVersionPrefix string    `xml:"targetversionprefix,attr,omitempty"`
	Status              string    `xml:"status,attr,omitempty"`
}

func (u *UpdateCheck) AddUrl(codebase string) *Url {
	if u.Urls == nil {
		u.Urls = new(Urls)
	}
	url := new(Url)
	url.CodeBase = codebase
	u.Urls.Urls = append(u.Urls.Urls, *url)
	return url
}

func (u *UpdateCheck) AddManifest(version string) *Manifest {
	u.Manifest = &Manifest{Version: version}
	return u.Manifest
}

type Ping struct {
	XMLName        xml.Name `xml:"ping" datastore:"-" json:"-"`
	LastReportDays string   `xml:"r,attr,omitempty"`
	Status         string   `xml:"status,attr,omitempty"`
}

type Os struct {
	XMLName  xml.Name `xml:"os" datastore:"-" json:"-"`
	Platform string   `xml:"platform,attr,omitempty"`
	Version  string   `xml:"version,attr,omitempty"`
	Sp       string   `xml:"sp,attr,omitempty"`
	Arch     string   `xml:"arch,attr,omitempty"`
}

func NewOs(platform string, version string, sp string, arch string) *Os {
	o := &Os{Version: version, Platform: platform, Sp: sp, Arch: arch}
	return o
}

type Event struct {
	XMLName         xml.Name `xml:"event" datastore:"-" json:"-"`
	Type            string   `xml:"eventtype,attr,omitempty"`
	Result          string   `xml:"eventresult,attr,omitempty"`
	PreviousVersion string   `xml:"previousversion,attr,omitempty"`
	ErrorCode       string   `xml:"errorcode,attr,omitempty"`
}

type Urls struct {
	XMLName xml.Name `xml:"urls" datastore:"-" json:"-"`
	Urls    []Url    `xml:"url" json:",omitempty"`
}

type Url struct {
	XMLName  xml.Name `xml:"url" datastore:"-" json:"-"`
	CodeBase string   `xml:"codebase,attr"`
}

type Manifest struct {
	XMLName  xml.Name `xml:"manifest" datastore:"-" json:"-"`
	Packages Packages `xml:"packages"`
	Actions  Actions  `xml:"actions"`
	Version  string   `xml:"version,attr"`
}

type Packages struct {
	XMLName  xml.Name  `xml:"packages" datastore:"-" json:"-"`
	Packages []Package `xml:"package" json:",omitempty"`
}

type Package struct {
	XMLName  xml.Name `xml:"package" datastore:"-" json:"-"`
	Hash     string   `xml:"hash,attr"`
	Name     string   `xml:"name,attr"`
	Size     string   `xml:"size,attr"`
	Required bool     `xml:"required,attr"`
}

func (m *Manifest) AddPackage(hash string, name string, size string, required bool) *Package {
	p := &Package{Hash: hash, Name: name, Size: size, Required: required}
	m.Packages.Packages = append(m.Packages.Packages, *p)
	return p
}

type Actions struct {
	XMLName xml.Name  `xml:"actions" datastore:"-" json:"-"`
	Actions []*Action `xml:"action" json:",omitempty"`
}

type Action struct {
	XMLName xml.Name `xml:"action" datastore:"-" json:"-"`
	Event   string   `xml:"event,attr"`

	// Extensions added by update_engine
	ChromeOSVersion       string `xml:"ChromeOSVersion,attr"`
	Sha256                string `xml:"sha256,attr"`
	NeedsAdmin            bool   `xml:"needsadmin,attr"`
	IsDelta               bool   `xml:"IsDelta,attr"`
	DisablePayloadBackoff bool   `xml:"DisablePayloadBackoff,attr,omitempty"`
	MetadataSignatureRsa  string `xml:"MetadataSignatureRsa,attr,omitempty"`
	MetadataSize          string `xml:"MetadataSize,attr,omitempty"`
	Deadline              string `xml:"deadline,attr,omitempty"`
}

func (m *Manifest) AddAction(event string) *Action {
	a := &Action{Event: event}
	m.Actions.Actions = append(m.Actions.Actions, a)
	return a
}

var EventTypes = map[int]string{
	0:   "unknown",
	1:   "download complete",
	2:   "install complete",
	3:   "update complete",
	4:   "uninstall",
	5:   "download started",
	6:   "install started",
	9:   "new application install started",
	10:  "setup started",
	11:  "setup finished",
	12:  "update application started",
	13:  "update download started",
	14:  "update download finished",
	15:  "update installer started",
	16:  "setup update begin",
	17:  "setup update complete",
	20:  "register product complete",
	30:  "OEM install first check",
	40:  "app-specific command started",
	41:  "app-specific command ended",
	100: "setup failure",
	102: "COM server failure",
	103: "setup update failure",
	800: "ping",
}

var EventResults = map[int]string{
	0:  "error",
	1:  "success",
	2:  "success reboot",
	3:  "success restart browser",
	4:  "cancelled",
	5:  "error installer MSI",
	6:  "error installer other",
	7:  "noupdate",
	8:  "error installer system",
	9:  "update deferred",
	10: "handoff error",
}
