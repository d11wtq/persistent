package vector

import (
	"testing"
)

func AssertContains(t *testing.T, vec *Vector, elems map[uint32]Value) {
	for k, v := range elems {
		x, err := vec.Get(k)
		if err != nil {
			t.Fatalf(`expected vec.Get(%d) to be ok, got %s`, k, err)
		}
		if x != v {
			t.Fatalf(`expected vec.Get(%d) == %s, got %s`, k, v, x)
		}
	}
}

func TestGet1Deep(t *testing.T) {
	vec := &Vector{
		Root: &Node{
			Elements: []Value{42, 21, 17},
			Shift:    0, // 5 * (1 - 1)
		},
		Length: 3,
	}

	x, err := vec.Get(0)
	if err != nil {
		t.Fatalf(`expected vec.Get(0) to be ok, got %s`, err)
	}
	if x != 42 {
		t.Fatalf(`expected vec.Get(0) == 42, got %s`, x)
	}

	y, err := vec.Get(1)
	if err != nil {
		t.Fatalf(`expected vec.Get(1) to be ok, got %s`, err)
	}
	if y != 21 {
		t.Fatalf(`expected vec.Get(1) == 21, got %s`, y)
	}

	z, err := vec.Get(2)
	if err != nil {
		t.Fatalf(`expected vec.Get(2) to be ok, got %s`, err)
	}
	if z != 17 {
		t.Fatalf(`expected vec.Get(2) == 17, got %s`, z)
	}

	_, err = vec.Get(3)
	if err == nil {
		t.Fatalf(`expected vec.Get(3) not to be ok, but was`)
	}
}

func TestGet2Deep(t *testing.T) {
	vec := &Vector{
		Root: &Node{
			Elements: []Value{
				&Node{Elements: []Value{42, 21, 17}},
			},
			Shift: 5, // 5 * (2 - 1)
		},
		Length: 3,
	}

	x, err := vec.Get(0)
	if err != nil {
		t.Fatalf(`expected vec.Get(0) to be ok, got %s`, err)
	}
	if x != 42 {
		t.Fatalf(`expected vec.Get(0) == 42, got %s`, x)
	}

	y, err := vec.Get(1)
	if err != nil {
		t.Fatalf(`expected vec.Get(1) to be ok, got %s`, err)
	}
	if y != 21 {
		t.Fatalf(`expected vec.Get(1) == 21, got %s`, y)
	}

	z, err := vec.Get(2)
	if err != nil {
		t.Fatalf(`expected vec.Get(2) to be ok, got %s`, err)
	}
	if z != 17 {
		t.Fatalf(`expected vec.Get(2) == 17, got %s`, z)
	}

	_, err = vec.Get(3)
	if err == nil {
		t.Fatalf(`expected vec.Get(3) not to be ok, but was`)
	}
}

func TestUpdateViaSet1Deep(t *testing.T) {
	vec := &Vector{
		Root: &Node{
			Elements: []Value{42, 21, 17},
			Shift:    0, // 5 * (1 - 1)
		},
		Length: 3,
	}

	cpy, err := vec.Set(1, 57)
	if err != nil {
		t.Fatalf(`expected vec.Set(1, ...) to be ok, got %s`, err)
	}

	AssertContains(
		t, cpy,
		map[uint32]Value{
			0: 42,
			1: 57,
			2: 17,
		},
	)

	AssertContains(
		t, vec,
		map[uint32]Value{
			0: 42,
			1: 21,
			2: 17,
		},
	)

	_, err = cpy.Get(3)
	if err == nil {
		t.Fatalf(`expected cpy.Get(3) not to be ok, but was`)
	}
}

func TestUpdateViaSet2Deep(t *testing.T) {
	vec := &Vector{
		Root: &Node{
			Elements: []Value{
				&Node{Elements: []Value{42, 21, 17}},
			},
			Shift: 5, // 5 * (2 - 1)
		},
		Length: 3,
	}

	cpy, err := vec.Set(1, 57)
	if err != nil {
		t.Fatalf(`expected vec.Set(1, ...) to be ok, got %s`, err)
	}

	AssertContains(
		t, cpy,
		map[uint32]Value{
			0: 42,
			1: 57,
			2: 17,
		},
	)

	AssertContains(
		t, vec,
		map[uint32]Value{
			0: 42,
			1: 21,
			2: 17,
		},
	)

	_, err = cpy.Get(3)
	if err == nil {
		t.Fatalf(`expected cpy.Get(3) not to be ok, but was`)
	}
}

func TestAppendViaSet1Deep(t *testing.T) {
	vec := &Vector{
		Root: &Node{
			Elements: []Value{42, 21, 17},
			Shift:    0, // 5 * (1 - 1)
		},
		Length: 3,
	}

	cpy, err := vec.Set(3, 57)
	if err != nil {
		t.Fatalf(`expected vec.Set(3, ...) to be ok, got %s`, err)
	}

	AssertContains(
		t, cpy,
		map[uint32]Value{
			0: 42,
			1: 21,
			2: 17,
			3: 57,
		},
	)

	_, err = cpy.Get(4)
	if err == nil {
		t.Fatalf(`expected cpy.Get(4) not to be ok, but was`)
	}

	_, err = vec.Get(3)
	if err == nil {
		t.Fatalf(`expected vec.Get(3) not to be ok, but was`)
	}
}

func TestAppendViaSet2Deep(t *testing.T) {
	vec := &Vector{
		Root: &Node{
			Elements: []Value{
				&Node{Elements: []Value{42, 21, 17}},
			},
			Shift: 5, // 5 * (2 - 1)
		},
		Length: 3,
	}

	cpy, err := vec.Set(3, 57)
	if err != nil {
		t.Fatalf(`expected vec.Set(3, ...) to be ok, got %s`, err)
	}

	AssertContains(
		t, cpy,
		map[uint32]Value{
			0: 42,
			1: 21,
			2: 17,
			3: 57,
		},
	)

	_, err = cpy.Get(4)
	if err == nil {
		t.Fatalf(`expected cpy.Get(4) not to be ok, but was`)
	}

	_, err = vec.Get(3)
	if err == nil {
		t.Fatalf(`expected vec.Get(3) not to be ok, but was`)
	}
}

func TestAppendOverflow1Deep(t *testing.T) {
	elems := make([]Value, 0, 32)
	for i := 0; i < 32; i += 1 {
		elems = append(elems, i)
	}
	vec := &Vector{
		Root:   &Node{Elements: elems, Shift: 0}, // 5 * (1 - 1)
		Length: 32,
	}

	cpy, err := vec.Set(32, 32)
	if err != nil {
		t.Fatalf(`expected vec.Set(32, ...) to be ok, got %s`, err)
	}

	AssertContains(
		t, cpy,
		map[uint32]Value{
			0:  0,
			31: 31,
			32: 32,
		},
	)

	_, err = cpy.Get(33)
	if err == nil {
		t.Fatalf(`expected cpy.Get(33) not to be ok, but was`)
	}
}

func TestAppendOverflow2Deep(t *testing.T) {
	nodes := make([]Value, 0, 32)
	for i := 0; i < 32; i += 1 {
		elems := make([]Value, 0, 32)
		for j := 0; j < 32; j += 1 {
			elems = append(elems, i*32+j)
		}
		nodes = append(nodes, &Node{Elements: elems})
	}

	vec := &Vector{
		Root:   &Node{Elements: nodes, Shift: 5}, // 5 * (2 - 1)
		Length: 1024,
	}

	cpy, err := vec.Set(1024, 1024)
	if err != nil {
		t.Fatalf(`expected vec.Set(1024, ...) to be ok, got %s`, err)
	}

	AssertContains(
		t, cpy,
		map[uint32]Value{
			0:    0,
			31:   31,
			32:   32,
			1023: 1023,
			1024: 1024,
		},
	)

	_, err = cpy.Get(1025)
	if err == nil {
		t.Fatalf(`expected cpy.Get(1025) not to be ok, but was`)
	}
}

func TestSetOutOfBounds1Deep(t *testing.T) {
	vec := &Vector{
		Root: &Node{
			Elements: []Value{42, 21, 17},
			Shift:    0, // 5 * (1 - 1)
		},
		Length: 3,
	}

	_, err := vec.Set(4, 57)
	if err == nil {
		t.Fatalf(`expected vec.Set(4, ...) not to be ok, but was`)
	}
}

func TestCount(t *testing.T) {
	vec := &Vector{
		Root:   &Node{Elements: []Value{}},
		Length: 0,
	}

	if vec.Count() != 0 {
		t.Fatalf(`expected vec.Count() == 0, got %s`, vec.Count())
	}

	vec, _ = vec.Set(0, 42)
	if vec.Count() != 1 {
		t.Fatalf(`expected vec.Count() == 1, got %s`, vec.Count())
	}

	vec, _ = vec.Set(0, 21)
	if vec.Count() != 1 {
		t.Fatalf(`expected vec.Count() == 1, got %s`, vec.Count())
	}

	vec, _ = vec.Set(1, 15)
	if vec.Count() != 2 {
		t.Fatalf(`expected vec.Count() == 2, got %s`, vec.Count())
	}
}

func TestAppend(t *testing.T) {
	vec := &Vector{
		Root:   &Node{Elements: []Value{}},
		Length: 0,
	}
	vec = vec.Append(42)
	vec = vec.Append(21)
	vec = vec.Append(17)

	AssertContains(
		t, vec,
		map[uint32]Value{
			0: 42,
			1: 21,
			2: 17,
		},
	)

	if vec.Count() != 3 {
		t.Fatalf(`expected vec.Count() == 3, got %s`, vec.Count())
	}
}

func TestNewWithoutArgs(t *testing.T) {
	vec := New()
	if vec.Count() != 0 {
		t.Fatalf(`expected vec.Count() == 0, got %s`, vec.Count())
	}
}