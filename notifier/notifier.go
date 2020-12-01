package notifier

type UsdtTransaction struct {
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TransactionHash string `json:"transactionHash"`
	// address which sent tokens
	From string `json:"from"`
	// address which received tokens
	To string `json:"to"`
	// amount of transferred tokens
	Amount uint64 `json:"amount"`
	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}

type Notifier interface {
	Notify(t UsdtTransaction)
}
