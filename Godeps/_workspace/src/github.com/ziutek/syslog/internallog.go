package syslog

// Logger is an interface for package internal (non fatal) logging
type Logger interface {
	Print(...interface{})
	Printf(format string, v ...interface{})
	Println(...interface{})
}

// FatalLogger is an interface for logging package internal fatal errors
type FatalLogger interface {
	Fatal(...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(...interface{})
}
