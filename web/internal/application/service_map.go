package app

import (
	"net/http"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/web/internal/misc"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// 这里的ErrorCount和SpanCount是为了展示当前节点的内部异常情况
type Node struct {
	Name       string `"json:"name"`
	SpanCount  int    `json:"span_count"`  // 收到的请求次数
	ErrorCount int    `json:"error_count"` // 错误次数
}

// 这里的ErrorCount和ReqCount是为了计算请求错误率
// A -> B
// Source: A, Target B
// A访问B的所有错误都计算在ErrorCount中，包含网络不通引起的错误
type NodeLink struct {
	Source      string `json:"source"`       // 起点app name
	Target      string `json:"target"`       // 终点app name
	AccessCount int    `json:"access_count"` // 请求次数
	ErrorCount  int    `json:"error_count"`  // 错误次数

	duration        int // 总耗时
	AverageDuration int `json:"avg"` // 平均耗时
}

type ServiceMapResult struct {
	Nodes []*Node     `json:"nodes"`
	Links []*NodeLink `json:"links"`
}

// 通过App name来查询service map
func QueryAPPServiceMap(c echo.Context) error {
	tname := c.FormValue("app_name")
	start, end, err := misc.StartEndDate(c)
	if err != nil {
		g.L.Info("日期参数不合法", zap.String("start", c.FormValue("start")), zap.String("end", c.FormValue("end")), zap.Error(err))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusOK,
			ErrCode: g.ParamInvalidC,
			Message: "日期参数不合法",
		})
	}

	nodeMap := make(map[string]*Node)
	linkMap := make(map[string]*NodeLink)

	// 获取当前APP的父应用信息
	q := misc.Cql.Query(`SELECT source_name,source_type,access_count,access_err_count,access_duration FROM service_map WHERE target_name = ? and input_date > ? and input_date < ?  ALLOW FILTERING`, tname, start.Unix(), end.Unix())
	iter := q.Iter()

	var sname string
	var stype, accessCount, accessErr, accessDuration int
	for iter.Scan(&sname, &stype, &accessCount, &accessErr, &accessDuration) {
		// 我们需要为每个APP区别不同的unknow，否则不同app对应的unkown请求涞源都会合并，数据就会不正确
		// 更新父节点信息
		_, ok := nodeMap[sname]
		if !ok {
			nodeMap[sname] = &Node{Name: sname}
		}
		// 更新子节点信息
		_, ok = nodeMap[tname]
		if !ok {
			nodeMap[tname] = &Node{Name: sname}
		}
		// 更新父子之间的link信息
		lname := sname + tname
		link, ok := linkMap[lname]
		if !ok {
			linkMap[lname] = &NodeLink{sname, tname, accessCount, accessErr, accessDuration, 0}
		} else {
			link.AccessCount += accessCount
			link.ErrorCount += accessErr
			link.duration += accessDuration
		}
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	// 获取当前APP的子应用信息
	q = misc.Cql.Query(`SELECT target_name,target_type,access_count,access_err_count,access_duration FROM service_map WHERE source_name = ? and input_date > ? and input_date < ?`, tname, start.Unix(), end.Unix())
	iter = q.Iter()

	// 当前应用变成了源应用
	sname = tname
	var ttype int
	for iter.Scan(&tname, &ttype, &accessCount, &accessErr, &accessDuration) {
		// 更新子节点信息
		_, ok := nodeMap[tname]
		if !ok {
			nodeMap[tname] = &Node{Name: tname}
		}

		// 更新父子之间的link信息
		lname := sname + tname
		link, ok := linkMap[lname]
		if !ok {
			linkMap[lname] = &NodeLink{sname, tname, accessCount, accessErr, accessDuration, 0}
		} else {
			link.AccessCount += accessCount
			link.ErrorCount += accessErr
			link.duration += accessDuration
		}
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	nodes := make([]*Node, 0, len(nodeMap))
	links := make([]*NodeLink, 0, len(linkMap))
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	for _, link := range linkMap {
		link.AverageDuration = link.duration / link.AccessCount
		links = append(links, link)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data: ServiceMapResult{
			Nodes: nodes,
			Links: links,
		},
	})
}

// 查询全局service map
func QueryServiceMap(c echo.Context) error {
	return nil
}
