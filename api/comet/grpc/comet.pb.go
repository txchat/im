// protoc -I=. -I=$GOPATH/src --go_out=plugins=grpc:. *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: comet.proto

package grpc

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
	Op_ProtoReady      Op = 11
	Op_ProtoFinish     Op = 12
	Op_Raw             Op = 13
	Op_SyncMsgReq      Op = 14
	Op_SyncMsgReply    Op = 15
	Op_RePush          Op = 16
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
		11: "ProtoReady",
		12: "ProtoFinish",
		13: "Raw",
		14: "SyncMsgReq",
		15: "SyncMsgReply",
		16: "RePush",
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
		"ProtoReady":      11,
		"ProtoFinish":     12,
		"Raw":             13,
		"SyncMsgReq":      14,
		"SyncMsgReply":    15,
		"RePush":          16,
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
	return file_comet_proto_enumTypes[0].Descriptor()
}

func (Op) Type() protoreflect.EnumType {
	return &file_comet_proto_enumTypes[0]
}

func (x Op) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Op.Descriptor instead.
func (Op) EnumDescriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{0}
}

type Proto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ver  int32  `protobuf:"varint,1,opt,name=ver,proto3" json:"ver,omitempty"`
	Op   int32  `protobuf:"varint,2,opt,name=op,proto3" json:"op,omitempty"`
	Seq  int32  `protobuf:"varint,3,opt,name=seq,proto3" json:"seq,omitempty"`
	Ack  int32  `protobuf:"varint,4,opt,name=ack,proto3" json:"ack,omitempty"`
	Body []byte `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Proto) Reset() {
	*x = Proto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Proto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Proto) ProtoMessage() {}

func (x *Proto) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[0]
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
	return file_comet_proto_rawDescGZIP(), []int{0}
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

func (x *Proto) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

// Proto 中 Op 为 OpAuth 时， body 必须可以反序列化为 AuthMsg
type AuthMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId string `protobuf:"bytes,1,opt,name=appId,proto3" json:"appId,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Ext   []byte `protobuf:"bytes,3,opt,name=ext,proto3" json:"ext,omitempty"` // 其它业务方可能需要的信息
}

func (x *AuthMsg) Reset() {
	*x = AuthMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthMsg) ProtoMessage() {}

func (x *AuthMsg) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthMsg.ProtoReflect.Descriptor instead.
func (*AuthMsg) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{1}
}

func (x *AuthMsg) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *AuthMsg) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *AuthMsg) GetExt() []byte {
	if x != nil {
		return x.Ext
	}
	return nil
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
		mi := &file_comet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncMsg) ProtoMessage() {}

func (x *SyncMsg) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[2]
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
	return file_comet_proto_rawDescGZIP(), []int{2}
}

func (x *SyncMsg) GetLogId() int64 {
	if x != nil {
		return x.LogId
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{3}
}

type PushMsgReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys    []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	ProtoOp int32    `protobuf:"varint,2,opt,name=protoOp,proto3" json:"protoOp,omitempty"`
	Proto   *Proto   `protobuf:"bytes,3,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *PushMsgReq) Reset() {
	*x = PushMsgReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReq) ProtoMessage() {}

func (x *PushMsgReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReq.ProtoReflect.Descriptor instead.
func (*PushMsgReq) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{4}
}

func (x *PushMsgReq) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *PushMsgReq) GetProtoOp() int32 {
	if x != nil {
		return x.ProtoOp
	}
	return 0
}

func (x *PushMsgReq) GetProto() *Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type PushMsgReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index map[string]int32 `protobuf:"bytes,1,rep,name=index,proto3" json:"index,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *PushMsgReply) Reset() {
	*x = PushMsgReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReply) ProtoMessage() {}

func (x *PushMsgReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReply.ProtoReflect.Descriptor instead.
func (*PushMsgReply) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{5}
}

func (x *PushMsgReply) GetIndex() map[string]int32 {
	if x != nil {
		return x.Index
	}
	return nil
}

type BroadcastReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProtoOp int32  `protobuf:"varint,1,opt,name=protoOp,proto3" json:"protoOp,omitempty"`
	Proto   *Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *BroadcastReq) Reset() {
	*x = BroadcastReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastReq) ProtoMessage() {}

func (x *BroadcastReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastReq.ProtoReflect.Descriptor instead.
func (*BroadcastReq) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{6}
}

func (x *BroadcastReq) GetProtoOp() int32 {
	if x != nil {
		return x.ProtoOp
	}
	return 0
}

func (x *BroadcastReq) GetProto() *Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type BroadcastReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastReply) Reset() {
	*x = BroadcastReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastReply) ProtoMessage() {}

func (x *BroadcastReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastReply.ProtoReflect.Descriptor instead.
func (*BroadcastReply) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{7}
}

type BroadcastGroupReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupID string `protobuf:"bytes,1,opt,name=groupID,proto3" json:"groupID,omitempty"`
	Proto   *Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *BroadcastGroupReq) Reset() {
	*x = BroadcastGroupReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastGroupReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastGroupReq) ProtoMessage() {}

func (x *BroadcastGroupReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastGroupReq.ProtoReflect.Descriptor instead.
func (*BroadcastGroupReq) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{8}
}

func (x *BroadcastGroupReq) GetGroupID() string {
	if x != nil {
		return x.GroupID
	}
	return ""
}

func (x *BroadcastGroupReq) GetProto() *Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type BroadcastGroupReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastGroupReply) Reset() {
	*x = BroadcastGroupReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastGroupReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastGroupReply) ProtoMessage() {}

func (x *BroadcastGroupReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastGroupReply.ProtoReflect.Descriptor instead.
func (*BroadcastGroupReply) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{9}
}

type JoinGroupsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Gid  []string `protobuf:"bytes,2,rep,name=gid,proto3" json:"gid,omitempty"`
}

func (x *JoinGroupsReq) Reset() {
	*x = JoinGroupsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinGroupsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinGroupsReq) ProtoMessage() {}

func (x *JoinGroupsReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinGroupsReq.ProtoReflect.Descriptor instead.
func (*JoinGroupsReq) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{10}
}

func (x *JoinGroupsReq) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *JoinGroupsReq) GetGid() []string {
	if x != nil {
		return x.Gid
	}
	return nil
}

type JoinGroupsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *JoinGroupsReply) Reset() {
	*x = JoinGroupsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinGroupsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinGroupsReply) ProtoMessage() {}

func (x *JoinGroupsReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinGroupsReply.ProtoReflect.Descriptor instead.
func (*JoinGroupsReply) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{11}
}

type LeaveGroupsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Gid  []string `protobuf:"bytes,2,rep,name=gid,proto3" json:"gid,omitempty"`
}

func (x *LeaveGroupsReq) Reset() {
	*x = LeaveGroupsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveGroupsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveGroupsReq) ProtoMessage() {}

func (x *LeaveGroupsReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaveGroupsReq.ProtoReflect.Descriptor instead.
func (*LeaveGroupsReq) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{12}
}

func (x *LeaveGroupsReq) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *LeaveGroupsReq) GetGid() []string {
	if x != nil {
		return x.Gid
	}
	return nil
}

type LeaveGroupsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LeaveGroupsReply) Reset() {
	*x = LeaveGroupsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveGroupsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveGroupsReply) ProtoMessage() {}

func (x *LeaveGroupsReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaveGroupsReply.ProtoReflect.Descriptor instead.
func (*LeaveGroupsReply) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{13}
}

type DelGroupsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gid []string `protobuf:"bytes,1,rep,name=gid,proto3" json:"gid,omitempty"`
}

func (x *DelGroupsReq) Reset() {
	*x = DelGroupsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelGroupsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelGroupsReq) ProtoMessage() {}

func (x *DelGroupsReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelGroupsReq.ProtoReflect.Descriptor instead.
func (*DelGroupsReq) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{14}
}

func (x *DelGroupsReq) GetGid() []string {
	if x != nil {
		return x.Gid
	}
	return nil
}

type DelGroupsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DelGroupsReply) Reset() {
	*x = DelGroupsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelGroupsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelGroupsReply) ProtoMessage() {}

func (x *DelGroupsReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelGroupsReply.ProtoReflect.Descriptor instead.
func (*DelGroupsReply) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{15}
}

var File_comet_proto protoreflect.FileDescriptor

var file_comet_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x69,
	0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x22, 0x61, 0x0a, 0x05, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x10, 0x0a, 0x03, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x76,
	0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02,
	0x6f, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x73, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x61, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x47, 0x0a, 0x07, 0x41, 0x75,
	0x74, 0x68, 0x4d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x65, 0x78, 0x74, 0x22, 0x1f, 0x0a, 0x07, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x12, 0x14,
	0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c,
	0x6f, 0x67, 0x49, 0x64, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x61, 0x0a,
	0x0a, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x6b,
	0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x12, 0x25, 0x0a, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f,
	0x6d, 0x65, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x81, 0x01, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x37, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x50, 0x75, 0x73, 0x68,
	0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x1a, 0x38, 0x0a, 0x0a, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x4f, 0x0a, 0x0c, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x12, 0x25,
	0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x10, 0x0a, 0x0e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x54, 0x0a, 0x11, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x12, 0x18, 0x0a, 0x07,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x49, 0x44, 0x12, 0x25, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x15, 0x0a,
	0x13, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x35, 0x0a, 0x0d, 0x4a, 0x6f, 0x69, 0x6e, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x67, 0x69, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x4a,
	0x6f, 0x69, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x36,
	0x0a, 0x0e, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71,
	0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x6b, 0x65, 0x79, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x03, 0x67, 0x69, 0x64, 0x22, 0x12, 0x0a, 0x10, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x20, 0x0a, 0x0c, 0x44, 0x65,
	0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x67, 0x69, 0x64, 0x22, 0x10, 0x0a, 0x0e,
	0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2a, 0x90,
	0x02, 0x0a, 0x02, 0x4f, 0x70, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e,
	0x65, 0x64, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x02, 0x12, 0x0d, 0x0a,
	0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x04,
	0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x10, 0x05,
	0x12, 0x13, 0x0a, 0x0f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x10, 0x06, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67,
	0x10, 0x07, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x10, 0x08, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d,
	0x73, 0x67, 0x10, 0x09, 0x12, 0x13, 0x0a, 0x0f, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x4d,
	0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x10, 0x0a, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x52, 0x65, 0x61, 0x64, 0x79, 0x10, 0x0b, 0x12, 0x0f, 0x0a, 0x0b, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x10, 0x0c, 0x12, 0x07, 0x0a, 0x03, 0x52, 0x61,
	0x77, 0x10, 0x0d, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x52, 0x65,
	0x71, 0x10, 0x0e, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x79, 0x6e, 0x63, 0x4d, 0x73, 0x67, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x10, 0x0f, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x65, 0x50, 0x75, 0x73, 0x68, 0x10,
	0x10, 0x32, 0x93, 0x03, 0x0a, 0x05, 0x43, 0x6f, 0x6d, 0x65, 0x74, 0x12, 0x37, 0x0a, 0x07, 0x50,
	0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x12, 0x14, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65,
	0x74, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x69,
	0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x3d, 0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73,
	0x74, 0x12, 0x16, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72, 0x6f,
	0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e, 0x69, 0x6d, 0x2e, 0x63,
	0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x4c, 0x0a, 0x0e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1b, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52,
	0x65, 0x71, 0x1a, 0x1d, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x40, 0x0a, 0x0a, 0x4a, 0x6f, 0x69, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12,
	0x17, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f,
	0x6d, 0x65, 0x74, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x43, 0x0a, 0x0b, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x12, 0x18, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x4c, 0x65,
	0x61, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x1a, 0x2e, 0x69,
	0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x3d, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x16, 0x2e, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x2e, 0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x18, 0x2e,
	0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x6c, 0x61,
	0x62, 0x2e, 0x33, 0x33, 0x2e, 0x63, 0x6e, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x69, 0x6d, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comet_proto_rawDescOnce sync.Once
	file_comet_proto_rawDescData = file_comet_proto_rawDesc
)

func file_comet_proto_rawDescGZIP() []byte {
	file_comet_proto_rawDescOnce.Do(func() {
		file_comet_proto_rawDescData = protoimpl.X.CompressGZIP(file_comet_proto_rawDescData)
	})
	return file_comet_proto_rawDescData
}

var file_comet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_comet_proto_msgTypes = make([]protoimpl.MessageInfo, 17)
var file_comet_proto_goTypes = []interface{}{
	(Op)(0),                     // 0: im.comet.Op
	(*Proto)(nil),               // 1: im.comet.Proto
	(*AuthMsg)(nil),             // 2: im.comet.AuthMsg
	(*SyncMsg)(nil),             // 3: im.comet.SyncMsg
	(*Empty)(nil),               // 4: im.comet.Empty
	(*PushMsgReq)(nil),          // 5: im.comet.PushMsgReq
	(*PushMsgReply)(nil),        // 6: im.comet.PushMsgReply
	(*BroadcastReq)(nil),        // 7: im.comet.BroadcastReq
	(*BroadcastReply)(nil),      // 8: im.comet.BroadcastReply
	(*BroadcastGroupReq)(nil),   // 9: im.comet.BroadcastGroupReq
	(*BroadcastGroupReply)(nil), // 10: im.comet.BroadcastGroupReply
	(*JoinGroupsReq)(nil),       // 11: im.comet.JoinGroupsReq
	(*JoinGroupsReply)(nil),     // 12: im.comet.JoinGroupsReply
	(*LeaveGroupsReq)(nil),      // 13: im.comet.LeaveGroupsReq
	(*LeaveGroupsReply)(nil),    // 14: im.comet.LeaveGroupsReply
	(*DelGroupsReq)(nil),        // 15: im.comet.DelGroupsReq
	(*DelGroupsReply)(nil),      // 16: im.comet.DelGroupsReply
	nil,                         // 17: im.comet.PushMsgReply.IndexEntry
}
var file_comet_proto_depIdxs = []int32{
	1,  // 0: im.comet.PushMsgReq.proto:type_name -> im.comet.Proto
	17, // 1: im.comet.PushMsgReply.index:type_name -> im.comet.PushMsgReply.IndexEntry
	1,  // 2: im.comet.BroadcastReq.proto:type_name -> im.comet.Proto
	1,  // 3: im.comet.BroadcastGroupReq.proto:type_name -> im.comet.Proto
	5,  // 4: im.comet.Comet.PushMsg:input_type -> im.comet.PushMsgReq
	7,  // 5: im.comet.Comet.Broadcast:input_type -> im.comet.BroadcastReq
	9,  // 6: im.comet.Comet.BroadcastGroup:input_type -> im.comet.BroadcastGroupReq
	11, // 7: im.comet.Comet.JoinGroups:input_type -> im.comet.JoinGroupsReq
	13, // 8: im.comet.Comet.LeaveGroups:input_type -> im.comet.LeaveGroupsReq
	15, // 9: im.comet.Comet.DelGroups:input_type -> im.comet.DelGroupsReq
	6,  // 10: im.comet.Comet.PushMsg:output_type -> im.comet.PushMsgReply
	8,  // 11: im.comet.Comet.Broadcast:output_type -> im.comet.BroadcastReply
	10, // 12: im.comet.Comet.BroadcastGroup:output_type -> im.comet.BroadcastGroupReply
	12, // 13: im.comet.Comet.JoinGroups:output_type -> im.comet.JoinGroupsReply
	14, // 14: im.comet.Comet.LeaveGroups:output_type -> im.comet.LeaveGroupsReply
	16, // 15: im.comet.Comet.DelGroups:output_type -> im.comet.DelGroupsReply
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_comet_proto_init() }
func file_comet_proto_init() {
	if File_comet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_comet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_comet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthMsg); i {
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
		file_comet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_comet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_comet_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReq); i {
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
		file_comet_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReply); i {
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
		file_comet_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastReq); i {
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
		file_comet_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastReply); i {
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
		file_comet_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastGroupReq); i {
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
		file_comet_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastGroupReply); i {
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
		file_comet_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinGroupsReq); i {
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
		file_comet_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinGroupsReply); i {
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
		file_comet_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaveGroupsReq); i {
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
		file_comet_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaveGroupsReply); i {
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
		file_comet_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelGroupsReq); i {
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
		file_comet_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelGroupsReply); i {
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
			RawDescriptor: file_comet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   17,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_comet_proto_goTypes,
		DependencyIndexes: file_comet_proto_depIdxs,
		EnumInfos:         file_comet_proto_enumTypes,
		MessageInfos:      file_comet_proto_msgTypes,
	}.Build()
	File_comet_proto = out.File
	file_comet_proto_rawDesc = nil
	file_comet_proto_goTypes = nil
	file_comet_proto_depIdxs = nil
}
