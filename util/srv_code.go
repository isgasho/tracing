package util

var SerCode map[int]string

func init() {
	SerCode = make(map[int]string)

	// activemq.client
	SerCode[8310] = "ACTIVEMQ_CLIENT"
	SerCode[8311] = "ACTIVEMQ_CLIENT_INTERNAL"

	// akka.http

	SerCode[1310] = "AKKA_HTTP_SERVER"
	SerCode[9998] = "1311"

	// arcus
	SerCode[8100] = "ARCUS"
	SerCode[8101] = "ARCUS_FUTURE_GET"
	SerCode[8102] = "ARCUS_EHCACHE_FUTURE_GET"

	SerCode[8103] = "ARCUS_INTERNAL"
	SerCode[8050] = "MEMCACHED"
	SerCode[8051] = "MEMCACHED_FUTURE_GET"

	// cassandra
	SerCode[2600] = "CASSANDRA"
	SerCode[2601] = "CASSANDRA_EXECUTE_QUERY"

	// cubrid
	SerCode[2400] = "CUBRID"
	SerCode[2401] = "CUBRID_EXECUTE_QUERY"

	// cxf
	SerCode[9080] = "CXF_CLIENT"

	// dbcp
	SerCode[6050] = "DBCP"
	// dbcp2
	SerCode[6052] = "DBCP2"

	// dubbo
	SerCode[1110] = "DUBBO_PROVIDER"
	SerCode[9110] = "DUBBO_CONSUMER"
	SerCode[9111] = "DUBBO"

	// httpclient
	SerCode[9054] = "GOOGLE_HTTP_CLIENT_INTERNAL"

	// gson
	SerCode[5010] = "GSON"

	// hikaricp
	SerCode[6060] = "HIKARICP"

	// httpclient3
	SerCode[9050] = "HTTP_CLIENT_3"

	// httpclient4
	SerCode[9052] = "HTTP_CLIENT_4"
	SerCode[9053] = "HTTP_CLIENT_4_INTERNAL"

	// hystrix
	SerCode[9120] = "HYSTRIX_COMMAND"
	SerCode[9121] = "HYSTRIX_COMMAND_INTERNAL"

	// ibatis
	SerCode[5500] = "IBATIS"
	SerCode[5501] = "IBATIS_SPRING"

	// jackson
	SerCode[5011] = "JACKSON"

	// jboss
	SerCode[1040] = "JBOSS"
	SerCode[1041] = "JBOSS_METHOD"

	// jdk.http
	SerCode[9055] = "JDK_HTTPURLCONNECTOR"

	// jetty
	SerCode[1030] = "JETTY"
	SerCode[1031] = "JETTY_METHOD"

	// json_lib
	SerCode[5012] = "JSON-LIB"

	// jsp
	SerCode[5005] = "JSP"

	// jtds
	SerCode[2200] = "MSSQLSERVER"
	SerCode[2201] = "MSSQL_EXECUTE_QUERY"
	// kafka
	SerCode[8660] = "KAFKA_CLIENT"
	SerCode[8661] = "KAFKA_CLIENT_INTERNAL"
	// mariadb
	SerCode[2150] = "MARIADB"
	SerCode[2151] = "MARIADB_EXECUTE_QUERY"
	// mybatis
	SerCode[5510] = "MYBATIS"

	// mysql
	SerCode[2100] = "MYSQL"
	SerCode[2101] = "MYSQL_EXECUTE_QUERY"

	// netty
	SerCode[9150] = "NETTY"
	SerCode[9151] = "NETTY_INTERNAL"
	SerCode[9152] = "NETTY_HTTP"

	// asynchttpclient
	SerCode[9056] = "ASYNC_HTTP_CLIENT"
	SerCode[9057] = "ASYNC_HTTP_CLIENT_INTERNAL"

	// okhttp
	SerCode[9058] = "OK_HTTP_CLIENT"
	SerCode[9059] = "OK_HTTP_CLIENT_INTERNAL"

	// oracle
	SerCode[2300] = "ORACLE"
	SerCode[2301] = "ORACLE_EXECUTE_QUERY"

	// php
	SerCode[1500] = "PHP"
	SerCode[1501] = "PHP_METHOD"
	SerCode[9700] = "PHP_REMOTE_METHOD"

	// postgresql
	SerCode[1500] = "PHP"
	SerCode[1501] = "PHP_METHOD"
	SerCode[9700] = "PHP_REMOTE_METHOD"

	// php
	SerCode[1500] = "PHP"
	SerCode[1501] = "PHP_METHOD"
	SerCode[9700] = "PHP_REMOTE_METHOD"

	// postgresql
	SerCode[2500] = "POSTGRESQL"
	SerCode[2501] = "POSTGRESQL_EXECUTE_QUERY"

	// rabbitmq.client
	SerCode[8300] = "RABBITMQ_CLIENT"
	SerCode[8301] = "RABBITMQ_CLIENT_INTERNAL"

}
