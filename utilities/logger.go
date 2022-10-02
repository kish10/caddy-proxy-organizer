package utilities

import (
	"log"
	"os"
)

// Reference: https://rollbar.com/blog/golang-error-logging-guide/#
var (
	InfoLog = log.New(os.Stdout, "\nINFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(os.Stdout, "\nWARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLog = log.New(os.Stdout, "\nERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)
