package service

import (
	"log"
	"sync"
	"time"

	"github.com/imdevlab/tracing/pkg/util"

	"github.com/imdevlab/tracing/pkg/sql"
	"go.uber.org/zap"
)

// Policys 策略s
type Policys struct {
	sync.RWMutex
	Policys map[string]*Policy
}

// Start 启动策略服务
func (p *Policys) Start() error {
	p.loadPolicys()
	return nil
}

// Start 启动策略服务
func (p *Policys) loadPolicys() error {
	query := gAlert.cql.Query(sql.LoadPolicys).Iter()
	defer func() {
		if err := query.Close(); err != nil {
			logger.Warn("close iter error:", zap.Error(err))
		}
	}()

	var name, owner, channel, group, policyID string
	var users []string
	var updateDate int64
	checkTime := time.Now().Unix()
	for query.Scan(&name, &owner, &channel, &group, &policyID, &updateDate, &users) {

		p.RLock()
		oldPolicy, ok := p.Policys[name]
		p.RUnlock()
		// 如果已经存在策略并且updatedate不相当，那么删除历史
		if ok {
			if oldPolicy.UpdateDate == updateDate {
				oldPolicy.checkTime = checkTime
				continue
			}
			p.Lock()
			delete(p.Policys, name)
			p.Unlock()
		}

		newPolicy := newPolicy()
		newPolicy.AppName = name
		newPolicy.Owner = owner
		newPolicy.Channel = channel
		newPolicy.Group = group
		newPolicy.ID = policyID
		newPolicy.UpdateDate = updateDate
		newPolicy.Users = users
		newPolicy.checkTime = checkTime
		// 根据alertid加具体载策略,如果policyID为null那么代表该模版的策略被删除，所以不用统计
		if len(policyID) == 0 {
			continue
		}

		log.Println(policyID)

		alertsQuery := gAlert.cql.Query(sql.LoadAlert, policyID)
		var tmpAlerts []*util.Alert
		if err := alertsQuery.Scan(&tmpAlerts); err != nil {
			logger.Warn("load alert scan error", zap.String("error", err.Error()), zap.String("sql", sql.LoadAlert))
			continue
		}
		if len(tmpAlerts) == 0 {
			continue
		}

		log.Println(len(tmpAlerts))
		for index, alert := range tmpAlerts {
			log.Println(index, alert)
		}

		// 保存策略
		p.Lock()
		p.Policys[name] = newPolicy
		p.Unlock()
	}

	// 对比模版checktime，发现checktime不相等，那么代表该模版已经被删除
	p.Lock()
	for name, poklicy := range p.Policys {
		if poklicy.checkTime != checkTime {
			delete(p.Policys, name)
		}
	}
	p.Unlock()

	return nil
}

func newPolicys() *Policys {
	return &Policys{
		Policys: make(map[string]*Policy),
	}
}

// Policy 策略
type Policy struct {
	AppName      string        // app名
	Owner        string        // owner
	ID           string        // policyid
	Group        string        // 组
	Channel      string        // 通道？？？？和组互斥？
	Users        []string      // 用户
	UpdateDate   int64         // 更新时间
	AlertsPolicy *AlertsPolicy // 策略模版
	checkTime    int64         // 上次检查时间
}

func newPolicy() *Policy {
	return &Policy{}
}

// AlertsPolicy 告警策略模版表
type AlertsPolicy struct {
	ID     string
	Name   string
	Owner  string
	Alerts []*AlertInfo
}

// AlertInfo 策略信息
type AlertInfo struct {
	Type     int     // 监控项类型
	Compare  int     // 比较类型 1: > 2:<  3:=
	Unit     int     // 单位：%、个
	Duration int     // 持续时间, 1 代表1分钟
	Value    float64 // 阀值
}
