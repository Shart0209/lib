package bitmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitmap(t *testing.T) {
	bitmap := New()
	bitmap.Set(0)
	require.Equal(t, "1: 1", fmt.Sprint(bitmap))

	pos := bitmap.Indexes()
	require.Equal(t, []int{0}, pos)

	require.True(t, bitmap.Check(0))

	bitmap.Set(1)
	bitmap.Clear(0)
	t.Log(bitmap)
	require.Equal(t, "10: 2", fmt.Sprint(bitmap))

	bitmap.Set(7)
	t.Log(bitmap)

	bitmap.ClearAll()
	require.Equal(t, "0: 0", fmt.Sprint(bitmap))
}
