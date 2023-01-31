// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: crd.proto

package crd

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ColumnType int32

const (
	// Unspecified column type
	ColumnType_CT_NONE ColumnType = 0
	// Non-floating-point numbers
	ColumnType_CT_INTEGER ColumnType = 1
	// Floating point numbers
	ColumnType_CT_NUMBER ColumnType = 2
	// Strings
	ColumnType_CT_STRING ColumnType = 3
	// true or false
	ColumnType_CT_BOOLEAN ColumnType = 4
	// Rendered differentially as time since this timestamp
	ColumnType_CT_DATE ColumnType = 5
)

// Enum value maps for ColumnType.
var (
	ColumnType_name = map[int32]string{
		0: "CT_NONE",
		1: "CT_INTEGER",
		2: "CT_NUMBER",
		3: "CT_STRING",
		4: "CT_BOOLEAN",
		5: "CT_DATE",
	}
	ColumnType_value = map[string]int32{
		"CT_NONE":    0,
		"CT_INTEGER": 1,
		"CT_NUMBER":  2,
		"CT_STRING":  3,
		"CT_BOOLEAN": 4,
		"CT_DATE":    5,
	}
)

func (x ColumnType) Enum() *ColumnType {
	p := new(ColumnType)
	*p = x
	return p
}

func (x ColumnType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ColumnType) Descriptor() protoreflect.EnumDescriptor {
	return file_crd_proto_enumTypes[0].Descriptor()
}

func (ColumnType) Type() protoreflect.EnumType {
	return &file_crd_proto_enumTypes[0]
}

func (x ColumnType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ColumnType.Descriptor instead.
func (ColumnType) EnumDescriptor() ([]byte, []int) {
	return file_crd_proto_rawDescGZIP(), []int{0}
}

type ColumnFormat int32

const (
	// Unspecified column format
	ColumnFormat_CF_NONE     ColumnFormat = 0
	ColumnFormat_CF_INT32    ColumnFormat = 1
	ColumnFormat_CF_INT64    ColumnFormat = 2
	ColumnFormat_CF_FLOAT    ColumnFormat = 3
	ColumnFormat_CF_DOUBLE   ColumnFormat = 4
	ColumnFormat_CF_BYTE     ColumnFormat = 5
	ColumnFormat_CF_DATE     ColumnFormat = 6
	ColumnFormat_CF_DATETIME ColumnFormat = 7
	ColumnFormat_CF_PASSWORD ColumnFormat = 8
)

// Enum value maps for ColumnFormat.
var (
	ColumnFormat_name = map[int32]string{
		0: "CF_NONE",
		1: "CF_INT32",
		2: "CF_INT64",
		3: "CF_FLOAT",
		4: "CF_DOUBLE",
		5: "CF_BYTE",
		6: "CF_DATE",
		7: "CF_DATETIME",
		8: "CF_PASSWORD",
	}
	ColumnFormat_value = map[string]int32{
		"CF_NONE":     0,
		"CF_INT32":    1,
		"CF_INT64":    2,
		"CF_FLOAT":    3,
		"CF_DOUBLE":   4,
		"CF_BYTE":     5,
		"CF_DATE":     6,
		"CF_DATETIME": 7,
		"CF_PASSWORD": 8,
	}
)

func (x ColumnFormat) Enum() *ColumnFormat {
	p := new(ColumnFormat)
	*p = x
	return p
}

func (x ColumnFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ColumnFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_crd_proto_enumTypes[1].Descriptor()
}

func (ColumnFormat) Type() protoreflect.EnumType {
	return &file_crd_proto_enumTypes[1]
}

func (x ColumnFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ColumnFormat.Descriptor instead.
func (ColumnFormat) EnumDescriptor() ([]byte, []int) {
	return file_crd_proto_rawDescGZIP(), []int{1}
}

type PrinterColumn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type        ColumnType   `protobuf:"varint,2,opt,name=type,proto3,enum=protoc_gen_crd.ColumnType" json:"type,omitempty"`
	Format      ColumnFormat `protobuf:"varint,3,opt,name=format,proto3,enum=protoc_gen_crd.ColumnFormat" json:"format,omitempty"`
	Description string       `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	JsonPath    string       `protobuf:"bytes,5,opt,name=json_path,json=jsonPath,proto3" json:"json_path,omitempty"`
	Priority    int32        `protobuf:"varint,6,opt,name=priority,proto3" json:"priority,omitempty"`
}

func (x *PrinterColumn) Reset() {
	*x = PrinterColumn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrinterColumn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrinterColumn) ProtoMessage() {}

func (x *PrinterColumn) ProtoReflect() protoreflect.Message {
	mi := &file_crd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrinterColumn.ProtoReflect.Descriptor instead.
func (*PrinterColumn) Descriptor() ([]byte, []int) {
	return file_crd_proto_rawDescGZIP(), []int{0}
}

func (x *PrinterColumn) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PrinterColumn) GetType() ColumnType {
	if x != nil {
		return x.Type
	}
	return ColumnType_CT_NONE
}

func (x *PrinterColumn) GetFormat() ColumnFormat {
	if x != nil {
		return x.Format
	}
	return ColumnFormat_CF_NONE
}

func (x *PrinterColumn) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PrinterColumn) GetJsonPath() string {
	if x != nil {
		return x.JsonPath
	}
	return ""
}

func (x *PrinterColumn) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

type K8SCRD struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiGroup   string   `protobuf:"bytes,1,opt,name=api_group,json=apiGroup,proto3" json:"api_group,omitempty"`
	Kind       string   `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Singular   string   `protobuf:"bytes,3,opt,name=singular,proto3" json:"singular,omitempty"`
	Plural     string   `protobuf:"bytes,4,opt,name=plural,proto3" json:"plural,omitempty"`
	ShortNames []string `protobuf:"bytes,5,rep,name=short_names,json=shortNames,proto3" json:"short_names,omitempty"`
	// list of grouped resources the custom resource belongs to
	Categories []string `protobuf:"bytes,6,rep,name=categories,proto3" json:"categories,omitempty"`
	// additional columns available in kubectl get
	AdditionalColumns []*PrinterColumn `protobuf:"bytes,7,rep,name=additional_columns,json=additionalColumns,proto3" json:"additional_columns,omitempty"`
}

func (x *K8SCRD) Reset() {
	*x = K8SCRD{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *K8SCRD) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*K8SCRD) ProtoMessage() {}

func (x *K8SCRD) ProtoReflect() protoreflect.Message {
	mi := &file_crd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use K8SCRD.ProtoReflect.Descriptor instead.
func (*K8SCRD) Descriptor() ([]byte, []int) {
	return file_crd_proto_rawDescGZIP(), []int{1}
}

func (x *K8SCRD) GetApiGroup() string {
	if x != nil {
		return x.ApiGroup
	}
	return ""
}

func (x *K8SCRD) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *K8SCRD) GetSingular() string {
	if x != nil {
		return x.Singular
	}
	return ""
}

func (x *K8SCRD) GetPlural() string {
	if x != nil {
		return x.Plural
	}
	return ""
}

func (x *K8SCRD) GetShortNames() []string {
	if x != nil {
		return x.ShortNames
	}
	return nil
}

func (x *K8SCRD) GetCategories() []string {
	if x != nil {
		return x.Categories
	}
	return nil
}

func (x *K8SCRD) GetAdditionalColumns() []*PrinterColumn {
	if x != nil {
		return x.AdditionalColumns
	}
	return nil
}

type K8SPatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MergeKey      string `protobuf:"bytes,1,opt,name=merge_key,json=mergeKey,proto3" json:"merge_key,omitempty"`
	MergeStrategy string `protobuf:"bytes,2,opt,name=merge_strategy,json=mergeStrategy,proto3" json:"merge_strategy,omitempty"`
}

func (x *K8SPatch) Reset() {
	*x = K8SPatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_crd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *K8SPatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*K8SPatch) ProtoMessage() {}

func (x *K8SPatch) ProtoReflect() protoreflect.Message {
	mi := &file_crd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use K8SPatch.ProtoReflect.Descriptor instead.
func (*K8SPatch) Descriptor() ([]byte, []int) {
	return file_crd_proto_rawDescGZIP(), []int{2}
}

func (x *K8SPatch) GetMergeKey() string {
	if x != nil {
		return x.MergeKey
	}
	return ""
}

func (x *K8SPatch) GetMergeStrategy() string {
	if x != nil {
		return x.MergeStrategy
	}
	return ""
}

var file_crd_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*K8SCRD)(nil),
		Field:         73394821,
		Name:          "protoc_gen_crd.k8s_crd",
		Tag:           "bytes,73394821,opt,name=k8s_crd",
		Filename:      "crd.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*K8SPatch)(nil),
		Field:         73394822,
		Name:          "protoc_gen_crd.k8s_patch",
		Tag:           "bytes,73394822,opt,name=k8s_patch",
		Filename:      "crd.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional protoc_gen_crd.K8sCRD k8s_crd = 73394821;
	E_K8SCrd = &file_crd_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional protoc_gen_crd.K8sPatch k8s_patch = 73394822;
	E_K8SPatch = &file_crd_proto_extTypes[1]
)

var File_crd_proto protoreflect.FileDescriptor

var file_crd_proto_rawDesc = []byte{
	0x0a, 0x09, 0x63, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x1a, 0x20, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe4, 0x01,
	0x0a, 0x0d, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63,
	0x72, 0x64, 0x2e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x34, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e,
	0x5f, 0x63, 0x72, 0x64, 0x2e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x6a,
	0x73, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6a, 0x73, 0x6f, 0x6e, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x22, 0xfc, 0x01, 0x0a, 0x06, 0x4b, 0x38, 0x73, 0x43, 0x52, 0x44, 0x12,
	0x1b, 0x0a, 0x09, 0x61, 0x70, 0x69, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x61, 0x70, 0x69, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x12, 0x0a, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69, 0x6e, 0x67, 0x75, 0x6c, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x69, 0x6e, 0x67, 0x75, 0x6c, 0x61, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6c,
	0x75, 0x72, 0x61, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x69, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x4c, 0x0a, 0x12, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63,
	0x72, 0x64, 0x2e, 0x50, 0x72, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e,
	0x52, 0x11, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x43, 0x6f, 0x6c, 0x75,
	0x6d, 0x6e, 0x73, 0x22, 0x4e, 0x0a, 0x08, 0x4b, 0x38, 0x73, 0x50, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x0e,
	0x6d, 0x65, 0x72, 0x67, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x53, 0x74, 0x72, 0x61, 0x74,
	0x65, 0x67, 0x79, 0x2a, 0x64, 0x0a, 0x0a, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x54, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0e,
	0x0a, 0x0a, 0x43, 0x54, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x43, 0x54, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x10, 0x02, 0x12, 0x0d, 0x0a,
	0x09, 0x43, 0x54, 0x5f, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a,
	0x43, 0x54, 0x5f, 0x42, 0x4f, 0x4f, 0x4c, 0x45, 0x41, 0x4e, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x54, 0x5f, 0x44, 0x41, 0x54, 0x45, 0x10, 0x05, 0x2a, 0x90, 0x01, 0x0a, 0x0c, 0x43, 0x6f,
	0x6c, 0x75, 0x6d, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x46,
	0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x46, 0x5f, 0x49, 0x4e,
	0x54, 0x33, 0x32, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x46, 0x5f, 0x49, 0x4e, 0x54, 0x36,
	0x34, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x46, 0x5f, 0x46, 0x4c, 0x4f, 0x41, 0x54, 0x10,
	0x03, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x46, 0x5f, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x10, 0x04,
	0x12, 0x0b, 0x0a, 0x07, 0x43, 0x46, 0x5f, 0x42, 0x59, 0x54, 0x45, 0x10, 0x05, 0x12, 0x0b, 0x0a,
	0x07, 0x43, 0x46, 0x5f, 0x44, 0x41, 0x54, 0x45, 0x10, 0x06, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x46,
	0x5f, 0x44, 0x41, 0x54, 0x45, 0x54, 0x49, 0x4d, 0x45, 0x10, 0x07, 0x12, 0x0f, 0x0a, 0x0b, 0x43,
	0x46, 0x5f, 0x50, 0x41, 0x53, 0x53, 0x57, 0x4f, 0x52, 0x44, 0x10, 0x08, 0x3a, 0x53, 0x0a, 0x07,
	0x6b, 0x38, 0x73, 0x5f, 0x63, 0x72, 0x64, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x85, 0xd5, 0xff, 0x22, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63,
	0x72, 0x64, 0x2e, 0x4b, 0x38, 0x73, 0x43, 0x52, 0x44, 0x52, 0x06, 0x6b, 0x38, 0x73, 0x43, 0x72,
	0x64, 0x3a, 0x57, 0x0a, 0x09, 0x6b, 0x38, 0x73, 0x5f, 0x70, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x86, 0xd5,
	0xff, 0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x5f,
	0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x4b, 0x38, 0x73, 0x50, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x08, 0x6b, 0x38, 0x73, 0x50, 0x61, 0x74, 0x63, 0x68, 0x42, 0x39, 0x0a, 0x23, 0x6c, 0x69,
	0x62, 0x72, 0x61, 0x72, 0x79, 0x2e, 0x67, 0x6f, 0x2e, 0x6b, 0x38, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x63, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x42, 0x03, 0x43, 0x52, 0x44, 0x50, 0x01, 0x5a, 0x0b, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x3b, 0x63, 0x72, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_crd_proto_rawDescOnce sync.Once
	file_crd_proto_rawDescData = file_crd_proto_rawDesc
)

func file_crd_proto_rawDescGZIP() []byte {
	file_crd_proto_rawDescOnce.Do(func() {
		file_crd_proto_rawDescData = protoimpl.X.CompressGZIP(file_crd_proto_rawDescData)
	})
	return file_crd_proto_rawDescData
}

var file_crd_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_crd_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_crd_proto_goTypes = []interface{}{
	(ColumnType)(0),                     // 0: protoc_gen_crd.ColumnType
	(ColumnFormat)(0),                   // 1: protoc_gen_crd.ColumnFormat
	(*PrinterColumn)(nil),               // 2: protoc_gen_crd.PrinterColumn
	(*K8SCRD)(nil),                      // 3: protoc_gen_crd.K8sCRD
	(*K8SPatch)(nil),                    // 4: protoc_gen_crd.K8sPatch
	(*descriptorpb.MessageOptions)(nil), // 5: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 6: google.protobuf.FieldOptions
}
var file_crd_proto_depIdxs = []int32{
	0, // 0: protoc_gen_crd.PrinterColumn.type:type_name -> protoc_gen_crd.ColumnType
	1, // 1: protoc_gen_crd.PrinterColumn.format:type_name -> protoc_gen_crd.ColumnFormat
	2, // 2: protoc_gen_crd.K8sCRD.additional_columns:type_name -> protoc_gen_crd.PrinterColumn
	5, // 3: protoc_gen_crd.k8s_crd:extendee -> google.protobuf.MessageOptions
	6, // 4: protoc_gen_crd.k8s_patch:extendee -> google.protobuf.FieldOptions
	3, // 5: protoc_gen_crd.k8s_crd:type_name -> protoc_gen_crd.K8sCRD
	4, // 6: protoc_gen_crd.k8s_patch:type_name -> protoc_gen_crd.K8sPatch
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	5, // [5:7] is the sub-list for extension type_name
	3, // [3:5] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_crd_proto_init() }
func file_crd_proto_init() {
	if File_crd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_crd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrinterColumn); i {
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
		file_crd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*K8SCRD); i {
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
		file_crd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*K8SPatch); i {
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
			RawDescriptor: file_crd_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_crd_proto_goTypes,
		DependencyIndexes: file_crd_proto_depIdxs,
		EnumInfos:         file_crd_proto_enumTypes,
		MessageInfos:      file_crd_proto_msgTypes,
		ExtensionInfos:    file_crd_proto_extTypes,
	}.Build()
	File_crd_proto = out.File
	file_crd_proto_rawDesc = nil
	file_crd_proto_goTypes = nil
	file_crd_proto_depIdxs = nil
}
