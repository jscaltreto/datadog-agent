// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

// copied from https://github.com/DataDog/dd-go/blob/prod/pb/logspb/internal_intake_request.pb.go

package fleetautomationextension

import (
	"reflect"
	"sync"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type InternalIntakeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*InternalIntakeRequest_Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *InternalIntakeRequest) Reset() {
	*x = InternalIntakeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_internal_intake_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InternalIntakeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InternalIntakeRequest) ProtoMessage() {}

func (x *InternalIntakeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_logs_internal_intake_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InternalIntakeRequest.ProtoReflect.Descriptor instead.
func (*InternalIntakeRequest) Descriptor() ([]byte, []int) {
	return file_logs_internal_intake_request_proto_rawDescGZIP(), []int{0}
}

func (x *InternalIntakeRequest) GetEvents() []*InternalIntakeRequest_Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type InternalIntakeRequest_Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload   []byte                `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	OrgId     *wrappers.Int64Value  `protobuf:"bytes,2,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Uuid      *wrappers.StringValue `protobuf:"bytes,3,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Timestamp *wrappers.Int64Value  `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *InternalIntakeRequest_Event) Reset() {
	*x = InternalIntakeRequest_Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_internal_intake_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InternalIntakeRequest_Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InternalIntakeRequest_Event) ProtoMessage() {}

func (x *InternalIntakeRequest_Event) ProtoReflect() protoreflect.Message {
	mi := &file_logs_internal_intake_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InternalIntakeRequest_Event.ProtoReflect.Descriptor instead.
func (*InternalIntakeRequest_Event) Descriptor() ([]byte, []int) {
	return file_logs_internal_intake_request_proto_rawDescGZIP(), []int{0, 0}
}

func (x *InternalIntakeRequest_Event) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *InternalIntakeRequest_Event) GetOrgId() *wrappers.Int64Value {
	if x != nil {
		return x.OrgId
	}
	return nil
}

func (x *InternalIntakeRequest_Event) GetUuid() *wrappers.StringValue {
	if x != nil {
		return x.Uuid
	}
	return nil
}

func (x *InternalIntakeRequest_Event) GetTimestamp() *wrappers.Int64Value {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

var File_logs_internal_intake_request_proto protoreflect.FileDescriptor

var file_logs_internal_intake_request_proto_rawDesc = []byte{
	0x0a, 0x22, 0x6c, 0x6f, 0x67, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f,
	0x69, 0x6e, 0x74, 0x61, 0x6b, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x02, 0x0a, 0x15, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x6e, 0x74, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x37, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70, 0x62, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x49, 0x6e, 0x74, 0x61, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0xc2, 0x01, 0x0a, 0x05,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x32, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x6f, 0x72,
	0x67, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x36, 0x34,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x44,
	0x61, 0x74, 0x61, 0x44, 0x6f, 0x67, 0x2f, 0x64, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x62, 0x2f,
	0x6c, 0x6f, 0x67, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_logs_internal_intake_request_proto_rawDescOnce sync.Once
	file_logs_internal_intake_request_proto_rawDescData = file_logs_internal_intake_request_proto_rawDesc
)

func file_logs_internal_intake_request_proto_rawDescGZIP() []byte {
	file_logs_internal_intake_request_proto_rawDescOnce.Do(func() {
		file_logs_internal_intake_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_logs_internal_intake_request_proto_rawDescData)
	})
	return file_logs_internal_intake_request_proto_rawDescData
}

var file_logs_internal_intake_request_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_logs_internal_intake_request_proto_goTypes = []any{
	(*InternalIntakeRequest)(nil),       // 0: pb.InternalIntakeRequest
	(*InternalIntakeRequest_Event)(nil), // 1: pb.InternalIntakeRequest.Event
	(*wrappers.Int64Value)(nil),         // 2: google.protobuf.Int64Value
	(*wrappers.StringValue)(nil),        // 3: google.protobuf.StringValue
}
var file_logs_internal_intake_request_proto_depIdxs = []int32{
	1, // 0: pb.InternalIntakeRequest.events:type_name -> pb.InternalIntakeRequest.Event
	2, // 1: pb.InternalIntakeRequest.Event.org_id:type_name -> google.protobuf.Int64Value
	3, // 2: pb.InternalIntakeRequest.Event.uuid:type_name -> google.protobuf.StringValue
	2, // 3: pb.InternalIntakeRequest.Event.timestamp:type_name -> google.protobuf.Int64Value
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_logs_internal_intake_request_proto_init() }
func file_logs_internal_intake_request_proto_init() {
	if File_logs_internal_intake_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_logs_internal_intake_request_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*InternalIntakeRequest); i {
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
		file_logs_internal_intake_request_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*InternalIntakeRequest_Event); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_logs_internal_intake_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_logs_internal_intake_request_proto_goTypes,
		DependencyIndexes: file_logs_internal_intake_request_proto_depIdxs,
		MessageInfos:      file_logs_internal_intake_request_proto_msgTypes,
	}.Build()
	File_logs_internal_intake_request_proto = out.File
	file_logs_internal_intake_request_proto_rawDesc = nil
	file_logs_internal_intake_request_proto_goTypes = nil
	file_logs_internal_intake_request_proto_depIdxs = nil
}
