package bank1

type message struct {
	amount int
	ok     chan bool
}

var deposits = make(chan int)
var balances = make(chan int)
var withdrawals = make(chan message)

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
		case msg := <-withdrawals:
			var ok bool

			if msg.amount > 0 && msg.amount <= balance {
				// If the withdrawal is valid, it will be processed and the balance updated.
				ok = true
				balance -= msg.amount
			}
			msg.ok <- ok
		}
	}
}

func Withdraw(amount int) bool {
	// This function sends a withdrawal request and waits for the response.
	msg := message{amount: amount, ok: make(chan bool)}
	withdrawals <- msg
	return <-msg.ok
}

func init() {
	go teller()
}
