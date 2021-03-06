package proconmongo

import (
	"context"
	"encoding/json"
	"fmt"
	"go_systems/apilogging"
	"go_systems/pr0config"
	"go_systems/procondata"
	"go_systems/proconutil"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type key string

const (
	// HostKey string
	HostKey = key("hostKey")
	// UsernameKey string
	UsernameKey = key("usernameKey")
	// PasswordKey string
	PasswordKey = key("passwordKey")
	// DatabaseKey string
	DatabaseKey = key("databaseKey")
)

var ctx context.Context
var cancel func()
var client *mongo.Client;
var err error

func init() {
	ctx = context.Background()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, HostKey, pr0config.MongoHost)
	ctx = context.WithValue(ctx, UsernameKey, pr0config.MongoUser)
	ctx = context.WithValue(ctx, PasswordKey, pr0config.MongoPassword)
	ctx = context.WithValue(ctx, DatabaseKey, pr0config.MongoDb)
	
	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		ctx.Value(UsernameKey).(string),
		ctx.Value(PasswordKey).(string),
		ctx.Value(HostKey).(string),
		ctx.Value(DatabaseKey).(string),
	)
	LogName := "init"
	respString := uri
	
	apilogging.Resplog(respString, LogName)
	clientoptions := options.Client().ApplyURI(uri)

	client, err = mongo.Connect(ctx, clientoptions)

	if err != nil {
		fmt.Println("Errror Connecting",err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Mongo Connected")

}

// CreateUser takes a create_user string and returns a status bool
func CreateUser(jsonCreateuser string, ws *websocket.Conn) bool {
	user := procondata.AUser{}
	err = json.Unmarshal([]byte(jsonCreateuser), &user)
	LogName := "CreateUser"
	if err != nil {
		fmt.Println("CreateUserin thismongo", err)
	}
	usr, pwd, err := proconutil.B64DecodeTryUser(jsonCreateuser)
	if err != nil {
		fmt.Println("Error with BB64Decode", err)
	}
	fmt.Println("Past line 81 in thismongo create user")
	           respString := ("Past line 81 in thismongo create user\n")
				apilogging.Resplog(respString, LogName)
	user.Email = string(usr)
	user.Password = string(pwd)

	collection := client.Database("api").Collection("users")

	//Check for a user
	
	var xdoc interface{}
	fmt.Println("In GWS.go Create User one ", string(usr), string(pwd))
	respString = ("In GWS.go Create User one"+ string(usr)+ string(pwd))
	apilogging.Resplog(respString, LogName)
	filter := bson.D{{"user", user.Email}}
	err = collection.FindOne(ctx, filter).Decode(&xdoc)
  
    respString = "User" + user.Email
	apilogging.Resplog(respString, LogName)
	if err != nil {
		respString = ("\nThis is the error generated in thismongo.go by collection.FindOne L#99 " + err.Error())
		apilogging.Resplog(respString, LogName)
		if xdoc == nil {
			fmt.Println("xdoc is  nil no message sent to Create User : ", xdoc)
			hp := proconutil.GenerateUserPassword(user.Password)
			fmt.Println("\nIn in thismongo if xdoc == nil, inserted password: #L104 ", hp)
			respString = ("\nIn in thismongo if xdoc == nil, inserted password: #L104: "+ hp)
	        apilogging.Resplog(respString, LogName)
			user.Password = hp
			user.Role = "Generic"

			insertResult, err := collection.InsertOne(ctx, &user)
			if err != nil {
				fmt.Println("Error Inserting Document with collection.InsertOne L#103", err)
				respString = ("\nError Inserting Document with collection.InsertOne L#103"+ err.Error())
				apilogging.Resplog(respString, LogName)
				return false
			}

			fmt.Println("Inserted User:  Line 108", insertResult.InsertedID)
			respString = ("\nInserted User:  Line 108 "+ fmt.Sprint(insertResult.InsertedID))
			apilogging.Resplog(respString, LogName)
			procondata.SendMsg("vAr", "toast-success", "user created successfully", ws)
			procondata.SendMsg("vAr", "user-created-successfully", "User created successfully", ws)

			return true
		}
		proconutil.SendMsg("vAr", "user-already-exists", "User Already Exists", ws)
		return false
	}
	proconutil.SendMsg("vAr", "user-already-exists", "User Already Exists", ws)
	fmt.Println("In in thismongo no if statements ran for Create User; meaning the user may already exist. ")
	        respString = ("\nIn in thismongo no if statements ran for Create User; meaning the user may already exist. \n"+ fmt.Sprint(ws))
			apilogging.Resplog(respString, LogName)

	return false
}

// MongoTryUser takes a username and password as a slice of bytes and returns bbool and Userstruct and error
func MongoTryUser(u []byte, p []byte) (bool, *procondata.AUser, error) {
	var xdoc procondata.AUser
	collection := client.Database("api").Collection("users")
	filter := bson.D{{"email", string(u)}}
	if err = collection.FindOne(ctx, filter).Decode(&xdoc); err != nil {
		    LogName := "CreateUser"
		    respString := ("L#150 error: \n "+ err.Error())
			apilogging.Resplog(respString, LogName)
		return false, nil, err
		    
	}
	bres, err := proconutil.ValidateUserPassword(p, []byte(xdoc.Password))
	if err != nil {
		LogName := "CreateUser"
		respString := ("L#150 error:\n "+ err.Error())
			apilogging.Resplog(respString, LogName)
		return false, nil, err
	
	}
	return bres, &xdoc, nil
}

// MongoGetUIComponent takes a string and Response Writer from http package
// It defines teh ui db api
func MongoGetUIComponent(component string, w http.ResponseWriter) {
	var xdoc map[string]interface{}
	collection := client.Database("api").Collection("ui")		
	LogName := "CreateUser"
	filter := bson.D{{"component", component}}
	fmt.Println("In MongoGetUIComponent", component)
	respString := ("In MongoGetUIComponent \n"+ component)
			apilogging.Resplog(respString, LogName)
	err = collection.FindOne(ctx, filter).Decode(&xdoc)
	if err != nil { 
		LogName := "CreateUser"
		fmt.Println("Error in This Mongo 145",err)
		respString = ("\nError in This Mongo 145: " + err.Error())
			apilogging.Resplog(respString, LogName)
	 }  
	json.NewEncoder(w).Encode(&xdoc)  
	
}
