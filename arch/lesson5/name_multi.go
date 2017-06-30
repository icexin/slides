package main

import "fmt"

func split(sum int) (x, y int) {
	x = sum / 10
	y = sum % 10
	return
}

func main() {
	fmt.Println(split(17))
}
