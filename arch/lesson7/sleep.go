package main

import (
	"fmt"
	"time"
)

func print(s string) {
	time.Sleep(time.Second)
	fmt.Println(s)
}

func main() {
	print("hello1")
	print("hello2")
	print("hello3")
}
