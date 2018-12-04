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

	Agent struct {
		VgoAddr          string
		KeepLiveInterval int
		UseEnv           bool
		ENV              string
		AppName          string
	}

	Pinpoint struct {
		InfoAddr          string // tcp addr for info
		StatAddr          string // udp addr for stat
		SpanAddr          string // udp addr for span
		SpanRportInterval int    // 全链路信息上报频率 单位毫秒
	}

	SkyWalking struct {
		HTTPAddr             string
		RPCAddr              string
		JVMReportInterval    int // jvm 信息上报频率
		JVMCollectorInterval int // jvm 采集频率控制
		JVMCacheLen          int // 缓存长度
		TraceReportInterval  int // 全链路信息上报频率 单位毫秒
		TraceCacheLen        int // 缓存长度
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
