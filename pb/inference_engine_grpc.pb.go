// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: inference_engine.proto

package pb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	GRPCInferenceService_ModelStreamInfer_FullMethodName  = "/language_inference.GRPCInferenceService/ModelStreamInfer"
	GRPCInferenceService_ModelFetchRequest_FullMethodName = "/language_inference.GRPCInferenceService/ModelFetchRequest"
	GRPCInferenceService_ModelSendResponse_FullMethodName = "/language_inference.GRPCInferenceService/ModelSendResponse"
)

// GRPCInferenceServiceClient is the client API for GRPCInferenceService gateway.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// protoc --proto_path=./pb --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative inference_engine.proto
// Inference Server GRPC endpoints.
type GRPCInferenceServiceClient interface {
	// 模型推理请求入口
	// 输入一个请求，流式返回多个response
	ModelStreamInfer(ctx context.Context, in *ModelInferRequest, opts ...grpc.CallOption) (GRPCInferenceService_ModelStreamInferClient, error)
	// 拉取一个请求，给inference server调用
	ModelFetchRequest(ctx context.Context, in *ModelFetchRequestParams, opts ...grpc.CallOption) (*ModelFetchRequestResult, error)
	// 发送请求的返回结果，给inference server调用
	// response是流式的发送
	ModelSendResponse(ctx context.Context, opts ...grpc.CallOption) (GRPCInferenceService_ModelSendResponseClient, error)
}

type gRPCInferenceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGRPCInferenceServiceClient(cc grpc.ClientConnInterface) GRPCInferenceServiceClient {
	return &gRPCInferenceServiceClient{cc}
}

func (c *gRPCInferenceServiceClient) ModelStreamInfer(ctx context.Context, in *ModelInferRequest, opts ...grpc.CallOption) (GRPCInferenceService_ModelStreamInferClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &GRPCInferenceService_ServiceDesc.Streams[0],
		GRPCInferenceService_ModelStreamInfer_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCInferenceServiceModelStreamInferClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GRPCInferenceService_ModelStreamInferClient interface {
	Recv() (*ModelInferResponse, error)
	grpc.ClientStream
}

type gRPCInferenceServiceModelStreamInferClient struct {
	grpc.ClientStream
}

func (x *gRPCInferenceServiceModelStreamInferClient) Recv() (*ModelInferResponse, error) {
	m := new(ModelInferResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gRPCInferenceServiceClient) ModelFetchRequest(ctx context.Context, in *ModelFetchRequestParams, opts ...grpc.CallOption) (*ModelFetchRequestResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ModelFetchRequestResult)
	err := c.cc.Invoke(ctx, GRPCInferenceService_ModelFetchRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gRPCInferenceServiceClient) ModelSendResponse(ctx context.Context, opts ...grpc.CallOption) (GRPCInferenceService_ModelSendResponseClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &GRPCInferenceService_ServiceDesc.Streams[1],
		GRPCInferenceService_ModelSendResponse_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &gRPCInferenceServiceModelSendResponseClient{ClientStream: stream}
	return x, nil
}

type GRPCInferenceService_ModelSendResponseClient interface {
	Send(*ModelInferResponse) error
	CloseAndRecv() (*ModelSendResponseResult, error)
	grpc.ClientStream
}

type gRPCInferenceServiceModelSendResponseClient struct {
	grpc.ClientStream
}

func (x *gRPCInferenceServiceModelSendResponseClient) Send(m *ModelInferResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gRPCInferenceServiceModelSendResponseClient) CloseAndRecv() (*ModelSendResponseResult, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ModelSendResponseResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GRPCInferenceServiceServer is the server API for GRPCInferenceService gateway.
// All implementations must embed UnimplementedGRPCInferenceServiceServer
// for forward compatibility
//
// protoc --proto_path=./pb --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative inference_engine.proto
// Inference Server GRPC endpoints.
type GRPCInferenceServiceServer interface {
	// 模型推理请求入口
	// 输入一个请求，流式返回多个response
	ModelStreamInfer(*ModelInferRequest, GRPCInferenceService_ModelStreamInferServer) error
	// 拉取一个请求，给inference server调用
	ModelFetchRequest(context.Context, *ModelFetchRequestParams) (*ModelFetchRequestResult, error)
	// 发送请求的返回结果，给inference server调用
	// response是流式的发送
	ModelSendResponse(GRPCInferenceService_ModelSendResponseServer) error
	mustEmbedUnimplementedGRPCInferenceServiceServer()
}

// UnimplementedGRPCInferenceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGRPCInferenceServiceServer struct {
}

func (UnimplementedGRPCInferenceServiceServer) ModelStreamInfer(*ModelInferRequest,
	GRPCInferenceService_ModelStreamInferServer) error {
	return status.Errorf(codes.Unimplemented, "method ModelStreamInfer not implemented")
}
func (UnimplementedGRPCInferenceServiceServer) ModelFetchRequest(context.Context, *ModelFetchRequestParams) (*ModelFetchRequestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModelFetchRequest not implemented")
}
func (UnimplementedGRPCInferenceServiceServer) ModelSendResponse(GRPCInferenceService_ModelSendResponseServer) error {
	return status.Errorf(codes.Unimplemented, "method ModelSendResponse not implemented")
}
func (UnimplementedGRPCInferenceServiceServer) mustEmbedUnimplementedGRPCInferenceServiceServer() {}

// UnsafeGRPCInferenceServiceServer may be embedded to opt out of forward compatibility for this gateway.
// Use of this interface is not recommended, as added methods to GRPCInferenceServiceServer will
// result in compilation errors.
type UnsafeGRPCInferenceServiceServer interface {
	mustEmbedUnimplementedGRPCInferenceServiceServer()
}

func RegisterGRPCInferenceServiceServer(s grpc.ServiceRegistrar, srv GRPCInferenceServiceServer) {
	s.RegisterService(&GRPCInferenceService_ServiceDesc, srv)
}

func _GRPCInferenceService_ModelStreamInfer_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ModelInferRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GRPCInferenceServiceServer).ModelStreamInfer(m, &gRPCInferenceServiceModelStreamInferServer{ServerStream: stream})
}

type GRPCInferenceService_ModelStreamInferServer interface {
	Send(*ModelInferResponse) error
	grpc.ServerStream
}

type gRPCInferenceServiceModelStreamInferServer struct {
	grpc.ServerStream
}

func (x *gRPCInferenceServiceModelStreamInferServer) Send(m *ModelInferResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GRPCInferenceService_ModelFetchRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModelFetchRequestParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCInferenceServiceServer).ModelFetchRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCInferenceService_ModelFetchRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCInferenceServiceServer).ModelFetchRequest(ctx, req.(*ModelFetchRequestParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _GRPCInferenceService_ModelSendResponse_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GRPCInferenceServiceServer).ModelSendResponse(&gRPCInferenceServiceModelSendResponseServer{ServerStream: stream})
}

type GRPCInferenceService_ModelSendResponseServer interface {
	SendAndClose(*ModelSendResponseResult) error
	Recv() (*ModelInferResponse, error)
	grpc.ServerStream
}

type gRPCInferenceServiceModelSendResponseServer struct {
	grpc.ServerStream
}

func (x *gRPCInferenceServiceModelSendResponseServer) SendAndClose(m *ModelSendResponseResult) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gRPCInferenceServiceModelSendResponseServer) Recv() (*ModelInferResponse, error) {
	m := new(ModelInferResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GRPCInferenceService_ServiceDesc is the grpc.ServiceDesc for GRPCInferenceService gateway.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GRPCInferenceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "language_inference.GRPCInferenceService",
	HandlerType: (*GRPCInferenceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ModelFetchRequest",
			Handler:    _GRPCInferenceService_ModelFetchRequest_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ModelStreamInfer",
			Handler:       _GRPCInferenceService_ModelStreamInfer_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ModelSendResponse",
			Handler:       _GRPCInferenceService_ModelSendResponse_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "inference_engine.proto",
}
