// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: many.proto

package server

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

type GetManyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetManyRequest) Reset() {
	*x = GetManyRequest{}
	mi := &file_many_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetManyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetManyRequest) ProtoMessage() {}

func (x *GetManyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_many_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetManyRequest.ProtoReflect.Descriptor instead.
func (*GetManyRequest) Descriptor() ([]byte, []int) {
	return file_many_proto_rawDescGZIP(), []int{0}
}

type GetManyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserData []*GetResponse `protobuf:"bytes,1,rep,name=UserData,proto3" json:"UserData,omitempty"`
}

func (x *GetManyResponse) Reset() {
	*x = GetManyResponse{}
	mi := &file_many_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetManyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetManyResponse) ProtoMessage() {}

func (x *GetManyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_many_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetManyResponse.ProtoReflect.Descriptor instead.
func (*GetManyResponse) Descriptor() ([]byte, []int) {
	return file_many_proto_rawDescGZIP(), []int{1}
}

func (x *GetManyResponse) GetUserData() []*GetResponse {
	if x != nil {
		return x.UserData
	}
	return nil
}

var File_many_proto protoreflect.FileDescriptor

var file_many_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x09, 0x67, 0x65,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x10, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4d, 0x61,
	0x6e, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3b, 0x0a, 0x0f, 0x47, 0x65, 0x74,
	0x4d, 0x61, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x55, 0x73,
	0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x42, 0x17, 0x5a, 0x15, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_many_proto_rawDescOnce sync.Once
	file_many_proto_rawDescData = file_many_proto_rawDesc
)

func file_many_proto_rawDescGZIP() []byte {
	file_many_proto_rawDescOnce.Do(func() {
		file_many_proto_rawDescData = protoimpl.X.CompressGZIP(file_many_proto_rawDescData)
	})
	return file_many_proto_rawDescData
}

var file_many_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_many_proto_goTypes = []any{
	(*GetManyRequest)(nil),  // 0: GetManyRequest
	(*GetManyResponse)(nil), // 1: GetManyResponse
	(*GetResponse)(nil),     // 2: GetResponse
}
var file_many_proto_depIdxs = []int32{
	2, // 0: GetManyResponse.UserData:type_name -> GetResponse
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_many_proto_init() }
func file_many_proto_init() {
	if File_many_proto != nil {
		return
	}
	file_get_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_many_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_many_proto_goTypes,
		DependencyIndexes: file_many_proto_depIdxs,
		MessageInfos:      file_many_proto_msgTypes,
	}.Build()
	File_many_proto = out.File
	file_many_proto_rawDesc = nil
	file_many_proto_goTypes = nil
	file_many_proto_depIdxs = nil
}
