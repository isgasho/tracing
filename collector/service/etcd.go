package service

import (
	"context"
	"time"

	"github.com/imdevlab/g"
	"github.com/shaocongcong/tracing/collector/misc"
	"go.etcd.io/etcd/clientv3"

	"go.uber.org/zap"
)

// Etcd struct,report udp listen addr
type Etcd struct {
	Client      *clientv3.Client
	Grant       *clientv3.LeaseGrantResponse
	ReportKey   string
	ReportValue string
	StopC       chan bool
}

// newEtcd new Etcd
func newEtcd() *Etcd {
	etcd := &Etcd{
		StopC: make(chan bool, 1),
	}
	return etcd
}

// Start start report thread
func (e *Etcd) Start() error {
	go e.registerWork()
	return nil
}

// Close stop ereport
func (e *Etcd) Close() {
	e.StopC <- true
	time.Sleep(1 * time.Second)

	// close channel
	close(e.StopC)

	if e.Client != nil {
		e.Client.Close()
	}
}

// Init init Etcd
func (e *Etcd) Init(addrs []string, reportKey, reportValue string) error {
	e.ReportKey = reportKey
	e.ReportValue = reportValue

	g.L.Info("Init", zap.String("@Key", e.ReportKey), zap.String("@Value", e.ReportValue))

	cfg := clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: time.Duration(misc.Conf.Etcd.TimeOut) * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return err
	}

	e.Client = client

	return nil
}

// registerWork register stream addr
func (e *Etcd) registerWork() {
	timeC := time.Tick(time.Duration(misc.Conf.Etcd.ReportTime) * time.Second)
	// 启动立刻注册一次
	e.Put(e.ReportKey, e.ReportValue, misc.Conf.Etcd.TTL)
	for {
		select {
		case <-e.StopC:
			g.L.Info("Etcd", zap.String("Close", "Ok"))
			return
		// Timing task
		case <-timeC:
			if err := e.Put(e.ReportKey, e.ReportValue, misc.Conf.Etcd.TTL); err != nil {
				g.L.Warn("Etcd", zap.String("error", err.Error()))
			}
			// g.L.Debug("register", zap.String("@Key", e.ReportKey), zap.String("addr", e.ReportValue))
			break
		}
	}
}

// Put put
func (e *Etcd) Put(key, value string, ttl int64) error {
	Grant, err := e.Client.Grant(context.TODO(), ttl)
	if err != nil {
		g.L.Error("Etcd", zap.Error(err), zap.Int64("@ReportTime", ttl))
		return err
	}
	_, err = e.Client.Put(context.TODO(), key, value, clientv3.WithLease(Grant.ID))
	if err != nil {
		g.L.Error("Put", zap.String("@key", key), zap.String("@value", value), zap.Error(err))
		return err
	}
	return nil
}
