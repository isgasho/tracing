

CREATE KEYSPACE vgo_v1_datacenter WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  
    AND durable_writes = true;

-- CREATE TABLE IF NOT EXISTS app_id (
--     name text,
--     id int,
--     PRIMARY KEY (name)
-- );


CREATE TYPE IF NOT EXISTS keyvalue (
    key             text,
    value_type      text,
    value_string    text,
    value_bool      boolean,
    value_long      bigint,
    value_double    double,
    value_binary    blob,
);


CREATE TYPE IF NOT EXISTS keyvalue_string (
    key         text,
    value       text,
);


CREATE TYPE IF NOT EXISTS log (
    ts      bigint,
    fields  list<frozen<keyvalue_string>>,
);

CREATE TYPE IF NOT EXISTS span_ref (
    trace_id        blob,
    span_id         int,
    ref_type        int,
);


CREATE TABLE IF NOT EXISTS jvms (
    app_name text,                  // app_name
    instance_id int,                // 实例id
    report_time bigint,             // 上报时间
    value blob,                     // jvm 数据
    PRIMARY KEY (app_name, instance_id, report_time)
);

CREATE TABLE IF NOT EXISTS traces (
    trace_id        blob,
    span_id         bigint,
    app_id          bigint,
    instance_id     bigint,
    span_type       int,
    span_layer      int,
    start_time      bigint,
    end_time        bigint,
    parent_span_id  int,
    operation_id    int,
    is_error        boolean,
    refs            list<frozen<span_ref>>,
    tags            list<frozen<keyvalue_string>>,
    logs            list<frozen<log>>,
    PRIMARY KEY (trace_id, span_id)
);