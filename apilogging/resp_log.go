package apilogging

import (
	"fmt"
	"io"
	"os"

	apilogger "go_systems/apilogging/alogger"
)

// Resplog is the logging for all API calls it takes a string for the LogName and the API response Data as a string.
func Resplog(respData string, name string) {
	fmt.Printf("see the log: %v.log\n", name)
	logname := name + ".log"
	f, err := os.OpenFile(logname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	//fmt.Println(os.Stdout, string(htmlData))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	io.WriteString(f, respData)
}

// APILog will take a string and log start of the server
func APILog(s string, l string) {
	logger := apilogger.GetInstance("URLs.log")

	logger.Printf("This is the URL for %v:  %v", l, s)
	// http.HandleFunc("/", sroot)
	// http.ListenAndServe(":8080", nil)
}
