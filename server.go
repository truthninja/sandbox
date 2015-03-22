package main

import (
	"encoding/json"
	//"errors"
	"io"
	//"net"
	"fmt"
	"net/http"
	//"net/url"
	"path"
	//"strconv"
	"strings"
	//"time"

	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api/latest"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/httplog"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/credentialprovider"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubelet"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubelet/dockertools"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubelet/server"
	"github.com/golang/glog"
	"github.com/google/cadvisor/info"
)

func handler(w http.ResponseWriter, req *http.Request) {
	//glog.V(1).Infof("Simple http request: %+v", req)
	components := strings.Split(strings.TrimPrefix(path.Clean(req.URL.Path), "/"), "/")
	glog.V(1).Infof("Components of stats %+v", components)
	//var stats *info.ContainerInfo
	var err error
	var query info.ContainerInfoRequest
	err = json.NewDecoder(req.Body).Decode(&query)
	if err != nil && err != io.EOF {
		fmt.Println(w, err)
		return
	}

	glog.V(1).Infof("Query for pod with namespace %+v and component %+v", components[1], components[2])
	/*
		pod, ok := s.host.GetPodByName(components[1], components[2])
		if !ok {
			http.Error(w, "Pod does not exist", http.StatusNotFound)
			return
		}
		glog.V(1).Infof("Query for pod with uid")
		stats, err = s.host.GetContainerInfo(GetPodFullName(pod), types.UID(components[3]), components[4], &query)
	*/

}

func main() {
	s := server.NewKubeletServer()
	fmt.Printf("New kubelet server %+v", s)
	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8001", nil)
	//credentialprovider.SetPreferredDockercfgPath(s.RootDirectory)

	kcfg := server.KubeletConfig{
		Address:                 s.Address,
		AllowPrivileged:         s.AllowPrivileged,
		HostnameOverride:        s.HostnameOverride,
		RootDirectory:           s.RootDirectory,
		ConfigFile:              s.Config,
		ManifestURL:             s.ManifestURL,
		FileCheckFrequency:      s.FileCheckFrequency,
		HTTPCheckFrequency:      s.HTTPCheckFrequency,
		PodInfraContainerImage:  s.PodInfraContainerImage,
		SyncFrequency:           s.SyncFrequency,
		RegistryPullQPS:         s.RegistryPullQPS,
		RegistryBurst:           s.RegistryBurst,
		MinimumGCAge:            s.MinimumGCAge,
		MaxContainerCount:       s.MaxContainerCount,
		ClusterDomain:           s.ClusterDomain,
		ClusterDNS:              s.ClusterDNS,
		Runonce:                 true, //s.RunOnce,
		Port:                    s.Port,
		CAdvisorPort:            s.CAdvisorPort,
		EnableServer:            s.EnableServer,
		EnableDebuggingHandlers: s.EnableDebuggingHandlers,
		DockerClient:            dockertools.ConnectToDockerOrDie(s.DockerEndpoint),
		KubeClient:              client.New(client.Config{}),
		EtcdClient:              kubelet.EtcdClientOrDie(s.EtcdServerList, s.EtcdConfigFile),
		MasterServiceNamespace:  s.MasterServiceNamespace,
		VolumePlugins:           server.ProbeVolumePlugins(),
	}
	fmt.Printf(kcfg)
	sever.RunKubelet(&kcfg)
}
