package goset_test

import (
	"regexp"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sfodje/goset"
)

var (
	testTypeSetStringRegex = regexp.MustCompile(`^Set\{(&goset_test.TestType\{ID:\d+, Name:".+", Importance:\d+},*)*}$`)
	intSetStringRegex      = regexp.MustCompile(`^Set{\d+, \d+, \d+}$`)
)

func TestSets(t *testing.T) {
	testCases := []struct {
		name   string
		newSet func(v ...int) goset.Set[int]
	}{
		{
			name:   "UnsafeSimpleSet",
			newSet: func(v ...int) goset.Set[int] { return goset.NewThreadUnsafeSet[int](v...) },
		},
		{
			name:   "SafeSimpleSet",
			newSet: func(v ...int) goset.Set[int] { return goset.NewSet[int](v...) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("Add", func(t *testing.T) {
				set := tc.newSet(1, 2, 3, 4, 5)
				set.Add(3, 4, 5, 6, 7, 8, 9, 10)

				expectedItems := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				actualItems := set.ToSlice()
				sort.Ints(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)
			})

			t.Run("Len", func(t *testing.T) {
				set := tc.newSet(1, 2, 3, 4)
				assert.Equal(t, 4, set.Len())

				set.Add(3, 4, 5, 6, 7)
				assert.Equal(t, 7, set.Len())
				assert.Len(t, set.ToSlice(), set.Len())
			})

			t.Run("Clear", func(t *testing.T) {
				set := tc.newSet(1, 2, 3, 4, 5, 6)
				assert.Equal(t, 6, set.Len())

				set.Clear()
				assert.Zero(t, set.Len())
			})

			t.Run("Clone", func(t *testing.T) {
				setA := tc.newSet(3, 4, 5, 6)
				setB := setA.Clone()

				assert.Equal(t, setA.Len(), setB.Len())
				for _, v := range setA.ToSlice() {
					assert.True(t, setB.Contains(v))
				}
			})

			t.Run("Contains", func(t *testing.T) {
				set := tc.newSet(4, 5, 6, 7, 8)
				assert.True(t, set.Contains(4, 5, 6))
				assert.True(t, set.Contains(7))
				assert.False(t, set.Contains(10, 11, 12))
				assert.False(t, set.Contains(13))
			})

			t.Run("Each", func(t *testing.T) {
				set := tc.newSet(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

				var even []int
				set.Each(func(v int) bool {
					if v%2 == 0 {
						even = append(even, v)
					}
					return true
				})
				sort.Ints(even)
				assert.EqualValues(t, []int{2, 4, 6, 8, 10}, even)

				count := 0
				var items []int
				set.Each(func(v int) bool {
					count += 1
					items = append(items, v)
					if count >= 3 {
						return false
					}
					return true
				})
				assert.Len(t, items, 3)
			})

			t.Run("Diff", func(t *testing.T) {
				setA := tc.newSet(1, 2, 3)
				setB := tc.newSet(2, 3, 4, 5)

				diff := setA.Diff(setB)
				expectedItems := []int{1}
				actualItems := diff.ToSlice()
				assert.EqualValues(t, expectedItems, actualItems)
			})

			t.Run("SymmetricDiff", func(t *testing.T) {
				setA := tc.newSet(1, 2, 3)
				setB := tc.newSet(2, 3, 4, 5)

				diff := setA.SymmetricDiff(setB)
				expectedItems := []int{1, 4, 5}
				actualItems := diff.ToSlice()
				sort.Ints(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)
			})

			t.Run("Equal", func(t *testing.T) {
				setA := tc.newSet(1, 2, 3)
				setB := tc.newSet(1, 2, 3)
				setC := tc.newSet(2, 4, 6, 7)
				setD := tc.newSet(7, 8, 9)

				assert.True(t, setA.Equal(setB))
				assert.True(t, setB.Equal(setA))
				assert.False(t, setA.Equal(setC))
				assert.False(t, setC.Equal(setB))
				assert.False(t, setD.Equal(setA))
			})

			t.Run("Intersect", func(t *testing.T) {
				setA := tc.newSet(1, 2, 3)
				setB := tc.newSet(1, 3, 4, 5)

				intersect := setA.Intersect(setB)
				expectedItems := []int{1, 3}
				actualItems := intersect.ToSlice()
				sort.Ints(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)

				intersect = setB.Intersect(setA)
				actualItems = intersect.ToSlice()
				sort.Ints(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)
			})

			t.Run("IsSubset/IsProperSubset/IsSuperset/IsProperSuperset", func(t *testing.T) {
				setA := tc.newSet(1, 2, 3)
				setB := tc.newSet(2, 3)
				setC := tc.newSet(4, 5, 6)
				assert.True(t, setB.IsSubset(setA))
				assert.False(t, setA.IsSubset(setB))
				assert.True(t, setA.IsSuperset(setB))
				assert.True(t, setA.IsProperSuperset(setB))
				assert.True(t, setB.IsProperSubset(setA))
				assert.True(t, setA.IsSuperset(setB))
				assert.False(t, setA.IsSubset(setC))
			})

			t.Run("Iter", func(t *testing.T) {
				items := []int{1, 2, 3, 4, 5, 6, 7}
				set := tc.newSet(items...)
				for item := range set.Iter() {
					assert.Contains(t, items, item)
				}
			})

			t.Run("Pop", func(t *testing.T) {
				set := tc.newSet()
				assert.Zero(t, set.Len())
				v, ok := set.Pop()
				assert.Zero(t, v)
				assert.False(t, ok)

				set.Add(1, 2, 3, 4)
				assert.Equal(t, 4, set.Len())
				v, ok = set.Pop()
				assert.True(t, ok)
				assert.LessOrEqual(t, v, 4)
				assert.Equal(t, 3, set.Len())
			})

			t.Run("Union", func(t *testing.T) {
				setA := tc.newSet(1, 3, 5, 7, 9)
				setB := tc.newSet(2, 4, 6, 8)

				unionA := setA.Union(setB)
				unionB := setB.Union(setA)

				expectedItems := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
				actualA := unionA.ToSlice()
				actualB := unionB.ToSlice()
				sort.Ints(actualA)
				sort.Ints(actualB)

				assert.EqualValues(t, expectedItems, actualA)
				assert.EqualValues(t, expectedItems, actualB)
			})

			t.Run("String", func(t *testing.T) {
				set := tc.newSet(1, 2, 3)
				assert.Regexp(t, intSetStringRegex, set.String())
			})
		})
	}
}

func TestResolvingSets(t *testing.T) {
	testCases := []struct {
		name   string
		newSet func() goset.Set[*TestType]
	}{
		{
			name: "UnsafeResolvingSet",
			newSet: func() goset.Set[*TestType] {
				return goset.NewThreadUnsafeResolvingSet(func(item *TestType) int {
					return item.ID
				}, func(foundItem, newItem *TestType) (*TestType, bool) {
					if newItem.Importance > foundItem.Importance {
						return newItem, true
					}
					return foundItem, false
				})
			},
		},
		{
			name: "SafeResolvingSet",
			newSet: func() goset.Set[*TestType] {
				return goset.NewResolvingSet(
					func(item *TestType) int {
						return item.ID
					},
					func(foundItem, newItem *TestType) (*TestType, bool) {
						if newItem.Importance > foundItem.Importance {
							return newItem, true
						}
						return foundItem, false
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("Add", func(t *testing.T) {
				set := tc.newSet()

				set.Add(testItems...)
				expectedItems := []*TestType{testItems[5], testItems[3], testItems[2]}
				actualItems := set.ToSlice()
				sortTestItems(actualItems)
				assert.Equal(t, 3, set.Len())
				assert.EqualValues(t, expectedItems, actualItems)
			})

			t.Run("Clear", func(t *testing.T) {
				set := tc.newSet()

				set.Add(testItems...)
				expectedItems := []*TestType{testItems[5], testItems[3], testItems[2]}
				actualItems := set.ToSlice()
				sortTestItems(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)

				set.Clear()
				assert.Equal(t, set.Len(), 0)
				assert.Empty(t, set.ToSlice())
			})

			t.Run("Clone", func(t *testing.T) {
				set := tc.newSet()

				set.Add(testItems...)
				clone := set.Clone()
				expectedItems := []*TestType{testItems[5], testItems[3], testItems[2]}
				setItems := set.ToSlice()
				clonedItems := clone.ToSlice()
				sortTestItems(setItems)
				sortTestItems(clonedItems)

				assert.Equal(t, set.Len(), clone.Len())
				assert.True(t, set.Equal(clone))
				assert.EqualValues(t, expectedItems, setItems)
				assert.EqualValues(t, expectedItems, clonedItems)
			})

			t.Run("Contains", func(t *testing.T) {
				set := tc.newSet()
				set.Add(testItems...)
				assert.True(t, set.Contains(testItems[5], testItems[3], testItems[2]))
				assert.True(t, set.Contains(testItems[3]))
				assert.True(t, set.Contains(testItems...))
				assert.False(t, set.Contains(&TestType{ID: 100, Name: "One Hundred", Importance: 1}))
			})

			t.Run("Diff/SymmetricDiff", func(t *testing.T) {
				setA := tc.newSet()
				setA.Add(testItems...)

				setB := tc.newSet()
				setB.Add(testItems[0], testItems[2], testItems[4], &TestType{ID: 100, Name: "One Hundred", Importance: 1})

				diff := setA.Diff(setB)
				expectedItems := []*TestType{testItems[3]}
				actualItems := diff.ToSlice()
				sortTestItems(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)

				diff = setB.Diff(setA)
				expectedItems = []*TestType{{ID: 100, Name: "One Hundred", Importance: 1}}
				actualItems = diff.ToSlice()
				sortTestItems(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)

				diff = setA.SymmetricDiff(setB)
				expectedItems = []*TestType{testItems[3], {ID: 100, Name: "One Hundred", Importance: 1}}
				actualItems = diff.ToSlice()
				sortTestItems(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)
			})

			t.Run("Each", func(t *testing.T) {
				var items []*TestType
				set := tc.newSet()
				set.Add(testItems...)

				set.Each(func(item *TestType) bool {
					if item.ID == 1 {
						items = append(items, item)
						return false
					}
					return true
				})

				assert.Len(t, items, 1)
				assert.Equal(t, items[0], testItems[5])

				set.Each(func(item *TestType) bool {
					items = append(items, item)
					return true
				})
				assert.Len(t, items, 4)

				expectedItems := []*TestType{testItems[5], testItems[5], testItems[3], testItems[2]}
				sortTestItems(items)
				assert.EqualValues(t, expectedItems, items)
			})

			t.Run("Intersect", func(t *testing.T) {
				setA := tc.newSet()
				setA.Add(testItems...)

				setB := tc.newSet()
				setB.Add(testItems[5], testItems[2], testItems[4], &TestType{ID: 100, Name: "One Hundred", Importance: 1},
					&TestType{ID: 200, Name: "Two Hundred", Importance: 1})

				intersect := setB.Intersect(setA)
				assert.Len(t, intersect.ToSlice(), 2)

				assert.Contains(t, intersect.ToSlice(), testItems[2])
				assert.Contains(t, intersect.ToSlice(), testItems[5])
			})

			t.Run("IsSubset/IsProperSubset/IsSuperset/IsProperSuperset", func(t *testing.T) {
				setA := tc.newSet()
				setA.Add(testItems...)

				setB := tc.newSet()
				setB.Add(testItems[5])

				setC := tc.newSet()
				setC.Add(testItems[5], testItems[3], &TestType{ID: 100, Name: "One Hundred", Importance: 1})
				assert.True(t, setB.IsProperSubset(setA))
				assert.True(t, setB.IsSubset(setA))
				assert.True(t, setA.IsProperSuperset(setB))
				assert.True(t, setA.IsSuperset(setB))
				assert.True(t, setA.IsSubset(setA))
				assert.True(t, setB.IsSubset(setB))
				assert.True(t, setC.IsProperSuperset(setB))
				assert.True(t, setC.IsSuperset(setB))
				assert.False(t, setC.IsSubset(setA))
				assert.False(t, setA.IsProperSubset(setA))
				assert.False(t, setB.IsProperSuperset(setA))
				assert.False(t, setA.IsSubset(setB))
			})

			t.Run("Equal", func(t *testing.T) {
				setA := tc.newSet()
				setA.Add(testItems...)

				setB := tc.newSet()
				setB.Add(testItems...)
				assert.True(t, setA.Equal(setA))
				assert.True(t, setA.Equal(setB))

				setB.Add(&TestType{ID: 100, Name: "One Hundred", Importance: 1})
				assert.False(t, setA.Equal(setB))

				_, _ = setB.Pop()
				setB.Add(&TestType{ID: 100, Name: "One Hundred", Importance: 1})
				assert.False(t, setA.Equal(setB))
			})

			t.Run("Pop", func(t *testing.T) {
				set := tc.newSet()
				v, ok := set.Pop()
				assert.False(t, ok)
				assert.Nil(t, v)
			})

			t.Run("Iter", func(t *testing.T) {
				set := tc.newSet()
				set.Add(testItems...)
				for item := range set.Iter() {
					assert.Contains(t, testItems, item)
				}
			})

			t.Run("Len", func(t *testing.T) {
				set := tc.newSet()
				assert.Equal(t, set.Len(), 0)

				set.Add(testItems...)
				assert.Equal(t, set.Len(), 3)
			})

			t.Run("Union", func(t *testing.T) {
				setA := tc.newSet()
				setA.Add(testItems...)

				setB := tc.newSet()
				setB.Add(testItems...)
				setB.Add(&TestType{ID: 100, Name: "One Hundred", Importance: 1})

				union := setA.Union(setB)
				expectedItems := []*TestType{testItems[5], testItems[3], testItems[2], {ID: 100, Name: "One Hundred", Importance: 1}}
				actualItems := union.ToSlice()
				sortTestItems(actualItems)
				assert.EqualValues(t, expectedItems, actualItems)
			})

			t.Run("String", func(t *testing.T) {
				set := tc.newSet()
				set.Add(testItems...)
				assert.Regexp(t, testTypeSetStringRegex, set.String())
			})

			t.Run("Remove", func(t *testing.T) {
				set := tc.newSet()
				set.Add(testItems...)
				expectedItems := []*TestType{testItems[5], testItems[3], testItems[2]}
				actualItems := set.ToSlice()
				sortTestItems(actualItems)
				assert.Equal(t, 3, set.Len())
				assert.EqualValues(t, expectedItems, actualItems)

				set.Remove(testItems[5], testItems[2])
				assert.Equal(t, 1, set.Len())
				assert.Equal(t, testItems[3], set.ToSlice()[0])
			})
		})
	}
}

type TestType struct {
	ID         int
	Name       string
	Importance int
}

var testItems = []*TestType{
	{
		ID:         1,
		Name:       "One",
		Importance: 1,
	},
	{
		ID:         2,
		Name:       "Two",
		Importance: 1,
	},
	{
		ID:         3,
		Name:       "Three",
		Importance: 1,
	},
	{
		ID:         2,
		Name:       "Two",
		Importance: 2,
	},
	{
		ID:         1,
		Name:       "One",
		Importance: 2,
	},
	{
		ID:         1,
		Name:       "One",
		Importance: 3,
	},
}

func sortTestItems(items []*TestType) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
}
