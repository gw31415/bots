// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        (unknown)
// source: cmdout.proto

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
type OutputMedia_MediaType int32

const (
	OutputMedia_UNKNOWN  OutputMedia_MediaType = 0
	OutputMedia_EXTEND   OutputMedia_MediaType = 1 //データタイプの決定を先伸ばし, 直後のデータと連結する
	OutputMedia_UTF8     OutputMedia_MediaType = 2 //文字列
	OutputMedia_FILE     OutputMedia_MediaType = 3 //ファイル
	OutputMedia_FILE_URL OutputMedia_MediaType = 4 //ファイルのURL
)

// Enum value maps for OutputMedia_MediaType.
var (
	OutputMedia_MediaType_name = map[int32]string{
		0: "UNKNOWN",
		1: "EXTEND",
		2: "UTF8",
		3: "FILE",
		4: "FILE_URL",
	}
	OutputMedia_MediaType_value = map[string]int32{
		"UNKNOWN":  0,
		"EXTEND":   1,
		"UTF8":     2,
		"FILE":     3,
		"FILE_URL": 4,
	}
)

func (x OutputMedia_MediaType) Enum() *OutputMedia_MediaType {
	p := new(OutputMedia_MediaType)
	*p = x
	return p
}

func (x OutputMedia_MediaType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OutputMedia_MediaType) Descriptor() protoreflect.EnumDescriptor {
	return file_cmdout_proto_enumTypes[0].Descriptor()
}

func (OutputMedia_MediaType) Type() protoreflect.EnumType {
	return &file_cmdout_proto_enumTypes[0]
}

func (x OutputMedia_MediaType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OutputMedia_MediaType.Descriptor instead.
func (OutputMedia_MediaType) EnumDescriptor() ([]byte, []int) {
	return file_cmdout_proto_rawDescGZIP(), []int{2, 0}
}

//実際にコマンドから出力されるオブジェクトの型
type Output struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msgs []*BotMsg `protobuf:"bytes,1,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (x *Output) Reset() {
	*x = Output{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmdout_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Output) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Output) ProtoMessage() {}

func (x *Output) ProtoReflect() protoreflect.Message {
	mi := &file_cmdout_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Output.ProtoReflect.Descriptor instead.
func (*Output) Descriptor() ([]byte, []int) {
	return file_cmdout_proto_rawDescGZIP(), []int{0}
}

func (x *Output) GetMsgs() []*BotMsg {
	if x != nil {
		return x.Msgs
	}
	return nil
}

//発言ひとつ分
type BotMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Medias           []*OutputMedia `protobuf:"bytes,1,rep,name=medias,proto3" json:"medias,omitempty"`
	EmbedRecommended bool           `protobuf:"varint,2,opt,name=embed_recommended,json=embedRecommended,proto3" json:"embed_recommended,omitempty"`
	//カラーコード、最上位のバイトは無視
	Color uint32 `protobuf:"varint,3,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *BotMsg) Reset() {
	*x = BotMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmdout_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BotMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BotMsg) ProtoMessage() {}

func (x *BotMsg) ProtoReflect() protoreflect.Message {
	mi := &file_cmdout_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BotMsg.ProtoReflect.Descriptor instead.
func (*BotMsg) Descriptor() ([]byte, []int) {
	return file_cmdout_proto_rawDescGZIP(), []int{1}
}

func (x *BotMsg) GetMedias() []*OutputMedia {
	if x != nil {
		return x.Medias
	}
	return nil
}

func (x *BotMsg) GetEmbedRecommended() bool {
	if x != nil {
		return x.EmbedRecommended
	}
	return false
}

func (x *BotMsg) GetColor() uint32 {
	if x != nil {
		return x.Color
	}
	return 0
}

//データ
type OutputMedia struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type OutputMedia_MediaType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.OutputMedia_MediaType" json:"type,omitempty"`
	//データ本体
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	//ファイル名(拡張子あり), ファイルタイプは拡張子で判断
	Filename string `protobuf:"bytes,3,opt,name=filename,proto3" json:"filename,omitempty"`
	//直後のデータに同フィールドが続くかどうか
	ExtendField bool `protobuf:"varint,4,opt,name=extend_field,json=extendField,proto3" json:"extend_field,omitempty"`
	//見出しレベル
	//0は通常, 1最も強調, 2,3,...と準じる
	Level int32 `protobuf:"varint,5,opt,name=level,proto3" json:"level,omitempty"`
	//エラーコード, 0ならエラーなし
	Error uint32 `protobuf:"varint,6,opt,name=error,proto3" json:"error,omitempty"`
	//隠せるなら隠すかどうか
	Spoiled bool `protobuf:"varint,7,opt,name=spoiled,proto3" json:"spoiled,omitempty"`
}

func (x *OutputMedia) Reset() {
	*x = OutputMedia{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmdout_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutputMedia) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutputMedia) ProtoMessage() {}

func (x *OutputMedia) ProtoReflect() protoreflect.Message {
	mi := &file_cmdout_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutputMedia.ProtoReflect.Descriptor instead.
func (*OutputMedia) Descriptor() ([]byte, []int) {
	return file_cmdout_proto_rawDescGZIP(), []int{2}
}

func (x *OutputMedia) GetType() OutputMedia_MediaType {
	if x != nil {
		return x.Type
	}
	return OutputMedia_UNKNOWN
}

func (x *OutputMedia) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *OutputMedia) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *OutputMedia) GetExtendField() bool {
	if x != nil {
		return x.ExtendField
	}
	return false
}

func (x *OutputMedia) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *OutputMedia) GetError() uint32 {
	if x != nil {
		return x.Error
	}
	return 0
}

func (x *OutputMedia) GetSpoiled() bool {
	if x != nil {
		return x.Spoiled
	}
	return false
}

var File_cmdout_proto protoreflect.FileDescriptor

var file_cmdout_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6d, 0x64, 0x6f, 0x75, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a, 0x06, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12,
	0x21, 0x0a, 0x04, 0x6d, 0x73, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x74, 0x4d, 0x73, 0x67, 0x52, 0x04, 0x6d, 0x73,
	0x67, 0x73, 0x22, 0x77, 0x0a, 0x06, 0x42, 0x6f, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x2a, 0x0a, 0x06,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x64, 0x69, 0x61,
	0x52, 0x06, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x73, 0x12, 0x2b, 0x0a, 0x11, 0x65, 0x6d, 0x62, 0x65,
	0x64, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x10, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x64, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x22, 0xa0, 0x02, 0x0a, 0x0b,
	0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x30, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x2e, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x70, 0x6f, 0x69, 0x6c, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x70, 0x6f, 0x69, 0x6c, 0x65, 0x64, 0x22, 0x46, 0x0a, 0x09, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x45, 0x58, 0x54, 0x45, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04,
	0x55, 0x54, 0x46, 0x38, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x03,
	0x12, 0x0c, 0x0a, 0x08, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x55, 0x52, 0x4c, 0x10, 0x04, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cmdout_proto_rawDescOnce sync.Once
	file_cmdout_proto_rawDescData = file_cmdout_proto_rawDesc
)

func file_cmdout_proto_rawDescGZIP() []byte {
	file_cmdout_proto_rawDescOnce.Do(func() {
		file_cmdout_proto_rawDescData = protoimpl.X.CompressGZIP(file_cmdout_proto_rawDescData)
	})
	return file_cmdout_proto_rawDescData
}

var file_cmdout_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_cmdout_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_cmdout_proto_goTypes = []interface{}{
	(OutputMedia_MediaType)(0), // 0: proto.OutputMedia.MediaType
	(*Output)(nil),             // 1: proto.Output
	(*BotMsg)(nil),             // 2: proto.BotMsg
	(*OutputMedia)(nil),        // 3: proto.OutputMedia
}
var file_cmdout_proto_depIdxs = []int32{
	2, // 0: proto.Output.msgs:type_name -> proto.BotMsg
	3, // 1: proto.BotMsg.medias:type_name -> proto.OutputMedia
	0, // 2: proto.OutputMedia.type:type_name -> proto.OutputMedia.MediaType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_cmdout_proto_init() }
func file_cmdout_proto_init() {
	if File_cmdout_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cmdout_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Output); i {
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
		file_cmdout_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BotMsg); i {
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
		file_cmdout_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutputMedia); i {
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
			RawDescriptor: file_cmdout_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cmdout_proto_goTypes,
		DependencyIndexes: file_cmdout_proto_depIdxs,
		EnumInfos:         file_cmdout_proto_enumTypes,
		MessageInfos:      file_cmdout_proto_msgTypes,
	}.Build()
	File_cmdout_proto = out.File
	file_cmdout_proto_rawDesc = nil
	file_cmdout_proto_goTypes = nil
	file_cmdout_proto_depIdxs = nil
}
