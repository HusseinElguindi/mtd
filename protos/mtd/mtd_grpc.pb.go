// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mtd

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

// MtdClient is the client API for Mtd service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MtdClient interface {
	RequestHTTPDownload(ctx context.Context, in *HTTPDownloadRequest, opts ...grpc.CallOption) (*HTTPDownloadResponse, error)
	RequestDownloadInfo(ctx context.Context, in *DownloadInfoRequest, opts ...grpc.CallOption) (*DownloadInfoResponse, error)
}

type mtdClient struct {
	cc grpc.ClientConnInterface
}

func NewMtdClient(cc grpc.ClientConnInterface) MtdClient {
	return &mtdClient{cc}
}

func (c *mtdClient) RequestHTTPDownload(ctx context.Context, in *HTTPDownloadRequest, opts ...grpc.CallOption) (*HTTPDownloadResponse, error) {
	out := new(HTTPDownloadResponse)
	err := c.cc.Invoke(ctx, "/mtd.mtd/RequestHTTPDownload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mtdClient) RequestDownloadInfo(ctx context.Context, in *DownloadInfoRequest, opts ...grpc.CallOption) (*DownloadInfoResponse, error) {
	out := new(DownloadInfoResponse)
	err := c.cc.Invoke(ctx, "/mtd.mtd/RequestDownloadInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MtdServer is the server API for Mtd service.
// All implementations must embed UnimplementedMtdServer
// for forward compatibility
type MtdServer interface {
	RequestHTTPDownload(context.Context, *HTTPDownloadRequest) (*HTTPDownloadResponse, error)
	RequestDownloadInfo(context.Context, *DownloadInfoRequest) (*DownloadInfoResponse, error)
	mustEmbedUnimplementedMtdServer()
}

// UnimplementedMtdServer must be embedded to have forward compatible implementations.
type UnimplementedMtdServer struct {
}

func (UnimplementedMtdServer) RequestHTTPDownload(context.Context, *HTTPDownloadRequest) (*HTTPDownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestHTTPDownload not implemented")
}
func (UnimplementedMtdServer) RequestDownloadInfo(context.Context, *DownloadInfoRequest) (*DownloadInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestDownloadInfo not implemented")
}
func (UnimplementedMtdServer) mustEmbedUnimplementedMtdServer() {}

// UnsafeMtdServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MtdServer will
// result in compilation errors.
type UnsafeMtdServer interface {
	mustEmbedUnimplementedMtdServer()
}

func RegisterMtdServer(s grpc.ServiceRegistrar, srv MtdServer) {
	s.RegisterService(&Mtd_ServiceDesc, srv)
}

func _Mtd_RequestHTTPDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HTTPDownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MtdServer).RequestHTTPDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mtd.mtd/RequestHTTPDownload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MtdServer).RequestHTTPDownload(ctx, req.(*HTTPDownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mtd_RequestDownloadInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MtdServer).RequestDownloadInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mtd.mtd/RequestDownloadInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MtdServer).RequestDownloadInfo(ctx, req.(*DownloadInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Mtd_ServiceDesc is the grpc.ServiceDesc for Mtd service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mtd_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mtd.mtd",
	HandlerType: (*MtdServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestHTTPDownload",
			Handler:    _Mtd_RequestHTTPDownload_Handler,
		},
		{
			MethodName: "RequestDownloadInfo",
			Handler:    _Mtd_RequestDownloadInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mtd.proto",
}