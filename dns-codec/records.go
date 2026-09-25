package dns

import (
	"encoding/binary"
	"net"
)

const (
	typeLenA    int = 4
	typeLenAAAA int = 16
)

type MXRecord struct {
	Preference uint16
	NameLen    int
}

type Record struct {
	raw      []byte
	nameEnd  int
	offset   int
	rdataLen int
}

type RecordIter struct {
	raw          []byte
	offset       int
	count        uint16
	firstNameEnd int
}

func (it *RecordIter) Next() (Record, bool) {
	if it.count == 0 {
		return Record{}, false
	}

	start := it.offset
	nameEnd := it.firstNameEnd
	if nameEnd == 0 {
		var err error
		nameEnd, err = skipName(it.raw, start)
		if err != nil {
			return Record{}, false
		}
	}

	it.firstNameEnd = 0
	if nameEnd+recordSize > len(it.raw) {
		return Record{}, false
	}

	rdataLen := int(binary.BigEndian.Uint16(
		it.raw[nameEnd+sizeTYPE+sizeCLASS+sizeTTL : nameEnd+recordSize],
	))

	rdataOff := nameEnd + recordSize
	if rdataOff+rdataLen > len(it.raw) {
		return Record{}, false
	}

	it.offset = rdataOff + rdataLen
	it.count--

	return Record{
		raw:      it.raw,
		offset:   start,
		nameEnd:  nameEnd,
		rdataLen: rdataLen,
	}, true
}

func (r Record) Name(dst []byte) (int, error) {
	return decodeName(r.raw, r.offset, dst)
}

func (r Record) Type() Type {
	return Type(binary.BigEndian.Uint16(r.raw[r.nameEnd : r.nameEnd+uint16Len]))
}

func (r Record) Class() ClassType {
	return ClassType(binary.BigEndian.Uint16(r.raw[r.nameEnd+uint16Len : r.nameEnd+uint32Len]))
}

func (r Record) TTL() uint32 {
	return binary.BigEndian.Uint32(r.raw[r.nameEnd+uint32Len : r.nameEnd+8])
}

func (r Record) NS(dst []byte) (int, error) {
	if r.Type() != TypeNS {
		return 0, ErrBadRR
	}

	return decodeName(r.raw, r.nameEnd+recordSize, dst)
}

func (r Record) CNAME(dst []byte) (int, error) {
	if r.Type() != TypeCNAME {
		return 0, ErrBadRR
	}

	return decodeName(r.raw, r.nameEnd+recordSize, dst)
}

func (r Record) PTR(dst []byte) (int, error) {
	if r.Type() != TypePTR {
		return 0, ErrBadRR
	}

	return decodeName(r.raw, r.nameEnd+recordSize, dst)
}

func (r Record) AAAA() (net.IP, error) {
	if r.Type() != TypeAAAA {
		return net.IP{}, ErrBadRR
	}

	if r.rdataLen != typeLenAAAA {
		return net.IP{}, ErrShortPacket
	}

	rdataOff := r.nameEnd + recordSize
	if rdataOff+typeLenAAAA > len(r.raw) {
		return net.IP{}, ErrShortPacket
	}

	return r.raw[rdataOff : rdataOff+typeLenAAAA], nil
}

func (r Record) A() (net.IP, error) {
	if r.Type() != TypeA {
		return net.IP{}, ErrBadRR
	}

	if r.rdataLen != typeLenA {
		return net.IP{}, ErrShortPacket
	}

	rdataOff := r.nameEnd + recordSize
	if rdataOff+typeLenA > len(r.raw) {
		return net.IP{}, ErrShortPacket
	}

	return r.raw[rdataOff : rdataOff+typeLenA], nil
}

func (r Record) MX(dst []byte) (MXRecord, error) {
	if r.Type() != TypeMX {
		return MXRecord{}, ErrBadRR
	}

	if r.rdataLen < 3 {
		return MXRecord{}, ErrShortPacket
	}

	rdataOff := r.nameEnd + recordSize
	if rdataOff+uint16Len > len(r.raw) {
		return MXRecord{}, ErrShortPacket
	}

	pref := binary.BigEndian.Uint16(r.raw[rdataOff : rdataOff+uint16Len])
	nameLen, err := decodeName(r.raw, rdataOff+uint16Len, dst)
	if err != nil {
		return MXRecord{}, err
	}

	return MXRecord{
		Preference: pref,
		NameLen:    nameLen,
	}, nil
}
