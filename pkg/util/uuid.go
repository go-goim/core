package util

import (
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/google/uuid"
)

// use base58 as default.
var uuidEncoder = base58Encoder{}

// UUID returns a short version of UUID v4.
func UUID() string {
	return uuidEncoder.Encode(uuid.New())
}

func GUID() string {
	return uuid.New().String()
}

type base58Encoder struct{}

func (base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func (base58Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}
