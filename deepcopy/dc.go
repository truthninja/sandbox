package main

import (
	"code.google.com/p/rog-go/exp/deepcopy"
	"fmt"
)

type embed struct {
	a int
}

type ifr struct {
	*embed
}

func main() {
	i := 10
	a := &ifr{&embed{i}}
	b := a
	c := *a
	d := deepcopy.Copy(a)
	a.a = 100
	fmt.Printf("\na %+v, b %+v, c %+v d %+v a %T b %T c %T d %T\n", a, b, c, d, a, b, c, d)
}
