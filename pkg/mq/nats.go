package mq

import (
	"fmt"

	"github.com/nats-io/nats"
	"go.uber.org/zap"
)

// Nats nats struct
type Nats struct {
	addrs  []string
	topic  string
	conn   *nats.Conn
	logger *zap.Logger
}

// NewNats return new nats
func NewNats() *Nats {
	return &Nats{}
}

// Start init && start nats
func (n *Nats) Start(addrs []string, topic string, logger *zap.Logger) error {
	n.addrs = addrs
	n.topic = topic
	n.logger = logger
	if err := n.start(); err != nil {
		n.logger.Warn("nats start", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// Close close nats
func (n *Nats) Close() error {
	err := n.Close()
	if err != nil {
		n.logger.Warn("nats close", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// start start nats
func (n *Nats) start() error {
	opts := nats.DefaultOptions
	opts.Servers = n.addrs
	nc, err := opts.Connect()
	if err != nil {
		n.logger.Warn("nats connect error", zap.String("error", err.Error()), zap.Strings("addrs", n.addrs))
		return err
	}

	// Setup callbacks to be notified on disconnects and reconnects
	nc.Opts.DisconnectedCB = func(nc *nats.Conn) {
		// log.Printf("%v got disconnected!\n", nc.ConnectedUrl())
	}

	// See who we are connected to on reconnect.
	nc.Opts.ReconnectedCB = func(nc *nats.Conn) {
		// log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
	}
	n.logger.Debug("Nats", zap.String("Topic", n.topic))

	n.conn = nc
	if nc.IsClosed() == true {
		n.logger.Warn("nats is closed", zap.String("error", err.Error()), zap.Strings("addrs", n.addrs))
		return fmt.Errorf("nats is closed")
	}
	return nil
}

// Subscribe ....
func (n *Nats) Subscribe(topic string, handler func(msg *nats.Msg)) error {
	// 普通订阅
	_, err := n.conn.Subscribe(topic, handler)
	if err != nil {
		n.logger.Warn("nats subscribe error", zap.String("error", err.Error()), zap.Strings("addrs", n.addrs))
		n.conn.Close()
		return err
	}
	return nil
}

// QueueSubscribe ....
func (n *Nats) QueueSubscribe(topic, queue string, handler func(msg *nats.Msg)) error {
	// 普通订阅
	_, err := n.conn.QueueSubscribe(topic, queue, handler)
	if err != nil {
		n.logger.Warn("nats subscribe error", zap.String("error", err.Error()), zap.Strings("addrs", n.addrs))
		n.conn.Close()
		return err
	}
	return nil
}

// Publish 发布
func (n *Nats) Publish(topic string, data []byte) error {
	return n.conn.Publish(topic, data)
}
