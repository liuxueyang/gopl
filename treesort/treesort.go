package treesort

import (
	"strconv"
)

type tree struct {
	value       int
	left, right *tree
}

func Add(t *tree, value int) *tree {
	if t == nil {
		t1 := new(tree)
		t1.value = value
		return t1
	}
	if value < t.value {
		t.left = Add(t.left, value)
	} else {
		t.right = Add(t.right, value)
	}
	return t
}

func appendValues(values []int, t *tree) []int {
	if t == nil {
		return values
	}
	values = appendValues(values, t.left)
	values = append(values, t.value)
	values = appendValues(values, t.right)
	return values
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = Add(root, v)
	}
	appendValues(values[:0], root)
}

func (t *tree) String() string {
	if t == nil {
		return ""
	}
	left := t.left.String()
	right := t.right.String()

	var res string
	if len(left) > 0 {
		res = left + ", "
	}
	res = res + strconv.FormatInt(int64(t.value), 10)

	if len(right) > 0 {
		res = res + ", " + right
	}

	return res
}
