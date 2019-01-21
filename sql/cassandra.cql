CREATE KEYSPACE vgo_v1_datacenter WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  
    AND durable_writes = true;

-- 用户表
CREATE TABLE IF NOT EXISTS account (
    id               text,
    name             text,
    mobile           text,
    email            text,
    message          text,
    last_login_date  text,
    PRIMARY KEY (id)
) WITH gc_grace_seconds = 10800;
-- 登录次数统计表
CREATE TABLE IF NOT EXISTS login_count (
    id               text,
    count          counter,
    PRIMARY KEY (id)
) WITH gc_grace_seconds = 10800;

-- 管理员表
CREATE TABLE IF NOT EXISTS admin （
    id text,
    priv text,
    PRIMARY KEY (id)
） WITH gc_grace_seconds = 10800;
INSERT INTO admin (id,priv) VALUES ('13269','super_admin');

-- 告警用户组表
CREATE TABLE IF NOT EXISTS alerts_group (
    id               UUID,          -- 唯一ID
    name             text,         -- 用户组名称
    channel          text,         -- 告警通道，支持'mobile' 'email' 'message'
    owner            text,         -- 当前组的所有者ID
    users            list<text>,   -- 该组的组员
    update_date      bigint,         -- 记录更新时间
    PRIMARY KEY (id,owner)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON alerts_group (owner) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;


-- 告警策略模版中的监控项
CREATE TYPE alert (
    name text, -- 监控项名称
    type text, -- 监控项类型： apm、system
    label text, -- 监控项描述
    compare tinyint, -- 比较类型 1: > 2:<  3:=
    unit text, -- 单位：%、个 
    duration tinyint, -- 持续时间, 1 代表1分钟
    value double,
);

-- 告警策略模版表
CREATE TABLE IF NOT EXISTS alerts_policy (
    id               UUID,          -- 唯一ID
    name             text,         -- 策略模版名
    owner            text,         -- 当前策略模版的所有者ID
    alerts           list<frozen<alert>>, -- 告警项列表
    update_date      bigint,         -- 记录更新时间
    PRIMARY KEY (id,owner)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON alerts_policy (update_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE CUSTOM INDEX IF NOT EXISTS ON alerts_policy (owner) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;

-- 应用告警表
CREATE TABLE IF NOT EXISTS app_alert (
    name             text,              -- 应用名
    owner            text,              -- 该应用告警的所有者ID
    policy           text,              -- 应用关联的告警策略模版
    group            text,              -- 应用关联的用户组
    channel          text,              -- 不使用用户组时，单独设置告警通道
    users            list<text>,        -- 不使用用户组时,单独设置通知的用户
    extra_alerts     map<text, frozen<alert>>,  -- 额外的url告警项设置(可以自行指定url，设置特殊的告警规则)
    update_date      bigint,            -- 记录更新时间
    PRIMARY KEY (name,owner)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON app_alert (update_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE CUSTOM INDEX IF NOT EXISTS ON app_alert (owner) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;

CREATE CUSTOM INDEX IF NOT EXISTS ON app_alert (policy) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;

CREATE CUSTOM INDEX IF NOT EXISTS ON app_alert (group) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;


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
    vgo_addr            text,
    agent_info          text,
    PRIMARY KEY (app_name, agent_id)
) WITH gc_grace_seconds = 10800;



CREATE TABLE IF NOT EXISTS app_apis (
    app_name            text,
    api_id              int,
    api_info            text,
    line                int,
    type                int,
    PRIMARY KEY (app_name, api_id)
) WITH gc_grace_seconds = 10800;  -- 3 hours of downtime acceptable on nodes



CREATE TABLE IF NOT EXISTS app_sqls (
    app_name            text,
    sql_id              int,
    sql_info            text,
    PRIMARY KEY (app_name, sql_id)
) WITH gc_grace_seconds = 10800;



CREATE TABLE IF NOT EXISTS app_strs (
    app_name            text,
    str_id              int,
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
    insert_date         bigint,
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
    PRIMARY KEY (trace_id, span_id, insert_date)
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
    insert_date     bigint,
    rpc             text,
    trace_id        text,
    span_id         bigint,
    PRIMARY KEY (app_name, agent_id, api_id, insert_date)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON app_operation_index (insert_date) 
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


CREATE TABLE IF NOT EXISTS  rpc_stats (
    app_name            text,
    input_date          bigint,
    url                 text,
    total_elapsed       int,
    max_elapsed         int,
    min_elapsed         int,
    average_elapsed     int,
    count               int,
    err_count           int,
    satisfaction        int,
    tolerate            int,
    PRIMARY KEY (app_name, url, input_date)
) WITH gc_grace_seconds = 10800;  



CREATE CUSTOM INDEX IF NOT EXISTS ON rpc_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


CREATE TABLE IF NOT EXISTS rpc_details_stats (
    app_name        text,
    url             text,
    api_id          int,
    input_date      bigint,
    ser_type        int,
    elapsed         int,
    ratio_elapsed   FLOAT,
    max_elapsed     int,
    min_elapsed     int,
    average_elapsed int,
    count           int,
    err_count       int,
    PRIMARY KEY (app_name, url, input_date, api_id)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON rpc_details_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE TABLE IF NOT EXISTS sql_stats (
    app_name        text,
    sql             int,
    input_date      bigint,
    elapsed         int,
    max_elapsed     int,
    min_elapsed     int,
    average_elapsed int,
    count           int,
    err_count       int,
    PRIMARY KEY (app_name, sql, input_date)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON sql_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE TABLE IF NOT EXISTS jvm_cpu_stats (
    app_name        text,
    agent_id        text,
    input_date      bigint,
    jvm             double,
    System          double,
    PRIMARY KEY (app_name, agent_id, input_date)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON jvm_cpu_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE TABLE IF NOT EXISTS jvm_memory_stats (
    app_name        text,
    agent_id        text,
    input_date      bigint,
    heap_used       bigint,
    non_heap        bigint,
    PRIMARY KEY (app_name, agent_id, input_date)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON jvm_memory_stats (input_date) 
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
TRUNCATE  TABLE  rpc_stats;
TRUNCATE  TABLE  rpc_details_stats;
TRUNCATE  TABLE  sql_stats;
TRUNCATE  TABLE  jvm_memory_stats;
TRUNCATE  TABLE  jvm_cpu_stats;

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
DROP  TABLE  rpc_stats;
DROP  TABLE  rpc_details_stats;
DROP  TABLE  sql_stats;
DROP  TABLE  jvm_memory_stats;
DROP  TABLE  jvm_cpu_stats;

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


-- CREATE TABLE IF NOT EXISTS  app_url_index (
--     app_name        text,
--     url_id          bigint,
--     insert_date     bigint,
--     agent_id        text,
--     api_id          int,
--     rpc             text,
--     trace_id        text,
--     span_id         bigint,
--     PRIMARY KEY (app_name, url_id, insert_date)
-- ) WITH gc_grace_seconds = 10800;  

-- CREATE CUSTOM INDEX IF NOT EXISTS ON app_url_index (insert_date) 
--     USING 'org.apache.cassandra.index.sasi.SASIIndex' 
--     WITH OPTIONS = {'mode': 'SPARSE'};
