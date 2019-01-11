package misc

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	Common struct {
		Version    string
		LogLevel   string
		AdminToken string
	}

	Analyze struct {
		LoadAppInterval int
	}

	Cassandra struct {
		Cluster  []string
		Keyspace string
		NumConns int
	}

	Cluster struct {
		Addr        string
		Port        int
		Name        string
		Seeds       []string
		HostUseTime bool
	}

	Etcd struct {
		Addrs      []string
		Dltimeout  int
		ReportDir  string
		ReportTime int
		WatchDir   string
		TTL        int64
	}

	Stats struct {
		Interval         int
		Range            int64
		SatisfactionTime int
		TolerateTime     int
	}
}

// Conf ...
var Conf *Config

// InitConfig ...
func InitConfig(path string) {
	conf := &Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("InitConfig:ioutil.ReadFile", err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("InitConfig:yaml.Unmarshal", err)
	}
	Conf = conf
	log.Println(Conf)
}
