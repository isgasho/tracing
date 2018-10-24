

CREATE KEYSPACE vgo_v1_datacenter WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3'}  
    AND durable_writes = true;

CREATE TABLE IF NOT EXISTS app_code (
    name text,
    code int,
    PRIMARY KEY (name)
);