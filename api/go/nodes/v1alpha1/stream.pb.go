// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: nodes/v1alpha1/stream.proto

package v1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StreamRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Request:
	//
	//	*StreamRequest_Watch
	Request       isStreamRequest_Request `protobuf_oneof:"request"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{0}
}

func (x *StreamRequest) GetRequest() isStreamRequest_Request {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *StreamRequest) GetWatch() *Watch {
	if x != nil {
		if x, ok := x.Request.(*StreamRequest_Watch); ok {
			return x.Watch
		}
	}
	return nil
}

type isStreamRequest_Request interface {
	isStreamRequest_Request()
}

type StreamRequest_Watch struct {
	Watch *Watch `protobuf:"bytes,1,opt,name=watch,proto3,oneof"`
}

func (*StreamRequest_Watch) isStreamRequest_Request() {}

type Watch struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DocumentWatch *DocumentWatch         `protobuf:"bytes,1,opt,name=document_watch,json=documentWatch,proto3" json:"document_watch,omitempty"`
	SectionWatch  *SectionWatch          `protobuf:"bytes,2,opt,name=section_watch,json=sectionWatch,proto3" json:"section_watch,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Watch) Reset() {
	*x = Watch{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Watch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Watch) ProtoMessage() {}

func (x *Watch) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Watch.ProtoReflect.Descriptor instead.
func (*Watch) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{1}
}

func (x *Watch) GetDocumentWatch() *DocumentWatch {
	if x != nil {
		return x.DocumentWatch
	}
	return nil
}

func (x *Watch) GetSectionWatch() *SectionWatch {
	if x != nil {
		return x.SectionWatch
	}
	return nil
}

type DocumentWatch struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WatchId       string                 `protobuf:"bytes,1,opt,name=watch_id,json=watchId,proto3" json:"watch_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DocumentWatch) Reset() {
	*x = DocumentWatch{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DocumentWatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentWatch) ProtoMessage() {}

func (x *DocumentWatch) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentWatch.ProtoReflect.Descriptor instead.
func (*DocumentWatch) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{2}
}

func (x *DocumentWatch) GetWatchId() string {
	if x != nil {
		return x.WatchId
	}
	return ""
}

type SectionWatch struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WatchId       string                 `protobuf:"bytes,1,opt,name=watch_id,json=watchId,proto3" json:"watch_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SectionWatch) Reset() {
	*x = SectionWatch{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SectionWatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SectionWatch) ProtoMessage() {}

func (x *SectionWatch) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SectionWatch.ProtoReflect.Descriptor instead.
func (*SectionWatch) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{3}
}

func (x *SectionWatch) GetWatchId() string {
	if x != nil {
		return x.WatchId
	}
	return ""
}

type StreamEvent struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Event:
	//
	//	*StreamEvent_InitializationComplete
	//	*StreamEvent_Load
	//	*StreamEvent_Change
	//	*StreamEvent_Delete
	Event         isStreamEvent_Event `protobuf_oneof:"event"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamEvent) Reset() {
	*x = StreamEvent{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamEvent) ProtoMessage() {}

func (x *StreamEvent) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamEvent.ProtoReflect.Descriptor instead.
func (*StreamEvent) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{4}
}

func (x *StreamEvent) GetEvent() isStreamEvent_Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *StreamEvent) GetInitializationComplete() *InitializationComplete {
	if x != nil {
		if x, ok := x.Event.(*StreamEvent_InitializationComplete); ok {
			return x.InitializationComplete
		}
	}
	return nil
}

func (x *StreamEvent) GetLoad() *Load {
	if x != nil {
		if x, ok := x.Event.(*StreamEvent_Load); ok {
			return x.Load
		}
	}
	return nil
}

func (x *StreamEvent) GetChange() *Change {
	if x != nil {
		if x, ok := x.Event.(*StreamEvent_Change); ok {
			return x.Change
		}
	}
	return nil
}

func (x *StreamEvent) GetDelete() *Delete {
	if x != nil {
		if x, ok := x.Event.(*StreamEvent_Delete); ok {
			return x.Delete
		}
	}
	return nil
}

type isStreamEvent_Event interface {
	isStreamEvent_Event()
}

type StreamEvent_InitializationComplete struct {
	InitializationComplete *InitializationComplete `protobuf:"bytes,1,opt,name=initialization_complete,json=initializationComplete,proto3,oneof"`
}

type StreamEvent_Load struct {
	Load *Load `protobuf:"bytes,2,opt,name=load,proto3,oneof"`
}

type StreamEvent_Change struct {
	Change *Change `protobuf:"bytes,3,opt,name=change,proto3,oneof"`
}

type StreamEvent_Delete struct {
	Delete *Delete `protobuf:"bytes,4,opt,name=delete,proto3,oneof"`
}

func (*StreamEvent_InitializationComplete) isStreamEvent_Event() {}

func (*StreamEvent_Load) isStreamEvent_Event() {}

func (*StreamEvent_Change) isStreamEvent_Event() {}

func (*StreamEvent_Delete) isStreamEvent_Event() {}

// Indicates the completion of the initialization of a watch request
type InitializationComplete struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WatchId       string                 `protobuf:"bytes,1,opt,name=watch_id,json=watchId,proto3" json:"watch_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InitializationComplete) Reset() {
	*x = InitializationComplete{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InitializationComplete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializationComplete) ProtoMessage() {}

func (x *InitializationComplete) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializationComplete.ProtoReflect.Descriptor instead.
func (*InitializationComplete) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{5}
}

func (x *InitializationComplete) GetWatchId() string {
	if x != nil {
		return x.WatchId
	}
	return ""
}

// Indicates that a node has been loaded as part of a watch request
type Load struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	WatchId string                 `protobuf:"bytes,1,opt,name=watch_id,json=watchId,proto3" json:"watch_id,omitempty"`
	// Types that are valid to be assigned to Node:
	//
	//	*Load_Document
	//	*Load_Section
	Node          isLoad_Node `protobuf_oneof:"node"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Load) Reset() {
	*x = Load{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Load) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Load) ProtoMessage() {}

func (x *Load) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Load.ProtoReflect.Descriptor instead.
func (*Load) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{6}
}

func (x *Load) GetWatchId() string {
	if x != nil {
		return x.WatchId
	}
	return ""
}

func (x *Load) GetNode() isLoad_Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *Load) GetDocument() *Document {
	if x != nil {
		if x, ok := x.Node.(*Load_Document); ok {
			return x.Document
		}
	}
	return nil
}

func (x *Load) GetSection() *Section {
	if x != nil {
		if x, ok := x.Node.(*Load_Section); ok {
			return x.Section
		}
	}
	return nil
}

type isLoad_Node interface {
	isLoad_Node()
}

type Load_Document struct {
	Document *Document `protobuf:"bytes,10,opt,name=document,proto3,oneof"`
}

type Load_Section struct {
	Section *Section `protobuf:"bytes,11,opt,name=section,proto3,oneof"`
}

func (*Load_Document) isLoad_Node() {}

func (*Load_Section) isLoad_Node() {}

// Indicates that a node has been created or updated
type Change struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	WatchId string                 `protobuf:"bytes,1,opt,name=watch_id,json=watchId,proto3" json:"watch_id,omitempty"`
	// Types that are valid to be assigned to Node:
	//
	//	*Change_Document
	//	*Change_Section
	Node          isChange_Node `protobuf_oneof:"node"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Change) Reset() {
	*x = Change{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Change) ProtoMessage() {}

func (x *Change) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Change.ProtoReflect.Descriptor instead.
func (*Change) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{7}
}

func (x *Change) GetWatchId() string {
	if x != nil {
		return x.WatchId
	}
	return ""
}

func (x *Change) GetNode() isChange_Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *Change) GetDocument() *Document {
	if x != nil {
		if x, ok := x.Node.(*Change_Document); ok {
			return x.Document
		}
	}
	return nil
}

func (x *Change) GetSection() *Section {
	if x != nil {
		if x, ok := x.Node.(*Change_Section); ok {
			return x.Section
		}
	}
	return nil
}

type isChange_Node interface {
	isChange_Node()
}

type Change_Document struct {
	Document *Document `protobuf:"bytes,10,opt,name=document,proto3,oneof"`
}

type Change_Section struct {
	Section *Section `protobuf:"bytes,11,opt,name=section,proto3,oneof"`
}

func (*Change_Document) isChange_Node() {}

func (*Change_Section) isChange_Node() {}

// Indicates that a node has been deleted
type Delete struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	WatchId string                 `protobuf:"bytes,1,opt,name=watch_id,json=watchId,proto3" json:"watch_id,omitempty"`
	// Types that are valid to be assigned to Node:
	//
	//	*Delete_DocumentId
	//	*Delete_SectionId
	Node          isDelete_Node `protobuf_oneof:"node"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Delete) Reset() {
	*x = Delete{}
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Delete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Delete) ProtoMessage() {}

func (x *Delete) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_stream_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Delete.ProtoReflect.Descriptor instead.
func (*Delete) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_stream_proto_rawDescGZIP(), []int{8}
}

func (x *Delete) GetWatchId() string {
	if x != nil {
		return x.WatchId
	}
	return ""
}

func (x *Delete) GetNode() isDelete_Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *Delete) GetDocumentId() string {
	if x != nil {
		if x, ok := x.Node.(*Delete_DocumentId); ok {
			return x.DocumentId
		}
	}
	return ""
}

func (x *Delete) GetSectionId() string {
	if x != nil {
		if x, ok := x.Node.(*Delete_SectionId); ok {
			return x.SectionId
		}
	}
	return ""
}

type isDelete_Node interface {
	isDelete_Node()
}

type Delete_DocumentId struct {
	DocumentId string `protobuf:"bytes,10,opt,name=document_id,json=documentId,proto3,oneof"`
}

type Delete_SectionId struct {
	SectionId string `protobuf:"bytes,11,opt,name=section_id,json=sectionId,proto3,oneof"`
}

func (*Delete_DocumentId) isDelete_Node() {}

func (*Delete_SectionId) isDelete_Node() {}

var File_nodes_v1alpha1_stream_proto protoreflect.FileDescriptor

var file_nodes_v1alpha1_stream_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x6e,
	0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x1a, 0x1a, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x4f, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x33, 0x0a, 0x05, 0x77, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1b, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x48, 0x00, 0x52, 0x05,
	0x77, 0x61, 0x74, 0x63, 0x68, 0x42, 0x09, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x9c, 0x01, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x4a, 0x0a, 0x0e, 0x64, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x77, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x0d, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x47, 0x0a, 0x0d, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x77, 0x61, 0x74, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e,
	0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x0c, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x22,
	0x2a, 0x0a, 0x0d, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x57, 0x61, 0x74, 0x63, 0x68,
	0x12, 0x19, 0x0a, 0x08, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x77, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x22, 0x29, 0x0a, 0x0c, 0x53,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x19, 0x0a, 0x08, 0x77,
	0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x77,
	0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x22, 0xa1, 0x02, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x67, 0x0a, 0x17, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f,
	0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x49,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x00, 0x52, 0x16, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x12,
	0x30, 0x0a, 0x04, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x36, 0x0a, 0x06, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x48,
	0x00, 0x52, 0x06, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x6f, 0x74, 0x65,
	0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x00, 0x52, 0x06, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x33, 0x0a, 0x16, 0x49, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x77, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x22,
	0xa2, 0x01, 0x0a, 0x04, 0x4c, 0x6f, 0x61, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x77, 0x61, 0x74, 0x63,
	0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x77, 0x61, 0x74, 0x63,
	0x68, 0x49, 0x64, 0x12, 0x3c, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e,
	0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63,
	0x75, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x39, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x48, 0x00, 0x52, 0x07, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x06, 0x0a, 0x04,
	0x6e, 0x6f, 0x64, 0x65, 0x22, 0xa4, 0x01, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x77, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64, 0x12, 0x3c, 0x0a, 0x08, 0x64, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6e,
	0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x08,
	0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6e, 0x6f, 0x74, 0x65,
	0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2e, 0x53, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x07, 0x73, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x6f, 0x0a, 0x06, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x77, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x77, 0x61, 0x74, 0x63, 0x68, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0b, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0a, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x09, 0x73, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x32, 0x65, 0x0a, 0x0b,
	0x4e, 0x6f, 0x64, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x06, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x23, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e,
	0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6e, 0x6f, 0x74,
	0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x28,
	0x01, 0x30, 0x01, 0x42, 0xca, 0x01, 0x0a, 0x18, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x6f, 0x74, 0x65,
	0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x42, 0x0b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x6f, 0x74, 0x65,
	0x64, 0x6f, 0x77, 0x6e, 0x6f, 0x72, 0x67, 0x2f, 0x6e, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0xa2, 0x02, 0x03, 0x4e, 0x4e, 0x58, 0xaa, 0x02, 0x14, 0x4e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77,
	0x6e, 0x2e, 0x4e, 0x64, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xca, 0x02, 0x14,
	0x4e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x5c, 0x4e, 0x64, 0x5c, 0x56, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x20, 0x4e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x5c,
	0x4e, 0x64, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x16, 0x4e, 0x6f, 0x74, 0x65, 0x64, 0x6f,
	0x77, 0x6e, 0x3a, 0x3a, 0x4e, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_nodes_v1alpha1_stream_proto_rawDescOnce sync.Once
	file_nodes_v1alpha1_stream_proto_rawDescData []byte
)

func file_nodes_v1alpha1_stream_proto_rawDescGZIP() []byte {
	file_nodes_v1alpha1_stream_proto_rawDescOnce.Do(func() {
		file_nodes_v1alpha1_stream_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_nodes_v1alpha1_stream_proto_rawDesc), len(file_nodes_v1alpha1_stream_proto_rawDesc)))
	})
	return file_nodes_v1alpha1_stream_proto_rawDescData
}

var file_nodes_v1alpha1_stream_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_nodes_v1alpha1_stream_proto_goTypes = []any{
	(*StreamRequest)(nil),          // 0: notedown.nd.v1alpha1.StreamRequest
	(*Watch)(nil),                  // 1: notedown.nd.v1alpha1.Watch
	(*DocumentWatch)(nil),          // 2: notedown.nd.v1alpha1.DocumentWatch
	(*SectionWatch)(nil),           // 3: notedown.nd.v1alpha1.SectionWatch
	(*StreamEvent)(nil),            // 4: notedown.nd.v1alpha1.StreamEvent
	(*InitializationComplete)(nil), // 5: notedown.nd.v1alpha1.InitializationComplete
	(*Load)(nil),                   // 6: notedown.nd.v1alpha1.Load
	(*Change)(nil),                 // 7: notedown.nd.v1alpha1.Change
	(*Delete)(nil),                 // 8: notedown.nd.v1alpha1.Delete
	(*Document)(nil),               // 9: notedown.nd.v1alpha1.Document
	(*Section)(nil),                // 10: notedown.nd.v1alpha1.Section
}
var file_nodes_v1alpha1_stream_proto_depIdxs = []int32{
	1,  // 0: notedown.nd.v1alpha1.StreamRequest.watch:type_name -> notedown.nd.v1alpha1.Watch
	2,  // 1: notedown.nd.v1alpha1.Watch.document_watch:type_name -> notedown.nd.v1alpha1.DocumentWatch
	3,  // 2: notedown.nd.v1alpha1.Watch.section_watch:type_name -> notedown.nd.v1alpha1.SectionWatch
	5,  // 3: notedown.nd.v1alpha1.StreamEvent.initialization_complete:type_name -> notedown.nd.v1alpha1.InitializationComplete
	6,  // 4: notedown.nd.v1alpha1.StreamEvent.load:type_name -> notedown.nd.v1alpha1.Load
	7,  // 5: notedown.nd.v1alpha1.StreamEvent.change:type_name -> notedown.nd.v1alpha1.Change
	8,  // 6: notedown.nd.v1alpha1.StreamEvent.delete:type_name -> notedown.nd.v1alpha1.Delete
	9,  // 7: notedown.nd.v1alpha1.Load.document:type_name -> notedown.nd.v1alpha1.Document
	10, // 8: notedown.nd.v1alpha1.Load.section:type_name -> notedown.nd.v1alpha1.Section
	9,  // 9: notedown.nd.v1alpha1.Change.document:type_name -> notedown.nd.v1alpha1.Document
	10, // 10: notedown.nd.v1alpha1.Change.section:type_name -> notedown.nd.v1alpha1.Section
	0,  // 11: notedown.nd.v1alpha1.NodeService.Stream:input_type -> notedown.nd.v1alpha1.StreamRequest
	4,  // 12: notedown.nd.v1alpha1.NodeService.Stream:output_type -> notedown.nd.v1alpha1.StreamEvent
	12, // [12:13] is the sub-list for method output_type
	11, // [11:12] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_nodes_v1alpha1_stream_proto_init() }
func file_nodes_v1alpha1_stream_proto_init() {
	if File_nodes_v1alpha1_stream_proto != nil {
		return
	}
	file_nodes_v1alpha1_types_proto_init()
	file_nodes_v1alpha1_stream_proto_msgTypes[0].OneofWrappers = []any{
		(*StreamRequest_Watch)(nil),
	}
	file_nodes_v1alpha1_stream_proto_msgTypes[4].OneofWrappers = []any{
		(*StreamEvent_InitializationComplete)(nil),
		(*StreamEvent_Load)(nil),
		(*StreamEvent_Change)(nil),
		(*StreamEvent_Delete)(nil),
	}
	file_nodes_v1alpha1_stream_proto_msgTypes[6].OneofWrappers = []any{
		(*Load_Document)(nil),
		(*Load_Section)(nil),
	}
	file_nodes_v1alpha1_stream_proto_msgTypes[7].OneofWrappers = []any{
		(*Change_Document)(nil),
		(*Change_Section)(nil),
	}
	file_nodes_v1alpha1_stream_proto_msgTypes[8].OneofWrappers = []any{
		(*Delete_DocumentId)(nil),
		(*Delete_SectionId)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_nodes_v1alpha1_stream_proto_rawDesc), len(file_nodes_v1alpha1_stream_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nodes_v1alpha1_stream_proto_goTypes,
		DependencyIndexes: file_nodes_v1alpha1_stream_proto_depIdxs,
		MessageInfos:      file_nodes_v1alpha1_stream_proto_msgTypes,
	}.Build()
	File_nodes_v1alpha1_stream_proto = out.File
	file_nodes_v1alpha1_stream_proto_goTypes = nil
	file_nodes_v1alpha1_stream_proto_depIdxs = nil
}
