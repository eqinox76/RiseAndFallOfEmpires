// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/command.proto

package riseandfall

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

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
	Empire               uint32           `protobuf:"varint,1,opt,name=empire,proto3" json:"empire,omitempty"`
	Orders               []*Command_Order `protobuf:"bytes,2,rep,name=orders" json:"orders,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_b3ea3e935678e99b, []int{0}
}
func (m *Command) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Command.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(dst, src)
}
func (m *Command) XXX_Size() int {
	return m.Size()
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

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
	Order                isCommand_Order_Order `protobuf_oneof:"order"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Command_Order) Reset()         { *m = Command_Order{} }
func (m *Command_Order) String() string { return proto.CompactTextString(m) }
func (*Command_Order) ProtoMessage()    {}
func (*Command_Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_b3ea3e935678e99b, []int{0, 0}
}
func (m *Command_Order) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Command_Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Command_Order.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Command_Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command_Order.Merge(dst, src)
}
func (m *Command_Order) XXX_Size() int {
	return m.Size()
}
func (m *Command_Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Command_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Command_Order proto.InternalMessageInfo

type isCommand_Order_Order interface {
	isCommand_Order_Order()
	MarshalTo([]byte) (int, error)
	Size() int
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
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
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
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type MovementOrder struct {
	Ship                 uint64   `protobuf:"varint,1,opt,name=ship,proto3" json:"ship,omitempty"`
	Start                uint32   `protobuf:"varint,2,opt,name=start,proto3" json:"start,omitempty"`
	Destination          uint32   `protobuf:"varint,3,opt,name=destination,proto3" json:"destination,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MovementOrder) Reset()         { *m = MovementOrder{} }
func (m *MovementOrder) String() string { return proto.CompactTextString(m) }
func (*MovementOrder) ProtoMessage()    {}
func (*MovementOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_b3ea3e935678e99b, []int{1}
}
func (m *MovementOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MovementOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MovementOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *MovementOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MovementOrder.Merge(dst, src)
}
func (m *MovementOrder) XXX_Size() int {
	return m.Size()
}
func (m *MovementOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_MovementOrder.DiscardUnknown(m)
}

var xxx_messageInfo_MovementOrder proto.InternalMessageInfo

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
func (m *Command) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Command) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Empire != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCommand(dAtA, i, uint64(m.Empire))
	}
	if len(m.Orders) > 0 {
		for _, msg := range m.Orders {
			dAtA[i] = 0x12
			i++
			i = encodeVarintCommand(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Command_Order) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Command_Order) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Order != nil {
		nn1, err := m.Order.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Command_Order_Move) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Move != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCommand(dAtA, i, uint64(m.Move.Size()))
		n2, err := m.Move.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}
func (m *MovementOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MovementOrder) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Ship != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCommand(dAtA, i, uint64(m.Ship))
	}
	if m.Start != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintCommand(dAtA, i, uint64(m.Start))
	}
	if m.Destination != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintCommand(dAtA, i, uint64(m.Destination))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintCommand(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Command) Size() (n int) {
	var l int
	_ = l
	if m.Empire != 0 {
		n += 1 + sovCommand(uint64(m.Empire))
	}
	if len(m.Orders) > 0 {
		for _, e := range m.Orders {
			l = e.Size()
			n += 1 + l + sovCommand(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Command_Order) Size() (n int) {
	var l int
	_ = l
	if m.Order != nil {
		n += m.Order.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Command_Order_Move) Size() (n int) {
	var l int
	_ = l
	if m.Move != nil {
		l = m.Move.Size()
		n += 1 + l + sovCommand(uint64(l))
	}
	return n
}
func (m *MovementOrder) Size() (n int) {
	var l int
	_ = l
	if m.Ship != 0 {
		n += 1 + sovCommand(uint64(m.Ship))
	}
	if m.Start != 0 {
		n += 1 + sovCommand(uint64(m.Start))
	}
	if m.Destination != 0 {
		n += 1 + sovCommand(uint64(m.Destination))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovCommand(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCommand(x uint64) (n int) {
	return sovCommand(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Command) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Command: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Command: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Empire", wireType)
			}
			m.Empire = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Empire |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Orders", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Orders = append(m.Orders, &Command_Order{})
			if err := m.Orders[len(m.Orders)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Command_Order) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Order: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Order: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Move", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthCommand
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &MovementOrder{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Order = &Command_Order_Move{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MovementOrder) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommand
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MovementOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MovementOrder: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ship", wireType)
			}
			m.Ship = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Ship |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Start", wireType)
			}
			m.Start = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Start |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Destination", wireType)
			}
			m.Destination = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Destination |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCommand(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommand
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCommand(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommand
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
					return 0, ErrIntOverflowCommand
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommand
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthCommand
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCommand
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipCommand(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthCommand = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommand   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("proto/command.proto", fileDescriptor_command_b3ea3e935678e99b) }

var fileDescriptor_command_b3ea3e935678e99b = []byte{
	// 235 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xce, 0xcf, 0xcd, 0x4d, 0xcc, 0x4b, 0xd1, 0x03, 0xf3, 0x84, 0xb8, 0x8b, 0x32,
	0x8b, 0x53, 0x13, 0xf3, 0x52, 0xd2, 0x12, 0x73, 0x72, 0x94, 0x66, 0x32, 0x72, 0xb1, 0x3b, 0x43,
	0xa4, 0x85, 0xc4, 0xb8, 0xd8, 0x52, 0x73, 0x0b, 0x32, 0x8b, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35,
	0x78, 0x83, 0xa0, 0x3c, 0x21, 0x23, 0x2e, 0xb6, 0xfc, 0xa2, 0x94, 0xd4, 0xa2, 0x62, 0x09, 0x26,
	0x05, 0x66, 0x0d, 0x6e, 0x23, 0x29, 0x3d, 0x24, 0x13, 0xf4, 0xa0, 0xba, 0xf5, 0xfc, 0x41, 0x4a,
	0x82, 0xa0, 0x2a, 0xa5, 0x9c, 0xb8, 0x58, 0xc1, 0x02, 0x42, 0x06, 0x5c, 0x2c, 0xb9, 0xf9, 0x65,
	0xa9, 0x12, 0xcc, 0x0a, 0x8c, 0x18, 0x5a, 0x7d, 0xf3, 0xcb, 0x52, 0x73, 0x53, 0xf3, 0x4a, 0xc0,
	0x2a, 0x3d, 0x18, 0x82, 0xc0, 0x2a, 0x9d, 0xd8, 0xb9, 0x58, 0xc1, 0x86, 0x28, 0x45, 0x73, 0xf1,
	0xa2, 0xa8, 0x10, 0x12, 0xe2, 0x62, 0x29, 0xce, 0xc8, 0x2c, 0x00, 0x3b, 0x8f, 0x25, 0x08, 0xcc,
	0x16, 0x12, 0xe1, 0x62, 0x2d, 0x2e, 0x49, 0x2c, 0x2a, 0x91, 0x60, 0x02, 0xbb, 0x19, 0xc2, 0x11,
	0x52, 0xe0, 0xe2, 0x4e, 0x49, 0x2d, 0x2e, 0xc9, 0xcc, 0x4b, 0x2c, 0xc9, 0xcc, 0xcf, 0x03, 0x5b,
	0xce, 0x1b, 0x84, 0x2c, 0xe4, 0x24, 0x70, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f,
	0x1e, 0xc9, 0x31, 0xce, 0x78, 0x2c, 0xc7, 0x90, 0xc4, 0x06, 0x0e, 0x1e, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x16, 0x7e, 0xc5, 0xdc, 0x35, 0x01, 0x00, 0x00,
}
