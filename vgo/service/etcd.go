package service

import (
	"context"
	"log"
	"time"

	"github.com/mafanr/vgo/vgo/misc"

	"go.etcd.io/etcd/clientv3"
)

// Etcd ...
type Etcd interface {
	Start() error
	Close()
	Watch()
}

// etcd etcd struct,report udp listen addr
type etcd struct {
	Client *clientv3.Client
	Grant  *clientv3.LeaseGrantResponse
	StopC  chan bool
}

// var etcd *etcd

// newEtcd new etcd
func newEtcd() *etcd {
	etcd := &etcd{
		StopC: make(chan bool, 1),
	}

	return etcd
}

// Start start etcd report thread
func (etcd *etcd) Start() error {
	cfg := clientv3.Config{
		Endpoints:   misc.Conf.Etcd.Addrs,
		DialTimeout: time.Duration(misc.Conf.Etcd.Dltimeout) * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return err
	}
	etcd.Client = client
	go etcd.Watch()
	return nil
}

// Close stop etcd report
func (etcd *etcd) Close() {
	etcd.StopC <- true
	time.Sleep(1 * time.Second)
	close(etcd.StopC)
}

// Watch watch
func (etcd *etcd) Watch() {
	for {
		rch := etcd.Client.Watch(context.Background(), misc.Conf.Etcd.WatchKey, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				// 上报
				if ev.Type == 0 {
					log.Println(string(ev.Kv.Key), string(ev.Kv.Value))
				} else {
					log.Println(string(ev.Kv.Key))
				}
			}
		}
	}
}
