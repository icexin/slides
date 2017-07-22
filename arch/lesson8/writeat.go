package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Create("a.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < 10; i++ {
		s := strconv.Itoa(i)
		f.WriteAt([]byte(s), int64(i))
	}
}
