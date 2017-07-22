package main

import (
	"fmt"
	"time"
)

func print(s string, sem chan int) {
	for {
		sem <- 1
		fmt.Println(s)
		time.Sleep(time.Second)
		<-sem
	}
}

func main() {
	sem := make(chan int, 1)
	go print("A", sem)
	go print("B", sem)
	select {}
}
