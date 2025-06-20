package popcount

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

// Exercise 2.3
func PopCountLoop(x uint64) int {
	var res int = 0

	for i := range 64 {
		if (x>>i)&1 == 1 {
			res++
		}
	}

	return res
}

// Exercise 2.4
func PopCountShift(x uint64) int {
	var res int = 0

	for x != 0 {
		if x&1 == 1 {
			res++
		}
		x >>= 1
	}
	return res
}

// Exercise 2.5
func PopCountClear(x uint64) int {
	var res int = 0

	for x != 0 {
		res++
		x = x & (x - 1)
	}

	return res
}
