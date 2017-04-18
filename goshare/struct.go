package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s %d", p.Name, p.Age)
}

type Stringer interface {
	String() string
}

func main() {
	var p = Person{"icexin", 20}
	var s Stringer = p
	fmt.Println(s.String())
}
