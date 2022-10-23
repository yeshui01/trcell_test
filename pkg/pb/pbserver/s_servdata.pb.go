// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: s_servdata.proto

package pbserver

import (
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

type ServDataTable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServTables []*DbTableData `protobuf:"bytes,1,rep,name=ServTables,proto3" json:"ServTables,omitempty"` // 服务器相关的数据表
}

func (x *ServDataTable) Reset() {
	*x = ServDataTable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServDataTable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServDataTable) ProtoMessage() {}

func (x *ServDataTable) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServDataTable.ProtoReflect.Descriptor instead.
func (*ServDataTable) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{0}
}

func (x *ServDataTable) GetServTables() []*DbTableData {
	if x != nil {
		return x.ServTables
	}
	return nil
}

type ServDataTableList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableID  int32              `protobuf:"varint,1,opt,name=TableID,proto3" json:"TableID,omitempty"` // 表id
	DataList []*DbTableListItem `protobuf:"bytes,2,rep,name=DataList,proto3" json:"DataList,omitempty"`
}

func (x *ServDataTableList) Reset() {
	*x = ServDataTableList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServDataTableList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServDataTableList) ProtoMessage() {}

func (x *ServDataTableList) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServDataTableList.ProtoReflect.Descriptor instead.
func (*ServDataTableList) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{1}
}

func (x *ServDataTableList) GetTableID() int32 {
	if x != nil {
		return x.TableID
	}
	return 0
}

func (x *ServDataTableList) GetDataList() []*DbTableListItem {
	if x != nil {
		return x.DataList
	}
	return nil
}

// ESMsgServDataLoadTables              = 1 // 加载数据表
type ESMsgServDataLoadTablesReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LoadOwner string `protobuf:"bytes,1,opt,name=LoadOwner,proto3" json:"LoadOwner,omitempty"` // 加载所属对象
}

func (x *ESMsgServDataLoadTablesReq) Reset() {
	*x = ESMsgServDataLoadTablesReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataLoadTablesReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataLoadTablesReq) ProtoMessage() {}

func (x *ESMsgServDataLoadTablesReq) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataLoadTablesReq.ProtoReflect.Descriptor instead.
func (*ESMsgServDataLoadTablesReq) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{2}
}

func (x *ESMsgServDataLoadTablesReq) GetLoadOwner() string {
	if x != nil {
		return x.LoadOwner
	}
	return ""
}

type ESMsgServDataLoadTablesRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataTable *ServDataTable `protobuf:"bytes,1,opt,name=DataTable,proto3" json:"DataTable,omitempty"`
}

func (x *ESMsgServDataLoadTablesRep) Reset() {
	*x = ESMsgServDataLoadTablesRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataLoadTablesRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataLoadTablesRep) ProtoMessage() {}

func (x *ESMsgServDataLoadTablesRep) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataLoadTablesRep.ProtoReflect.Descriptor instead.
func (*ESMsgServDataLoadTablesRep) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{3}
}

func (x *ESMsgServDataLoadTablesRep) GetDataTable() *ServDataTable {
	if x != nil {
		return x.DataTable
	}
	return nil
}

// ESMsgServDataPushTablesPartial = 2 // 推送数据表分片数据
type ESMsgServDataPushTablesPartialNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataTable *ServDataTable `protobuf:"bytes,1,opt,name=DataTable,proto3" json:"DataTable,omitempty"`
}

func (x *ESMsgServDataPushTablesPartialNotify) Reset() {
	*x = ESMsgServDataPushTablesPartialNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataPushTablesPartialNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataPushTablesPartialNotify) ProtoMessage() {}

func (x *ESMsgServDataPushTablesPartialNotify) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataPushTablesPartialNotify.ProtoReflect.Descriptor instead.
func (*ESMsgServDataPushTablesPartialNotify) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{4}
}

func (x *ESMsgServDataPushTablesPartialNotify) GetDataTable() *ServDataTable {
	if x != nil {
		return x.DataTable
	}
	return nil
}

// ESMsgServDataSaveTables        = 3 // 保存数据表
type ESMsgServDataSaveTablesReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataTable *ServDataTable `protobuf:"bytes,1,opt,name=DataTable,proto3" json:"DataTable,omitempty"`
}

func (x *ESMsgServDataSaveTablesReq) Reset() {
	*x = ESMsgServDataSaveTablesReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataSaveTablesReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataSaveTablesReq) ProtoMessage() {}

func (x *ESMsgServDataSaveTablesReq) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataSaveTablesReq.ProtoReflect.Descriptor instead.
func (*ESMsgServDataSaveTablesReq) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{5}
}

func (x *ESMsgServDataSaveTablesReq) GetDataTable() *ServDataTable {
	if x != nil {
		return x.DataTable
	}
	return nil
}

type ESMsgServDataSaveTablesRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ESMsgServDataSaveTablesRep) Reset() {
	*x = ESMsgServDataSaveTablesRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataSaveTablesRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataSaveTablesRep) ProtoMessage() {}

func (x *ESMsgServDataSaveTablesRep) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataSaveTablesRep.ProtoReflect.Descriptor instead.
func (*ESMsgServDataSaveTablesRep) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{6}
}

// ESMsgServDataLoadTableList        = 4 // 加载数据表(列表)
type ESMsgServDataLoadTableListReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LoadOwner string `protobuf:"bytes,1,opt,name=LoadOwner,proto3" json:"LoadOwner,omitempty"`
}

func (x *ESMsgServDataLoadTableListReq) Reset() {
	*x = ESMsgServDataLoadTableListReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataLoadTableListReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataLoadTableListReq) ProtoMessage() {}

func (x *ESMsgServDataLoadTableListReq) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataLoadTableListReq.ProtoReflect.Descriptor instead.
func (*ESMsgServDataLoadTableListReq) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{7}
}

func (x *ESMsgServDataLoadTableListReq) GetLoadOwner() string {
	if x != nil {
		return x.LoadOwner
	}
	return ""
}

type ESMsgServDataLoadTableListRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableList []*ServDataTableList `protobuf:"bytes,1,rep,name=TableList,proto3" json:"TableList,omitempty"`
}

func (x *ESMsgServDataLoadTableListRep) Reset() {
	*x = ESMsgServDataLoadTableListRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataLoadTableListRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataLoadTableListRep) ProtoMessage() {}

func (x *ESMsgServDataLoadTableListRep) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataLoadTableListRep.ProtoReflect.Descriptor instead.
func (*ESMsgServDataLoadTableListRep) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{8}
}

func (x *ESMsgServDataLoadTableListRep) GetTableList() []*ServDataTableList {
	if x != nil {
		return x.TableList
	}
	return nil
}

// ESMsgServDataSaveTableList        = 5 // 保存数据表(列表)
type ESMsgServDataSaveTableListReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableList []*ServDataTableList `protobuf:"bytes,1,rep,name=TableList,proto3" json:"TableList,omitempty"`
}

func (x *ESMsgServDataSaveTableListReq) Reset() {
	*x = ESMsgServDataSaveTableListReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataSaveTableListReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataSaveTableListReq) ProtoMessage() {}

func (x *ESMsgServDataSaveTableListReq) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataSaveTableListReq.ProtoReflect.Descriptor instead.
func (*ESMsgServDataSaveTableListReq) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{9}
}

func (x *ESMsgServDataSaveTableListReq) GetTableList() []*ServDataTableList {
	if x != nil {
		return x.TableList
	}
	return nil
}

type ESMsgServDataSaveTableListRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ESMsgServDataSaveTableListRep) Reset() {
	*x = ESMsgServDataSaveTableListRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataSaveTableListRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataSaveTableListRep) ProtoMessage() {}

func (x *ESMsgServDataSaveTableListRep) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataSaveTableListRep.ProtoReflect.Descriptor instead.
func (*ESMsgServDataSaveTableListRep) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{10}
}

// ESMsgServDataPushTableListPartial = 6 // 推送数据列表 分片数据
type ESMsgServDataPushTableListPartialNotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableList []*ServDataTableList `protobuf:"bytes,1,rep,name=TableList,proto3" json:"TableList,omitempty"`
}

func (x *ESMsgServDataPushTableListPartialNotify) Reset() {
	*x = ESMsgServDataPushTableListPartialNotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_s_servdata_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ESMsgServDataPushTableListPartialNotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ESMsgServDataPushTableListPartialNotify) ProtoMessage() {}

func (x *ESMsgServDataPushTableListPartialNotify) ProtoReflect() protoreflect.Message {
	mi := &file_s_servdata_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ESMsgServDataPushTableListPartialNotify.ProtoReflect.Descriptor instead.
func (*ESMsgServDataPushTableListPartialNotify) Descriptor() ([]byte, []int) {
	return file_s_servdata_proto_rawDescGZIP(), []int{11}
}

func (x *ESMsgServDataPushTableListPartialNotify) GetTableList() []*ServDataTableList {
	if x != nil {
		return x.TableList
	}
	return nil
}

var File_s_servdata_proto protoreflect.FileDescriptor

var file_s_servdata_proto_rawDesc = []byte{
	0x0a, 0x10, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x70, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x1a, 0x0a, 0x73, 0x5f,
	0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76,
	0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x35, 0x0a, 0x0a, 0x53, 0x65, 0x72,
	0x76, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x70, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x44, 0x62, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x0a, 0x53, 0x65, 0x72, 0x76, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73,
	0x22, 0x64, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x49, 0x44, 0x12,
	0x35, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x70, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x44, 0x62, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x08, 0x44, 0x61,
	0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x3a, 0x0a, 0x1a, 0x45, 0x53, 0x4d, 0x73, 0x67, 0x53,
	0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x6f, 0x61, 0x64, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x4c, 0x6f, 0x61, 0x64, 0x4f, 0x77, 0x6e, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4c, 0x6f, 0x61, 0x64, 0x4f, 0x77, 0x6e,
	0x65, 0x72, 0x22, 0x53, 0x0a, 0x1a, 0x45, 0x53, 0x4d, 0x73, 0x67, 0x53, 0x65, 0x72, 0x76, 0x44,
	0x61, 0x74, 0x61, 0x4c, 0x6f, 0x61, 0x64, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x70,
	0x12, 0x35, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x09, 0x44, 0x61,
	0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x5d, 0x0a, 0x24, 0x45, 0x53, 0x4d, 0x73, 0x67,
	0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x50, 0x75, 0x73, 0x68, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12,
	0x35, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x09, 0x44, 0x61, 0x74,
	0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x53, 0x0a, 0x1a, 0x45, 0x53, 0x4d, 0x73, 0x67, 0x53,
	0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x53, 0x61, 0x76, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x12, 0x35, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x09, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x1c, 0x0a, 0x1a, 0x45,
	0x53, 0x4d, 0x73, 0x67, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x53, 0x61, 0x76, 0x65,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x70, 0x22, 0x3d, 0x0a, 0x1d, 0x45, 0x53, 0x4d,
	0x73, 0x67, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x6f, 0x61, 0x64, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x4c, 0x6f,
	0x61, 0x64, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x4c,
	0x6f, 0x61, 0x64, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x22, 0x5a, 0x0a, 0x1d, 0x45, 0x53, 0x4d, 0x73,
	0x67, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x6f, 0x61, 0x64, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x70, 0x12, 0x39, 0x0a, 0x09, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x09, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x5a, 0x0a, 0x1d, 0x45, 0x53, 0x4d, 0x73, 0x67, 0x53, 0x65, 0x72,
	0x76, 0x44, 0x61, 0x74, 0x61, 0x53, 0x61, 0x76, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x12, 0x39, 0x0a, 0x09, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x62, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x09, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74,
	0x22, 0x1f, 0x0a, 0x1d, 0x45, 0x53, 0x4d, 0x73, 0x67, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61, 0x74,
	0x61, 0x53, 0x61, 0x76, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x70, 0x22, 0x64, 0x0a, 0x27, 0x45, 0x53, 0x4d, 0x73, 0x67, 0x53, 0x65, 0x72, 0x76, 0x44, 0x61,
	0x74, 0x61, 0x50, 0x75, 0x73, 0x68, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x39, 0x0a, 0x09,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x70, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x44,
	0x61, 0x74, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x09, 0x54, 0x61,
	0x62, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_s_servdata_proto_rawDescOnce sync.Once
	file_s_servdata_proto_rawDescData = file_s_servdata_proto_rawDesc
)

func file_s_servdata_proto_rawDescGZIP() []byte {
	file_s_servdata_proto_rawDescOnce.Do(func() {
		file_s_servdata_proto_rawDescData = protoimpl.X.CompressGZIP(file_s_servdata_proto_rawDescData)
	})
	return file_s_servdata_proto_rawDescData
}

var file_s_servdata_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_s_servdata_proto_goTypes = []interface{}{
	(*ServDataTable)(nil),                           // 0: pbserver.ServDataTable
	(*ServDataTableList)(nil),                       // 1: pbserver.ServDataTableList
	(*ESMsgServDataLoadTablesReq)(nil),              // 2: pbserver.ESMsgServDataLoadTablesReq
	(*ESMsgServDataLoadTablesRep)(nil),              // 3: pbserver.ESMsgServDataLoadTablesRep
	(*ESMsgServDataPushTablesPartialNotify)(nil),    // 4: pbserver.ESMsgServDataPushTablesPartialNotify
	(*ESMsgServDataSaveTablesReq)(nil),              // 5: pbserver.ESMsgServDataSaveTablesReq
	(*ESMsgServDataSaveTablesRep)(nil),              // 6: pbserver.ESMsgServDataSaveTablesRep
	(*ESMsgServDataLoadTableListReq)(nil),           // 7: pbserver.ESMsgServDataLoadTableListReq
	(*ESMsgServDataLoadTableListRep)(nil),           // 8: pbserver.ESMsgServDataLoadTableListRep
	(*ESMsgServDataSaveTableListReq)(nil),           // 9: pbserver.ESMsgServDataSaveTableListReq
	(*ESMsgServDataSaveTableListRep)(nil),           // 10: pbserver.ESMsgServDataSaveTableListRep
	(*ESMsgServDataPushTableListPartialNotify)(nil), // 11: pbserver.ESMsgServDataPushTableListPartialNotify
	(*DbTableData)(nil),                             // 12: pbserver.DbTableData
	(*DbTableListItem)(nil),                         // 13: pbserver.DbTableListItem
}
var file_s_servdata_proto_depIdxs = []int32{
	12, // 0: pbserver.ServDataTable.ServTables:type_name -> pbserver.DbTableData
	13, // 1: pbserver.ServDataTableList.DataList:type_name -> pbserver.DbTableListItem
	0,  // 2: pbserver.ESMsgServDataLoadTablesRep.DataTable:type_name -> pbserver.ServDataTable
	0,  // 3: pbserver.ESMsgServDataPushTablesPartialNotify.DataTable:type_name -> pbserver.ServDataTable
	0,  // 4: pbserver.ESMsgServDataSaveTablesReq.DataTable:type_name -> pbserver.ServDataTable
	1,  // 5: pbserver.ESMsgServDataLoadTableListRep.TableList:type_name -> pbserver.ServDataTableList
	1,  // 6: pbserver.ESMsgServDataSaveTableListReq.TableList:type_name -> pbserver.ServDataTableList
	1,  // 7: pbserver.ESMsgServDataPushTableListPartialNotify.TableList:type_name -> pbserver.ServDataTableList
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_s_servdata_proto_init() }
func file_s_servdata_proto_init() {
	if File_s_servdata_proto != nil {
		return
	}
	file_s_db_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_s_servdata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServDataTable); i {
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
		file_s_servdata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServDataTableList); i {
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
		file_s_servdata_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataLoadTablesReq); i {
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
		file_s_servdata_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataLoadTablesRep); i {
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
		file_s_servdata_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataPushTablesPartialNotify); i {
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
		file_s_servdata_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataSaveTablesReq); i {
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
		file_s_servdata_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataSaveTablesRep); i {
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
		file_s_servdata_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataLoadTableListReq); i {
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
		file_s_servdata_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataLoadTableListRep); i {
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
		file_s_servdata_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataSaveTableListReq); i {
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
		file_s_servdata_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataSaveTableListRep); i {
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
		file_s_servdata_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ESMsgServDataPushTableListPartialNotify); i {
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
			RawDescriptor: file_s_servdata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_s_servdata_proto_goTypes,
		DependencyIndexes: file_s_servdata_proto_depIdxs,
		MessageInfos:      file_s_servdata_proto_msgTypes,
	}.Build()
	File_s_servdata_proto = out.File
	file_s_servdata_proto_rawDesc = nil
	file_s_servdata_proto_goTypes = nil
	file_s_servdata_proto_depIdxs = nil
}