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
	str := `{"Id":2,"Name":"alice"}`
	var s Student
	err := json.Unmarshal([]byte(str), &s)
	if err != nil {
		log.Fatalf("unmarshal error:%s", err)
	}
	fmt.Println(s)
}
