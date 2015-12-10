package vector

const (
	// The number of bits to read for each sub key
	PARTITION_BITS = 5
	// The bits we're interested in for each sub key index
	PARTITION_MASK = 1<<PARTITION_BITS - 1
)

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
var empty = &Vector{
	Root:   &Node{Elements: []Value{}, Shift: 0},
	Length: 0,
}

// Return a new empty vector.
func New(elements ...Value) *Vector {
	acc := empty
	// FIXME: Bulk load in chunks of 32 elements
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
		// this should never happen
		panic(err)
	}

	return vec
}