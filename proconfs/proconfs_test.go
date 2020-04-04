package proconfs

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

///var/www/go_systems/

// Test_preconfs tests create test.txt
func Test_preconfs(t *testing.T) {

	CreateFile("./", "test.txt")
	fn, r, err := ReadFile("./test.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(fn))
	fmt.Printf("fn is %s : \n r is %v : \n", fn, r)
	y := time.Now()
	x := "I Replaced you at :" + y.String()
	WriteFile("./test.txt", []byte(x))
	fmt.Printf("fn is %s : \n r is %v : \n", fn, strconv.FormatBool(r))
}
