package dns

import "fmt"

func TypeString(t Type) string {
	if name, ok := typeToString[t]; ok {
		return name
	}

	return fmt.Sprintf("TYPE:%d", t)
}

func ClassString(t ClassType) string {
	if name, ok := classToString[t]; ok {
		return name
	}

	return fmt.Sprintf("CLASS:%d", t)
}

func RCodeString(code RcodeType) string {
	if data, ok := rcodeToString[code]; ok {
		return data
	}

	return fmt.Sprintf("RCODE:%d", code)
}

func OpcodeString(code OpcodeType) string {
	if data, ok := opcodeToStr[code]; ok {
		return data
	}

	return fmt.Sprintf("OPCODE:%d", code)
}
