# Persistent Data Structures

This package provides persistent data structures to Go.

Persistent data structures are immutable, efficient data structures built
using pointers to shared data, only copying what is modified.

## Install

    $ go get github.com/d11wtq/persistent

### Vector

This implements the same data structure used by Clojure for its Vectors. The
underlying data structure is a bit-partitioned Hash Array Mapped Trie as
described by [Phil Bagwell][1].

The data structure remains immutable during appends, prepends, pops, shifts and
updates because the tree-based nature of the data structure allows paths to
stored values to be copied, without copying the entire tree.

Most operations perform in O(log(n)) time, however they are effectively
constant time due to the fact the implementation uses O(log32(n)).

``` go
import (
	"fmt"
	"github.com/d11wtq/persistent/vector"
)

// initialize a new vector with some values
vec0 := vector.New(42, 7, 12)

// append some values to the vector
vec1 := vec0.Append(12).Append(19)

// read the number of elements
fmt.Println("vec1 has %d elements", vec1.Count())

// read the value at key 1
x, err := vec1.Get(1)
if err != nil {
	panic(err)
}
fmt.Println("vec1[1] == %s", x) // 7

// set key 1 to a new value
vec2, err := vec1.Set(1, 99)
if err != nil {
	panic(err)
}

// read the new value at key 1
y, err := vec2.Get(1)
if err != nil {
	panic(err)
}
fmt.Println("vec2[1] == %s", y) // 99

// verify that the original vector is unchanged
z, err := vec1.Get(1)
if err != nil {
	panic(err)
}
fmt.Println("vec1[1] == %s", z) // 7

// pop the last element off the end of vec1
vec3 := vec1.Pop()
fmt.Println("vec3 has %d elements", vec3.Count()) // 4

// verify that the original vector is unchanged
fmt.Println("vec1 has %d elements", vec1.Count()) // 5

// prepend some values to the start of the vector
vec4 := vec3.Prepend(9).Prepend(13)

// verify that the original vector is unchanged
a, err := vec3.Get(1)
if err != nil {
	panic(err)
}
fmt.Println("vec3[1] == %s", a) // 7

// read the new value at key 0
b, err := vec4.Get(0)
if err != nil {
	panic(err)
}
fmt.Println("vec4[0] == %s", b) // 13

// shift the first element off the start of vec4
vec5 := vec4.Shift()
fmt.Println("vec5 has %d elements", vec5.Count()) // 5
```

##### Vector operations

``` go
// Create a new vector containing elems.
// Complexity: O(n)
func New(elems ...interface{}) *Vector

// Methods defined on the *Vector type.
type interface {
  // Push an element onto the end of the vector.
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Append(interface{}) *Vector

  // Push an element onto the start of the vector.
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Prepend(interface{}) *Vector

  // Remove the last element from the vector.
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Pop() *Vector

  // Remove the first element from the vector (get the tail).
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Shift() *Vector

  // Truncate the vector to at most length n.
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Truncate(uint32) *Vector

  // Remove the first n elements from the Vector.
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Drop(uint32) *Vector

  // Set the value of the element at index i.
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Set(uint32, interface{}) (*Vector, error)

  // Get the value of the element at index i.
  // Complexity: O(log(n))
  // Effectively: O(1)
  func Get(uint32) (interface{}, error)

  // Get the length of the vector.
  // Complexity: O(1)
  func Count() uint32
}
```

  [1]: http://lampwww.epfl.ch/papers/idealhashtrees.pdf
