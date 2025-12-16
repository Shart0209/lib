package bloom_filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	test1 = "test_1"
	test2 = "test_2"
)

func TestUnit_Add(t *testing.T) {
	t.Parallel()

	fb, err := New(50_000, EnableOptimal(0.01))
	require.NoError(t, err)

	fb.Add([]byte(test1))
	require.True(t, fb.Check([]byte(test1)))
}

func TestUnit_Check(t *testing.T) {
	t.Parallel()

	fb, err := New(50_000, EnableOptimal(0.01))
	require.NoError(t, err)

	fb.Add([]byte(test1))
	require.True(t, fb.Check([]byte(test1)))

	require.False(t, fb.Check([]byte(test2)))
}

func TestUnit_Clear(t *testing.T) {
	t.Parallel()

	fb, err := New(50_000)
	require.NoError(t, err)

	fb.Add([]byte(test1))
	t.Log(fb.String([]byte(test1)))

	fb.Clear([]byte(test1))
	res2 := fb.String([]byte(test1))
	t.Log(res2)

	require.False(t, fb.Check([]byte(test1)))
}

func TestUnit_ClearAll(t *testing.T) {
	t.Parallel()

	fb, err := New(50_000)
	require.NoError(t, err)

	fb.Add([]byte(test1))
	fb.Add([]byte(test2))

	fb.ClearAll()
	require.False(t, fb.Check([]byte(test1)))
	require.False(t, fb.Check([]byte(test2)))
}

// cpu: Intel(R) Core(TM) Ultra 7 165H
// BenchmarkBloomFilter_Check_Success-22    	60473558	        25.86 ns/op	       0 B/op	       0 allocs/op
func BenchmarkBloomFilter_Check_Success(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		fb, err := New(10_000_000)
		require.NoError(b, err)
		fb.Add([]byte(test1))
		for pb.Next() {
			fb.Check([]byte(test1))
		}
	})
}

// cpu: Intel(R) Core(TM) Ultra 7 165H
// BenchmarkBloomFilter_Check_Fail-22    	93488436	        16.21 ns/op	       0 B/op	       0 allocs/op
func BenchmarkBloomFilter_Check_Fail(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		fb, err := New(10_000_000)
		require.NoError(b, err)
		fb.Add([]byte(test1))
		for pb.Next() {
			fb.Check([]byte(test2))
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
			fb.Add([]byte(test1))
			fb.Check([]byte(test2))
		}
	})
}
