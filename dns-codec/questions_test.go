package dns

import (
	"bytes"
	"testing"
)

func TestQuestionsIter(t *testing.T) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	tests := []struct {
		name  string
		typ   Type
		class ClassType
	}{
		{
			name:  "www.example.com.",
			typ:   Type(1),
			class: ClassType(1),
		},
		{
			name:  "mail.example.com.",
			typ:   Type(28),
			class: ClassType(1),
		},
	}

	it := msg.Questions()

	var buf [255]byte

	for i, want := range tests {
		t.Run(want.name, func(t *testing.T) {
			got, ok := it.Next()
			if !ok {
				t.Fatal("Next() = false, want true")
			}

			n, err := got.Name(buf[:])
			if err != nil {
				t.Fatalf("Name() error = %v", err)
			}

			if !bytes.Equal(buf[:n], []byte(want.name)) {
				t.Errorf(
					"Name() = %q, want %q",
					buf[:n],
					want.name,
				)
			}

			if got.Type() != want.typ {
				t.Errorf(
					"Type() = %d, want %d",
					got.Type(),
					want.typ,
				)
			}

			if got.Class() != want.class {
				t.Errorf(
					"ClassType() = %d, want %d",
					got.Class(),
					want.class,
				)
			}

			// Проверяем, что offset следующего вопроса
			// действительно рассчитывается относительно текущего.
			if i == len(tests)-1 {
				if _, ok := it.Next(); ok {
					t.Fatal("Next() = true after last question, want false")
				}
			}
		})
	}
}

func TestQuestionName(t *testing.T) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	it := msg.Questions()

	q, ok := it.Next()
	if !ok {
		t.Fatal("Next() = false, want true")
	}

	var buf [255]byte

	n, err := q.Name(buf[:])
	if err != nil {
		t.Fatalf("Name() error = %v", err)
	}

	want := []byte("www.example.com.")

	if !bytes.Equal(buf[:n], want) {
		t.Errorf("Name() = %q, want %q", buf[:n], want)
	}
}

func TestQuestionNameShortBuffer(t *testing.T) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	it := msg.Questions()

	q, ok := it.Next()
	if !ok {
		t.Fatal("Next() = false, want true")
	}

	buf := make([]byte, 3)

	if _, err := q.Name(buf); err == nil {
		t.Fatal("Name() error = nil, want error")
	}
}

func TestQuestionIterExhausted(t *testing.T) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	it := msg.Questions()

	for {
		if _, ok := it.Next(); !ok {
			break
		}
	}

	if _, ok := it.Next(); ok {
		t.Fatal("Next() = true after iterator exhausted, want false")
	}

	if _, ok := it.Next(); ok {
		t.Fatal("Next() = true on repeated call after exhaustion, want false")
	}
}

func testDNSPacketQuestion() []byte {
	return []byte{
		// header
		0x12, 0x34, // ID
		0x01, 0x00, // Flags
		0x00, 0x02, // QDCOUNT = 2
		0x00, 0x00, // ANCOUNT = 0
		0x00, 0x00, // NSCOUNT = 0
		0x00, 0x00, // ARCOUNT = 0

		// Question 1: www.example.com IN A
		0x03, 'w', 'w', 'w',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x01, // QTYPE = A
		0x00, 0x01, // QCLASS = IN

		// Question 2: mail.example.com IN AAAA
		0x04, 'm', 'a', 'i', 'l',
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,
		0x00, 0x1c, // QTYPE = AAAA
		0x00, 0x01, // QCLASS = IN
	}
}

var (
	benchQuestion      Question
	benchQuestionNameN int
	benchQuestionType  Type
	benchQuestionClass ClassType
)

func BenchmarkQuestionNextOne(b *testing.B) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := msg.Questions()

		q, ok := it.Next()
		if !ok {
			b.Fatal("Next() = false")
		}

		benchQuestion = q
	}
}

func BenchmarkQuestionNextAll(b *testing.B) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := msg.Questions()

		for {
			q, ok := it.Next()
			if !ok {
				break
			}

			benchQuestion = q
		}
	}
}

func BenchmarkQuestionName(b *testing.B) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Questions()

	q, ok := it.Next()
	if !ok {
		b.Fatal("Next() = false")
	}

	var buf [255]byte

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := q.Name(buf[:])
		if err != nil {
			b.Fatal(err)
		}

		benchQuestionNameN = n
	}
}

func BenchmarkQuestionType(b *testing.B) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Questions()

	q, ok := it.Next()
	if !ok {
		b.Fatal("Next() = false")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchQuestionType = q.Type()
	}
}

func BenchmarkQuestionClass(b *testing.B) {
	packet := testDNSPacketQuestion()

	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Questions()

	q, ok := it.Next()
	if !ok {
		b.Fatal("Next() = false")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchQuestionClass = q.Class()
	}
}
