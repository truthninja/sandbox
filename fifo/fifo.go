package main

/*
This is the main package for an etcd watch event fifo manager (aka a reflector,
or ewefm). What follows is the spec:
	1. It is capable of watching etcd
		* Needs an embedded etcd client
		- It's capable of watching particular fields alone
			* Needs an embedded apiserver client instead of a raw etcd client.
	2. It is capable of enqueuing watch notifications in fifo order
		- fifo-ing is based on resourceVersion, or etcd's modifiedIndex
			* Needs to deserialize watch event and retrieve resourceVersion
		- It can't hold up the main thread on the watch
			* Needs something that decodes bits received over a streaming
				connection into watch events and shoves them into a channel
				the invoker of the watch can select over.
		- It can't leak the stream watcher
			* 2 channels are needed for bidirectional communication. One way
				for stop signals to the stream watcher, and the other for
				stream events.
	3. The fifo needs semantics that avoid an overwrite
		- Needs to check if an element exists in the queue in constant time
			* Needs a hashmap as well as a list.
			* Needs a keying strategy for the hash map.
			* Needs to go between elemnts popped from the queue and the map.
	4. Questions for the queue itself:
		- What is the size?
		- How does it behave once it's full?
		- Can one requeue elements?
			* The item is removed from the queue and the store on Pop, so one needs
				to requeue it with Add if processing fails.
		- How does one index it?

A couple of assumptions in this package:
	1. You already have an apiserver running on localhost and listening on 8080.
	2. You already have etcd instance running and listening on some port (like 4001),
		that the apiserver is aware of. Since this package proxies everything
		through the apiserver, the etcd port doesn't matter as long as the apiserver
		works. It won't work without etcd.
	The easiest way to achieve this is by running local-up-cluster.sh from the
	kubernetes project.
*/

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/cache"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/truthninja/sandbox/kubeclient"
	"time"
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
	listwatcher := cache.NewListWatchFromClient(
		kubeClient.Client, "pods", api.NamespaceAll,
		fields.Set{"DesiredState.Host": ""}.AsSelector())
	fifo := cache.NewFIFO(cache.MetaNamespaceKeyFunc)
	reflector := cache.NewReflector(listwatcher, &api.Pod{}, fifo, 0)
	reflector.Run()
	for {
		fmt.Println("Listing fifo")
		fmt.Println(fifo.List())
		time.Sleep(2 * time.Second)
	}
}
