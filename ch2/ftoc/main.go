package main

import "fmt"

const freezingF = 32.0
const boilingF = 212.0

func main() {
	fmt.Printf("%g F = %g C\n", freezingF, fToC(freezingF))
	fmt.Printf("%g F = %g C\n", boilingF, fToC(boilingF))
}

func fToC(f float64) float64 {
	return (f - freezingF) * 5 / 9
}
