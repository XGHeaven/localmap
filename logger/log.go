package logger

import (
	"log"
	"os"
)

const (
	ERROR = iota
	WARNING
	INFO
	DEBUG
)

var (
	infoLog    = log.New(os.Stdout, "[info]\t", log.LstdFlags)
	debugLog   = log.New(os.Stdout, "[debug]\t", log.LstdFlags)
	errLog     = log.New(os.Stderr, "[error]\t", log.LstdFlags)
	warninglog = log.New(os.Stderr, "[warning]", log.LstdFlags)

	LogLevel = INFO
)

func Info(v ...interface{}) {
	if LogLevel < INFO {
		return
	}
	infoLog.Println(v...)
}

func Infof(format string, v ...interface{}) {
	if LogLevel < INFO {
		return
	}
	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}
	infoLog.Printf(format, v...)
}

func Debug(v ...interface{}) {
	if LogLevel < DEBUG {
		return
	}
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
	if LogLevel < WARNING {
		return
	}
	warninglog.Println(v...)
}
