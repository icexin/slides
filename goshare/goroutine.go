package main

import (
	"fmt"
	"time"
)

// START OMIT
func slowfunc1(ch chan int) {
	time.Sleep(time.Second)
	ch <- 1
	fmt.Println("func1 done")
}

func slowfunc2(ch chan int) {
	time.Sleep(time.Second)
	ch <- 2
	fmt.Println("func2 done")
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go slowfunc1(ch1)
	go slowfunc2(ch2)
	fmt.Println(<-ch1, <-ch2)
}

// END OMIT
