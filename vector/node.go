package vector

const (
	// The number of bits to read for each sub key
	BITS = 5
	// The bits we're interested in for each sub key index
	MASK = 1<<BITS - 1
)

// Representation of a persistent vector node.
// Boundary checks are not performed, as it is assumed the consumer is aware of
// the length of the vector.
type Node struct {
	// The elements stored in this node
	Elements []Value
	// The maximum number of elements
	Shift uint32
}

// Create a new empty root node.
func EmptyNode() *Node {
	return &Node{Elements: []Value{}, Shift: 0}
}

// Find the element at a given key starting from this node.
func (node *Node) Get(key uint32) Value {
	for node.Shift > 0 {
		node = node.Elements[((key >> node.Shift) & MASK)].(*Node)
	}

	return node.Elements[(key & MASK)]
}

// Set key in vector to value, returning a new root node.
// Attempting to set a key beyond the current length is an OutOfBounds error.
func (node *Node) Set(key uint32, value Value) (into *Node) {
	into = node.NewRoot(key)
	node = into

	for node.Shift > 0 {
		node = node.CopySubKey((key >> node.Shift) & MASK)
	}

	node.SetSubKey((key & MASK), value)

	return
}

// Truncate the length of this node (and its children).
// This discards all branches to the right of the length.
func (node *Node) Truncate(length uint32) (into *Node) {
	if length == 0 {
		return EmptyNode()
	}

	var (
		key uint32 = length
		idx uint32
	)

	into = node.Copy()
	node = into

	for node.Shift > 0 {
		idx = (key >> node.Shift) & MASK
		node.Elements = append([]Value(nil), node.Elements[:idx+1]...)
		node = node.CopySubKey(idx)
	}

	node.Elements = append([]Value(nil), node.Elements[:(key&MASK)]...)

	return into.Flatten()
}

// Get the number of elements in this node.
func (node *Node) Width() uint32 {
	return uint32(len(node.Elements))
}

// Make a shallow copy of this node.
// This copies the node and its internal slice, but not its branches or values.
func (node *Node) Copy() *Node {
	return &Node{
		Elements: append([]Value(nil), node.Elements...),
		Shift:    node.Shift,
	}
}

// Return a copy of the root, or a new root if key overflows this root.
// A new root has an increased shift size.
func (node *Node) NewRoot(key uint32) *Node {
	if (1 << node.Shift) < (key >> BITS) {
		return node.Copy()
	} else {
		return &Node{
			Shift:    node.Shift + BITS,
			Elements: []Value{node},
		}
	}
}

// Set the direct subkey in node to a copy of itself and return the copy.
// If the subkey is effectively an append, generate a new node.
// Mutates, on the assumption that node is a copy.
func (node *Node) CopySubKey(key uint32) (into *Node) {
	switch key {
	case node.Width():
		into = &Node{
			Shift:    (node.Shift - BITS),
			Elements: make([]Value, 0),
		}
		node.Elements = append(node.Elements, into)
	default:
		into = node.Elements[key].(*Node).Copy()
		node.Elements[key] = into
	}

	return
}

// Set the direct subkey in this node to a new value.
// If the subkey is effectively an append, generate a new slot.
// Mutates, on the assumption that node is a copy.
func (node *Node) SetSubKey(key uint32, value Value) {
	switch key {
	case node.Width():
		node.Elements = append(node.Elements, value)
	default:
		node.Elements[key] = value
	}
}

// Kill any redundant root nodes and return the first real root.
func (node *Node) Flatten() *Node {
	for node.Shift > 0 && node.Width() == 1 {
		node = node.Elements[0].(*Node)
	}
	return node
}
