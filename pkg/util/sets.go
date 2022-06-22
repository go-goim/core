package util

type Set[T comparable] struct {
	m map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]bool),
	}
}

func (s *Set[T]) Add(elem ...T) *Set[T] {
	for _, e := range elem {
		s.m[e] = true
	}
}

func (s *Set[T]) Remove(elem ...T) *Set[T] {
	for _, e := range elem {
		delete(s.m, e)
	}
}

func (s *Set[T]) Contains(elem T) bool {
	return s.m[elem]
}
