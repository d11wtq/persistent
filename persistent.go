package persistent

import (
	"./vector"
)

// Return a new persistent vector with specified elements.
func Vector(elements ...vector.Value) *vector.Vector {
	return vector.New(elements...)
}
