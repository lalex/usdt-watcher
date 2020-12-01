package notifier

import (
	"fmt"
)

type ConsoleNotifier struct{}

func NewConsoleNotifier() Notifier {
	return &ConsoleNotifier{}
}

func (n *ConsoleNotifier) Notify(t UsdtTransaction) {
	fmt.Println(fmt.Sprintf("[%s] Received %d USDT from %s", t.TransactionHash, t.Amount, t.From))
}
