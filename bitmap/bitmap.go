package bitmap

import (
	"fmt"
	"math/bits"
)

type Bitmap interface {
	Toggle(bit uint64)
	Set(bit uint64)
	Check(bit uint64) bool
	Clear(bit uint64)
	ClearAll()
	Indexes() []int
	String() string
}

type bitmap struct {
	bm uint64
}

func New() Bitmap {
	return &bitmap{}
}

func (b *bitmap) Toggle(bit uint64) {
	b.bm ^= 1 << bit
}

func (b *bitmap) Set(bit uint64) {
	b.bm |= 1 << bit
}

func (b *bitmap) Check(bit uint64) bool {
	return (b.bm & (1 << bit)) > 0
}

func (b *bitmap) ClearAll() {
	b.bm &^= 1<<64 - 1
}

func (b *bitmap) Clear(bit uint64) {
	b.bm &^= bit
}

func (b *bitmap) String() string {
	return fmt.Sprintf("%064b", b.bm)
}

func (b *bitmap) Indexes() []int {
	return indexes(b.bm)
}

func indexes(bit uint64) []int {
	var ind = make([]int, bits.OnesCount64(bit))
	pos := 0
	for bit != 0 {
		ind[pos] = bits.TrailingZeros64(bit)
		bit &= bit - 1
		pos += 1
	}
	return ind
}
