package main

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/conversion"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/master"
	etcdgeneric "github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic/etcd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/GoogleCloudPlatform/kubernetes/plugin/cmd/kube-scheduler/app"
	"github.com/GoogleCloudPlatform/kubernetes/plugin/pkg/scheduler"
	"github.com/GoogleCloudPlatform/kubernetes/plugin/pkg/scheduler/factory"
	"github.com/coreos/go-etcd/etcd"
	//	"github.com/davecgh/go-spew/spew"
	"reflect"
	"time"
)

const VERBOSE = false

// Dumb example to modify the given pod
func pods(kubeClient *client.Client, podName, podNamespace string) {
	if podName == "" {
		podName = "my-nginx-fplln"
	}
	if podNamespace == "" {
		podNamespace = "default"
	}
	podsClient := kubeClient.Pods(podNamespace)
	fmt.Printf("Kubeclient %+v, pods %+v", kubeClient, podsClient)
	pod, _ := podsClient.Get(podName)
	pod.Status.Phase = "Failed"
	pod, _ = podsClient.Update(pod)
	fmt.Printf("\n\nPod %+v\n\n", pod)
}

func rc(kubeClient *client.Client) {
	rcClient := kubeClient.ReplicationControllers("default")
	rc, _ := rcClient.Get("my-nginx")
	printRc(rc)

	obj := createEmptyObj(&api.ReplicationController{})

	ctx := api.NewDefaultContext()
	key, _ := etcdgeneric.NamespaceKeyFunc(ctx, "/registry/controllers", "my-nginx")
	fmt.Println("Entire etcd key ", key)

	helper, _ := master.NewEtcdHelper(etcd.NewClient(util.StringList{"http://127.0.0.1:4001"}), "")
	h := TheRealEtcdHelper{helper}
	_, index, _ := h.theRealBodyAndExtractObj(key, obj, false)
	fmt.Println("Retrieved object index ", index)
	prevRc := obj.(*api.ReplicationController)
	printRc(prevRc)
}

func newKubeClient(s *app.SchedulerServer) *client.Client {
	s = app.NewSchedulerServer()
	s.ClientConfig.Host = "http://127.0.0.1:8080"
	kubeClient, err := client.New(&s.ClientConfig)
	if err != nil {
		fmt.Printf("Failed: %+v", err)
		return nil
	}
	return kubeClient
}

func newSchedulerConfig(kubeClient *client.Client, s *app.SchedulerServer) *scheduler.Config {
	configFactory := factory.NewConfigFactory(kubeClient)
	config, _ := configFactory.CreateFromProvider("DefaultProvider")
	fmt.Println(configFactory.PodQueue)
	fmt.Println(configFactory.PodLister)
	fmt.Println(configFactory.NodeLister)
	return config
}

func createEmptyObj(newType runtime.Object) runtime.Object {
	v, _ := conversion.EnforcePtr(newType)
	reflectType := reflect.New(v.Type())
	iface := reflectType.Interface()
	obj := iface.(runtime.Object)
	if VERBOSE {
		fmt.Println("===Parsing the line reflect.New(v.Type()).Interface().(runtime.Object)===")
		fmt.Printf(" Conversion of &api.ReplicationController returned %+v\n", v)
		fmt.Printf(" Reflect.new(%+v) = %+v \n .Interface() type %T \n final object type = %T\n",
			v.Type(), reflectType, iface, iface.(runtime.Object))
		fmt.Println("===Parsing the line reflect.New(v.Type()).Interface().(runtime.Object)===")
	}
	return obj
}

func printRc(rc *api.ReplicationController) {
	fmt.Printf("\nRetrieved an rc name: %+v version %+v status.replicas %+v labels %q \n",
		rc.Name, rc.ResourceVersion, rc.Status.Replicas, rc.Labels)
}

func main() {
	s := app.NewSchedulerServer()
	kubeClient := newKubeClient(s)

	// pods(kubeClient, "", "")
	// rc(kubeClient)
	c := newSchedulerConfig(kubeClient, s)
	scheduler := scheduler.New(c)
	fmt.Printf("WAiting on next pod for scheduler %+v", scheduler)
	for i := 0; ; i++ {
		pod := c.NextPod()
		fmt.Printf("\n[%v] delay %v", pod.Name, time.Since(pod.CreationTimestamp.Time))
	}
}
