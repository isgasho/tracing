// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ApplicationRegisterService.proto

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

type Application struct {
	ApplicationCode      string   `protobuf:"bytes,1,opt,name=applicationCode,proto3" json:"applicationCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Application) Reset()         { *m = Application{} }
func (m *Application) String() string { return proto.CompactTextString(m) }
func (*Application) ProtoMessage()    {}
func (*Application) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a88ee0d7366d2ac, []int{0}
}

func (m *Application) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Application.Unmarshal(m, b)
}
func (m *Application) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Application.Marshal(b, m, deterministic)
}
func (m *Application) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Application.Merge(m, src)
}
func (m *Application) XXX_Size() int {
	return xxx_messageInfo_Application.Size(m)
}
func (m *Application) XXX_DiscardUnknown() {
	xxx_messageInfo_Application.DiscardUnknown(m)
}

var xxx_messageInfo_Application proto.InternalMessageInfo

func (m *Application) GetApplicationCode() string {
	if m != nil {
		return m.ApplicationCode
	}
	return ""
}

type ApplicationMapping struct {
	Application          *KeyWithIntegerValue `protobuf:"bytes,1,opt,name=application,proto3" json:"application,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ApplicationMapping) Reset()         { *m = ApplicationMapping{} }
func (m *ApplicationMapping) String() string { return proto.CompactTextString(m) }
func (*ApplicationMapping) ProtoMessage()    {}
func (*ApplicationMapping) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a88ee0d7366d2ac, []int{1}
}

func (m *ApplicationMapping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApplicationMapping.Unmarshal(m, b)
}
func (m *ApplicationMapping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApplicationMapping.Marshal(b, m, deterministic)
}
func (m *ApplicationMapping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApplicationMapping.Merge(m, src)
}
func (m *ApplicationMapping) XXX_Size() int {
	return xxx_messageInfo_ApplicationMapping.Size(m)
}
func (m *ApplicationMapping) XXX_DiscardUnknown() {
	xxx_messageInfo_ApplicationMapping.DiscardUnknown(m)
}

var xxx_messageInfo_ApplicationMapping proto.InternalMessageInfo

func (m *ApplicationMapping) GetApplication() *KeyWithIntegerValue {
	if m != nil {
		return m.Application
	}
	return nil
}

func init() {
	proto.RegisterType((*Application)(nil), "Application")
	proto.RegisterType((*ApplicationMapping)(nil), "ApplicationMapping")
}

func init() { proto.RegisterFile("ApplicationRegisterService.proto", fileDescriptor_9a88ee0d7366d2ac) }

var fileDescriptor_9a88ee0d7366d2ac = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x70, 0x2c, 0x28, 0xc8,
	0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0x0b, 0x4a, 0x4d, 0xcf, 0x2c, 0x2e, 0x49, 0x2d, 0x0a,
	0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0x92, 0xf4, 0x4e,
	0xad, 0x0c, 0xcf, 0x2c, 0xc9, 0xf0, 0xcc, 0x2b, 0x49, 0x4d, 0x4f, 0x2d, 0x0a, 0x4b, 0xcc, 0x29,
	0x85, 0x4a, 0x29, 0x99, 0x73, 0x71, 0x23, 0x69, 0x17, 0xd2, 0xe0, 0xe2, 0x4f, 0x44, 0x70, 0x9d,
	0xf3, 0x53, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xd0, 0x85, 0x95, 0x7c, 0xb8, 0x84,
	0x90, 0x34, 0xfa, 0x26, 0x16, 0x14, 0x64, 0xe6, 0xa5, 0x0b, 0x99, 0x71, 0x71, 0x23, 0x29, 0x04,
	0xeb, 0xe5, 0x36, 0x12, 0xd1, 0xc3, 0x62, 0x7f, 0x10, 0xb2, 0x42, 0xa3, 0x18, 0x2e, 0x29, 0xdc,
	0xbe, 0x10, 0xb2, 0xe3, 0x12, 0x47, 0xb3, 0x1e, 0xa6, 0x42, 0x88, 0x47, 0x0f, 0x49, 0x9f, 0x94,
	0xb0, 0x1e, 0xa6, 0x9b, 0x94, 0x18, 0x9c, 0x3c, 0xb8, 0xd4, 0xf3, 0x8b, 0xd2, 0xf5, 0x12, 0x0b,
	0x12, 0x93, 0x33, 0x52, 0xf5, 0x8a, 0xb3, 0x2b, 0xcb, 0x13, 0x73, 0xb2, 0x33, 0xf3, 0x40, 0x22,
	0xb9, 0x7a, 0x79, 0xa9, 0x25, 0xe5, 0xf9, 0x45, 0xd9, 0x90, 0xf0, 0x08, 0x60, 0x5c, 0xc5, 0x24,
	0x15, 0x9c, 0x5d, 0x19, 0x0e, 0x55, 0xe0, 0x07, 0x91, 0x0c, 0x00, 0xc9, 0x25, 0xe7, 0xe7, 0x24,
	0xb1, 0x81, 0x55, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x89, 0x96, 0x6a, 0xb4, 0x74, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ApplicationRegisterServiceClient is the client API for ApplicationRegisterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApplicationRegisterServiceClient interface {
	ApplicationCodeRegister(ctx context.Context, in *Application, opts ...grpc.CallOption) (*ApplicationMapping, error)
}

type applicationRegisterServiceClient struct {
	cc *grpc.ClientConn
}

func NewApplicationRegisterServiceClient(cc *grpc.ClientConn) ApplicationRegisterServiceClient {
	return &applicationRegisterServiceClient{cc}
}

func (c *applicationRegisterServiceClient) ApplicationCodeRegister(ctx context.Context, in *Application, opts ...grpc.CallOption) (*ApplicationMapping, error) {
	out := new(ApplicationMapping)
	err := c.cc.Invoke(ctx, "/ApplicationRegisterService/applicationCodeRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApplicationRegisterServiceServer is the server API for ApplicationRegisterService service.
type ApplicationRegisterServiceServer interface {
	ApplicationCodeRegister(context.Context, *Application) (*ApplicationMapping, error)
}

func RegisterApplicationRegisterServiceServer(s *grpc.Server, srv ApplicationRegisterServiceServer) {
	s.RegisterService(&_ApplicationRegisterService_serviceDesc, srv)
}

func _ApplicationRegisterService_ApplicationCodeRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Application)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApplicationRegisterServiceServer).ApplicationCodeRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ApplicationRegisterService/ApplicationCodeRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApplicationRegisterServiceServer).ApplicationCodeRegister(ctx, req.(*Application))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApplicationRegisterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ApplicationRegisterService",
	HandlerType: (*ApplicationRegisterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "applicationCodeRegister",
			Handler:    _ApplicationRegisterService_ApplicationCodeRegister_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ApplicationRegisterService.proto",
}
