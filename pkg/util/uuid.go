package util

import (
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

var uuidEncoder shortuuid.Encoder

func init() {
	// use base58 as default.
	uuidEncoder = base58Encoder{}
}

// UUID returns a short version of UUID v4.
func UUID() string {
	return shortuuid.NewWithEncoder(uuidEncoder) // 6R7VqaQHbzC1xwA5UueGe6
}

type base58Encoder struct{}

func (base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func (base58Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}
