package types

import (
	"testing"

	"github.com/go-goim/core/pkg/types/snowflake"
)

func TestID_Fromats(t *testing.T) {
	n, err := snowflake.NewNode(1)
	if err != nil {
		t.Fatal(err)
	}

	id := ID(n.Generate())
	if id == 0 {
		t.Fatal("id should not be 0")
	}

	t.Log(id.String())
	t.Log(snowflake.TimeFromID(id.Int64()))
	t.Log("base2", id.Base2())
	t.Log("base8", id.Base8())
	t.Log("base16", id.Base16())
	t.Log("base32", id.Base32())
	t.Log("base36", id.Base36())
	t.Log("base58", id.Base58())
	t.Log("base64", id.Base64())
	id2, err := ParseBase58([]byte(id.Base58()))
	if err != nil {
		t.Fatal(err)
	}
	if id != id2 {
		t.Fatal("id should be equal to id2")
	}
	for i := 0; i < 10; i++ {
		t.Log(n.Generate())
	}
}
