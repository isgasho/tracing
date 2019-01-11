package service

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"

	"go.uber.org/zap"
)

// Analyze ...
type Analyze struct {
	cql         *g.Cassandra
	stats       *Stats
	blink       *Blink
	appStore    *AppStore
	hash        *g.Hash
	etcd        *Etcd
	analyzes    map[string]string
	clusterName string
	// cluster  *Cluster
}

var gAnalyze *Analyze

// New ...
func New() *Analyze {
	analyze := &Analyze{
		stats:    NewStats(),
		blink:    NewBlink(),
		cql:      g.NewCassandra(),
		hash:     g.NewHash(),
		etcd:     NewEtcd(),
		analyzes: make(map[string]string),
		// cluster: NewCluster(),
	}
	gAnalyze = analyze
	return analyze
}

// Start ...
func (analyze *Analyze) Start() error {

	g.L.Info("Conf", zap.Any("conf", misc.Conf))

	if err := analyze.etcd.Start(); err != nil {
		g.L.Fatal("Start etcd.Start", zap.String("error", err.Error()))
	}

	// if err := analyze.cluster.Start(); err != nil {
	// 	g.L.Fatal("Start cluster.Start", zap.String("error", err.Error()))
	// }

	if err := analyze.cql.Init(misc.Conf.Cassandra.NumConns, misc.Conf.Cassandra.Keyspace, misc.Conf.Cassandra.Cluster); err != nil {
		g.L.Fatal("Start Init", zap.String("error", err.Error()))
	}

	appStore := NewAppStore(analyze.cql)
	analyze.appStore = appStore

	if err := analyze.appStore.Start(); err != nil {
		g.L.Fatal("Start appStore", zap.String("error", err.Error()))
	}

	if err := analyze.blink.Start(); err != nil {
		g.L.Fatal("Start blink", zap.String("error", err.Error()))
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

	if analyze.cql != nil {
		analyze.cql.Close()
	}

	g.L.Info("Close ok!")
	return nil
}
