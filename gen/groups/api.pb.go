// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: groups/api.proto

package vk_messages

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetFavoritesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Count int32 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *GetFavoritesRequest) Reset() {
	*x = GetFavoritesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_groups_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavoritesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavoritesRequest) ProtoMessage() {}

func (x *GetFavoritesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_groups_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavoritesRequest.ProtoReflect.Descriptor instead.
func (*GetFavoritesRequest) Descriptor() ([]byte, []int) {
	return file_groups_api_proto_rawDescGZIP(), []int{0}
}

func (x *GetFavoritesRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetFavoritesRequest) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type GetFavoriteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Groups []*FavoriteGroup `protobuf:"bytes,1,rep,name=groups,proto3" json:"groups,omitempty"`
}

func (x *GetFavoriteResponse) Reset() {
	*x = GetFavoriteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_groups_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavoriteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavoriteResponse) ProtoMessage() {}

func (x *GetFavoriteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_groups_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavoriteResponse.ProtoReflect.Descriptor instead.
func (*GetFavoriteResponse) Descriptor() ([]byte, []int) {
	return file_groups_api_proto_rawDescGZIP(), []int{1}
}

func (x *GetFavoriteResponse) GetGroups() []*FavoriteGroup {
	if x != nil {
		return x.Groups
	}
	return nil
}

type FavoriteGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar string `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
}

func (x *FavoriteGroup) Reset() {
	*x = FavoriteGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_groups_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavoriteGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteGroup) ProtoMessage() {}

func (x *FavoriteGroup) ProtoReflect() protoreflect.Message {
	mi := &file_groups_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavoriteGroup.ProtoReflect.Descriptor instead.
func (*FavoriteGroup) Descriptor() ([]byte, []int) {
	return file_groups_api_proto_rawDescGZIP(), []int{2}
}

func (x *FavoriteGroup) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FavoriteGroup) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FavoriteGroup) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

type AddFavoriteGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []int32 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (x *AddFavoriteGroupRequest) Reset() {
	*x = AddFavoriteGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_groups_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddFavoriteGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddFavoriteGroupRequest) ProtoMessage() {}

func (x *AddFavoriteGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_groups_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddFavoriteGroupRequest.ProtoReflect.Descriptor instead.
func (*AddFavoriteGroupRequest) Descriptor() ([]byte, []int) {
	return file_groups_api_proto_rawDescGZIP(), []int{3}
}

func (x *AddFavoriteGroupRequest) GetIds() []int32 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type RemoveGroupFromFavoriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []int32 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (x *RemoveGroupFromFavoriteRequest) Reset() {
	*x = RemoveGroupFromFavoriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_groups_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveGroupFromFavoriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveGroupFromFavoriteRequest) ProtoMessage() {}

func (x *RemoveGroupFromFavoriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_groups_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveGroupFromFavoriteRequest.ProtoReflect.Descriptor instead.
func (*RemoveGroupFromFavoriteRequest) Descriptor() ([]byte, []int) {
	return file_groups_api_proto_rawDescGZIP(), []int{4}
}

func (x *RemoveGroupFromFavoriteRequest) GetIds() []int32 {
	if x != nil {
		return x.Ids
	}
	return nil
}

var File_groups_api_proto protoreflect.FileDescriptor

var file_groups_api_proto_rawDesc = []byte{
	0x0a, 0x10, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a,
	0x13, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3d,
	0x0a, 0x13, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x22, 0x4b, 0x0a,
	0x0d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x22, 0x2b, 0x0a, 0x17, 0x41, 0x64,
	0x64, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x05, 0x52, 0x03, 0x69, 0x64, 0x73, 0x22, 0x32, 0x0a, 0x1e, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x03, 0x69, 0x64, 0x73, 0x32, 0xb7, 0x02, 0x0a, 0x0d,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x60, 0x0a,
	0x0b, 0x41, 0x64, 0x64, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x41,
	0x64, 0x64, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1f,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x3a, 0x01, 0x2a, 0x22, 0x14, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x12,
	0x70, 0x0a, 0x17, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x46, 0x72,
	0x6f, 0x6d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x1c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x3a, 0x01, 0x2a, 0x22, 0x11,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x6c, 0x65, 0x61, 0x76,
	0x65, 0x12, 0x52, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x73, 0x12, 0x14, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x10, 0x3a, 0x01, 0x2a, 0x22, 0x0b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x6f, 0x6f, 0x64, 0x68, 0x64, 0x73, 0x2f, 0x76, 0x6b, 0x2e, 0x70,
	0x6f, 0x73, 0x74, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x3b, 0x76, 0x6b, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_groups_api_proto_rawDescOnce sync.Once
	file_groups_api_proto_rawDescData = file_groups_api_proto_rawDesc
)

func file_groups_api_proto_rawDescGZIP() []byte {
	file_groups_api_proto_rawDescOnce.Do(func() {
		file_groups_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_groups_api_proto_rawDescData)
	})
	return file_groups_api_proto_rawDescData
}

var file_groups_api_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_groups_api_proto_goTypes = []interface{}{
	(*GetFavoritesRequest)(nil),            // 0: GetFavoritesRequest
	(*GetFavoriteResponse)(nil),            // 1: GetFavoriteResponse
	(*FavoriteGroup)(nil),                  // 2: FavoriteGroup
	(*AddFavoriteGroupRequest)(nil),        // 3: AddFavoriteGroupRequest
	(*RemoveGroupFromFavoriteRequest)(nil), // 4: RemoveGroupFromFavoriteRequest
	(*emptypb.Empty)(nil),                  // 5: google.protobuf.Empty
}
var file_groups_api_proto_depIdxs = []int32{
	2, // 0: GetFavoriteResponse.groups:type_name -> FavoriteGroup
	3, // 1: GroupsService.AddFavorite:input_type -> AddFavoriteGroupRequest
	4, // 2: GroupsService.RemoveGroupFromFavorite:input_type -> RemoveGroupFromFavoriteRequest
	0, // 3: GroupsService.GetFavorites:input_type -> GetFavoritesRequest
	5, // 4: GroupsService.AddFavorite:output_type -> google.protobuf.Empty
	5, // 5: GroupsService.RemoveGroupFromFavorite:output_type -> google.protobuf.Empty
	1, // 6: GroupsService.GetFavorites:output_type -> GetFavoriteResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_groups_api_proto_init() }
func file_groups_api_proto_init() {
	if File_groups_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_groups_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavoritesRequest); i {
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
		file_groups_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavoriteResponse); i {
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
		file_groups_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FavoriteGroup); i {
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
		file_groups_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddFavoriteGroupRequest); i {
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
		file_groups_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveGroupFromFavoriteRequest); i {
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
			RawDescriptor: file_groups_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_groups_api_proto_goTypes,
		DependencyIndexes: file_groups_api_proto_depIdxs,
		MessageInfos:      file_groups_api_proto_msgTypes,
	}.Build()
	File_groups_api_proto = out.File
	file_groups_api_proto_rawDesc = nil
	file_groups_api_proto_goTypes = nil
	file_groups_api_proto_depIdxs = nil
}
