package logger

import (
  "os"
  
  "github.com/Sirupsen/logrus"
)

var Log = logrus.New()

func init(){
  Log.Formatter = new(StdOutFormatter)

  logLevel := os.Getenv("LOG")
  if logLevel != "" {
    if level, err := logrus.ParseLevel(logLevel); err == nil {
      Log.Level = level
    }
  }  
}