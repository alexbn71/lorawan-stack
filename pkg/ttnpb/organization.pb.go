// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/TheThingsNetwork/ttn/api/organization.proto

package ttnpb

import proto "github.com/gogo/protobuf/proto"
import golang_proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import _ "github.com/gogo/protobuf/types"

import time "time"

import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

type Organization struct {
	// id is the unique and immutable Organization's ID.
	// This ID shares namespace with the users' IDs.
	OrganizationIdentifier `protobuf:"bytes,1,opt,name=id,embedded=id" json:"id"`
	// name is the organization's name.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// description is an organization's description.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// url is the URL of the organization website.
	URL string `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	// location is the organization's location, e.g. "Amsterdam, Europe".
	Location string `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	// email is a generic contact email address that is shown as contact email.
	Email string `protobuf:"bytes,6,opt,name=email,proto3" json:"email,omitempty"`
	// created_at denotes when the user was created.
	// This is a read-only field.
	CreatedAt time.Time `protobuf:"bytes,7,opt,name=created_at,json=createdAt,stdtime" json:"created_at"`
	// updated_at is the last time the user was updated.
	// This is a read-only field.
	UpdatedAt time.Time `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,stdtime" json:"updated_at"`
}

func (m *Organization) Reset()                    { *m = Organization{} }
func (m *Organization) String() string            { return proto.CompactTextString(m) }
func (*Organization) ProtoMessage()               {}
func (*Organization) Descriptor() ([]byte, []int) { return fileDescriptorOrganization, []int{0} }

func (m *Organization) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Organization) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Organization) GetURL() string {
	if m != nil {
		return m.URL
	}
	return ""
}

func (m *Organization) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Organization) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Organization) GetCreatedAt() time.Time {
	if m != nil {
		return m.CreatedAt
	}
	return time.Time{}
}

func (m *Organization) GetUpdatedAt() time.Time {
	if m != nil {
		return m.UpdatedAt
	}
	return time.Time{}
}

type OrganizationMember struct {
	// organization_id is the organization's ID.
	OrganizationIdentifier `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,embedded=organization_id" json:"organization_id"`
	// user_id is the user's ID.
	UserIdentifier `protobuf:"bytes,2,opt,name=user_id,json=userId,embedded=user_id" json:"user_id"`
	// rights is the list of rights the user bears to the organization.
	Rights []Right `protobuf:"varint,3,rep,packed,name=rights,enum=ttn.v3.Right" json:"rights,omitempty"`
}

func (m *OrganizationMember) Reset()                    { *m = OrganizationMember{} }
func (m *OrganizationMember) String() string            { return proto.CompactTextString(m) }
func (*OrganizationMember) ProtoMessage()               {}
func (*OrganizationMember) Descriptor() ([]byte, []int) { return fileDescriptorOrganization, []int{1} }

func (m *OrganizationMember) GetRights() []Right {
	if m != nil {
		return m.Rights
	}
	return nil
}

func init() {
	proto.RegisterType((*Organization)(nil), "ttn.v3.Organization")
	golang_proto.RegisterType((*Organization)(nil), "ttn.v3.Organization")
	proto.RegisterType((*OrganizationMember)(nil), "ttn.v3.OrganizationMember")
	golang_proto.RegisterType((*OrganizationMember)(nil), "ttn.v3.OrganizationMember")
}
func (this *Organization) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*Organization)
	if !ok {
		that2, ok := that.(Organization)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *Organization")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *Organization but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *Organization but is not nil && this == nil")
	}
	if !this.OrganizationIdentifier.Equal(&that1.OrganizationIdentifier) {
		return fmt.Errorf("OrganizationIdentifier this(%v) Not Equal that(%v)", this.OrganizationIdentifier, that1.OrganizationIdentifier)
	}
	if this.Name != that1.Name {
		return fmt.Errorf("Name this(%v) Not Equal that(%v)", this.Name, that1.Name)
	}
	if this.Description != that1.Description {
		return fmt.Errorf("Description this(%v) Not Equal that(%v)", this.Description, that1.Description)
	}
	if this.URL != that1.URL {
		return fmt.Errorf("URL this(%v) Not Equal that(%v)", this.URL, that1.URL)
	}
	if this.Location != that1.Location {
		return fmt.Errorf("Location this(%v) Not Equal that(%v)", this.Location, that1.Location)
	}
	if this.Email != that1.Email {
		return fmt.Errorf("Email this(%v) Not Equal that(%v)", this.Email, that1.Email)
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return fmt.Errorf("CreatedAt this(%v) Not Equal that(%v)", this.CreatedAt, that1.CreatedAt)
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return fmt.Errorf("UpdatedAt this(%v) Not Equal that(%v)", this.UpdatedAt, that1.UpdatedAt)
	}
	return nil
}
func (this *Organization) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*Organization)
	if !ok {
		that2, ok := that.(Organization)
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
	if !this.OrganizationIdentifier.Equal(&that1.OrganizationIdentifier) {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if this.URL != that1.URL {
		return false
	}
	if this.Location != that1.Location {
		return false
	}
	if this.Email != that1.Email {
		return false
	}
	if !this.CreatedAt.Equal(that1.CreatedAt) {
		return false
	}
	if !this.UpdatedAt.Equal(that1.UpdatedAt) {
		return false
	}
	return true
}
func (this *OrganizationMember) VerboseEqual(that interface{}) error {
	if that == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that == nil && this != nil")
	}

	that1, ok := that.(*OrganizationMember)
	if !ok {
		that2, ok := that.(OrganizationMember)
		if ok {
			that1 = &that2
		} else {
			return fmt.Errorf("that is not of type *OrganizationMember")
		}
	}
	if that1 == nil {
		if this == nil {
			return nil
		}
		return fmt.Errorf("that is type *OrganizationMember but is nil && this != nil")
	} else if this == nil {
		return fmt.Errorf("that is type *OrganizationMember but is not nil && this == nil")
	}
	if !this.OrganizationIdentifier.Equal(&that1.OrganizationIdentifier) {
		return fmt.Errorf("OrganizationIdentifier this(%v) Not Equal that(%v)", this.OrganizationIdentifier, that1.OrganizationIdentifier)
	}
	if !this.UserIdentifier.Equal(&that1.UserIdentifier) {
		return fmt.Errorf("UserIdentifier this(%v) Not Equal that(%v)", this.UserIdentifier, that1.UserIdentifier)
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
func (this *OrganizationMember) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*OrganizationMember)
	if !ok {
		that2, ok := that.(OrganizationMember)
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
	if !this.OrganizationIdentifier.Equal(&that1.OrganizationIdentifier) {
		return false
	}
	if !this.UserIdentifier.Equal(&that1.UserIdentifier) {
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
func (m *Organization) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Organization) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(m.OrganizationIdentifier.Size()))
	n1, err := m.OrganizationIdentifier.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Description) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Description)))
		i += copy(dAtA[i:], m.Description)
	}
	if len(m.URL) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.URL)))
		i += copy(dAtA[i:], m.URL)
	}
	if len(m.Location) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Location)))
		i += copy(dAtA[i:], m.Location)
	}
	if len(m.Email) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(len(m.Email)))
		i += copy(dAtA[i:], m.Email)
	}
	dAtA[i] = 0x3a
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)))
	n2, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CreatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	dAtA[i] = 0x42
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)))
	n3, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpdatedAt, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func (m *OrganizationMember) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrganizationMember) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(m.OrganizationIdentifier.Size()))
	n4, err := m.OrganizationIdentifier.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	dAtA[i] = 0x12
	i++
	i = encodeVarintOrganization(dAtA, i, uint64(m.UserIdentifier.Size()))
	n5, err := m.UserIdentifier.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	if len(m.Rights) > 0 {
		dAtA7 := make([]byte, len(m.Rights)*10)
		var j6 int
		for _, num := range m.Rights {
			for num >= 1<<7 {
				dAtA7[j6] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j6++
			}
			dAtA7[j6] = uint8(num)
			j6++
		}
		dAtA[i] = 0x1a
		i++
		i = encodeVarintOrganization(dAtA, i, uint64(j6))
		i += copy(dAtA[i:], dAtA7[:j6])
	}
	return i, nil
}

func encodeVarintOrganization(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func NewPopulatedOrganization(r randyOrganization, easy bool) *Organization {
	this := &Organization{}
	v1 := NewPopulatedOrganizationIdentifier(r, easy)
	this.OrganizationIdentifier = *v1
	this.Name = randStringOrganization(r)
	this.Description = randStringOrganization(r)
	this.URL = randStringOrganization(r)
	this.Location = randStringOrganization(r)
	this.Email = randStringOrganization(r)
	v2 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.CreatedAt = *v2
	v3 := github_com_gogo_protobuf_types.NewPopulatedStdTime(r, easy)
	this.UpdatedAt = *v3
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

func NewPopulatedOrganizationMember(r randyOrganization, easy bool) *OrganizationMember {
	this := &OrganizationMember{}
	v4 := NewPopulatedOrganizationIdentifier(r, easy)
	this.OrganizationIdentifier = *v4
	v5 := NewPopulatedUserIdentifier(r, easy)
	this.UserIdentifier = *v5
	v6 := r.Intn(10)
	this.Rights = make([]Right, v6)
	for i := 0; i < v6; i++ {
		this.Rights[i] = Right([]int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 17, 18, 31, 32, 33, 34, 35, 36, 37, 38, 39, 51, 52, 53, 54, 55, 56, 57, 58, 71, 72, 73, 74, 75, 76, 77, 78, 79}[r.Intn(40)])
	}
	if !easy && r.Intn(10) != 0 {
	}
	return this
}

type randyOrganization interface {
	Float32() float32
	Float64() float64
	Int63() int64
	Int31() int32
	Uint32() uint32
	Intn(n int) int
}

func randUTF8RuneOrganization(r randyOrganization) rune {
	ru := r.Intn(62)
	if ru < 10 {
		return rune(ru + 48)
	} else if ru < 36 {
		return rune(ru + 55)
	}
	return rune(ru + 61)
}
func randStringOrganization(r randyOrganization) string {
	v7 := r.Intn(100)
	tmps := make([]rune, v7)
	for i := 0; i < v7; i++ {
		tmps[i] = randUTF8RuneOrganization(r)
	}
	return string(tmps)
}
func randUnrecognizedOrganization(r randyOrganization, maxFieldNumber int) (dAtA []byte) {
	l := r.Intn(5)
	for i := 0; i < l; i++ {
		wire := r.Intn(4)
		if wire == 3 {
			wire = 5
		}
		fieldNumber := maxFieldNumber + r.Intn(100)
		dAtA = randFieldOrganization(dAtA, r, fieldNumber, wire)
	}
	return dAtA
}
func randFieldOrganization(dAtA []byte, r randyOrganization, fieldNumber int, wire int) []byte {
	key := uint32(fieldNumber)<<3 | uint32(wire)
	switch wire {
	case 0:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		v8 := r.Int63()
		if r.Intn(2) == 0 {
			v8 *= -1
		}
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(v8))
	case 1:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	case 2:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		ll := r.Intn(100)
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(ll))
		for j := 0; j < ll; j++ {
			dAtA = append(dAtA, byte(r.Intn(256)))
		}
	default:
		dAtA = encodeVarintPopulateOrganization(dAtA, uint64(key))
		dAtA = append(dAtA, byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)), byte(r.Intn(256)))
	}
	return dAtA
}
func encodeVarintPopulateOrganization(dAtA []byte, v uint64) []byte {
	for v >= 1<<7 {
		dAtA = append(dAtA, uint8(v&0x7f|0x80))
		v >>= 7
	}
	dAtA = append(dAtA, uint8(v))
	return dAtA
}
func (m *Organization) Size() (n int) {
	var l int
	_ = l
	l = m.OrganizationIdentifier.Size()
	n += 1 + l + sovOrganization(uint64(l))
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.URL)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.Location)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = len(m.Email)
	if l > 0 {
		n += 1 + l + sovOrganization(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CreatedAt)
	n += 1 + l + sovOrganization(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpdatedAt)
	n += 1 + l + sovOrganization(uint64(l))
	return n
}

func (m *OrganizationMember) Size() (n int) {
	var l int
	_ = l
	l = m.OrganizationIdentifier.Size()
	n += 1 + l + sovOrganization(uint64(l))
	l = m.UserIdentifier.Size()
	n += 1 + l + sovOrganization(uint64(l))
	if len(m.Rights) > 0 {
		l = 0
		for _, e := range m.Rights {
			l += sovOrganization(uint64(e))
		}
		n += 1 + sovOrganization(uint64(l)) + l
	}
	return n
}

func sovOrganization(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozOrganization(x uint64) (n int) {
	return sovOrganization((x << 1) ^ uint64((int64(x) >> 63)))
}
func (m *Organization) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrganization
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
			return fmt.Errorf("proto: Organization: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Organization: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrganizationIdentifier", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OrganizationIdentifier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field URL", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.URL = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Location = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Email", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Email = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.CreatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpdatedAt", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UpdatedAt, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOrganization(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOrganization
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
func (m *OrganizationMember) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrganization
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
			return fmt.Errorf("proto: OrganizationMember: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrganizationMember: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrganizationIdentifier", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OrganizationIdentifier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserIdentifier", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrganization
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
				return ErrInvalidLengthOrganization
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.UserIdentifier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType == 0 {
				var v Right
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowOrganization
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
						return ErrIntOverflowOrganization
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
					return ErrInvalidLengthOrganization
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v Right
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowOrganization
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
			skippy, err := skipOrganization(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOrganization
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
func skipOrganization(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOrganization
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
					return 0, ErrIntOverflowOrganization
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
					return 0, ErrIntOverflowOrganization
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
				return 0, ErrInvalidLengthOrganization
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowOrganization
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
				next, err := skipOrganization(dAtA[start:])
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
	ErrInvalidLengthOrganization = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOrganization   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/organization.proto", fileDescriptorOrganization)
}
func init() {
	golang_proto.RegisterFile("github.com/TheThingsNetwork/ttn/api/organization.proto", fileDescriptorOrganization)
}

var fileDescriptorOrganization = []byte{
	// 553 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x3f, 0x6c, 0xd4, 0x4c,
	0x10, 0xc5, 0x77, 0xee, 0x92, 0x4b, 0xb2, 0xf9, 0xbe, 0x20, 0xad, 0x10, 0x32, 0x87, 0x34, 0x77,
	0x8a, 0x84, 0x14, 0x0a, 0x7c, 0x28, 0x11, 0x08, 0xca, 0x84, 0x0a, 0x89, 0x3f, 0xc2, 0x4a, 0x1a,
	0x9a, 0xc8, 0x77, 0xb7, 0xf1, 0xad, 0x72, 0xf6, 0x5a, 0xf6, 0x1a, 0x04, 0x55, 0xca, 0x94, 0x29,
	0xe9, 0xa0, 0x4c, 0x99, 0x32, 0x65, 0xca, 0xd0, 0xa5, 0x41, 0x4a, 0x15, 0xe2, 0x75, 0x93, 0x32,
	0x65, 0x4a, 0xe4, 0xb5, 0x1d, 0x0c, 0x42, 0xe2, 0xa0, 0xf2, 0xce, 0xbc, 0xf9, 0x3d, 0x6b, 0x9e,
	0xd7, 0xf4, 0x91, 0x27, 0xd4, 0x28, 0xe9, 0xdb, 0x03, 0xe9, 0xf7, 0xd6, 0x47, 0x7c, 0x7d, 0x24,
	0x02, 0x2f, 0x7e, 0xc9, 0xd5, 0x3b, 0x19, 0x6d, 0xf7, 0x94, 0x0a, 0x7a, 0x6e, 0x28, 0x7a, 0x32,
	0xf2, 0xdc, 0x40, 0x7c, 0x70, 0x95, 0x90, 0x81, 0x1d, 0x46, 0x52, 0x49, 0xd6, 0x52, 0x2a, 0xb0,
	0xdf, 0xae, 0xb4, 0xed, 0x49, 0x78, 0x37, 0x51, 0xa3, 0x82, 0x6b, 0x3f, 0x9c, 0x64, 0x5e, 0x0c,
	0x79, 0xa0, 0xc4, 0x96, 0xe0, 0x51, 0x5c, 0x62, 0x0f, 0x26, 0xc1, 0x22, 0xe1, 0x8d, 0x54, 0x45,
	0xdc, 0xaf, 0x11, 0x9e, 0xf4, 0x64, 0xcf, 0xb4, 0xfb, 0xc9, 0x96, 0xa9, 0x4c, 0x61, 0x4e, 0xe5,
	0xf8, 0x1d, 0x4f, 0x4a, 0x6f, 0xcc, 0x7f, 0x4c, 0x71, 0x3f, 0x54, 0xef, 0x4b, 0xb1, 0xf3, 0xab,
	0xa8, 0x84, 0xcf, 0x63, 0xe5, 0xfa, 0x61, 0x31, 0xb0, 0xf8, 0xb5, 0x41, 0xff, 0x7b, 0x55, 0x0b,
	0x89, 0x3d, 0xa6, 0x0d, 0x31, 0xb4, 0xa0, 0x0b, 0x4b, 0xf3, 0xcb, 0x68, 0x17, 0x59, 0xd9, 0xf5,
	0x89, 0x67, 0xd7, 0x2b, 0xae, 0xcd, 0x1e, 0x9f, 0x75, 0xc8, 0xc9, 0x59, 0x07, 0x9c, 0x86, 0x18,
	0x32, 0x46, 0xa7, 0x02, 0xd7, 0xe7, 0x56, 0xa3, 0x0b, 0x4b, 0x73, 0x8e, 0x39, 0xb3, 0x2e, 0x9d,
	0x1f, 0xf2, 0x78, 0x10, 0x89, 0x30, 0x47, 0xad, 0xa6, 0x91, 0xea, 0x2d, 0x76, 0x9b, 0x36, 0x93,
	0x68, 0x6c, 0x4d, 0xe5, 0xca, 0xda, 0x8c, 0x3e, 0xeb, 0x34, 0x37, 0x9c, 0xe7, 0x4e, 0xde, 0x63,
	0x6d, 0x3a, 0x3b, 0x96, 0x03, 0xf3, 0x52, 0x6b, 0xda, 0x90, 0xd7, 0x35, 0xbb, 0x49, 0xa7, 0xb9,
	0xef, 0x8a, 0xb1, 0xd5, 0x32, 0x42, 0x51, 0xb0, 0xa7, 0x94, 0x0e, 0x22, 0xee, 0x2a, 0x3e, 0xdc,
	0x74, 0x95, 0x35, 0x63, 0x96, 0x68, 0xdb, 0x45, 0x06, 0x76, 0x95, 0x81, 0xbd, 0x5e, 0x65, 0x50,
	0x2c, 0xb0, 0xf7, 0xad, 0x03, 0xce, 0x5c, 0xc9, 0xad, 0xaa, 0xdc, 0x24, 0x09, 0x87, 0x95, 0xc9,
	0xec, 0xdf, 0x98, 0x94, 0xdc, 0xaa, 0x5a, 0xfc, 0x02, 0x94, 0xd5, 0x53, 0x7b, 0xc1, 0xfd, 0x3e,
	0x8f, 0xd8, 0x6b, 0x7a, 0xa3, 0x7e, 0x25, 0x37, 0xff, 0x21, 0xea, 0x05, 0xf9, 0xd3, 0x04, 0x7b,
	0x42, 0x67, 0x92, 0x98, 0x47, 0xb9, 0x55, 0xc3, 0x58, 0xdd, 0xaa, 0xac, 0x36, 0x62, 0x1e, 0xfd,
	0xd6, 0xa2, 0x95, 0x18, 0x85, 0xdd, 0xa5, 0xad, 0xe2, 0xe6, 0x59, 0xcd, 0x6e, 0x73, 0x69, 0x61,
	0xf9, 0xff, 0x8a, 0x74, 0xf2, 0xae, 0x53, 0x8a, 0x6b, 0x9f, 0xe0, 0x38, 0x45, 0x38, 0x49, 0x11,
	0x4e, 0x53, 0x84, 0xf3, 0x14, 0xe1, 0x22, 0x45, 0x72, 0x99, 0x22, 0xb9, 0x4a, 0x11, 0x76, 0x34,
	0x92, 0x5d, 0x8d, 0x64, 0x5f, 0x23, 0x1c, 0x68, 0x24, 0x87, 0x1a, 0xe1, 0x48, 0x23, 0x1c, 0x6b,
	0x84, 0x13, 0x8d, 0x70, 0xaa, 0x91, 0x9c, 0x6b, 0x84, 0x0b, 0x8d, 0xe4, 0x52, 0x23, 0x5c, 0x69,
	0x24, 0x3b, 0x19, 0x92, 0xdd, 0x0c, 0x61, 0x2f, 0x43, 0xf2, 0x31, 0x43, 0xf8, 0x9c, 0x21, 0xd9,
	0xcf, 0x90, 0x1c, 0x64, 0x08, 0x87, 0x19, 0xc2, 0x51, 0x86, 0xf0, 0xe6, 0xde, 0x9f, 0xfe, 0x9c,
	0x70, 0xdb, 0xcb, 0x9f, 0x61, 0xbf, 0xdf, 0x32, 0x9f, 0x65, 0xe5, 0x7b, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x07, 0xf6, 0x2b, 0x6b, 0x14, 0x04, 0x00, 0x00,
}
