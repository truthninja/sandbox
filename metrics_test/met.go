package main

import (
	"fmt"
	kube_metrics "github.com/GoogleCloudPlatform/kubernetes/pkg/metrics"
	//"github.com/rcrowley/go-metrics"
	//"log"
	//h"net"
	//	"os"
	"time"
)

func main() {
	kube_metrics.Init()
	for {
		fmt.Printf("%s\n", time.Now())
		kube_metrics.Count(1, "test_count")
		time.Sleep(1 * time.Second)
	}
}
