package dns

import "encoding/binary"

type Question struct {
	raw     []byte
	offset  int
	nameEnd int
}

type QuestionIter struct {
	raw          []byte
	offset       int
	count        uint16
	firstNameEnd int
}

func (it *QuestionIter) Next() (Question, bool) {
	if it.count == 0 {
		return Question{}, false
	}

	nameEnd := it.firstNameEnd
	it.firstNameEnd = 0
	if nameEnd == 0 {
		var err error
		nameEnd, err = skipName(it.raw, it.offset)
		if err != nil {
			return Question{}, false
		}
	}

	if nameEnd+questionsSize > len(it.raw) {
		return Question{}, false
	}

	q := Question{
		raw:     it.raw,
		offset:  it.offset,
		nameEnd: nameEnd,
	}

	it.offset = nameEnd + questionsSize
	it.count--

	return q, true
}

func (q Question) Name(dst []byte) (int, error) {
	return decodeName(q.raw, q.offset, dst)
}

func (q Question) Type() Type {
	return Type(binary.BigEndian.Uint16(q.raw[q.nameEnd : q.nameEnd+uint16Len]))
}

func (q Question) Class() ClassType {
	return ClassType(binary.BigEndian.Uint16(q.raw[q.nameEnd+uint16Len : q.nameEnd+uint32Len]))
}
