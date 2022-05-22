package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	uuid := UUID()
	t.Log(uuid)
	t.Log(GUID())
	assert.Equal(t, 22, len(uuid))
}

/*
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkUUID
BenchmarkUUID-12    	 1110436	      1008 ns/op
PASS
*/
func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID()
	}
}
