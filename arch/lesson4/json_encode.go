package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Student struct {
	Id   int
	Name string
}

func main() {
	s := Student{
		Id:   2,
		Name: "alice",
	}
	buf, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("marshal error:%s", err)
	}
	fmt.Println(string(buf))
}
