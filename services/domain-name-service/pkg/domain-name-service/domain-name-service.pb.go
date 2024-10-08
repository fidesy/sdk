// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.3
// source: api/domain-name-service/domain-name-service.proto

package domain_name_service

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

type GetAddressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
}

func (x *GetAddressRequest) Reset() {
	*x = GetAddressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAddressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAddressRequest) ProtoMessage() {}

func (x *GetAddressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAddressRequest.ProtoReflect.Descriptor instead.
func (*GetAddressRequest) Descriptor() ([]byte, []int) {
	return file_api_domain_name_service_domain_name_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAddressRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

type GetAddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *GetAddressResponse) Reset() {
	*x = GetAddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAddressResponse) ProtoMessage() {}

func (x *GetAddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAddressResponse.ProtoReflect.Descriptor instead.
func (*GetAddressResponse) Descriptor() ([]byte, []int) {
	return file_api_domain_name_service_domain_name_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetAddressResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type UpdateAddressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName string `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Address     string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *UpdateAddressRequest) Reset() {
	*x = UpdateAddressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAddressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAddressRequest) ProtoMessage() {}

func (x *UpdateAddressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAddressRequest.ProtoReflect.Descriptor instead.
func (*UpdateAddressRequest) Descriptor() ([]byte, []int) {
	return file_api_domain_name_service_domain_name_service_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateAddressRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *UpdateAddressRequest) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type UpdateAddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateAddressResponse) Reset() {
	*x = UpdateAddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAddressResponse) ProtoMessage() {}

func (x *UpdateAddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_domain_name_service_domain_name_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAddressResponse.ProtoReflect.Descriptor instead.
func (*UpdateAddressResponse) Descriptor() ([]byte, []int) {
	return file_api_domain_name_service_domain_name_service_proto_rawDescGZIP(), []int{3}
}

var File_api_domain_name_service_domain_name_service_proto protoreflect.FileDescriptor

var file_api_domain_name_service_domain_name_service_proto_rawDesc = []byte{
	0x0a, 0x31, 0x61, 0x70, 0x69, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x6e, 0x61, 0x6d,
	0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x2d, 0x6e, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x13, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0x36, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a,
	0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x2e, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x22, 0x53, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x17, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xda,
	0x01, 0x0a, 0x11, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x5d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x26, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x66, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x29, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2a, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3f, 0x5a, 0x3d, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x64, 0x65, 0x73, 0x79,
	0x2d, 0x70, 0x61, 0x79, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x6e, 0x61, 0x6d, 0x65,
	0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3b, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_domain_name_service_domain_name_service_proto_rawDescOnce sync.Once
	file_api_domain_name_service_domain_name_service_proto_rawDescData = file_api_domain_name_service_domain_name_service_proto_rawDesc
)

func file_api_domain_name_service_domain_name_service_proto_rawDescGZIP() []byte {
	file_api_domain_name_service_domain_name_service_proto_rawDescOnce.Do(func() {
		file_api_domain_name_service_domain_name_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_domain_name_service_domain_name_service_proto_rawDescData)
	})
	return file_api_domain_name_service_domain_name_service_proto_rawDescData
}

var file_api_domain_name_service_domain_name_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_domain_name_service_domain_name_service_proto_goTypes = []interface{}{
	(*GetAddressRequest)(nil),     // 0: domain_name_service.GetAddressRequest
	(*GetAddressResponse)(nil),    // 1: domain_name_service.GetAddressResponse
	(*UpdateAddressRequest)(nil),  // 2: domain_name_service.UpdateAddressRequest
	(*UpdateAddressResponse)(nil), // 3: domain_name_service.UpdateAddressResponse
}
var file_api_domain_name_service_domain_name_service_proto_depIdxs = []int32{
	0, // 0: domain_name_service.DomainNameService.GetAddress:input_type -> domain_name_service.GetAddressRequest
	2, // 1: domain_name_service.DomainNameService.UpdateAddress:input_type -> domain_name_service.UpdateAddressRequest
	1, // 2: domain_name_service.DomainNameService.GetAddress:output_type -> domain_name_service.GetAddressResponse
	3, // 3: domain_name_service.DomainNameService.UpdateAddress:output_type -> domain_name_service.UpdateAddressResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_domain_name_service_domain_name_service_proto_init() }
func file_api_domain_name_service_domain_name_service_proto_init() {
	if File_api_domain_name_service_domain_name_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_domain_name_service_domain_name_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAddressRequest); i {
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
		file_api_domain_name_service_domain_name_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAddressResponse); i {
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
		file_api_domain_name_service_domain_name_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAddressRequest); i {
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
		file_api_domain_name_service_domain_name_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAddressResponse); i {
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
			RawDescriptor: file_api_domain_name_service_domain_name_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_domain_name_service_domain_name_service_proto_goTypes,
		DependencyIndexes: file_api_domain_name_service_domain_name_service_proto_depIdxs,
		MessageInfos:      file_api_domain_name_service_domain_name_service_proto_msgTypes,
	}.Build()
	File_api_domain_name_service_domain_name_service_proto = out.File
	file_api_domain_name_service_domain_name_service_proto_rawDesc = nil
	file_api_domain_name_service_domain_name_service_proto_goTypes = nil
	file_api_domain_name_service_domain_name_service_proto_depIdxs = nil
}
