package bank

var deposits = make(chan int)
var balances = make(chan int)

func Deposit(amount int) {
	deposits <- amount
}

func Balance() int {
	return <-balances
}

func teller() {
	var balance int

	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
			// This case sends the current balance to the balances channel.
			// It will block until someone reads from the balances channel.
			// If no one is reading, it will wait until a read occurs.
			// This ensures that the balance is only sent when it can be received.
		}
	}
}

func init() {
	go teller()
}
