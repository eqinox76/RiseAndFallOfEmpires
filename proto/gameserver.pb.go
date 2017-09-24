// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/gameserver.proto

package riseandfall

import proto "github.com/gogo/protobuf/proto"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GameServer service

type GameServerClient interface {
	// get all game states streamed
	Observe(ctx context.Context, in *ID, opts ...grpc.CallOption) (GameServer_ObserveClient, error)
	CurrentGameState(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Space, error)
	// control one or multiple clients and get the current game state
	StrategyClient(ctx context.Context, opts ...grpc.CallOption) (GameServer_StrategyClientClient, error)
}

type gameServerClient struct {
	cc *grpc.ClientConn
}

func NewGameServerClient(cc *grpc.ClientConn) GameServerClient {
	return &gameServerClient{cc}
}

func (c *gameServerClient) Observe(ctx context.Context, in *ID, opts ...grpc.CallOption) (GameServer_ObserveClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GameServer_serviceDesc.Streams[0], c.cc, "/riseandfall.GameServer/Observe", opts...)
	if err != nil {
		return nil, err
	}
	x := &gameServerObserveClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GameServer_ObserveClient interface {
	Recv() (*Space, error)
	grpc.ClientStream
}

type gameServerObserveClient struct {
	grpc.ClientStream
}

func (x *gameServerObserveClient) Recv() (*Space, error) {
	m := new(Space)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gameServerClient) CurrentGameState(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Space, error) {
	out := new(Space)
	err := grpc.Invoke(ctx, "/riseandfall.GameServer/CurrentGameState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameServerClient) StrategyClient(ctx context.Context, opts ...grpc.CallOption) (GameServer_StrategyClientClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GameServer_serviceDesc.Streams[1], c.cc, "/riseandfall.GameServer/StrategyClient", opts...)
	if err != nil {
		return nil, err
	}
	x := &gameServerStrategyClientClient{stream}
	return x, nil
}

type GameServer_StrategyClientClient interface {
	Send(*Command) error
	Recv() (*Space, error)
	grpc.ClientStream
}

type gameServerStrategyClientClient struct {
	grpc.ClientStream
}

func (x *gameServerStrategyClientClient) Send(m *Command) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gameServerStrategyClientClient) Recv() (*Space, error) {
	m := new(Space)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for GameServer service

type GameServerServer interface {
	// get all game states streamed
	Observe(*ID, GameServer_ObserveServer) error
	CurrentGameState(context.Context, *ID) (*Space, error)
	// control one or multiple clients and get the current game state
	StrategyClient(GameServer_StrategyClientServer) error
}

func RegisterGameServerServer(s *grpc.Server, srv GameServerServer) {
	s.RegisterService(&_GameServer_serviceDesc, srv)
}

func _GameServer_Observe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ID)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GameServerServer).Observe(m, &gameServerObserveServer{stream})
}

type GameServer_ObserveServer interface {
	Send(*Space) error
	grpc.ServerStream
}

type gameServerObserveServer struct {
	grpc.ServerStream
}

func (x *gameServerObserveServer) Send(m *Space) error {
	return x.ServerStream.SendMsg(m)
}

func _GameServer_CurrentGameState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServerServer).CurrentGameState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/riseandfall.GameServer/CurrentGameState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServerServer).CurrentGameState(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameServer_StrategyClient_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GameServerServer).StrategyClient(&gameServerStrategyClientServer{stream})
}

type GameServer_StrategyClientServer interface {
	Send(*Space) error
	Recv() (*Command, error)
	grpc.ServerStream
}

type gameServerStrategyClientServer struct {
	grpc.ServerStream
}

func (x *gameServerStrategyClientServer) Send(m *Space) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gameServerStrategyClientServer) Recv() (*Command, error) {
	m := new(Command)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _GameServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "riseandfall.GameServer",
	HandlerType: (*GameServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CurrentGameState",
			Handler:    _GameServer_CurrentGameState_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Observe",
			Handler:       _GameServer_Observe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StrategyClient",
			Handler:       _GameServer_StrategyClient_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/gameserver.proto",
}

func init() { proto.RegisterFile("proto/gameserver.proto", fileDescriptorGameserver) }

var fileDescriptorGameserver = []byte{
	// 225 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x4f, 0xcc, 0x4d, 0x2d, 0x4e, 0x2d, 0x2a, 0x4b, 0x2d, 0xd2, 0x03, 0x0b, 0x08,
	0x71, 0x17, 0x65, 0x16, 0xa7, 0x26, 0xe6, 0xa5, 0xa4, 0x25, 0xe6, 0xe4, 0x48, 0x09, 0x42, 0x14,
	0x15, 0x17, 0x24, 0x26, 0xa7, 0x42, 0xe4, 0xa5, 0x84, 0x21, 0x42, 0xc9, 0xf9, 0xb9, 0xb9, 0x89,
	0x79, 0x29, 0x10, 0x41, 0xa3, 0xbd, 0x8c, 0x5c, 0x5c, 0xee, 0x89, 0xb9, 0xa9, 0xc1, 0x60, 0x93,
	0x84, 0x8c, 0xb8, 0xd8, 0xfd, 0x93, 0xc0, 0xa6, 0x0a, 0xf1, 0xeb, 0x21, 0x99, 0xa7, 0xe7, 0xe9,
	0x22, 0x25, 0x84, 0x22, 0x10, 0x0c, 0x32, 0x59, 0x89, 0xc1, 0x80, 0x51, 0xc8, 0x92, 0x4b, 0xc0,
	0xb9, 0xb4, 0xa8, 0x28, 0x35, 0xaf, 0x04, 0x6c, 0x50, 0x49, 0x62, 0x09, 0xb1, 0x9a, 0x85, 0x1c,
	0xb8, 0xf8, 0x82, 0x4b, 0x8a, 0x12, 0x4b, 0x52, 0xd3, 0x2b, 0x9d, 0x73, 0x32, 0x53, 0xf3, 0x4a,
	0x84, 0x44, 0x50, 0xd4, 0x39, 0x43, 0xdc, 0x8a, 0x5d, 0xb7, 0x06, 0xa3, 0x01, 0xa3, 0x93, 0xce,
	0x85, 0x87, 0x72, 0x0c, 0x37, 0x1e, 0xca, 0x31, 0x7c, 0x78, 0x28, 0xc7, 0xd8, 0xf0, 0x48, 0x8e,
	0x71, 0xc5, 0x23, 0x39, 0xc6, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48,
	0x8e, 0xf1, 0xc5, 0x23, 0x39, 0x86, 0x0f, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0x48, 0x62,
	0x03, 0x7b, 0xda, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xf9, 0xa6, 0x5d, 0x6c, 0x43, 0x01, 0x00,
	0x00,
}