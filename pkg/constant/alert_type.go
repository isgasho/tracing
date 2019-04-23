package constant

var Alert map[string]int
var AlertInfo map[int]string

const (
	APM_APDEX_COUNT         = 1
	APM_http_code_ratio     = 2
	APM_http_code_count     = 3
	APM_api_error_ratio     = 4
	APM_sql_error_ratio     = 5
	APM_api_duration        = 6
	APM_sql_duration        = 7
	APM_jvm_fullgc_count    = 8
	APM_api_count           = 9
	SYSTEM_cpu_used_ratio   = 10
	SYSTEM_load_count       = 11
	SYSTEM_mem_used_ratio   = 12
	SYSTEM_disk_used_ratio  = 13
	SYSTEM_syn_recv_count   = 14
	SYSTEM_time_wait_count  = 15
	SYSTEM_diskio_ratio     = 16
	SYSTEM_ifstat_out_speed = 17
	SYSTEM_close_wait_count = 18
	SYSTEM_ifstat_in_speed  = 19
	SYSTEM_estab_count      = 20

	ALERT_TYPE_API = 1000
	ALERT_TYPE_SQL = 1001
)

func initAlertType() {
	Alert = make(map[string]int)
	AlertInfo = make(map[int]string)

	Alert["apm.apdex.count"] = 1
	AlertInfo[1] = "综合健康指数Apdex"

	Alert["apm.http_code.ratio"] = 2
	AlertInfo[2] = "错误HTTP CODE比率"

	Alert["apm.http_code.count"] = 3
	AlertInfo[3] = "错误HTTP CODE次数"

	Alert["apm.api_error.ratio"] = 4
	AlertInfo[4] = "接口错误率"

	Alert["apm.sql_error.ratio"] = 5
	AlertInfo[5] = "sql错误率"

	Alert["apm.api.duration"] = 6
	AlertInfo[6] = "接口平均耗时"

	Alert["apm.sql.duration"] = 7
	AlertInfo[7] = "sql平均耗时"

	Alert["apm.jvm_fullgc.count"] = 8
	AlertInfo[8] = "JVMFullGC报警"

	// Alert["apm.api.count"] = 9
	// AlertInfo[9] = "接口访问次数"

	Alert["system.cpu_used.ratio"] = 10
	AlertInfo[10] = "cpu使用率"

	Alert["system.load.count"] = 11
	AlertInfo[11] = "系统Load"

	Alert["system.mem_used.ratio"] = 12
	AlertInfo[12] = "内存使用率"

	Alert["system.disk_used.ratio"] = 13
	AlertInfo[13] = "硬盘使用率"

	Alert["system.syn_recv.count"] = 14
	AlertInfo[14] = "sync_recv数"

	Alert["system.time_wait.count"] = 15
	AlertInfo[15] = "time_wait数"

	Alert["system.diskio.ratio"] = 16
	AlertInfo[16] = "diskio利用率"

	Alert["system.ifstat_out.speed"] = 17
	AlertInfo[17] = "网络out速度"

	Alert["system.close_wait.count"] = 18
	AlertInfo[18] = "close_wait数"

	Alert["system.ifstat_in.speed"] = 19
	AlertInfo[19] = "网络in速度"

	Alert["system.estab.count"] = 20
	AlertInfo[20] = "建立长链接数"
}

// AlertType 通过描述获取类型
func AlertType(desc string) (int, bool) {
	alertType, ok := Alert[desc]
	return alertType, ok
}

// AlertDesc 通过类型获取描述
func AlertDesc(alertType int) (string, bool) {
	desc, ok := AlertInfo[alertType]
	return desc, ok
}
