// Code generated by protoc-gen-go.
// source: kitfw.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	kitfw.proto

It has these top-level messages:
	KitfwRequest
	KitfwReply
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type KitfwRequest struct {
	Protoid int32  `protobuf:"varint,1,opt,name=protoid" json:"protoid,omitempty"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *KitfwRequest) Reset()                    { *m = KitfwRequest{} }
func (m *KitfwRequest) String() string            { return proto.CompactTextString(m) }
func (*KitfwRequest) ProtoMessage()               {}
func (*KitfwRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type KitfwReply struct {
	Protoid int32  `protobuf:"varint,1,opt,name=protoid" json:"protoid,omitempty"`
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *KitfwReply) Reset()                    { *m = KitfwReply{} }
func (m *KitfwReply) String() string            { return proto.CompactTextString(m) }
func (*KitfwReply) ProtoMessage()               {}
func (*KitfwReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*KitfwRequest)(nil), "pb.KitfwRequest")
	proto.RegisterType((*KitfwReply)(nil), "pb.KitfwReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Kitfw service

type KitfwClient interface {
	Process(ctx context.Context, in *KitfwRequest, opts ...grpc.CallOption) (*KitfwReply, error)
}

type kitfwClient struct {
	cc *grpc.ClientConn
}

func NewKitfwClient(cc *grpc.ClientConn) KitfwClient {
	return &kitfwClient{cc}
}

func (c *kitfwClient) Process(ctx context.Context, in *KitfwRequest, opts ...grpc.CallOption) (*KitfwReply, error) {
	out := new(KitfwReply)
	err := grpc.Invoke(ctx, "/pb.Kitfw/Process", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Kitfw service

type KitfwServer interface {
	Process(context.Context, *KitfwRequest) (*KitfwReply, error)
}

func RegisterKitfwServer(s *grpc.Server, srv KitfwServer) {
	s.RegisterService(&_Kitfw_serviceDesc, srv)
}

func _Kitfw_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KitfwRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KitfwServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Kitfw/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KitfwServer).Process(ctx, req.(*KitfwRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Kitfw_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Kitfw",
	HandlerType: (*KitfwServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _Kitfw_Process_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 137 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xce, 0x2c, 0x49,
	0x2b, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x72, 0xe2, 0xe2, 0xf1,
	0x06, 0x09, 0x05, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x49, 0x70, 0xb1, 0x83, 0x25, 0x33,
	0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0x60, 0x5c, 0xb0, 0x4c, 0x62, 0x65, 0x4e, 0x7e,
	0x62, 0x8a, 0x04, 0x33, 0x50, 0x86, 0x27, 0x08, 0xc6, 0x55, 0x72, 0xe0, 0xe2, 0x82, 0x9a, 0x51,
	0x90, 0x53, 0x49, 0x8e, 0x09, 0x46, 0x66, 0x5c, 0xac, 0x60, 0x13, 0x84, 0x74, 0xb9, 0xd8, 0x03,
	0x8a, 0xf2, 0x93, 0x53, 0x8b, 0x8b, 0x85, 0x04, 0xf4, 0x0a, 0x92, 0xf4, 0x90, 0xdd, 0x26, 0xc5,
	0x87, 0x24, 0x02, 0xb4, 0x49, 0x89, 0x21, 0x89, 0x0d, 0x6c, 0xb4, 0x31, 0x20, 0x00, 0x00, 0xff,
	0xff, 0x51, 0xf5, 0x6f, 0xac, 0xd7, 0x00, 0x00, 0x00,
}