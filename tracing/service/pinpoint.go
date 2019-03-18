package service

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/proto/pinpoint/thrift"
	"github.com/imdevlab/tracing/proto/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/proto/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/util"
	"github.com/imdevlab/tracing/tracing/misc"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
)

// Pinpoint ...
type Pinpoint struct {
}

// NewPinpoint ...
func NewPinpoint() *Pinpoint {
	return &Pinpoint{}
}

// dealUpload 处理pinpoint上行数据
func (p *Pinpoint) dealUpload(conn net.Conn, inPacket *util.TracingPacket) error {
	packet := &util.PinpointData{}
	if err := msgpack.Unmarshal(inPacket.Payload, packet); err != nil {
		g.L.Warn("msgpack.Unmarshal", zap.String("error", err.Error()))
		return err
	}

	switch packet.Type {
	case util.TypeOfTCPData:
		for _, value := range packet.Payload {
			switch value.Type {
			case util.TypeOfRegister:
				agentInfo := util.NewAgentInfo()
				if err := msgpack.Unmarshal(value.Spans, agentInfo); err != nil {
					g.L.Warn("msgpack.Unmarshal", zap.String("error", err.Error()))
					return err
				}

				if !gVgo.appStore.checkApp(agentInfo.AppName) {
					query := gVgo.storage.cql.Query(
						misc.InsertApp,
						agentInfo.AppName,
					)
					if err := query.Exec(); err != nil {
						g.L.Warn("inster apps error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
						return err
					}
				}

				if err := gVgo.storage.AgentStore(agentInfo); err != nil {
					g.L.Warn("storage.AgentStore", zap.String("error", err.Error()))
					return err
				}

				// 内存缓存Agent信息
				gVgo.appStore.StoreAgent(agentInfo, conn)

				// 注册信息原样返回
				if _, err := conn.Write(inPacket.Encode()); err != nil {
					g.L.Warn("conn.Write", zap.String("error", err.Error()))
					return err
				}

				// gVgo.appStore.Apps

				g.L.Info("agentInfo", zap.String("appName", agentInfo.AppName), zap.String("agentID", agentInfo.AgentID), zap.Bool("isLive", agentInfo.IsLive))
				break
			case util.TypeOfAgentOffline:

				// Agent下线处理
				agentInfo := util.NewAgentInfo()
				if err := msgpack.Unmarshal(value.Spans, agentInfo); err != nil {
					g.L.Warn("msgpack.Unmarshal", zap.String("error", err.Error()))
					return err
				}

				g.L.Info("AgentOffline", zap.String("appName", agentInfo.AppName), zap.String("agentID", agentInfo.AgentID), zap.Bool("isLive", agentInfo.IsLive))

				// 清理内存缓存Agent信息
				gVgo.appStore.RemoveAgent(agentInfo)

				// 数据库中下线标志
				if err := gVgo.storage.AgentOffline(packet.AppName, packet.AgentID, agentInfo.StartTimestamp, agentInfo.EndTimestamp, agentInfo.IsLive); err != nil {
					g.L.Warn("storage.AgentOffline", zap.String("error", err.Error()))
					return err
				}
				// 注册信息原样返回
				if _, err := conn.Write(inPacket.Encode()); err != nil {
					g.L.Warn("conn.Write", zap.String("error", err.Error()))
					return err
				}

				g.L.Info("agentInfo", zap.String("appName", agentInfo.AppName), zap.String("agentID", agentInfo.AgentID), zap.Bool("isLive", agentInfo.IsLive))

				break
			case util.TypeOfAgentInfo, util.TypeOfSQLMetaData, util.TypeOfAPIMetaData, util.TypeOfStringMetaData:
				if err := p.DealTCPRequestResponse(packet, value.Spans); err != nil {
					g.L.Warn("DealTCPRequestResponse", zap.String("error", err.Error()))
					return err
				}
				break
			default:
				g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", value.Type)), zap.Uint16("type", value.Type))
				break
			}
		}
		break
	case util.TypeOfUDPData:
		for _, value := range packet.Payload {
			p.DealUDPRequestResponse(packet.AppName, packet.AgentID, value.Spans)
		}
		break
	default:
		g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", packet.Type)))
	}
	return nil
}

// DealUDPRequestResponse ...
func (p *Pinpoint) DealUDPRequestResponse(appName, agentID string, data []byte) {
	tStruct := thrift.Deserialize(data)
	switch m := tStruct.(type) {
	case *trace.TSpan:
		gVgo.storage.spanChan <- m
		break
	case *trace.TSpanChunk:
		gVgo.storage.spanChunkChan <- m
		break
	case *pinpoint.TAgentStat:
		if err := gVgo.storage.writeAgentStat(appName, agentID, m, data); err != nil {
			g.L.Warn("writeAgentStat error", zap.String("error", err.Error()))
		}
		break
	case *pinpoint.TAgentStatBatch:
		if err := gVgo.storage.writeAgentStatBatch(appName, agentID, m, data); err != nil {
			g.L.Warn("writeAgentStatBatch error", zap.String("error", err.Error()))
		}
		break
	default:
		g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", m)))
	}
}

// DealTCPRequestResponse ...
func (p *Pinpoint) DealTCPRequestResponse(packet *util.PinpointData, message []byte) error {
	tStruct := thrift.Deserialize(message)
	switch m := tStruct.(type) {
	case *pinpoint.TAgentInfo:
		agentInfo, err := json.Marshal(m)
		if err != nil {
			g.L.Warn("json.Marshal", zap.String("error", err.Error()))
			return err
		}
		if err := gVgo.storage.AgentInfoStore(packet.AppName, packet.AgentID, m.StartTimestamp, agentInfo); err != nil {
			g.L.Warn("AgentInfoStore", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TSqlMetaData:
		if err := gVgo.storage.AppSQLStore(packet.AppName, m); err != nil {
			g.L.Warn("AppSQLStore", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TApiMetaData:
		if err := gVgo.storage.AppAPIStore(packet.AppName, m); err != nil {
			g.L.Warn("AppAPIStore", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TStringMetaData:
		if err := gVgo.storage.AppStringStore(packet.AppName, m); err != nil {
			g.L.Warn("AppStringStore", zap.String("error", err.Error()))
			return err
		}
		break
	default:
		g.L.Warn("unknown type", zap.String("type", fmt.Sprintf("%t", m)))
	}
	return nil
}
