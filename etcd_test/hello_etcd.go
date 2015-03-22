package main

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/coreos/go-etcd/etcd"
	"time"
)

func main() {
	c := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	kvs := map[string]string{"foo": fmt.Sprintf("%s", time.Now())}
	res, _ := c.Set("/bar/foo", kvs["foo"], 0)
	res, _ = c.Set("/bar/baz", kvs["foo"], 0)
	fmt.Printf("response %+v\n", res)

	value, _ := c.Get("foo", false, false)
	fmt.Printf("[%d] get response: %+v\n", "foo", value.Node.Value)
	testTime := util.Now()
	fmt.Printf("%s", time.Since(testTime.Time))
	//for i, res := range value {
	//}
}
