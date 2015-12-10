package vector

import (
	"fmt"
)

// Error type returned when accessing an invalid index
type OutOfBounds struct {
	Key uint32
}

func (e *OutOfBounds) Error() string {
	return fmt.Sprintf("key %d out of bounds", e.Key)
}
