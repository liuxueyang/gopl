package main

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for i := range 10 {
			naturals <- i
		}
		close(naturals)
	}()

	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	for x := range squares {
		println(x)
	}
}
