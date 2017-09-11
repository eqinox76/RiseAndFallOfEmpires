// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/command.proto

/*
Package riseandfall is a generated protocol buffer package.

It is generated from these files:
	proto/command.proto
	proto/space.proto

It has these top-level messages:
	Command
	MovementOrder
	Space
	Planet
	Ship
	Movement
	Empire
*/
package riseandfall

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Command struct {
	Empire uint32           `protobuf:"varint,1,opt,name=empire" json:"empire,omitempty"`
	Orders []*Command_Order `protobuf:"bytes,2,rep,name=orders" json:"orders,omitempty"`
}

func (m *Command) Reset()                    { *m = Command{} }
func (m *Command) String() string            { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()               {}
func (*Command) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Command) GetEmpire() uint32 {
	if m != nil {
		return m.Empire
	}
	return 0
}

func (m *Command) GetOrders() []*Command_Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type Command_Order struct {
	// Types that are valid to be assigned to Order:
	//	*Command_Order_Move
	Order isCommand_Order_Order `protobuf_oneof:"order"`
}

func (m *Command_Order) Reset()                    { *m = Command_Order{} }
func (m *Command_Order) String() string            { return proto.CompactTextString(m) }
func (*Command_Order) ProtoMessage()               {}
func (*Command_Order) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type isCommand_Order_Order interface {
	isCommand_Order_Order()
}

type Command_Order_Move struct {
	Move *MovementOrder `protobuf:"bytes,3,opt,name=move,oneof"`
}

func (*Command_Order_Move) isCommand_Order_Order() {}

func (m *Command_Order) GetOrder() isCommand_Order_Order {
	if m != nil {
		return m.Order
	}
	return nil
}

func (m *Command_Order) GetMove() *MovementOrder {
	if x, ok := m.GetOrder().(*Command_Order_Move); ok {
		return x.Move
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Command_Order) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Command_Order_OneofMarshaler, _Command_Order_OneofUnmarshaler, _Command_Order_OneofSizer, []interface{}{
		(*Command_Order_Move)(nil),
	}
}

func _Command_Order_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Command_Order)
	// order
	switch x := m.Order.(type) {
	case *Command_Order_Move:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Move); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Command_Order.Order has unexpected type %T", x)
	}
	return nil
}

func _Command_Order_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Command_Order)
	switch tag {
	case 3: // order.move
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MovementOrder)
		err := b.DecodeMessage(msg)
		m.Order = &Command_Order_Move{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Command_Order_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Command_Order)
	// order
	switch x := m.Order.(type) {
	case *Command_Order_Move:
		s := proto.Size(x.Move)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type MovementOrder struct {
	Ship        uint64 `protobuf:"varint,1,opt,name=ship" json:"ship,omitempty"`
	Start       uint32 `protobuf:"varint,2,opt,name=start" json:"start,omitempty"`
	Destination uint32 `protobuf:"varint,3,opt,name=destination" json:"destination,omitempty"`
}

func (m *MovementOrder) Reset()                    { *m = MovementOrder{} }
func (m *MovementOrder) String() string            { return proto.CompactTextString(m) }
func (*MovementOrder) ProtoMessage()               {}
func (*MovementOrder) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MovementOrder) GetShip() uint64 {
	if m != nil {
		return m.Ship
	}
	return 0
}

func (m *MovementOrder) GetStart() uint32 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *MovementOrder) GetDestination() uint32 {
	if m != nil {
		return m.Destination
	}
	return 0
}

func init() {
	proto.RegisterType((*Command)(nil), "riseandfall.Command")
	proto.RegisterType((*Command_Order)(nil), "riseandfall.Command.Order")
	proto.RegisterType((*MovementOrder)(nil), "riseandfall.MovementOrder")
}

func init() { proto.RegisterFile("proto/command.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x14, 0x84, 0x71, 0x9b, 0xa4, 0xd2, 0x8b, 0xb2, 0x3c, 0x10, 0xb2, 0x3a, 0x59, 0x9d, 0x3a, 0x19,
	0x54, 0xfe, 0x41, 0x59, 0x58, 0x10, 0x92, 0x57, 0x26, 0x83, 0x1f, 0xc2, 0x52, 0x6d, 0x47, 0xb6,
	0x95, 0xff, 0xc2, 0xbf, 0x45, 0xbc, 0x64, 0x08, 0xea, 0xe6, 0x3b, 0x7f, 0x77, 0x3e, 0x19, 0x6e,
	0xc7, 0x9c, 0x6a, 0x7a, 0xf8, 0x4c, 0x21, 0xd8, 0xe8, 0x34, 0x2b, 0xec, 0xb3, 0x2f, 0x64, 0xa3,
	0xfb, 0xb2, 0x97, 0xcb, 0xe1, 0x47, 0xc0, 0xee, 0x79, 0xbe, 0xc6, 0x7b, 0xe8, 0x28, 0x8c, 0x3e,
	0x93, 0x14, 0x4a, 0x1c, 0x07, 0xb3, 0x28, 0x3c, 0x41, 0x97, 0xb2, 0xa3, 0x5c, 0xe4, 0x46, 0x6d,
	0x8f, 0xfd, 0x69, 0xaf, 0x57, 0x0d, 0x7a, 0x49, 0xeb, 0xb7, 0x3f, 0xc4, 0x2c, 0xe4, 0xfe, 0x0c,
	0x2d, 0x1b, 0xf8, 0x08, 0x4d, 0x48, 0x13, 0xc9, 0xad, 0x12, 0x57, 0xd1, 0xd7, 0x34, 0x51, 0xa0,
	0x58, 0x99, 0x7c, 0xb9, 0x31, 0x4c, 0x9e, 0x77, 0xd0, 0x72, 0xc9, 0xe1, 0x1d, 0x86, 0x7f, 0x04,
	0x22, 0x34, 0xe5, 0xdb, 0x8f, 0x3c, 0xaf, 0x31, 0x7c, 0xc6, 0x3b, 0x68, 0x4b, 0xb5, 0xb9, 0xca,
	0x0d, 0x6f, 0x9e, 0x05, 0x2a, 0xe8, 0x1d, 0x95, 0xea, 0xa3, 0xad, 0x3e, 0x45, 0x7e, 0x7c, 0x30,
	0x6b, 0xeb, 0xa3, 0xe3, 0xcf, 0x78, 0xfa, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x7b, 0xd2, 0x79,
	0x23, 0x01, 0x00, 0x00,
}