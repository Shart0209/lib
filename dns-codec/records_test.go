package dns

import (
	"bytes"
	"errors"
	"testing"
)

func TestRecordsIter(t *testing.T) {
	packet := testDNSRecordsPacket()

	msg, err := Decode(packet)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	tests := []struct {
		name  string
		typ   Type
		class ClassType
		ttl   uint32
	}{
		{"www.example.com.", TypeA, ClassIN, 60},
		{"mail.example.com.", TypeAAAA, ClassIN, 120},
		{"api.example.com.", TypeCNAME, ClassIN, 60},
		{"example.com.", TypeNS, ClassIN, 60},
		{"4.3.2.1.in-addr.arpa.", TypePTR, ClassIN, 60},
		{"example.com.", TypeMX, ClassIN, 60},
		{"api2.example.com.", TypeA, ClassIN, 60},
	}

	it := msg.Answers()
	var buf [255]byte

	for i, want := range tests {
		record, ok := it.Next()
		if !ok {
			t.Fatalf("record %d: Next() = false, want true", i)
		}

		n, err := record.Name(buf[:])
		if err != nil {
			t.Fatalf("record %d: Name() error = %v", i, err)
		}
		if got := string(buf[:n]); got != want.name {
			t.Errorf("record %d: Name() = %q, want %q", i, got, want.name)
		}
		if got := record.Type(); got != want.typ {
			t.Errorf("record %d: Type() = %d, want %d", i, got, want.typ)
		}
		if got := record.Class(); got != want.class {
			t.Errorf("record %d: ClassType() = %d, want %d", i, got, want.class)
		}
		if got := record.TTL(); got != want.ttl {
			t.Errorf("record %d: TTL() = %d, want %d", i, got, want.ttl)
		}
	}

	if _, ok := it.Next(); ok {
		t.Fatal("Next() = true after iterator exhausted")
	}
}

func TestRecordA(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		index int
		want  [4]byte
	}{
		{0, [4]byte{1, 2, 3, 4}},
		{6, [4]byte{10, 20, 30, 40}},
	}

	for _, tt := range tests {
		it := msg.Answers()
		var record Record
		var ok bool

		for i := 0; i <= tt.index; i++ {
			record, ok = it.Next()
			if !ok {
				t.Fatalf("record %d not found", tt.index)
			}
		}

		got, err := record.A()
		if err != nil {
			t.Fatalf("A() error = %v", err)
		}
		if [4]byte(got) != tt.want {
			t.Errorf("A().IP = %v, want %v", got, tt.want)
		}
	}
}

func TestRecordAAAA(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	it := msg.Answers()
	if _, ok := it.Next(); !ok {
		t.Fatal("A record not found")
	}
	record, ok := it.Next()
	if !ok {
		t.Fatal("AAAA record not found")
	}

	got, err := record.AAAA()
	if err != nil {
		t.Fatalf("AAAA() error = %v", err)
	}

	want := [16]byte{
		0x20, 0x01, 0x0d, 0xb8,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1,
	}

	if [16]byte(got) != want {
		t.Errorf("AAAA().IP = %v, want %v", got, want)
	}
}

func TestRecordCNAME(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 2; i++ {
		if _, ok := it.Next(); !ok {
			t.Fatalf("record %d not found", i)
		}
	}
	record, ok := it.Next()
	if !ok {
		t.Fatal("CNAME record not found")
	}

	var buf [255]byte
	n, err := record.CNAME(buf[:])
	if err != nil {
		t.Fatalf("CNAME() error = %v", err)
	}
	if got := string(buf[:n]); got != "www.example.com." {
		t.Errorf("CNAME() = %q, want %q", got, "www.example.com.")
	}
}

func TestRecordNS(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 3; i++ {
		if _, ok := it.Next(); !ok {
			t.Fatalf("record %d not found", i)
		}
	}
	record, ok := it.Next()
	if !ok {
		t.Fatal("NS record not found")
	}

	var buf [255]byte
	n, err := record.NS(buf[:])
	if err != nil {
		t.Fatalf("NS() error = %v", err)
	}
	if got := string(buf[:n]); got != "ns1.example.com." {
		t.Errorf("NS() = %q, want %q", got, "ns1.example.com")
	}
}

func TestRecordPTR(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 4; i++ {
		if _, ok := it.Next(); !ok {
			t.Fatalf("record %d not found", i)
		}
	}
	record, ok := it.Next()
	if !ok {
		t.Fatal("PTR record not found")
	}

	var buf [255]byte
	n, err := record.PTR(buf[:])
	if err != nil {
		t.Fatalf("PTR() error = %v", err)
	}
	if got := string(buf[:n]); got != "www.example.com." {
		t.Errorf("PTR() = %q, want %q", got, "www.example.com.")
	}
}

func TestRecordMX(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 5; i++ {
		if _, ok := it.Next(); !ok {
			t.Fatalf("record %d not found", i)
		}
	}
	record, ok := it.Next()
	if !ok {
		t.Fatal("MX record not found")
	}

	var buf [255]byte
	mx, err := record.MX(buf[:])
	if err != nil {
		t.Fatalf("MX() error = %v", err)
	}
	if mx.Preference != 10 {
		t.Errorf("MX().Preference = %d, want 10", mx.Preference)
	}
	if !bytes.Equal(buf[:mx.NameLen], []byte("mail.example.com.")) {
		t.Errorf("MX() name = %q, want %q", buf[:mx.NameLen], "mail.example.com.")
	}
}

func TestRecordWrongType(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	records := msg.Answers()
	record, ok := records.Next()
	if !ok {
		t.Fatal("record not found")
	}

	var buf [255]byte

	if _, err := record.AAAA(); !errors.Is(err, ErrBadRR) {
		t.Errorf("AAAA() error = %v, want ErrBadRR", err)
	}
	if _, err := record.NS(buf[:]); !errors.Is(err, ErrBadRR) {
		t.Errorf("NS() error = %v, want ErrBadRR", err)
	}
	if _, err := record.CNAME(buf[:]); !errors.Is(err, ErrBadRR) {
		t.Errorf("CNAME() error = %v, want ErrBadRR", err)
	}
	if _, err := record.PTR(buf[:]); !errors.Is(err, ErrBadRR) {
		t.Errorf("PTR() error = %v, want ErrBadRR", err)
	}
	if _, err := record.MX(buf[:]); !errors.Is(err, ErrBadRR) {
		t.Errorf("MX() error = %v, want ErrBadRR", err)
	}
}

func TestRecordNameShortBuffer(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}
	records := msg.Answers()
	record, ok := records.Next()
	if !ok {
		t.Fatal("record not found")
	}

	var buf [3]byte
	if _, err := record.Name(buf[:]); err == nil {
		t.Fatal("Name() error = nil, want error")
	}
}

func TestRecordIterExhausted(t *testing.T) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		t.Fatal(err)
	}

	it := msg.Answers()

	for i := 0; i < 7; i++ {
		if _, ok := it.Next(); !ok {
			t.Fatalf("record %d: Next() = false, want true", i)
		}
	}
	if _, ok := it.Next(); ok {
		t.Fatal("Next() = true after exhaustion")
	}
	if _, ok := it.Next(); ok {
		t.Fatal("Next() = true on repeated exhausted call")
	}
}

func testDNSRecordsPacket() []byte {
	return []byte{
		// header
		0x12, 0x34,
		0x81, 0x80, // response, NOERROR
		0x00, 0x01, // QDCOUNT = 1
		0x00, 0x07, // ANCOUNT = 7
		0x00, 0x00, // NSCOUNT
		0x00, 0x00, // ARCOUNT

		// Question: www.example.com IN A
		0x03, 'w', 'w', 'w',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x01,
		0x00, 0x01,

		// Answer 1: www.example.com IN A 1.2.3.4
		0x03, 'w', 'w', 'w',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x01,
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x3c,
		0x00, 0x04,
		0x01, 0x02, 0x03, 0x04,

		// Answer 2: mail.example.com IN AAAA 2001:db8::1
		0x04, 'm', 'a', 'i', 'l',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x1c,
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x78,
		0x00, 0x10,
		0x20, 0x01, 0x0d, 0xb8,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01,

		// Answer 3: api.example.com IN CNAME www.example.com
		0x03, 'a', 'p', 'i',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x05, // TYPE CNAME
		0x00, 0x01, // CLASS IN
		0x00, 0x00, 0x00, 0x3c, // TTL 60
		0x00, 0x11, // RDLENGTH = 17

		// www.example.com
		0x03, 'w', 'w', 'w',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,

		// Answer 4: example.com IN NS ns1.example.com
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x02, // TYPE NS
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x3c,
		0x00, 0x11, // RDLENGTH = 17

		// ns1.example.com
		0x03, 'n', 's', '1',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,

		// Answer 5: 4.3.2.1.in-addr.arpa IN PTR www.example.com
		0x01, '4',
		0x01, '3',
		0x01, '2',
		0x01, '1',
		0x07, 'i', 'n', '-', 'a', 'd', 'd', 'r',
		0x04, 'a', 'r', 'p', 'a',
		0x00,
		0x00, 0x0c, // TYPE PTR
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x3c,
		0x00, 0x11,

		// www.example.com
		0x03, 'w', 'w', 'w',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,

		// Answer 6: example.com IN MX 10 mail.example.com
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x0f, // TYPE MX
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x3c,
		0x00, 0x14, // RDLENGTH = 20

		// Preference = 10
		0x00, 0x0a,

		// mail.example.com
		0x04, 'm', 'a', 'i', 'l',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,

		// Answer 7: api2.example.com IN A 10.20.30.40
		0x04, 'a', 'p', 'i', '2',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x01,
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x3c,
		0x00, 0x04,
		0x0a, 0x14, 0x1e, 0x28,

		// Answer 7: api2.example.com IN A 10.20.30.40
		0x04, 'a', 'p', 'i', '2',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x01,
		0x00, 0x01,
		0x00, 0x00, 0x00, 0x3c,
		0x00, 0x04,
		0x0a, 0x14, 0x1e, 0x28,
	}
}
