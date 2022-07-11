package snowflake

import (
	"sync"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	n, err := NewNode(1)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		n.Generate()
	}
}

func TestNode_Generate(t *testing.T) {
	n, err := NewNode(1)
	if err != nil {
		t.Fatal(err)
	}

	id := n.Generate()
	if id == 0 {
		t.Fatal("id should not be 0")
	}

	t.Log(id.String())
	t.Log(id.Time())
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

func TestIdDuplicate(t *testing.T) {
	n, err := NewNode(1)
	if err != nil {
		t.Fatal(err)
	}

	n2, err := NewNode(2)
	if err != nil {
		t.Fatal(err)
	}
	var (
		m  sync.Map
		wg sync.WaitGroup
	)
	f := func(node *Node) {
		defer wg.Done()
		for i := 0; i < 500_000; i++ {
			id := node.Generate()
			if _, ok := m.Load(id); ok {
				t.Fatal("id should be unique")
			}
			m.Store(id, id)
		}
	}

	wg.Add(2)
	go f(n)
	go f(n2)
	wg.Wait()

	t.Log("ok")
}
