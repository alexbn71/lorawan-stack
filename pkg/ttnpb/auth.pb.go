// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/auth.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strconv "strconv"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// GrantType enum defines the OAuth2 flows a third-party client can use to get
// access to a token.
type GrantType int32

const (
	// Grant type used to exchange an authorization code for an access token.
	GRANT_AUTHORIZATION_CODE GrantType = 0
	// Grant type used to exchange an user ID and password for an access token.
	GRANT_PASSWORD GrantType = 1
	// Grant type used to exchange a refresh token for an access token.
	GRANT_REFRESH_TOKEN GrantType = 2
)

var GrantType_name = map[int32]string{
	0: "GRANT_AUTHORIZATION_CODE",
	1: "GRANT_PASSWORD",
	2: "GRANT_REFRESH_TOKEN",
}
var GrantType_value = map[string]int32{
	"GRANT_AUTHORIZATION_CODE": 0,
	"GRANT_PASSWORD":           1,
	"GRANT_REFRESH_TOKEN":      2,
}

func (GrantType) EnumDescriptor() ([]byte, []int) { return fileDescriptorAuth, []int{0} }

type APIKey struct {
	// key is the API key itself which is generated by the Identity Server.
	// It is unique and immutable.
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// name is the immutable and unique (in the namespace of the entity) API key name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// rights are the rights this API key bears.
	Rights []Right `protobuf:"varint,3,rep,packed,name=rights,enum=ttn.v3.Right" json:"rights,omitempty"`
}

func (m *APIKey) Reset()                    { *m = APIKey{} }
func (m *APIKey) String() string            { return proto.CompactTextString(m) }
func (*APIKey) ProtoMessage()               {}
func (*APIKey) Descriptor() ([]byte, []int) { return fileDescriptorAuth, []int{0} }

func (m *APIKey) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *APIKey) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *APIKey) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

func init() {
	proto.RegisterType((*APIKey)(nil), "ttn.v3.APIKey")
	golang_proto.RegisterType((*APIKey)(nil), "ttn.v3.APIKey")
	proto.RegisterEnum("ttn.v3.GrantType", GrantType_name, GrantType_value)
	golang_proto.RegisterEnum("ttn.v3.GrantType", GrantType_name, GrantType_value)
}
func (x GrantType) String() string {
	s, ok := GrantType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *APIKey) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*APIKey)
	if !ok {
		that2, ok := that.(APIKey)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *APIKey")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *APIKey but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *APIKey but is not nil && this == nil")
	}
	if this.Key != that1.Key {
		return fmt.Errorf("Key this(%v) Not Equal that(%v)", this.Key, that1.Key)
	}
	if this.Name != that1.Name {
		return fmt.Errorf("Name this(%v) Not Equal that(%v)", this.Name, that1.Name)
	}
	if len(this.Rights) != len(that1.Rights) {
		return fmt.Errorf("Rights this(%v) Not Equal that(%v)", len(this.Rights), len(that1.Rights))
	}
	for i := range this.Rights {
		if this.Rights[i] != that1.Rights[i] {
			return fmt.Errorf("Rights this[%v](%v) Not Equal that[%v](%v)", i, this.Rights[i], i, that1.Rights[i])
		}
	}
	return nil
}
func (this *APIKey) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*APIKey)
	if !ok {
		that2, ok := that.(APIKey)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Key != that1.Key {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if len(this.Rights) != len(that1.Rights) {
		return false
	}
	for i := range this.Rights {
		if this.Rights[i] != that1.Rights[i] {
			return false
		}
	}
	return true
}
func (m *APIKey) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *APIKey) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAuth(dAtA, i, uint64(len(m.Key)))
		i += copy(dAtA[i:], m.Key)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAuth(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Rights) > 0 {
		dAtA2 := make([]byte, len(m.Rights)*10)
		var j1 int
		for _, num := range m.Rights {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		dAtA[i] = 0x1a
		i++
		i = encodeVarintAuth(dAtA, i, uint64(j1))
		i += copy(dAtA[i:], dAtA2[:j1])
	}
	return i, nil
}

func encodeVarintAuth(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedAPIKey(r randyAuth, easy bool) *APIKey {
	this := &APIKey{}
	this.Key = randStringAuth(r)
	this.Name = randStringAuth(r)
	v1 := r.Intn(10)
	this.Rights = make([]Right, v1)
	for i := 0; i < v1; i++ {
		this.Rights[i] = Right([]int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 17, 18, 31, 32, 33, 34, 35, 36, 37, 38, 39, 51, 52, 53, 54, 55, 56, 57, 58, 71, 72, 73, 74, 75, 76, 77, 78, 79}[r.Intn(40)])
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyAuth interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneAuth(r randyAuth) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringAuth(r randyAuth) string {
	v2 := r.Intn(100)
	tmps := make([]rune, v2)
	for i := 0; i < v2; i++ {
		tmps[i] = randUTF8RuneAuth(r)
	}
	return string(tmps)
}
func randUnrecognizedAuth(r randyAuth, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldAuth(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldAuth(dAtA []byte, r randyAuth, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateAuth(dAtA, uint64(key))
		v3 := r.Int63()
		if r.Intn(2) == 0 {
			v3 *= -1
		}
		dAtA = encodeVarintPopulateAuth(dAtA, uint64(v3))
	case 1:
		dAtA = encodeVarintPopulateAuth(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateAuth(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateAuth(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateAuth(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateAuth(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *APIKey) Size() (n int) {
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovAuth(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovAuth(uint64(l))
	}
	if len(m.Rights) > 0 {
		l = 0
		for _, e := range m.Rights {
			l += sovAuth(uint64(e))
		}
		n += 1 + sovAuth(uint64(l)) + l
	}
	return n
}

func sovAuth(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozAuth(x uint64) (n int) {
	return sovAuth((x << 1) ^ uint64((int64(x) >> 63)))
}
func (m *APIKey) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuth
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
			return fmt.Errorf("proto: APIKey: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: APIKey: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuth
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuth
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuth
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAuth
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v Right
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowAuth
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (Right(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Rights = append(m.Rights, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowAuth
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthAuth
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v Right
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAuth
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (Right(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Rights = append(m.Rights, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Rights", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAuth(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAuth
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
func skipAuth(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAuth
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
					return 0, ErrIntOverflowAuth
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
					return 0, ErrIntOverflowAuth
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
				return 0, ErrInvalidLengthAuth
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowAuth
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
				next, err := skipAuth(dAtA[start:])
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
	ErrInvalidLengthAuth = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAuth   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/auth.proto", fileDescriptorAuth) }
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/auth.proto", fileDescriptorAuth)
}

var fileDescriptorAuth = []byte{
	// 393 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xb1, 0x6b, 0xdb, 0x40,
	0x14, 0xc6, 0xef, 0xc5, 0x41, 0x90, 0x83, 0x06, 0x71, 0x1d, 0x2a, 0x4c, 0x79, 0x84, 0x42, 0x21,
	0x2d, 0x54, 0x2a, 0xcd, 0x5f, 0xa0, 0x36, 0x6a, 0x62, 0x02, 0x52, 0x38, 0x2b, 0x14, 0x42, 0xc1,
	0x48, 0x45, 0x95, 0x84, 0x89, 0x24, 0xe4, 0x73, 0x8b, 0x37, 0x8f, 0x1e, 0x3b, 0x76, 0x6b, 0xa1,
	0x8b, 0x47, 0x8f, 0x1e, 0x3d, 0x7a, 0xf4, 0xe8, 0xd1, 0xba, 0x5b, 0x3c, 0x7a, 0xf4, 0x58, 0x2c,
	0x75, 0xe8, 0xd6, 0x4e, 0xf7, 0x7d, 0xbf, 0xef, 0x3e, 0xee, 0x78, 0x8f, 0x9a, 0x71, 0x2a, 0x92,
	0x61, 0x68, 0x7e, 0xca, 0x1f, 0x2c, 0x3f, 0x89, 0xfc, 0x24, 0xcd, 0xe2, 0x81, 0x1b, 0x89, 0xaf,
	0x79, 0xd9, 0xb7, 0x84, 0xc8, 0xac, 0xa0, 0x48, 0xad, 0x60, 0x28, 0x12, 0xb3, 0x28, 0x73, 0x91,
	0x33, 0x4d, 0x88, 0xcc, 0xfc, 0x72, 0xd1, 0x7e, 0xfd, 0x3f, 0xbd, 0x32, 0x8d, 0x13, 0x31, 0x68,
	0x9a, 0xed, 0x57, 0x7f, 0x35, 0xe2, 0x3c, 0xce, 0xad, 0x1a, 0x87, 0xc3, 0xcf, 0xb5, 0xab, 0x4d,
	0xad, 0x9a, 0xeb, 0xcf, 0xee, 0xa8, 0x66, 0xdf, 0x76, 0x6e, 0xa2, 0x11, 0xd3, 0x69, 0xab, 0x1f,
	0x8d, 0x0c, 0x38, 0x83, 0xf3, 0x13, 0x7e, 0x90, 0x8c, 0xd1, 0xe3, 0x2c, 0x78, 0x88, 0x8c, 0xa3,
	0x1a, 0xd5, 0x9a, 0x3d, 0xa7, 0x5a, 0xf3, 0x9c, 0xd1, 0x3a, 0x6b, 0x9d, 0x9f, 0xbe, 0x79, 0x64,
	0x36, 0x3f, 0x35, 0xf9, 0x81, 0xf2, 0x3f, 0xe1, 0xcb, 0x8f, 0xf4, 0xe4, 0xaa, 0x0c, 0x32, 0xe1,
	0x8f, 0x8a, 0x88, 0x3d, 0xa5, 0xc6, 0x15, 0xb7, 0x5d, 0xbf, 0x67, 0xdf, 0xf9, 0xd7, 0x1e, 0xef,
	0xdc, 0xdb, 0x7e, 0xc7, 0x73, 0x7b, 0xef, 0xbc, 0x4b, 0x47, 0x27, 0x8c, 0xd1, 0xd3, 0x26, 0xbd,
	0xb5, 0xbb, 0xdd, 0x0f, 0x1e, 0xbf, 0xd4, 0x81, 0x3d, 0xa1, 0x8f, 0x1b, 0xc6, 0x9d, 0xf7, 0xdc,
	0xe9, 0x5e, 0xf7, 0x7c, 0xef, 0xc6, 0x71, 0xf5, 0xa3, 0xf6, 0xf1, 0xe4, 0x17, 0x92, 0xb7, 0x3f,
	0x60, 0x59, 0x21, 0xac, 0x2a, 0x84, 0x75, 0x85, 0xb0, 0xa9, 0x10, 0xb6, 0x15, 0x92, 0x5d, 0x85,
	0x64, 0x5f, 0x21, 0x8c, 0x25, 0x92, 0x89, 0x44, 0x32, 0x95, 0x08, 0x33, 0x89, 0x64, 0x2e, 0x11,
	0x16, 0x12, 0x61, 0x29, 0x11, 0x56, 0x12, 0x61, 0x2d, 0x91, 0x6c, 0x24, 0xc2, 0x56, 0x22, 0xd9,
	0x49, 0x84, 0xbd, 0x44, 0x32, 0x56, 0x48, 0x26, 0x0a, 0xe1, 0x9b, 0x42, 0xf2, 0x5d, 0x21, 0xfc,
	0x54, 0x48, 0xa6, 0x0a, 0xc9, 0x4c, 0x21, 0xcc, 0x15, 0xc2, 0x42, 0x21, 0xdc, 0xbf, 0xf8, 0xd7,
	0x32, 0x8a, 0x7e, 0x7c, 0x38, 0x8b, 0x30, 0xd4, 0xea, 0xe9, 0x5e, 0xfc, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0xc1, 0xc9, 0x6f, 0xdc, 0xf8, 0x01, 0x00, 0x00,
}
