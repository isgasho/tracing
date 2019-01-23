package service

import (
	"fmt"
	"os"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"

	"go.uber.org/zap"
)

// Analyze ...
type Analyze struct {
	stats        Stats        // 离线统计
	blink        Blink        // 实时计算
	serDiscovery SerDiscovery // 服务发现
	cql          *g.Cassandra
	appStore     *AppStore
	hash         *g.Hash
	analyzes     map[string]string
	clusterName  string
}

var gAnalyze *Analyze

// New ...
func New() *Analyze {
	analyze := &Analyze{
		stats:        newStats(),
		blink:        newBlink(),
		cql:          g.NewCassandra(),
		hash:         g.NewHash(),
		serDiscovery: newEtcd(),
		analyzes:     make(map[string]string),
	}
	gAnalyze = analyze
	return analyze
}

// Start ...
func (analyze *Analyze) Start() error {

	g.L.Info("Conf", zap.Any("conf", misc.Conf))
	watchDir := initDir(misc.Conf.Etcd.WatchDir)

	reportValue, _ := analyzeName()
	gAnalyze.clusterName = reportValue

	reportDir := initDir(misc.Conf.Etcd.ReportDir)

	if err := analyze.serDiscovery.Init(reportDir+reportValue, reportValue, watchDir); err != nil {
		g.L.Fatal("Start etcd.Start", zap.String("error", err.Error()))
	}

	if err := analyze.serDiscovery.Start(); err != nil {
		g.L.Fatal("Start etcd.Start", zap.String("error", err.Error()))
	}

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

	if analyze.serDiscovery != nil {
		analyze.serDiscovery.Close()
	}

	g.L.Info("Close ok!")
	return nil
}

func initDir(dir string) string {
	dirLen := len(dir)
	if dirLen > 0 && dir[dirLen-1] != '/' {
		return dir + "/"
	}
	return dir
}

// analyzeName get key
func analyzeName() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d", host, os.Getpid()), nil
}

func reportKey(dir string) (string, error) {
	value, err := analyzeName()
	if err != nil {
		return "", err
	}

	dirLen := len(dir)
	if dirLen > 0 && dir[dirLen-1] != '/' {
		return dir + "/" + value, nil
	}
	return dir + value, nil
}
