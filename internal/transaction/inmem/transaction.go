package inmem

import (
	"context"

	"github.com/quolpr/booking/internal/transaction"
)

var _ transaction.Transaction = &Transaction{}

type Transaction struct {
	compensations []func(ctx context.Context) error
}

func (t *Transaction) Commit(_ context.Context) error {
	return nil
}

func (t *Transaction) Rollback(ctx context.Context) error {
	for i := len(t.compensations) - 1; i >= 0; i-- {
		err := t.compensations[i](ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Transaction) RegisterCompensation(ctx context.Context, fn func(ctx context.Context) error) {
	t.compensations = append(t.compensations, fn)
}

type TransactionManager struct{}

func (c *TransactionManager) CreateTr() transaction.Transaction {
	return &Transaction{}
}
