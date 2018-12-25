package service

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/thrift"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/pinpoint"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
	"github.com/mafanr/vgo/util"
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
func (p *Pinpoint) dealUpload(conn net.Conn, inPacket *util.VgoPacket) error {
	packet := &util.PinpointData{}
	if err := msgpack.Unmarshal(inPacket.Payload, packet); err != nil {
		g.L.Warn("dealUpload:msgpack.Unmarshal", zap.String("error", err.Error()))
		return err
	}

	switch packet.Type {
	case util.TypeOfTCPData:
		for _, value := range packet.Payload {
			switch value.Type {
			case util.TypeOfRegister:
				agentInfo := util.NewAgentInfo()
				if err := msgpack.Unmarshal(value.Spans, agentInfo); err != nil {
					g.L.Warn("dealUpload:msgpack.Unmarshal", zap.String("error", err.Error()))
					return err
				}

				if !gVgo.appStore.checkApp(agentInfo.AppName) {
					insertApp := `
					INSERT
					INTO apps(app_name)
					VALUES (?)`
					if err := gVgo.storage.session.Query(
						insertApp,
						agentInfo.AppName,
					).Exec(); err != nil {
						g.L.Warn("inster apps error", zap.String("error", err.Error()), zap.String("SQL", insertApp))
						return err
					}
				}

				if err := gVgo.storage.AgentStore(agentInfo); err != nil {
					g.L.Warn("dealUpload:storage.AgentStore", zap.String("error", err.Error()))
					return err
				}
				// 注册信息原样返回
				if _, err := conn.Write(inPacket.Encode()); err != nil {
					g.L.Warn("dealUpload:conn.Write", zap.String("error", err.Error()))
					return err
				}
				g.L.Info("agentInfo", zap.String("appName", agentInfo.AppName), zap.String("agentID", agentInfo.AgentID), zap.Bool("isLive", agentInfo.IsLive))
				break
			case util.TypeOfAgentOffline:
				// Agent下线处理
				agentInfo := util.NewAgentInfo()
				if err := msgpack.Unmarshal(value.Spans, agentInfo); err != nil {
					g.L.Warn("dealUpload:msgpack.Unmarshal", zap.String("error", err.Error()))
					return err
				}
				if err := gVgo.storage.AgentOffline(packet.AppName, packet.AgentID, agentInfo.StartTimestamp, agentInfo.EndTimestamp, agentInfo.IsLive); err != nil {
					g.L.Warn("dealUpload:storage.AgentOffline", zap.String("error", err.Error()))
					return err
				}
				// 注册信息原样返回
				if _, err := conn.Write(inPacket.Encode()); err != nil {
					g.L.Warn("dealUpload:conn.Write", zap.String("error", err.Error()))
					return err
				}
				g.L.Info("agentInfo", zap.String("appName", agentInfo.AppName), zap.String("agentID", agentInfo.AgentID), zap.Bool("isLive", agentInfo.IsLive))

				break
			case util.TypeOfAgentInfo, util.TypeOfSQLMetaData, util.TypeOfAPIMetaData, util.TypeOfStringMetaData:
				if err := p.DealTCPRequestResponse(packet, value.Spans); err != nil {
					g.L.Warn("dealUpload:p.DealTCPRequestResponse", zap.String("error", err.Error()))
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
			g.L.Warn("DealTCPRequestResponse:json.Marshal", zap.String("error", err.Error()))
			return err
		}
		if err := gVgo.storage.AgentInfoStore(packet.AppName, packet.AgentID, m.StartTimestamp, agentInfo); err != nil {
			g.L.Warn("DealTCPRequestResponse:gVgo.storage.AgentInfoStore", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TSqlMetaData:
		if err := gVgo.storage.AgentSQLStore(packet.AppName, packet.AgentID, m); err != nil {
			g.L.Warn("DealTCPRequestResponse:gVgo.storage.AgentSQLStore", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TApiMetaData:
		if err := gVgo.storage.AgentAPIStore(packet.AppName, packet.AgentID, m); err != nil {
			g.L.Warn("DealTCPRequestResponse:gVgo.storage.AgentAPIStore", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TStringMetaData:
		if err := gVgo.storage.AgentStringStore(packet.AppName, packet.AgentID, m); err != nil {
			g.L.Warn("DealTCPRequestResponse:gVgo.storage.AgentStringStore", zap.String("error", err.Error()))
			return err
		}
		break
	default:
		g.L.Warn("unknown type", zap.String("type", fmt.Sprintf("%t", m)))
	}
	return nil
}
