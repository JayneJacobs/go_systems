package proconjwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

//  GenerateJWT returns a jwt token given the sign in data name  alieas, email, role
func GenerateJWT(privkeyfile *rsa.PrivateKey, name string, alias string, email string, role string) (string, error) {
	fmt.Println()
	token := jwtgo.New(jwtgo.SigningMethodES256)
	in10m := time.Now().Add(time.Duration(30) * time.Minute).Unix()
	token.Claims = jwtgo.MapClaims{
		"iss":    "pr0con.selfmanagedmusician.com",
		"aud":    "pr0con.selfmanagedmusician.com",
		"exp":    in10m,
		"jti":    "Unique",
		"nbf":    2,
		"sub":    "subject",
		"scopes": "api:read.api:write",
		"name":   name,
		"alies":  alias,
		"email":  email,
		"role":   role,
	}

	tokenString, err := token.SignedString((privkeyfile))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT parses the token and validates it then returns status of true or false
func ValidateJWT(publickeyfile *rsa.PublicKey, jwt string) (bool, error) {
	token, err := jwtgo.Parse(jwt, func(token *jwtgo.Token) (interface{}, error) {
		return publickeyfile, nil
	})

	if err != nil {
		return false, err
	}
	if token.Valid && err == nil {
		return true, nil
	}

	return false, err

}
