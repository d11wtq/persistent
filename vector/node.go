package vector

const (
	// The number of bits to read for each sub key
	BITS = 5
	// The bits we're interested in for each sub key index
	MASK = 1<<BITS - 1
)

// Representation of a persistent vector node
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
// If the key does not exist, it is an OutOfBounds error.
func (node *Node) Get(key uint32) (val Value, err error) {
	for node.Shift > 0 {
		val, err = node.ReadSubKey((key >> node.Shift) & MASK)
		if err != nil {
			return nil, &OutOfBounds{key}
		}

		node = val.(*Node)
	}

	val, err = node.ReadSubKey(key & MASK)
	if err != nil {
		return nil, &OutOfBounds{key}
	}

	return
}

// Set key in vector to value, returning a new root node.
// Attempting to set a key beyond the current length is an OutOfBounds error.
func (node *Node) Set(key uint32, value Value) (into *Node, err error) {
	into = node.NewRoot(key)
	node = into

	for node.Shift > 0 {
		node = node.CopySubKey((key >> node.Shift) & MASK)
	}

	if node.SetSubKey((key & MASK), value) {
		return into, nil
	}

	return nil, &OutOfBounds{key}
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

	for into.Shift > 0 && into.Width() == 1 {
		into = into.Elements[0].(*Node)
	}

	return
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

// Get the direct subkey in node, or return OutOfBounds if does not exist.
func (node *Node) ReadSubKey(key uint32) (Value, error) {
	if node.Width() > key {
		return node.Elements[key], nil
	}

	return nil, &OutOfBounds{key}
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
func (node *Node) SetSubKey(key uint32, value Value) bool {
	switch {
	case node.Width() == key:
		node.Elements = append(node.Elements, value)
	case node.Width() > key:
		node.Elements[key] = value
	default:
		return false
	}

	return true
}
