package sql

// 加载所有策略
var LoadPolicys string = `SELECT name, owner, api_alerts, channel, group,
 policy_id, update_date, users FROM alerts_app ;`

// 加载策略详情
var LoadAlert string = `SELECT alerts FROM alerts_policy WHERE id=?;`
