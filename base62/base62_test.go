package base62

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_Check(t *testing.T) {
	b := New("")
	for i := 0; i < 10_000; i++ {
		enc := b.Encode(uint64(i))
		dec := b.Decode(enc)
		require.Equal(t, uint64(i), dec)
	}
}

// cpu: Intel(R) Core(TM) Ultra 7 165H
// BenchmarkBase62_Encode_Decode
// BenchmarkBase62_Encode_Decode-22    	163143406	         7.571 ns/op	       0 B/op	       0 allocs/op
func BenchmarkBase62_Encode_Decode(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		bs := New("")
		for pb.Next() {
			enc := bs.Encode(math.MaxUint64)
			bs.Decode(enc)
		}
	})
}
