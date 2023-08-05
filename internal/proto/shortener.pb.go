// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.20.3
// source: pkg/proto/shortener.proto

package proto

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

type URLData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Id  string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *URLData) Reset() {
	*x = URLData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *URLData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*URLData) ProtoMessage() {}

func (x *URLData) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use URLData.ProtoReflect.Descriptor instead.
func (*URLData) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{0}
}

func (x *URLData) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *URLData) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CreateLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *CreateLinkRequest) Reset() {
	*x = CreateLinkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLinkRequest) ProtoMessage() {}

func (x *CreateLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLinkRequest.ProtoReflect.Descriptor instead.
func (*CreateLinkRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{1}
}

func (x *CreateLinkRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type CreateLinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CreateLinkResponse) Reset() {
	*x = CreateLinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLinkResponse) ProtoMessage() {}

func (x *CreateLinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLinkResponse.ProtoReflect.Descriptor instead.
func (*CreateLinkResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{2}
}

func (x *CreateLinkResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CreateLinkBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []string `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *CreateLinkBatchRequest) Reset() {
	*x = CreateLinkBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLinkBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLinkBatchRequest) ProtoMessage() {}

func (x *CreateLinkBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLinkBatchRequest.ProtoReflect.Descriptor instead.
func (*CreateLinkBatchRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{3}
}

func (x *CreateLinkBatchRequest) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

type CreateLinkBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []*URLData `protobuf:"bytes,2,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *CreateLinkBatchResponse) Reset() {
	*x = CreateLinkBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLinkBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLinkBatchResponse) ProtoMessage() {}

func (x *CreateLinkBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLinkBatchResponse.ProtoReflect.Descriptor instead.
func (*CreateLinkBatchResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{4}
}

func (x *CreateLinkBatchResponse) GetUrls() []*URLData {
	if x != nil {
		return x.Urls
	}
	return nil
}

type GetURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetURLRequest) Reset() {
	*x = GetURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLRequest) ProtoMessage() {}

func (x *GetURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLRequest.ProtoReflect.Descriptor instead.
func (*GetURLRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{5}
}

func (x *GetURLRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *GetURLResponse) Reset() {
	*x = GetURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLResponse) ProtoMessage() {}

func (x *GetURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLResponse.ProtoReflect.Descriptor instead.
func (*GetURLResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{6}
}

func (x *GetURLResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetAllURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAllURLRequest) Reset() {
	*x = GetAllURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllURLRequest) ProtoMessage() {}

func (x *GetAllURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllURLRequest.ProtoReflect.Descriptor instead.
func (*GetAllURLRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{7}
}

type GetAllURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []*URLData `protobuf:"bytes,2,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *GetAllURLResponse) Reset() {
	*x = GetAllURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllURLResponse) ProtoMessage() {}

func (x *GetAllURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllURLResponse.ProtoReflect.Descriptor instead.
func (*GetAllURLResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{8}
}

func (x *GetAllURLResponse) GetUrls() []*URLData {
	if x != nil {
		return x.Urls
	}
	return nil
}

type DeleteURLBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
}

func (x *DeleteURLBatchRequest) Reset() {
	*x = DeleteURLBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteURLBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteURLBatchRequest) ProtoMessage() {}

func (x *DeleteURLBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteURLBatchRequest.ProtoReflect.Descriptor instead.
func (*DeleteURLBatchRequest) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteURLBatchRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

type DeleteURLBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteURLBatchResponse) Reset() {
	*x = DeleteURLBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_shortener_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteURLBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteURLBatchResponse) ProtoMessage() {}

func (x *DeleteURLBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_shortener_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteURLBatchResponse.ProtoReflect.Descriptor instead.
func (*DeleteURLBatchResponse) Descriptor() ([]byte, []int) {
	return file_pkg_proto_shortener_proto_rawDescGZIP(), []int{10}
}

var File_pkg_proto_shortener_proto protoreflect.FileDescriptor

var file_pkg_proto_shortener_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x68, 0x6f,
	0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x22, 0x2b, 0x0a, 0x07, 0x55, 0x52, 0x4c, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x25, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x24, 0x0a, 0x12, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x2c, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72,
	0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x41,
	0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x75, 0x72, 0x6c,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x2e, 0x55, 0x52, 0x4c, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x75, 0x72, 0x6c,
	0x73, 0x22, 0x1f, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x22, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x12, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c,
	0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3b, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x26, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x55, 0x52, 0x4c, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x22, 0x29, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x55, 0x52, 0x4c, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69,
	0x64, 0x73, 0x22, 0x18, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x8e, 0x03, 0x0a,
	0x09, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x12, 0x49, 0x0a, 0x0a, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x1c, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c,
	0x69, 0x6e, 0x6b, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x21, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69,
	0x6e, 0x6b, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3d, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x18, 0x2e, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e,
	0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46,
	0x0a, 0x09, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x52, 0x4c, 0x12, 0x1b, 0x2e, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x52,
	0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x55, 0x52, 0x4c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x55, 0x52, 0x4c, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x20, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x73, 0x68, 0x6f,
	0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x52, 0x4c,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x35, 0x5a,
	0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x76, 0x61, 0x6e,
	0x70, 0x6f, 0x64, 0x67, 0x6f, 0x72, 0x6e, 0x79, 0x2f, 0x75, 0x72, 0x6c, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_shortener_proto_rawDescOnce sync.Once
	file_pkg_proto_shortener_proto_rawDescData = file_pkg_proto_shortener_proto_rawDesc
)

func file_pkg_proto_shortener_proto_rawDescGZIP() []byte {
	file_pkg_proto_shortener_proto_rawDescOnce.Do(func() {
		file_pkg_proto_shortener_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_shortener_proto_rawDescData)
	})
	return file_pkg_proto_shortener_proto_rawDescData
}

var file_pkg_proto_shortener_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_pkg_proto_shortener_proto_goTypes = []interface{}{
	(*URLData)(nil),                 // 0: shortener.URLData
	(*CreateLinkRequest)(nil),       // 1: shortener.CreateLinkRequest
	(*CreateLinkResponse)(nil),      // 2: shortener.CreateLinkResponse
	(*CreateLinkBatchRequest)(nil),  // 3: shortener.CreateLinkBatchRequest
	(*CreateLinkBatchResponse)(nil), // 4: shortener.CreateLinkBatchResponse
	(*GetURLRequest)(nil),           // 5: shortener.GetURLRequest
	(*GetURLResponse)(nil),          // 6: shortener.GetURLResponse
	(*GetAllURLRequest)(nil),        // 7: shortener.GetAllURLRequest
	(*GetAllURLResponse)(nil),       // 8: shortener.GetAllURLResponse
	(*DeleteURLBatchRequest)(nil),   // 9: shortener.DeleteURLBatchRequest
	(*DeleteURLBatchResponse)(nil),  // 10: shortener.DeleteURLBatchResponse
}
var file_pkg_proto_shortener_proto_depIdxs = []int32{
	0,  // 0: shortener.CreateLinkBatchResponse.urls:type_name -> shortener.URLData
	0,  // 1: shortener.GetAllURLResponse.urls:type_name -> shortener.URLData
	1,  // 2: shortener.Shortener.CreateLink:input_type -> shortener.CreateLinkRequest
	3,  // 3: shortener.Shortener.CreateLinkBatch:input_type -> shortener.CreateLinkBatchRequest
	5,  // 4: shortener.Shortener.GetURL:input_type -> shortener.GetURLRequest
	7,  // 5: shortener.Shortener.GetAllURL:input_type -> shortener.GetAllURLRequest
	9,  // 6: shortener.Shortener.DeleteURLBatch:input_type -> shortener.DeleteURLBatchRequest
	2,  // 7: shortener.Shortener.CreateLink:output_type -> shortener.CreateLinkResponse
	4,  // 8: shortener.Shortener.CreateLinkBatch:output_type -> shortener.CreateLinkBatchResponse
	6,  // 9: shortener.Shortener.GetURL:output_type -> shortener.GetURLResponse
	8,  // 10: shortener.Shortener.GetAllURL:output_type -> shortener.GetAllURLResponse
	10, // 11: shortener.Shortener.DeleteURLBatch:output_type -> shortener.DeleteURLBatchResponse
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_proto_shortener_proto_init() }
func file_pkg_proto_shortener_proto_init() {
	if File_pkg_proto_shortener_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_shortener_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*URLData); i {
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
		file_pkg_proto_shortener_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLinkRequest); i {
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
		file_pkg_proto_shortener_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLinkResponse); i {
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
		file_pkg_proto_shortener_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLinkBatchRequest); i {
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
		file_pkg_proto_shortener_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateLinkBatchResponse); i {
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
		file_pkg_proto_shortener_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetURLRequest); i {
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
		file_pkg_proto_shortener_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetURLResponse); i {
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
		file_pkg_proto_shortener_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllURLRequest); i {
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
		file_pkg_proto_shortener_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllURLResponse); i {
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
		file_pkg_proto_shortener_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteURLBatchRequest); i {
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
		file_pkg_proto_shortener_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteURLBatchResponse); i {
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
			RawDescriptor: file_pkg_proto_shortener_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_proto_shortener_proto_goTypes,
		DependencyIndexes: file_pkg_proto_shortener_proto_depIdxs,
		MessageInfos:      file_pkg_proto_shortener_proto_msgTypes,
	}.Build()
	File_pkg_proto_shortener_proto = out.File
	file_pkg_proto_shortener_proto_rawDesc = nil
	file_pkg_proto_shortener_proto_goTypes = nil
	file_pkg_proto_shortener_proto_depIdxs = nil
}