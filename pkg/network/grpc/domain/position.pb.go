// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/network/protocol/domain/position.proto

package domain

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Position struct {
	X                    float64  `protobuf:"fixed64,1,opt,name=X,proto3" json:"X,omitempty"`
	Y                    float64  `protobuf:"fixed64,2,opt,name=Y,proto3" json:"Y,omitempty"`
	Z                    float64  `protobuf:"fixed64,3,opt,name=Z,proto3" json:"Z,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f1bd2d837e5d7e1, []int{0}
}

func (m *Position) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Position.Unmarshal(m, b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Position.Marshal(b, m, deterministic)
}
func (m *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(m, src)
}
func (m *Position) XXX_Size() int {
	return xxx_messageInfo_Position.Size(m)
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Position) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Position) GetZ() float64 {
	if m != nil {
		return m.Z
	}
	return 0
}

func init() {
	proto.RegisterType((*Position)(nil), "domain.Position")
}

func init() {
	proto.RegisterFile("pkg/network/protocol/domain/position.proto", fileDescriptor_0f1bd2d837e5d7e1)
}

var fileDescriptor_0f1bd2d837e5d7e1 = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2a, 0xc8, 0x4e, 0xd7,
	0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf,
	0xd1, 0x4f, 0xc9, 0xcf, 0x4d, 0xcc, 0xcc, 0xd3, 0x2f, 0xc8, 0x2f, 0xce, 0x2c, 0xc9, 0xcc, 0xcf,
	0xd3, 0x03, 0x4b, 0x08, 0xb1, 0x41, 0x84, 0x95, 0x4c, 0xb8, 0x38, 0x02, 0xa0, 0x32, 0x42, 0x3c,
	0x5c, 0x8c, 0x11, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x8c, 0x41, 0x8c, 0x11, 0x20, 0x5e, 0xa4, 0x04,
	0x13, 0x84, 0x17, 0x09, 0xe2, 0x45, 0x49, 0x30, 0x43, 0x78, 0x51, 0x4e, 0x26, 0x51, 0x46, 0xe9,
	0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x79, 0x99, 0xc9, 0xd9, 0x49, 0x45,
	0x95, 0x89, 0x79, 0xfa, 0xe9, 0xf9, 0x65, 0xf9, 0xa9, 0xfa, 0xc8, 0xce, 0x48, 0x2f, 0x2a, 0x48,
	0x86, 0x3a, 0x21, 0x89, 0x0d, 0x6c, 0xb5, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x57, 0xdb, 0x99,
	0xc2, 0xa8, 0x00, 0x00, 0x00,
}