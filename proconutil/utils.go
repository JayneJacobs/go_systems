package proconutil

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"go_systems/procondata"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

//B64DecodeTryUser take a test user string and returns 64 bit encoded email and password and error
func B64DecodeTryUser(test_user_json_str string) ([]byte, []byte, error) {
	var tu procondata.TryUser
	err := json.Unmarshal([]byte(test_user_json_str), &tu)
	if err != nil {
		return nil, nil, err
	}
	tuUsDec, _ := b64.StdEncoding.DecodeString(string(tu.Email))
	tuPsDec, _ := b64.StdEncoding.DecodeString(string(tu.Password))

	return tuUsDec, tuPsDec, nil

}

// GenerateUserPassword takes pwdstr string ane returns 64bit encoded string
func GenerateUserPassword(pwdstr string) string {
	hp, err := bcrypt.GenerateFromPassword([]byte(pwdstr), 0)
	if err != nil {
		fmt.Printf("Generate BCrypt Error: %s", err)
	}
	return string(hp)
}

// ValidateUserPassword takes a password and byteHash
// as a slice of bytes and returns a bool and error
func ValidateUserPassword(tryPass []byte, byteHash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(byteHash, tryPass)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SendMsg takes three strings a json, text and data and a websocker
func SendMsg(j string, t string, d string, c *websocket.Conn) {
	m := procondata.Msg{j, t, d}

	if err := c.WriteJSON(m); err != nil {
		fmt.Println("Error writing Json in utils.go", err)
	}
}
