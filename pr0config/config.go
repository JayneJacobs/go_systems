package pr0config

import (
	"crypto/rsa"
	"fmt"
	"go_systems/proconfs"

	jwtgo "github.com/dgrijalva/jwt-go"
)

var (
	// PubKeyFile for RSA pub key
	PubKeyFile *rsa.PublicKey
	// PrivKeyFile rsa Key
	PrivKeyFile *rsa.PrivateKey
)

const (
	// PKPWD password for key
	PKPWD = "HARDTOGUESS"
	// FileStoragePath for uploads
	FileStoragePath = "/var/www/uploads/"
	// KeyCertPath  for key
	KeyCertPath = "/var/www/keycertz/"
	// PrivKeyPath pem
	PrivKeyPath = "/var/www/keycertz/mykey.pem"
	// PubKeyPath pubb
	PubKeyPath = "/var/www/keycertz/mykey.pub"
	// MongoHost is the address of the db
	MongoHost = "127.0.0.1"
	// MongoUser username string
	MongoUser = "mongod"
	// MongoPassword the password string
	MongoPassword = "ThisPassword"
	// Mongodb name of db string
	MongoDb = "admin"
)

func init() {
	f, ok, err := proconfs.ReadFile(PubKeyPath)
	if !ok || err != nil {
		fmt.Println(err)
	}
	PubKeyFile, err = jwtgo.ParseRSAPublicKeyFromPEM(f)
	if err != nil {
		fmt.Println(err)
	}
	f, ok, err = proconfs.ReadFile(PrivKeyPath)
	if !ok || err != nil {
		fmt.Println(err)
	}
	PrivKeyFile, err = jwtgo.ParseRSAPrivateKeyFromPEMWithPassword(f, PKPWD)
	if err != nil {
		fmt.Println(err)
	}
}
