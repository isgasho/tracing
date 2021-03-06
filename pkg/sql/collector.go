package sql

// App 名 信息入库
var InsertApp string = `
INSERT
INTO apps(app_name)
VALUES (?)`

// agent 信息入库
var InsertAgent string = `INSERT INTO agents (app_name, agent_id, service_type, 
	host_name, ip, start_time, end_time, is_container, operating_env, tracing_addr, is_live) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

// 更新agent 在线信息
var UpdateAgentState string = `UPDATE agents  SET is_live=?, end_time=? WHERE app_name =? AND agent_id =?;`

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

// insert runtime stat 信息入库
var InsertRuntimeStat string = `
	INSERT
	INTO agent_runtime(app_name, agent_id, input_date, metrics, runtime_type)
	VALUES (?, ?, ?, ?, ?);`

// agent stat 信息入库 + 过期时间
// var InsertAgentStatWithTTL string = `
// 	INSERT
// 	INTO agent_runtime(app_name, agent_id, input_date, metrics, runtime_type)
// 	VALUES (?, ?, ?, ?, 1) USING TTL ?;`

// 插入span
var InsertSpan string = `
	INSERT
	INTO traces(trace_id, span_id, agent_id, app_name, agent_start_time, parent_id,
		input_date, elapsed, api, service_type, end_point, remote_addr, annotations, error,
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
	INTO app_operation_index(app_name, agent_id, method_id, input_date, elapsed, trace_id, api, span_id, error, remote_addr)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

var InsertAPIs string = `INSERT INTO app_apis (app_name, api) VALUES (?, ?) ;`

// 插入服务类型
var InsertSrvType string = `
INSERT
INTO service_type(service_type, info)
VALUES (?, ?) ;`

// API 记录语句
var InsertAPIStats string = `INSERT INTO api_stats (app_name, input_date, api, total_elapsed, max_elapsed, min_elapsed, count, err_count, satisfaction, tolerate)
 VALUES (?,?,?,?,?,?,?,?,?,?);`

// InsertMethodStats ...
var InsertMethodStats string = ` INSERT INTO method_stats (app_name, api, input_date,
	 method_id, service_type, elapsed, max_elapsed, 
	 min_elapsed, count, err_count) VALUES (?,?,?,?,?,?,?,?,?,?);`

//  InserSQLStats ...
var InsertSQLStats string = `INSERT INTO sql_stats (app_name, sql, 
	input_date, elapsed, max_elapsed, min_elapsed, count, err_count) 
VALUES (?,?,?,?,?,?,?,?);`

// InsertExceptionStats ....
var InsertExceptionStats string = `INSERT INTO exception_stats (app_name, method_id, class_id, input_date, total_elapsed, max_elapsed, 
	min_elapsed, count, service_type) VALUES (?,?,?,?,?,?,?,?,?);`

// 父节点应用拓扑图入库
// var InsertParentMap string = `INSERT INTO parent_map (app_name, input_date,
// 	 service_type, parent_name, parent_type, req_recv_count, access_err_count, exception_count, duration)
// 	VALUES (?,?,?,?,?,?,?,?,?);`

var InsertParentMap string = `INSERT INTO service_map (
	source_name, 
	source_type,

	target_name,
	target_type,

	access_count,     
	access_err_count, 
	access_duration,  

	target_count,    
	target_err_count,
	input_date)
   VALUES (?,?,?,?,?,?,?,?,?,?);`

var InsertTargetMap string = `INSERT INTO service_map (
	source_name, 
	source_type,

	target_name,
	target_type,

	access_count,     
	access_err_count, 
	access_duration,  

	target_count,    
	target_err_count,
	input_date)
   VALUES (?,?,?,?,?,?,?,?,?,?);`

var InsertUnknowParentMap string = `INSERT INTO service_map (
	source_name, 
	source_type,

	target_name,
	target_type,

	access_count,     
	access_err_count, 
	access_duration,  

	target_count,    
	target_err_count,
	input_date)
   VALUES (?,?,?,?,?,?,?,?,?,?);`

// Api被调用统计信息
// var InsertAPICallStats string = `INSERT INTO api_call_stats (app_name, input_date, service_type, api_id, parent_name, req_recv_count, err_count, duration)
// VALUES (?,?,?,?,?,?,?,?);`

// Api被调用统计信息
var InsertAPIMapStats string = `INSERT INTO api_map (source_name, source_type, target_name, target_type, access_count, access_err_count, access_duration, api, input_date)
VALUES (?,?,?,?,?,?,?,?,?);`

var LoadAgents string = `SELECT service_type, agent_id, start_time, ip, is_live, host_name FROM agents WHERE app_name=?;`
var LoadApps string = `SELECT app_name FROM apps ;`

// //  InsertRuntimeStats ...
// var InsertRuntimeStats string = `INSERT INTO runtime_stats (app_name, agent_id,
// 	input_date, jvm_cpu_load, system_cpu_load, heap_used, non_heap, count)
// VALUES (?,?,?,?,?,?,?,?);`
