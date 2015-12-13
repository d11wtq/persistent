package vector

// Values storable in the vector
type Value interface{}

// Pointer to the root node and its length
type Vector struct {
	// The root node of the vector
	Root *Node
	// The number of elements in the vector
	Length uint32
}

// Value for the empty vector
var empty = &Vector{Root: EmptyNode(), Length: 0}

// Return a new empty vector.
func New(elements ...Value) *Vector {
	acc := empty
	for _, v := range elements {
		acc = acc.Append(v)
	}
	return acc
}

// Return the number of elements in this vector.
func (vec *Vector) Count() uint32 {
	return vec.Length
}

// Get the value for a given key in the vector.
// Access to a key that is not in the vector is an OutOfBounds error.
func (vec *Vector) Get(key uint32) (Value, error) {
	return vec.Root.Get(key)
}

// Set a given key in the vector.
// Allowed indices are those already set, and that in the append position.
// Attempts to set key > length is an OutOfBounds error.
// A new vector is returned, sharing memory with the original.
func (vec *Vector) Set(key uint32, value Value) (*Vector, error) {
	if key > vec.Length {
		return nil, &OutOfBounds{key}
	}

	newRoot, err := vec.Root.Set(key, value)
	if err != nil {
		return nil, err
	}

	newLength := vec.Length
	if key == newLength {
		newLength += 1
	}

	return &Vector{
		Root:   newRoot,
		Length: newLength,
	}, nil
}

// Append a value to the end of this vector.
// A new vector is returned, sharing memory with the original.
func (vec *Vector) Append(value Value) *Vector {
	vec, err := vec.Set(vec.Length, value)
	if err != nil {
		panic(err)
	}

	return vec
}

// Return the vector with the last element removed.
// A new vector is returned, sharing memory with the original.
func (vec *Vector) Pop() *Vector {
	if vec.Length == 0 {
		return vec
	}

	newLength := vec.Length - 1

	return &Vector{
		Root:   vec.Root.Truncate(newLength),
		Length: newLength,
	}
}

// Return the sub-vector that lies between indices start and end.
// Attempting to slice outside of the vector is an OutOfBounds error.
// A new vector is returned, sharing memory with the original.
func (vec *Vector) Slice(start, end uint32) (into *Vector, err error) {
	return vec, nil
}
