package policy

import (
	"sync"
	"time"

	"github.com/imdevlab/tracing/pkg/alert"

	"github.com/gocql/gocql"
	"github.com/imdevlab/tracing/alert/misc"
	"github.com/imdevlab/tracing/alert/policy/ticker"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/util"

	"github.com/imdevlab/tracing/pkg/sql"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Policys 策略s
type Policys struct {
	sync.RWMutex
	tickers *ticker.Tickers
	cql     *gocql.Session
	Policys map[string]*Policy
}

// NewPolicys new policys
func NewPolicys(l *zap.Logger) *Policys {
	logger = l
	return &Policys{
		Policys: make(map[string]*Policy),
		tickers: ticker.NewTickers(10, misc.Conf.Analyze.Interval, logger),
	}
}

// LoadPolicys 加载策略
func (p *Policys) LoadPolicys(cql *gocql.Session) error {
	query := cql.Query(sql.LoadPolicys).Iter()
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
			oldPolicy.close()
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

		alertsQuery := cql.Query(sql.LoadAlert, policyID)
		var tmpAlerts []*util.Alert
		if err := alertsQuery.Scan(&tmpAlerts); err != nil {
			logger.Warn("load alert scan error", zap.String("error", err.Error()), zap.String("sql", sql.LoadAlert))
			continue
		}
		if len(tmpAlerts) == 0 {
			continue
		}

		for _, tmpAlert := range tmpAlerts {
			alert := newAlertInfo()
			alert.Compare = tmpAlert.Compare
			alert.Duration = tmpAlert.Duration
			alert.Value = tmpAlert.Value
			alertType, ok := constant.AlertType(tmpAlert.Name)
			if !ok {
				logger.Warn("alertType unfind error", zap.String("name", tmpAlert.Name))
				continue
			}
			alert.Type = alertType
			newPolicy.Alerts[alertType] = alert
		}

		// 保存策略
		p.Lock()
		p.Policys[name] = newPolicy
		newPolicy.start()
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

// Policy 策略
type Policy struct {
	AppName    string             // app名
	Owner      string             // owner
	ID         string             // policyid
	Group      string             // 组
	Channel    string             // 通道？？？？和组互斥？
	Users      []string           // 用户
	UpdateDate int64              // 更新时间
	Alerts     map[int]*AlertInfo // 策略模版
	checkTime  int64              // 上次检查时间
	apiC       chan *alert.API    //
}

func (p *Policy) analyze() {
	// for {

	// }
}

func (p *Policy) start() {
	logger.Info("policy start", zap.String("appName", p.AppName))
	// 启动计算线程
	go p.analyze()
}

func (p *Policy) close() {
	logger.Info("policy close", zap.String("appName", p.AppName))
}

func newPolicy() *Policy {
	return &Policy{
		Alerts: make(map[int]*AlertInfo),
	}
}

func newAlertInfo() *AlertInfo {
	return &AlertInfo{}
}

// AlertInfo 策略信息
type AlertInfo struct {
	Type     int     // 监控项类型
	Compare  int     // 比较类型 1: > 2:<  3:=
	Duration int     // 持续时间, 1 代表1分钟
	Value    float64 // 阀值
}
