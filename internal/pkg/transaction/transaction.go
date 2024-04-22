package transaction

import (
	"context"
	"errors"
	"fmt"
)

type ctxKeys string

const (
	transactionKey        ctxKeys = "transaction"
	transactionManagerKey ctxKeys = "transaction_manager"
)

var ErrRollback = errors.New("rollback")

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	RegisterCompensation(ctx context.Context, fn func(ctx context.Context) error)
}

type TransactionManager interface {
	CreateTr() Transaction
}

func WithTransactionManager(ctx context.Context, manager TransactionManager) context.Context {
	return context.WithValue(ctx, transactionManagerKey, manager)
}

func InTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	if ctx.Value(transactionKey) != nil {
		return f(ctx)
	}

	manager, ok := ctx.Value(transactionManagerKey).(TransactionManager)

	if !ok {
		panic("Transaction manager not found in context")
	}

	tr := manager.CreateTr()

	ctx = context.WithValue(ctx, transactionKey, tr)

	err := f(ctx)
	if err != nil {
		rollbackErr := tr.Rollback(ctx)
		if rollbackErr != nil {
			err = fmt.Errorf("%w: %w", ErrRollback, rollbackErr)
		}

		return err
	}

	return tr.Commit(ctx)
}

func RegisterTrCompensation(ctx context.Context, fn func(ctx context.Context) error) {
	val := ctx.Value(transactionKey)

	if val == nil {
		return
	}

	tr, ok := val.(Transaction)

	if !ok {
		// It's ok if we are not in transaction
		return
	}

	tr.RegisterCompensation(ctx, fn)
}
