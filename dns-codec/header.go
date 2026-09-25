package dns

type header struct {
	id      uint16
	bits    uint16
	qdCount uint16
	anCount uint16
	nsCount uint16
	arCount uint16
}

func (h header) QR() bool {
	return (h.bits & maskQR) != 0
}
func (h header) AA() bool {
	return (h.bits & maskAA) != 0
}
func (h header) TC() bool {
	return (h.bits & maskTC) != 0
}
func (h header) RD() bool {
	return (h.bits & maskRD) != 0
}
func (h header) RA() bool {
	return (h.bits & maskRA) != 0
}
func (h header) AD() bool {
	return (h.bits & maskAD) != 0
}
func (h header) CD() bool {
	return (h.bits & maskCD) != 0
}

func (h header) Opcode() OpcodeType {
	return OpcodeType((h.bits >> maskShift) & mask)
}
func (h header) RCode() RcodeType {
	return RcodeType(h.bits & mask)
}

func (h header) FlagRaw() uint16 {
	return h.bits
}

func (h header) qdcount() uint16 {
	return h.qdCount
}
func (h header) ancount() uint16 {
	return h.anCount
}
func (h header) nscount() uint16 {
	return h.nsCount
}
func (h header) arcount() uint16 {
	return h.arCount
}
