package genkeys

import (
	"fmt"
	"go_systems/pr0config"
	"go_systems/proconfs"
	"testing"
)

// Test_proGenKeys TESTS Encryption of the Private Key
func Test_proGenKeys(t *testing.T) {

	// Need  a public key
	// openssl rsa -in mykey.pem -pubout > mykey.pub

	pk, err := EncryptionPrivateKey(1028, "HARDTOGUESS")
	if err != nil {
		fmt.Println(err)
	}
	f, ok, err := proconfs.CreateFile(pr0config.KeyCertPath, "mykey.pem")
	if !ok {
		fmt.Println(err)
	}
	proconfs.WriteFile(f, pk)
}
