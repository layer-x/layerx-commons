package lxdatabase

import (
	"github.com/coreos/etcd/client"
	"github.com/layer-x/layerx-commons/lxerrors"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var c client.Client
var m sync.Mutex

func Init(etcdEndpoints []string) error {
	cfg := client.Config{
		Endpoints:               etcdEndpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	var err error
	c, err = client.New(cfg)
	if err != nil {
		return lxerrors.New("initialize etcd", err)
	}
	m = &sync.Mutex{}
	return nil
}

func Get(key string) (string, error) {
	m.Lock()
	defer m.Unlock()
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Get(context.Background(), key, nil)
	if err != nil {
		return "", lxerrors.New("getting key/val", err)
	}
	return resp.Node.Value
}

func Set(key string, value string) error {
	m.Lock()
	defer m.Unlock()
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Set(context.Background(), key, value, nil)
	if err != nil {
		return lxerrors.New("setting key/val pair", err)
	}
	if resp.Node.Key != key || resp.Node.Value != value {
		return lxerrors.New("key/value pair not set as expected", nil)
	}
	return nil
}

func Rm(key string) error {
	m.Lock()
	defer m.Unlock()
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Delete(context.Background(), key, nil)
	if err != nil {
		return lxerrors.New("deleting key/val pair", err)
	}
	if resp.Node.Key != key {
		return lxerrors.New("removed pair does not have expected key", nil)
	}
	return nil
}

func Mkdir(dir string) error {
	m.Lock()
	defer m.Unlock()
	kapi := client.NewKeysAPI(c)
	opts := client.SetOptions{
		Dir: true,
	}
	resp, err := kapi.Set(context.Background(), dir, nil, opts)
	if err != nil {
		return lxerrors.New("making directory", err)
	}
	if resp.Node.Key != dir || !resp.Node.Dir {
		return lxerrors.New("directory not created as expected", nil)
	}
	return nil
}

func Rmdir(dir string) error {
	m.Lock()
	defer m.Unlock()
	kapi := client.NewKeysAPI(c)
	opts := client.SetOptions{
		Dir: true,
	}
	resp, err := kapi.Set(context.Background(), dir, nil, opts)
	if err != nil {
		return lxerrors.New("removing directory", err)
	}
	if resp.Node.Key != dir || !resp.Node.Dir {
		return lxerrors.New("directory not created as expected", nil)
	}
	return nil
}

func Ls(dir string) (map[string]string, error) {
	m.Lock()
	defer m.Unlock()
	kapi := client.NewKeysAPI(c)
	resp, err := kapi.Get(context.Background(), dir, nil)
	if err != nil {
		return "", lxerrors.New("getting key/vals for dir", err)
	}
	if !resp.Node.Dir {
		return "", lxerrors.New("ls used on a non-dir key", err)
	}
	result := make(map[string]string)
	for _, node := range resp.Node.Nodes {
		if !node.Dir {
			result[node.Key] = node.Value
		} //ignore directories
	}
	return result
}
