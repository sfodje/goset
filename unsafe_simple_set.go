package goset

import (
	"fmt"
	"strings"
)

type unsafeSimpleSet[T comparable] map[T]struct{}

// Assert concrete type:unsafeSimpleSet adheres to Set interface.
var _ Set[string] = (*unsafeSimpleSet[string])(nil)

func newUnsafeSimpleSet[T comparable]() *unsafeSimpleSet[T] {
	set := make(unsafeSimpleSet[T])
	return &set
}

func (s *unsafeSimpleSet[T]) add(v ...T) {
	for _, val := range v {
		(*s)[val] = struct{}{}
	}
}

func (s *unsafeSimpleSet[T]) Add(v ...T) bool {
	prevLen := s.Len()
	s.add(v...)
	return prevLen != s.Len()
}

func (s *unsafeSimpleSet[T]) Len() int {
	return len(*s)
}

func (s *unsafeSimpleSet[T]) Clear() {
	*s = make(unsafeSimpleSet[T])
}

func (s *unsafeSimpleSet[T]) Clone() Set[T] {
	clone := make(unsafeSimpleSet[T], s.Len())
	for elem := range *s {
		clone.add(elem)
	}
	return &clone
}

func (s *unsafeSimpleSet[T]) contains(v T) bool {
	_, ok := (*s)[v]
	return ok
}

func (s *unsafeSimpleSet[T]) Contains(v ...T) bool {
	for _, val := range v {
		if !s.contains(val) {
			return false
		}
	}
	return true
}

func (s *unsafeSimpleSet[T]) Each(fn func(T) bool) {
	for elem := range *s {
		if !fn(elem) {
			break
		}
	}
}

func (s *unsafeSimpleSet[T]) Diff(other Set[T]) Set[T] {
	o := other.(*unsafeSimpleSet[T])
	diff := newUnsafeSimpleSet[T]()
	for elem := range *s {
		if !o.contains(elem) {
			diff.Add(elem)
		}
	}
	return diff
}

func (s *unsafeSimpleSet[T]) SymmetricDiff(other Set[T]) Set[T] {
	o := other.(*unsafeSimpleSet[T])
	diff := o.Diff(s)
	for elem := range *s {
		if !o.contains(elem) {
			diff.Add(elem)
		}
	}
	return diff
}

func (s *unsafeSimpleSet[T]) Equal(other Set[T]) bool {
	o := other.(*unsafeSimpleSet[T])
	if s.Len() != other.Len() {
		return false
	}
	for elem := range *s {
		if !o.contains(elem) {
			return false
		}
	}
	return true
}

func (s *unsafeSimpleSet[T]) Intersect(other Set[T]) Set[T] {
	o := other.(*unsafeSimpleSet[T])
	intersection := newUnsafeSimpleSet[T]()

	smallerSet := s
	largerSet := o
	if o.Len() < s.Len() {
		smallerSet = o
		largerSet = s
	}

	for elem := range *smallerSet {
		if largerSet.contains(elem) {
			intersection.Add(elem)
		}
	}
	return intersection
}

func (s *unsafeSimpleSet[T]) IsSubset(other Set[T]) bool {
	o := other.(*unsafeSimpleSet[T])
	if s.Len() > other.Len() {
		return false
	}
	for elem := range *s {
		if !o.contains(elem) {
			return false
		}
	}
	return true
}

func (s *unsafeSimpleSet[T]) IsProperSubset(other Set[T]) bool {
	return s.Len() < other.Len() && s.IsSubset(other)
}

func (s *unsafeSimpleSet[T]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

func (s *unsafeSimpleSet[T]) IsProperSuperset(other Set[T]) bool {
	return s.Len() > other.Len() && s.IsSuperset(other)
}

func (s *unsafeSimpleSet[T]) Iter() <-chan T {
	ch := make(chan T, s.Len())

	go func() {
		defer close(ch)
		for elem := range *s {
			ch <- elem
		}
	}()
	return ch
}

func (s *unsafeSimpleSet[T]) Pop() (T, bool) {
	for elem := range *s {
		s.Remove(elem)
		return elem, true
	}
	var zeroElem T
	return zeroElem, false
}

func (s *unsafeSimpleSet[T]) Remove(v ...T) {
	for _, val := range v {
		delete(*s, val)
	}
}

func (s *unsafeSimpleSet[T]) Union(other Set[T]) Set[T] {
	o := other.(*unsafeSimpleSet[T])
	union := newUnsafeSimpleSet[T]()

	for elem := range *s {
		union.Add(elem)
	}
	for elem := range *o {
		union.Add(elem)
	}
	return union
}

func (s *unsafeSimpleSet[T]) ToSlice() []T {
	var elems []T
	for elem := range *s {
		elems = append(elems, elem)
	}
	return elems
}

func (s *unsafeSimpleSet[T]) String() string {
	var items []string
	for elem := range *s {
		items = append(items, fmt.Sprintf("%#v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}
