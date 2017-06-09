package main

import (
	"fmt"
	"io/ioutil"
)

// START OMIT
func printFile(name string) {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buf))
}

// END OMIT
