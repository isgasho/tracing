package service

import (
	"context"
	"time"

	"github.com/mafanr/vgo/analyze/misc"

	"github.com/mafanr/g"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
)

// SrvDiscovery ...
type SrvDiscovery interface {
	Init(reportKey, reportValue, watchDir string) error
	Start() error
	Close()
	Put(key, value string, ttl int64) error
	Watch()
}

// Etcd etcd struct,report udp listen addr
type etcd struct {
	Client      *clientv3.Client
	Grant       *clientv3.LeaseGrantResponse
	ReportKey   string
	ReportValue string
	WatchDir    string
	StopC       chan bool
}

// newEtcd new Etcd
func newEtcd() *etcd {
	etcd := &etcd{
		StopC: make(chan bool, 1),
	}
	return etcd
}

// Start start etcd report thread
func (etcd *etcd) Start() error {
	go etcd.registerWork()
	go etcd.Watch()
	return nil
}

// Close stop etcd report
func (etcd *etcd) Close() {
	etcd.StopC <- true
	time.Sleep(1 * time.Second)

	// close channel
	close(etcd.StopC)

	if etcd.Client != nil {
		etcd.Client.Close()
	}
}

// Init init Etcd
func (etcd *etcd) Init(reportKey, reportValue, watchDir string) error {
	etcd.ReportKey = reportKey
	etcd.ReportValue = reportValue
	etcd.WatchDir = watchDir

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
func (etcd *etcd) registerWork() {
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
			if err := etcd.Put(etcd.ReportKey, etcd.ReportValue, misc.Conf.Etcd.TTL); err != nil {
				g.L.Warn("etcd", zap.String("error", err.Error()))
			}
			g.L.Debug("register", zap.String("@Key", etcd.ReportKey), zap.String("addr", etcd.ReportValue))
			break
		}
	}
}

// Put put
func (etcd *etcd) Put(key, value string, ttl int64) error {
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

// Watch watch
func (etcd *etcd) Watch() {
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
