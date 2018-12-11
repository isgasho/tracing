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
    PRIMARY KEY (trace_id, span_id)
);


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