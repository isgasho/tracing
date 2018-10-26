package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

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
}

// NewSkyWalking ...
func NewSkyWalking() *SkyWalking {

	return &SkyWalking{
		appRegisterC: make(chan *appRegister, 10),
	}
}

// Start ...
func (sky *SkyWalking) Start() error {

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

	// protocol.RegisterNetworkAddressRegisterServiceServer(s, &NetworkAddressRegisterService{})
	// protocol.RegisterJVMMetricsServiceServer(s, &JVMMetricsService{})
	// protocol.RegisterTraceSegmentServiceServer(s, &TraceSegmentService{})

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
