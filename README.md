# Go-Set

A set implementation library for golang. Based on the [golang-set](https://github.com/deckarep/golang-set),
this library provides more configurable options for sets of structs, allowing you to specify how uniqueness is defined 
and the priority between elements of the same "uniqueness".

## Usage

```go
// Initiate new thread-safe set. 
set := goset.NewSet(1, 2, 3, 4)

// or
set = goset.NewSet()
set.Add(1, 2, 3, 4)

// Initiate a thread unsafe set (for faster performance in a non-concurrent environment) 
threadUnsafeSet := goset.NewThreadUnsafeSet("one", "two", "three")

type TestStruct struct {
    ID int
    Name string
    Importance int
}
structs := []TestStruct {
    {
        ID: 1,
        Name: "One",
        Importance: 1,
    },
    {
        ID: 2,
        Name: "Two",
        Importance: 2,
    },
    {
        ID: 1,
        Name: "One",
        Importance: 2,
    },
}

// Initiate a set of structs
structSet := goset.NewSet(structs...)

// Initiate a priority set of structs
prioritySet := goset.NewPrioritySet(
	func(item TestStruct) int {
        return item.ID
    }, 
	func(foundItem, newItem TestStruct) int {
        if newItem.Importance > foundItem.Importance {
        return 1
    }
    if foundItem == newItem {
        return 0
    }
    return -1
	})

// The NewPrioritySet function takes a keyGetter and a Comparer function
// The KeyGetter function returns the attribute by which set of structs will be determined unique
// The Comparer function returns -1, 0, 1 to determine the level of importance of each attribute in the set.
// When adding new elements, the more important elements will always replace the less important ones.
prioritySet.Add(structs...)

fmt.Println(prioritySet.String())

// Output: Set{TestStruct{ID:1, Name:"One", Importance:2}, TestStruct{ID:2, Name:"Two", Importance:2}}
```

```go

seta := goset.NewSet(1, 2, 3, 4)
setb := goset.NewSet(2,3,4,5)
setc := goset.NewSet(5, 6)

union := seta.Union(setb).Union(setc)

fmt.Println(union)
// Output: Set{1, 2, 3, 4, 5, 6}

seta.Contains(1)
// Output: true

setb.Contains(2, 3, 4)
// Output: true

setc.Contains(1, 5, 6)
// Output: false

diff := setb.Diff(setc)
fmt.Println(diff)
// Output: Set{2, 3, 4}

```