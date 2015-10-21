// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	calldepth = 2
)

var (
	logger = log.New(os.Stderr, "", 0)
	debug  = false
)

func EnableTimestamps() {
	logger.SetFlags(logger.Flags() | log.Ldate | log.Ltime)
}

func EnableDebug() {
	debug = true
}

func Debug(v ...interface{}) {
	if debug {
		logger.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
	}
}

func Debugf(format string, v ...interface{}) {
	if debug {
		logger.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
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

func Warning(v ...interface{}) {
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
