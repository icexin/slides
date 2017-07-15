package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(1000 * time.Millisecond)
	boom := time.After(5000 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("滴答...")
		case <-boom:
			fmt.Println("嘣!!!")
			return
		default:
			fmt.Println("吃一口面")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
