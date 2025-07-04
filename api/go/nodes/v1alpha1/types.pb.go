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
// source: nodes/v1alpha1/types.proto

package v1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
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

type Node struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Node:
	//
	//	*Node_Document
	//	*Node_Section
	Node          isNode_Node `protobuf_oneof:"node"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Node) Reset() {
	*x = Node{}
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_types_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetNode() isNode_Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *Node) GetDocument() *Document {
	if x != nil {
		if x, ok := x.Node.(*Node_Document); ok {
			return x.Document
		}
	}
	return nil
}

func (x *Node) GetSection() *Section {
	if x != nil {
		if x, ok := x.Node.(*Node_Section); ok {
			return x.Section
		}
	}
	return nil
}

type isNode_Node interface {
	isNode_Node()
}

type Node_Document struct {
	Document *Document `protobuf:"bytes,1,opt,name=document,proto3,oneof"`
}

type Node_Section struct {
	Section *Section `protobuf:"bytes,2,opt,name=section,proto3,oneof"`
}

func (*Node_Document) isNode_Node() {}

func (*Node_Section) isNode_Node() {}

type Document struct {
	state         protoimpl.MessageState    `protogen:"open.v1"`
	Id            string                    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Workspace     string                    `protobuf:"bytes,2,opt,name=workspace,proto3" json:"workspace,omitempty"`
	Children      []string                  `protobuf:"bytes,3,rep,name=children,proto3" json:"children,omitempty"`
	Metadata      []*Document_MetadataEntry `protobuf:"bytes,4,rep,name=metadata,proto3" json:"metadata,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Document) Reset() {
	*x = Document{}
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Document) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Document) ProtoMessage() {}

func (x *Document) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Document.ProtoReflect.Descriptor instead.
func (*Document) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_types_proto_rawDescGZIP(), []int{1}
}

func (x *Document) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Document) GetWorkspace() string {
	if x != nil {
		return x.Workspace
	}
	return ""
}

func (x *Document) GetChildren() []string {
	if x != nil {
		return x.Children
	}
	return nil
}

func (x *Document) GetMetadata() []*Document_MetadataEntry {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type Section struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Parent        string                 `protobuf:"bytes,2,opt,name=parent,proto3" json:"parent,omitempty"`
	Title         string                 `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Level         int32                  `protobuf:"varint,4,opt,name=level,proto3" json:"level,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Section) Reset() {
	*x = Section{}
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Section) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Section) ProtoMessage() {}

func (x *Section) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Section.ProtoReflect.Descriptor instead.
func (*Section) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_types_proto_rawDescGZIP(), []int{2}
}

func (x *Section) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Section) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (x *Section) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Section) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

// Use array to maintain order
type Document_MetadataEntry struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         *anypb.Any             `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Document_MetadataEntry) Reset() {
	*x = Document_MetadataEntry{}
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Document_MetadataEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Document_MetadataEntry) ProtoMessage() {}

func (x *Document_MetadataEntry) ProtoReflect() protoreflect.Message {
	mi := &file_nodes_v1alpha1_types_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Document_MetadataEntry.ProtoReflect.Descriptor instead.
func (*Document_MetadataEntry) Descriptor() ([]byte, []int) {
	return file_nodes_v1alpha1_types_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Document_MetadataEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Document_MetadataEntry) GetValue() *anypb.Any {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_nodes_v1alpha1_types_proto protoreflect.FileDescriptor

var file_nodes_v1alpha1_types_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31,
	0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x6e, 0x6f,
	0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x87, 0x01,
	0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64,
	0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e,
	0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x07, 0x73, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x06, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0xed, 0x01, 0x0a, 0x08, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x12, 0x48,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2c, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x4d, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x5d, 0x0a, 0x07, 0x53, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x42, 0xc9, 0x01, 0x0a, 0x18, 0x63, 0x6f, 0x6d, 0x2e, 0x6e,
	0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x6e, 0x64, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x42, 0x0a, 0x54, 0x79, 0x70, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x6f,
	0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x6f, 0x72, 0x67, 0x2f, 0x6e, 0x64, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x3b, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0xa2, 0x02, 0x03, 0x4e, 0x4e, 0x58, 0xaa, 0x02, 0x14, 0x4e, 0x6f, 0x74, 0x65, 0x64,
	0x6f, 0x77, 0x6e, 0x2e, 0x4e, 0x64, 0x2e, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xca,
	0x02, 0x14, 0x4e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77, 0x6e, 0x5c, 0x4e, 0x64, 0x5c, 0x56, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0xe2, 0x02, 0x20, 0x4e, 0x6f, 0x74, 0x65, 0x64, 0x6f, 0x77,
	0x6e, 0x5c, 0x4e, 0x64, 0x5c, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x16, 0x4e, 0x6f, 0x74, 0x65,
	0x64, 0x6f, 0x77, 0x6e, 0x3a, 0x3a, 0x4e, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_nodes_v1alpha1_types_proto_rawDescOnce sync.Once
	file_nodes_v1alpha1_types_proto_rawDescData []byte
)

func file_nodes_v1alpha1_types_proto_rawDescGZIP() []byte {
	file_nodes_v1alpha1_types_proto_rawDescOnce.Do(func() {
		file_nodes_v1alpha1_types_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_nodes_v1alpha1_types_proto_rawDesc), len(file_nodes_v1alpha1_types_proto_rawDesc)))
	})
	return file_nodes_v1alpha1_types_proto_rawDescData
}

var file_nodes_v1alpha1_types_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_nodes_v1alpha1_types_proto_goTypes = []any{
	(*Node)(nil),                   // 0: notedown.nd.v1alpha1.Node
	(*Document)(nil),               // 1: notedown.nd.v1alpha1.Document
	(*Section)(nil),                // 2: notedown.nd.v1alpha1.Section
	(*Document_MetadataEntry)(nil), // 3: notedown.nd.v1alpha1.Document.MetadataEntry
	(*anypb.Any)(nil),              // 4: google.protobuf.Any
}
var file_nodes_v1alpha1_types_proto_depIdxs = []int32{
	1, // 0: notedown.nd.v1alpha1.Node.document:type_name -> notedown.nd.v1alpha1.Document
	2, // 1: notedown.nd.v1alpha1.Node.section:type_name -> notedown.nd.v1alpha1.Section
	3, // 2: notedown.nd.v1alpha1.Document.metadata:type_name -> notedown.nd.v1alpha1.Document.MetadataEntry
	4, // 3: notedown.nd.v1alpha1.Document.MetadataEntry.value:type_name -> google.protobuf.Any
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_nodes_v1alpha1_types_proto_init() }
func file_nodes_v1alpha1_types_proto_init() {
	if File_nodes_v1alpha1_types_proto != nil {
		return
	}
	file_nodes_v1alpha1_types_proto_msgTypes[0].OneofWrappers = []any{
		(*Node_Document)(nil),
		(*Node_Section)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_nodes_v1alpha1_types_proto_rawDesc), len(file_nodes_v1alpha1_types_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_nodes_v1alpha1_types_proto_goTypes,
		DependencyIndexes: file_nodes_v1alpha1_types_proto_depIdxs,
		MessageInfos:      file_nodes_v1alpha1_types_proto_msgTypes,
	}.Build()
	File_nodes_v1alpha1_types_proto = out.File
	file_nodes_v1alpha1_types_proto_goTypes = nil
	file_nodes_v1alpha1_types_proto_depIdxs = nil
}
