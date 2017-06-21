package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	infos, err := f.Readdir(-1)
	names, err := f.Readdirnames(-1)
}
