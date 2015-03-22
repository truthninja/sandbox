package main

import (
	"fmt"
)

//	"github.com/golang/glog"
type Person struct {
	Name  string
	Numbr int64
}

func (p *Person) String() string {
	return p.Name
}

func main() {
	p := &Person{Name: "vim-go"}
	glog.V(1).Infof("This is a print statement")
	fmt.Printf("p %s\n", p)
}
