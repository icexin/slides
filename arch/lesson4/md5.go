package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	data := []byte("hello")
	md5sum := md5.Sum(data)
	fmt.Printf("%x\n", md5sum)
}
