//
//Run from lakeFS root:
//protoc --proto_path=pkg/plugins/diff --go_out=pkg/plugins/diff --go_opt=paths=source_relative \
//--go-grpc_out=pkg/plugins/diff --go-grpc_opt=paths=source_relative diff.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: diff.proto

package diff

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

type DiffPaths struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeftPath      string  `protobuf:"bytes,1,opt,name=left_path,json=leftPath,proto3" json:"left_path,omitempty"`
	RightPath     string  `protobuf:"bytes,2,opt,name=right_path,json=rightPath,proto3" json:"right_path,omitempty"`
	BaseTablePath *string `protobuf:"bytes,3,opt,name=base_table_path,json=baseTablePath,proto3,oneof" json:"base_table_path,omitempty"`
}

func (x *DiffPaths) Reset() {
	*x = DiffPaths{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diff_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiffPaths) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiffPaths) ProtoMessage() {}

func (x *DiffPaths) ProtoReflect() protoreflect.Message {
	mi := &file_diff_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiffPaths.ProtoReflect.Descriptor instead.
func (*DiffPaths) Descriptor() ([]byte, []int) {
	return file_diff_proto_rawDescGZIP(), []int{0}
}

func (x *DiffPaths) GetLeftPath() string {
	if x != nil {
		return x.LeftPath
	}
	return ""
}

func (x *DiffPaths) GetRightPath() string {
	if x != nil {
		return x.RightPath
	}
	return ""
}

func (x *DiffPaths) GetBaseTablePath() string {
	if x != nil && x.BaseTablePath != nil {
		return *x.BaseTablePath
	}
	return ""
}

type GatewayConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Secret   string `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	Endpoint string `protobuf:"bytes,3,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
}

func (x *GatewayConfig) Reset() {
	*x = GatewayConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diff_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayConfig) ProtoMessage() {}

func (x *GatewayConfig) ProtoReflect() protoreflect.Message {
	mi := &file_diff_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayConfig.ProtoReflect.Descriptor instead.
func (*GatewayConfig) Descriptor() ([]byte, []int) {
	return file_diff_proto_rawDescGZIP(), []int{1}
}

func (x *GatewayConfig) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GatewayConfig) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

func (x *GatewayConfig) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

type DiffRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Paths         *DiffPaths     `protobuf:"bytes,1,opt,name=paths,proto3" json:"paths,omitempty"`
	GatewayConfig *GatewayConfig `protobuf:"bytes,2,opt,name=gateway_config,json=gatewayConfig,proto3" json:"gateway_config,omitempty"`
}

func (x *DiffRequest) Reset() {
	*x = DiffRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diff_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiffRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiffRequest) ProtoMessage() {}

func (x *DiffRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diff_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiffRequest.ProtoReflect.Descriptor instead.
func (*DiffRequest) Descriptor() ([]byte, []int) {
	return file_diff_proto_rawDescGZIP(), []int{2}
}

func (x *DiffRequest) GetPaths() *DiffPaths {
	if x != nil {
		return x.Paths
	}
	return nil
}

func (x *DiffRequest) GetGatewayConfig() *GatewayConfig {
	if x != nil {
		return x.GatewayConfig
	}
	return nil
}

type Diff struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version     uint32                 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Timestamp   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Description string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Content     map[string]string      `protobuf:"bytes,4,rep,name=content,proto3" json:"content,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Diff) Reset() {
	*x = Diff{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diff_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Diff) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Diff) ProtoMessage() {}

func (x *Diff) ProtoReflect() protoreflect.Message {
	mi := &file_diff_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Diff.ProtoReflect.Descriptor instead.
func (*Diff) Descriptor() ([]byte, []int) {
	return file_diff_proto_rawDescGZIP(), []int{3}
}

func (x *Diff) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Diff) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Diff) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Diff) GetContent() map[string]string {
	if x != nil {
		return x.Content
	}
	return nil
}

type DiffResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Diffs []*Diff `protobuf:"bytes,1,rep,name=diffs,proto3" json:"diffs,omitempty"`
}

func (x *DiffResponse) Reset() {
	*x = DiffResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diff_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiffResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiffResponse) ProtoMessage() {}

func (x *DiffResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diff_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiffResponse.ProtoReflect.Descriptor instead.
func (*DiffResponse) Descriptor() ([]byte, []int) {
	return file_diff_proto_rawDescGZIP(), []int{4}
}

func (x *DiffResponse) GetDiffs() []*Diff {
	if x != nil {
		return x.Diffs
	}
	return nil
}

var File_diff_proto protoreflect.FileDescriptor

var file_diff_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x69,
	0x66, 0x66, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x88, 0x01, 0x0a, 0x09, 0x44, 0x69, 0x66, 0x66, 0x50, 0x61, 0x74, 0x68,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x66, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x65, 0x66, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1d,
	0x0a, 0x0a, 0x72, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x72, 0x69, 0x67, 0x68, 0x74, 0x50, 0x61, 0x74, 0x68, 0x12, 0x2b, 0x0a,
	0x0f, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0d, 0x62, 0x61, 0x73, 0x65, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x88, 0x01, 0x01, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x62,
	0x61, 0x73, 0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x22, 0x55,
	0x0a, 0x0d, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x70, 0x0a, 0x0b, 0x44, 0x69, 0x66, 0x66, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x44, 0x69, 0x66, 0x66, 0x50,
	0x61, 0x74, 0x68, 0x73, 0x52, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x12, 0x3a, 0x0a, 0x0e, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x47, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xeb, 0x01, 0x0a, 0x04, 0x44, 0x69, 0x66, 0x66,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x31, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x44,
	0x69, 0x66, 0x66, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x1a, 0x3a, 0x0a, 0x0c, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x30, 0x0a, 0x0c, 0x44, 0x69, 0x66, 0x66, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x05, 0x64, 0x69, 0x66, 0x66, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x64, 0x69, 0x66, 0x66, 0x2e, 0x44, 0x69, 0x66, 0x66,
	0x52, 0x05, 0x64, 0x69, 0x66, 0x66, 0x73, 0x32, 0x37, 0x0a, 0x06, 0x44, 0x69, 0x66, 0x66, 0x65,
	0x72, 0x12, 0x2d, 0x0a, 0x04, 0x44, 0x69, 0x66, 0x66, 0x12, 0x11, 0x2e, 0x64, 0x69, 0x66, 0x66,
	0x2e, 0x44, 0x69, 0x66, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x64,
	0x69, 0x66, 0x66, 0x2e, 0x44, 0x69, 0x66, 0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74,
	0x72, 0x65, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x2f, 0x6c, 0x61, 0x6b, 0x65, 0x66, 0x73, 0x2f,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x73, 0x2f, 0x64, 0x69, 0x66, 0x66, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_diff_proto_rawDescOnce sync.Once
	file_diff_proto_rawDescData = file_diff_proto_rawDesc
)

func file_diff_proto_rawDescGZIP() []byte {
	file_diff_proto_rawDescOnce.Do(func() {
		file_diff_proto_rawDescData = protoimpl.X.CompressGZIP(file_diff_proto_rawDescData)
	})
	return file_diff_proto_rawDescData
}

var file_diff_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_diff_proto_goTypes = []interface{}{
	(*DiffPaths)(nil),             // 0: diff.DiffPaths
	(*GatewayConfig)(nil),         // 1: diff.GatewayConfig
	(*DiffRequest)(nil),           // 2: diff.DiffRequest
	(*Diff)(nil),                  // 3: diff.Diff
	(*DiffResponse)(nil),          // 4: diff.DiffResponse
	nil,                           // 5: diff.Diff.ContentEntry
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_diff_proto_depIdxs = []int32{
	0, // 0: diff.DiffRequest.paths:type_name -> diff.DiffPaths
	1, // 1: diff.DiffRequest.gateway_config:type_name -> diff.GatewayConfig
	6, // 2: diff.Diff.timestamp:type_name -> google.protobuf.Timestamp
	5, // 3: diff.Diff.content:type_name -> diff.Diff.ContentEntry
	3, // 4: diff.DiffResponse.diffs:type_name -> diff.Diff
	2, // 5: diff.Differ.Diff:input_type -> diff.DiffRequest
	4, // 6: diff.Differ.Diff:output_type -> diff.DiffResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_diff_proto_init() }
func file_diff_proto_init() {
	if File_diff_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_diff_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiffPaths); i {
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
		file_diff_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayConfig); i {
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
		file_diff_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiffRequest); i {
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
		file_diff_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Diff); i {
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
		file_diff_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiffResponse); i {
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
	file_diff_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_diff_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_diff_proto_goTypes,
		DependencyIndexes: file_diff_proto_depIdxs,
		MessageInfos:      file_diff_proto_msgTypes,
	}.Build()
	File_diff_proto = out.File
	file_diff_proto_rawDesc = nil
	file_diff_proto_goTypes = nil
	file_diff_proto_depIdxs = nil
}
