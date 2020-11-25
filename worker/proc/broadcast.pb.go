// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.14.0
// source: broadcast.proto

package proc

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

type Broadcast struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Head *Head `protobuf:"bytes,1,req,name=head" json:"head,omitempty"`
	Body *Body `protobuf:"bytes,2,opt,name=body" json:"body,omitempty"`
}

func (x *Broadcast) Reset() {
	*x = Broadcast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broadcast_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Broadcast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Broadcast) ProtoMessage() {}

func (x *Broadcast) ProtoReflect() protoreflect.Message {
	mi := &file_broadcast_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Broadcast.ProtoReflect.Descriptor instead.
func (*Broadcast) Descriptor() ([]byte, []int) {
	return file_broadcast_proto_rawDescGZIP(), []int{0}
}

func (x *Broadcast) GetHead() *Head {
	if x != nil {
		return x.Head
	}
	return nil
}

func (x *Broadcast) GetBody() *Body {
	if x != nil {
		return x.Body
	}
	return nil
}

type Body struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//服务器主机
	Server *string `protobuf:"bytes,1,req,name=server" json:"server,omitempty"`
	//端口号4字节 65536
	Port *uint32 `protobuf:"varint,2,req,name=port" json:"port,omitempty"`
}

func (x *Body) Reset() {
	*x = Body{}
	if protoimpl.UnsafeEnabled {
		mi := &file_broadcast_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Body) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Body) ProtoMessage() {}

func (x *Body) ProtoReflect() protoreflect.Message {
	mi := &file_broadcast_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Body.ProtoReflect.Descriptor instead.
func (*Body) Descriptor() ([]byte, []int) {
	return file_broadcast_proto_rawDescGZIP(), []int{1}
}

func (x *Body) GetServer() string {
	if x != nil && x.Server != nil {
		return *x.Server
	}
	return ""
}

func (x *Body) GetPort() uint32 {
	if x != nil && x.Port != nil {
		return *x.Port
	}
	return 0
}

var File_broadcast_proto protoreflect.FileDescriptor

var file_broadcast_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x70, 0x72, 0x6f, 0x63, 0x1a, 0x0a, 0x68, 0x65, 0x61, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x12, 0x1e, 0x0a, 0x04, 0x68, 0x65, 0x61, 0x64, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x52, 0x04, 0x68, 0x65, 0x61, 0x64,
	0x12, 0x1e, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x22, 0x32, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0d, 0x52, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x50, 0x00,
}

var (
	file_broadcast_proto_rawDescOnce sync.Once
	file_broadcast_proto_rawDescData = file_broadcast_proto_rawDesc
)

func file_broadcast_proto_rawDescGZIP() []byte {
	file_broadcast_proto_rawDescOnce.Do(func() {
		file_broadcast_proto_rawDescData = protoimpl.X.CompressGZIP(file_broadcast_proto_rawDescData)
	})
	return file_broadcast_proto_rawDescData
}

var file_broadcast_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_broadcast_proto_goTypes = []interface{}{
	(*Broadcast)(nil), // 0: proc.Broadcast
	(*Body)(nil),      // 1: proc.Body
	(*Head)(nil),      // 2: proc.Head
}
var file_broadcast_proto_depIdxs = []int32{
	2, // 0: proc.Broadcast.head:type_name -> proc.Head
	1, // 1: proc.Broadcast.body:type_name -> proc.Body
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_broadcast_proto_init() }
func file_broadcast_proto_init() {
	if File_broadcast_proto != nil {
		return
	}
	file_head_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_broadcast_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Broadcast); i {
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
		file_broadcast_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Body); i {
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
			RawDescriptor: file_broadcast_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_broadcast_proto_goTypes,
		DependencyIndexes: file_broadcast_proto_depIdxs,
		MessageInfos:      file_broadcast_proto_msgTypes,
	}.Build()
	File_broadcast_proto = out.File
	file_broadcast_proto_rawDesc = nil
	file_broadcast_proto_goTypes = nil
	file_broadcast_proto_depIdxs = nil
}
