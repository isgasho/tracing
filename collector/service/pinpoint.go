package service

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/pkg/network"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
)

// pinpointPacket 处理agent上报的监控数据
func pinpointPacket(conn net.Conn, tracePack *network.TracePack) error {
	packet := &network.SpansPacket{}
	if err := msgpack.Unmarshal(tracePack.Payload, packet); err != nil {
		g.L.Warn("msgpack.Unmarshal", zap.String("error", err.Error()))
		return err
	}

	switch packet.Type {
	case constant.TypeOfTCPData:
		for _, value := range packet.Payload {
			switch value.Type {
			case constant.TypeOfRegister:
				agentInfo := network.NewAgentInfo()
				if err := msgpack.Unmarshal(value.Spans, agentInfo); err != nil {
					g.L.Warn("msgpack Unmarshal", zap.String("error", err.Error()))
					return err
				}

				// 检查内存中是否存在app信息，不存在存入数据库中
				if !gCollector.apps.isAppExist(agentInfo.AppName) {
					if err := gCollector.storage.AppNameStore(agentInfo.AppName); err != nil {
						g.L.Warn("insert apps error", zap.String("error", err.Error()))
						return err
					}
				}

				if err := gCollector.storage.AgentStore(agentInfo); err != nil {
					g.L.Warn("agent Store", zap.String("error", err.Error()))
					return err
				}

				// 内存缓存Agent信息
				gCollector.apps.storeAgent(agentInfo.AppName, agentInfo.AgentID, agentInfo.StartTimestamp)

				// 注册信息原样返回
				if _, err := conn.Write(tracePack.Encode()); err != nil {
					g.L.Warn("conn.Write", zap.String("error", err.Error()))
					return err
				}

				break
			case constant.TypeOfAgentOffline:

				// 			// Agent下线处理
				// 			agentInfo := util.NewAgentInfo()
				// 			if err := msgpack.Unmarshal(value.Spans, agentInfo); err != nil {
				// 				g.L.Warn("msgpack.Unmarshal", zap.String("error", err.Error()))
				// 				return err
				// 			}

				// 			g.L.Info("AgentOffline", zap.String("appName", agentInfo.AppName), zap.String("agentID", agentInfo.AgentID), zap.Bool("isLive", agentInfo.IsLive))

				// 			// 清理内存缓存Agent信息
				// 			gVgo.appStore.RemoveAgent(agentInfo)

				// 			// 数据库中下线标志
				// 			if err := gVgo.storage.AgentOffline(packet.AppName, packet.AgentID, agentInfo.StartTimestamp, agentInfo.EndTimestamp, agentInfo.IsLive); err != nil {
				// 				g.L.Warn("storage.AgentOffline", zap.String("error", err.Error()))
				// 				return err
				// 			}
				// 			// 注册信息原样返回
				// 			if _, err := conn.Write(inPacket.Encode()); err != nil {
				// 				g.L.Warn("conn.Write", zap.String("error", err.Error()))
				// 				return err
				// 			}

				// 			g.L.Info("agentInfo", zap.String("appName", agentInfo.AppName), zap.String("agentID", agentInfo.AgentID), zap.Bool("isLive", agentInfo.IsLive))

				break
			case constant.TypeOfAgentInfo, constant.TypeOfSQLMetaData, constant.TypeOfAPIMetaData, constant.TypeOfStringMetaData:
				if err := tcpRequestResponse(packet, value.Spans); err != nil {
					g.L.Warn("tcpRequestResponse", zap.String("error", err.Error()))
					return err
				}
				break
			default:
				g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", value.Type)), zap.Uint16("type", value.Type))
				break
			}
		}
		break
	case constant.TypeOfUDPData:
		for _, value := range packet.Payload {
			udpRequest(packet.AppName, packet.AgentID, value.Spans)
		}
		break
	default:
		g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", packet.Type)), zap.Uint16("value", packet.Type))
	}
	return nil
}

// udpRequest udp报文处理
func udpRequest(appName, agentID string, data []byte) {
	tStruct := thrift.Deserialize(data)
	switch m := tStruct.(type) {
	case *trace.TSpan:
		gCollector.storage.SpanStore(m)
		gCollector.apps.routerSapn(appName, agentID, m)
		break
	case *trace.TSpanChunk:
		gCollector.storage.SpanChunkStore(m)
		break
	case *pinpoint.TAgentStat:
		if err := gCollector.storage.WriteAgentStat(appName, agentID, m, data); err != nil {
			g.L.Warn("agent stat", zap.String("error", err.Error()))
		}
		break
	case *pinpoint.TAgentStatBatch:
		if err := gCollector.storage.WriteAgentStatBatch(appName, agentID, m, data); err != nil {
			g.L.Warn("stat batch", zap.String("error", err.Error()))
		}
		break
	default:
		g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", m)))
	}
}

// tcpRequestResponse ...
func tcpRequestResponse(packet *network.SpansPacket, message []byte) error {
	tStruct := thrift.Deserialize(message)
	switch m := tStruct.(type) {
	case *pinpoint.TAgentInfo:
		agentInfo, err := json.Marshal(m)
		if err != nil {
			g.L.Warn("json.Marshal", zap.String("error", err.Error()))
			return err
		}
		if err := gCollector.storage.AgentInfoStore(packet.AppName, packet.AgentID, m.StartTimestamp, agentInfo); err != nil {
			g.L.Warn("agent info store", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TSqlMetaData:
		if err := gCollector.storage.AppSQLStore(packet.AppName, m); err != nil {
			g.L.Warn("sql store", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TApiMetaData:
		if err := gCollector.storage.AppMethodStore(packet.AppName, m); err != nil {
			g.L.Warn("api store", zap.String("error", err.Error()))
			return err
		}
		break
	case *trace.TStringMetaData:
		if err := gCollector.storage.AppStringStore(packet.AppName, m); err != nil {
			g.L.Warn("string store", zap.String("error", err.Error()))
			return err
		}
		break
	default:
		g.L.Warn("unknown type", zap.String("type", fmt.Sprintf("%t", m)))
	}
	return nil
}
