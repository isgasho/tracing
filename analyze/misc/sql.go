package misc

var QueryAgentStat string = `SELECT timestamp, stat_info  FROM agent_stats WHERE app_name=? AND  agent_id=? and timestamp>? and timestamp<=?;`

var InsertCPULoadRecord string = `INSERT INTO jvm_cpu_stats (app_name, agent_id, input_date, jvm, system) VALUES (?,?,?,?,?);`

var InsertJVMMemoryRecord string = `INSERT INTO jvm_memory_stats (app_name , agent_id , input_date , heap_used , non_heap ) VALUES (?,?,?,?,?);`

var QueryStartTime string = `SELECT start_time  FROM agent_stats  WHERE app_name=? LIMIT 1;`

var QueryAgents string = `SELECT agent_id, is_live FROM agents WHERE app_name=?;`

var InserAPIRecord string = `INSERT INTO api_stats (app_name, input_date, api, total_elapsed, max_elapsed, min_elapsed, average_elapsed, count, err_count, satisfaction, tolerate)
 VALUES (?,?,?,?,?,?,?,?,?,?,?);`

var CounterQuerySpan string = `SELECT app_name, input_date, api, elapsed,  service_type, parent_app_name,
parent_app_type, span_event_list, err, agent_id
FROM traces WHERE trace_id=? AND span_id=?;`

var ChunkEventsIterTrace string = `SELECT span_event_list FROM traces_chunk WHERE trace_id=? AND  span_id=?;`

var UpdateLastCounterTime string = `UPDATE apps SET last_count_time=? WHERE app_name=?;`

var InsertAPIs string = `INSERT INTO app_apis (app_name, api) VALUES (?, ?) ;`

var QueryTraceID string = `SELECT trace_id, span_id FROM app_operation_index WHERE app_name=? and input_date>? and input_date<=?;`

var InserRPCDetailsRecord string = ` INSERT INTO api_details_stats (app_name, api, input_date, method_id, service_type, elapsed, max_elapsed, min_elapsed, average_elapsed, count, err_count) VALUES (?,?,?,?,?,?,?,?,?,?,?);`

var InserExceptionRecord string = `INSERT INTO exception_stats (app_name, method_id, exception_info, input_date, total_elapsed, max_elapsed, 
	min_elapsed, count) VALUES (?,?,?,?,?,?,?,?);`

var InserSQLRecord string = `INSERT INTO sql_stats (app_name, sql, input_date, elapsed, max_elapsed, min_elapsed, average_elapsed, count, err_count) VALUES (?,?,?,?,?,?,?,?,?);`
