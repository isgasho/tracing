package misc

var AgentInsert string = `INSERT INTO agents (app_name, agent_id, ser_type, socket_id, host_name, ip,
	pid, version, start_time, end_time, is_live, is_container, operating_env) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

var AgentInofInsert string = `INSERT INTO agents (app_name, agent_id, start_time, agent_info) 
VALUES ( ?, ?, ?, ?);`

var AgentOfflineInsert string = `INSERT INTO agents (app_name, agent_id, end_time, is_live) 
VALUES ( ?, ?, ?, ?);`

var AppAPIInsert string = `INSERT INTO app_apis (app_name, api_id, api_info, line, type) 
VALUES (?, ?, ?, ?, ?);`

var APIInsert string = `INSERT INTO apis (api_id, api_info, line, type) 
VALUES (?, ?, ?, ?);`

var AppSQLInsert string = `INSERT INTO app_sqls (app_name, sql_id, sql_info) 
VALUES (?, ?, ?);`

var AgentStrInsert string = `INSERT INTO app_strs (app_name, str_id, str_info) 
VALUES (?, ?, ?);`

var InsertSpan string = `
INSERT
INTO traces(trace_id, span_id, agent_id, app_name, agent_start_time, parent_id,
	insert_date, elapsed, rpc, service_type, end_point, remote_addr, annotations, err,
	span_event_list, parent_app_name, parent_app_type, acceptor_host, app_service_type, exception_info, api_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

var InsertSpanChunk string = `
INSERT
INTO traces_chunk(trace_id, span_id, agent_id, app_name, service_type, end_point,
	span_event_list, app_service_type, key_time, version)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

var InsertOperIndex string = `
	INSERT
	INTO app_operation_index(app_name, agent_id, api_id, insert_date, trace_id, rpc, span_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

var InsertAgentStat string = `
	INSERT
	INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
	VALUES (?, ?, ?, ?, ?);`

var InsertAgentStatTTL string = `
	INSERT
	INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
	VALUES (?, ?, ?, ?, ?) USING TTL ?;`
