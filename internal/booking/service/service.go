package service

import (
	"context"

	"github.com/quolpr/booking/internal/booking/model"
)

type OrderCreator interface {
	Create(ctx context.Context, order model.Order) (model.Order, error)
}
