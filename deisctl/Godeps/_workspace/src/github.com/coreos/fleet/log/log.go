/*
   Copyright 2014 CoreOS, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

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
