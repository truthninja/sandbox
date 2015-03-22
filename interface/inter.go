package main

import (
	"fmt"
)

type Human interface {
	Name() string
}

// Both android and person implement the human interface.
type Person struct {
	name string
}

func (p Person) Name() string {
	return p.name
}

type Android struct {
	Person
	model string
	name  string
}

func (a Android) Model() string {
	return a.model
}

func duck_type(p Person) Human {
	return p
}

func main() {
	p := &Person{name: "beeps"}
	a := &Android{model: "prashanth", name: "foobar"}
	//var h Human
	//h = p
	h := duck_type(*p)
	fmt.Printf("\nThis is h: %+v", h)
	fmt.Printf("\nThis is person %s", p.Name())
	fmt.Printf("\nThis is android %s", a.Model())
	fmt.Printf("\nThis is the person in the android %s\n", h.Name())
}
