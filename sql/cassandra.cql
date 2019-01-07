CREATE KEYSPACE vgo_v1_datacenter WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  
    AND durable_writes = true;

CREATE TABLE IF NOT EXISTS agents (
    app_name            text,
    agent_id            text,
    ser_type            int,
    socket_id           int,
    host_name           text,
    ip                  text,
    pid                 int,
    version             text,
    start_time          bigint,
    end_time            bigint,
    is_live             BOOLEAN,
    is_container        BOOLEAN,
    operating_env       int,
    agent_info          text,
    PRIMARY KEY (app_name, agent_id, start_time)
) WITH gc_grace_seconds = 10800;



CREATE TABLE IF NOT EXISTS app_apis (
    app_name            text,
    api_id              int,
    start_time          bigint,
    api_info            text,
    line                int,
    type                int,
    PRIMARY KEY (app_name, api_id)
) WITH gc_grace_seconds = 10800;  -- 3 hours of downtime acceptable on nodes



CREATE TABLE IF NOT EXISTS app_sqls (
    app_name            text,
    sql_id              int,
    start_time          bigint,
    sql_info            text,
    PRIMARY KEY (app_name, sql_id)
) WITH gc_grace_seconds = 10800;


CREATE TABLE IF NOT EXISTS app_strs (
    app_name            text,
    str_id              int,
    start_time          bigint,
    str_info            text,
    PRIMARY KEY (app_name, str_id)
)WITH gc_grace_seconds = 10800; 



CREATE TABLE IF NOT EXISTS traces (
    trace_id            text,
    span_id             bigint,
    app_name            text,
    agent_id            text,
    agent_start_time    bigint,
    parent_id           bigint,
    start_time          bigint,
    elapsed             int,
    rpc                 text,
    service_type        int,
    end_point           text,
    remote_addr         text,
    annotations         BLOB,
    err                 int,
    span_event_list     blob,
    parent_app_name     text,
    parent_app_type     int,
    acceptor_host       text,
    api_id              int,
    exception_info      blob,
    app_service_type    int,
    PRIMARY KEY (trace_id, span_id, start_time)
) WITH gc_grace_seconds = 10800;


CREATE TABLE IF NOT EXISTS traces_chunk (
    trace_id            text,
    span_id             bigint,
    agent_id            text,
    app_name            text,
    service_type        int,
    end_point           text,
    span_event_list     blob,
    app_service_type    int,
    key_time            bigint,
    version             int,
    PRIMARY KEY (trace_id, span_id)
)WITH gc_grace_seconds = 10800; 


CREATE TABLE IF NOT EXISTS apps (
    app_name        text,
    last_count_time BIGINT,
    PRIMARY KEY (app_name)
) WITH gc_grace_seconds = 10800;  


CREATE TABLE IF NOT EXISTS  app_operation_index (
    app_name        text,
    agent_id        text,
    api_id          int,
    start_time      bigint,
    rpc             text,
    trace_id        text,
    span_id         bigint,
    PRIMARY KEY (app_name, agent_id, api_id, start_time)
) WITH gc_grace_seconds = 10800;  



CREATE CUSTOM INDEX IF NOT EXISTS ON app_operation_index (start_time) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


CREATE TABLE IF NOT EXISTS agent_stats (
    app_name            text,
    agent_id            text,
    start_time          bigint,
    timestamp           bigint,
    stat_info           blob,
    PRIMARY KEY (app_name, agent_id, timestamp)
) WITH gc_grace_seconds = 10800; 


CREATE CUSTOM INDEX IF NOT EXISTS ON agent_stats (timestamp) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


CREATE TABLE IF NOT EXISTS  rpc_record (
    app_name        text,
    start_time      bigint,
    url             text,
    elapsed         int,
    max_elapsed     int,
    min_elapsed     int,
    average_elapsed int,
    count           int,
    err_count       int,
    PRIMARY KEY (app_name, url, start_time)
) WITH gc_grace_seconds = 10800;  



CREATE CUSTOM INDEX IF NOT EXISTS ON rpc_record (start_time) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


CREATE TABLE IF NOT EXISTS rpc_details_record (
    app_name        text,
    url             text,
    api_id          int,
    start_time      bigint,
    ser_type        int,
    elapsed         int,
    ratio_elapsed   FLOAT,
    max_elapsed     int,
    min_elapsed     int,
    average_elapsed int,
    count           int,
    err_count       int,
    PRIMARY KEY (app_name, url, start_time, api_id)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON rpc_details_record (start_time) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE TABLE IF NOT EXISTS sql_record (
    app_name        text,
    sql             int,
    start_time      bigint,
    elapsed         int,
    max_elapsed     int,
    min_elapsed     int,
    average_elapsed int,
    count           int,
    err_count       int,
    PRIMARY KEY (app_name, sql, start_time)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON sql_record (start_time) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE TABLE IF NOT EXISTS cpu_load_record (
    app_name        text,
    agent_id        text,
    start_time      bigint,
    jvm             double,
    System          double,
    PRIMARY KEY (app_name, agent_id, start_time)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON cpu_load_record (start_time) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE TABLE IF NOT EXISTS jvm_memory_record (
    app_name        text,
    agent_id        text,
    start_time      bigint,
    heap_used       bigint,
    non_heap        bigint,
    PRIMARY KEY (app_name, agent_id, start_time)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON jvm_memory_record (start_time) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

TRUNCATE  TABLE  app_strs;
TRUNCATE  TABLE  agent_stats;
TRUNCATE  TABLE  traces_chunk;
TRUNCATE  TABLE  agents;
TRUNCATE  TABLE  traces;
TRUNCATE  TABLE  app_apis;
TRUNCATE  TABLE  app_sqls;
TRUNCATE  TABLE  app_operation_index;
TRUNCATE  TABLE  apps;
TRUNCATE  TABLE  rpc_record;
TRUNCATE  TABLE  rpc_details_record;
TRUNCATE  TABLE  sql_record;
-- TRUNCATE  TABLE apis;


DROP  TABLE  app_strs;
DROP  TABLE  agent_stats;
DROP  TABLE  traces_chunk;
DROP  TABLE  agents;
DROP  TABLE  traces;
DROP  TABLE  app_apis;
DROP  TABLE  app_sqls;
DROP  TABLE  app_operation_index;
DROP  TABLE  apps;
DROP  TABLE  rpc_record;
DROP  TABLE  rpc_details_record;
DROP  TABLE  sql_record;
-- DROP  TABLE  apis;
 

-- CREATE TABLE IF NOT EXISTS operation_apis (
--     app_name        text,
--     agent_id        text,
--     api_id          int,
--     PRIMARY KEY ((app_name, agent_id), api_id)
-- );


-- CREATE TABLE IF NOT EXISTS app_names (
--     app_name        text,
--     agent_id        text,
--     PRIMARY KEY (app_name, agent_id)
-- ) WITH gc_grace_seconds = 10800;  


-- CREATE TABLE IF NOT EXISTS operation_apis (
--     app_name        text,
--     agent_id        text,
--     api_id          int,
--     start_time      bigint,
--     PRIMARY KEY (app_name, agent_id, api_id, start_time)
-- ) WITH gc_grace_seconds = 10800;  


-- CREATE CUSTOM INDEX IF NOT EXISTS ON operation_apis (start_time) 
--     USING 'org.apache.cassandra.index.sasi.SASIIndex' 
--     WITH OPTIONS = {'mode': 'SPARSE'};





-- CREATE TABLE IF NOT EXISTS agent_apis (
--     app_name            text,
--     agent_id            text,
--     api_id              int,
--     start_time          bigint,
--     api_info            text,
--     line                int,
--     type                int,
--     PRIMARY KEY (app_name, agent_id, api_id, start_time)
-- ) WITH gc_grace_seconds = 10800;  -- 3 hours of downtime acceptable on nodes

