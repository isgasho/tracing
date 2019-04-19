package service

import (
	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/alert/misc"
	"github.com/imdevlab/tracing/pkg/mq"
	"github.com/nats-io/nats"
	"go.uber.org/zap"
)

// Alert 告警服务
type Alert struct {
	mq *mq.Nats
}

// New new alert
func New() *Alert {
	return &Alert{
		mq: mq.NewNats(),
	}
}

// Start start server
func (a *Alert) Start() error {
	// 初始化sql
	// 加载各种策略
	// 启动消息队列服务
	if err := a.mq.Start(misc.Conf.MQ.Addrs, misc.Conf.MQ.Topic, g.L); err != nil {
		g.L.Warn("mq start error", zap.String("error", err.Error()))
		return err
	}

	// 订阅mq信息
	if err := a.mq.Subscribe(misc.Conf.MQ.Topic, msgHandle); err != nil {
		g.L.Warn("mq subscribe error", zap.String("error", err.Error()))
		return err
	}

	g.L.Info("start ok", zap.Any("config", misc.Conf))
	return nil
}

// Close stop server
func (a *Alert) Close() error {
	return nil
}

func msgHandle(msg *nats.Msg) {
	g.L.Info("msgHandle", zap.String("msg", string(msg.Data)))
}
