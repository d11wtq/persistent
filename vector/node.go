package vector

// Representation of a persistent vector node
type Node struct {
	// The elements stored in this node
	Elements []Value
	// The maximum number of elements
	Shift uint32
}

// Find the element at a given key starting from this node.
// If the key does not exist, it is an OutOfBounds error.
func (node *Node) Get(key uint32) (Value, error) {
	for node.Shift > 0 {
		node = node.Elements[((key >> node.Shift) & PARTITION_MASK)].(*Node)
	}

	idx := (key & PARTITION_MASK)
	if node.Width() > idx {
		return node.Elements[idx], nil
	}

	return nil, &OutOfBounds{key}
}

// Set key in vector to value, returning a new root node.
// Attempting to set a key beyond the current length is an OutOfBounds error.
func (node *Node) Set(key uint32, value Value) (into *Node, err error) {
	into = node.NewRoot(key)
	node = into

	for node.Shift > 0 {
		node = node.CopySubKey((key >> node.Shift) & PARTITION_MASK)
	}

	if node.SetSubKey((key & PARTITION_MASK), value) {
		return into, nil
	}

	return nil, &OutOfBounds{key}
}

func (node *Node) Pop() (into *Node) {
	if node.Width() == 0 {
		return node
	}

	into = node.Copy()
	node = into

	for node.Shift > 0 {
		node = node.CopySubKey(node.Width() - 1)
	}

	node.Elements = node.Elements[:node.Width()-1]

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
	if (1 << node.Shift) < (key >> PARTITION_BITS) {
		return node.Copy()
	} else {
		return &Node{
			Shift:    node.Shift + PARTITION_BITS,
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
			Shift:    (node.Shift - PARTITION_BITS),
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
