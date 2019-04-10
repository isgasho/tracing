package service

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/imdevlab/g"
)

// Tickers 定时器集合
type Tickers struct {
	sync.Mutex
	id      int64              // 自增任务编号
	hash    *g.Hash            // 一致性hash工具
	tickers map[string]*Ticker // 定时器集合
}

func newTickers(tickerNum int) *Tickers {
	tickers := &Tickers{
		hash:    g.NewHash(),
		tickers: make(map[string]*Ticker),
	}

	rand.Seed(time.Now().UnixNano())

	for index := 0; index < tickerNum; index++ {
		tickers.hash.Add(fmt.Sprintf("%d", index))
		ticker := newTicker()
		// 延迟启动时间
		deferTime := rand.Intn(30)
		g.L.Info("deferTime", zap.Int("deferTime", deferTime))
		// 启动定时器
		go ticker.start(deferTime)

		tickers.tickers[fmt.Sprintf("%d", index)] = ticker
	}

	return tickers
}

func (t *Tickers) getID() int64 {
	t.Lock()
	// id自增
	id := t.id
	t.id++
	t.Unlock()
	return id
}

func (t *Tickers) addTask(id int64, channel chan bool) {
	// id 通过hash计算出来key
	key, err := t.hash.Get(fmt.Sprintf("%d", id))
	if err != nil {
		g.L.Warn("hash get", zap.String("error", err.Error()))
		return
	}
	// 加入任务
	t.tickers[key].addTask(id, channel)
}

// Ticker 定时器
type Ticker struct {
	sync.RWMutex
	tasks map[int64]*Task
}

func newTicker() *Ticker {
	return &Ticker{
		tasks: make(map[int64]*Task),
	}
}

func (t *Ticker) start(deferTime int) {
	// 延迟启动
	time.Sleep(time.Second * time.Duration(deferTime))
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			t.RLock()
			for _, task := range t.tasks {
				log.Println("任务触发， task.id is", task.id)
				task.channel <- true
			}
			t.RUnlock()
			break
		}
	}
}

func (t *Ticker) addTask(id int64, channel chan bool) {
	t.Lock()
	t.tasks[id] = newTask(id, channel)
	t.Unlock()
}

// Task 任务
type Task struct {
	id        int64     // 编号
	channel   chan bool // 通知管道
	inputTime time.Time // 插入时间
}

// newTask 新任务
func newTask(id int64, channel chan bool) *Task {
	return &Task{
		id:        id,
		channel:   channel,
		inputTime: time.Now(),
	}
}
