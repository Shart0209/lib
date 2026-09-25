package dns

import "errors"

var (
	ErrShortPacket        = errors.New("short packet")
	ErrInvalidName        = errors.New("invalid name")
	errCalcLen            = errors.New("insufficient data for calculated length type")
	ErrBadName            = errors.New("malformed name")
	ErrBadRR              = errors.New("malformed resource record")
	errTooManyQuestions   = errors.New("too many Questions to pack (>65535)")
	errTooManyAnswers     = errors.New("too many Answers to pack (>65535)")
	errTooManyAuthorities = errors.New("too many Authorities to pack (>65535)")
	errTooManyAdditionals = errors.New("too many Additionals to pack (>65535)")
)
