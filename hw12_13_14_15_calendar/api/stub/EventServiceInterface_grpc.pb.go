// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: EventServiceInterface.proto

package api

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

// CalendarClient is the client API for Calendar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalendarClient interface {
	InsertEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error)
	UpdateEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error)
	DeleteEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error)
	LookupEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error)
	ListEventsV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error)
}

type calendarClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarClient(cc grpc.ClientConnInterface) CalendarClient {
	return &calendarClient{cc}
}

func (c *calendarClient) InsertEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error) {
	out := new(ReplyV1)
	err := c.cc.Invoke(ctx, "/api.Calendar/InsertEventV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) UpdateEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error) {
	out := new(ReplyV1)
	err := c.cc.Invoke(ctx, "/api.Calendar/UpdateEventV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) DeleteEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error) {
	out := new(ReplyV1)
	err := c.cc.Invoke(ctx, "/api.Calendar/DeleteEventV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) LookupEventV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error) {
	out := new(ReplyV1)
	err := c.cc.Invoke(ctx, "/api.Calendar/LookupEventV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) ListEventsV1(ctx context.Context, in *RequestV1, opts ...grpc.CallOption) (*ReplyV1, error) {
	out := new(ReplyV1)
	err := c.cc.Invoke(ctx, "/api.Calendar/ListEventsV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServer is the server API for Calendar service.
// All implementations must embed UnimplementedCalendarServer
// for forward compatibility
type CalendarServer interface {
	InsertEventV1(context.Context, *RequestV1) (*ReplyV1, error)
	UpdateEventV1(context.Context, *RequestV1) (*ReplyV1, error)
	DeleteEventV1(context.Context, *RequestV1) (*ReplyV1, error)
	LookupEventV1(context.Context, *RequestV1) (*ReplyV1, error)
	ListEventsV1(context.Context, *RequestV1) (*ReplyV1, error)
	mustEmbedUnimplementedCalendarServer()
}

// UnimplementedCalendarServer must be embedded to have forward compatible implementations.
type UnimplementedCalendarServer struct {
}

func (UnimplementedCalendarServer) InsertEventV1(context.Context, *RequestV1) (*ReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertEventV1 not implemented")
}
func (UnimplementedCalendarServer) UpdateEventV1(context.Context, *RequestV1) (*ReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEventV1 not implemented")
}
func (UnimplementedCalendarServer) DeleteEventV1(context.Context, *RequestV1) (*ReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEventV1 not implemented")
}
func (UnimplementedCalendarServer) LookupEventV1(context.Context, *RequestV1) (*ReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LookupEventV1 not implemented")
}
func (UnimplementedCalendarServer) ListEventsV1(context.Context, *RequestV1) (*ReplyV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsV1 not implemented")
}
func (UnimplementedCalendarServer) mustEmbedUnimplementedCalendarServer() {}

// UnsafeCalendarServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServer will
// result in compilation errors.
type UnsafeCalendarServer interface {
	mustEmbedUnimplementedCalendarServer()
}

func RegisterCalendarServer(s grpc.ServiceRegistrar, srv CalendarServer) {
	s.RegisterService(&Calendar_ServiceDesc, srv)
}

func _Calendar_InsertEventV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).InsertEventV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Calendar/InsertEventV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).InsertEventV1(ctx, req.(*RequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_UpdateEventV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).UpdateEventV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Calendar/UpdateEventV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).UpdateEventV1(ctx, req.(*RequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_DeleteEventV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).DeleteEventV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Calendar/DeleteEventV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).DeleteEventV1(ctx, req.(*RequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_LookupEventV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).LookupEventV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Calendar/LookupEventV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).LookupEventV1(ctx, req.(*RequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_ListEventsV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).ListEventsV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Calendar/ListEventsV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).ListEventsV1(ctx, req.(*RequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

// Calendar_ServiceDesc is the grpc.ServiceDesc for Calendar service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calendar_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Calendar",
	HandlerType: (*CalendarServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InsertEventV1",
			Handler:    _Calendar_InsertEventV1_Handler,
		},
		{
			MethodName: "UpdateEventV1",
			Handler:    _Calendar_UpdateEventV1_Handler,
		},
		{
			MethodName: "DeleteEventV1",
			Handler:    _Calendar_DeleteEventV1_Handler,
		},
		{
			MethodName: "LookupEventV1",
			Handler:    _Calendar_LookupEventV1_Handler,
		},
		{
			MethodName: "ListEventsV1",
			Handler:    _Calendar_ListEventsV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "EventServiceInterface.proto",
}