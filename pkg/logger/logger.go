package logger

import (
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	level        = INFO
	levelMapping = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}
)

func SetLevel(logLevel Level) {
	level = logLevel
}

func Debug(v ...interface{}) {
	if level <= DEBUG {
		log.Print("[", levelMapping[DEBUG], "] ", v)
	}
}

func Info(v ...interface{}) {
	if level <= INFO {
		log.Print("[", levelMapping[INFO], "] ", v)
	}
}

func Warn(v ...interface{}) {
	if level <= WARN {
		log.Print("[", levelMapping[WARN], "] ", v)
	}
}

func Error(v ...interface{}) {
	if level <= ERROR {
		log.Print("[", levelMapping[ERROR], "] ", v)
	}
}

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}
