package dns

import (
	"testing"
)

import (
	"net"
)

func TestDecode(t *testing.T) {
	msg := Message{
		Header: Header{
			ID:                 1234,
			Response:           true,
			Opcode:             OpcodeQUERY,
			Authoritative:      true,
			RecursionDesired:   true,
			RecursionAvailable: true,
			RCode:              RCodeNOERROR,
		},
		Questions: []Questions{
			{Name: "example.com.", Type: TypeA, Class: ClassIN},
		},
		Answers: []RR{
			{
				Name:  "example.com.",
				Type:  TypeA,
				Class: ClassIN,
				TTL:   300,
				RData: &AResource{A: net.IPv4(192, 0, 2, 1)},
			},
		},
	}

	packed, err := msg.Pack()
	if err != nil {
		t.Fatalf("Pack() unexpected error: %v", err)
	}

	want, err := Decode(packed)
	if err != nil {
		t.Fatalf("Decode() unexpected error: %v", err)
	}

	if want.header.id != msg.Header.ID {
		t.Errorf("header.id = %d, want %d", want.header.id, msg.Header.ID)
	}
	if want.header.QR() != msg.Header.Response {
		t.Errorf("header.response = %t, want %t", want.header.QR(), msg.Header.Response)
	}
	if want.header.Opcode() != msg.Header.Opcode {
		t.Errorf("header.Opcode = %d, want %d", want.header.Opcode(), msg.Header.Opcode)
	}
	if want.header.AA() != msg.Header.Authoritative {
		t.Errorf("header.Authoritative = %t, want %t", want.header.AA(), msg.Header.Authoritative)
	}
	if want.header.RD() != msg.Header.RecursionDesired {
		t.Errorf("header.RecursionDesired = %t, want %t", want.header.RD(), msg.Header.RecursionDesired)
	}
	if want.header.RA() != msg.Header.RecursionAvailable {
		t.Errorf("header.RecursionAvailable = %t, want %t", want.header.RA(), msg.Header.RecursionAvailable)
	}
	if want.header.RCode() != msg.Header.RCode {
		t.Errorf("header.RCode = %d, want %d", want.header.RCode(), msg.Header.RCode)
	}
	if got, want := want.header.qdcount(), uint16(len(msg.Questions)); got != want {
		t.Errorf("qdcount = %d, want %d", got, want)
	}
	if got, want := want.header.ancount(), uint16(len(msg.Answers)); got != want {
		t.Errorf("ancount = %d, want %d", got, want)
	}
	if got, want := want.header.nscount(), uint16(len(msg.Authorities)); got != want {
		t.Errorf("nscount = %d, want %d", got, want)
	}
	if got, want := want.header.arcount(), uint16(len(msg.Additionals)); got != want {
		t.Errorf("arcount = %d, want %d", got, want)
	}

	var buf [255]byte
	var n int
	qIter := want.Questions()
	for i := 0; i < int(want.header.qdcount()); i++ {
		item, ok := qIter.Next()
		if !ok {
			t.Errorf("question iter failed")
		}

		n, err = item.Name(buf[:])
		if err != nil {
			t.Errorf("failed to get name question")
		}
		qName := string(buf[:n])

		if qName != msg.Questions[0].Name {
			t.Errorf("name = %s, want %s", qName, msg.Questions[0].Name)
		}
		if item.Type() != msg.Questions[0].Type {
			t.Errorf("type = %d, want %d", item.Type(), msg.Questions[0].Type)
		}
		if item.Class() != msg.Questions[0].Class {
			t.Errorf("class = %d, want %d", item.Class(), msg.Questions[0].Class)
		}
	}

	aIter := want.Answers()
	for i := 0; i < int(want.header.ancount()); i++ {
		item, ok := aIter.Next()
		if !ok {
			t.Errorf("answer iter failed")
		}

		n, err = item.Name(buf[:])
		if err != nil {
			t.Errorf("failed to get name answer")
		}

		name := string(buf[:n])
		rdata, err := item.A()
		if err != nil {
			t.Errorf("failed get rdata")
		}

		if name != msg.Answers[0].Name {
			t.Errorf("name = %s, want %s", name, msg.Answers[0].Name)
		}
		if item.Type() != msg.Answers[0].Type {
			t.Errorf("type = %d, want %d", item.Type(), msg.Answers[0].Type)
		}
		if item.Class() != msg.Answers[0].Class {
			t.Errorf("class = %d, want %d", item.Class(), msg.Answers[0].Class)
		}
		rdataMsg, ok := msg.Answers[i].RData.(*AResource)
		if !ok {
			t.Errorf("failed get AResource")
		}
		if !net.IP.Equal(rdata, rdataMsg.A) {
			t.Errorf("rdata = %s, want %s", rdata.String(), rdataMsg.A.String())
		}
	}

}
