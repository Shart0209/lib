package base62

import (
	"bytes"
)

type Base62 interface {
	Encode(src uint64) string
	Decode(src string) uint64
}

const (
	BASE62 string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base   uint64 = 62
)

type base62 struct {
	enc []byte
}

// New if the encoder is empty then BASE62 is used
func New(encoder string) Base62 {
	if encoder == "" {
		encoder = BASE62
	}

	enc := make([]byte, len(encoder))
	copy(enc[:], encoder)

	return &base62{
		enc: enc,
	}
}

func (b *base62) Encode(src uint64) string {
	res := make([]byte, 0, 11)
	for src > 0 {
		res = append(res, b.enc[src%base])
		src = src / base
	}

	return string(res)
}

func (b *base62) Decode(src string) (intermediate uint64) {
	for i := len(src) - 1; i >= 0; i-- {
		intermediate = (base * intermediate) + uint64(bytes.IndexRune(b.enc, rune(src[i])))
	}

	return
}
