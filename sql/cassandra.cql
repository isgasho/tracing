CREATE KEYSPACE vgo_v1_datacenter WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  
    AND durable_writes = true;




CREATE TYPE IF NOT EXISTS value (
    key     int,
    value   frozen<value>,
);


CREATE TYPE IF NOT EXISTS annotation (
    key     int,
    value   frozen<value>,
);



CREATE TABLE IF NOT EXISTS agents_test (
    app_name            text,
    agent_id            text,
    agent_info          blob,
    PRIMARY KEY (app_name, agent_id)
);


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
    last_point_time     bigint,
    is_live             BOOLEAN,
    is_container        BOOLEAN,
    operating_env       int,
    agent_info          text,
    PRIMARY KEY ((app_name, agent_id), start_time)
) WITH CLUSTERING ORDER BY (start_time DESC);


CREATE TABLE IF NOT EXISTS agent_apis (
    app_name            text,
    agent_id            text,
    api_id              int,
    start_time          bigint,
    api_info            text,
    line                int,
    type                int,
    PRIMARY KEY ((app_name, agent_id, api_id),start_time)
) WITH CLUSTERING ORDER BY (start_time DESC);


CREATE TABLE IF NOT EXISTS agent_sqls (
    app_name            text,
    agent_id            text,
    sql_id              int,
    start_time          bigint,
    sql_info            text,
    PRIMARY KEY ((app_name, agent_id, sql_id), start_time)
) WITH CLUSTERING ORDER BY (start_time DESC);


CREATE TABLE IF NOT EXISTS agent_strs (
    app_name            text,
    agent_id            text,
    str_id              int,
    start_time          bigint,
    str_info            text,
    PRIMARY KEY ((app_name, agent_id, str_id), start_time)
)WITH CLUSTERING ORDER BY (start_time DESC);



CREATE TABLE IF NOT EXISTS traces (
    trace_id            text,
    span_id             bigint,
    agent_id            text,
    app_name            text,
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
    PRIMARY KEY ((trace_id, span_id), start_time)
) WITH CLUSTERING ORDER BY (start_time DESC);


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
);

CREATE TABLE IF NOT EXISTS app_names (
    app_name        text,
    agent_id        text,
    PRIMARY KEY (app_name, agent_id)
);

-- CREATE TABLE IF NOT EXISTS operation_apis (
--     app_name        text,
--     agent_id        text,
--     api_id          int,
--     PRIMARY KEY ((app_name, agent_id), api_id)
-- );


CREATE TABLE IF NOT EXISTS operation_apis (
    app_name        text,
    agent_id        text,
    api_id          int,
    start_time      bigint,
    PRIMARY KEY ((app_name, agent_id, api_id), start_time)
) WITH CLUSTERING ORDER BY (start_time DESC);


CREATE TABLE IF NOT EXISTS  app_operation_index (
    app_name        text,
    agent_id        text,
    api_id          int,
    start_time      bigint,
    trace_id        text,
    PRIMARY KEY ((app_name, agent_id, api_id), start_time)
) WITH CLUSTERING ORDER BY (start_time DESC);



CREATE TABLE IF NOT EXISTS agent_stats (
    app_name            text,
    agent_id            text,
    start_time          bigint,
    timestamp           bigint,
    stat_info           blob,
    PRIMARY KEY ((app_name, agent_id, start_time), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);


TRUNCATE  TABLE  agent_strs;
TRUNCATE  TABLE  agent_stats;
TRUNCATE  TABLE  traces_chunk;
TRUNCATE  TABLE  agents;
TRUNCATE  TABLE  traces;
TRUNCATE  TABLE  agent_apis;
TRUNCATE  TABLE  operation_apis;
TRUNCATE  TABLE  agent_sqls;
TRUNCATE  TABLE  app_operation_index;
TRUNCATE  TABLE  app_names;
            
          