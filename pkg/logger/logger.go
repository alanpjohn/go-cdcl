package dbg

import (
	"log"
	"os"
)

var Verbosity bool

var infoLogger *log.Logger
var errorLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(message string) {
	if Verbosity {
		infoLogger.Println(message)
	}
}

func Error(err error) {
	errorLogger.Fatal(err)
}
