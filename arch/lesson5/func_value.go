package main

import (
	"fmt"
	"os"
	"strconv"
)

func add(m, n int) int {
	return m + n
}
func sub(m, n int) int {
	return m - n
}

func main() {
	funcmap := map[string]func(int, int) int{
		"+": add,
		"-": sub,
	}

	m, _ := strconv.Atoi(os.Args[1])
	n, _ := strconv.Atoi(os.Args[3])

	f := funcmap[os.Args[2]]
	if f != nil {
		fmt.Println(f(m, n))
	}
}
