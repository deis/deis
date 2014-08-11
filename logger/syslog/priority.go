package syslog

type Facility byte

// The following is a list of Facilities as defined by RFC 3164.
const (
	Kern Facility = iota // kernel messages
	User                 // user-level messages
	Mail                 // mail system
	Daemon               // system daemons
	Auth                 // security/authorization messages
	Syslog               // messages internal to syslogd
	Lpr                  // line printer subsystem
	News                 // newtork news subsystem
	Uucp                 // UUCP subsystem
	Cron                 // cron messages
	Authpriv             // security/authorization messages
	System0              // historically FTP daemon
	System1              // historically NTP subsystem
	System2              // historically log audit
	System3              // historically log alert
	System4              // historically clock daemon, some operating systems use this for cron
	Local0               // local use 0
	Local1               // local use 1
	Local2               // local use 2
	Local3               // local use 3
	Local4               // local use 4
	Local5               // local use 5
	Local6               // local use 6
	Local7               // local use 7
)

var facToStr = [...]string{
	"kern",
	"user",
	"mail",
	"daemon",
	"auth",
	"syslog",
	"lpr",
	"news",
	"uucp",
	"cron",
	"authpriv",
	"system0",
	"system1",
	"system2",
	"system3",
	"system4",
	"local0",
	"local1",
	"local2",
	"local3",
	"local4",
	"local5",
	"local6",
	"local7",
}

// String returns a string representation of the Facility. This satisfies the
// fmt.Stringer interface.
func (f Facility) String() string {
	if f > Local7 {
		return "unknown"
	}
	return facToStr[f]
}

type Severity byte

const (
	Emerg Severity = iota // Emergency: system is unusable
	Alert                 // immediate action required
	Crit                  // critical conditions
	Err                   // error conditions
	Warning               // warning conditions
	Notice                // normal but significant condition
	Info                  // information message
	Debug                 // debug-level message
)

var sevToStr = [...]string{
	"emerg",
	"alert",
	"crit",
	"err",
	"waining",
	"notice",
	"info",
	"debug",
}

// String returns a string representation of the Severity. This satisfies the
// fmt.Stringer interface.
func (s Severity) String() string {
	if s > Debug {
		return "unknown"
	}
	return sevToStr[s]
}
