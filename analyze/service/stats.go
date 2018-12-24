package service

import (
	"log"
	"sync"
	"time"

	"github.com/mafanr/g"
)

// Stats 离线计算
type Stats struct {
}

// NewStats ...
func NewStats() *Stats {
	return &Stats{}
}

// Start ...
func (s *Stats) Start() error {
	g.L.Info("Start Stats")

	return nil
}

// Close ...
func (s *Stats) Close() error {
	g.L.Info("Close Stats")

	return nil
}

func (s *Stats) counter(app *App, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, agent := range app.Agents {
		log.Println("----->>>", agent.AgentID, agent.lastPointTime, agent.startTime)
		var queryStartTime int64
		var queryEndTime int64
		if agent.lastPointTime == 0 {
			queryStartTime = agent.startTime
		} else {
			queryStartTime = agent.lastPointTime
		}
		queryEndTime = queryStartTime + 60*1000

		log.Println("queryStartTime", time.Unix(0, queryStartTime*1e6).String(), queryStartTime)
		log.Println("queryEndTime", time.Unix(0, queryEndTime*1e6).String(), queryEndTime)
		// query:= fmt.Sprintf("", a ...interface{})
	}
}
