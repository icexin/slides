package main

import "fmt"

func main() {
	// FOR1 OMIT
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}

	// END OMIT

	// FOR2 OMIT
	i := 5
	for i < 7 {
		fmt.Println(i)
		i = i + 1
	}

	// END OMIT

	// FOR3 OMIT
	i = 8
	// 等价于 while true
	for {
		i = i + 1
		fmt.Println(i)
		if i > 10 {
			break
		}
	}

	// END OMIT
}
