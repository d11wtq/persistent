package vector

const (
	// The number of bits to read for each sub key
	BITS = 5
	// The bits we're interested in for each sub key index
	MASK = 1<<BITS - 1
	// The size of each node
	SIZE = 1 << BITS
)

// Sentinel type unset values
type nullSentinel struct{}

// Sentinel for unset values
var Null = &nullSentinel{}

// Representation of a persistent vector node.
// Boundary checks are not performed, as it is assumed the consumer is aware of
// the length of the vector.
type Node struct {
	// The elements stored in this node
	Elements []Value
	// The number of bits to shift off at this level
	Shift uint32
}

// Create a new empty root node.
// Complexity: O(1)
func EmptyNode() *Node {
	return NewNode(0)
}

// Fill elements up to the expected capacity.
// Complexity: O(1)
func Fill(elements ...Value) []Value {
	elements = append(make([]Value, 0, SIZE), elements...)
	for len(elements) < cap(elements) {
		elements = append(elements, Null)
	}
	return elements
}

// Create a new node at shift depth, filled with elements.
// Complexity: O(1)
func NewNode(shift uint32, elements ...Value) *Node {
	return &Node{
		Elements: Fill(elements...),
		Shift:    shift,
	}
}

// Find the element at a given key starting from this node.
// Complexity: O(log(n))
// Effectively: O(1)
func (node *Node) Get(key uint32) Value {
	for node.Shift > 0 {
		node = node.Elements[((key >> node.Shift) & MASK)].(*Node)
	}

	return node.Elements[(key & MASK)]
}

// Set key in vector to value, returning a new root node.
// Attempting to set a key beyond the current length is an OutOfBounds error.
// Complexity: O(log(n))
// Effectively: O(1)
func (node *Node) Set(key uint32, value Value) (into *Node) {
	into = node.NewRoot(key)
	node = into

	for node.Shift > 0 {
		node = node.CopySubKey((key >> node.Shift) & MASK)
	}

	node.Elements[(key & MASK)] = value

	return
}

// Truncate the length of this node (and its children).
// This discards all branches to the right of the length.
// Complexity: O(log(n))
// Effectively: O(1)
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
		for i := idx + 1; i < SIZE; i++ {
			node.Elements[i] = Null
		}
		node = node.CopySubKey(idx)
	}

	for i := (key & MASK); i < SIZE; i++ {
		node.Elements[i] = Null
	}

	// Root node with only one child
	for into.Shift > 0 && length < (1<<(into.Shift+BITS)) {
		into = into.Elements[0].(*Node)
	}

	return
}

// Erase memory at the start of this node.
// Access to elements where idx < length are invalid.
// Complexity: O(log(n))
// Effectively: O(1)
func (node *Node) EraseTo(length uint32) (into *Node) {
	if length == 0 {
		return node
	}

	var (
		key uint32 = length
		idx uint32
	)

	into = node.Copy()
	node = into

	for node.Shift > 0 {
		idx = (key >> node.Shift) & MASK

		for i := idx; i > 0; i-- {
			node.Elements[i-1] = Null
		}
		node = node.CopySubKey(idx)
	}

	for i := (key & MASK); i > 0; i-- {
		node.Elements[i-1] = Null
	}

	return
}

// Make a shallow copy of this node.
// This copies the node and its internal slice, but not its branches or values.
// Complexity: O(1)
func (node *Node) Copy() *Node {
	return NewNode(node.Shift, node.Elements...)
}

// Return a copy of the root, or a new root if key overflows this root.
// A new root has an increased shift size.
// Complexity: O(1)
func (node *Node) NewRoot(key uint32) *Node {
	if (1 << node.Shift) < (key >> BITS) {
		return node.Copy()
	} else {
		return NewNode(node.Shift+BITS, node)
	}
}

// Set the direct subkey in node to a copy of itself and return the copy.
// If the subkey is effectively an append, generate a new node.
// Mutates, on the assumption that node is a copy.
// Complexity: O(1)
func (node *Node) CopySubKey(key uint32) (into *Node) {
	if node.Elements[key] == Null {
		into = NewNode(node.Shift - BITS)
		node.Elements[key] = into
	} else {
		into = node.Elements[key].(*Node).Copy()
		node.Elements[key] = into
	}

	return
}

// Allocate space to the left of the current node.
// Returns a new node and the new offset of the existing data.
// Complexity: O(1)
func (node *Node) AllocLeft() (*Node, uint32) {
	var (
		into = NewNode(node.Shift + BITS)
		half = uint32((1 << (into.Shift + BITS)) / 2)
	)
	into.Elements[(half>>into.Shift)&MASK] = node
	return into, half
}
