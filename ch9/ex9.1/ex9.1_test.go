package bank1

import "testing"

func TestWithdraw(t *testing.T) {
	Deposit(100)
	if !Withdraw(50) {
		t.Error("Expected withdrawal to succeed")
	}
	if Balance() != 50 {
		t.Errorf("Expected balance to be 50, got %d", Balance())
	}

	if Withdraw(60) {
		t.Error("Expected withdrawal to fail")
	}
	if Balance() != 50 {
		t.Errorf("Expected balance to remain 50, got %d", Balance())
	}

	if Withdraw(0) {
		t.Error("Expected withdrawal of zero to fail")
	}
	if Balance() != 50 {
		t.Errorf("Expected balance to remain 50, got %d", Balance())
	}
}
