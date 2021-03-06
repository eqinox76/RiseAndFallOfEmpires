// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/command.proto

package riseandfall

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Command struct {
	Empire               uint32           `protobuf:"varint,1,opt,name=empire,proto3" json:"empire,omitempty"`
	Orders               []*Command_Order `protobuf:"bytes,2,rep,name=orders" json:"orders,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Command) Reset()      { *m = Command{} }
func (*Command) ProtoMessage() {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_05a114ae39c7d553, []int{0}
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
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Command_Order) Reset()      { *m = Command_Order{} }
func (*Command_Order) ProtoMessage() {}
func (*Command_Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_05a114ae39c7d553, []int{0, 0}
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
	Equal(interface{}) bool
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
	XXX_sizecache        int32    `json:"-"`
}

func (m *MovementOrder) Reset()      { *m = MovementOrder{} }
func (*MovementOrder) ProtoMessage() {}
func (*MovementOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_command_05a114ae39c7d553, []int{1}
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
func (this *Command) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Command)
	if !ok {
		that2, ok := that.(Command)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Empire != that1.Empire {
		return false
	}
	if len(this.Orders) != len(that1.Orders) {
		return false
	}
	for i := range this.Orders {
		if !this.Orders[i].Equal(that1.Orders[i]) {
			return false
		}
	}
	return true
}
func (this *Command_Order) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Command_Order)
	if !ok {
		that2, ok := that.(Command_Order)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if that1.Order == nil {
		if this.Order != nil {
			return false
		}
	} else if this.Order == nil {
		return false
	} else if !this.Order.Equal(that1.Order) {
		return false
	}
	return true
}
func (this *Command_Order_Move) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Command_Order_Move)
	if !ok {
		that2, ok := that.(Command_Order_Move)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Move.Equal(that1.Move) {
		return false
	}
	return true
}
func (this *MovementOrder) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MovementOrder)
	if !ok {
		that2, ok := that.(MovementOrder)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Ship != that1.Ship {
		return false
	}
	if this.Start != that1.Start {
		return false
	}
	if this.Destination != that1.Destination {
		return false
	}
	return true
}
func (this *Command) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&riseandfall.Command{")
	s = append(s, "Empire: "+fmt.Sprintf("%#v", this.Empire)+",\n")
	if this.Orders != nil {
		s = append(s, "Orders: "+fmt.Sprintf("%#v", this.Orders)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Command_Order) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&riseandfall.Command_Order{")
	if this.Order != nil {
		s = append(s, "Order: "+fmt.Sprintf("%#v", this.Order)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Command_Order_Move) GoString() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&riseandfall.Command_Order_Move{` +
		`Move:` + fmt.Sprintf("%#v", this.Move) + `}`}, ", ")
	return s
}
func (this *MovementOrder) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&riseandfall.MovementOrder{")
	s = append(s, "Ship: "+fmt.Sprintf("%#v", this.Ship)+",\n")
	s = append(s, "Start: "+fmt.Sprintf("%#v", this.Start)+",\n")
	s = append(s, "Destination: "+fmt.Sprintf("%#v", this.Destination)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringCommand(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
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
	return n
}

func (m *Command_Order) Size() (n int) {
	var l int
	_ = l
	if m.Order != nil {
		n += m.Order.Size()
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
func (this *Command) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Command{`,
		`Empire:` + fmt.Sprintf("%v", this.Empire) + `,`,
		`Orders:` + strings.Replace(fmt.Sprintf("%v", this.Orders), "Command_Order", "Command_Order", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Command_Order) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Command_Order{`,
		`Order:` + fmt.Sprintf("%v", this.Order) + `,`,
		`}`,
	}, "")
	return s
}
func (this *Command_Order_Move) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Command_Order_Move{`,
		`Move:` + strings.Replace(fmt.Sprintf("%v", this.Move), "MovementOrder", "MovementOrder", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MovementOrder) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MovementOrder{`,
		`Ship:` + fmt.Sprintf("%v", this.Ship) + `,`,
		`Start:` + fmt.Sprintf("%v", this.Start) + `,`,
		`Destination:` + fmt.Sprintf("%v", this.Destination) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringCommand(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
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

func init() { proto.RegisterFile("proto/command.proto", fileDescriptor_command_05a114ae39c7d553) }

var fileDescriptor_command_05a114ae39c7d553 = []byte{
	// 264 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x31, 0x4e, 0xc3, 0x30,
	0x18, 0x85, 0xf3, 0xb7, 0x49, 0x2a, 0xfd, 0x51, 0x16, 0x83, 0x50, 0xd4, 0xe1, 0x57, 0xd4, 0xa9,
	0x03, 0x0a, 0xa8, 0xdc, 0x20, 0x2c, 0x2c, 0x08, 0x29, 0x2b, 0x53, 0x20, 0x46, 0x44, 0x6a, 0xec,
	0xc8, 0xb6, 0x3a, 0x73, 0x04, 0x56, 0x6e, 0xc0, 0x51, 0x18, 0x3b, 0x32, 0x12, 0xb3, 0x30, 0xf6,
	0x08, 0x88, 0xbf, 0x19, 0x8a, 0xba, 0xf9, 0x3d, 0x7f, 0xef, 0xf9, 0xc9, 0x78, 0xd2, 0x1b, 0xed,
	0xf4, 0xc5, 0xa3, 0xee, 0xba, 0x5a, 0x35, 0x05, 0x2b, 0x91, 0x98, 0xd6, 0xca, 0x5a, 0x35, 0x4f,
	0xf5, 0x7a, 0xbd, 0x78, 0x03, 0x9c, 0x5d, 0xef, 0xaf, 0xc5, 0x19, 0xc6, 0xb2, 0xeb, 0x5b, 0x23,
	0x33, 0xc8, 0x61, 0x99, 0x56, 0xa3, 0x12, 0x2b, 0x8c, 0xb5, 0x69, 0xa4, 0xb1, 0xd9, 0x24, 0x9f,
	0x2e, 0x93, 0xd5, 0xbc, 0x38, 0x68, 0x28, 0xc6, 0x74, 0x71, 0xf7, 0x87, 0x54, 0x23, 0x39, 0x2f,
	0x31, 0x62, 0x43, 0x5c, 0x62, 0xd8, 0xe9, 0x8d, 0xcc, 0xa6, 0x39, 0x1c, 0x45, 0x6f, 0xf5, 0x46,
	0x76, 0x52, 0x39, 0x26, 0x6f, 0x82, 0x8a, 0xc9, 0x72, 0x86, 0x11, 0x97, 0x2c, 0xee, 0x31, 0xfd,
	0x47, 0x08, 0x81, 0xa1, 0x7d, 0x6e, 0x7b, 0x9e, 0x17, 0x56, 0x7c, 0x16, 0xa7, 0x18, 0x59, 0x57,
	0x1b, 0x97, 0x4d, 0x78, 0xf3, 0x5e, 0x88, 0x1c, 0x93, 0x46, 0x5a, 0xd7, 0xaa, 0xda, 0xb5, 0x5a,
	0xf1, 0xe3, 0x69, 0x75, 0x68, 0x95, 0xe7, 0xdb, 0x81, 0x82, 0xcf, 0x81, 0x82, 0xdd, 0x40, 0xf0,
	0xe2, 0x09, 0xde, 0x3d, 0xc1, 0x87, 0x27, 0xd8, 0x7a, 0x82, 0x2f, 0x4f, 0xf0, 0xe3, 0x29, 0xd8,
	0x79, 0x82, 0xd7, 0x6f, 0x0a, 0x1e, 0x62, 0xfe, 0xba, 0xab, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xfe, 0x5c, 0xe3, 0x5c, 0x51, 0x01, 0x00, 0x00,
}
