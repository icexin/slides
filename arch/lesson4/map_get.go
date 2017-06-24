package main

import "fmt"

func main() {
	ages := map[string]int{
		"a": 1,
		"b": 2,
	}
	fmt.Println(ages["a"])
	ages["a"] = ages["b"] + 2

	c, ok := ages["c"]
	if ok {
		fmt.Println(c)
	} else {
		fmt.Println("not found")
	}

	if c, ok := ages["c"]; ok {
		fmt.Println(c)
	}
}
