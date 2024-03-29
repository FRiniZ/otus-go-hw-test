// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: EventService.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID          *int64                 `protobuf:"varint,1,opt,name=ID,proto3,oneof" json:"ID,omitempty"`
	UserID      *int64                 `protobuf:"varint,2,opt,name=UserID,proto3,oneof" json:"UserID,omitempty"`
	Title       *string                `protobuf:"bytes,3,opt,name=Title,proto3,oneof" json:"Title,omitempty"`
	Description *string                `protobuf:"bytes,4,opt,name=Description,proto3,oneof" json:"Description,omitempty"`
	OnTime      *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=OnTime,proto3,oneof" json:"OnTime,omitempty"`
	OffTime     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=OffTime,proto3,oneof" json:"OffTime,omitempty"`
	NotifyTime  *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=NotifyTime,proto3,oneof" json:"NotifyTime,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_EventService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_EventService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_EventService_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetID() int64 {
	if x != nil && x.ID != nil {
		return *x.ID
	}
	return 0
}

func (x *Event) GetUserID() int64 {
	if x != nil && x.UserID != nil {
		return *x.UserID
	}
	return 0
}

func (x *Event) GetTitle() string {
	if x != nil && x.Title != nil {
		return *x.Title
	}
	return ""
}

func (x *Event) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Event) GetOnTime() *timestamppb.Timestamp {
	if x != nil {
		return x.OnTime
	}
	return nil
}

func (x *Event) GetOffTime() *timestamppb.Timestamp {
	if x != nil {
		return x.OffTime
	}
	return nil
}

func (x *Event) GetNotifyTime() *timestamppb.Timestamp {
	if x != nil {
		return x.NotifyTime
	}
	return nil
}

type ReqByEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event *Event `protobuf:"bytes,1,opt,name=event,proto3,oneof" json:"event,omitempty"`
}

func (x *ReqByEvent) Reset() {
	*x = ReqByEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_EventService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqByEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqByEvent) ProtoMessage() {}

func (x *ReqByEvent) ProtoReflect() protoreflect.Message {
	mi := &file_EventService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqByEvent.ProtoReflect.Descriptor instead.
func (*ReqByEvent) Descriptor() ([]byte, []int) {
	return file_EventService_proto_rawDescGZIP(), []int{1}
}

func (x *ReqByEvent) GetEvent() *Event {
	if x != nil {
		return x.Event
	}
	return nil
}

type ReqByID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID *int64 `protobuf:"varint,1,opt,name=ID,proto3,oneof" json:"ID,omitempty"`
}

func (x *ReqByID) Reset() {
	*x = ReqByID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_EventService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqByID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqByID) ProtoMessage() {}

func (x *ReqByID) ProtoReflect() protoreflect.Message {
	mi := &file_EventService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqByID.ProtoReflect.Descriptor instead.
func (*ReqByID) Descriptor() ([]byte, []int) {
	return file_EventService_proto_rawDescGZIP(), []int{2}
}

func (x *ReqByID) GetID() int64 {
	if x != nil && x.ID != nil {
		return *x.ID
	}
	return 0
}

type ReqByUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID *int64 `protobuf:"varint,1,opt,name=UserID,proto3,oneof" json:"UserID,omitempty"`
}

func (x *ReqByUser) Reset() {
	*x = ReqByUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_EventService_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqByUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqByUser) ProtoMessage() {}

func (x *ReqByUser) ProtoReflect() protoreflect.Message {
	mi := &file_EventService_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqByUser.ProtoReflect.Descriptor instead.
func (*ReqByUser) Descriptor() ([]byte, []int) {
	return file_EventService_proto_rawDescGZIP(), []int{3}
}

func (x *ReqByUser) GetUserID() int64 {
	if x != nil && x.UserID != nil {
		return *x.UserID
	}
	return 0
}

type ReqByUserByDate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID *int64                 `protobuf:"varint,1,opt,name=UserID,proto3,oneof" json:"UserID,omitempty"`
	Date   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=Date,proto3,oneof" json:"Date,omitempty"`
}

func (x *ReqByUserByDate) Reset() {
	*x = ReqByUserByDate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_EventService_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqByUserByDate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqByUserByDate) ProtoMessage() {}

func (x *ReqByUserByDate) ProtoReflect() protoreflect.Message {
	mi := &file_EventService_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqByUserByDate.ProtoReflect.Descriptor instead.
func (*ReqByUserByDate) Descriptor() ([]byte, []int) {
	return file_EventService_proto_rawDescGZIP(), []int{4}
}

func (x *ReqByUserByDate) GetUserID() int64 {
	if x != nil && x.UserID != nil {
		return *x.UserID
	}
	return 0
}

func (x *ReqByUserByDate) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

type RepID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID *int64 `protobuf:"varint,1,opt,name=ID,proto3,oneof" json:"ID,omitempty"`
}

func (x *RepID) Reset() {
	*x = RepID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_EventService_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepID) ProtoMessage() {}

func (x *RepID) ProtoReflect() protoreflect.Message {
	mi := &file_EventService_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepID.ProtoReflect.Descriptor instead.
func (*RepID) Descriptor() ([]byte, []int) {
	return file_EventService_proto_rawDescGZIP(), []int{5}
}

func (x *RepID) GetID() int64 {
	if x != nil && x.ID != nil {
		return *x.ID
	}
	return 0
}

type RepEvents struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event []*Event `protobuf:"bytes,2,rep,name=event,proto3" json:"event,omitempty"`
}

func (x *RepEvents) Reset() {
	*x = RepEvents{}
	if protoimpl.UnsafeEnabled {
		mi := &file_EventService_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepEvents) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepEvents) ProtoMessage() {}

func (x *RepEvents) ProtoReflect() protoreflect.Message {
	mi := &file_EventService_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepEvents.ProtoReflect.Descriptor instead.
func (*RepEvents) Descriptor() ([]byte, []int) {
	return file_EventService_proto_rawDescGZIP(), []int{6}
}

func (x *RepEvents) GetEvent() []*Event {
	if x != nil {
		return x.Event
	}
	return nil
}

var File_EventService_proto protoreflect.FileDescriptor

var file_EventService_proto_rawDesc = []byte{
	0x0a, 0x12, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x82, 0x03, 0x0a, 0x05, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x13, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x48, 0x00, 0x52, 0x02, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x25, 0x0a, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x0b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x37, 0x0a, 0x06, 0x4f, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x48, 0x04, 0x52, 0x06, 0x4f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x39, 0x0a, 0x07, 0x4f, 0x66, 0x66, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x05,
	0x52, 0x07, 0x4f, 0x66, 0x66, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x3f, 0x0a, 0x0a,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x06, 0x52, 0x0a,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a,
	0x03, 0x5f, 0x49, 0x44, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x44, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4f, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x4f, 0x66, 0x66, 0x54, 0x69, 0x6d, 0x65,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x22,
	0x3d, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x42, 0x79, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a,
	0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x25,
	0x0a, 0x07, 0x52, 0x65, 0x71, 0x42, 0x79, 0x49, 0x44, 0x12, 0x13, 0x0a, 0x02, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x02, 0x49, 0x44, 0x88, 0x01, 0x01, 0x42, 0x05,
	0x0a, 0x03, 0x5f, 0x49, 0x44, 0x22, 0x33, 0x0a, 0x09, 0x52, 0x65, 0x71, 0x42, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x1b, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x48, 0x00, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x88, 0x01, 0x01, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x77, 0x0a, 0x0f, 0x52, 0x65,
	0x71, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a,
	0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52,
	0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x04, 0x44, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52, 0x04, 0x44, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x44,
	0x61, 0x74, 0x65, 0x22, 0x23, 0x0a, 0x05, 0x52, 0x65, 0x70, 0x49, 0x44, 0x12, 0x13, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x02, 0x49, 0x44, 0x88, 0x01,
	0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x49, 0x44, 0x22, 0x2d, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x20, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x73, 0x74, 0x75,
	0x62, 0x2f, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_EventService_proto_rawDescOnce sync.Once
	file_EventService_proto_rawDescData = file_EventService_proto_rawDesc
)

func file_EventService_proto_rawDescGZIP() []byte {
	file_EventService_proto_rawDescOnce.Do(func() {
		file_EventService_proto_rawDescData = protoimpl.X.CompressGZIP(file_EventService_proto_rawDescData)
	})
	return file_EventService_proto_rawDescData
}

var file_EventService_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_EventService_proto_goTypes = []interface{}{
	(*Event)(nil),                 // 0: api.Event
	(*ReqByEvent)(nil),            // 1: api.ReqByEvent
	(*ReqByID)(nil),               // 2: api.ReqByID
	(*ReqByUser)(nil),             // 3: api.ReqByUser
	(*ReqByUserByDate)(nil),       // 4: api.ReqByUserByDate
	(*RepID)(nil),                 // 5: api.RepID
	(*RepEvents)(nil),             // 6: api.RepEvents
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_EventService_proto_depIdxs = []int32{
	7, // 0: api.Event.OnTime:type_name -> google.protobuf.Timestamp
	7, // 1: api.Event.OffTime:type_name -> google.protobuf.Timestamp
	7, // 2: api.Event.NotifyTime:type_name -> google.protobuf.Timestamp
	0, // 3: api.ReqByEvent.event:type_name -> api.Event
	7, // 4: api.ReqByUserByDate.Date:type_name -> google.protobuf.Timestamp
	0, // 5: api.RepEvents.event:type_name -> api.Event
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_EventService_proto_init() }
func file_EventService_proto_init() {
	if File_EventService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_EventService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_EventService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqByEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_EventService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqByID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_EventService_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqByUser); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_EventService_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqByUserByDate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_EventService_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_EventService_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepEvents); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_EventService_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_EventService_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_EventService_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_EventService_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_EventService_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_EventService_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_EventService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_EventService_proto_goTypes,
		DependencyIndexes: file_EventService_proto_depIdxs,
		MessageInfos:      file_EventService_proto_msgTypes,
	}.Build()
	File_EventService_proto = out.File
	file_EventService_proto_rawDesc = nil
	file_EventService_proto_goTypes = nil
	file_EventService_proto_depIdxs = nil
}
