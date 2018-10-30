

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


CREATE TABLE IF NOT EXISTS jvm (
    app_name text,                  // app_name
    instance_id int,                // 实例id
    report_time bigint,             // 上报时间
    value blob,                     // jvm 数据
    PRIMARY KEY (app_name, instance_id, report_time)
);