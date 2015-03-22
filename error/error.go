package main

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubelet/dockertools"
)

func ContainerPodError() (string, error) {
	return "This is a value", dockertools.ErrNoContainersInPod
}

func ContainerCannotRunError() (string, error) {
	return "This is a value", dockertools.ErrContainerCannotRun
}

type Foo map[int]string

func main() {
	i := 1
	f := make(Foo)
	f[123] = "bar"
	fmt.Printf("lenght of f is %d", len(f))
	switch i {
	case 1, 2:
		fmt.Printf("one")
		return
	case 3:
		fmt.Printf("all")
		return
	}
	v, err := ContainerPodError()
	if err != nil {
		if err == dockertools.ErrContainerCannotRun {
			fmt.Printf("This is a err container cannot run.")
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf(v)
}
