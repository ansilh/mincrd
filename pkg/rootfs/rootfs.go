/*
Package rootfs ...
Package to support contianer root filesystem creation
*/
package rootfs

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//CreateRootDir ...
//Create the base directory for storing container images
func CreateRootDir(path string) error {
	if len(path) > 0 {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

//CheckRoot ...
//Returns true if exists
func CheckRoot(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//RemoveRoot ...
//Function to remove the base directory if there is no contents
func RemoveRoot(path string) bool {

	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//GenCtrID ...
//Generate container name with supplied name and a random Id
func GenCtrID(path string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprint(path, "-", strconv.FormatInt(rand.Int63(), 10))
}
