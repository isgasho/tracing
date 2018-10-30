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

	Vgo struct {
		ListenAddr   string
		AgentTimeout int
	}

	Storage struct {
		Cluster           []string
		Keyspace          string
		NumConns          int
		JVMCacheLen       int // 缓存长度
		JVMStoreInterval  int // 毫秒
		JVMStoreLen       int // 时间未到超过该长度即可插入
		SpanCacheLen      int // 缓存长度
		SpanStoreLen      int // 时间未到超过该长度即可插入
		SpanStoreInterval int // 毫秒
	}

	Mysql struct {
		Addr     string
		Port     string
		Database string
		Acc      string
		Pw       string
	}
}

// Conf ...
var Conf *Config

// InitConfig ...
func InitConfig(path string) {
	conf := &Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("read config error :", err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("yaml decode error :", err)
	}
	Conf = conf
	log.Println(Conf)
}
