package order

import (
	"context"
	"sync"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/repository"

	"github.com/quolpr/booking/internal/app/transaction"
	"github.com/quolpr/booking/pkg/randomstr"
)

var _ repository.OrdersRepo = (*Repo)(nil)

type Repo struct {
	orders map[string]model.Order
	lock   sync.RWMutex
}

func NewRepo() *Repo {
	return &Repo{
		orders: make(map[string]model.Order),
	}
}

func (r *Repo) Create(ctx context.Context, order model.Order) (model.Order, error) {
	defer r.lock.Unlock()
	r.lock.Lock()

	order.ID = randomstr.Generate(20)

	r.orders[order.ID] = order

	transaction.RegisterTrCompensation(ctx, func(ctx context.Context) error {
		return r.RemoveOrder(ctx, order.ID)
	})

	return order, nil
}

func (r *Repo) GetAll(ctx context.Context) ([]model.Order, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	values := make([]model.Order, 0, len(r.orders))

	for _, order := range r.orders {
		values = append(values, order)
	}

	return values, nil
}

func (r *Repo) Get(ctx context.Context, id string) (model.Order, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, order := range r.orders {
		if order.ID == id {
			return order, nil
		}
	}

	return model.Order{}, repository.ErrNotFound
}

func (r *Repo) RemoveOrder(ctx context.Context, id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.orders[id]

	if ok {
		delete(r.orders, id)

		return nil
	}

	return repository.ErrNotFound
}
