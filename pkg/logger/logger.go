// The logger package contains configured logging tools for debugging
package logger

import (
	"log"
	"os"
)

// Verbosity is set by the CLI flags and can be used to enable/disable global logging
var Verbosity bool

var infoLogger *log.Logger  // infoLogger for general logging
var errorLogger *log.Logger // errorlogger for exclusive error reporting

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// exported functions to simplify infoLogger
func Info(message string) {
	if Verbosity {
		infoLogger.Println(message)
	}
}

// exported functions to simplify errorLogger
func Error(err error) {
	errorLogger.Fatal(err)
}
