package main

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go counter(naturals)
	go square(naturals, squares)
	printer(squares)
}

func counter(out chan<- int) {
	for i := range 10 {
		out <- i
	}
	close(out)
}

func square(in <-chan int, out chan<- int) {
	for i := range in {
		out <- i * i
	}
	close(out)
}

func printer(in <-chan int) {
	for i := range in {
		println(i)
	}
}
