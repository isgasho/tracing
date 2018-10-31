package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/agent/misc"
	"github.com/mafanr/vgo/agent/protocol"
	"github.com/mafanr/vgo/util"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// SkyWalking ...
type SkyWalking struct {
	appRegisterC chan *appRegister
	jvmChan      chan *protocol.JVMMetrics
	upTracChan   chan *protocol.UpstreamSegment
}

// NewSkyWalking ...
func NewSkyWalking() *SkyWalking {
	return &SkyWalking{
		appRegisterC: make(chan *appRegister, 10),
		jvmChan:      make(chan *protocol.JVMMetrics, misc.Conf.SkyWalking.JVMCacheLen+100),
		upTracChan:   make(chan *protocol.UpstreamSegment, misc.Conf.SkyWalking.TraceCacheLen+100),
	}
}

// JVMCollector jvm 信息采集，JVMReportInterval秒上报一次
func (sky *SkyWalking) JVMCollector() {
	jvmQueue := &util.JVMS{
		AppName: gAgent.appName,
		JVMs:    make([]*util.JVM, 0),
	}

	isHaveInstanceID := false
	if gAgent.appInstanceID != 0 {
		isHaveInstanceID = true
	}

	ticker := time.NewTicker(time.Duration(misc.Conf.SkyWalking.JVMReportInterval) * time.Second)
	start := time.Now()
	for {
		select {
		case jvmPack, ok := <-sky.jvmChan:
			if ok {
				end := time.Now()
				// 采集频率控制
				if end.Sub(start).Seconds() >= float64(misc.Conf.SkyWalking.JVMReportInterval) {
					// 时间更新
					start = end
					if !isHaveInstanceID {
						jvmQueue.InstanceID = jvmPack.ApplicationInstanceId
						// 本地也缓存一次
						gAgent.appInstanceID = jvmPack.ApplicationInstanceId
						// 保证只进入一次
						isHaveInstanceID = true
					}

					jvmQueue.Time = end.UnixNano() / 1e6
					// analysis jvm
					for _, metric := range jvmPack.Metrics {
						reportJVM := analysisJVM(metric)
						jvmQueue.JVMs = append(jvmQueue.JVMs, reportJVM)
					}

					if len(jvmQueue.JVMs) > misc.Conf.SkyWalking.JVMCacheLen {
						// 发送
						if err := sky.sendJVMs(jvmQueue); err != nil {
							g.L.Warn("JVMCollector:sky.sendJVMs", zap.String("error", err.Error()))
						}
						// 清空缓存
						jvmQueue.JVMs = jvmQueue.JVMs[:0]
					}
				}
			}
			break
		case <-ticker.C:
			if len(jvmQueue.JVMs) > 0 {
				// 发送
				if err := sky.sendJVMs(jvmQueue); err != nil {
					g.L.Warn("JVMCollector:sky.sendJVMs", zap.String("error", err.Error()))
				}
				// 清空缓存
				jvmQueue.JVMs = jvmQueue.JVMs[:0]
			}
			break
		}
	}
}

// analysisSpan ...
func analysisSpan(appID, instanceID int32, tracID string, skySpan *protocol.SpanObject) *util.Span {
	var refs []*util.SpanRef
	for _, skyRef := range skySpan.Refs {
		if len(skyRef.ParentTraceSegmentId.IdParts) != 3 {
			continue
		}
		tracID := fmt.Sprintf("%d.%d.%d", skyRef.ParentTraceSegmentId.IdParts[0], skyRef.ParentTraceSegmentId.IdParts[1], skyRef.ParentTraceSegmentId.IdParts[2])
		ref := &util.SpanRef{
			TraceID: tracID,
			SpanID:  skyRef.ParentSpanId,
			RefType: util.RefType(skyRef.RefType),
		}
		refs = append(refs, ref)
	}

	var tags []*util.KeyWithStringValue
	for _, skyTag := range skySpan.Tags {
		tag := &util.KeyWithStringValue{
			Key:   skyTag.Key,
			Value: skyTag.Value,
		}
		tags = append(tags, tag)
	}

	var newLogs []*util.LogMessage
	for _, skylog := range skySpan.Logs {
		var datas []*util.KeyWithStringValue
		for _, skyData := range skylog.Data {
			data := &util.KeyWithStringValue{
				Key:   skyData.Key,
				Value: skyData.Value,
			}
			datas = append(datas, data)
		}
		newLog := &util.LogMessage{
			Time: skylog.Time,
			Data: datas,
		}
		newLogs = append(newLogs, newLog)
	}

	span := &util.Span{
		TraceID:         tracID,
		SpanID:          skySpan.SpanId,
		AppID:           appID,
		InstanceID:      instanceID,
		SpanType:        util.SpanType(skySpan.SpanType),
		SpanLayer:       util.SpanLayer(skySpan.SpanLayer),
		Refs:            refs,
		StartTime:       skySpan.StartTime,
		EndTime:         skySpan.EndTime,
		ParentSpanID:    skySpan.ParentSpanId,
		OperationNameID: skySpan.OperationNameId,
		IsError:         skySpan.IsError,
		Tags:            tags,
		Logs:            newLogs,
	}
	return span
}

// TracCollector  全链路信息采集
func (sky *SkyWalking) TracCollector() {
	// 缓存
	var spanQueue []*util.Span
	// 定时器
	ticker := time.NewTicker(time.Duration(misc.Conf.SkyWalking.TraceReportInterval) * time.Millisecond)
	for {
		select {
		case reqPack, ok := <-sky.upTracChan:
			if ok {
				spansInfo := &protocol.TraceSegmentObject{}
				if err := proto.Unmarshal(reqPack.Segment, spansInfo); err != nil {
					g.L.Warn("Collect:proto.Unmarshal", zap.String("error", err.Error()))
					break
				}

				if !gAgent.isGetID {
					gAgent.appID = spansInfo.ApplicationId
					gAgent.appInstanceID = spansInfo.ApplicationInstanceId
				}

				// GlobalTraceIds 这里可能被MQ类似的服务批量处理，所以有多个tracID
				for _, tid := range reqPack.GlobalTraceIds {
					if len(tid.IdParts) != 3 {
						continue
					}
					tracID := fmt.Sprintf("%d.%d.%d", tid.IdParts[0], tid.IdParts[1], tid.IdParts[2])
					for _, skySpan := range spansInfo.Spans {
						newSpan := analysisSpan(spansInfo.ApplicationId, spansInfo.ApplicationInstanceId, tracID, skySpan)
						spanQueue = append(spanQueue, newSpan)
					}
				}
				if len(spanQueue) >= misc.Conf.SkyWalking.TraceCacheLen {
					// 发送
					if err := sky.sendSpans(spanQueue); err != nil {
						g.L.Warn("TracCollector:sky.sendSpans", zap.String("error", err.Error()))
					}
					// 清空缓存
					spanQueue = spanQueue[:0]
				}
			}
			break
		case <-ticker.C:
			if len(spanQueue) > 0 {
				// 发送
				if err := sky.sendSpans(spanQueue); err != nil {
					g.L.Warn("TracCollector:sky.sendSpans", zap.String("error", err.Error()))
				}
				// 清空缓存
				spanQueue = spanQueue[:0]
			}
			break
		}
	}
}

// sendSpans ...
func (sky *SkyWalking) sendSpans(spans []*util.Span) error {
	buf, err := msgpack.Marshal(spans)
	if err != nil {
		g.L.Warn("sendSpans:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	skp := util.SkywalkingPacket{
		Type:    util.TypeOfTraceSegment,
		Payload: buf,
	}

	payload, err := msgpack.Marshal(skp)
	if err != nil {
		g.L.Warn("sendSpans:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	packet := &util.VgoPacket{
		Type:       util.TypeOfSkywalking,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncNo,
		IsCompress: util.TypeOfCompressYes,
		Payload:    payload,
	}

	if err := gAgent.client.WritePacket(packet); err != nil {
		g.L.Warn("sendSpans:gAgent.client.WritePacket", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// sendJVMs ...
func (sky *SkyWalking) sendJVMs(jvms *util.JVMS) error {
	buf, err := msgpack.Marshal(jvms)
	if err != nil {
		g.L.Warn("sendJVMs:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	skp := util.SkywalkingPacket{
		Type:    util.TypeOfJVMMetrics,
		Payload: buf,
	}

	payload, err := msgpack.Marshal(skp)
	if err != nil {
		g.L.Warn("sendJVMs:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	packet := &util.VgoPacket{
		Type:       util.TypeOfSkywalking,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncNo,
		IsCompress: util.TypeOfCompressYes,
		Payload:    payload,
	}

	if err := gAgent.client.WritePacket(packet); err != nil {
		g.L.Warn("sendJVMs:gAgent.client.WritePacket", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// analysisJVM ...
func analysisJVM(old *protocol.JVMMetric) *util.JVM {
	cpu := &util.CPU{}
	memorys := []*util.Memory{}
	memoryPools := []*util.MemoryPool{}
	gcs := []*util.GC{}

	if old.Cpu != nil {
		cpu.UsagePercent = old.Cpu.UsagePercent
	}

	for _, value := range old.Memory {
		memorys = append(memorys, &util.Memory{
			IsHeap:    value.IsHeap,
			Init:      value.Init,
			Max:       value.Max,
			Used:      value.Used,
			Committed: value.Committed,
		})
	}

	for _, value := range old.MemoryPool {
		memoryPools = append(memoryPools, &util.MemoryPool{
			Type:     util.PoolType(value.Type),
			Init:     value.Init,
			Max:      value.Max,
			Used:     value.Used,
			Commited: value.Commited,
		})
	}

	for _, value := range old.Gc {
		gcs = append(gcs, &util.GC{
			Phrase: util.GCPhrase(value.Phrase),
			Count:  value.Count,
			Time:   value.Time,
		})
	}

	newJVM := &util.JVM{
		Time:       old.Time,
		CPU:        cpu,
		Memory:     memorys,
		MemoryPool: memoryPools,
		Gc:         gcs,
	}
	return newJVM
}

// TraceSegmentCollector trace 信息采集
func (sky *SkyWalking) TraceSegmentCollector() {
	ticker := time.NewTicker(time.Duration(misc.Conf.SkyWalking.TraceReportInterval) * time.Second)
	for {
		select {
		case jvmPack, ok := <-sky.jvmChan:
			if ok {
				log.Println(jvmPack)
			}
			break
		case <-ticker.C:
			log.Println("我是定时器 Trace")
			break
		}
	}
}

// Start ...
func (sky *SkyWalking) Start() error {

	// jvm上报
	go sky.JVMCollector()

	// trace信息上报
	go sky.TracCollector()
	// start http server
	sky.httpSer()

	// init grpc
	sky.grpcSer()

	return nil
}

func (sky *SkyWalking) grpcSer() error {
	lis, err := net.Listen("tcp", misc.Conf.SkyWalking.RPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	protocol.RegisterApplicationRegisterServiceServer(s, &appRegister{})
	protocol.RegisterInstanceDiscoveryServiceServer(s, &instanceDiscoveryService{})
	protocol.RegisterServiceNameDiscoveryServiceServer(s, &serviceNameDiscoveryService{})
	protocol.RegisterNetworkAddressRegisterServiceServer(s, &addressRegister{})

	protocol.RegisterJVMMetricsServiceServer(s, &jvmMetricsService{})
	protocol.RegisterTraceSegmentServiceServer(s, &traceSegmentService{})

	s.Serve(lis)
	return nil
}

func (sky *SkyWalking) httpSer() error {
	e := echo.New()
	e.GET("/agent/gRPC", rpcAddr)
	go e.Start(misc.Conf.SkyWalking.HTTPAddr)
	g.L.Info("Start:sky.httpSer", zap.String("httpSer", "ok"))
	return nil
}

// Close ...
func (sky *SkyWalking) Close() error {
	return nil
}

// rpcAddr ...
func rpcAddr(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{misc.Conf.SkyWalking.RPCAddr})
}

// appRegister ...
type appRegister struct {
}

// ApplicationCodeRegister ...
func (a *appRegister) ApplicationCodeRegister(ctx context.Context, al *protocol.Application) (*protocol.ApplicationMapping, error) {
	g.L.Info("ApplicationCodeRegister", zap.String("ApplicationCode", al.ApplicationCode), zap.String("AppName", gAgent.appName))
	// registerPacker
	registerPacker := &util.KeyWithStringValue{
		Key:   "appName",
		Value: gAgent.appName,
	}

	buf, err := msgpack.Marshal(registerPacker)
	if err != nil {
		g.L.Warn("ApplicationCodeRegister:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	skp := util.SkywalkingPacket{
		Type:    util.TypeOfAppRegister,
		Payload: buf,
	}

	payload, err := msgpack.Marshal(skp)
	if err != nil {
		g.L.Warn("ApplicationCodeRegister:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	// 获取ID
	id := gAgent.getSyncID()

	packet := &util.VgoPacket{
		Type:       util.TypeOfSkywalking,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncYes,
		IsCompress: util.TypeOfCompressNo,
		ID:         id,
		Payload:    payload,
	}

	if err := gAgent.client.WritePacket(packet); err != nil {
		g.L.Warn("ApplicationCodeRegister:gAgent.client.WritePacket", zap.String("error", err.Error()))
		return nil, err
	}

	// 创建chan
	if _, ok := gAgent.syncCall.newChan(id, 10); !ok {
		g.L.Warn("ApplicationCodeRegister:gAgent.syncCall.newChan", zap.String("error", "创建sync chan失败"))
		return nil, err
	}

	// 阻塞同步等待，并关闭chan
	respPakcet, err := gAgent.syncCall.syncRead(id, 10, true)
	if err != nil {
		g.L.Warn("ApplicationCodeRegister:gAgent.syncCall.syncRead", zap.String("error", err.Error()))
		return nil, err
	}

	// 非Skywalking返回错误
	if respPakcet.Type != util.TypeOfSkywalking {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("ApplicationCodeRegister:.", zap.Error(err))
		return nil, err
	}

	skypacket := &util.SkywalkingPacket{}
	if err := msgpack.Unmarshal(respPakcet.Payload, skypacket); err != nil {
		g.L.Warn("ApplicationCodeRegister:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	if skypacket.Type != util.TypeOfAppRegister {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("ApplicationCodeRegister:.", zap.Error(err))
		return nil, err
	}

	appIDPacket := &util.KeyWithIntegerValue{}
	if err := msgpack.Unmarshal(skypacket.Payload, appIDPacket); err != nil {
		g.L.Warn("ApplicationCodeRegister:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	// 保存AppCode
	gAgent.appID = appIDPacket.Value

	g.L.Info("ApplicationCodeRegister", zap.String("AppName", gAgent.appName), zap.Int32("appID", gAgent.appID))
	// log.Println("ApplicationCodeRegister 获取服务端packet", respRegister)

	kv := &protocol.KeyWithIntegerValue{
		Key:   gAgent.appName, // al.ApplicationCode,
		Value: gAgent.appID,
	}

	appMapping := &protocol.ApplicationMapping{
		Application: kv,
	}

	return appMapping, nil
}

// instanceDiscoveryService ...
type instanceDiscoveryService struct {
}

// RegisterInstance ...
func (i *instanceDiscoveryService) RegisterInstance(ctx context.Context, in *protocol.ApplicationInstance) (*protocol.ApplicationInstanceMapping, error) {
	ips, err := json.Marshal(&in.Osinfo.Ipv4S)
	if err != nil {
		g.L.Warn("RegisterInstance:json.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	agentInfo := &util.AgentInfo{
		AppID:        in.ApplicationId,
		AgentUUID:    in.AgentUUID,
		AppName:      gAgent.appName,
		OsName:       in.Osinfo.OsName,
		Ipv4S:        string(ips),
		RegisterTime: in.RegisterTime,
		ProcessID:    in.Osinfo.ProcessNo,
		HostName:     in.Osinfo.Hostname,
	}

	buf, err := msgpack.Marshal(agentInfo)
	if err != nil {
		g.L.Warn("RegisterInstance:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	skp := util.SkywalkingPacket{
		Type:    util.TypeOfAppRegisterInstance,
		Payload: buf,
	}

	payload, err := msgpack.Marshal(skp)
	if err != nil {
		g.L.Warn("RegisterInstance:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	// 获取ID
	id := gAgent.getSyncID()

	packet := &util.VgoPacket{
		Type:       util.TypeOfSkywalking,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncYes,
		IsCompress: util.TypeOfCompressNo,
		ID:         id,
		Payload:    payload,
	}

	if err := gAgent.client.WritePacket(packet); err != nil {
		g.L.Warn("RegisterInstance:gAgent.client.WritePacket", zap.String("error", err.Error()))
		return nil, err
	}

	// 创建chan
	if _, ok := gAgent.syncCall.newChan(id, 10); !ok {
		g.L.Warn("RegisterInstance:gAgent.syncCall.newChan", zap.String("error", "创建sync chan失败"))
		return nil, err
	}

	// 阻塞同步等待，并关闭chan
	respPakcet, err := gAgent.syncCall.syncRead(id, 10, true)
	if err != nil {
		g.L.Warn("RegisterInstance:gAgent.syncCall.syncRead", zap.String("error", err.Error()))
		return nil, err
	}

	// 非Skywalking返回错误
	if respPakcet.Type != util.TypeOfSkywalking {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("RegisterInstance:.", zap.Error(err))
		return nil, err
	}

	skypacket := &util.SkywalkingPacket{}
	if err := msgpack.Unmarshal(respPakcet.Payload, skypacket); err != nil {
		g.L.Warn("RegisterInstance:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	if skypacket.Type != util.TypeOfAppRegisterInstance {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("RegisterInstance:.", zap.Error(err))
		return nil, err
	}

	respPacket := &util.KeyWithIntegerValue{}
	if err := msgpack.Unmarshal(skypacket.Payload, respPacket); err != nil {
		g.L.Warn("RegisterInstance:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	gAgent.appInstanceID = respPacket.Value

	// 注册完成
	gAgent.isGetID = false

	am := &protocol.ApplicationInstanceMapping{
		ApplicationId:         gAgent.appID,
		ApplicationInstanceId: gAgent.appInstanceID,
	}

	return am, nil
}

// Heartbeat ...
func (i *instanceDiscoveryService) Heartbeat(ctx context.Context, in *protocol.ApplicationInstanceHeartbeat) (*protocol.Downstream, error) {
	log.Println("Heartbeat", in)
	dm := &protocol.Downstream{}
	return dm, nil
}

// serviceNameDiscoveryService ...
type serviceNameDiscoveryService struct {
}

// Discovery ...
func (s *serviceNameDiscoveryService) Discovery(ctx context.Context, in *protocol.ServiceNameCollection) (*protocol.ServiceNameMappingCollection, error) {
	if len(in.Elements) < 0 {
		return nil, fmt.Errorf("Elements is nil")
	}

	log.Println(in)

	reqPacket := &util.SerNameDiscoveryServices{
		SerNames: make([]*util.API, 0),
	}

	for _, value := range in.Elements {
		ser := &util.API{
			AppID:    value.ApplicationId,
			SerName:  value.ServiceName,
			SpanType: int32(value.SrcSpanType),
		}
		reqPacket.SerNames = append(reqPacket.SerNames, ser)
	}

	buf, err := msgpack.Marshal(reqPacket)
	if err != nil {
		g.L.Warn("Discovery:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	skp := util.SkywalkingPacket{
		Type:    util.TypeOfSerNameDiscoveryService,
		Payload: buf,
	}

	payload, err := msgpack.Marshal(skp)
	if err != nil {
		g.L.Warn("Discovery:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	// 获取ID
	id := gAgent.getSyncID()

	packet := &util.VgoPacket{
		Type:       util.TypeOfSkywalking,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncYes,
		IsCompress: util.TypeOfCompressNo,
		ID:         id,
		Payload:    payload,
	}

	if err := gAgent.client.WritePacket(packet); err != nil {
		g.L.Warn("Discovery:gAgent.client.WritePacket", zap.String("error", err.Error()))
		return nil, err
	}

	// 创建chan
	if _, ok := gAgent.syncCall.newChan(id, 10); !ok {
		g.L.Warn("Discovery:gAgent.syncCall.newChan", zap.String("error", "创建sync chan失败"))
		return nil, err
	}

	// 阻塞同步等待，并关闭chan
	respPakcet, err := gAgent.syncCall.syncRead(id, 10, true)
	if err != nil {
		g.L.Warn("Discovery:gAgent.syncCall.syncRead", zap.String("error", err.Error()))
		return nil, err
	}

	// 非Skywalking返回错误
	if respPakcet.Type != util.TypeOfSkywalking {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("Discovery:.", zap.Error(err))
		return nil, err
	}

	skypacket := &util.SkywalkingPacket{}
	if err := msgpack.Unmarshal(respPakcet.Payload, skypacket); err != nil {
		g.L.Warn("Discovery:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	if skypacket.Type != util.TypeOfSerNameDiscoveryService {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("Discovery:.", zap.Error(err))
		return nil, err
	}

	respRegister := &util.SerNameDiscoveryServices{}
	if err := msgpack.Unmarshal(skypacket.Payload, respRegister); err != nil {
		g.L.Warn("Discovery:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	// 组包 返回给sdk
	elements := make([]*protocol.ServiceNameMappingElement, 0)
	for _, value := range respRegister.SerNames {
		el := &protocol.ServiceNameElement{
			ServiceName:   value.SerName,
			ApplicationId: value.AppID,
			SrcSpanType:   protocol.SpanType(value.SpanType),
		}

		element := &protocol.ServiceNameMappingElement{
			ServiceId: value.SerID,
			Element:   el,
		}
		elements = append(elements, element)
	}
	mc := &protocol.ServiceNameMappingCollection{
		Elements: elements,
	}

	return mc, nil
}

// 地址注册服务
// addressRegisterService ...
type addressRegister struct {
}

// BatchRegister ...
func (a *addressRegister) BatchRegister(ctx context.Context, in *protocol.NetworkAddresses) (*protocol.NetworkAddressMappings, error) {
	return nil, nil

	log.Println("BatchRegister", in)

	if len(in.Addresses) < 0 {
		return nil, fmt.Errorf("addrs is nil")
	}

	reqPacket := &util.RegisterAddrs{
		// Skywalking设计问题，这里注册不带appid，重启就注册失败如果本地appid为0，
		AppID: gAgent.appID,
		Addrs: make([]*util.KeyWithIntegerValue, 0),
	}

	for _, value := range in.Addresses {
		addr := &util.KeyWithIntegerValue{
			Key: value,
		}
		reqPacket.Addrs = append(reqPacket.Addrs, addr)
	}

	buf, err := msgpack.Marshal(reqPacket)
	if err != nil {
		g.L.Warn("BatchRegister:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	skp := util.SkywalkingPacket{
		Type:    util.TypeOfNewworkAddrRegister,
		Payload: buf,
	}

	payload, err := msgpack.Marshal(skp)
	if err != nil {
		g.L.Warn("BatchRegister:msgpack.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	// 获取ID
	id := gAgent.getSyncID()

	packet := &util.VgoPacket{
		Type:       util.TypeOfSkywalking,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncYes,
		IsCompress: util.TypeOfCompressNo,
		ID:         id,
		Payload:    payload,
	}

	if err := gAgent.client.WritePacket(packet); err != nil {
		g.L.Warn("BatchRegister:gAgent.client.WritePacket", zap.String("error", err.Error()))
		return nil, err
	}

	// 创建chan
	if _, ok := gAgent.syncCall.newChan(id, 10); !ok {
		g.L.Warn("BatchRegister:gAgent.syncCall.newChan", zap.String("error", "创建sync chan失败"))
		return nil, err
	}

	// 阻塞同步等待，并关闭chan
	respPakcet, err := gAgent.syncCall.syncRead(id, 10, true)
	if err != nil {
		g.L.Warn("BatchRegister:gAgent.syncCall.syncRead", zap.String("error", err.Error()))
		return nil, err
	}

	// 非Skywalking返回错误
	if respPakcet.Type != util.TypeOfSkywalking {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("BatchRegister:.", zap.Error(err))
		return nil, err
	}

	skypacket := &util.SkywalkingPacket{}
	if err := msgpack.Unmarshal(respPakcet.Payload, skypacket); err != nil {
		g.L.Warn("BatchRegister:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	if skypacket.Type != util.TypeOfNewworkAddrRegister {
		err := fmt.Errorf("unknow type %d", respPakcet.Type)
		g.L.Warn("BatchRegister:.", zap.Error(err))
		return nil, err
	}

	respRegister := &util.RegisterAddrs{}
	if err := msgpack.Unmarshal(skypacket.Payload, respRegister); err != nil {
		g.L.Warn("BatchRegister:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	// protocol.KeyWithIntegerValue{}

	// 组包 返回给sdk
	mc := &protocol.NetworkAddressMappings{
		AddressIds: make([]*protocol.KeyWithIntegerValue, 0),
	}

	for _, value := range respRegister.Addrs {
		kv := &protocol.KeyWithIntegerValue{
			Key:   value.Key,
			Value: value.Value,
		}
		mc.AddressIds = append(mc.AddressIds, kv)
	}

	return nil, nil
}

// jvmMetricsService ...
type jvmMetricsService struct {
}

// Collect ...
func (j *jvmMetricsService) Collect(ctx context.Context, in *protocol.JVMMetrics) (*protocol.Downstream, error) {
	gAgent.skyWalk.jvmChan <- in
	return &protocol.Downstream{}, nil
}

// traceSegmentService ...
type traceSegmentService struct {
}

// Collect ...
func (c *traceSegmentService) Collect(in protocol.TraceSegmentService_CollectServer) error {

	u, err := in.Recv()
	if err != nil {
		log.Println(err)
	}
	// log.Println("u.GlobalTraceIds", u.GlobalTraceIds)
	// tarr := &protocol.TraceSegmentObject{}

	// if err := proto.Unmarshal(u.Segment, tarr); err != nil {
	// 	g.L.Warn("Collect:proto.Unmarshal", zap.String("error", err.Error()))
	// 	goto SEND_CLOSE
	// }

	gAgent.skyWalk.upTracChan <- u

	// g.L.Debug("Collect:u.Segment", zap.Any("u", u))
	// g.L.Debug("Collect:u.TraceSegmentObject", zap.Any("tarr", tarr))
	// for _, value := range tarr.Spans {
	// 	log.Println(tarr.TraceSegmentId, "value.SpanType", value.SpanType, "----------------------------->>>>", value, tarr)
	// }

	// SEND_CLOSE:
	if err := in.SendAndClose(&protocol.Downstream{}); err != nil {
		g.L.Warn("Collect:in.SendAndClose", zap.String("error", err.Error()))
		return err
	}

	return nil
}
