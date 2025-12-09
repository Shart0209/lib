package bitmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitmap(t *testing.T) {
	bitmap := New()
	bitmap.Set(0)
	require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001: 1", fmt.Sprint(bitmap))

	pos := bitmap.Indexes()
	require.Equal(t, []int{0}, pos)

	require.True(t, bitmap.Check(0))

	bitmap.Clear(1)
	fmt.Println(bitmap)
	require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000: 0", fmt.Sprint(bitmap))

	bitmap.Set(7)
	fmt.Println(bitmap)
	
	bitmap.ClearAll()
	require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000: 0", fmt.Sprint(bitmap))
}
