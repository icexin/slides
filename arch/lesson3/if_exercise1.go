package main

import (
	"fmt"
	"math/rand"
)

func main() {
	x := rand.Intn(10)
	fmt.Print("guess a number 1-10:")
	var n int
	fmt.Scanf("%d", &n)

	// 补全代码，如果n==x 输出"正确"
	// 如果n>x输出"过大"
	// 如果n<x输出"过小"
}
