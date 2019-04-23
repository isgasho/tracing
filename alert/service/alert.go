package service

import (
	"github.com/gocql/gocql"
	"github.com/imdevlab/tracing/alert/misc"
	"github.com/imdevlab/tracing/pkg/mq"
	"github.com/nats-io/nats"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Alert 告警服务
type Alert struct {
	mq      *mq.Nats
	cql     *gocql.Session
	policys *Policys
}

var gAlert *Alert

// New new alert
func New(l *zap.Logger) *Alert {
	logger = l
	gAlert = &Alert{
		mq:      mq.NewNats(logger),
		policys: newPolicys(),
	}
	return gAlert
}

// Start start server
func (a *Alert) Start() error {
	// 初始化db
	if err := a.initDB(); err != nil {
		logger.Warn("init db error", zap.String("error", err.Error()))
		return err
	}
	// 加载各种策略
	if err := a.policys.Start(); err != nil {
		logger.Warn("policy start error", zap.String("error", err.Error()))
		return err
	}
	// 启动消息队列服务
	if err := a.mq.Start(misc.Conf.MQ.Addrs, misc.Conf.MQ.Topic); err != nil {
		logger.Warn("mq start error", zap.String("error", err.Error()))
		return err
	}

	// 订阅mq信息
	if err := a.mq.Subscribe(misc.Conf.MQ.Topic, msgHandle); err != nil {
		logger.Warn("mq subscribe error", zap.String("error", err.Error()))
		return err
	}

	logger.Info("start ok", zap.Any("config", misc.Conf))
	return nil
}

// Close stop server
func (a *Alert) Close() error {
	return nil
}

func msgHandle(msg *nats.Msg) {
	logger.Info("msgHandle", zap.String("msg", string(msg.Data)))
}

// init 初始化存储
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
