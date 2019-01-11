package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mafanr/vgo/analyze/misc"

	"github.com/mafanr/g"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
)

// Etcd etcd struct,report udp listen addr
type Etcd struct {
	Client      *clientv3.Client
	Grant       *clientv3.LeaseGrantResponse
	ReportKey   string
	ReportValue string
	WatchDir    string
	StopC       chan bool
}

// NewEtcd new Etcd
func NewEtcd() *Etcd {
	etcd := &Etcd{
		StopC: make(chan bool, 1),
	}
	return etcd
}

// Start start etcd report thread
func (etcd *Etcd) Start() error {
	err := etcd.Init()
	if err != nil {
		return err
	}

	go etcd.registerWork()
	go etcd.Watch()

	return nil
}

// Close stop etcd report
func (etcd *Etcd) Close() {
	etcd.StopC <- true
	time.Sleep(1 * time.Second)

	// close channel
	close(etcd.StopC)

	if etcd.Client != nil {
		etcd.Client.Close()
	}
}

// Init init Etcd
func (etcd *Etcd) Init() error {
	reportKey, err := GetRegisterKey()
	if err != nil {
		return err
	}
	keylen := len(misc.Conf.Etcd.ReportDir)
	if keylen > 0 && misc.Conf.Etcd.ReportDir[keylen-1] != '/' {
		etcd.ReportKey = misc.Conf.Etcd.ReportDir + "/" + reportKey
	} else {
		etcd.ReportKey = misc.Conf.Etcd.ReportDir + reportKey
	}
	gAnalyze.clusterName = reportKey
	etcd.ReportValue = reportKey

	{
		wdlen := len(misc.Conf.Etcd.WatchDir)
		if wdlen > 0 && misc.Conf.Etcd.ReportDir[wdlen-1] != '/' {
			etcd.WatchDir = misc.Conf.Etcd.WatchDir + "/"
		} else {
			etcd.WatchDir = misc.Conf.Etcd.WatchDir
		}
	}

	g.L.Info("Init", zap.String("@Key", etcd.ReportKey), zap.String("@Value", etcd.ReportValue))

	cfg := clientv3.Config{
		Endpoints:   misc.Conf.Etcd.Addrs,
		DialTimeout: time.Duration(misc.Conf.Etcd.Dltimeout) * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return err
	}

	etcd.Client = client

	return nil
}

// registerWork register stream addr
func (etcd *Etcd) registerWork() {
	timeC := time.Tick(time.Duration(misc.Conf.Etcd.ReportTime) * time.Second)
	// 启动立刻注册一次
	etcd.Put(etcd.ReportKey, etcd.ReportValue, misc.Conf.Etcd.TTL)
	for {
		select {
		case <-etcd.StopC:
			g.L.Info("Etcd", zap.String("Close", "Ok"))
			return
		// Timing task
		case <-timeC:
			etcd.Put(etcd.ReportKey, etcd.ReportValue, misc.Conf.Etcd.TTL)
			// g.L.Debug("register", zap.String("@Key", etcd.ReportKey), zap.String("addr", etcd.ReportValue))
			break
		}
	}
}

// Put put
func (etcd *Etcd) Put(key, value string, ttl int64) error {
	Grant, err := etcd.Client.Grant(context.TODO(), ttl)
	if err != nil {
		g.L.Error("Etcd", zap.Error(err), zap.Int64("@ReportTime", ttl))
		return err
	}
	_, err = etcd.Client.Put(context.TODO(), key, value, clientv3.WithLease(Grant.ID))
	if err != nil {
		g.L.Error("Put", zap.String("@key", key), zap.String("@value", value), zap.Error(err))
		return err
	}
	return nil
}

// GetRegisterKey get key
func GetRegisterKey() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d", host, os.Getpid()), nil
}

// Watch watch
func (etcd *Etcd) Watch() {
	for {
		rch := etcd.Client.Watch(context.Background(), etcd.WatchDir, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				// 上报
				if ev.Type == 0 {
					if _, ok := gAnalyze.analyzes[string(ev.Kv.Key)]; !ok {
						gAnalyze.hash.Add(string(ev.Kv.Value))
						gAnalyze.analyzes[string(ev.Kv.Key)] = string(ev.Kv.Value)
						g.L.Info("Notify Join", zap.String("Key", string(ev.Kv.Key)), zap.String("Value", string(ev.Kv.Value)))
					}
				} else {
					if value, ok := gAnalyze.analyzes[string(ev.Kv.Key)]; ok {
						gAnalyze.hash.Remove(value)
						delete(gAnalyze.analyzes, string(ev.Kv.Key))
						g.L.Info("Notify Leave", zap.String("Key", string(ev.Kv.Key)), zap.String("Value", value))
					}
				}
			}
		}
	}
}
