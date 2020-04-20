package apilogging

import (
	"fmt"
	apilogger "go_systems/apilogging/alogger"
	"net/http"
)

// Hydra will take a string and log start of the server

func sroot(w http.ResponseWriter, r *http.Request) {
	logger := apilogger.GetInstance("URLs")
	homeMessage := "Welcome to the API Test Tool"
	fmt.Fprint(w, homeMessage)

	logger.Println("Received an API Request")
}
