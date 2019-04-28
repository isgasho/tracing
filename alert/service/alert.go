package service

import (
	"github.com/gocql/gocql"
	"github.com/imdevlab/tracing/alert/misc"
	"github.com/imdevlab/tracing/alert/ticker"
	"github.com/nats-io/nats"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Alert 告警服务
type Alert struct {
	apps    *Apps
	cql     *gocql.Session
	tickers *ticker.Tickers
}

var gAlert *Alert

// New new alert
func New(l *zap.Logger) *Alert {
	logger = l
	gAlert = &Alert{
		apps:    newApps(),
		tickers: ticker.NewTickers(10, misc.Conf.Analyze.Interval, logger),
	}
	return gAlert
}

// Start start server
func (a *Alert) Start() error {

	if err := a.initDB(); err != nil {
		logger.Warn("int db", zap.String("error", err.Error()))
		return err
	}

	if err := a.apps.start(); err != nil {
		logger.Warn("apps start", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// Close stop server
func (a *Alert) Close() error {
	return nil
}

func msgHandle(msg *nats.Msg) {
	logger.Info("msgHandle", zap.String("msg", string(msg.Data)))
}

// initDB 初始化存储
func (a *Alert) initDB() error {
	// connect to the cluster
	cluster := gocql.NewCluster(misc.Conf.DB.Cluster...)
	cluster.Keyspace = misc.Conf.DB.Keyspace
	cluster.Consistency = gocql.Quorum
	//设置连接池的数量,默认是2个（针对每一个host,都建立起NumConns个连接）
	cluster.NumConns = misc.Conf.DB.NumConns

	session, err := cluster.CreateSession()
	if err != nil {
		logger.Warn("create session", zap.String("error", err.Error()))
		return err
	}
	a.cql = session
	return nil
}

// GetCql ...
func (a *Alert) GetCql() *gocql.Session {
	if a.cql != nil {
		return a.cql
	}
	return nil
}
