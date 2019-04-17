package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/web/internal/misc"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Trace struct {
	ID        string `json:"id"`
	API       string `json:"api"`
	Elapsed   int    `json:"y"`
	AgentID   string `json:"agent_id"`
	InputDate int64  `json:"x"`
	Error     int    `json:"error"`
}

type Traces []*Trace

func (a Traces) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a Traces) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a Traces) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Elapsed < a[i].Elapsed
}

type ChartTraces struct {
	Suc   bool    `json:"is_suc"`
	Xaxis []int64 `json:"timeXticks"`
	Title string  `json:"subTitle"`

	Series []*TraceSeries `json:"series"`
}

type TraceSeries struct {
	Name  string   `json:"name"`
	Color string   `json:"color"`
	Data  []*Trace `json:"data"`
}

// traceSeries":[{"name":"success","color":"rgb(18, 147, 154,.5)","data":[{"x":1545200556716,"y":7,"traceId":"yunbaoParkApp3^1545036617750^4217","agentId":"agencyBookKeep3","startTime":"1545200556716","url":"/agencyBookKeep/financialstatementscs/getOneAccountingDataByParams","traceIp":"127.0.0.1"},
func QueryTraces(c echo.Context) error {
	appName := c.FormValue("app_name")
	api := c.FormValue("api")
	min, _ := strconv.Atoi(c.FormValue("min_elapsed"))
	max, _ := strconv.Atoi(c.FormValue("max_elapsed"))
	limit, err := strconv.Atoi(c.FormValue("limit"))
	if err != nil {
		limit = 50
	}

	start, end, _ := misc.StartEndDate(c)

	var q *gocql.Query
	if api == "" {
		if max == 0 {
			q = misc.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date,error FROM app_operation_index WHERE app_name=? and input_date > ? and input_date < ? and elapsed >= ? ALLOW FILTERING`, appName, start.Unix()*1000, end.Unix()*1000, min)
		} else {
			q = misc.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date,error FROM app_operation_index WHERE app_name=? and input_date > ? and input_date < ? and elapsed >= ? and elapsed <= ? ALLOW FILTERING`, appName, start.Unix()*1000, end.Unix()*1000, min, max)
		}
	} else {
		if max == 0 {
			q = misc.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date,error FROM app_operation_index WHERE app_name=? and api=?  and input_date > ? and input_date < ? and elapsed >= ? ALLOW FILTERING`, appName, api, start.Unix()*1000, end.Unix()*1000, min)
		} else {
			q = misc.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date,error FROM app_operation_index WHERE app_name=? and api=?  and input_date > ? and input_date < ? and elapsed >= ? and elapsed <= ? ALLOW FILTERING`, appName, api, start.Unix()*1000, end.Unix()*1000, min, max)
		}
	}

	iter := q.Iter()

	var elapsed, isError int
	var inputDate int64
	var tid, agentID string

	traceMap := make(map[string]*Trace)
	for iter.Scan(&tid, &api, &elapsed, &agentID, &inputDate, &isError) {
		traceMap[tid] = &Trace{tid, api, elapsed, agentID, inputDate, isError}
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	traces := make(Traces, 0, len(traceMap))
	for _, t := range traceMap {
		traces = append(traces, t)
	}

	sort.Sort(traces)

	// 取出耗时最高的limit数量的trace
	if limit < len(traces) {
		traces = traces[:limit]
	}

	ct := &ChartTraces{}
	if len(traces) == 0 {
		ct.Suc = false
	} else {
		ct.Suc = true
		ct.Xaxis = []int64{start.Unix() / 1e6, end.Unix() / 1e6}
		ct.Title = fmt.Sprintf("success: %d, error: %d", len(traces), 0)

		var sucTraces Traces
		var errTraces Traces

		for _, t := range traces {
			if t.Error == 0 {
				sucTraces = append(sucTraces, t)
			} else {
				errTraces = append(errTraces, t)
			}
		}
		sucData := &TraceSeries{
			Name:  "success",
			Color: "rgb(18, 147, 154,.5)",
			Data:  sucTraces,
		}

		errData := &TraceSeries{
			Name:  "error",
			Color: "rgba(223, 83, 83, .5)",
			Data:  errTraces,
		}

		ct.Series = []*TraceSeries{sucData, errData}
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ct,
	})
}

func QueryTrace(c echo.Context) error {
	tid := c.FormValue("trace_id")

	// 加载trace的所有span
	spans := make(traceSpans, 0)
	err := spans.load(tid)
	if err != nil {
		return err
	}

	// 把span按照请求链路的顺序排序，例如A、B、C三个服务，请求顺序A -> B -> C，那么span的排列顺序也应该是span(A) -> span(B) -> span(C)
	spans.sort()

	// 将span和events组合成链路tree
	tree := make(TraceTree, 0)
	for _, span := range spans {
		tree.addSpan(span)

		for _, event := range span.events {
			tree.addEvent(event, span)
		}
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   tree,
	})
}

/*-----------------------------数据结构和方法定义------------------------------------*/
//
// trace span
//
type traceSpan struct {
	id          int64 // span id
	pid         int64 // the parent span id
	appName     string
	agentID     string
	serviceType int
	startTime   int64 // ms timestamp
	events      []*SpanEvent
	duration    int
	api         string // 接口url
	methodID    int
	remoteAddr  string
	annotations []*TempTag
}

type traceSpans []*traceSpan

func (o traceSpans) Len() int {
	return len(o)
}

func (o traceSpans) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o traceSpans) Less(i, j int) bool {
	return o[i].startTime < o[j].startTime
}

func (spans *traceSpans) load(tid string) error {
	q := misc.Cql.Query(`SELECT span_id,parent_id,app_name,agent_id,input_date,elapsed,api,service_type,
	end_point,remote_addr,err,span_event_list,method_id,annotations,exception_info from traces where trace_id=?`, tid)
	iter := q.Iter()

	var spanID, pid, inputDate int64
	var elapsed, serviceType, isErr, methodID int
	var appName, agentID, api, endPoint, remoteAddr, events, annotations, exception string

	// parse span
	for iter.Scan(&spanID, &pid, &appName, &agentID, &inputDate, &elapsed, &api, &serviceType,
		&endPoint, &remoteAddr, &isErr, &events, &methodID, &annotations, &exception) {
		// 首先把span本身转为segment
		var tags []*TempTag
		json.Unmarshal([]byte(annotations), &tags)

		span := &traceSpan{
			appName:     appName,
			agentID:     agentID,
			duration:    elapsed,
			api:         api,
			serviceType: serviceType,
			startTime:   inputDate,
			// startTime:   misc.Timestamp2TimeString(inputDate),
			id:          spanID,
			pid:         pid,
			remoteAddr:  remoteAddr,
			annotations: tags,
			methodID:    methodID,
		}
		fmt.Println("span error:", isErr, exception)

		// 解析span的events，并根据sequence进行排序(从小到大)
		var spanEvents SpanEvents
		json.Unmarshal([]byte(events), &spanEvents)
		// 加载span chunk
		q1 := misc.Cql.Query(`SELECT span_event_list from traces_chunk where trace_id=? and span_id=?`, tid, span.id)
		iter1 := q1.Iter()
		var eventsChunkS string
		for iter1.Scan(&eventsChunkS) {
			var eventsChunk SpanEvents
			json.Unmarshal([]byte(eventsChunkS), &eventsChunk)

			spanEvents = append(spanEvents, eventsChunk...)
		}
		sort.Sort(spanEvents)
		span.events = spanEvents

		*spans = append(*spans, span)
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
		return err
	}

	return nil
}

func (spans *traceSpans) sort() {
	bucket := make(map[int64]int)
	for i, span := range *spans {
		bucket[span.id] = i
	}

	// 找到起始span
	// span之间可能存在以下几种关系
	// 通过父子串联
	// 若不存在父子关系(离散服务的span，特殊情况)，则通过时间排序(不同服务器时间不同，因此不够可靠)？
	noParents := make([]int, 0)
	for i, span := range *spans {
		// 若没有父节点，则加入noParents
		_, ok := bucket[span.pid]
		if !ok {
			noParents = append(noParents, i)
		}
	}

	var nspans traceSpans
	if len(noParents) == 1 {
		//只有一个节点没有父节点，证明可以串起来
		nspans = append(nspans, (*spans)[noParents[0]])
		for {
			if len(nspans) >= len(*spans) {
				break
			}

			currentSpan := nspans[len(nspans)-1]
			for _, span := range *spans {
				if span.pid == currentSpan.id {
					nspans = append(nspans, span)
					break
				}
			}
		}
		*spans = nspans
	} else {
		// 多个父节点，通过时间排序
		sort.Sort(spans)
	}
}

//
// 为了形成全链路，我们需要把trace的span和event组成一个tree结构,span和event对应的是tree node
// Tags、Exceptions都将转换为node进行展示
type TraceTreeNode struct {
	ID          string      `json:"id"` // seg的id,p-0/p-0-1类的形式，通过这种形式形成层级tree
	Sequence    int         `json:"seq"`
	Depth       int         `json:"depth"`      //node在完整的链路树上所处的层级，绝对层级 Tree Depth
	SpanDepth   int         `json:"span_depth"` //node对应的event在对应span中的层级,相对层级 span depth
	AppName     string      `json:"app_name"`
	MethodID    int         `json:"method_id"`
	Method      string      `json:"method"`
	Duration    int         `json:"duration"` // 耗时，-1 代表不显示耗时信息
	Params      string      `json:"params"`
	ServiceType string      `json:"service_type"`
	AgentID     string      `json:"agent_id"`
	Class       string      `json:"class"`
	StartTime   string      `json:"start_time"`
	Tags        []*TraceTag `json:"tags"`     // Seg标签
	Icon        string      `json:"icon"`     // 有些节点会显示特殊的icon，作为样式
	IsError     bool        `json:"is_error"` // 是否是错误/异常，
}

// 初始化node
func (node *TraceTreeNode) init(appName, agentID, params string, service, methodID int) {
	node.AppName = appName
	node.Params = params
	node.AgentID = agentID

	node.ServiceType = constant.ServiceType[service]

	// 通过method id 查询method
	method := misc.GetMethodByID(appName, methodID)
	node.Class, node.Method = misc.SplitMethod(method)
	node.MethodID = methodID
}

func (node *TraceTreeNode) setDepth(spanDepth int, treeDepth int) {
	node.SpanDepth = spanDepth
	node.Depth = treeDepth
}

// 若nid 为 p-0-1，则当前node的ID为p-0-2
func (node *TraceTreeNode) setNeighborID(s string) {
	sep := strings.LastIndex(s, "-")
	s2 := s[sep+1:]
	i, _ := strconv.Atoi(s2)
	i = i + 1
	node.ID = s[:sep+1] + strconv.Itoa(i)
}

// 若传入父id为p-0-1，则设置id为p-0-1-0
func (node *TraceTreeNode) setChildID(s string) {
	node.ID = s + "-0"
}

// 解析annotations，转换为span tag
func (node *TraceTreeNode) setTags(tags []*TempTag) {
	for _, tag := range tags {
		if (tag.Key == constant.STRING_ID) || (tag.Key <= constant.CACHE_ARGS0 && tag.Key >= constant.CACHE_ARGSN) {
			// 添加method_id : method的tag
			methodID := int(tag.Value.IntValue)
			method := misc.GetMethodByID(node.AppName, methodID)
			stag := &TraceTag{constant.AnnotationKeys[tag.Key], method}
			node.Tags = append(node.Tags, stag)
		} else if tag.Key == constant.SQL_ID {
			// {"key":20,"value":{"intStringStringValue":{"intValue":1,"stringValue1":"0","stringValue2":"testC, testC, 2019-04-15 08:43:03.713, null, null, testCNickName, testC_64b3def7-1a76-4ed7-bf21-67f5afc440fc, E10ADC3949BA59ABBE56E057F20F883E, null, 0"}}}
			sqlID := int(tag.Value.IntStringStringValue.IntValue)
			sqlS := misc.GetSqlByID(node.AppName, sqlID)
			sql, _ := g.B64.DecodeString(sqlS)
			// 添加sqlID: sql的tag
			stag1 := &TraceTag{constant.AnnotationKeys[tag.Key], string(sql)}
			// 添加sql bind value的tag
			stag2 := &TraceTag{constant.AnnotationKeys[constant.SQL_BINDVALUE], tag.Value.IntStringStringValue.StringValue2}
			node.Tags = append(node.Tags, stag1, stag2)
		} else {
			var stag *TraceTag
			switch tag.Key {
			case constant.HTTP_STATUS_CODE:
				stag = &TraceTag{constant.AnnotationKeys[tag.Key], strconv.Itoa(int(tag.Value.IntValue))}
			default:
				fmt.Println("tag key", tag.Key, constant.AnnotationKeys[tag.Key])
				stag = &TraceTag{constant.AnnotationKeys[tag.Key], tag.Value.StringValue}
			}

			node.Tags = append(node.Tags, stag)
		}
	}
}

type TraceTree []*TraceTreeNode

func (tree *TraceTree) addSpan(span *traceSpan) {
	n := &TraceTreeNode{}
	n.init(span.appName, span.agentID, span.api, span.serviceType, span.methodID)
	n.setTags(span.annotations)
	n.Duration = span.duration
	// span本身一定是http/dubbo/rpc服务的入口，因此要做特殊标示
	n.Icon = "hand"
	//remote addr -> tag
	n.Tags = append(n.Tags, &TraceTag{"remote_addr", span.remoteAddr})

	//@test
	n.Sequence = 0

	//set node id
	if len(*tree) == 0 {
		n.setDepth(0, 0)
		// 第一个span也是第一个node
		n.ID = "p-0"
	} else {
		lastn := (*tree)[len(*tree)-1]
		// 若上一个节点的depth是-1，那该span是它的兄弟节点
		// 否则，该span是它的子节点
		if lastn.SpanDepth == -1 {
			n.setDepth(0, lastn.Depth)
			n.setNeighborID(lastn.ID)
		} else {
			n.setDepth(0, lastn.Depth+1)
			n.setChildID(lastn.ID)
		}
	}

	*tree = append(*tree, n)

	// tags -> node
	for _, tag := range n.Tags {
		en := &TraceTreeNode{}
		en.setDepth(-1, n.Depth+1)
		en.setChildID(n.ID)

		en.Params = tag.Value
		// 获取exception id
		en.Method = tag.Key

		en.Duration = -1
		en.Icon = "info"
		*tree = append(*tree, en)
	}
}

func (tree *TraceTree) addEvent(event *SpanEvent, span *traceSpan) {
	n := &TraceTreeNode{}
	n.init(span.appName, span.agentID, event.DestinationID, event.ServiceType, event.MethodID)
	n.setTags(event.Annotations)
	n.Duration = event.EndElapsed

	//@test
	n.Sequence = event.Sequence

	// 若当前event的span depth为-1，则该event为叶子node
	//     我们要找到上一个不是叶子的节点，然后把该event作为该节点的最新的叶子node
	// 若当前event的span depth不为-1
	//     我们要找到depth-1的节点，然后该节点是当前event的父节点
	if event.Depth == -1 {
		lastn := (*tree)[len(*tree)-1]
		if lastn.SpanDepth == -1 {
			// 上一个节点也是叶子节点
			// 因此该event是上一个节点的邻节点
			n.setDepth(event.Depth, lastn.Depth)
			// lastn的邻节点，因此id + 1: 例如p-0-1 -> p-0-2
			n.setNeighborID(lastn.ID)
		} else {
			// 当前的event是上一个节点的子节点
			n.setDepth(event.Depth, lastn.Depth+1)
			n.setChildID(lastn.ID)
		}
	} else {
		// 寻找该event的兄弟节点或者父节点
		for i := len(*tree) - 1; i >= 0; i-- {
			// 先寻找兄弟节点：depth相同
			if (*tree)[i].SpanDepth == event.Depth {
				n.setDepth(event.Depth, (*tree)[i].Depth)
				n.setNeighborID((*tree)[i].ID)
				break
			}
			// 再寻找父节点
			if (*tree)[i].SpanDepth == event.Depth-1 {
				n.setDepth(event.Depth, (*tree)[i].Depth+1)
				n.setChildID((*tree)[i].ID)
				break
			}
		}
	}

	*tree = append(*tree, n)

	// 将exception转为当前event node的叶子node
	// 叶子节点的span depth = -1
	if event.ExceptionInfo != nil {
		en := &TraceTreeNode{}
		en.setDepth(-1, n.Depth+1)
		en.setChildID(n.ID)

		en.Params = event.ExceptionInfo.StringValue
		// 获取exception id
		en.Method = misc.GetClassByID(n.AppName, int(event.ExceptionInfo.IntValue))
		fmt.Println(event.ExceptionInfo, n.AppName, en.Method)
		en.IsError = true
		en.Duration = -1
		en.Icon = "bug"

		fmt.Println("enid:", event.Sequence, en.ID)
		*tree = append(*tree, en)
	}

	// tags -> node
	for _, tag := range n.Tags {
		en := &TraceTreeNode{}
		en.setDepth(-1, n.Depth+1)
		en.setChildID(n.ID)

		en.Params = tag.Value
		// 获取exception id
		en.Method = tag.Key

		en.Duration = -1
		en.Icon = "info"
		*tree = append(*tree, en)
	}
}

// 标签，原数据为annotations，统一转换为tag
type TraceTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// 原数据annotations使用的格式
type TempTag struct {
	Key   int       `json:"key"`
	Value *TagValue `json:"value"`
}

type TagValue struct {
	StringValue                   string                         `json:"stringValue,omitempty"`
	BoolValue                     bool                           `json:"boolValue,omitempty"`
	IntValue                      int32                          `json:"intValue,omitempty"`
	LongValue                     int64                          `json:"longValue,omitempty"`
	ShortValue                    int16                          `json:"shortValue,omitempty"`
	DoubleValue                   float64                        `json:"doubleValue,omitempty"`
	BinaryValue                   []byte                         `json:"binaryValue,omitempty"`
	ByteValue                     int8                           `json:"byteValue,omitempty"`
	IntStringValue                *IntStringValue                `json:"intStringValue,omitempty"`
	IntStringStringValue          *IntStringStringValue          `json:"intStringStringValue,omitempty"`
	LongIntIntByteByteStringValue *LongIntIntByteByteStringValue `json:"longIntIntByteByteStringValue,omitempty"`
	IntBooleanIntBooleanValue     *IntBooleanIntBooleanValue     `json:"intBooleanIntBooleanValue,omitempty"`
}

type IntStringValue struct {
	IntValue    int32  `json:"intValue"`
	StringValue string `json:"stringValue,omitempty"`
}

type IntStringStringValue struct {
	IntValue     int32  `json:"intValue"`
	StringValue1 string `json:"stringValue1,omitempty"`
	StringValue2 string `json:"stringValue2,omitempty"`
}

type LongIntIntByteByteStringValue struct {
	LongValue   int64  `json:"longValue"`
	IntValue1   int32  `json:"intValue1"`
	IntValue2   int32  `json:"intValue2,omitempty"`
	ByteValue1  int8   `json:"byteValue1,omitempty"`
	ByteValue2  int8   `json:"byteValue2,omitempty"`
	StringValue string `json:"stringValue,omitempty"`
}

type IntBooleanIntBooleanValue struct {
	IntValue1  int32 `json:"intValue1"`
	BoolValue1 bool  `json:"boolValue1"`
	IntValue2  int32 `json:"intValue2"`
	BoolValue2 bool  `json:"boolValue2"`
}

type SpanEvent struct {
	Sequence      int             `json:"sequence"`
	StartElapsed  int             `json:"startElapsed"`
	EndElapsed    int             `json:"endElapsed"`
	ServiceType   int             `json:"serviceType"`
	EndPoint      string          `json:"endPoint"`
	Annotations   []*TempTag      `json:"annotations"`
	Depth         int             `json:"depth"`
	NextSpanID    int             `json:"nextSpanId"`
	DestinationID string          `json:"destinationId"`
	MethodID      int             `json:"apiId"`
	ExceptionInfo *IntStringValue `json:"exceptionInfo"`
}
type SpanEvents []*SpanEvent

func (o SpanEvents) Len() int {
	return len(o)
}

func (o SpanEvents) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o SpanEvents) Less(i, j int) bool {
	return o[i].Sequence < o[j].Sequence
}
