package misc

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Common struct {
		Version    string
		LogLevel   string
		AdminToken string
	}
	MQ struct {
		Addrs []string // mq地址
		Topic string   // 主题
	}

	App struct {
		LoadInterval int
	}

	DB struct {
		Cluster  []string
		Keyspace string
		NumConns int
	}

	Analyze struct {
		Interval int
	}
}

var Conf *Config

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
}
