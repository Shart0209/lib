package dns

import "testing"

var (
	benchRecordSink Record

	benchIntSink   int
	benchTypeSink  Type
	benchClassSink ClassType
	benchTTLSink   uint32

	benchASink    [4]byte
	benchAAAASink [16]byte

	benchErrSink error
)

func BenchmarkRecordNextOne(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := msg.Answers()
		record, ok := it.Next()
		if !ok {
			b.Fatal("Next() = false")
		}
		benchRecordSink = record
	}
}

func BenchmarkRecordNextAll(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := msg.Answers()
		for {
			record, ok := it.Next()
			if !ok {
				break
			}
			benchRecordSink = record
		}
	}
}

func BenchmarkRecordName(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	records := msg.Answers()
	record, ok := records.Next()
	if !ok {
		b.Fatal("Next() = false")
	}

	var buf [255]byte

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := record.Name(buf[:])
		if err != nil {
			b.Fatal(err)
		}
		benchIntSink = n
	}
}

func BenchmarkRecordAccessors(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	records := msg.Answers()
	record, ok := records.Next()
	if !ok {
		b.Fatal("Next() = false")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchTypeSink = record.Type()
		benchClassSink = record.Class()
		benchTTLSink = record.TTL()
	}
}

func BenchmarkRecordA(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	records := msg.Answers()
	record, ok := records.Next()
	if !ok {
		b.Fatal("A record not found")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value, err := record.A()
		if err != nil {
			b.Fatal(err)
		}
		benchASink = [4]byte(value)
	}
}

func BenchmarkRecordAAAA(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Answers()
	if _, ok := it.Next(); !ok {
		b.Fatal("A record not found")
	}
	record, ok := it.Next()
	if !ok {
		b.Fatal("AAAA record not found")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value, err := record.AAAA()
		if err != nil {
			b.Fatal(err)
		}
		benchAAAASink = [16]byte(value)
	}
}

func BenchmarkRecordCNAME(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 2; i++ {
		if _, ok := it.Next(); !ok {
			b.Fatal("record not found")
		}
	}
	record, ok := it.Next()
	if !ok {
		b.Fatal("CNAME record not found")
	}

	var buf [255]byte

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := record.CNAME(buf[:])
		if err != nil {
			b.Fatal(err)
		}
		benchIntSink = n
	}
}

func BenchmarkRecordNS(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 3; i++ {
		if _, ok := it.Next(); !ok {
			b.Fatal("record not found")
		}
	}
	record, ok := it.Next()
	if !ok {
		b.Fatal("NS record not found")
	}

	var buf [255]byte

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := record.NS(buf[:])
		if err != nil {
			b.Fatal(err)
		}
		benchIntSink = n
	}
}

func BenchmarkRecordPTR(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 4; i++ {
		if _, ok := it.Next(); !ok {
			b.Fatal("record not found")
		}
	}
	record, ok := it.Next()
	if !ok {
		b.Fatal("PTR record not found")
	}

	var buf [255]byte

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := record.PTR(buf[:])
		if err != nil {
			b.Fatal(err)
		}
		benchIntSink = n
	}
}

func BenchmarkRecordMX(b *testing.B) {
	packet := testDNSRecordsPacket()
	msg, err := Decode(packet)
	if err != nil {
		b.Fatal(err)
	}

	it := msg.Answers()
	for i := 0; i < 5; i++ {
		if _, ok := it.Next(); !ok {
			b.Fatal("record not found")
		}
	}
	record, ok := it.Next()
	if !ok {
		b.Fatal("MX record not found")
	}

	var buf [255]byte

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value, err := record.MX(buf[:])
		if err != nil {
			b.Fatal(err)
		}
		benchIntSink = int(value.Preference)
	}
}
