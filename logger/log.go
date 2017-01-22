package logger

import (
  "log"
  "os"
)

var (
  infoLog = log.New(os.Stdout, "[info]\t", log.LstdFlags)
  debugLog = log.New(os.Stdout, "[debug]\t", log.LstdFlags)
  errLog = log.New(os.Stderr, "[error]\t", log.LstdFlags)
  warninglog = log.New(os.Stderr, "[warning]", log.LstdFlags)
)

func Info(v ...interface{}) {
  infoLog.Println(v...)
}

func Infof(format string, v ...interface{}) {
  if format[len(format)-1] != '\n' {
    format = format + "\n"
  }
  infoLog.Printf(format, v...)
}

func Debug(v ...interface{}) {
  debugLog.Println(v...)
}

func Error(v ...interface{}) {
  errLog.Println(v...)
}

func Fatal(v ...interface{}) {
  errLog.Fatalln(v...)
}

func Panic(v ...interface{}) {
  errLog.Panicln(v...)
}

func Warning(v ...interface{}) {
  warninglog.Println(v...)
}
