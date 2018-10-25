package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/mafanr/vgo/util"
	"github.com/vmihailenco/msgpack"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/agent/misc"
	"github.com/mafanr/vgo/agent/protocol"
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
	// protocol.RegisterServiceNameDiscoveryServiceServer(s, &ServiceNameDiscoveryService{})
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

	g.L.Info("ApplicationCodeRegister", zap.String("ApplicationCode", al.ApplicationCode))
	// AppRegister
	registerPacker := &util.AppRegister{
		Name: gAgent.appName,
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

	respRegister := &util.AppRegister{}
	if err := msgpack.Unmarshal(skypacket.Payload, respRegister); err != nil {
		g.L.Warn("ApplicationCodeRegister:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	// 保存AppCode
	gAgent.appCode = respRegister.Code
	gAgent.agentInfo.AppCode = gAgent.appCode

	log.Println("ApplicationCodeRegister 获取服务端packet", respRegister)

	kv := &protocol.KeyWithIntegerValue{
		Key:   gAgent.appName, // al.ApplicationCode,
		Value: respRegister.Code,
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
	// agent重启需要保存app code
	if gAgent.appCode == -1 {
		gAgent.appCode = in.ApplicationId
	}

	gAgent.agentInfo.AppCode = gAgent.appCode
	gAgent.agentInfo.AppName = gAgent.appName
	gAgent.agentInfo.AgentUUID = in.AgentUUID
	gAgent.agentInfo.RegisterTime = in.RegisterTime
	gAgent.agentInfo.OsName = in.Osinfo.OsName
	gAgent.agentInfo.HostName = in.Osinfo.Hostname
	gAgent.agentInfo.ProcessID = in.Osinfo.ProcessNo

	ips, err := json.Marshal(&in.Osinfo.Ipv4S)
	if err != nil {
		g.L.Warn("RegisterInstance:json.Marshal", zap.String("error", err.Error()))
		return nil, err
	}

	gAgent.agentInfo.Ipv4S = string(ips)

	buf, err := msgpack.Marshal(gAgent.agentInfo)
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

	respRegister := &util.KeyWithIntegerValue{}
	if err := msgpack.Unmarshal(skypacket.Payload, respRegister); err != nil {
		g.L.Warn("RegisterInstance:msgpack.Unmarshal", zap.String("error", err.Error()))
		return nil, err
	}

	log.Println("RegisterInstance 获取服务端packet", respRegister)

	gAgent.agentInfo.ID = respRegister.Value

	am := &protocol.ApplicationInstanceMapping{
		ApplicationId:         gAgent.appCode,
		ApplicationInstanceId: gAgent.agentInfo.ID,
	}
	log.Println(am)

	log.Println("appCode--->>>>>", gAgent.appCode)
	log.Println("ID--->>>>>", gAgent.agentInfo.ID)
	log.Println("AgentUUID--->>>>>", gAgent.agentInfo.AgentUUID)
	log.Println("AppCode--->>>>>", gAgent.agentInfo.AppCode)
	log.Println("OsName--->>>>>", gAgent.agentInfo.OsName)
	log.Println("Ipv4S--->>>>>", gAgent.agentInfo.Ipv4S)
	log.Println("RegisterTime--->>>>>", gAgent.agentInfo.RegisterTime)
	log.Println("ProcessID--->>>>>", gAgent.agentInfo.ProcessID)
	log.Println("HostName--->>>>>", gAgent.agentInfo.HostName)
	return nil, nil
}

// Heartbeat ...
func (i *instanceDiscoveryService) Heartbeat(ctx context.Context, in *protocol.ApplicationInstanceHeartbeat) (*protocol.Downstream, error) {
	log.Println("Heartbeat", in)
	dm := &protocol.Downstream{}
	return dm, nil
}
