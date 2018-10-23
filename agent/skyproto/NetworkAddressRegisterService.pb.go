// Code generated by protoc-gen-go. DO NOT EDIT.
// source: NetworkAddressRegisterService.proto

package skyproto

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NetworkAddresses struct {
	Addresses            []string `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NetworkAddresses) Reset()         { *m = NetworkAddresses{} }
func (m *NetworkAddresses) String() string { return proto.CompactTextString(m) }
func (*NetworkAddresses) ProtoMessage()    {}
func (*NetworkAddresses) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc5fb5cfec0cc919, []int{0}
}

func (m *NetworkAddresses) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkAddresses.Unmarshal(m, b)
}
func (m *NetworkAddresses) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkAddresses.Marshal(b, m, deterministic)
}
func (m *NetworkAddresses) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkAddresses.Merge(m, src)
}
func (m *NetworkAddresses) XXX_Size() int {
	return xxx_messageInfo_NetworkAddresses.Size(m)
}
func (m *NetworkAddresses) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkAddresses.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkAddresses proto.InternalMessageInfo

func (m *NetworkAddresses) GetAddresses() []string {
	if m != nil {
		return m.Addresses
	}
	return nil
}

type NetworkAddressMappings struct {
	AddressIds           []*KeyWithIntegerValue `protobuf:"bytes,1,rep,name=addressIds,proto3" json:"addressIds,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *NetworkAddressMappings) Reset()         { *m = NetworkAddressMappings{} }
func (m *NetworkAddressMappings) String() string { return proto.CompactTextString(m) }
func (*NetworkAddressMappings) ProtoMessage()    {}
func (*NetworkAddressMappings) Descriptor() ([]byte, []int) {
	return fileDescriptor_fc5fb5cfec0cc919, []int{1}
}

func (m *NetworkAddressMappings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NetworkAddressMappings.Unmarshal(m, b)
}
func (m *NetworkAddressMappings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NetworkAddressMappings.Marshal(b, m, deterministic)
}
func (m *NetworkAddressMappings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkAddressMappings.Merge(m, src)
}
func (m *NetworkAddressMappings) XXX_Size() int {
	return xxx_messageInfo_NetworkAddressMappings.Size(m)
}
func (m *NetworkAddressMappings) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkAddressMappings.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkAddressMappings proto.InternalMessageInfo

func (m *NetworkAddressMappings) GetAddressIds() []*KeyWithIntegerValue {
	if m != nil {
		return m.AddressIds
	}
	return nil
}

func init() {
	proto.RegisterType((*NetworkAddresses)(nil), "NetworkAddresses")
	proto.RegisterType((*NetworkAddressMappings)(nil), "NetworkAddressMappings")
}

func init() {
	proto.RegisterFile("NetworkAddressRegisterService.proto", fileDescriptor_fc5fb5cfec0cc919)
}

var fileDescriptor_fc5fb5cfec0cc919 = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xf6, 0x4b, 0x2d, 0x29,
	0xcf, 0x2f, 0xca, 0x76, 0x4c, 0x49, 0x29, 0x4a, 0x2d, 0x2e, 0x0e, 0x4a, 0x4d, 0xcf, 0x2c, 0x2e,
	0x49, 0x2d, 0x0a, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97,
	0x92, 0xf4, 0x4e, 0xad, 0x0c, 0xcf, 0x2c, 0xc9, 0xf0, 0xcc, 0x2b, 0x49, 0x4d, 0x4f, 0x2d, 0x0a,
	0x4b, 0xcc, 0x29, 0x85, 0x4a, 0x29, 0x19, 0x70, 0x09, 0xa0, 0x9a, 0x90, 0x5a, 0x2c, 0x24, 0xc3,
	0xc5, 0x99, 0x08, 0xe3, 0x48, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x06, 0x21, 0x04, 0x94, 0xfc, 0xb8,
	0xc4, 0x50, 0x75, 0xf8, 0x26, 0x16, 0x14, 0x64, 0xe6, 0xa5, 0x17, 0x0b, 0x99, 0x70, 0x71, 0x41,
	0x95, 0x79, 0xa6, 0x40, 0x34, 0x72, 0x1b, 0x89, 0xe8, 0x61, 0xb1, 0x3b, 0x08, 0x49, 0x9d, 0x51,
	0x1c, 0x97, 0x2c, 0x5e, 0x3f, 0x08, 0xd9, 0x72, 0xf1, 0x26, 0x25, 0x96, 0x24, 0x67, 0xc0, 0xc4,
	0x85, 0x04, 0xf5, 0xd0, 0x9d, 0x2c, 0x25, 0xae, 0x87, 0xdd, 0x4d, 0x4a, 0x0c, 0x4e, 0x1e, 0x5c,
	0xea, 0xf9, 0x45, 0xe9, 0x7a, 0x89, 0x05, 0x89, 0xc9, 0x19, 0xa9, 0x7a, 0xc5, 0xd9, 0x95, 0xe5,
	0x89, 0x39, 0xd9, 0x99, 0x79, 0x20, 0x91, 0x5c, 0xbd, 0x3c, 0x88, 0x2e, 0x48, 0x60, 0x04, 0x30,
	0xae, 0x62, 0x92, 0x0a, 0xce, 0xae, 0x0c, 0x87, 0x2a, 0x80, 0x1a, 0x19, 0x00, 0x92, 0x4b, 0xce,
	0xcf, 0x49, 0x62, 0x03, 0xab, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x27, 0xd6, 0x46, 0x99,
	0x74, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NetworkAddressRegisterServiceClient is the client API for NetworkAddressRegisterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NetworkAddressRegisterServiceClient interface {
	BatchRegister(ctx context.Context, in *NetworkAddresses, opts ...grpc.CallOption) (*NetworkAddressMappings, error)
}

type networkAddressRegisterServiceClient struct {
	cc *grpc.ClientConn
}

func NewNetworkAddressRegisterServiceClient(cc *grpc.ClientConn) NetworkAddressRegisterServiceClient {
	return &networkAddressRegisterServiceClient{cc}
}

func (c *networkAddressRegisterServiceClient) BatchRegister(ctx context.Context, in *NetworkAddresses, opts ...grpc.CallOption) (*NetworkAddressMappings, error) {
	out := new(NetworkAddressMappings)
	err := c.cc.Invoke(ctx, "/NetworkAddressRegisterService/batchRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkAddressRegisterServiceServer is the server API for NetworkAddressRegisterService service.
type NetworkAddressRegisterServiceServer interface {
	BatchRegister(context.Context, *NetworkAddresses) (*NetworkAddressMappings, error)
}

func RegisterNetworkAddressRegisterServiceServer(s *grpc.Server, srv NetworkAddressRegisterServiceServer) {
	s.RegisterService(&_NetworkAddressRegisterService_serviceDesc, srv)
}

func _NetworkAddressRegisterService_BatchRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkAddresses)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkAddressRegisterServiceServer).BatchRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NetworkAddressRegisterService/BatchRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkAddressRegisterServiceServer).BatchRegister(ctx, req.(*NetworkAddresses))
	}
	return interceptor(ctx, in, info, handler)
}

var _NetworkAddressRegisterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "NetworkAddressRegisterService",
	HandlerType: (*NetworkAddressRegisterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "batchRegister",
			Handler:    _NetworkAddressRegisterService_BatchRegister_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "NetworkAddressRegisterService.proto",
}
