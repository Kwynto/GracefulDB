// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: config/proto/graceful_service.proto

package graceful_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GracefulService_VQuery_FullMethodName = "/graceful_service.GracefulService/VQuery"
	GracefulService_SQuery_FullMethodName = "/graceful_service.GracefulService/SQuery"
)

// GracefulServiceClient is the client API for GracefulService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GracefulServiceClient interface {
	VQuery(ctx context.Context, in *VRequest, opts ...grpc.CallOption) (*VResponse, error)
	SQuery(ctx context.Context, in *SRequest, opts ...grpc.CallOption) (*SResponse, error)
}

type gracefulServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGracefulServiceClient(cc grpc.ClientConnInterface) GracefulServiceClient {
	return &gracefulServiceClient{cc}
}

func (c *gracefulServiceClient) VQuery(ctx context.Context, in *VRequest, opts ...grpc.CallOption) (*VResponse, error) {
	out := new(VResponse)
	err := c.cc.Invoke(ctx, GracefulService_VQuery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gracefulServiceClient) SQuery(ctx context.Context, in *SRequest, opts ...grpc.CallOption) (*SResponse, error) {
	out := new(SResponse)
	err := c.cc.Invoke(ctx, GracefulService_SQuery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GracefulServiceServer is the server API for GracefulService service.
// All implementations must embed UnimplementedGracefulServiceServer
// for forward compatibility
type GracefulServiceServer interface {
	VQuery(context.Context, *VRequest) (*VResponse, error)
	SQuery(context.Context, *SRequest) (*SResponse, error)
	mustEmbedUnimplementedGracefulServiceServer()
}

// UnimplementedGracefulServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGracefulServiceServer struct {
}

func (UnimplementedGracefulServiceServer) VQuery(context.Context, *VRequest) (*VResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VQuery not implemented")
}
func (UnimplementedGracefulServiceServer) SQuery(context.Context, *SRequest) (*SResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SQuery not implemented")
}
func (UnimplementedGracefulServiceServer) mustEmbedUnimplementedGracefulServiceServer() {}

// UnsafeGracefulServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GracefulServiceServer will
// result in compilation errors.
type UnsafeGracefulServiceServer interface {
	mustEmbedUnimplementedGracefulServiceServer()
}

func RegisterGracefulServiceServer(s grpc.ServiceRegistrar, srv GracefulServiceServer) {
	s.RegisterService(&GracefulService_ServiceDesc, srv)
}

func _GracefulService_VQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GracefulServiceServer).VQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GracefulService_VQuery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GracefulServiceServer).VQuery(ctx, req.(*VRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GracefulService_SQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GracefulServiceServer).SQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GracefulService_SQuery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GracefulServiceServer).SQuery(ctx, req.(*SRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GracefulService_ServiceDesc is the grpc.ServiceDesc for GracefulService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GracefulService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "graceful_service.GracefulService",
	HandlerType: (*GracefulServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VQuery",
			Handler:    _GracefulService_VQuery_Handler,
		},
		{
			MethodName: "SQuery",
			Handler:    _GracefulService_SQuery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "config/proto/graceful_service.proto",
}