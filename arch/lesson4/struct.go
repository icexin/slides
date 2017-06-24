package main

import "fmt"

type Student struct {
	Id   int
	Name string
}

func main() {
	var s Student
	s.Id = 1
	s.Name = "jack"
	fmt.Println(s)

	s1 := Student{
		Id:   2,
		Name: "alice",
	}
	fmt.Println(s1)

	s1 = s
	fmt.Println(s1)
}
