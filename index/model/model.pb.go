// Code generated by protoc-gen-go. DO NOT EDIT.
// source: model.proto

package model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Entry_Type int32

const (
	Entry_TREE Entry_Type = 0
	Entry_BLOB Entry_Type = 1
)

var Entry_Type_name = map[int32]string{
	0: "TREE",
	1: "BLOB",
}

var Entry_Type_value = map[string]int32{
	"TREE": 0,
	"BLOB": 1,
}

func (x Entry_Type) String() string {
	return proto.EnumName(Entry_Type_name, int32(x))
}

func (Entry_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{2, 0}
}

// Data model
type Repo struct {
	ClientId             string   `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	RepoId               string   `protobuf:"bytes,2,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`
	DefaultBranch        string   `protobuf:"bytes,3,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	PartialCommitRatio   float32  `protobuf:"fixed32,4,opt,name=partial_commit_ratio,json=partialCommitRatio,proto3" json:"partial_commit_ratio,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Repo) Reset()         { *m = Repo{} }
func (m *Repo) String() string { return proto.CompactTextString(m) }
func (*Repo) ProtoMessage()    {}
func (*Repo) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{0}
}

func (m *Repo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Repo.Unmarshal(m, b)
}
func (m *Repo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Repo.Marshal(b, m, deterministic)
}
func (m *Repo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Repo.Merge(m, src)
}
func (m *Repo) XXX_Size() int {
	return xxx_messageInfo_Repo.Size(m)
}
func (m *Repo) XXX_DiscardUnknown() {
	xxx_messageInfo_Repo.DiscardUnknown(m)
}

var xxx_messageInfo_Repo proto.InternalMessageInfo

func (m *Repo) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *Repo) GetRepoId() string {
	if m != nil {
		return m.RepoId
	}
	return ""
}

func (m *Repo) GetDefaultBranch() string {
	if m != nil {
		return m.DefaultBranch
	}
	return ""
}

func (m *Repo) GetPartialCommitRatio() float32 {
	if m != nil {
		return m.PartialCommitRatio
	}
	return 0
}

type Blob struct {
	Blocks               []string `protobuf:"bytes,4,rep,name=blocks,proto3" json:"blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Blob) Reset()         { *m = Blob{} }
func (m *Blob) String() string { return proto.CompactTextString(m) }
func (*Blob) ProtoMessage()    {}
func (*Blob) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{1}
}

func (m *Blob) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Blob.Unmarshal(m, b)
}
func (m *Blob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Blob.Marshal(b, m, deterministic)
}
func (m *Blob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Blob.Merge(m, src)
}
func (m *Blob) XXX_Size() int {
	return xxx_messageInfo_Blob.Size(m)
}
func (m *Blob) XXX_DiscardUnknown() {
	xxx_messageInfo_Blob.DiscardUnknown(m)
}

var xxx_messageInfo_Blob proto.InternalMessageInfo

func (m *Blob) GetBlocks() []string {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type Entry struct {
	Name                 string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Address              string            `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Type                 Entry_Type        `protobuf:"varint,3,opt,name=type,proto3,enum=Entry_Type" json:"type,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Timestamp            int64             `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{2}
}

func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Entry) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Entry) GetType() Entry_Type {
	if m != nil {
		return m.Type
	}
	return Entry_TREE
}

func (m *Entry) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Entry) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type Commit struct {
	Tree                 string            `protobuf:"bytes,1,opt,name=tree,proto3" json:"tree,omitempty"`
	Parents              []string          `protobuf:"bytes,2,rep,name=parents,proto3" json:"parents,omitempty"`
	Committer            string            `protobuf:"bytes,3,opt,name=committer,proto3" json:"committer,omitempty"`
	Message              string            `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	Timestamp            int64             `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,6,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Commit) Reset()         { *m = Commit{} }
func (m *Commit) String() string { return proto.CompactTextString(m) }
func (*Commit) ProtoMessage()    {}
func (*Commit) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{3}
}

func (m *Commit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Commit.Unmarshal(m, b)
}
func (m *Commit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Commit.Marshal(b, m, deterministic)
}
func (m *Commit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Commit.Merge(m, src)
}
func (m *Commit) XXX_Size() int {
	return xxx_messageInfo_Commit.Size(m)
}
func (m *Commit) XXX_DiscardUnknown() {
	xxx_messageInfo_Commit.DiscardUnknown(m)
}

var xxx_messageInfo_Commit proto.InternalMessageInfo

func (m *Commit) GetTree() string {
	if m != nil {
		return m.Tree
	}
	return ""
}

func (m *Commit) GetParents() []string {
	if m != nil {
		return m.Parents
	}
	return nil
}

func (m *Commit) GetCommitter() string {
	if m != nil {
		return m.Committer
	}
	return ""
}

func (m *Commit) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Commit) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Commit) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Branch struct {
	Commit               string   `protobuf:"bytes,1,opt,name=commit,proto3" json:"commit,omitempty"`
	CommitRoot           string   `protobuf:"bytes,2,opt,name=commit_root,json=commitRoot,proto3" json:"commit_root,omitempty"`
	WorkspaceRoot        string   `protobuf:"bytes,3,opt,name=workspace_root,json=workspaceRoot,proto3" json:"workspace_root,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Branch) Reset()         { *m = Branch{} }
func (m *Branch) String() string { return proto.CompactTextString(m) }
func (*Branch) ProtoMessage()    {}
func (*Branch) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{4}
}

func (m *Branch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Branch.Unmarshal(m, b)
}
func (m *Branch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Branch.Marshal(b, m, deterministic)
}
func (m *Branch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Branch.Merge(m, src)
}
func (m *Branch) XXX_Size() int {
	return xxx_messageInfo_Branch.Size(m)
}
func (m *Branch) XXX_DiscardUnknown() {
	xxx_messageInfo_Branch.DiscardUnknown(m)
}

var xxx_messageInfo_Branch proto.InternalMessageInfo

func (m *Branch) GetCommit() string {
	if m != nil {
		return m.Commit
	}
	return ""
}

func (m *Branch) GetCommitRoot() string {
	if m != nil {
		return m.CommitRoot
	}
	return ""
}

func (m *Branch) GetWorkspaceRoot() string {
	if m != nil {
		return m.WorkspaceRoot
	}
	return ""
}

type Tombstone struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tombstone) Reset()         { *m = Tombstone{} }
func (m *Tombstone) String() string { return proto.CompactTextString(m) }
func (*Tombstone) ProtoMessage()    {}
func (*Tombstone) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{5}
}

func (m *Tombstone) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tombstone.Unmarshal(m, b)
}
func (m *Tombstone) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tombstone.Marshal(b, m, deterministic)
}
func (m *Tombstone) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tombstone.Merge(m, src)
}
func (m *Tombstone) XXX_Size() int {
	return xxx_messageInfo_Tombstone.Size(m)
}
func (m *Tombstone) XXX_DiscardUnknown() {
	xxx_messageInfo_Tombstone.DiscardUnknown(m)
}

var xxx_messageInfo_Tombstone proto.InternalMessageInfo

type Object struct {
	Blob                 *Blob             `protobuf:"bytes,1,opt,name=blob,proto3" json:"blob,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Timestamp            int64             `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Size                 int64             `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Object) Reset()         { *m = Object{} }
func (m *Object) String() string { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()    {}
func (*Object) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{6}
}

func (m *Object) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Object.Unmarshal(m, b)
}
func (m *Object) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Object.Marshal(b, m, deterministic)
}
func (m *Object) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object.Merge(m, src)
}
func (m *Object) XXX_Size() int {
	return xxx_messageInfo_Object.Size(m)
}
func (m *Object) XXX_DiscardUnknown() {
	xxx_messageInfo_Object.DiscardUnknown(m)
}

var xxx_messageInfo_Object proto.InternalMessageInfo

func (m *Object) GetBlob() *Blob {
	if m != nil {
		return m.Blob
	}
	return nil
}

func (m *Object) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Object) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Object) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

type WorkspaceEntry struct {
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// Types that are valid to be assigned to Data:
	//	*WorkspaceEntry_Tombstone
	//	*WorkspaceEntry_Object
	Data                 isWorkspaceEntry_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *WorkspaceEntry) Reset()         { *m = WorkspaceEntry{} }
func (m *WorkspaceEntry) String() string { return proto.CompactTextString(m) }
func (*WorkspaceEntry) ProtoMessage()    {}
func (*WorkspaceEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{7}
}

func (m *WorkspaceEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkspaceEntry.Unmarshal(m, b)
}
func (m *WorkspaceEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkspaceEntry.Marshal(b, m, deterministic)
}
func (m *WorkspaceEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkspaceEntry.Merge(m, src)
}
func (m *WorkspaceEntry) XXX_Size() int {
	return xxx_messageInfo_WorkspaceEntry.Size(m)
}
func (m *WorkspaceEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkspaceEntry.DiscardUnknown(m)
}

var xxx_messageInfo_WorkspaceEntry proto.InternalMessageInfo

func (m *WorkspaceEntry) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type isWorkspaceEntry_Data interface {
	isWorkspaceEntry_Data()
}

type WorkspaceEntry_Tombstone struct {
	Tombstone *Tombstone `protobuf:"bytes,2,opt,name=tombstone,proto3,oneof"`
}

type WorkspaceEntry_Object struct {
	Object *Object `protobuf:"bytes,3,opt,name=object,proto3,oneof"`
}

func (*WorkspaceEntry_Tombstone) isWorkspaceEntry_Data() {}

func (*WorkspaceEntry_Object) isWorkspaceEntry_Data() {}

func (m *WorkspaceEntry) GetData() isWorkspaceEntry_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *WorkspaceEntry) GetTombstone() *Tombstone {
	if x, ok := m.GetData().(*WorkspaceEntry_Tombstone); ok {
		return x.Tombstone
	}
	return nil
}

func (m *WorkspaceEntry) GetObject() *Object {
	if x, ok := m.GetData().(*WorkspaceEntry_Object); ok {
		return x.Object
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*WorkspaceEntry) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*WorkspaceEntry_Tombstone)(nil),
		(*WorkspaceEntry_Object)(nil),
	}
}

func init() {
	proto.RegisterEnum("Entry_Type", Entry_Type_name, Entry_Type_value)
	proto.RegisterType((*Repo)(nil), "Repo")
	proto.RegisterType((*Blob)(nil), "Blob")
	proto.RegisterType((*Entry)(nil), "Entry")
	proto.RegisterMapType((map[string]string)(nil), "Entry.MetadataEntry")
	proto.RegisterType((*Commit)(nil), "Commit")
	proto.RegisterMapType((map[string]string)(nil), "Commit.MetadataEntry")
	proto.RegisterType((*Branch)(nil), "Branch")
	proto.RegisterType((*Tombstone)(nil), "Tombstone")
	proto.RegisterType((*Object)(nil), "Object")
	proto.RegisterMapType((map[string]string)(nil), "Object.MetadataEntry")
	proto.RegisterType((*WorkspaceEntry)(nil), "WorkspaceEntry")
}

func init() { proto.RegisterFile("model.proto", fileDescriptor_4c16552f9fdb66d8) }

var fileDescriptor_4c16552f9fdb66d8 = []byte{
	// 558 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x5f, 0x8b, 0xd3, 0x4e,
	0x14, 0x6d, 0x9a, 0x34, 0x6d, 0x6e, 0xd8, 0x52, 0x86, 0xfe, 0xf6, 0x17, 0x57, 0x71, 0x6b, 0x40,
	0x28, 0x3e, 0x84, 0xb5, 0xbe, 0x88, 0xbe, 0x55, 0x0a, 0xbb, 0xa0, 0x2c, 0x0c, 0x05, 0x1f, 0xcb,
	0x24, 0x19, 0x6d, 0x6c, 0x92, 0x09, 0x93, 0xbb, 0x4a, 0xf5, 0x73, 0xf8, 0x91, 0xfc, 0x56, 0x82,
	0x32, 0x7f, 0xd2, 0xba, 0x82, 0x22, 0xec, 0xdb, 0x3d, 0xe7, 0xcc, 0x9f, 0x73, 0xcf, 0xdc, 0x04,
	0xc2, 0x4a, 0xe4, 0xbc, 0x4c, 0x1a, 0x29, 0x50, 0xc4, 0x5f, 0x1d, 0xf0, 0x28, 0x6f, 0x04, 0xb9,
	0x0f, 0x41, 0x56, 0x16, 0xbc, 0xc6, 0x4d, 0x91, 0x47, 0xce, 0xcc, 0x99, 0x07, 0x74, 0x64, 0x88,
	0xab, 0x9c, 0xfc, 0x0f, 0x43, 0xc9, 0x1b, 0xa1, 0xa4, 0xbe, 0x96, 0x7c, 0x05, 0xaf, 0x72, 0xf2,
	0x18, 0xc6, 0x39, 0x7f, 0xc7, 0x6e, 0x4a, 0xdc, 0xa4, 0x92, 0xd5, 0xd9, 0x36, 0x72, 0xb5, 0x7e,
	0x62, 0xd9, 0xa5, 0x26, 0xc9, 0x05, 0x4c, 0x1b, 0x26, 0xb1, 0x60, 0xe5, 0x26, 0x13, 0x55, 0x55,
	0xe0, 0x46, 0x32, 0x2c, 0x44, 0xe4, 0xcd, 0x9c, 0x79, 0x9f, 0x12, 0xab, 0xbd, 0xd2, 0x12, 0x55,
	0x4a, 0xfc, 0x10, 0xbc, 0x65, 0x29, 0x52, 0x72, 0x0a, 0x7e, 0x5a, 0x8a, 0x6c, 0xd7, 0x46, 0xde,
	0xcc, 0x55, 0x17, 0x1b, 0x14, 0xff, 0x70, 0x60, 0xb0, 0xaa, 0x51, 0xee, 0x09, 0x01, 0xaf, 0x66,
	0x15, 0xb7, 0x9e, 0x75, 0x4d, 0x22, 0x18, 0xb2, 0x3c, 0x97, 0xbc, 0x6d, 0xad, 0xdf, 0x0e, 0x92,
	0x73, 0xf0, 0x70, 0xdf, 0x70, 0x6d, 0x73, 0xbc, 0x08, 0x13, 0x7d, 0x46, 0xb2, 0xde, 0x37, 0x9c,
	0x6a, 0x81, 0x5c, 0xc0, 0xa8, 0xe2, 0xc8, 0x72, 0x86, 0x4c, 0x5f, 0x19, 0x2e, 0xa6, 0x76, 0xd1,
	0x1b, 0x4b, 0x6b, 0x44, 0x0f, 0xab, 0xc8, 0x03, 0x08, 0xb0, 0xa8, 0x78, 0x8b, 0xac, 0x6a, 0xa2,
	0xc1, 0xcc, 0x99, 0xbb, 0xf4, 0x48, 0x9c, 0xbd, 0x84, 0x93, 0x5b, 0x1b, 0xc9, 0x04, 0xdc, 0x1d,
	0xdf, 0x5b, 0xbb, 0xaa, 0x24, 0x53, 0x18, 0x7c, 0x64, 0xe5, 0x0d, 0xb7, 0x5e, 0x0d, 0x78, 0xd1,
	0x7f, 0xee, 0xc4, 0x67, 0xe0, 0x29, 0x6b, 0x64, 0x04, 0xde, 0x9a, 0xae, 0x56, 0x93, 0x9e, 0xaa,
	0x96, 0xaf, 0xaf, 0x97, 0x13, 0x27, 0xfe, 0xee, 0x80, 0x6f, 0x12, 0x53, 0x11, 0xa0, 0xe4, 0x87,
	0x08, 0x54, 0xad, 0x22, 0x68, 0x98, 0xe4, 0x35, 0xaa, 0x08, 0x54, 0x72, 0x1d, 0x54, 0x7e, 0xcd,
	0x23, 0x20, 0x97, 0xf6, 0xb9, 0x8e, 0x84, 0xda, 0x57, 0xf1, 0xb6, 0x65, 0xef, 0xb9, 0x7e, 0x9d,
	0x80, 0x76, 0xf0, 0xef, 0x7d, 0x92, 0xa7, 0xbf, 0xe4, 0xe6, 0xeb, 0xdc, 0xfe, 0x4b, 0x8c, 0xbd,
	0x3f, 0x05, 0x77, 0xb7, 0x68, 0xb6, 0xe0, 0xdb, 0xe1, 0x3a, 0x05, 0xdf, 0xd8, 0xb7, 0x1b, 0x2d,
	0x22, 0xe7, 0x10, 0x76, 0xc3, 0x26, 0x04, 0xda, 0x13, 0xc0, 0x50, 0x54, 0x08, 0x54, 0xc3, 0xfb,
	0x49, 0xc8, 0x5d, 0xdb, 0xb0, 0x8c, 0x9b, 0x35, 0x76, 0x78, 0x0f, 0xac, 0x5a, 0x16, 0x87, 0x10,
	0xac, 0x45, 0x95, 0xb6, 0x28, 0x6a, 0x1e, 0x7f, 0x73, 0xc0, 0xbf, 0x4e, 0x3f, 0xf0, 0x0c, 0xc9,
	0x3d, 0xf0, 0xd2, 0x52, 0xa4, 0xfa, 0xd6, 0x70, 0x31, 0x48, 0xd4, 0xbc, 0x52, 0x4d, 0xdd, 0x0a,
	0xa3, 0x6f, 0xc3, 0x30, 0xbb, 0xfe, 0x6d, 0x8a, 0xdc, 0xdf, 0xd3, 0x25, 0xe0, 0xb5, 0xc5, 0x67,
	0xf3, 0x24, 0x2e, 0xd5, 0xf5, 0xdd, 0xe2, 0xfb, 0x02, 0xe3, 0xb7, 0x5d, 0x97, 0x87, 0xef, 0xa8,
	0x61, 0xb8, 0xed, 0x86, 0x48, 0xd5, 0xe4, 0x09, 0x04, 0xd8, 0xb5, 0xae, 0xcf, 0x08, 0x17, 0x90,
	0x1c, 0xc2, 0xb8, 0xec, 0xd1, 0xa3, 0x4c, 0x1e, 0x81, 0x2f, 0x74, 0x8b, 0xda, 0x7d, 0xb8, 0x18,
	0xda, 0x8e, 0x2f, 0x7b, 0xd4, 0x0a, 0x4b, 0x1f, 0x3c, 0xe5, 0x36, 0xf5, 0xf5, 0xbf, 0xe7, 0xd9,
	0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x00, 0x34, 0xe9, 0x8a, 0x04, 0x00, 0x00,
}
