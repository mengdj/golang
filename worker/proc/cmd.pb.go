// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.14.0
// source: cmd.proto

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

//支持ping 截图
type ContentType int32

const (
	//option allow_alias=true;
	ContentType_PING      ContentType = 1
	ContentType_CAPTURE   ContentType = 2
	ContentType_CHAT      ContentType = 3
	ContentType_BROADCAST ContentType = 4
)

// Enum value maps for ContentType.
var (
	ContentType_name = map[int32]string{
		1: "PING",
		2: "CAPTURE",
		3: "CHAT",
		4: "BROADCAST",
	}
	ContentType_value = map[string]int32{
		"PING":      1,
		"CAPTURE":   2,
		"CHAT":      3,
		"BROADCAST": 4,
	}
)

func (x ContentType) Enum() *ContentType {
	p := new(ContentType)
	*p = x
	return p
}

func (x ContentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContentType) Descriptor() protoreflect.EnumDescriptor {
	return file_cmd_proto_enumTypes[0].Descriptor()
}

func (ContentType) Type() protoreflect.EnumType {
	return &file_cmd_proto_enumTypes[0]
}

func (x ContentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *ContentType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = ContentType(num)
	return nil
}

// Deprecated: Use ContentType.Descriptor instead.
func (ContentType) EnumDescriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{0}
}

//消息来源
type Source int32

const (
	Source_SERVER Source = 1
	Source_CLIENT Source = 2
)

// Enum value maps for Source.
var (
	Source_name = map[int32]string{
		1: "SERVER",
		2: "CLIENT",
	}
	Source_value = map[string]int32{
		"SERVER": 1,
		"CLIENT": 2,
	}
)

func (x Source) Enum() *Source {
	p := new(Source)
	*p = x
	return p
}

func (x Source) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Source) Descriptor() protoreflect.EnumDescriptor {
	return file_cmd_proto_enumTypes[1].Descriptor()
}

func (Source) Type() protoreflect.EnumType {
	return &file_cmd_proto_enumTypes[1]
}

func (x Source) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *Source) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = Source(num)
	return nil
}

// Deprecated: Use Source.Descriptor instead.
func (Source) EnumDescriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{1}
}

type Cmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Head    *Head    `protobuf:"bytes,1,req,name=head" json:"head,omitempty"`
	Content *Content `protobuf:"bytes,2,req,name=content" json:"content,omitempty"`
}

func (x *Cmd) Reset() {
	*x = Cmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cmd) ProtoMessage() {}

func (x *Cmd) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cmd.ProtoReflect.Descriptor instead.
func (*Cmd) Descriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{0}
}

func (x *Cmd) GetHead() *Head {
	if x != nil {
		return x.Head
	}
	return nil
}

func (x *Cmd) GetContent() *Content {
	if x != nil {
		return x.Content
	}
	return nil
}

type Content struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//内容类型
	Type *ContentType `protobuf:"varint,1,req,name=type,enum=proc.ContentType,def=1" json:"type,omitempty"`
	//内容来源
	Source *Source `protobuf:"varint,2,req,name=source,enum=proc.Source,def=1" json:"source,omitempty"`
	//内容体
	//
	// Types that are assignable to Param:
	//	*Content_Ping
	//	*Content_Capture
	//	*Content_Chat
	Param isContent_Param `protobuf_oneof:"param"`
	//附加数据（可选）
	Extra []byte `protobuf:"bytes,8,opt,name=extra" json:"extra,omitempty"`
}

// Default values for Content fields.
const (
	Default_Content_Type   = ContentType_PING
	Default_Content_Source = Source_SERVER
)

func (x *Content) Reset() {
	*x = Content{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content) ProtoMessage() {}

func (x *Content) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content.ProtoReflect.Descriptor instead.
func (*Content) Descriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{1}
}

func (x *Content) GetType() ContentType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return Default_Content_Type
}

func (x *Content) GetSource() Source {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return Default_Content_Source
}

func (m *Content) GetParam() isContent_Param {
	if m != nil {
		return m.Param
	}
	return nil
}

func (x *Content) GetPing() *Ping {
	if x, ok := x.GetParam().(*Content_Ping); ok {
		return x.Ping
	}
	return nil
}

func (x *Content) GetCapture() *Capture {
	if x, ok := x.GetParam().(*Content_Capture); ok {
		return x.Capture
	}
	return nil
}

func (x *Content) GetChat() *Chat {
	if x, ok := x.GetParam().(*Content_Chat); ok {
		return x.Chat
	}
	return nil
}

func (x *Content) GetExtra() []byte {
	if x != nil {
		return x.Extra
	}
	return nil
}

type isContent_Param interface {
	isContent_Param()
}

type Content_Ping struct {
	Ping *Ping `protobuf:"bytes,3,opt,name=ping,oneof"`
}

type Content_Capture struct {
	Capture *Capture `protobuf:"bytes,4,opt,name=capture,oneof"`
}

type Content_Chat struct {
	Chat *Chat `protobuf:"bytes,5,opt,name=chat,oneof"`
}

func (*Content_Ping) isContent_Param() {}

func (*Content_Capture) isContent_Param() {}

func (*Content_Chat) isContent_Param() {}

type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//主机名
	Name *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	//架构
	Arch *string `protobuf:"bytes,2,opt,name=arch" json:"arch,omitempty"`
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{2}
}

func (x *Ping) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Ping) GetArch() string {
	if x != nil && x.Arch != nil {
		return *x.Arch
	}
	return ""
}

//截图
type Capture struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//数据
	Data []byte `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
	//当前数据包编号
	Id *uint32 `protobuf:"varint,2,opt,name=id,def=0" json:"id,omitempty"`
	//数据包切分序列号
	Seq *uint32 `protobuf:"varint,3,opt,name=seq,def=0" json:"seq,omitempty"`
	//确认（传入确认号，确认号为服务器反馈收到的序号 id:seq）
	Ack *uint32 `protobuf:"varint,4,opt,name=ack,def=0" json:"ack,omitempty"`
	//是否有更多数据(分段，大数据)
	More *bool `protobuf:"varint,5,opt,name=more,def=0" json:"more,omitempty"`
	//数据原始大小（整个包）
	Size *uint32 `protobuf:"varint,6,opt,name=size,def=0" json:"size,omitempty"`
	//是否是压缩数据 snappy
	Compress *bool `protobuf:"varint,7,opt,name=compress,def=0" json:"compress,omitempty"`
}

// Default values for Capture fields.
const (
	Default_Capture_Id       = uint32(0)
	Default_Capture_Seq      = uint32(0)
	Default_Capture_Ack      = uint32(0)
	Default_Capture_More     = bool(false)
	Default_Capture_Size     = uint32(0)
	Default_Capture_Compress = bool(false)
)

func (x *Capture) Reset() {
	*x = Capture{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Capture) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Capture) ProtoMessage() {}

func (x *Capture) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Capture.ProtoReflect.Descriptor instead.
func (*Capture) Descriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{3}
}

func (x *Capture) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Capture) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return Default_Capture_Id
}

func (x *Capture) GetSeq() uint32 {
	if x != nil && x.Seq != nil {
		return *x.Seq
	}
	return Default_Capture_Seq
}

func (x *Capture) GetAck() uint32 {
	if x != nil && x.Ack != nil {
		return *x.Ack
	}
	return Default_Capture_Ack
}

func (x *Capture) GetMore() bool {
	if x != nil && x.More != nil {
		return *x.More
	}
	return Default_Capture_More
}

func (x *Capture) GetSize() uint32 {
	if x != nil && x.Size != nil {
		return *x.Size
	}
	return Default_Capture_Size
}

func (x *Capture) GetCompress() bool {
	if x != nil && x.Compress != nil {
		return *x.Compress
	}
	return Default_Capture_Compress
}

type Chat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//聊天
	Text *string `protobuf:"bytes,1,req,name=text" json:"text,omitempty"`
}

func (x *Chat) Reset() {
	*x = Chat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmd_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chat) ProtoMessage() {}

func (x *Chat) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chat.ProtoReflect.Descriptor instead.
func (*Chat) Descriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{4}
}

func (x *Chat) GetText() string {
	if x != nil && x.Text != nil {
		return *x.Text
	}
	return ""
}

var File_cmd_proto protoreflect.FileDescriptor

var file_cmd_proto_rawDesc = []byte{
	0x0a, 0x09, 0x63, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x70, 0x72, 0x6f,
	0x63, 0x1a, 0x0a, 0x68, 0x65, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4e, 0x0a,
	0x03, 0x43, 0x6d, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x68, 0x65, 0x61, 0x64, 0x18, 0x01, 0x20, 0x02,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x52, 0x04,
	0x68, 0x65, 0x61, 0x64, 0x12, 0x27, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x02, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0xf2, 0x01,
	0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x43,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x3a, 0x04, 0x50, 0x49, 0x4e, 0x47,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x18, 0x02, 0x20, 0x02, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x3a, 0x06, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x52, 0x06, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x70, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x48, 0x00,
	0x52, 0x04, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x29, 0x0a, 0x07, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x43,
	0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x48, 0x00, 0x52, 0x07, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x20, 0x0a, 0x04, 0x63, 0x68, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x48, 0x00, 0x52, 0x04, 0x63,
	0x68, 0x61, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78, 0x74, 0x72, 0x61, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x65, 0x78, 0x74, 0x72, 0x61, 0x42, 0x07, 0x0a, 0x05, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x22, 0x2e, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x61, 0x72, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72,
	0x63, 0x68, 0x22, 0xaf, 0x01, 0x0a, 0x07, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x11, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01,
	0x30, 0x52, 0x02, 0x69, 0x64, 0x12, 0x13, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x03, 0x73, 0x65, 0x71, 0x12, 0x13, 0x0a, 0x03, 0x61, 0x63,
	0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x03, 0x61, 0x63, 0x6b, 0x12,
	0x19, 0x0a, 0x04, 0x6d, 0x6f, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x05, 0x66,
	0x61, 0x6c, 0x73, 0x65, 0x52, 0x04, 0x6d, 0x6f, 0x72, 0x65, 0x12, 0x15, 0x0a, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x21, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x08, 0x3a, 0x05, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x22, 0x1a, 0x0a, 0x04, 0x43, 0x68, 0x61, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x2a, 0x3d, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x41, 0x50,
	0x54, 0x55, 0x52, 0x45, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x43, 0x48, 0x41, 0x54, 0x10, 0x03,
	0x12, 0x0d, 0x0a, 0x09, 0x42, 0x52, 0x4f, 0x41, 0x44, 0x43, 0x41, 0x53, 0x54, 0x10, 0x04, 0x2a,
	0x20, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x45, 0x52,
	0x56, 0x45, 0x52, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x4c, 0x49, 0x45, 0x4e, 0x54, 0x10,
	0x02, 0x50, 0x00,
}

var (
	file_cmd_proto_rawDescOnce sync.Once
	file_cmd_proto_rawDescData = file_cmd_proto_rawDesc
)

func file_cmd_proto_rawDescGZIP() []byte {
	file_cmd_proto_rawDescOnce.Do(func() {
		file_cmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_cmd_proto_rawDescData)
	})
	return file_cmd_proto_rawDescData
}

var file_cmd_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_cmd_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_cmd_proto_goTypes = []interface{}{
	(ContentType)(0), // 0: proc.ContentType
	(Source)(0),      // 1: proc.Source
	(*Cmd)(nil),      // 2: proc.Cmd
	(*Content)(nil),  // 3: proc.Content
	(*Ping)(nil),     // 4: proc.Ping
	(*Capture)(nil),  // 5: proc.Capture
	(*Chat)(nil),     // 6: proc.Chat
	(*Head)(nil),     // 7: proc.Head
}
var file_cmd_proto_depIdxs = []int32{
	7, // 0: proc.Cmd.head:type_name -> proc.Head
	3, // 1: proc.Cmd.content:type_name -> proc.Content
	0, // 2: proc.Content.type:type_name -> proc.ContentType
	1, // 3: proc.Content.source:type_name -> proc.Source
	4, // 4: proc.Content.ping:type_name -> proc.Ping
	5, // 5: proc.Content.capture:type_name -> proc.Capture
	6, // 6: proc.Content.chat:type_name -> proc.Chat
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_cmd_proto_init() }
func file_cmd_proto_init() {
	if File_cmd_proto != nil {
		return
	}
	file_head_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_cmd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cmd); i {
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
		file_cmd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content); i {
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
		file_cmd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
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
		file_cmd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Capture); i {
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
		file_cmd_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chat); i {
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
	file_cmd_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Content_Ping)(nil),
		(*Content_Capture)(nil),
		(*Content_Chat)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cmd_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cmd_proto_goTypes,
		DependencyIndexes: file_cmd_proto_depIdxs,
		EnumInfos:         file_cmd_proto_enumTypes,
		MessageInfos:      file_cmd_proto_msgTypes,
	}.Build()
	File_cmd_proto = out.File
	file_cmd_proto_rawDesc = nil
	file_cmd_proto_goTypes = nil
	file_cmd_proto_depIdxs = nil
}
