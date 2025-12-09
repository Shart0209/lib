package bloom_filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_Check(t *testing.T) {
	fb, err := New(50, EnableOptimal(0.01))
	require.NoError(t, err)

	fb.Add([]byte("test_1"))
	require.True(t, fb.Check([]byte("test_1")))

	fb.Add([]byte("test_12"))
	require.True(t, fb.Check([]byte("test_12")))

	fb.Add([]byte("test_30"))
	require.True(t, fb.Check([]byte("test_30")))

	require.False(t, fb.Check([]byte("test_510")))
	require.False(t, fb.Check([]byte("test_230")))
	require.False(t, fb.Check([]byte("test_131")))
}

// cpu: Intel(R) Core(TM) Ultra 7 165H
// BenchmarkBloomFilter_Check
// BenchmarkBloomFilter_Check-22    	83264004	        14.00 ns/op	       0 B/op	       0 allocs/op
func BenchmarkBloomFilter_Check(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		fb, err := New(10_000_000)
		require.NoError(b, err)
		fb.Add([]byte("test_5"))
		for pb.Next() {
			fb.Check([]byte("test_505"))
		}
	})
}

// cpu: Intel(R) Core(TM) Ultra 7 165H
// BenchmarkBloomFilter
// BenchmarkBloomFilter-22    	   21447	     53252 ns/op	 1253591 B/op	       5 allocs/op
func BenchmarkBloomFilter(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fb, err := New(10_000_000)
			require.NoError(b, err)
			fb.Add([]byte("test_5"))
			fb.Check([]byte("test_505"))
		}
	})
}
