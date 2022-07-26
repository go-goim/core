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
		m    sync.Map
		wg   sync.WaitGroup
		fail bool
	)
	f := func(node *Node) {
		defer wg.Done()
		for i := 0; i < 500_000; i++ {
			id := node.Generate()
			if _, ok := m.Load(id); ok {
				fail = true
				break
			}
			m.Store(id, id)
		}
	}

	wg.Add(2)
	go f(n)
	go f(n2)
	wg.Wait()
	if fail {
		t.Fatal("duplicate id found")
	}

	t.Log("ok")
}
