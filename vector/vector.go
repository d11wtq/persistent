package vector

// Values storable in the vector
type Value interface{}

// Pointer to the root node and its length
type Vector struct {
	// The root node of the vector
	Root *Node
	// The number of elements in the vector
	Length uint32
	// The key at which the vector starts
	Offset uint32
}

// Value for the empty vector
var empty = &Vector{
	Root:   EmptyNode(),
	Length: 0,
	Offset: 0,
}

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
	if vec.Length > key {
		return vec.Root.Get(vec.Offset + key), nil
	}

	return nil, &OutOfBounds{key}
}

// Set a given key in the vector.
// Allowed indices are those already set, and that in the append position.
// Attempts to set key > length is an OutOfBounds error.
// A new vector is returned, sharing memory with the original.
func (vec *Vector) Set(key uint32, value Value) (*Vector, error) {
	if key > vec.Length {
		return nil, &OutOfBounds{key}
	}

	newLength := vec.Length
	if key == newLength {
		newLength += 1
	}

	return &Vector{
		Root:   vec.Root.Set(key+vec.Offset, value),
		Length: newLength,
		Offset: vec.Offset,
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

// Prepend a value to the start of this vector.
// A new vector is returned, sharing memory with the original.
func (vec *Vector) Prepend(value Value) *Vector {
	if vec.Offset == 0 {
		// FIXME: Allocate more space to the left in O(1) time
		return vec
	}

	return &Vector{
		Root:   vec.Root.Set((vec.Offset - 1), value),
		Length: vec.Length + 1,
		Offset: vec.Offset - 1,
	}
}

// Return the vector with all elements > length removed.
// A new vector is returned, sharing memory with the original.
// Attempting to truncate to a length > the current length returns itself.
func (vec *Vector) Truncate(length uint32) *Vector {
	if length < vec.Length {
		return &Vector{
			Root:   vec.Root.Truncate(vec.Offset + length),
			Length: length,
			Offset: vec.Offset,
		}
	}

	return vec
}

// Return the vector with all elements < length removed.
// A new vector is returned, sharing memory with the original.
// Attempting to drop length > the current length returns the empty vector.
func (vec *Vector) Drop(length uint32) *Vector {
	if length >= vec.Length {
		return empty
	} else if length > 0 {
		return &Vector{
			Root:   vec.Root.EraseTo(vec.Offset + length),
			Length: vec.Length - length,
			Offset: vec.Offset + length,
		}
	}

	return vec
}

// Return the vector with the last element removed.
// A new vector is returned, sharing memory with the original.
// Attempting to pop an empty vector returns itself.
func (vec *Vector) Pop() *Vector {
	if vec.Length == 0 {
		return vec
	}

	return vec.Truncate(vec.Length - 1)
}

// Return the vector with the first element removed.
// A new vector is returned, sharing memory with the original.
// Attempting to shift an empty vector returns itself.
func (vec *Vector) Shift() *Vector {
	if vec.Length == 0 {
		return vec
	}

	return vec.Drop(1)
}
