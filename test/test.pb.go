// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: test/test.proto

package test

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

type Person struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Phone    string  `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	Siblings int32   `protobuf:"varint,3,opt,name=siblings,proto3" json:"siblings,omitempty"`
	Spouse   bool    `protobuf:"varint,4,opt,name=spouse,proto3" json:"spouse,omitempty"`
	Money    float32 `protobuf:"fixed32,5,opt,name=money,proto3" json:"money,omitempty"`
}

func (x *Person) Reset() {
	*x = Person{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Person) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person) ProtoMessage() {}

func (x *Person) ProtoReflect() protoreflect.Message {
	mi := &file_test_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person.ProtoReflect.Descriptor instead.
func (*Person) Descriptor() ([]byte, []int) {
	return file_test_test_proto_rawDescGZIP(), []int{0}
}

func (x *Person) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Person) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *Person) GetSiblings() int32 {
	if x != nil {
		return x.Siblings
	}
	return 0
}

func (x *Person) GetSpouse() bool {
	if x != nil {
		return x.Spouse
	}
	return false
}

func (x *Person) GetMoney() float32 {
	if x != nil {
		return x.Money
	}
	return 0
}

var File_test_test_proto protoreflect.FileDescriptor

var file_test_test_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x69, 0x6f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x22,
	0x7c, 0x0a, 0x06, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68,
	0x6f, 0x6e, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x69, 0x62, 0x6c, 0x69, 0x6e, 0x67, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x69, 0x62, 0x6c, 0x69, 0x6e, 0x67, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x70, 0x6f, 0x75, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x73, 0x70, 0x6f, 0x75, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x6e, 0x65, 0x79,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x42, 0x23, 0x5a,
	0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x62, 0x6e,
	0x79, 0x62, 0x65, 0x72, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x69, 0x6f, 0x2f, 0x74, 0x65,
	0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_test_proto_rawDescOnce sync.Once
	file_test_test_proto_rawDescData = file_test_test_proto_rawDesc
)

func file_test_test_proto_rawDescGZIP() []byte {
	file_test_test_proto_rawDescOnce.Do(func() {
		file_test_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_test_proto_rawDescData)
	})
	return file_test_test_proto_rawDescData
}

var file_test_test_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_test_test_proto_goTypes = []interface{}{
	(*Person)(nil), // 0: protoio.test.Person
}
var file_test_test_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_test_proto_init() }
func file_test_test_proto_init() {
	if File_test_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Person); i {
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
			RawDescriptor: file_test_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_test_proto_goTypes,
		DependencyIndexes: file_test_test_proto_depIdxs,
		MessageInfos:      file_test_test_proto_msgTypes,
	}.Build()
	File_test_test_proto = out.File
	file_test_test_proto_rawDesc = nil
	file_test_test_proto_goTypes = nil
	file_test_test_proto_depIdxs = nil
}
