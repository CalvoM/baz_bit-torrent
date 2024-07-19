package udp

import "math/rand"

type (
	TransactionID uint32
	Transaction   struct {
		id TransactionID
	}
)

func (t Transaction) ID() TransactionID {
	return t.id
}

func (t *Transaction) Refresh() {
	t.id = TransactionID(rand.Uint32())
}

func (t *Transaction) New() TransactionID {
	t.id = TransactionID(rand.Uint32())
	return t.id
}
