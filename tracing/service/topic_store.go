package service

// import (
// 	"log"
// 	"sync"

// 	"github.com/mafanr/g"
// )

// // TopicStore topic 缓存
// type TopicStore struct {
// 	sync.RWMutex
// 	topics   map[string]string // key is host_pid, value is topic
// 	appCache map[string]string // key is appname, value is topic
// 	hash     *g.Hash
// }

// // NewTopicStore new topic store
// func NewTopicStore() *TopicStore {
// 	etcdTopics := &TopicStore{
// 		topics:   make(map[string]string),
// 		appCache: make(map[string]string),
// 		hash:     g.NewHash(),
// 	}
// 	return etcdTopics
// }

// // Add add topic
// func (ts *TopicStore) Add(key string, topic string) {
// 	ts.RLock()
// 	// 将地址加入hash列表
// 	_, ok := ts.topics[key]
// 	ts.RUnlock()
// 	if ok {
// 		return
// 	}
// 	// add topic into hash
// 	ts.hash.Add(topic)

// 	// 只添加topic不添加appCache
// 	ts.Lock()
// 	ts.topics[key] = topic
// 	//新加
// 	ts.appCache = make(map[string]string)
// 	ts.Unlock()

// 	log.Println("Add", ts.topics)
// 	log.Println("Add", ts.appCache)
// }

// // Remove delete topic from cache
// func (ts *TopicStore) Remove(key string) {
// 	ts.RLock()
// 	topic, ok := ts.topics[key]
// 	ts.RUnlock()
// 	if !ok {
// 		return
// 	}

// 	log.Println("remove:", key, topic)

// 	// 将地址从hash列表中删除
// 	ts.hash.Remove(topic)

// 	ts.Lock()
// 	delete(ts.topics, key)
// 	ts.appCache = make(map[string]string)
// 	ts.Unlock()
// }

// // Get get topic by appName
// func (ts *TopicStore) Get(appName string) (string, bool) {
// 	ts.RLock()
// 	topic, ok := ts.appCache[appName]
// 	ts.RUnlock()
// 	if ok {
// 		return topic, ok
// 	}

// 	newTopic, err := ts.hash.Get(appName)
// 	if err != nil {
// 		return "", false
// 	}
// 	ts.Lock()
// 	ts.appCache[appName] = newTopic
// 	ts.Unlock()

// 	return newTopic, true
// }
