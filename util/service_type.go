package util

var ServiceType map[int]string

func init() {
	ServiceType = make(map[int]string)

	// activemq.client
	ServiceType[8310] = "ACTIVEMQ_CLIENT"
	ServiceType[8311] = "ACTIVEMQ_CLIENT_INTERNAL"

	// akka.http

	ServiceType[1310] = "AKKA_HTTP_SERVER"
	ServiceType[9998] = "1311"

	// arcus
	ServiceType[8100] = "ARCUS"
	ServiceType[8101] = "ARCUS_FUTURE_GET"
	ServiceType[8102] = "ARCUS_EHCACHE_FUTURE_GET"

	ServiceType[8103] = "ARCUS_INTERNAL"
	ServiceType[8050] = "MEMCACHED"
	ServiceType[8051] = "MEMCACHED_FUTURE_GET"

	// cassandra
	ServiceType[2600] = "CASSANDRA"
	ServiceType[2601] = "CASSANDRA_EXECUTE_QUERY"

	// cubrid
	ServiceType[2400] = "CUBRID"
	ServiceType[2401] = "CUBRID_EXECUTE_QUERY"

	// cxf
	ServiceType[9080] = "CXF_CLIENT"

	// dbcp
	ServiceType[6050] = "DBCP"
	// dbcp2
	ServiceType[6052] = "DBCP2"

	// dubbo
	ServiceType[1110] = "DUBBO_PROVIDER"
	ServiceType[9110] = "DUBBO_CONSUMER"
	ServiceType[9111] = "DUBBO"

	// httpclient
	ServiceType[9054] = "GOOGLE_HTTP_CLIENT_INTERNAL"

	// gson
	ServiceType[5010] = "GSON"

	// hikaricp
	ServiceType[6060] = "HIKARICP"

	// httpclient3
	ServiceType[9050] = "HTTP_CLIENT_3"

	// httpclient4
	ServiceType[9052] = "HTTP_CLIENT_4"
	ServiceType[9053] = "HTTP_CLIENT_4_INTERNAL"

	// hystrix
	ServiceType[9120] = "HYSTRIX_COMMAND"
	ServiceType[9121] = "HYSTRIX_COMMAND_INTERNAL"

	// ibatis
	ServiceType[5500] = "IBATIS"
	ServiceType[5501] = "IBATIS_SPRING"

	// jackson
	ServiceType[5011] = "JACKSON"

	// jboss
	ServiceType[1040] = "JBOSS"
	ServiceType[1041] = "JBOSS_METHOD"

	// jdk.http
	ServiceType[9055] = "JDK_HTTPURLCONNECTOR"

	// jetty
	ServiceType[1030] = "JETTY"
	ServiceType[1031] = "JETTY_METHOD"

	// json_lib
	ServiceType[5012] = "JSON-LIB"

	// jsp
	ServiceType[5005] = "JSP"

	// jtds
	ServiceType[2200] = "MSSQLSERVER"
	ServiceType[2201] = "MSSQL_EXECUTE_QUERY"
	// kafka
	ServiceType[8660] = "KAFKA_CLIENT"
	ServiceType[8661] = "KAFKA_CLIENT_INTERNAL"
	// mariadb
	ServiceType[2150] = "MARIADB"
	ServiceType[2151] = "MARIADB_EXECUTE_QUERY"
	// mybatis
	ServiceType[5510] = "MYBATIS"

	// mysql
	ServiceType[2100] = "MYSQL"
	ServiceType[2101] = "MYSQL_EXECUTE_QUERY"

	// netty
	ServiceType[9150] = "NETTY"
	ServiceType[9151] = "NETTY_INTERNAL"
	ServiceType[9152] = "NETTY_HTTP"

	// asynchttpclient
	ServiceType[9056] = "ASYNC_HTTP_CLIENT"
	ServiceType[9057] = "ASYNC_HTTP_CLIENT_INTERNAL"

	// okhttp
	ServiceType[9058] = "OK_HTTP_CLIENT"
	ServiceType[9059] = "OK_HTTP_CLIENT_INTERNAL"

	// oracle
	ServiceType[2300] = "ORACLE"
	ServiceType[2301] = "ORACLE_EXECUTE_QUERY"

	// php
	ServiceType[1500] = "PHP"
	ServiceType[1501] = "PHP_METHOD"
	ServiceType[9700] = "PHP_REMOTE_METHOD"

	// postgresql
	ServiceType[1500] = "PHP"
	ServiceType[1501] = "PHP_METHOD"
	ServiceType[9700] = "PHP_REMOTE_METHOD"

	// php
	ServiceType[1500] = "PHP"
	ServiceType[1501] = "PHP_METHOD"
	ServiceType[9700] = "PHP_REMOTE_METHOD"

	// postgresql
	ServiceType[2500] = "POSTGRESQL"
	ServiceType[2501] = "POSTGRESQL_EXECUTE_QUERY"

	// rabbitmq.client
	ServiceType[8300] = "RABBITMQ_CLIENT"
	ServiceType[8301] = "RABBITMQ_CLIENT_INTERNAL"

	// redis
	ServiceType[8200] = "REDIS"

	// resin
	ServiceType[1200] = "RESIN"
	ServiceType[1201] = "RESIN_METHOD"

	// resttemplate
	ServiceType[9140] = "REST_TEMPLATE"

	// rxjava
	ServiceType[6500] = "RX_JAVA"
	ServiceType[6501] = "RX_JAVA_INTERNAL"

	// spring.beans
	ServiceType[5071] = "SPRING_BEAN"
	// spring.async
	ServiceType[5052] = "SPRING_ASYNC"
	// spring.web
	ServiceType[5051] = "SPRING_MVC"

	// spring.boot
	ServiceType[1210] = "NAME"

	ServiceType[1100] = "THRIFT_SERVER"
	ServiceType[9100] = "THRIFT_CLIENT"
	ServiceType[1101] = "THRIFT_SERVER_INTERNAL"
	ServiceType[9101] = "THRIFT_CLIENT_INTERNAL"

	// tomcat
	ServiceType[1010] = "TOMCAT"
	ServiceType[1011] = "TOMCAT_METHOD"

	// tomcat
	ServiceType[1010] = "TOMCAT"
	ServiceType[1011] = "TOMCAT_METHOD"

	// undertow
	ServiceType[1120] = "UNDERTOW"
	ServiceType[1121] = "UNDERTOW_METHOD"

	// vertx
	ServiceType[1050] = "VERTX"
	ServiceType[1051] = "VERTX_INTERNAL"
	ServiceType[1052] = "VERTX_HTTP_SERVER"
	ServiceType[1053] = "VERTX_HTTP_SERVER_INTERNAL"
	ServiceType[9130] = "VERTX_HTTP_CLIENT"
	ServiceType[9131] = "VERTX_HTTP_CLIENT_INTERNAL"

	// weblogic
	ServiceType[1070] = "WEBLOGIC"
	ServiceType[1071] = "WEBLOGIC_METHOD"

	// weblogic
	ServiceType[1060] = "WEBSPHERE"
	ServiceType[1061] = "WEBSPHERE_METHOD"
}
