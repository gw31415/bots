// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: cmdin.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

//データタイプ
type InputMedia_MediaType int32

const (
	InputMedia_UNKNOWN InputMedia_MediaType = 0
	InputMedia_EXTEND  InputMedia_MediaType = 1 //データタイプの決定を先伸ばし, 直後のデータと連結する
	InputMedia_UTF8    InputMedia_MediaType = 2 //文字列
	InputMedia_FILE    InputMedia_MediaType = 3 //ファイル
)

// Enum value maps for InputMedia_MediaType.
var (
	InputMedia_MediaType_name = map[int32]string{
		0: "UNKNOWN",
		1: "EXTEND",
		2: "UTF8",
		3: "FILE",
	}
	InputMedia_MediaType_value = map[string]int32{
		"UNKNOWN": 0,
		"EXTEND":  1,
		"UTF8":    2,
		"FILE":    3,
	}
)

func (x InputMedia_MediaType) Enum() *InputMedia_MediaType {
	p := new(InputMedia_MediaType)
	*p = x
	return p
}

func (x InputMedia_MediaType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InputMedia_MediaType) Descriptor() protoreflect.EnumDescriptor {
	return file_cmdin_proto_enumTypes[0].Descriptor()
}

func (InputMedia_MediaType) Type() protoreflect.EnumType {
	return &file_cmdin_proto_enumTypes[0]
}

func (x InputMedia_MediaType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InputMedia_MediaType.Descriptor instead.
func (InputMedia_MediaType) EnumDescriptor() ([]byte, []int) {
	return file_cmdin_proto_rawDescGZIP(), []int{1, 0}
}

type Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Media []*InputMedia `protobuf:"bytes,1,rep,name=media,proto3" json:"media,omitempty"`
	//prefixの文字列, help用
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"`
}

func (x *Input) Reset() {
	*x = Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmdin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Input) ProtoMessage() {}

func (x *Input) ProtoReflect() protoreflect.Message {
	mi := &file_cmdin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Input.ProtoReflect.Descriptor instead.
func (*Input) Descriptor() ([]byte, []int) {
	return file_cmdin_proto_rawDescGZIP(), []int{0}
}

func (x *Input) GetMedia() []*InputMedia {
	if x != nil {
		return x.Media
	}
	return nil
}

func (x *Input) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

type InputMedia struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type InputMedia_MediaType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.InputMedia_MediaType" json:"type,omitempty"`
	//データ本体
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	//ファイル名(拡張子あり), ファイルタイプは拡張子で判断
	Filename string `protobuf:"bytes,3,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *InputMedia) Reset() {
	*x = InputMedia{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmdin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputMedia) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputMedia) ProtoMessage() {}

func (x *InputMedia) ProtoReflect() protoreflect.Message {
	mi := &file_cmdin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputMedia.ProtoReflect.Descriptor instead.
func (*InputMedia) Descriptor() ([]byte, []int) {
	return file_cmdin_proto_rawDescGZIP(), []int{1}
}

func (x *InputMedia) GetType() InputMedia_MediaType {
	if x != nil {
		return x.Type
	}
	return InputMedia_UNKNOWN
}

func (x *InputMedia) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *InputMedia) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

var File_cmdin_proto protoreflect.FileDescriptor

var file_cmdin_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x6d, 0x64, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a, 0x05, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x27, 0x0a,
	0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52,
	0x05, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x22, 0xad,
	0x01, 0x0a, 0x0a, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x2f, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x4d,
	0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x38,
	0x0a, 0x09, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x45, 0x58, 0x54, 0x45,
	0x4e, 0x44, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x55, 0x54, 0x46, 0x38, 0x10, 0x02, 0x12, 0x08,
	0x0a, 0x04, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x04, 0x10, 0x08, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cmdin_proto_rawDescOnce sync.Once
	file_cmdin_proto_rawDescData = file_cmdin_proto_rawDesc
)

func file_cmdin_proto_rawDescGZIP() []byte {
	file_cmdin_proto_rawDescOnce.Do(func() {
		file_cmdin_proto_rawDescData = protoimpl.X.CompressGZIP(file_cmdin_proto_rawDescData)
	})
	return file_cmdin_proto_rawDescData
}

var file_cmdin_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_cmdin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_cmdin_proto_goTypes = []interface{}{
	(InputMedia_MediaType)(0), // 0: proto.InputMedia.MediaType
	(*Input)(nil),             // 1: proto.Input
	(*InputMedia)(nil),        // 2: proto.InputMedia
}
var file_cmdin_proto_depIdxs = []int32{
	2, // 0: proto.Input.media:type_name -> proto.InputMedia
	0, // 1: proto.InputMedia.type:type_name -> proto.InputMedia.MediaType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_cmdin_proto_init() }
func file_cmdin_proto_init() {
	if File_cmdin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cmdin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Input); i {
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
		file_cmdin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputMedia); i {
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
			RawDescriptor: file_cmdin_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cmdin_proto_goTypes,
		DependencyIndexes: file_cmdin_proto_depIdxs,
		EnumInfos:         file_cmdin_proto_enumTypes,
		MessageInfos:      file_cmdin_proto_msgTypes,
	}.Build()
	File_cmdin_proto = out.File
	file_cmdin_proto_rawDesc = nil
	file_cmdin_proto_goTypes = nil
	file_cmdin_proto_depIdxs = nil
}
