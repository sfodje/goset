package goset

type KeyGetter[T any, U comparable] func(v T) U

type Comparator[T any] func(foundItem, newItem T) int

// Set represents an unordered set of data the operations that can be applied to it.
type Set[T any] interface {
	// Add adds one or more elements to a set
	Add(v ...T) bool

	// Len returns the number of elements in the set
	Len() int

	// Clear removes all elements from the set, resulting in an empty set
	Clear()

	// Clone returns a copy of the set
	Clone() Set[T]

	// Contains returns a boolean indicating if any of the given items are in the set
	Contains(v ...T) bool

	// Each iterates over items in the set applying the given function on each element.
	// Breaks iteration if the given function returns false
	Each(fn func(T) bool)

	// Diff returns a new set containing all items in this set, but not in the other
	Diff(other Set[T]) Set[T]

	// SymmetricDiff returns a new set containing all items that are not common to both sets.
	SymmetricDiff(other Set[T]) Set[T]

	// Equal returns a boolean indicating if both sets are equal.
	// That is, both have the same number of elements and the same elements.
	Equal(other Set[T]) bool

	// Intersect returns a new set containing only elements that exist in both sets
	Intersect(other Set[T]) Set[T]

	// IsSubset returns a boolean indicating if all elements in this set are also in the other set.
	IsSubset(other Set[T]) bool

	// IsProperSubset returns a boolean indicating if all elements in this set are also in the other set,
	// but both sets are unequal
	IsProperSubset(other Set[T]) bool

	// IsSuperset returns a boolean indicating if all elements in the other set are also in this set
	IsSuperset(other Set[T]) bool

	// IsProperSuperset returns a boolean indicating if all elements in the other set are also in this set,
	// but the sets are not equal
	IsProperSuperset(other Set[T]) bool

	// Iter returns a channel of all the elements in the set which allows the caller to range over the elements
	Iter() <-chan T

	// Pop removes and returns an arbitrary item from the set
	Pop() (T, bool)

	// Remove removes the given item from the set
	Remove(v ...T)

	// Union returns a new set containing all elements from both sets
	Union(other Set[T]) Set[T]

	// ToSlice returns a slice containing all elements in the set
	ToSlice() []T

	// String returns a string representation of the set
	String() string
}

func NewSet[T comparable](v ...T) Set[T] {
	set := newSafeSimpleSet[T]()
	set.Add(v...)
	return set
}

func NewThreadUnsafeSet[T comparable](v ...T) Set[T] {
	set := newUnsafeSimpleSet[T]()
	set.Add(v...)
	return set
}

func NewPrioritySet[T any, U comparable](keyGetter KeyGetter[T, U], comparator Comparator[T]) Set[T] {
	return newSafePrioritySet(keyGetter, comparator)
}

func NewThreadUnsafePrioritySet[T any, U comparable](keyGetter KeyGetter[T, U], comparator Comparator[T]) Set[T] {
	return newUnsafePrioritySet(keyGetter, comparator)
}
