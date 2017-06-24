package main

import "fmt"

func main() {
	ages := map[string]int{
		"a": 1,
		"b": 2,
	}
	for name, age := range ages {
		fmt.Println("name", name, "age", age)
	}

	for name := range ages {
		fmt.Println(name)
	}
}
