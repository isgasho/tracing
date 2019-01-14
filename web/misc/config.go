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

	Etcd struct {
		Addrs     []string
		Dltimeout int
		WatchKey  string
	}

	Login struct {
		SsoLogin  string
		SsoLogout string
	}

	Storage struct {
		Cluster  []string
		Keyspace string
		NumConns int
	}

	Web struct {
		Addr string
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
