package syslog

type Facility byte

const (
	Kern Facility = iota
	User
	Mail
	Daemon
	Auth
	Syslog
	Lpr
	News
	Uucp
	Cron
	Authpriv
	System0
	System1
	System2
	System3
	System4
	Local0
	Local1
	Local2
	Local3
	Local4
	Local5
	Local6
	Local7
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

func (f Facility) String() string {
	if f > Local7 {
		return "unknown"
	}
	return facToStr[f]
}

type Severity byte

const (
	Emerg Severity = iota
	Alert
	Crit
	Err
	Warning
	Notice
	Info
	Debug
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

func (s Severity) String() string {
	if s > Debug {
		return "unknown"
	}
	return sevToStr[s]
}
