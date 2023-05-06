package goset

import (
	"sync"
)

type safeSet[T any, U comparable] struct {
	sync.RWMutex
	set Set[T]
}

// Assert concrete type:safeSet adheres to Set interface.
var _ Set[int] = (*safeSet[int, string])(nil)

func newSafeResolvingSet[T any, U comparable](keyGetter KeyGetter[T, U], comparator Resolver[T]) *safeSet[T, U] {
	set := newUnsafeResolvingSet(keyGetter, comparator)
	return &safeSet[T, U]{
		set: set,
	}
}

func newSafeSimpleSet[T comparable]() *safeSet[T, struct{}] {
	set := newUnsafeSimpleSet[T]()
	return &safeSet[T, struct{}]{
		set: set,
	}
}

func (s *safeSet[T, U]) Add(v ...T) bool {
	s.Lock()
	defer s.Unlock()
	return s.set.Add(v...)
}

func (s *safeSet[T, U]) Len() int {
	s.RLock()
	defer s.RUnlock()
	return s.set.Len()
}

func (s *safeSet[T, U]) Clear() {
	s.Lock()
	defer s.Unlock()
	s.set.Clear()
}

func (s *safeSet[T, U]) Clone() Set[T] {
	s.RLock()
	defer s.RUnlock()
	unsafeClone := s.set.Clone()
	return &safeSet[T, U]{set: unsafeClone}
}

func (s *safeSet[T, U]) Contains(v ...T) bool {
	s.RLock()
	defer s.RUnlock()
	return s.set.Contains(v...)
}

func (s *safeSet[T, U]) Each(fn func(T) bool) {
	s.RLock()
	defer s.RUnlock()
	s.set.Each(fn)
}

func (s *safeSet[T, U]) Diff(other Set[T]) Set[T] {
	o := other.(*safeSet[T, U])
	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	unsafeDiff := s.set.Diff(o.set)
	return &safeSet[T, U]{set: unsafeDiff}
}

func (s *safeSet[T, U]) SymmetricDiff(other Set[T]) Set[T] {
	o := other.(*safeSet[T, U])
	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	unsafeDiff := s.set.SymmetricDiff(o.set)
	return &safeSet[T, U]{set: unsafeDiff}
}

func (s *safeSet[T, U]) Equal(other Set[T]) bool {
	o := other.(*safeSet[T, U])
	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	return s.set.Equal(o.set)
}

func (s *safeSet[T, U]) Intersect(other Set[T]) Set[T] {
	o := other.(*safeSet[T, U])
	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	unsafeIntersection := s.set.Intersect(o.set)
	return &safeSet[T, U]{set: unsafeIntersection}
}

func (s *safeSet[T, U]) IsSubset(other Set[T]) bool {
	o := other.(*safeSet[T, U])
	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	return s.set.IsSubset(o.set)
}

func (s *safeSet[T, U]) IsProperSubset(other Set[T]) bool {
	o := other.(*safeSet[T, U])
	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	return s.set.IsProperSubset(o.set)
}

func (s *safeSet[T, U]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

func (s *safeSet[T, U]) IsProperSuperset(other Set[T]) bool {
	return other.IsProperSubset(s)
}

func (s *safeSet[T, U]) Iter() <-chan T {
	s.RLock()
	defer s.RUnlock()
	return s.set.Iter()
}

func (s *safeSet[T, U]) Pop() (T, bool) {
	s.Lock()
	defer s.Unlock()
	return s.set.Pop()
}

func (s *safeSet[T, U]) Remove(v ...T) {
	s.Lock()
	defer s.Unlock()
	s.set.Remove(v...)
}

func (s *safeSet[T, U]) Union(other Set[T]) Set[T] {
	o := other.(*safeSet[T, U])
	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	return s.set.Union(o.set)

}

func (s *safeSet[T, U]) ToSlice() []T {
	s.RLock()
	defer s.RUnlock()
	return s.set.ToSlice()
}

func (s *safeSet[T, U]) String() string {
	s.RLock()
	defer s.RUnlock()
	return s.set.String()
}
