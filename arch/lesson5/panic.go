package main

import "fmt"

func print() {
	var p *int
	fmt.Println(*p)
}

func main() {
	var n int
	fmt.Println(10 / n)
	print()

	var slice [3]int
	fmt.Println(slice[3])
}
