package bank2

import (
	"sync"
	"testing"
)

func TestDeposit(t *testing.T) {
	Deposit(100)
	if Balance() != 100 {
		t.Errorf("Expected balance to be 100, got %d", Balance())
	}

	Deposit(50)
	if Balance() != 150 {
		t.Errorf("Expected balance to be 150, got %d", Balance())
	}
}

func TestDeposit1(t *testing.T) {
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Deposit(i)
		}(i + 1)
	}
	wg.Wait()

	if Balance() != 5050 {
		t.Errorf("Expected balance to be 5050, got %d", Balance())
	}
}
