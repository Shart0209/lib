package dns

type (
	Type       uint16
	ClassType  uint16
	OpcodeType uint16
	RcodeType  uint16
)

const (
	RCodeNOERROR  RcodeType = 0
	RCodeFORMERR  RcodeType = 1
	RCodeSERVFAIL RcodeType = 2
	RcodeNXDOMAIN RcodeType = 3
	RCodeNOTIMP   RcodeType = 4
	RCodeREFUSED  RcodeType = 5
)

const (
	OpcodeQUERY  OpcodeType = 0
	OpcodeIQUERY OpcodeType = 1
	OpcodeSTATUS OpcodeType = 2
	OpcodeNOTIFY OpcodeType = 4
	OpcodeUPDATE OpcodeType = 5
)

const (
	TypeA     Type = 1
	TypeNS    Type = 2
	TypeCNAME Type = 5
	TypePTR   Type = 12
	TypeMX    Type = 15
	TypeAAAA  Type = 28
)

const (
	ClassIN  ClassType = 1
	ClassCS  ClassType = 2
	ClassCH  ClassType = 3
	ClassHS  ClassType = 4
	ClassANY ClassType = 255
)

const (
	uint16Len = 2
	uint32Len = 4

	headerSize   = 6 * uint16Len
	maxLabelSize = 63

	sizeTYPE     = uint16Len
	sizeCLASS    = uint16Len
	sizeTTL      = uint32Len
	sizeRDLENGTH = uint16Len

	recordSize    = sizeTYPE + sizeCLASS + sizeTTL + sizeRDLENGTH
	questionsSize = sizeTYPE + sizeCLASS
)

const (
	maskQR    uint16 = 1 << 15
	maskAA    uint16 = 1 << 10
	maskTC    uint16 = 1 << 9
	maskRD    uint16 = 1 << 8
	maskRA    uint16 = 1 << 7
	maskAD    uint16 = 1 << 5
	maskCD    uint16 = 1 << 4
	mask      uint16 = 0x0F
	maskShift        = 11
)

var rcodeToString = map[RcodeType]string{
	RCodeNOERROR:  "NOERROR",
	RCodeFORMERR:  "FORMERR",
	RCodeSERVFAIL: "SERVFAIL",
	RcodeNXDOMAIN: "NXDOMAIN",
	RCodeNOTIMP:   "NOTIMP",
	RCodeREFUSED:  "REFUSED",
}

var typeToString = map[Type]string{
	TypeA:     "A",
	TypeNS:    "NS",
	TypeCNAME: "CNAME",
	TypePTR:   "PTR",
	TypeMX:    "MX",
	TypeAAAA:  "AAAA",
}

var classToString = map[ClassType]string{
	ClassIN:  "IN",
	ClassCS:  "CS",
	ClassCH:  "CH",
	ClassHS:  "HS",
	ClassANY: "ANY",
}

var opcodeToStr = map[OpcodeType]string{
	OpcodeQUERY:  "QUERY",
	OpcodeIQUERY: "IQUERY",
	OpcodeSTATUS: "STATUS",
	OpcodeNOTIFY: "NOTIFY",
	OpcodeUPDATE: "UPDATE",
}
