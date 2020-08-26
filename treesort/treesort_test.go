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
	a = a.Add(10)
	a = a.Add(9)
	a = a.Add(1)

	t.Logf("t=%s", a)
}
