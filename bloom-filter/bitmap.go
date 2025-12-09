package bloom_filter

import "fmt"

type bitmap uint64

func (b *bitmap) Set(bit uint64) {
	*b |= 1 << bit
}

func (b *bitmap) Check(bit uint64) bool {
	return (*b & (1 << bit)) > 0
}

func (b *bitmap) String() string {
	return fmt.Sprintf("%064b", *b)
}
