CREATE KEYSPACE tracing_v1_datacenter WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  
    AND durable_writes = true;

USE tracing_v1_datacenter;

-- 用户表
CREATE TABLE IF NOT EXISTS account (
    id               text,
    name             text,
    mobile           text,
    email            text,
    message          text,
    app_show         tinyint, -- 1代表显示所有应用，2代表显示自己设定的应用
    app_names        text, -- 用户自定义显示的应用列表
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
CREATE TABLE IF NOT EXISTS admin (
    id text,
    priv text,
    PRIMARY KEY (id)
) WITH gc_grace_seconds = 10800;

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
    keys text, -- 为一些特殊指标使用，例如http code告警，此处就是code list
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
CREATE TABLE IF NOT EXISTS alerts_app (
    name             text,              -- 应用名
    owner            text,              -- 该应用告警的所有者ID
    policy_id        text,           -- 应用关联的告警策略模版
    group            text,              -- 应用关联的用户组
    channel          text,              -- 不使用用户组时，单独设置告警通道
    users            list<text>,        -- 不使用用户组时,单独设置通知的用户
    api_alerts     text,              -- 额外的api告警项设置(可以自行指定url，设置特殊的告警规则)
    update_date      bigint,            -- 记录更新时间
    PRIMARY KEY (name,owner)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON alerts_app (update_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

CREATE CUSTOM INDEX IF NOT EXISTS ON alerts_app (owner) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;

CREATE CUSTOM INDEX IF NOT EXISTS ON alerts_app (policy_id) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;

CREATE CUSTOM INDEX IF NOT EXISTS ON alerts_app (group) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' ;


CREATE TABLE IF NOT EXISTS agents (
    app_name            text,
    agent_id            text,
    service_type        int,
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
    tracing_addr        text,
    agent_info          text,
    PRIMARY KEY (app_name, agent_id)
) WITH gc_grace_seconds = 10800;



CREATE TABLE IF NOT EXISTS app_methods (
    app_name            text,
    method_id           int,
    method_info         text,
    line                int,
    type                int,
    PRIMARY KEY (app_name, method_id)
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



-- 全链路数据表
CREATE TABLE IF NOT EXISTS traces (
    trace_id            text,
    span_id             bigint,
    app_name            text,
    agent_id            text,
    agent_start_time    bigint,
    parent_id           bigint,
    input_date          bigint,
    elapsed             int,
    api                 text,
    service_type        int,
    end_point           text,
    remote_addr         text,
    annotations         BLOB,
    error               int,
    span_event_list     blob,
    parent_app_name     text,
    parent_app_type     int,
    acceptor_host       text,
    method_id           int,
    exception_info      blob,
    app_service_type    int,
    PRIMARY KEY (trace_id, span_id, input_date)
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
    -- last_count_time BIGINT,
    PRIMARY KEY (app_name)
) WITH gc_grace_seconds = 10800;  


CREATE TABLE IF NOT EXISTS  app_operation_index (
    app_name        text,
    agent_id        text,
    method_id       int,
    input_date      bigint,
    api             text,
    elapsed         int,
    trace_id        text,
    span_id         bigint,
    error           tinyint,            
    PRIMARY KEY (app_name, api, input_date, method_id, agent_id)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON app_operation_index (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};

-- agent runtime 信息表
CREATE TABLE IF NOT EXISTS agent_runtime (
    app_name            text,
    agent_id            text,
    runtime_type        int,
    input_date          bigint,
    metrics             blob,
    PRIMARY KEY (app_name, agent_id, runtime_type, input_date)
) WITH gc_grace_seconds = 10800; 


CREATE CUSTOM INDEX IF NOT EXISTS ON agent_runtime (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


CREATE TABLE IF NOT EXISTS  api_stats (
    app_name            text,
    input_date          bigint,
    api                 text,
    total_elapsed       int,
    max_elapsed         int,
    min_elapsed         int,
    count               int,
    err_count           int,
    satisfaction        int,    -- 满意
    tolerate            int,    -- 可容忍
    PRIMARY KEY (app_name, api, input_date)
) WITH gc_grace_seconds = 10800;  



CREATE CUSTOM INDEX IF NOT EXISTS ON api_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};



CREATE TABLE IF NOT EXISTS  exception_stats (
    app_name            text,
    method_id           int,
    class_id            INT,
    input_date          bigint,
    service_type        int,
    total_elapsed       int,
    max_elapsed         int,
    min_elapsed         int,
    count               int,
    PRIMARY KEY (app_name, method_id, class_id, input_date)
) WITH gc_grace_seconds = 10800;  


CREATE CUSTOM INDEX IF NOT EXISTS ON exception_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


-- app api表
CREATE TABLE IF NOT EXISTS  app_apis (
    app_name            text, -- app name
    api                 text, -- api信息
    PRIMARY KEY (app_name, api)
) WITH gc_grace_seconds = 10800;  



CREATE TABLE IF NOT EXISTS method_stats (
    app_name        text,
    api             text,
    method_id       int,
    input_date      bigint,
    service_type    int,
    elapsed         int,
    max_elapsed     int,
    min_elapsed     int,
    count           int,
    err_count       int,
    PRIMARY KEY (app_name, api, input_date, method_id)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON method_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


CREATE TABLE IF NOT EXISTS sql_stats (
    app_name        text,
    sql             int,
    input_date      bigint,
    elapsed         int,
    max_elapsed     int,
    min_elapsed     int,
    count           int,
    err_count       int,
    PRIMARY KEY (app_name, sql, input_date)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON sql_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


CREATE TABLE IF NOT EXISTS service_type (
    service_type        int,
    info                text,
    PRIMARY KEY (service_type)
) WITH gc_grace_seconds = 10800;
 

CREATE TABLE IF NOT EXISTS system (
    app_name        text,
    agent_id        text,
    input_date      bigint,
    metrics         blob,
    PRIMARY KEY (app_name, agent_id, input_date)
) WITH gc_grace_seconds = 10800;
 
-- 父节点应用拓扑图，这里的父节点都是接入了监控的
CREATE TABLE IF NOT EXISTS parent_map (
    app_name            text,               -- 应用名
    input_date          bigint,             -- 插入时间
    service_type        int,                -- 应用类型
    parent_name         text,               -- 父应用名
    parent_type         int,                -- 父应用类型
    req_recv_count      int,                -- 被请求次数
    err_count           int,                -- 错误次数
    total_elapsed       int,                -- 总耗时
    PRIMARY KEY (app_name, input_date, parent_name)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON parent_map (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


-- 子节点应用拓扑图，这里的子节点可能是数据库等信息
CREATE TABLE IF NOT EXISTS child_map (
    app_name            text,               -- 应用名
    input_date          bigint,             -- 插入时间
    service_type        int,                -- 应用类型
    child_type          int,                -- 子应用类型
    destinations        text,               -- 访问目标
    req_send_count      int,                -- 请求次数
    err_count           int,                -- 错误次数
    total_elapsed       int,                -- 总耗时
    PRIMARY KEY (app_name, input_date, child_type, destinations)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON child_map (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};




-- 未知类型或者未接入监控的parent拓扑图
CREATE TABLE IF NOT EXISTS unknow_parent_map (
    app_name            text,               -- 应用名
    input_date          bigint,             -- 插入时间
    service_type        int,                -- 应用类型
    req_recv_count      int,                -- 被请求次数
    err_count           int,                -- 错误次数
    total_elapsed       int,                -- 总耗时
    PRIMARY KEY (app_name, input_date)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON unknow_parent_map (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


-- api被应用调用统计表
CREATE TABLE IF NOT EXISTS api_call_stats (
    app_name            text,               -- 应用名
    input_date          bigint,             -- 插入时间
    service_type        int,                -- 应用类型
    api_id              int,                -- api id
    parent_name         text,               -- 访问者/父节点名称
    req_recv_count      int,                -- 被请求次数
    err_count           int,                -- 错误次数
    total_elapsed       int,                -- 总耗时
    PRIMARY KEY (app_name, input_date, api_id, parent_name)
) WITH gc_grace_seconds = 10800;

CREATE CUSTOM INDEX IF NOT EXISTS ON api_call_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};



CREATE TABLE IF NOT EXISTS runtime_stats (
    app_name                    text,
    agent_id                    text,
    input_date                  bigint,
    jvm_cpu_load                double,
    system_cpu_load             double,
    heap_used                   BIGINT,
    non_heap                    BIGINT,
    count                       int,
    PRIMARY KEY (app_name, agent_id, input_date)
) WITH gc_grace_seconds = 10800;  

CREATE CUSTOM INDEX IF NOT EXISTS ON runtime_stats (input_date) 
    USING 'org.apache.cassandra.index.sasi.SASIIndex' 
    WITH OPTIONS = {'mode': 'SPARSE'};


TRUNCATE TABLE alerts_group  ;
TRUNCATE TABLE alerts_policy ;
TRUNCATE TABLE account       ;
TRUNCATE TABLE admin         ;


TRUNCATE TABLE agent_stats   ;
TRUNCATE TABLE agents        ;
TRUNCATE TABLE login_count  ;


TRUNCATE TABLE api_stats          ;
TRUNCATE TABLE alerts_app          ;
TRUNCATE TABLE app_apis           ;
TRUNCATE TABLE app_methods        ;
TRUNCATE TABLE app_operation_index;
TRUNCATE TABLE app_sqls           ;

TRUNCATE TABLE app_strs         ;
TRUNCATE TABLE apps             ;
TRUNCATE TABLE api_call_stats           ;
TRUNCATE TABLE exception_stats  ;
TRUNCATE TABLE runtime_stats    ;
TRUNCATE TABLE unknow_parent_map;

TRUNCATE TABLE api_call_stats ;

TRUNCATE TABLE method_stats ;
TRUNCATE TABLE parent_map  ;
TRUNCATE TABLE service_type ;
TRUNCATE TABLE sql_stats    ;

TRUNCATE TABLE traces ;
TRUNCATE TABLE traces_chunk ;
