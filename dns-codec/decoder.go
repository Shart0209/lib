package dns

import (
	"encoding/binary"
)

type Parser struct {
	raw        []byte
	header     header
	questions  section
	answers    section
	authority  section
	additional section

	questionsNameEnd int
	answersNameEnd   int
}

type section struct {
	offset int
	count  uint16
}

func Decode(data []byte) (*Parser, error) {
	if len(data) < headerSize {
		return nil, ErrShortPacket
	}

	p := &Parser{raw: data}
	p.decodeHeader()
	p.questions = section{offset: headerSize, count: p.header.qdcount()}

	offset, err := p.skipQuestions(headerSize, p.questions.count)
	if err != nil {
		return nil, err
	}

	p.answers = section{
		offset: offset,
		count:  p.header.ancount(),
	}

	offset, p.answersNameEnd, err = p.skipRecords(offset, p.answers.count)
	if err != nil {
		return nil, err
	}

	p.authority = section{
		offset: offset,
		count:  p.header.nscount(),
	}

	offset, _, err = p.skipRecords(offset, p.authority.count)
	if err != nil {
		return nil, err
	}

	p.additional = section{
		offset: offset,
		count:  p.header.arcount(),
	}

	offset, _, err = p.skipRecords(offset, p.authority.count)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Parser) decodeHeader() {
	p.header.id = binary.BigEndian.Uint16(p.raw[0:2])
	p.header.bits = binary.BigEndian.Uint16(p.raw[2:4])
	p.header.qdCount = binary.BigEndian.Uint16(p.raw[4:6])
	p.header.anCount = binary.BigEndian.Uint16(p.raw[6:8])
	p.header.nsCount = binary.BigEndian.Uint16(p.raw[8:10])
	p.header.arCount = binary.BigEndian.Uint16(p.raw[10:12])
}

func (p *Parser) skipQuestions(offset int, count uint16) (int, error) {
	for i := uint16(0); i < count; i++ {
		var err error

		offset, err = skipName(p.raw, offset)
		if err != nil {
			return 0, err
		}

		if i == 0 {
			p.answersNameEnd = offset
		}

		if offset+questionsSize > len(p.raw) {
			return 0, ErrShortPacket
		}
		offset += questionsSize
	}

	return offset, nil
}

func (p *Parser) skipRecords(offset int, count uint16) (int, int, error) {
	var firstNameEnd int

	for i := uint16(0); i < count; i++ {
		nameEnd, err := skipName(p.raw, offset)
		if err != nil {
			return 0, 0, err
		}

		if i == 0 {
			firstNameEnd = nameEnd
		}

		if nameEnd+recordSize > len(p.raw) {
			return 0, 0, ErrShortPacket
		}

		rdLength := int(binary.BigEndian.Uint16(
			p.raw[nameEnd+sizeTYPE+sizeCLASS+sizeTTL : nameEnd+recordSize],
		))

		offset = nameEnd + recordSize
		if offset+rdLength > len(p.raw) {
			return 0, 0, ErrShortPacket
		}
		offset += rdLength
	}

	return offset, firstNameEnd, nil
}

func skipName(data []byte, offset int) (int, error) {
	if offset >= len(data) {
		return 0, ErrShortPacket
	}

	for {
		if offset >= len(data) {
			return 0, ErrShortPacket
		}

		b := data[offset]

		// Compression pointer.
		if b&0xC0 == 0xC0 {
			if offset+uint16Len > len(data) {
				return 0, ErrShortPacket
			}

			return offset + uint16Len, nil
		}

		// Reserved label types.
		if b&0xC0 != 0 {
			return 0, ErrBadName
		}

		// End of name.
		if b == 0 {
			return offset + 1, nil
		}

		// Label length.
		length := int(b)

		// DNS labels max 63 bytes.
		if length > maxLabelSize {
			return 0, ErrBadName
		}

		offset++
		if offset+length > len(data) {
			return 0, ErrShortPacket
		}

		offset += length
	}
}

func decodeName(raw []byte, offset int, dst []byte) (int, error) {
	pos := offset
	n := 0
	jumps := 0

	for {
		if pos >= len(raw) {
			return 0, ErrShortPacket
		}

		b := raw[pos]
		switch b & 0xc0 {
		case 0:
			pos++

			if b == 0 {
				if n >= len(dst) {
					return 0, ErrInvalidName
				}
				dst[n] = '.'
				n++
				return n, nil
			}

			labelLen := int(b)
			if labelLen > maxLabelSize || pos+labelLen > len(raw) {
				return 0, ErrInvalidName
			}

			if n > 0 {
				if n >= len(dst) {
					return 0, ErrInvalidName
				}
				dst[n] = '.'
				n++
			}

			if n+labelLen > len(dst) {
				return 0, ErrInvalidName
			}

			copy(dst[n:], raw[pos:pos+labelLen])
			n += labelLen
			pos += labelLen

		case 0xc0:
			if pos+1 >= len(raw) {
				return 0, ErrShortPacket
			}

			ptr := int(b&0x3f)<<8 | int(raw[pos+1])
			if ptr >= len(raw) {
				return 0, ErrInvalidName
			}

			pos = ptr
			jumps++
			if jumps > len(raw) {
				return 0, ErrInvalidName
			}

		default:
			return 0, ErrInvalidName
		}
	}
}

func (p *Parser) Questions() QuestionIter {
	return QuestionIter{
		raw:          p.raw,
		offset:       p.questions.offset,
		count:        p.questions.count,
		firstNameEnd: p.questionsNameEnd,
	}
}

func (p *Parser) Answers() RecordIter {
	return RecordIter{
		raw:          p.raw,
		offset:       p.answers.offset,
		count:        p.answers.count,
		firstNameEnd: p.answersNameEnd,
	}
}

func (p *Parser) Authority() RecordIter {
	return RecordIter{
		raw:    p.raw,
		offset: p.authority.offset,
		count:  p.authority.count,
	}
}

func (p *Parser) Additional() RecordIter {
	return RecordIter{
		raw:    p.raw,
		offset: p.additional.offset,
		count:  p.additional.count,
	}
}
