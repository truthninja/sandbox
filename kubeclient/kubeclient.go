package kubeclient

/*
Assortment of methods that require a kube client.
*/

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

type KubeSandbox struct {
	Client *client.Client
}

// NewKubeClient creates an insecure kubernets client to the apiserver running on localhost.
func NewKubeClient() (*KubeSandbox, error) {
	c, err := client.New(&client.Config{Host: "http://127.0.0.1:8080"})
	if err != nil {
		return nil, err
	}
	return &KubeSandbox{c}, nil
}

// Watcher creates a watch interface for the given resource
func (k *KubeSandbox) Watcher(resource string, fieldSelector fields.Selector) (watch.Interface, error) {
	return k.Client.
		Get().
		Prefix("watch").
		Namespace("default").
		Resource(resource).
		FieldsSelectorParam(
		api.FieldSelectorQueryParam(k.Client.APIVersion()),
		fieldSelector).
		Watch()
}

// WatchSinglePodEvent watches all pods for a single event
// curl http://127.0.0.1:8080/api/v1beta1/watch/pods?namespace=default
func (k *KubeSandbox) WatchSinglePodEvent() (*watch.Event, error) {
	w, err := k.Watcher("pods", fields.Everything())
	if err != nil {
		fmt.Printf("Failed to watch all pods %v\n", err)
		return nil, err
	}
	event := <-w.ResultChan()
	fmt.Printf("Watch results %#v\n", event)
	return &event, nil
}
