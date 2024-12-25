package testdata

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CpuStat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number int32  `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	State  string `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *CpuStat) Reset() {
	*x = CpuStat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_cpu_stat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CpuStat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CpuStat) ProtoMessage() {}

func (x *CpuStat) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_cpu_stat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CpuStat.ProtoReflect.Descriptor instead.
func (*CpuStat) Descriptor() ([]byte, []int) {
	return file_testdata_cpu_stat_proto_rawDescGZIP(), []int{0}
}

func (x *CpuStat) GetNumber() int32 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *CpuStat) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

var File_testdata_cpu_stat_proto protoreflect.FileDescriptor

var file_testdata_cpu_stat_proto_rawDesc = []byte{
	0x0a, 0x17, 0x66, 0x75, 0x7a, 0x7a, 0x73, 0x74, 0x61, 0x6e, 0x2f, 0x63, 0x70, 0x75, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x66, 0x75, 0x7a, 0x7a, 0x73,
	0x74, 0x61, 0x6e, 0x22, 0x37, 0x0a, 0x07, 0x43, 0x70, 0x75, 0x53, 0x74, 0x61, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0x5e, 0x0a, 0x15,
	0x6f, 0x72, 0x67, 0x2e, 0x66, 0x75, 0x7a, 0x7a, 0x73, 0x74, 0x61, 0x6e, 0x2e, 0x66, 0x75, 0x7a,
	0x7a, 0x73, 0x74, 0x61, 0x6e, 0x42, 0x0c, 0x43, 0x70, 0x75, 0x53, 0x74, 0x61, 0x74, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x66, 0x75, 0x7a, 0x7a, 0x73, 0x74, 0x61, 0x6e, 0x2f, 0x66, 0x75, 0x7a, 0x7a, 0x73,
	0x74, 0x61, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x66, 0x75, 0x7a, 0x7a, 0x73,
	0x74, 0x61, 0x6e, 0x3b, 0x66, 0x75, 0x7a, 0x7a, 0x73, 0x74, 0x61, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_testdata_cpu_stat_proto_rawDescOnce sync.Once
	file_testdata_cpu_stat_proto_rawDescData = file_testdata_cpu_stat_proto_rawDesc
)

func file_testdata_cpu_stat_proto_rawDescGZIP() []byte {
	file_testdata_cpu_stat_proto_rawDescOnce.Do(func() {
		file_testdata_cpu_stat_proto_rawDescData = protoimpl.X.CompressGZIP(file_testdata_cpu_stat_proto_rawDescData)
	})
	return file_testdata_cpu_stat_proto_rawDescData
}

var file_testdata_cpu_stat_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_testdata_cpu_stat_proto_goTypes = []interface{}{
	(*CpuStat)(nil), // 0: testdata.CpuStat
}
var file_testdata_cpu_stat_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_testdata_cpu_stat_proto_init() }
func file_testdata_cpu_stat_proto_init() {
	if File_testdata_cpu_stat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_testdata_cpu_stat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CpuStat); i {
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
			RawDescriptor: file_testdata_cpu_stat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_testdata_cpu_stat_proto_goTypes,
		DependencyIndexes: file_testdata_cpu_stat_proto_depIdxs,
		MessageInfos:      file_testdata_cpu_stat_proto_msgTypes,
	}.Build()
	File_testdata_cpu_stat_proto = out.File
	file_testdata_cpu_stat_proto_rawDesc = nil
	file_testdata_cpu_stat_proto_goTypes = nil
	file_testdata_cpu_stat_proto_depIdxs = nil
}
