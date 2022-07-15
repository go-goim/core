package types

import (
	"fmt"

	"github.com/go-goim/core/pkg/util/snowflake"
)

// A JSONSyntaxError is returned from UnmarshalJSON if an invalid ID is provided.
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

// ID represents a unique snowflake ID.
//  It contains original snowflake ID and base58 encoded snowflake ID.
//  ID implements the json.Marshaler, json.Unmarshaler, sql.Scanner, driver.Valuer,
//   fmt.Stringer interfaces.
type ID struct {
	snowflake.ID
	base58 string
}

// Note: If you want to get string of base58, use Base58() instead.
//  ID.String() returns string(int64) and Base58() returns base58 string.

func NewID(id int64) *ID {
	sid := snowflake.ID(id)
	return &ID{
		ID:     sid,
		base58: sid.Base58(),
	}
}

// GenerateID generates a new snowflake ID.
func GenerateID() *ID {
	sid := snowflake.Generate()
	return &ID{
		ID:     sid,
		base58: sid.Base58(),
	}
}

func (id *ID) Equal(other *ID) bool {
	return id.ID == other.ID
}

func (id *ID) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, len(id.base58)+2)
	buff = append(buff, '"')
	buff = append(buff, id.base58...)
	buff = append(buff, '"')
	return buff, nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var err error

	if len(data) < 2 {
		return JSONSyntaxError{data}
	}
	if data[0] != '"' || data[len(data)-1] != '"' {
		return JSONSyntaxError{data}
	}
	id.base58 = string(data[1 : len(data)-1])
	id.ID, err = snowflake.ParseBase58(data[1 : len(data)-1])
	if err != nil {
		return err
	}

	return nil
}
