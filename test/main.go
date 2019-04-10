package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
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
		deferTime := rand.Intn(30)
		// 启动定时器
		go ticker.start(deferTime)
		g.L.Info("defer start", zap.Int("deferTime", deferTime))
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
		log.Println(err.Error())
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

func helloTask(id int64, channel chan bool) {
	for {
		select {
		case <-channel:
			log.Println("任务开始执行, id为", id)
			break
			// default:
			// 	log.Println("hello")
			// 	break
		}
	}
}

func tickertest() {
	tickers := newTickers(10)

	for index := 0; index < 1000; index++ {
		id := tickers.getID()
		channel := make(chan bool, 1)
		tickers.addTask(id, channel)
		go helloTask(id, channel)
		time.Sleep(1 * time.Second)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func hostname() {
	// n, e := os.Hostname()
	// if e != nil {
	// 	return
	// }
	// log.Println(n)
}

type stats struct {
	index  Index
	points map[int64]int
}

func newstats() *stats {
	return &stats{
		points: make(map[int64]int),
	}
}

// Index ....
type Index []int64

func (index Index) Swap(i, j int)      { index[i], index[j] = index[j], index[i] }
func (index Index) Len() int           { return len(index) }
func (index Index) Less(i, j int) bool { return index[i] < index[j] }

func stattest() {
	stats := newstats()
	stats.points[1000201200212] = 1232111
	stats.points[10200201200212] = 9878978
	stats.points[100020120021] = 987897
	stats.points[10002012002126] = 32345
	stats.points[1000201202] = 565678
	stats.points[10002012] = 98786
	stats.points[10002012022121223] = 7657

	start := time.Now()
	for key := range stats.points {
		stats.index = append(stats.index, key)
	}

	sort.Sort(stats.index)

	for _, key := range stats.index {
		log.Println(stats.points[key])
	}

	log.Println("耗时", time.Now().Sub(start).Nanoseconds()/1000000)

}

func main() {
	stattest()
}
