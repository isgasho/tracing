package service

import (
	"fmt"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
	"github.com/mafanr/vgo/util"
	"github.com/sunface/talent"
	"go.uber.org/zap"
)

// Storage ...
type Storage struct {
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{}
}

// Start ...
func (s *Storage) Start() error {
	return nil
}

// Close ...
func (s *Storage) Close() error {
	return nil
}

// AgentStore ...
func (s *Storage) AgentStore(agentInfo *util.AgentInfo) error {
	query := fmt.Sprintf(util.AgentInsert,
		agentInfo.AppName, agentInfo.AgentID, agentInfo.ServiceType, agentInfo.HostName,
		agentInfo.IP4S, agentInfo.Pid, agentInfo.Version, agentInfo.StartTimestamp, agentInfo.EndTimestamp,
		agentInfo.IsLive, agentInfo.IsContainer, agentInfo.OperatingEnv)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentInfoStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf(util.AgentUpdate, agentInfo.ServiceType, agentInfo.HostName,
			agentInfo.IP4S, agentInfo.Pid, agentInfo.Version, agentInfo.StartTimestamp, agentInfo.EndTimestamp,
			agentInfo.IsLive, agentInfo.IsContainer, agentInfo.OperatingEnv, agentInfo.AppName, agentInfo.AgentID)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}

// AgentInfoStore ...
func (s *Storage) AgentInfoStore(appName, agentID string, agentInfo []byte) error {
	query := fmt.Sprintf(util.AgentInfoInsert, agentInfo, appName, agentID)
	_, err := g.DB.Query(query)
	if err != nil {
		g.L.Warn("AgentInfoStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		return err
	}
	return nil
}

// AgentOffline ...
func (s *Storage) AgentOffline(appName, agentID string, endTime int64, isLive bool) error {
	query := fmt.Sprintf(util.AgentOffLine, isLive, endTime, appName, agentID)
	_, err := g.DB.Query(query)
	if err != nil {
		g.L.Warn("AgentOffline:g.DB.Query", zap.String("query", query), zap.Error(err))
		return err
	}
	return nil
}

// AgentAPIStore ...
func (s *Storage) AgentAPIStore(appName, agentID string, apiInfo *trace.TApiMetaData) error {
	query := fmt.Sprintf("insert into agent_api (app_name, agent_id, api_id, api_info, agent_start_time) values ('%s','%q','%d','%s',%d);",
		appName, agentID, apiInfo.ApiId, apiInfo.ApiInfo, apiInfo.AgentStartTime)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentAPIStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf("update  agent_api set api_info='%q', agent_start_time=%d where app_name='%s' and agent_id='%s' and api_id='%d';",
			apiInfo.ApiInfo, apiInfo.AgentStartTime, appName, agentID, apiInfo.ApiId)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentAPIStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}

// AgentSQLStore ...
func (s *Storage) AgentSQLStore(appName, agentID string, apiInfo *trace.TSqlMetaData) error {
	newSQL := g.B64.EncodeToString(talent.String2Bytes(apiInfo.Sql))
	query := fmt.Sprintf("insert into agent_sql (app_name, agent_id, sql_id, sql_info, agent_start_time) values ('%s','%s','%d','%q',%d);",
		appName, agentID, apiInfo.SqlId, newSQL, apiInfo.AgentStartTime)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentSQLStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf("update  agent_sql set sql_info='%q', agent_start_time=%d where app_name='%s' and agent_id='%s' and sql_id='%d';",
			newSQL, apiInfo.AgentStartTime, appName, agentID, apiInfo.SqlId)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentSQLStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}

// AgentStringStore ...
func (s *Storage) AgentStringStore(appName, agentID string, apiInfo *trace.TStringMetaData) error {
	query := fmt.Sprintf("insert into agent_str (app_name, agent_id, str_id, str_info, agent_start_time) values ('%s','%s','%d','%q',%d);",
		appName, agentID, apiInfo.StringId, apiInfo.StringValue, apiInfo.AgentStartTime)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentStringStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf("update  agent_str set str_info='%q', agent_start_time=%d where app_name='%s' and agent_id='%s' and str_id='%d';",
			apiInfo.StringValue, apiInfo.AgentStartTime, appName, agentID, apiInfo.StringId)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentStringStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}
