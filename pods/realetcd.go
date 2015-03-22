package main

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/conversion"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/tools"
	"github.com/coreos/go-etcd/etcd"
	"github.com/davecgh/go-spew/spew"
	"reflect"
)

type TheRealEtcdHelper struct {
	tools.EtcdHelper
}

func (h *TheRealEtcdHelper) theRealBodyAndExtractObj(key string, objPtr runtime.Object, ignoreNotFound bool) (body string, modifiedIndex uint64, err error) {
	response, err := h.Client.Get(key, false, false)
	spew.Dump(response)
	if err != nil && !tools.IsEtcdNotFound(err) {
		return "", 0, err
	}
	return h.theRealExtractObj(response, err, objPtr, ignoreNotFound, false)
}

func (h *TheRealEtcdHelper) theRealExtractObj(response *etcd.Response, inErr error, objPtr runtime.Object, ignoreNotFound, prevNode bool) (body string, modifiedIndex uint64, err error) {
	var node *etcd.Node
	if response != nil {
		if prevNode {
			node = response.PrevNode
		} else {
			node = response.Node
		}
	}
	if inErr != nil || node == nil || len(node.Value) == 0 {
		if ignoreNotFound {
			v, err := conversion.EnforcePtr(objPtr)
			if err != nil {
				return "", 0, err
			}
			v.Set(reflect.Zero(v.Type()))
			return "", 0, nil
		} else if inErr != nil {
			return "", 0, inErr
		}
		return "", 0, fmt.Errorf("unable to locate a value on the response: %#v", response)
	}
	body = node.Value
	err = h.Codec.DecodeInto([]byte(body), objPtr)
	if h.Versioner != nil {
		_ = h.Versioner.UpdateObject(objPtr, node)
		// being unable to set the version does not prevent the object from being extracted
	}
	return body, node.ModifiedIndex, err
}
