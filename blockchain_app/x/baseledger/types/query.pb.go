// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: baseledger/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// this line is used by starport scaffolding # 3
type QueryGetBaseledgerTransactionRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *QueryGetBaseledgerTransactionRequest) Reset()         { *m = QueryGetBaseledgerTransactionRequest{} }
func (m *QueryGetBaseledgerTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetBaseledgerTransactionRequest) ProtoMessage()    {}
func (*QueryGetBaseledgerTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_27b223a452b70176, []int{0}
}
func (m *QueryGetBaseledgerTransactionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetBaseledgerTransactionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetBaseledgerTransactionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetBaseledgerTransactionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetBaseledgerTransactionRequest.Merge(m, src)
}
func (m *QueryGetBaseledgerTransactionRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetBaseledgerTransactionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetBaseledgerTransactionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetBaseledgerTransactionRequest proto.InternalMessageInfo

func (m *QueryGetBaseledgerTransactionRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type QueryGetBaseledgerTransactionResponse struct {
	BaseledgerTransaction *BaseledgerTransaction `protobuf:"bytes,1,opt,name=BaseledgerTransaction,proto3" json:"BaseledgerTransaction,omitempty"`
}

func (m *QueryGetBaseledgerTransactionResponse) Reset()         { *m = QueryGetBaseledgerTransactionResponse{} }
func (m *QueryGetBaseledgerTransactionResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetBaseledgerTransactionResponse) ProtoMessage()    {}
func (*QueryGetBaseledgerTransactionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_27b223a452b70176, []int{1}
}
func (m *QueryGetBaseledgerTransactionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetBaseledgerTransactionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetBaseledgerTransactionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetBaseledgerTransactionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetBaseledgerTransactionResponse.Merge(m, src)
}
func (m *QueryGetBaseledgerTransactionResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetBaseledgerTransactionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetBaseledgerTransactionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetBaseledgerTransactionResponse proto.InternalMessageInfo

func (m *QueryGetBaseledgerTransactionResponse) GetBaseledgerTransaction() *BaseledgerTransaction {
	if m != nil {
		return m.BaseledgerTransaction
	}
	return nil
}

type QueryAllBaseledgerTransactionRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllBaseledgerTransactionRequest) Reset()         { *m = QueryAllBaseledgerTransactionRequest{} }
func (m *QueryAllBaseledgerTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllBaseledgerTransactionRequest) ProtoMessage()    {}
func (*QueryAllBaseledgerTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_27b223a452b70176, []int{2}
}
func (m *QueryAllBaseledgerTransactionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllBaseledgerTransactionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllBaseledgerTransactionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllBaseledgerTransactionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllBaseledgerTransactionRequest.Merge(m, src)
}
func (m *QueryAllBaseledgerTransactionRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllBaseledgerTransactionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllBaseledgerTransactionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllBaseledgerTransactionRequest proto.InternalMessageInfo

func (m *QueryAllBaseledgerTransactionRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

type QueryAllBaseledgerTransactionResponse struct {
	BaseledgerTransaction []*BaseledgerTransaction `protobuf:"bytes,1,rep,name=BaseledgerTransaction,proto3" json:"BaseledgerTransaction,omitempty"`
	Pagination            *query.PageResponse      `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllBaseledgerTransactionResponse) Reset()         { *m = QueryAllBaseledgerTransactionResponse{} }
func (m *QueryAllBaseledgerTransactionResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllBaseledgerTransactionResponse) ProtoMessage()    {}
func (*QueryAllBaseledgerTransactionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_27b223a452b70176, []int{3}
}
func (m *QueryAllBaseledgerTransactionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllBaseledgerTransactionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllBaseledgerTransactionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllBaseledgerTransactionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllBaseledgerTransactionResponse.Merge(m, src)
}
func (m *QueryAllBaseledgerTransactionResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllBaseledgerTransactionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllBaseledgerTransactionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllBaseledgerTransactionResponse proto.InternalMessageInfo

func (m *QueryAllBaseledgerTransactionResponse) GetBaseledgerTransaction() []*BaseledgerTransaction {
	if m != nil {
		return m.BaseledgerTransaction
	}
	return nil
}

func (m *QueryAllBaseledgerTransactionResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryGetBaseledgerTransactionRequest)(nil), "unibrightio.baseledger.baseledger.QueryGetBaseledgerTransactionRequest")
	proto.RegisterType((*QueryGetBaseledgerTransactionResponse)(nil), "unibrightio.baseledger.baseledger.QueryGetBaseledgerTransactionResponse")
	proto.RegisterType((*QueryAllBaseledgerTransactionRequest)(nil), "unibrightio.baseledger.baseledger.QueryAllBaseledgerTransactionRequest")
	proto.RegisterType((*QueryAllBaseledgerTransactionResponse)(nil), "unibrightio.baseledger.baseledger.QueryAllBaseledgerTransactionResponse")
}

func init() { proto.RegisterFile("baseledger/query.proto", fileDescriptor_27b223a452b70176) }

var fileDescriptor_27b223a452b70176 = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4b, 0x4a, 0x2c, 0x4e,
	0xcd, 0x49, 0x4d, 0x49, 0x4f, 0x2d, 0xd2, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x52, 0x2c, 0xcd, 0xcb, 0x4c, 0x2a, 0xca, 0x4c, 0xcf, 0x28, 0xc9, 0xcc, 0xd7,
	0x43, 0xa8, 0x41, 0x62, 0x4a, 0xc9, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7, 0xa4, 0xea, 0x27, 0x16, 0x64,
	0xea, 0x27, 0xe6, 0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7, 0x15, 0x43, 0x0c, 0x90, 0xd2,
	0x4a, 0xce, 0x2f, 0xce, 0xcd, 0x2f, 0xd6, 0x07, 0x69, 0x80, 0x98, 0xac, 0x5f, 0x66, 0x98, 0x94,
	0x5a, 0x92, 0x68, 0xa8, 0x5f, 0x90, 0x98, 0x9e, 0x99, 0x07, 0x56, 0x0c, 0x55, 0xab, 0x86, 0xe4,
	0x08, 0x27, 0x38, 0x33, 0xa4, 0x28, 0x31, 0xaf, 0x38, 0x31, 0x19, 0xa1, 0x4e, 0xc9, 0x8c, 0x4b,
	0x25, 0x10, 0x64, 0x92, 0x7b, 0x6a, 0x09, 0x56, 0x65, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25,
	0x42, 0x7c, 0x5c, 0x4c, 0x99, 0x29, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x4c, 0x99, 0x29,
	0x4a, 0xd3, 0x19, 0xb9, 0x54, 0x09, 0x68, 0x2c, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x15, 0xca, 0xe3,
	0x12, 0xc5, 0xaa, 0x00, 0x6c, 0x18, 0xb7, 0x91, 0x85, 0x1e, 0xc1, 0x60, 0xd1, 0xc3, 0x6e, 0x01,
	0x76, 0x63, 0x95, 0xf2, 0xa0, 0x3e, 0x72, 0xcc, 0xc9, 0xc1, 0xeb, 0x23, 0x37, 0x2e, 0x2e, 0x44,
	0xa8, 0x41, 0x1d, 0xa3, 0xa6, 0x07, 0x09, 0x62, 0xb0, 0xe5, 0x7a, 0x90, 0xc8, 0x83, 0x06, 0xb1,
	0x5e, 0x40, 0x62, 0x7a, 0x2a, 0x54, 0x6f, 0x10, 0x92, 0x4e, 0xa5, 0x07, 0xb0, 0x90, 0xc0, 0x6d,
	0x21, 0xe1, 0x90, 0x60, 0xa6, 0x41, 0x48, 0x08, 0xb9, 0xa3, 0xf8, 0x90, 0x09, 0xec, 0x43, 0x75,
	0x82, 0x3e, 0x84, 0x38, 0x16, 0xd9, 0x8b, 0x46, 0x6f, 0x99, 0xb9, 0x58, 0xc1, 0x5e, 0x14, 0xfa,
	0xc8, 0xc8, 0x85, 0xcb, 0x32, 0x22, 0x5c, 0x4f, 0x4c, 0x4a, 0x93, 0xf2, 0xa0, 0xdc, 0x20, 0x88,
	0x17, 0x94, 0x5c, 0x9b, 0x2e, 0x3f, 0x99, 0xcc, 0x64, 0x2f, 0x64, 0xab, 0x8f, 0x64, 0xa2, 0x3e,
	0x52, 0xc6, 0x20, 0x94, 0x47, 0xf4, 0xab, 0x33, 0x53, 0x6a, 0x85, 0xde, 0x33, 0x72, 0x49, 0x60,
	0x95, 0x76, 0xcc, 0xc9, 0x21, 0xde, 0xdb, 0x04, 0x92, 0x23, 0xf1, 0xde, 0x26, 0x94, 0xcc, 0x94,
	0x1c, 0xc0, 0xde, 0xb6, 0x12, 0xb2, 0x20, 0xd7, 0xdb, 0x4e, 0x7e, 0x27, 0x1e, 0xc9, 0x31, 0x5e,
	0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31,
	0xdc, 0x78, 0x2c, 0xc7, 0x10, 0x65, 0x92, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f,
	0x8b, 0xcb, 0xf4, 0x0a, 0x64, 0x4e, 0x49, 0x65, 0x41, 0x6a, 0x71, 0x12, 0x1b, 0xb8, 0xac, 0x31,
	0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x0f, 0x17, 0x4a, 0xd8, 0x1a, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Queries a BaseledgerTransaction by id.
	BaseledgerTransaction(ctx context.Context, in *QueryGetBaseledgerTransactionRequest, opts ...grpc.CallOption) (*QueryGetBaseledgerTransactionResponse, error)
	// Queries a list of BaseledgerTransaction items.
	BaseledgerTransactionAll(ctx context.Context, in *QueryAllBaseledgerTransactionRequest, opts ...grpc.CallOption) (*QueryAllBaseledgerTransactionResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) BaseledgerTransaction(ctx context.Context, in *QueryGetBaseledgerTransactionRequest, opts ...grpc.CallOption) (*QueryGetBaseledgerTransactionResponse, error) {
	out := new(QueryGetBaseledgerTransactionResponse)
	err := c.cc.Invoke(ctx, "/unibrightio.baseledger.baseledger.Query/BaseledgerTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) BaseledgerTransactionAll(ctx context.Context, in *QueryAllBaseledgerTransactionRequest, opts ...grpc.CallOption) (*QueryAllBaseledgerTransactionResponse, error) {
	out := new(QueryAllBaseledgerTransactionResponse)
	err := c.cc.Invoke(ctx, "/unibrightio.baseledger.baseledger.Query/BaseledgerTransactionAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries a BaseledgerTransaction by id.
	BaseledgerTransaction(context.Context, *QueryGetBaseledgerTransactionRequest) (*QueryGetBaseledgerTransactionResponse, error)
	// Queries a list of BaseledgerTransaction items.
	BaseledgerTransactionAll(context.Context, *QueryAllBaseledgerTransactionRequest) (*QueryAllBaseledgerTransactionResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) BaseledgerTransaction(ctx context.Context, req *QueryGetBaseledgerTransactionRequest) (*QueryGetBaseledgerTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BaseledgerTransaction not implemented")
}
func (*UnimplementedQueryServer) BaseledgerTransactionAll(ctx context.Context, req *QueryAllBaseledgerTransactionRequest) (*QueryAllBaseledgerTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BaseledgerTransactionAll not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_BaseledgerTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetBaseledgerTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BaseledgerTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/unibrightio.baseledger.baseledger.Query/BaseledgerTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BaseledgerTransaction(ctx, req.(*QueryGetBaseledgerTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BaseledgerTransactionAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllBaseledgerTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BaseledgerTransactionAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/unibrightio.baseledger.baseledger.Query/BaseledgerTransactionAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BaseledgerTransactionAll(ctx, req.(*QueryAllBaseledgerTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "unibrightio.baseledger.baseledger.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BaseledgerTransaction",
			Handler:    _Query_BaseledgerTransaction_Handler,
		},
		{
			MethodName: "BaseledgerTransactionAll",
			Handler:    _Query_BaseledgerTransactionAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "baseledger/query.proto",
}

func (m *QueryGetBaseledgerTransactionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetBaseledgerTransactionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetBaseledgerTransactionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetBaseledgerTransactionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetBaseledgerTransactionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetBaseledgerTransactionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BaseledgerTransaction != nil {
		{
			size, err := m.BaseledgerTransaction.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllBaseledgerTransactionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllBaseledgerTransactionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllBaseledgerTransactionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllBaseledgerTransactionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllBaseledgerTransactionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllBaseledgerTransactionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.BaseledgerTransaction) > 0 {
		for iNdEx := len(m.BaseledgerTransaction) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.BaseledgerTransaction[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryGetBaseledgerTransactionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryGetBaseledgerTransactionResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.BaseledgerTransaction != nil {
		l = m.BaseledgerTransaction.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllBaseledgerTransactionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllBaseledgerTransactionResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.BaseledgerTransaction) > 0 {
		for _, e := range m.BaseledgerTransaction {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryGetBaseledgerTransactionRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetBaseledgerTransactionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetBaseledgerTransactionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetBaseledgerTransactionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetBaseledgerTransactionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetBaseledgerTransactionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseledgerTransaction", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BaseledgerTransaction == nil {
				m.BaseledgerTransaction = &BaseledgerTransaction{}
			}
			if err := m.BaseledgerTransaction.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAllBaseledgerTransactionRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllBaseledgerTransactionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllBaseledgerTransactionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAllBaseledgerTransactionResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllBaseledgerTransactionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllBaseledgerTransactionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseledgerTransaction", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaseledgerTransaction = append(m.BaseledgerTransaction, &BaseledgerTransaction{})
			if err := m.BaseledgerTransaction[len(m.BaseledgerTransaction)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
