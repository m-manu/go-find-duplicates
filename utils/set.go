package utils

type Set[T comparable] map[T]struct{}

func (s Set[T]) Exists(entry T) bool {
	_, exists := s[entry]
	return exists
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

func (s Set[T]) Add(entry T) {
	s[entry] = struct{}{}
}

func (s Set[T]) Delete(e T) {
	delete(s, e)
}
