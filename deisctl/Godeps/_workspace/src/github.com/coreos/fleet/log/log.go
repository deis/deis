package log

import (
	"log"
	"os"
	"strconv"
	"sync/atomic"
)

var (
	verbosity = VLevel(0)

	iLog = log.New(os.Stdout, "INFO ", log.Lshortfile)
	eLog = log.New(os.Stdout, "ERROR ", log.Lshortfile)
	wLog = log.New(os.Stdout, "WARN ", log.Lshortfile)
	fLog = log.New(os.Stdout, "FATAL ", log.Lshortfile)
)

func SetVerbosity(lvl int) {
	verbosity.set(int32(lvl))
}

type VLevel int32

func (l *VLevel) get() VLevel {
	return VLevel(atomic.LoadInt32((*int32)(l)))
}

func (l *VLevel) String() string {
	return strconv.FormatInt(int64(*l), 10)
}

func (l *VLevel) Get() interface{} {
	return l.get()
}

func (l *VLevel) Set(val string) error {
	vi, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	l.set(int32(vi))
	return nil
}

func (l *VLevel) set(lvl int32) {
	atomic.StoreInt32((*int32)(l), lvl)
}

type VLogger bool

func V(level VLevel) VLogger {
	return VLogger(verbosity.get() >= level)
}

func (v VLogger) Info(args ...interface{}) {
	if v {
		iLog.Print(args...)
	}
}

func (v VLogger) Infof(format string, args ...interface{}) {
	if v {
		iLog.Printf(format, args...)
	}
}

func Info(args ...interface{}) {
	iLog.Print(args...)
}

func Infof(fmt string, args ...interface{}) {
	iLog.Printf(fmt, args...)
}

func Error(args ...interface{}) {
	eLog.Print(args...)
}

func Errorf(fmt string, args ...interface{}) {
	eLog.Printf(fmt, args...)
}

func Warning(fmt string, args ...interface{}) {
	wLog.Print(args...)
}

func Warningf(fmt string, args ...interface{}) {
	wLog.Printf(fmt, args...)
}

func Fatal(args ...interface{}) {
	fLog.Fatal(args...)
}

func Fatalf(fmt string, args ...interface{}) {
	fLog.Fatalf(fmt, args...)
}
