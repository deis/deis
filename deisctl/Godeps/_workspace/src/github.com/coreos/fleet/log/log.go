package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync/atomic"
)

const (
	calldepth = 2
)

var (
	logger    = log.New(os.Stderr, "", 0)
	verbosity = VLevel(0)
)

func EnableTimestamps() {
	logger.SetFlags(logger.Flags() | log.Ldate | log.Ltime)
}

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

func (vl VLogger) Info(v ...interface{}) {
	if vl {
		logger.Output(calldepth, header("INFO", fmt.Sprint(v...)))
	}
}

func (vl VLogger) Infof(format string, v ...interface{}) {
	if vl {
		logger.Output(calldepth, header("INFO", fmt.Sprintf(format, v...)))
	}
}

func Info(v ...interface{}) {
	logger.Output(calldepth, header("INFO", fmt.Sprint(v...)))
}

func Infof(format string, v ...interface{}) {
	logger.Output(calldepth, header("INFO", fmt.Sprintf(format, v...)))
}

func Error(v ...interface{}) {
	logger.Output(calldepth, header("ERROR", fmt.Sprint(v...)))
}

func Errorf(format string, v ...interface{}) {
	logger.Output(calldepth, header("ERROR", fmt.Sprintf(format, v...)))
}

func Warning(format string, v ...interface{}) {
	logger.Output(calldepth, header("WARN", fmt.Sprint(v...)))
}

func Warningf(format string, v ...interface{}) {
	logger.Output(calldepth, header("WARN", fmt.Sprintf(format, v...)))
}

func Fatal(v ...interface{}) {
	logger.Output(calldepth, header("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	logger.Output(calldepth, header("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func header(lvl, msg string) string {
	_, file, line, ok := runtime.Caller(calldepth)
	if ok {
		file = filepath.Base(file)
	}

	if len(file) == 0 {
		file = "???"
	}

	if line < 0 {
		line = 0
	}

	return fmt.Sprintf("%s %s:%d: %s", lvl, file, line, msg)
}
