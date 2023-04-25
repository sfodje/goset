package goset

import (
	"fmt"
	"strings"
)

type unsafePrioritySet[T any, U comparable] struct {
	set        map[U]T
	keyGetter  KeyGetter[T, U]
	comparator Comparator[T]
}

// Assert concrete type:unsafePrioritySet adheres to Set interface.
var _ Set[int] = (*unsafePrioritySet[int, string])(nil)

func newUnsafePrioritySet[T any, U comparable](keyGetter KeyGetter[T, U], comparator Comparator[T]) *unsafePrioritySet[T, U] {
	return &unsafePrioritySet[T, U]{
		set:        make(map[U]T),
		keyGetter:  keyGetter,
		comparator: comparator,
	}
}

func (s *unsafePrioritySet[T, U]) Add(v ...T) bool {
	var ret bool
	for _, val := range v {
		key := s.keyGetter(val)
		if item, ok := s.set[key]; !ok || (s.comparator != nil && s.comparator(item, val) > 0) {
			s.set[key] = val
			ret = true
		}
	}
	return ret
}

func (s *unsafePrioritySet[T, U]) Len() int {
	return len(s.set)
}

func (s *unsafePrioritySet[T, U]) Clear() {
	s.set = make(map[U]T)
}

func (s *unsafePrioritySet[T, U]) Clone() Set[T] {
	clonedSet := newUnsafePrioritySet(s.keyGetter, s.comparator)
	for _, elem := range s.set {
		clonedSet.Add(elem)
	}
	return clonedSet
}

func (s *unsafePrioritySet[T, U]) contains(v T) bool {
	key := s.keyGetter(v)
	elem, ok := s.set[key]
	return ok && (s.comparator == nil || s.comparator(elem, v) == 0)
}

func (s *unsafePrioritySet[T, U]) Contains(v ...T) bool {
	for _, val := range v {
		if !s.contains(val) {
			return false
		}
	}
	return true
}

func (s *unsafePrioritySet[T, U]) Diff(other Set[T]) Set[T] {
	o := other.(*unsafePrioritySet[T, U])
	diff := newUnsafePrioritySet(s.keyGetter, s.comparator)
	for _, elem := range s.set {
		if !o.contains(elem) {
			diff.Add(elem)
		}
	}
	return diff
}

func (s *unsafePrioritySet[T, U]) SymmetricDiff(other Set[T]) Set[T] {
	o := other.(*unsafePrioritySet[T, U])
	diff := o.Diff(s)
	for _, elem := range s.set {
		if !o.contains(elem) {
			diff.Add(elem)
		}
	}
	return diff
}

func (s *unsafePrioritySet[T, U]) Each(fn func(T) bool) {
	for _, elem := range s.set {
		if !fn(elem) {
			break
		}
	}
}

func (s *unsafePrioritySet[T, U]) Equal(other Set[T]) bool {
	o := other.(*unsafePrioritySet[T, U])
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

func (s *unsafePrioritySet[T, U]) Intersect(other Set[T]) Set[T] {
	o := other.(*unsafePrioritySet[T, U])
	intersection := newUnsafePrioritySet(s.keyGetter, s.comparator)

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

func (s *unsafePrioritySet[T, U]) IsSubset(other Set[T]) bool {
	o := other.(*unsafePrioritySet[T, U])
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
func (s *unsafePrioritySet[T, U]) IsProperSubset(other Set[T]) bool {
	return s.Len() < other.Len() && s.IsSubset(other)
}

func (s *unsafePrioritySet[T, U]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

func (s *unsafePrioritySet[T, U]) IsProperSuperset(other Set[T]) bool {
	return s.Len() > other.Len() && s.IsSuperset(other)
}

func (s *unsafePrioritySet[T, U]) Iter() <-chan T {
	ch := make(chan T, s.Len())

	go func() {
		defer close(ch)
		for _, elem := range s.set {
			ch <- elem
		}
	}()
	return ch
}

func (s *unsafePrioritySet[T, U]) Remove(v ...T) {
	for _, val := range v {
		key := s.keyGetter(val)
		delete(s.set, key)
	}
}

func (s *unsafePrioritySet[T, U]) Pop() (T, bool) {
	for _, elem := range s.set {
		s.Remove(elem)
		return elem, true
	}
	var zeroElem T
	return zeroElem, false
}

func (s *unsafePrioritySet[T, U]) String() string {
	var items []string
	for _, elem := range s.set {
		items = append(items, fmt.Sprintf("%#v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

func (s *unsafePrioritySet[T, U]) Union(other Set[T]) Set[T] {
	o := other.(*unsafePrioritySet[T, U])
	union := newUnsafePrioritySet(s.keyGetter, s.comparator)

	for _, elem := range s.set {
		union.Add(elem)
	}
	for _, elem := range o.set {
		union.Add(elem)
	}
	return union
}
func (s *unsafePrioritySet[T, U]) ToSlice() []T {
	var elems []T
	for _, elem := range s.set {
		elems = append(elems, elem)
	}
	return elems
}
