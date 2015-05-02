package main

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/coreos/go-etcd/etcd"
	"time"
)

// If recursive is set to true the watch returns the first change under the given
// prefix since the given index.
//
// If recursive is set to false the watch returns the first change to the given key
// since the given index.
//
// To watch for the latest change, set waitIndex = 0.
//
// If a receiver channel is given, it will be a long-term watch. Watch will block at the
//channel. After someone receives the channel, it will go on to watch that
// prefix.  If a stop channel is given, the client can close long-term watch using
// the stop channel.
func main() {
	c := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	kvs := map[string]string{"foo": fmt.Sprintf("%s", time.Now())}
	ch := make(chan *etcd.Response)
	res, _ := c.Set("/bar/foo", kvs["foo"], 0)
	go func() {
		r, e := c.Watch("/bar", 0, true, ch, nil)
		if e != nil {
			fmt.Printf("\nFailed to watch %+v", e)
		}
		fmt.Printf("\nresponse %+v\n", r)
	}()
	res, _ = c.Set("/bar/baz", kvs["foo"], 0)
	fmt.Printf("response %+v\n", res)
	for {
		select {
		case v, _ := <-ch:
			fmt.Printf("\nGot watch update %+v", v)
		case <-time.After(10 * time.Minute):
			fmt.Printf("\nTimeout\n")
		}
	}
	value, _ := c.Get("/bar/foo", false, false)
	fmt.Printf("[%d] get response: %+v\n", "foo", value.Node.Value)
	testTime := util.Now()
	fmt.Printf("%s", time.Since(testTime.Time))
}
