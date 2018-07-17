// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat.proto

package api

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
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	From                 string   `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	To                   string   `protobuf:"bytes,2,opt,name=To,proto3" json:"To,omitempty"`
	Body                 []byte   `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_chat_c96ecb2c79ebc921, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Message) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *Message) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "api.Message")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChatClient interface {
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (Chat_SubscribeClient, error)
}

type chatClient struct {
	cc *grpc.ClientConn
}

func NewChatClient(cc *grpc.ClientConn) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (Chat_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Chat_serviceDesc.Streams[0], "/api.Chat/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatSubscribeClient{stream}
	return x, nil
}

type Chat_SubscribeClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type chatSubscribeClient struct {
	grpc.ClientStream
}

func (x *chatSubscribeClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatSubscribeClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServer is the server API for Chat service.
type ChatServer interface {
	Subscribe(Chat_SubscribeServer) error
}

func RegisterChatServer(s *grpc.Server, srv ChatServer) {
	s.RegisterService(&_Chat_serviceDesc, srv)
}

func _Chat_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServer).Subscribe(&chatSubscribeServer{stream})
}

type Chat_SubscribeServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type chatSubscribeServer struct {
	grpc.ServerStream
}

func (x *chatSubscribeServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatSubscribeServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Chat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Chat_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "chat.proto",
}

func init() { proto.RegisterFile("chat.proto", fileDescriptor_chat_c96ecb2c79ebc921) }

var fileDescriptor_chat_c96ecb2c79ebc921 = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0xce, 0x48, 0x2c,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x54, 0x72, 0xe4, 0x62, 0xf7,
	0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x12, 0xe2, 0x62, 0x71, 0x2b, 0xca, 0xcf, 0x95, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x85, 0xf8, 0xb8, 0x98, 0x42, 0xf2, 0x25, 0x98, 0xc0,
	0x22, 0x4c, 0x21, 0xf9, 0x20, 0x35, 0x4e, 0xf9, 0x29, 0x95, 0x12, 0xcc, 0x0a, 0x8c, 0x1a, 0x3c,
	0x41, 0x60, 0xb6, 0x91, 0x31, 0x17, 0x8b, 0x73, 0x46, 0x62, 0x89, 0x90, 0x36, 0x17, 0x67, 0x70,
	0x69, 0x52, 0x71, 0x72, 0x51, 0x66, 0x52, 0xaa, 0x10, 0x8f, 0x5e, 0x62, 0x41, 0xa6, 0x1e, 0xd4,
	0x68, 0x29, 0x14, 0x9e, 0x06, 0xa3, 0x01, 0x63, 0x12, 0x1b, 0xd8, 0x0d, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x6b, 0x9f, 0xa8, 0x3b, 0x91, 0x00, 0x00, 0x00,
}