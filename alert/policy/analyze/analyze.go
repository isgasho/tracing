package analyze

import (
	"sync"

	"github.com/imdevlab/tracing/pkg/alert"
	"github.com/imdevlab/tracing/pkg/constant"

	"go.uber.org/zap"

	"github.com/imdevlab/tracing/alert/policy"
)

// // 获取任务ID
// a.taskID = gCollector.ticker.NewID()
// logger.Info("app start", zap.String("name", a.name), zap.Int64("taskID", a.taskID))
// // 加入定时模块
// gCollector.ticker.AddTask(a.taskID, a.tickerC)
var logger *zap.Logger

// Analyze 数据分析&实时计算
type Analyze struct {
	sync.RWMutex
	Policys *policy.Policys
}

// New new analyze
func New(l *zap.Logger) *Analyze {
	logger = l
	return &Analyze{
		Policys: policy.NewPolicys(logger),
	}
}

// Start .
func (a *Analyze) Start() error {
	// // 加载策略
	// if err := a.loadPolicys(); err != nil {
	// 	logger.Warn("load policys error", zap.String("error", err.Error()))
	// 	return err
	// }
	return nil
}

// Write recv data
func (a *Analyze) Write(data *alert.Data) error {
	switch data.Type {
	case constant.ALERT_TYPE_API:
		break
	case constant.ALERT_TYPE_SQL:
		break
	}
	return nil
}
