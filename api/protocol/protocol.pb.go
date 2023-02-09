// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: protocol/protocol.proto

package protocol

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

type Op int32

const (
	Op_Undefined       Op = 0
	Op_Auth            Op = 1
	Op_AuthReply       Op = 2
	Op_Heartbeat       Op = 3
	Op_HeartbeatReply  Op = 4
	Op_Disconnect      Op = 5
	Op_DisconnectReply Op = 6
	Op_SendMsg         Op = 7
	Op_SendMsgReply    Op = 8 //客户端回复消息已收到
	Op_ReceiveMsg      Op = 9
	Op_ReceiveMsgReply Op = 10
	Op_SyncMsgReq      Op = 11
	Op_SyncMsgReply    Op = 12
	Op_ProtoReady      Op = 20
	Op_ProtoFinish     Op = 21
	Op_ProtoResend     Op = 22
	Op_Raw             Op = 23
)

// Enum value maps for Op.
var (
	Op_name = map[int32]string{
		0:  "Undefined",
		1:  "Auth",
		2:  "AuthReply",
		3:  "Heartbeat",
		4:  "HeartbeatReply",
		5:  "Disconnect",
		6:  "DisconnectReply",
		7:  "SendMsg",
		8:  "SendMsgReply",
		9:  "ReceiveMsg",
		10: "ReceiveMsgReply",
		11: "SyncMsgReq",
		12: "SyncMsgReply",
		20: "ProtoReady",
		21: "ProtoFinish",
		22: "ProtoResend",
		23: "Raw",
	}
	Op_value = map[string]int32{
		"Undefined":       0,
		"Auth":            1,
		"AuthReply":       2,
		"Heartbeat":       3,
		"HeartbeatReply":  4,
		"Disconnect":      5,
		"DisconnectReply": 6,
		"SendMsg":         7,
		"SendMsgReply":    8,
		"ReceiveMsg":      9,
		"ReceiveMsgReply": 10,
		"SyncMsgReq":      11,
		"SyncMsgReply":    12,
		"ProtoReady":      20,
		"ProtoFinish":     21,
		"ProtoResend":     22,
		"Raw":             23,
	}
)

func (x Op) Enum() *Op {
	p := new(Op)
	*p = x
	return p
}

func (x Op) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Op) Descriptor() protoreflect.EnumDescriptor {
	return file_protocol_protocol_proto_enumTypes[0].Descriptor()
}

func (Op) Type() protoreflect.EnumType {
	return &file_protocol_protocol_proto_enumTypes[0]
}

func (x Op) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Op.Descriptor instead.
func (Op) EnumDescriptor() ([]byte, []int) {
	return file_protocol_protocol_proto_rawDescGZIP(), []int{0}
}

type Channel int32

const (
	Channel_Private Channel = 0
	Channel_Group   Channel = 1
)

// Enum value maps for Channel.
var (
	Channel_name = map[int32]string{
		0: "Private",
		1: "Group",
	}
	Channel_value = map[string]int32{
		"Private": 0,
		"Group":   1,
	}
)

func (x Channel) Enum() *Channel {
	p := new(Channel)
	*p = x
	return p
}

func (x Channel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Channel) Descriptor() protoreflect.EnumDescriptor {
	return file_protocol_protocol_proto_enumTypes[1].Descriptor()
}

func (Channel) Type() protoreflect.EnumType {
	return &file_protocol_protocol_proto_enumTypes[1]
}

func (x Channel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Channel.Descriptor instead.
func (Channel) EnumDescriptor() ([]byte, []int) {
	return file_protocol_protocol_proto_rawDescGZIP(), []int{1}
}

type SendMsgReplyBody_ErrorType int32

const (
	SendMsgReplyBody_IsOk               SendMsgReplyBody_ErrorType = 0
	SendMsgReplyBody_NotFriend          SendMsgReplyBody_ErrorType = 1
	SendMsgReplyBody_NotGroupMember     SendMsgReplyBody_ErrorType = 2
	SendMsgReplyBody_UnsupportedChannel SendMsgReplyBody_ErrorType = 3
)

// Enum value maps for SendMsgReplyBody_ErrorType.
var (
	SendMsgReplyBody_ErrorType_name = map[int32]string{
		0: "IsOk",
		1: "NotFriend",
		2: "NotGroupMember",
		3: "UnsupportedChannel",
	}
	SendMsgReplyBody_ErrorType_value = map[string]int32{
		"IsOk":               0,
		"NotFriend":          1,
		"NotGroupMember":     2,
		"UnsupportedChannel": 3,
	}
)

func (x SendMsgReplyBody_ErrorType) Enum() *SendMsgReplyBody_ErrorType {
	p := new(SendMsgReplyBody_ErrorType)
	*p = x
	return p
}

func (x SendMsgReplyBody_ErrorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SendMsgReplyBody_ErrorType) Descriptor() protoreflect.EnumDescriptor {
	return file_protocol_protocol_proto_enumTypes[2].Descriptor()
}

func (SendMsgReplyBody_ErrorType) Type() protoreflect.EnumType {
	return &file_protocol_protocol_proto_enumTypes[2]
}

func (x SendMsgReplyBody_ErrorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SendMsgReplyBody_ErrorType.Descriptor instead.
func (SendMsgReplyBody_ErrorType) EnumDescriptor() ([]byte, []int) {
	return file_protocol_protocol_proto_rawDescGZIP(), []int{1, 0}
}

// Proto 中 Op 为 OpAuth 时， body 必须可以反序列化为 AuthMsg
type AuthBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId string `protobuf:"bytes,1,opt,name=appId,proto3" json:"appId,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Ext   []byte `protobuf:"bytes,3,opt,name=ext,proto3" json:"ext,omitempty"` // 其它业务方可能需要的信息
}

func (x *AuthBody) Reset() {
	*x = AuthBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_protocol_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthBody) ProtoMessage() {}

func (x *AuthBody) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_protocol_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthBody.ProtoReflect.Descriptor instead.
func (*AuthBody) Descriptor() ([]byte, []int) {
	return file_protocol_protocol_proto_rawDescGZIP(), []int{0}
}

func (x *AuthBody) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *AuthBody) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *AuthBody) GetExt() []byte {
	if x != nil {
		return x.Ext
	}
	return nil
}

type SendMsgReplyBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type SendMsgReplyBody_ErrorType `protobuf:"varint,1,opt,name=type,proto3,enum=im.protocol.SendMsgReplyBody_ErrorType" json:"type,omitempty"`
	Msg  string                     `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *SendMsgReplyBody) Reset() {
	*x = SendMsgReplyBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_protocol_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMsgReplyBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMsgReplyBody) ProtoMessage() {}

func (x *SendMsgReplyBody) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_protocol_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMsgReplyBody.ProtoReflect.Descriptor instead.
func (*SendMsgReplyBody) Descriptor() ([]byte, []int) {
	return file_protocol_protocol_proto_rawDescGZIP(), []int{1}
}

func (x *SendMsgReplyBody) GetType() SendMsgReplyBody_ErrorType {
	if x != nil {
		return x.Type
	}
	return SendMsgReplyBody_IsOk
}

func (x *SendMsgReplyBody) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

// Proto 中 Op 为 SyncMsgReply 时， body 必须可以反序列化为 SyncMsg
type SyncMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogId int64 `protobuf:"varint,1,opt,name=logId,proto3" json:"logId,omitempty"`
}

func (x *SyncMsg) Reset() {
	*x = SyncMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_protocol_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncMsg) ProtoMessage() {}

func (x *SyncMsg) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_protocol_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncMsg.ProtoReflect.Descriptor instead.
func (*SyncMsg) Descriptor() ([]byte, []int) {
	return file_protocol_protocol_proto_rawDescGZIP(), []int{2}
}

func (x *SyncMsg) GetLogId() int64 {
	if x != nil {
		return x.LogId
	}
	return 0
}

type Proto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver     int32   `protobuf:"varint,1,opt,name=ver,proto3" json:"ver,omitempty"`
	Op      int32   `protobuf:"varint,2,opt,name=op,proto3" json:"op,omitempty"`
	Seq     int32   `protobuf:"varint,3,opt,name=seq,proto3" json:"seq,omitempty"`
	Ack     int32   `protobuf:"varint,4,opt,name=ack,proto3" json:"ack,omitempty"`
	Mid     int64   `protobuf:"varint,5,opt,name=mid,proto3" json:"mid,omitempty"`
	Channel Channel `protobuf:"varint,6,opt,name=channel,proto3,enum=im.protocol.Channel" json:"channel,omitempty"`
	Target  string  `protobuf:"bytes,7,opt,name=target,proto3" json:"target,omitempty"`
	Time    int64   `protobuf:"varint,8,opt,name=time,proto3" json:"time,omitempty"`
	Body    []byte  `protobuf:"bytes,9,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Proto) Reset() {
	*x = Proto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_protocol_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Proto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Proto) ProtoMessage() {}

func (x *Proto) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_protocol_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Proto.ProtoReflect.Descriptor instead.
func (*Proto) Descriptor() ([]byte, []int) {
	return file_protocol_protocol_proto_rawDescGZIP(), []int{3}
}

func (x *Proto) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *Proto) GetOp() int32 {
	if x != nil {
		return x.Op
	}
	return 0
}

func (x *Proto) GetSeq() int32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *Proto) GetAck() int32 {
	if x != nil {
		return x.Ack
	}
	return 0
}

func (x *Proto) GetMid() int64 {
	if x != nil {
		return x.Mid
	}
	return 0
}

func (x *Proto) GetChannel() Channel {
	if x != nil {
		return x.Channel
	}
	return Channel_Private
}

func (x *Proto) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *Proto) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Proto) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_protocol_protocol_proto protoreflect.FileDescriptor

var file_protocol_protocol_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x48, 0x0a, 0x08, 0x41, 0x75, 0x74, 0x68, 0x42, 0x6f,
	0x64, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x65, 0x78, 0x74,
	0x22, 0xb3, 0x01, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x3b, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x69, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x6f,
	0x64, 0x79, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x22, 0x50, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x73, 0x4f, 0x6b, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x4e,
	0x6f, 0x74, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x4e, 0x6f,
	0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x10, 0x02, 0x12, 0x16,
	0x0a, 0x12, 0x55, 0x6e, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x10, 0x03, 0x22, 0x1f, 0x0a, 0x07, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73,
	0x67, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x6c, 0x6f, 0x67, 0x49, 0x64, 0x22, 0xcf, 0x01, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x76, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x6f, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x73, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x61, 0x63, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6d, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x07, 0x63, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x69, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x2a, 0x95, 0x02, 0x0a, 0x02, 0x4f, 0x70,
	0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x10, 0x00, 0x12,
	0x08, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x75, 0x74,
	0x68, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x72, 0x74,
	0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x10, 0x05, 0x12, 0x13, 0x0a, 0x0f, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x06,
	0x12, 0x0b, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x10, 0x07, 0x12, 0x10, 0x0a,
	0x0c, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x08, 0x12,
	0x0e, 0x0a, 0x0a, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x73, 0x67, 0x10, 0x09, 0x12,
	0x13, 0x0a, 0x0f, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x10, 0x0a, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x52,
	0x65, 0x71, 0x10, 0x0b, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x10, 0x0c, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52,
	0x65, 0x61, 0x64, 0x79, 0x10, 0x14, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x46,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x10, 0x15, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x52, 0x65, 0x73, 0x65, 0x6e, 0x64, 0x10, 0x16, 0x12, 0x07, 0x0a, 0x03, 0x52, 0x61, 0x77, 0x10,
	0x17, 0x2a, 0x21, 0x0a, 0x07, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x0b, 0x0a, 0x07,
	0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x10, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x74, 0x78, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x69, 0x6d, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protocol_protocol_proto_rawDescOnce sync.Once
	file_protocol_protocol_proto_rawDescData = file_protocol_protocol_proto_rawDesc
)

func file_protocol_protocol_proto_rawDescGZIP() []byte {
	file_protocol_protocol_proto_rawDescOnce.Do(func() {
		file_protocol_protocol_proto_rawDescData = protoimpl.X.CompressGZIP(file_protocol_protocol_proto_rawDescData)
	})
	return file_protocol_protocol_proto_rawDescData
}

var file_protocol_protocol_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_protocol_protocol_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protocol_protocol_proto_goTypes = []interface{}{
	(Op)(0),                         // 0: im.protocol.Op
	(Channel)(0),                    // 1: im.protocol.Channel
	(SendMsgReplyBody_ErrorType)(0), // 2: im.protocol.SendMsgReplyBody.ErrorType
	(*AuthBody)(nil),                // 3: im.protocol.AuthBody
	(*SendMsgReplyBody)(nil),        // 4: im.protocol.SendMsgReplyBody
	(*SyncMsg)(nil),                 // 5: im.protocol.SyncMsg
	(*Proto)(nil),                   // 6: im.protocol.Proto
}
var file_protocol_protocol_proto_depIdxs = []int32{
	2, // 0: im.protocol.SendMsgReplyBody.type:type_name -> im.protocol.SendMsgReplyBody.ErrorType
	1, // 1: im.protocol.Proto.channel:type_name -> im.protocol.Channel
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_protocol_protocol_proto_init() }
func file_protocol_protocol_proto_init() {
	if File_protocol_protocol_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protocol_protocol_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthBody); i {
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
		file_protocol_protocol_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMsgReplyBody); i {
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
		file_protocol_protocol_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncMsg); i {
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
		file_protocol_protocol_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Proto); i {
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
			RawDescriptor: file_protocol_protocol_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protocol_protocol_proto_goTypes,
		DependencyIndexes: file_protocol_protocol_proto_depIdxs,
		EnumInfos:         file_protocol_protocol_proto_enumTypes,
		MessageInfos:      file_protocol_protocol_proto_msgTypes,
	}.Build()
	File_protocol_protocol_proto = out.File
	file_protocol_protocol_proto_rawDesc = nil
	file_protocol_protocol_proto_goTypes = nil
	file_protocol_protocol_proto_depIdxs = nil
}
