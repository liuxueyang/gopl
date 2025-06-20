package treesort

import "testing"

func TestSort(t *testing.T) {
	a := []int{
		3, 2, 1, 5, 4,
	}
	Sort(a)
	t.Logf("a=%#v", a)
}

func TestTreeString(t *testing.T) {
	var a *tree
	a = Add(a, 10)
	a = Add(a, 9)
	a = Add(a, 1)

	t.Logf("t=%s", a)
}
