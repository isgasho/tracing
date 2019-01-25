package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
)

func queryTraces(c echo.Context) error {
	appName := c.FormValue("app_name")
	api := c.FormValue("api")
	minResp := c.FormValue("min_resp")
	maxResp := c.FormValue("max_resp")
	limit := c.FormValue("limit")

	// iter := web.Cql.Query(`SELECT id,name,owner,channel,users FROM alerts_group WHERE owner=?`, li.ID).Iter()

	// var id, name, owner, channel string
	// var users []string

	// groups := make([]*Group, 0)
	// for iter.Scan(&id, &name, &owner, &channel, &users) {
	// 	ownerNameR, _ := web.usersMap.Load(owner)
	// 	groups = append(groups, &Group{id, name, owner, ownerNameR.(*User).Name, channel, users})
	// }

	// CREATE TABLE IF NOT EXISTS  app_operation_index (
	// 	app_name        text,
	// 	agent_id        text,
	// 	api_id          int,
	// 	insert_date     bigint,
	// 	rpc             text,
	// 	trace_id        text,
	// 	span_id         bigint,
	// 	PRIMARY KEY (app_name, agent_id, api_id, insert_date)
	// ) WITH gc_grace_seconds = 10800;

	fmt.Println("an:", appName, "api:", api, "min:", minResp, "max:", maxResp, "limit:", limit)
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   tracesTest,
	})
}

func queryTrace(c echo.Context) error {
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   traceTest,
	})
}
