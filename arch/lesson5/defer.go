package main

import "fmt"

func print() {
	defer func() {
		fmt.Println("defer")
	}()
	fmt.Println("hello")
}

func main() {
	print()
}
