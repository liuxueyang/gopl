package bank3

import "sync"

var (
	mu      sync.Mutex
	balance int
)

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	balance += amount
}

func Balance() int {
	mu.Lock()
	defer mu.Unlock()
	b := balance
	return b
}
