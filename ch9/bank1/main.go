package main

import (
	"ch9/bank1/bank"
	"fmt"
	"sync"
)

func main() {
	bank.Deposit(200)
	balance := bank.Balance()
	println("Balance:", balance)

	bank.Deposit(100)
	balance = bank.Balance()
	println("Balance:", balance)

	bank.Deposit(50)
	balance = bank.Balance()
	println("Balance:", balance)

	// Simulate concurrent deposits
	// This will demonstrate that the teller goroutine can handle multiple deposits.
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(amount int) {
			defer wg.Done()
			bank.Deposit(amount * 10)
			fmt.Println("Deposited:", amount*10)
		}(i + 1)
	}
	wg.Wait()

	balance = bank.Balance()
	println("Final Balance:", balance)
}
