package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go_systems/pr0config"
	"go_systems/pr0conpty"
	"go_systems/procondata"
	"go_systems/proconjwt"
	"go_systems/proconmongo"
	"go_systems/proconutil"
	"go_systems/profilesystem"
	"strings"

	"net/http"

	"go_systems/proconasyncq"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:1200", "http service address")
var upgrader = websocket.Upgrader{} // default options 

func handleAPI(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WTF @HandleAPI Ws Upgrader Error in handlAPI ", err)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("WTF is up here in handleAPI", err)
	}

	c.UUID = "ws-" + id.String()
	fmt.Println("going into loop")

Loop:
	for {
		fmt.Println("in gws loop")
		in := procondata.Msg{}
		fmt.Println(&in)
		err := c.ReadJSON(&in)
		if err != nil {
			fmt.Println("In Loop Error in ReadJson", in.Data, err)
			c.Close()
			break Loop
		}
		
		fmt.Println("This is the in.Type", in.Type)
		
		switch in.Type {
		case "register-client-message":
			fmt.Println("message received: register-client-message")
			procondata.SendMsg("^vAr^", "server-ws-connect-success-msg", c.UUID, c)
		case "test-jwt-message":
			valid, err := proconjwt.ValidateJWT(pr0config.PubKeyFile, in.Jwt)
			if err != nil {
				fmt.Println("Error in gws case test-jwt-message", err)
				procondata.SendMsg("^vAr^", "jwt-token-invalid", err.Error(), c)
			}
			if err == nil && valid {
				fmt.Println("VALID JWT")
			}
		case "create-user":
			res := proconmongo.CreateUser(in.Data, c)
			fmt.Println("Mongo Function Result/Error for create-user: ", res)
		case "user-created-successfully":
			fmt.Println("User Created Successfully: ")
		case "login-user":
			usr, pwd, err := proconutil.B64DecodeTryUser(in.Data)
			if err != nil {
				fmt.Println("Error in gws proconutil.B64DecodeTryUser", err)
			}
			vres, auser, err := proconmongo.MongoTryUser(usr, pwd)
			if err != nil {
				fmt.Println("Error in gws  case proconmongo.MongoTryUser", err)
			}
			//fmt.Println("In gws, login-user", vres, auser.Email, auser.Password)
			auser.Password = "F00"
				fmt.Println("in gws case login-user", vres, auser.Password)
				jauser, err := json.Marshal(auser)
				if err != nil {
					fmt.Println("Error in gws switch marshaling auser login-user", err)
				}
				jwt, err := proconjwt.GenerateJWT(pr0config.PrivKeyFile, auser.Name, "@"+auser.Name, auser.Email, auser.Role)
				if err != nil {
					fmt.Println("JWT Generate error in gws.go switch case login-user", err)
				}
				if vres == false {
					procondata.SendMsg("^vAr^", "server-ws-connect-login-failure", string(jauser), c)
					fmt.Println("User Not found or invalid credentials: in gws case userlogin vres = false")
				}
				procondata.SendMsg(jwt, "server-ws-connect-success-jwt", string(jauser), c)
		case "validate-jwt":
			fallthrough
		case "validate-stored-jwt":
			valid, err := proconjwt.ValidateJWT(pr0config.PubKeyFile, in.Jwt)
			fmt.Println(in.Jwt)
			if err != nil {
				fmt.Println("Error in gws case validate-stored-jwt", err)
				if in.Type  == "validate-jwt" {
					procondata.SendMsg("^vAr^", "jwt-token-invalid", err.Error(), c)
				}
				if in.Type  == "validate-stored-jwt" {
					procondata.SendMsg("^vAr^", "stored-jwt-token-invalid", err.Error(), c)
				}
			}
			if err == nil && valid {
				fmt.Println("VALID JWT")
				if in.Type == "validate-jwt" {
					procondata.SendMsg("^vAr^", "server-ws-connect-jwt-verified", "noop", c)
				}
				if in.Type == "validate-stored-jwt" {
					procondata.SendMsg("^vAr^", "server-ws-connect-stored-jwt-verified", "noop", c)
				}
			}
		case "get-fs-path":
			if strings.HasPrefix(in.Data, "/var/www/VFS/") {
				tobj := profilesystem.NewGetFileSystemTask(in.Data, c)
				proconasyncq.TaskQueue <- tobj
				break
			}
			
		default:
			fmt.Println("Default case: No switch statemens in gws true")
			break
		}
	}
}

func handleUI(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	component := params["component"]

	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(component)

	proconmongo.MongoGetUIComponent(component, w)	
}

func main() {
	fmt.Println("This is from the Go Main Fumction")
	flag.Parse()

	// look into subrouter
	r := mux.NewRouter()

	
	
	//Rest API
	r.HandleFunc("/rest/api/ui/{component}", handleUI)

	//Websocket API
	r.HandleFunc("/ws", handleAPI)
	r.HandleFunc("/ws", pr0conpty.HandlePty)

	http.ListenAndServeTLS(*addr, "/etc/letsencrypt/live/pr0con.selfmanagedmusician.com/cert.pem", "/etc/letsencrypt/live/pr0con.selfmanagedmusician.com/privkey.pem", r)

}
