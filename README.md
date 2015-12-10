# Persistent Data Structures

This package provides persistent data structures to Go.

Persistent data structures are immutable, efficient data structures built
using pointers to shared data, only copying what is modified.

## Install

    $ go get github.com/d11wtq/persistent

### Data Structures Provided

#### Vector

This implements the same data structure used by Clojure for its Vectors. The
underlying data structure is a bit-partitioned Hash Array Mapped Trie as
described by [Phil Bagwell][1].

The data structure remains immutable during appends, pops and updates because
the tree-based nature of the data structure allows paths to stored values to be
copied, without copying the entire tree.

``` go
import (
	"fmt"
	"github.com/d11wtq/persistent/vector"
)

// initialize a new vector with some values
vec1 := vector.New().Append(42).Append(7)

// read the number of elements
fmt.Println("vec1 has %d elements", vec1.Count())

// read the value at key 0
x, err := vec1.Get(0)
if err != nil {
	panic(err)
}
fmt.Println("vec1[0] == %s", x) // 42

// set key 0 to a new value
vec2, err := vec1.Set(0, 99)
if err != nil {
	panic(err)
}

// read the new value at key 0
y, err := vec2.Get(0)
if err != nil {
	panic(err)
}
fmt.Println("vec2[0] == %s", y) // 99

// verify that the original vector is unchanged
z, err := vec1.Get(0)
if err != nil {
	panic(err)
}
fmt.Println("vec1[0] == %s", z) // 42

// pop the last element off the env of vec1
vec3 := vec1.Pop()
fmt.Println("vec3 has %d elements", vec3.Count())

// verify that the original vector is unchanged
fmt.Println("vec1 has %d elements", vec1.Count())
```

I plan to experiment with changing the underlying data structure to an RRB Tree
as described by [Bagwell and Rompf][2], since there are some apparent benefits
for join/split operations.

  [1]: http://lampwww.epfl.ch/papers/idealhashtrees.pdf
  [2]: http://infoscience.epfl.ch/record/169879/files/RMTrees.pdf
