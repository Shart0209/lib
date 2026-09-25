package dns

import (
	"bytes"
	"net"
	"testing"
)

func TestMessage_Pack(t *testing.T) {
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

	got, err := msg.Pack()
	if err != nil {
		t.Fatalf("Pack() unexpected error: %v", err)
	}

	// Expected wire format for the message above, per RFC 1035:
	//   header (12 bytes):
	//     ID=1234 (0x04d2)
	//     flags=0x8580  -> QR=1 AA=1 RD=1 RA=1 RCODE=0 (NOERROR), OPCODE=0 (QUERY)
	//     QDCOUNT=1 ANCOUNT=1 NSCOUNT=0 ARCOUNT=0
	//   question: "example.com." QTYPE=A(1) QCLASS=IN(1)
	//   answer: name is a pointer to the question's name (0xc00c),
	//     TYPE=A(1) CLASS=IN(1) TTL=300 RDLENGTH=4 RDATA=192.0.2.1
	want := []byte{
		0x04, 0xd2, // ID = 1234
		0x85, 0x80, // flags: QR|AA|RD|RA, OPCODE=QUERY, RCODE=NOERROR
		0x00, 0x01, // QDCOUNT = 1
		0x00, 0x01, // ANCOUNT = 1
		0x00, 0x00, // NSCOUNT = 0
		0x00, 0x00, // ARCOUNT = 0
		// Question: example.com.
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,       // root label
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN
		// Answer
		0xc0, 0x0c, // NAME = pointer to offset 12 (the question's name)
		0x00, 0x01, // TYPE = A
		0x00, 0x01, // CLASS = IN
		0x00, 0x00, 0x01, 0x2c, // TTL = 300
		0x00, 0x04, // RDLENGTH = 4
		192, 0, 2, 1, // RDATA
	}

	if !bytes.Equal(got, want) {
		t.Errorf("Pack() = % x, want % x", got, want)
	}
}

// BenchmarkMessage_Pack-22    	 2172373	       554.0 ns/op	    1360 B/op	       7 allocs/op
func BenchmarkMessage_Pack(b *testing.B) {
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

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := msg.Pack(); err != nil {
			b.Fatalf("Pack() unexpected error: %v", err)
		}
	}
}
