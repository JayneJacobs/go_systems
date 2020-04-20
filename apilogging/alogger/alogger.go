package apilogger

import (
	"fmt"
	"log"
	"os"
)

// APILogger methods use a file name to create a log
type APILogger struct {
	*log.Logger
	filename string
	message  string
}

var alogger *APILogger

//GetInstance create a singleton instance of the hydra logger
func GetInstance(name string) *APILogger {

	alogger = createLogger(name)

	return alogger
}

//Create a logger instance
func createLogger(fname string) *APILogger {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Printf("error opening %v: %v \n", fname, err)
	}
	//defer file.Close()
	return &APILogger{
		filename: fname,
		Logger:   log.New(file, fname+" ", log.Lshortfile),
	}
}
