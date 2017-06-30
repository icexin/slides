package main

import "fmt"

func addn(n int) func(int) int {
	return func(m int) int {
		return m + n
	}
}

func main() {
	f := addn(3)
	fmt.Println(f(2))
}
