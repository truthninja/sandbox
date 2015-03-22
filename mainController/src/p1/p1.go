package p1

import (
	"errors"
)

var (
	Err1 = errors.New("This is a new error")
)

type Foo struct {
	arob int
	AROB string
}

type P1 interface{}
