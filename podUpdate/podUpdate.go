package main

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/truthninja/sandbox/kubeclient"
)

func main() {

	// Things needed for a reflector:
	//	1. kubernetes client: Interface to kubernetes apiserver.
	//		- Construct urls to REST resources that include field and label selectors
	//	2. listwatcher: Interface that knows how to list/watch a resource.
	//	3. fifo queue: Interface capable of storing the output of the listwatcher.

	kubeClient, err := kubeclient.NewKubeClient()
	if err != nil {
		fmt.Printf("Unable to create kubeclient %v", err)
		return
	}
	// Setting host is hard because of checks in the binding, set status instead.
	kubeClient.UpdatePodStatus("static-pod-from-spec", api.PodFailed)
	return
}
