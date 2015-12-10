package vector

const (
	// The number of bits to read for each sub key
	PARTITION_BITS = 5
	// The bits we're interested in for each sub key index
	PARTITION_MASK = 1<<PARTITION_BITS - 1
)

// What storable values
type Value interface{}

// Representation of a persistent vector node
type Vector struct {
	// The elements stored in this node
	Elements []Value
	// The maximum number of elements
	Shift uint32
}

// Sentinel value for the empty vector
var EmptyVector = &Vector{
	Elements: []Value{},
	Shift:    0,
}

// Return a new empty vector.
func New() *Vector {
	return EmptyVector
}

// Get the number of elements in this node.
func (vec *Vector) Width() uint32 {
	return uint32(len(vec.Elements))
}

// Find the element at a given key in the vector.
// If the key does not exist, (nil, false) is returned.
func (vec *Vector) Find(key uint32) (Value, bool) {
	for vec.Shift > 0 {
		vec = vec.Elements[((key >> vec.Shift) & PARTITION_MASK)].(*Vector)
	}

	idx := (key & PARTITION_MASK)
	if vec.Width() > idx {
		return vec.Elements[idx], true
	}

	return nil, false
}

// Set key in vector to value, returning a new vector.
func (vec *Vector) Set(key uint32, value Value) (cpy *Vector, ok bool) {
	cpy = vec.NewRoot(key)
	vec = cpy

	for vec.Shift > 0 {
		vec = vec.CopySubKey((key >> vec.Shift) & PARTITION_MASK)
	}

	if vec.SetSubKey((key & PARTITION_MASK), value) {
		return cpy, true
	}

	return nil, false
}

// Make a shallow copy of this node.
// This copies the node and its internal slice, but not its branches or values.
func (vec *Vector) Copy() *Vector {
	return &Vector{
		Elements: append([]Value(nil), vec.Elements...),
		Shift:    vec.Shift,
	}
}

// Return a copy of the root, or a new root if key overflows this root.
// A new root has an increased shift size.
func (vec *Vector) NewRoot(key uint32) *Vector {
	if (1 << vec.Shift) < (key >> PARTITION_BITS) {
		return vec.Copy()
	} else {
		return &Vector{
			Shift:    vec.Shift + PARTITION_BITS,
			Elements: []Value{vec},
		}
	}
}

// Set the direct subkey in vec to a copy of itself and return the copy.
// If the subkey is effectively an append, generate a new node.
// Mutates, on the assumption that vec is a copy.
func (vec *Vector) CopySubKey(key uint32) (cpy *Vector) {
	switch key {
	case vec.Width():
		cpy = &Vector{
			Shift:    (vec.Shift - PARTITION_BITS),
			Elements: make([]Value, 0),
		}
		vec.Elements = append(vec.Elements, cpy)
	default:
		cpy = vec.Elements[key].(*Vector).Copy()
		vec.Elements[key] = cpy
	}

	return
}

// Set the direct subkey in vec to a new value.
// If the subkey is effectively an append, generate a new slot.
// Mutates, on the assumption that vec is a copy.
func (vec *Vector) SetSubKey(key uint32, value Value) bool {
	switch {
	case vec.Width() == key:
		vec.Elements = append(vec.Elements, value)
	case vec.Width() > key:
		vec.Elements[key] = value
	default:
		return false
	}

	return true
}
