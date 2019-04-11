package sql

// App 名 信息入库
var InsertApp string = `
INSERT
INTO apps(app_name)
VALUES (?)`

// agent 信息入库
var InsertAgent string = `INSERT INTO agents (app_name, agent_id, service_type, 
	host_name, ip, start_time, end_time, is_live, is_container, operating_env, tracing_addr) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

// agent info 信息入库
var InsertAgentInfo string = `INSERT INTO agents (app_name, agent_id, start_time, agent_info) 
VALUES ( ?, ?, ?, ?);`

// sql语句 信息入库
var InsertSQL string = `INSERT INTO app_sqls (app_name, sql_id, sql_info) 
VALUES (?, ?, ?);`

// app method 信息入库
var InsertMethod string = `INSERT INTO app_methods (app_name, method_id, method_info, line, type) 
VALUES (?, ?, ?, ?, ?);`

// string 信息入库
var InsertString string = `INSERT INTO app_strs (app_name, str_id, str_info) 
VALUES (?, ?, ?);`

// agent stat 信息入库
var InsertAgentStat string = `
	INSERT
	INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
	VALUES (?, ?, ?, ?, ?);`

// agent stat 信息入库 + 过期时间
var InsertAgentStatWithTTL string = `
	INSERT
	INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
	VALUES (?, ?, ?, ?, ?) USING TTL ?;`

// 插入span
var InsertSpan string = `
	INSERT
	INTO traces(trace_id, span_id, agent_id, app_name, agent_start_time, parent_id,
		input_date, elapsed, api, service_type, end_point, remote_addr, annotations, err,
		span_event_list, parent_app_name, parent_app_type, acceptor_host, app_service_type, exception_info, method_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// 插入span chunk
var InsertSpanChunk string = `
INSERT
INTO traces_chunk(trace_id, span_id, agent_id, app_name, service_type, end_point,
	span_event_list, app_service_type, key_time, version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// 插入api索引
var InsertOperIndex string = `
	INSERT
	INTO app_operation_index(app_name, agent_id, method_id, input_date, elapsed, trace_id, api, span_id, error)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

var InsertAPIs string = `INSERT INTO app_apis (app_name, api) VALUES (?, ?) ;`

// 插入服务类型
var InsertSrvType string = `
INSERT
INTO service_type(service_type, info)
VALUES (?, ?) ;`
