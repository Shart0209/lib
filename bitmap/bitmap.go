package bitmap

import (
	"fmt"
	"math/bits"
)

type IBitmap interface {
	Toggle(bit uint64)
	Set(bit uint64)
	Check(bit uint64) bool
	Clear(bit uint64)
	ClearAll()
	Indexes() []int
	String() string
}

type Bitmap uint64

func New() IBitmap {
	var bitmap Bitmap
	return &bitmap
}

func (b *Bitmap) Toggle(bit uint64) {
	*b ^= 1 << bit
}

func (b *Bitmap) Set(bit uint64) {
	*b |= 1 << bit
}

func (b *Bitmap) Check(bit uint64) bool {
	return (*b & (1 << bit)) > 0
}

func (b *Bitmap) ClearAll() {
	*b &^= 1<<64 - 1
}

func (b *Bitmap) Clear(bit uint64) {
	*b &^= 1 << bit
}

func (b *Bitmap) String() string {
	return fmt.Sprintf("%b: %d", *b, *b)
}

func (b *Bitmap) Indexes() []int {
	return indexes(*b)
}

func indexes(bit Bitmap) []int {
	var ind = make([]int, bits.OnesCount64(uint64(bit)))
	pos := 0
	for bit != 0 {
		ind[pos] = bits.TrailingZeros64(uint64(bit))
		bit &= bit - 1
		pos += 1
	}
	return ind
}
