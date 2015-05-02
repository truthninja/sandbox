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
		FieldsSelectorParam(fieldSelector).
		Watch()
}

func (k *KubeSandbox) Poster(resource string, obj interface{}) error {
	return k.Client.
		Post().
		Namespace("default").
		Resource(resource).
		Body(obj).
		Do().
		Error()
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

func (k *KubeSandbox) GetPod(name string) (*api.Pod, error) {
	podClient := k.Client.Pods("default")
	return podClient.Get(name)
}

func (k *KubeSandbox) UpdatePod(pod *api.Pod) (*api.Pod, error) {
	podClient := k.Client.Pods("default")
	return podClient.Update(pod)
}

func (k *KubeSandbox) UpdatePodHost(name string, host string) error {
	pod, err := k.GetPod(name)
	if err != nil {
		fmt.Printf("\n Failed to get pod %#v", name)
		return err
	}
	fmt.Printf("The preudpate pod has host %+v\n", pod.Spec.Host)
	pod.Spec.Host = host
	pod, err = k.UpdatePod(pod)
	if err != nil {
		fmt.Printf("Failed to update pod %#v", pod)
	}
	fmt.Printf("The updated pod has host %v\n", pod.Spec.Host)
	return nil
}

func (k *KubeSandbox) UpdatePodStatus(name string, status api.PodPhase) error {
	//k.Client.Pods("default").UpdateStatus(name, &api.PodStatus{Phase: status})
	pod, err := k.GetPod(name)
	if err != nil {
		fmt.Printf("\n Failed to get pod %#v", name)
		return err
	}
	fmt.Printf("The preudpate pod has host %+v\n", pod.Status.Phase)
	pod.Status.Phase = status
	//pod, err = k.Client.Pods("default").UpdateStatus(name, &pod.Status)
	pod, err = k.Client.Pods("default").UpdateStatus(pod)
	if err != nil {
		fmt.Printf("Failed to update pod %#v", pod.Status.Phase)
	}
	fmt.Printf("The updated pod has host %v\n", pod.Spec.Host)
	return nil
}

func (k *KubeSandbox) BindPodHost(name string, host string) error {
	pod, err := k.GetPod(name)
	if err != nil {
		fmt.Printf("\n Failed to get pod %#v", name)
		return err
	}
	b := &api.Binding{
		ObjectMeta: api.ObjectMeta{Namespace: pod.Namespace, Name: pod.Name},
		Target: api.ObjectReference{
			Kind: "Node",
			Name: host,
		},
	}
	return k.Poster("bindings", b)
}
