

CREATE KEYSPACE vgo_v1_datacenter WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  
    AND durable_writes = true;


CREATE TYPE IF NOT EXISTS ser_info (
	ser_info_str        TEXT,
	vm_args             LIST<TEXT>,
	ser_infos           frozen<ser_info>,
);

CREATE TYPE IF NOT EXISTS ser_meta_data (
	ser_info_str        TEXT,
	vm_args             LIST<TEXT>,
	ser_infos           LIST<frozen<ser_info>>,
);


CREATE TABLE IF NOT EXISTS agents (
    app_name     TEXT,
    agent_id     TEXT,                 
    ser_type     int,            
    socket_id    int,               
    host_name    TEXT,
    ip           TEXT,
    pid          INT,
    version      TEXT,
    start_time   bigint,
    end_time     bigint,
    is_live      BOOLEAN,
    is_container BOOLEAN,
    PRIMARY KEY (app_name, agent_id)
);


CREATE TABLE IF NOT EXISTS agent_info (
    app_name                TEXT,
    agent_id                TEXT,
    ser_meta_data           frozen<ser_meta_data>,
    PRIMARY KEY (app_name, agent_id)
);

-- CREATE TYPE IF NOT EXISTS keyvalue (
--     key             text,
--     value_type      text,
--     value_string    text,
--     value_bool      boolean,
--     value_long      bigint,
--     value_double    double,
--     value_binary    blob,
-- );


-- CREATE TYPE IF NOT EXISTS keyvalue_string (
--     key         text,
--     value       text,
-- );


-- CREATE TYPE IF NOT EXISTS log (
--     ts      bigint,
--     fields  list<frozen<keyvalue_string>>,
-- );

-- CREATE TYPE IF NOT EXISTS span_ref (
--     trace_segment_id    blob,
--     span_id             int,
--     ref_type            int,
-- );

-- CREATE TYPE IF NOT EXISTS cpu (
--     usage_percent   double,
-- );

-- CREATE TYPE IF NOT EXISTS memory (
--     is_heap     boolean,
-- 	init        bigint,
-- 	max         bigint,
-- 	used        bigint,
-- 	committed   bigint,
-- );


-- CREATE TYPE IF NOT EXISTS memory_pool (
-- 	type        int,
-- 	init        bigint,
-- 	max         bigint,
-- 	used        bigint,
-- 	commited    bigint,
-- );

-- CREATE TYPE IF NOT EXISTS gc (
-- 	phrase      int,
-- 	count       bigint,
-- 	time        bigint,
-- );

-- CREATE TYPE IF NOT EXISTS jvm (
--     time            bigint,
--     cpu             frozen<cpu>,
--     memory          list<frozen<memory>>,
--     memory_pool     list<frozen<memory_pool>>,
--     gc              list<frozen<gc>>,
-- );



-- CREATE TABLE IF NOT EXISTS jvms (
--     app_name text,                  // app_name
--     instance_id int,                // 实例id
--     report_time bigint,             // 上报时间
--     jvms  list<frozen<jvm>>,
--     PRIMARY KEY (app_name, instance_id, report_time)
-- );

-- CREATE TABLE IF NOT EXISTS traces (
--     trace_id            blob,
--     trace_segment_id    blob,
--     span_id             bigint,
--     app_id              bigint,
--     instance_id         bigint,
--     span_type           int,
--     span_layer          int,
--     start_time          bigint,
--     end_time            bigint,
--     parent_span_id      int,
--     operation_id        int,
--     is_error            boolean,
--     refs                list<frozen<span_ref>>,
--     tags                list<frozen<keyvalue_string>>,
--     logs                list<frozen<log>>,
--     PRIMARY KEY (trace_id, trace_segment_id, span_id)
-- );