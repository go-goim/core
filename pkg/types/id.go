package types

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

// Create maps for decoding Base58/Base32.
// This speeds up the process tremendously.
func init() {
	for i := 0; i < len(decodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(decodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// A JSONSyntaxError is returned from UnmarshalJSON if an invalid ID is provided.
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

// ErrInvalidBase58 is returned by ParseBase58 when given an invalid []byte
var ErrInvalidBase58 = errors.New("invalid base58")

// ErrInvalidBase32 is returned by ParseBase32 when given an invalid []byte
var ErrInvalidBase32 = errors.New("invalid base32")

// ID represents a unique snowflake ID.
type ID int64

// NewID returns a new snowflake ID
func NewID() ID {
	assertDefaultNode()
	return ID(defaultNode.Generate())
}

// Note: If you want to get string of base58, use Base58() instead.
//  ID.String() returns string(int64) and Base58() returns base58 string.

// Int64 returns an int64 of the snowflake ID
func (f ID) Int64() int64 {
	return int64(f)
}

// ParseInt64 converts an int64 into a snowflake ID
func ParseInt64(id int64) ID {
	return ID(id)
}

// String returns a string of the snowflake ID
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// ParseString converts a string into a snowflake ID
func ParseString(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	return ID(i), err

}

// Base2 returns a string base2 of the snowflake ID
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// ParseBase2 converts a Base2 string into a snowflake ID
func ParseBase2(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 2, 64)
	return ID(i), err
}

// Base8 returns a string base8 of the snowflake ID
func (f ID) Base8() string {
	return strconv.FormatInt(int64(f), 8)
}

// ParseBase8 converts a Base8 string into a snowflake ID
func ParseBase8(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 8, 64)
	return ID(i), err
}

// Base16 returns a string base16 of the snowflake ID
func (f ID) Base16() string {
	return strconv.FormatInt(int64(f), 16)
}

// ParseBase16 converts a Base16 string into a snowflake ID
func ParseBase16(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 16, 64)
	return ID(i), err
}

// Base32 uses the z-base-32 character set but encodes and decodes similar
// to base58, allowing it to create an even smaller result string.
// NOTE: There are many base32 implementations so be careful when
// doing any interoperation.
func (f ID) Base32() string {
	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase32 parses a base32 []byte into a snowflake ID
// NOTE: There are many base32 implementations so be careful when
// doing any interoperation.
func ParseBase32(b []byte) (ID, error) {
	var id int64

	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID(id), nil
}

// Base36 returns a base36 string of the snowflake ID
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// ParseBase36 converts a Base36 string into a snowflake ID
func ParseBase36(id string) (ID, error) {
	i, err := strconv.ParseInt(id, 36, 64)
	return ID(i), err
}

// Base58 returns a base58 string of the snowflake ID
func (f ID) Base58() string {
	if f < 58 {
		return string(encodeBase58Map[f])
	}

	b := make([]byte, 0, 11)
	for f >= 58 {
		b = append(b, encodeBase58Map[f%58])
		f /= 58
	}
	b = append(b, encodeBase58Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

// ParseBase58 parses a base58 []byte into a snowflake ID
func ParseBase58(b []byte) (ID, error) {
	var id int64

	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase58
		}
		id = id*58 + int64(decodeBase58Map[b[i]])
	}

	return ID(id), nil
}

// Base64 returns a base64 string of the snowflake ID
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

// ParseBase64 converts a base64 string into a snowflake ID
func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)

}

// Bytes returns a byte slice of the snowflake ID
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

// ParseBytes converts a byte slice into a snowflake ID
func ParseBytes(id []byte) (ID, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return ID(i), err
}

// IntBytes returns an array of bytes of the snowflake ID, encoded as a
// big endian integer.
func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

// ParseIntBytes converts an array of bytes encoded as big endian integer as
// a snowflake ID
func ParseIntBytes(id [8]byte) ID {
	return ID(int64(binary.BigEndian.Uint64(id[:])))
}

// We use int64 in internal services and use base58 for external services(like clients).
// So those marshal/unmarshal methods are converts to/from base58.

// MarshalJSON returns a json byte array string of the snowflake ID.
// Note: this is not regular MarshalJSON, it converts the ID to a base58 string,
// instead of a regular strconv.FormatInt(id,10).
func (f ID) MarshalJSON() ([]byte, error) {
	b58 := f.Base58()
	buff := make([]byte, 0, len(b58)+2)
	buff = append(buff, '"')
	buff = append(buff, b58...)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON converts a json byte array of a snowflake ID into an ID type.
// Note: this is not regular UnmarshalJSON, it converts the ID from a base58,
// instead of a regular strconv.ParseInt(id,10).
func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	id, err := ParseBase58(b[1 : len(b)-1])
	if err != nil {
		return err
	}

	*f = id
	return nil
}
