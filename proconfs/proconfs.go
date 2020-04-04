package proconfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// CreateFile detects if the file exoists and if not creates it with an input string as the name.
//   It returns an error and a bool and tne name(string)of the file
func CreateFile(p string, fn string) (string, bool, error) {
	// /strubg ij err
	// detect uf fuke exusts
	var _, err = os.Stat(p + fn)
	// create file if not exist
	if os.IsNotExist(err) {
		file, err := os.Create(p + fn)
		if err != nil {
			return "", false, err
		}

		defer file.Close()
	}
	return p + fn, true, nil
}

// ReadFile takes a string and returns a slice of bytes, bool and error
func ReadFile(fp string) ([]byte, bool, error) {
	//re-open file
	file, err := os.OpenFile(fp, os.O_RDWR, 0644)
	if err != nil {
		return []byte(""), false, err
	}

	defer file.Close()

	// read file, line bby line
	data := make([]byte, 1024)
	for {
		_, err := file.Read(data)

		if err == io.EOF {
			break
		}

		//break if error occured
		if err != nil && err != io.EOF {
			//    return []byte(""),false,err
			break
		}
	}
	fmt.Println("==> done reading from file")
	return data, true, nil
}

// WriteFile takes a string and a slice of bytes
func WriteFile(fp string, data []byte) error {
	// Open file using READ  WRITE permission
	// fmt.Println(data)
	file, err := os.OpenFile(fp, os.O_RDWR, 0544)
	if err != nil {
		return err
	}
	defer file.Close()
	// clear then set cursor at pos 0
	file.Truncate(0)
	file.Seek(0, 0)

	// Write some text line-by-line to file
	err = ioutil.WriteFile(file.Name(), data, 0544)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = file.Sync()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("==> done writing to file")
	return nil
}

// DeleteFile takes a string and deletes thefile and returns an error
func DeleteFile(p string) error {
	err := os.Remove(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("++> done deleting File")
	return nil
}
