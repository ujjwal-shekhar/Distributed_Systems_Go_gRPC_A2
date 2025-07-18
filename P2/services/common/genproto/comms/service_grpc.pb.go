// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.3
// source: service.proto

package mapreduce

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FileTransfer_SendToMapper_FullMethodName  = "/FileTransfer/SendToMapper"
	FileTransfer_SendToReducer_FullMethodName = "/FileTransfer/SendToReducer"
	FileTransfer_Vomit_FullMethodName         = "/FileTransfer/Vomit"
	FileTransfer_Close_FullMethodName         = "/FileTransfer/Close"
)

// FileTransferClient is the client API for FileTransfer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileTransferClient interface {
	SendToMapper(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, FileInfo], error)
	SendToReducer(ctx context.Context, in *FileInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Vomit(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Close(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type fileTransferClient struct {
	cc grpc.ClientConnInterface
}

func NewFileTransferClient(cc grpc.ClientConnInterface) FileTransferClient {
	return &fileTransferClient{cc}
}

func (c *fileTransferClient) SendToMapper(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, FileInfo], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileTransfer_ServiceDesc.Streams[0], FileTransfer_SendToMapper_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileChunk, FileInfo]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransfer_SendToMapperClient = grpc.ClientStreamingClient[FileChunk, FileInfo]

func (c *fileTransferClient) SendToReducer(ctx context.Context, in *FileInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, FileTransfer_SendToReducer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) Vomit(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, FileTransfer_Vomit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferClient) Close(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, FileTransfer_Close_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileTransferServer is the server API for FileTransfer service.
// All implementations must embed UnimplementedFileTransferServer
// for forward compatibility.
type FileTransferServer interface {
	SendToMapper(grpc.ClientStreamingServer[FileChunk, FileInfo]) error
	SendToReducer(context.Context, *FileInfo) (*emptypb.Empty, error)
	Vomit(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedFileTransferServer()
}

// UnimplementedFileTransferServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFileTransferServer struct{}

func (UnimplementedFileTransferServer) SendToMapper(grpc.ClientStreamingServer[FileChunk, FileInfo]) error {
	return status.Errorf(codes.Unimplemented, "method SendToMapper not implemented")
}
func (UnimplementedFileTransferServer) SendToReducer(context.Context, *FileInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendToReducer not implemented")
}
func (UnimplementedFileTransferServer) Vomit(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vomit not implemented")
}
func (UnimplementedFileTransferServer) Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedFileTransferServer) mustEmbedUnimplementedFileTransferServer() {}
func (UnimplementedFileTransferServer) testEmbeddedByValue()                      {}

// UnsafeFileTransferServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileTransferServer will
// result in compilation errors.
type UnsafeFileTransferServer interface {
	mustEmbedUnimplementedFileTransferServer()
}

func RegisterFileTransferServer(s grpc.ServiceRegistrar, srv FileTransferServer) {
	// If the following call pancis, it indicates UnimplementedFileTransferServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FileTransfer_ServiceDesc, srv)
}

func _FileTransfer_SendToMapper_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileTransferServer).SendToMapper(&grpc.GenericServerStream[FileChunk, FileInfo]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileTransfer_SendToMapperServer = grpc.ClientStreamingServer[FileChunk, FileInfo]

func _FileTransfer_SendToReducer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).SendToReducer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileTransfer_SendToReducer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).SendToReducer(ctx, req.(*FileInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_Vomit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).Vomit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileTransfer_Vomit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).Vomit(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransfer_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileTransfer_Close_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServer).Close(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// FileTransfer_ServiceDesc is the grpc.ServiceDesc for FileTransfer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileTransfer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "FileTransfer",
	HandlerType: (*FileTransferServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendToReducer",
			Handler:    _FileTransfer_SendToReducer_Handler,
		},
		{
			MethodName: "Vomit",
			Handler:    _FileTransfer_Vomit_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _FileTransfer_Close_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendToMapper",
			Handler:       _FileTransfer_SendToMapper_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}
