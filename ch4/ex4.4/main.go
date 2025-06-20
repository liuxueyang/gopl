package main

func rotate_1pass(lst []int, n int) []int {
	ans := make([]int, len(lst))
	sz := len(lst)

	for i := range sz {
		ans[(i-n+sz)%sz] = lst[i]
	}
	return ans
}

func main() {
	lst := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	n := 3

	rotated_lst := rotate_1pass(lst, n)
	for _, v := range rotated_lst {
		print(v, " ")
	}
	println()
}
