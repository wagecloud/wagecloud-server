// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: account/v1/account.proto

package accountv1

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

// Account user message
type Account struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email         string                 `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Type          AccountType            `protobuf:"varint,5,opt,name=type,proto3,enum=account.v1.AccountType" json:"type,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     int64                  `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Account) Reset() {
	*x = Account{}
	mi := &file_account_v1_account_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_account_v1_account_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_account_v1_account_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Account) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Account) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Account) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Account) GetType() AccountType {
	if x != nil {
		return x.Type
	}
	return AccountType_ACCOUNT_TYPE_UNSPECIFIED
}

func (x *Account) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Account) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// Get account request
type GetUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            *int64                 `protobuf:"varint,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	Username      *string                `protobuf:"bytes,2,opt,name=username,proto3,oneof" json:"username,omitempty"`
	Email         *string                `protobuf:"bytes,3,opt,name=email,proto3,oneof" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	mi := &file_account_v1_account_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_account_v1_account_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_account_v1_account_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserRequest) GetId() int64 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *GetUserRequest) GetUsername() string {
	if x != nil && x.Username != nil {
		return *x.Username
	}
	return ""
}

func (x *GetUserRequest) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

// Account response message
type GetUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       *Account               `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	mi := &file_account_v1_account_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_account_v1_account_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserResponse.ProtoReflect.Descriptor instead.
func (*GetUserResponse) Descriptor() ([]byte, []int) {
	return file_account_v1_account_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserResponse) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

// Login request
type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            *int64                 `protobuf:"varint,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	Username      *string                `protobuf:"bytes,2,opt,name=username,proto3,oneof" json:"username,omitempty"`
	Email         *string                `protobuf:"bytes,3,opt,name=email,proto3,oneof" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_account_v1_account_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_account_v1_account_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_account_v1_account_proto_rawDescGZIP(), []int{3}
}

func (x *LoginRequest) GetId() int64 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *LoginRequest) GetUsername() string {
	if x != nil && x.Username != nil {
		return *x.Username
	}
	return ""
}

func (x *LoginRequest) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Login response
type LoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Account       *Account               `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_account_v1_account_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_account_v1_account_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_account_v1_account_proto_rawDescGZIP(), []int{4}
}

func (x *LoginResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *LoginResponse) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

// Register request
type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_account_v1_account_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_account_v1_account_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_account_v1_account_proto_rawDescGZIP(), []int{5}
}

func (x *RegisterRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Register response
type RegisterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Account       *Account               `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_account_v1_account_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_account_v1_account_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_account_v1_account_proto_rawDescGZIP(), []int{6}
}

func (x *RegisterResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *RegisterResponse) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

var File_account_v1_account_proto protoreflect.FileDescriptor

const file_account_v1_account_proto_rawDesc = "" +
	"\n" +
	"\x18account/v1/account.proto\x12\n" +
	"account.v1\x1a\x17account/v1/common.proto\"\xca\x01\n" +
	"\aAccount\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x1a\n" +
	"\busername\x18\x02 \x01(\tR\busername\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12\x14\n" +
	"\x05email\x18\x04 \x01(\tR\x05email\x12+\n" +
	"\x04type\x18\x05 \x01(\x0e2\x17.account.v1.AccountTypeR\x04type\x12\x1d\n" +
	"\n" +
	"created_at\x18\x06 \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\a \x01(\x03R\tupdatedAt\"\x7f\n" +
	"\x0eGetUserRequest\x12\x13\n" +
	"\x02id\x18\x01 \x01(\x03H\x00R\x02id\x88\x01\x01\x12\x1f\n" +
	"\busername\x18\x02 \x01(\tH\x01R\busername\x88\x01\x01\x12\x19\n" +
	"\x05email\x18\x03 \x01(\tH\x02R\x05email\x88\x01\x01B\x05\n" +
	"\x03_idB\v\n" +
	"\t_usernameB\b\n" +
	"\x06_email\"@\n" +
	"\x0fGetUserResponse\x12-\n" +
	"\aaccount\x18\x01 \x01(\v2\x13.account.v1.AccountR\aaccount\"\x99\x01\n" +
	"\fLoginRequest\x12\x13\n" +
	"\x02id\x18\x01 \x01(\x03H\x00R\x02id\x88\x01\x01\x12\x1f\n" +
	"\busername\x18\x02 \x01(\tH\x01R\busername\x88\x01\x01\x12\x19\n" +
	"\x05email\x18\x03 \x01(\tH\x02R\x05email\x88\x01\x01\x12\x1a\n" +
	"\bpassword\x18\x04 \x01(\tR\bpasswordB\x05\n" +
	"\x03_idB\v\n" +
	"\t_usernameB\b\n" +
	"\x06_email\"T\n" +
	"\rLoginResponse\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12-\n" +
	"\aaccount\x18\x02 \x01(\v2\x13.account.v1.AccountR\aaccount\"s\n" +
	"\x0fRegisterRequest\x12\x1a\n" +
	"\busername\x18\x01 \x01(\tR\busername\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\x12\x12\n" +
	"\x04name\x18\x04 \x01(\tR\x04name\"W\n" +
	"\x10RegisterResponse\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12-\n" +
	"\aaccount\x18\x02 \x01(\v2\x13.account.v1.AccountR\aaccount2\xdf\x01\n" +
	"\x0eAccountService\x12D\n" +
	"\aGetUser\x12\x1a.account.v1.GetUserRequest\x1a\x1b.account.v1.GetUserResponse\"\x00\x12>\n" +
	"\x05Login\x12\x18.account.v1.LoginRequest\x1a\x19.account.v1.LoginResponse\"\x00\x12G\n" +
	"\bRegister\x12\x1b.account.v1.RegisterRequest\x1a\x1c.account.v1.RegisterResponse\"\x00B\xaa\x01\n" +
	"\x0ecom.account.v1B\fAccountProtoP\x01ZAgithub.com/wagecloud/wagecloud-server/gen/pb/account/v1;accountv1\xa2\x02\x03AXX\xaa\x02\n" +
	"Account.V1\xca\x02\n" +
	"Account\\V1\xe2\x02\x16Account\\V1\\GPBMetadata\xea\x02\vAccount::V1b\x06proto3"

var (
	file_account_v1_account_proto_rawDescOnce sync.Once
	file_account_v1_account_proto_rawDescData []byte
)

func file_account_v1_account_proto_rawDescGZIP() []byte {
	file_account_v1_account_proto_rawDescOnce.Do(func() {
		file_account_v1_account_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_account_v1_account_proto_rawDesc), len(file_account_v1_account_proto_rawDesc)))
	})
	return file_account_v1_account_proto_rawDescData
}

var file_account_v1_account_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_account_v1_account_proto_goTypes = []any{
	(*Account)(nil),          // 0: account.v1.Account
	(*GetUserRequest)(nil),   // 1: account.v1.GetUserRequest
	(*GetUserResponse)(nil),  // 2: account.v1.GetUserResponse
	(*LoginRequest)(nil),     // 3: account.v1.LoginRequest
	(*LoginResponse)(nil),    // 4: account.v1.LoginResponse
	(*RegisterRequest)(nil),  // 5: account.v1.RegisterRequest
	(*RegisterResponse)(nil), // 6: account.v1.RegisterResponse
	(AccountType)(0),         // 7: account.v1.AccountType
}
var file_account_v1_account_proto_depIdxs = []int32{
	7, // 0: account.v1.Account.type:type_name -> account.v1.AccountType
	0, // 1: account.v1.GetUserResponse.account:type_name -> account.v1.Account
	0, // 2: account.v1.LoginResponse.account:type_name -> account.v1.Account
	0, // 3: account.v1.RegisterResponse.account:type_name -> account.v1.Account
	1, // 4: account.v1.AccountService.GetUser:input_type -> account.v1.GetUserRequest
	3, // 5: account.v1.AccountService.Login:input_type -> account.v1.LoginRequest
	5, // 6: account.v1.AccountService.Register:input_type -> account.v1.RegisterRequest
	2, // 7: account.v1.AccountService.GetUser:output_type -> account.v1.GetUserResponse
	4, // 8: account.v1.AccountService.Login:output_type -> account.v1.LoginResponse
	6, // 9: account.v1.AccountService.Register:output_type -> account.v1.RegisterResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_account_v1_account_proto_init() }
func file_account_v1_account_proto_init() {
	if File_account_v1_account_proto != nil {
		return
	}
	file_account_v1_common_proto_init()
	file_account_v1_account_proto_msgTypes[1].OneofWrappers = []any{}
	file_account_v1_account_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_account_v1_account_proto_rawDesc), len(file_account_v1_account_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_account_v1_account_proto_goTypes,
		DependencyIndexes: file_account_v1_account_proto_depIdxs,
		MessageInfos:      file_account_v1_account_proto_msgTypes,
	}.Build()
	File_account_v1_account_proto = out.File
	file_account_v1_account_proto_goTypes = nil
	file_account_v1_account_proto_depIdxs = nil
}
