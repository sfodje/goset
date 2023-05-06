package goset

import (
	"fmt"
	"strings"
)

type unsafeResolvingSet[T any, U comparable] struct {
	set       map[U]T
	keyGetter KeyGetter[T, U]
	resolver  Resolver[T]
}

// Assert concrete type:unsafeResolvingSet adheres to Set interface.
var _ Set[int] = (*unsafeResolvingSet[int, string])(nil)

func newUnsafeResolvingSet[T any, U comparable](keyGetter KeyGetter[T, U], resolver Resolver[T]) *unsafeResolvingSet[T, U] {
	return &unsafeResolvingSet[T, U]{
		set:       make(map[U]T),
		keyGetter: keyGetter,
		resolver:  resolver,
	}
}

func (s *unsafeResolvingSet[T, U]) Add(v ...T) bool {
	var ret bool
	for _, val := range v {
		key := s.keyGetter(val)
		foundItem, ok := s.set[key]
		// if item already exists in set, resolve and add
		if ok && s.resolver != nil {
			if newItem, ok := s.resolver(foundItem, val); ok {
				s.set[key] = newItem
				ret = true
			}
		}
		// if item not in set, add
		if !ok {
			s.set[key] = val
			ret = true
		}
	}
	return ret
}

func (s *unsafeResolvingSet[T, U]) Len() int {
	return len(s.set)
}

func (s *unsafeResolvingSet[T, U]) Clear() {
	s.set = make(map[U]T)
}

func (s *unsafeResolvingSet[T, U]) Clone() Set[T] {
	clonedSet := newUnsafeResolvingSet(s.keyGetter, s.resolver)
	for _, elem := range s.set {
		clonedSet.Add(elem)
	}
	return clonedSet
}

func (s *unsafeResolvingSet[T, U]) contains(v T) bool {
	key := s.keyGetter(v)
	_, ok := s.set[key]
	// if ok && s.resolver != nil {
	// 	_, itemIsDifferent := s.resolver(foundItem, v)
	// 	return !itemIsDifferent
	// }
	return ok
}

func (s *unsafeResolvingSet[T, U]) Contains(v ...T) bool {
	for _, val := range v {
		if !s.contains(val) {
			return false
		}
	}
	return true
}

func (s *unsafeResolvingSet[T, U]) Diff(other Set[T]) Set[T] {
	o := other.(*unsafeResolvingSet[T, U])
	diff := newUnsafeResolvingSet(s.keyGetter, s.resolver)
	for _, elem := range s.set {
		if !o.contains(elem) {
			diff.Add(elem)
		}
	}
	return diff
}

func (s *unsafeResolvingSet[T, U]) SymmetricDiff(other Set[T]) Set[T] {
	o := other.(*unsafeResolvingSet[T, U])
	diff := o.Diff(s)
	for _, elem := range s.set {
		if !o.contains(elem) {
			diff.Add(elem)
		}
	}
	return diff
}

func (s *unsafeResolvingSet[T, U]) Each(fn func(T) bool) {
	for _, elem := range s.set {
		if !fn(elem) {
			break
		}
	}
}

func (s *unsafeResolvingSet[T, U]) Equal(other Set[T]) bool {
	o := other.(*unsafeResolvingSet[T, U])
	if s.Len() != other.Len() {
		return false
	}
	for _, elem := range s.set {
		if !o.contains(elem) {
			return false
		}
	}
	return true
}

func (s *unsafeResolvingSet[T, U]) Intersect(other Set[T]) Set[T] {
	o := other.(*unsafeResolvingSet[T, U])
	intersection := newUnsafeResolvingSet(s.keyGetter, s.resolver)

	smallerSet := s
	largerSet := o
	if o.Len() < s.Len() {
		smallerSet = o
		largerSet = s
	}

	for _, elem := range smallerSet.set {
		if largerSet.contains(elem) {
			intersection.Add(elem)
		}
	}

	return intersection
}

func (s *unsafeResolvingSet[T, U]) IsSubset(other Set[T]) bool {
	o := other.(*unsafeResolvingSet[T, U])
	if s.Len() > other.Len() {
		return false
	}
	for _, elem := range s.set {
		if !o.contains(elem) {
			return false
		}
	}
	return true
}
func (s *unsafeResolvingSet[T, U]) IsProperSubset(other Set[T]) bool {
	return s.Len() < other.Len() && s.IsSubset(other)
}

func (s *unsafeResolvingSet[T, U]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

func (s *unsafeResolvingSet[T, U]) IsProperSuperset(other Set[T]) bool {
	return s.Len() > other.Len() && s.IsSuperset(other)
}

func (s *unsafeResolvingSet[T, U]) Iter() <-chan T {
	ch := make(chan T, s.Len())

	go func() {
		defer close(ch)
		for _, elem := range s.set {
			ch <- elem
		}
	}()
	return ch
}

func (s *unsafeResolvingSet[T, U]) Remove(v ...T) {
	for _, val := range v {
		key := s.keyGetter(val)
		delete(s.set, key)
	}
}

func (s *unsafeResolvingSet[T, U]) Pop() (T, bool) {
	for _, elem := range s.set {
		s.Remove(elem)
		return elem, true
	}
	var zeroElem T
	return zeroElem, false
}

func (s *unsafeResolvingSet[T, U]) String() string {
	var items []string
	for _, elem := range s.set {
		items = append(items, fmt.Sprintf("%#v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

func (s *unsafeResolvingSet[T, U]) Union(other Set[T]) Set[T] {
	o := other.(*unsafeResolvingSet[T, U])
	union := newUnsafeResolvingSet(s.keyGetter, s.resolver)

	for _, elem := range s.set {
		union.Add(elem)
	}
	for _, elem := range o.set {
		union.Add(elem)
	}
	return union
}
func (s *unsafeResolvingSet[T, U]) ToSlice() []T {
	var elems []T
	for _, elem := range s.set {
		elems = append(elems, elem)
	}
	return elems
}
