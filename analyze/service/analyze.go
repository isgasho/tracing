package service

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"
	"github.com/mafanr/vgo/analyze/service/blink"
	"github.com/mafanr/vgo/analyze/service/stats"
	"go.uber.org/zap"
)

// Analyze ...
type Analyze struct {
	db       *g.Cassandra
	stats    *stats.Stats
	blink    *blink.Blink
	appStore *AppStore
}

// New ...
func New() *Analyze {
	analyze := &Analyze{
		stats: stats.New(),
		blink: blink.New(),
		db:    g.NewCassandra(),
	}

	return analyze
}

// Start ...
func (analyze *Analyze) Start() error {

	g.L.Info("Conf", zap.Any("conf", misc.Conf))

	if err := analyze.db.Init(misc.Conf.Cassandra.NumConns, misc.Conf.Cassandra.Keyspace, misc.Conf.Cassandra.Cluster); err != nil {
		g.L.Fatal("Start Init", zap.String("error", err.Error()))
	}

	appStore := NewAppStore(analyze.db)
	analyze.appStore = appStore

	if err := analyze.appStore.Start(); err != nil {
		g.L.Fatal("Start appStore", zap.String("error", err.Error()))
	}

	if err := analyze.blink.Start(); err != nil {
		g.L.Fatal("Start blink", zap.String("error", err.Error()))
	}

	if err := analyze.stats.Start(); err != nil {
		g.L.Fatal("Start stats", zap.String("error", err.Error()))
	}

	g.L.Info("Start ok!")
	return nil
}

// Close ...
func (analyze *Analyze) Close() error {

	if analyze.blink != nil {
		analyze.blink.Close()
	}

	if analyze.stats != nil {
		analyze.stats.Close()
	}

	if analyze.db != nil {
		analyze.db.Close()
	}

	g.L.Info("Close ok!")
	return nil
}
