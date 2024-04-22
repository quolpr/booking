package repository

import (
	"context"
	"errors"
	"time"

	"github.com/quolpr/booking/internal/model"
)

var ErrNotFound = errors.New("not found")

type OrdersRepo interface {
	Create(ctx context.Context, order model.Order) (model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Get(ctx context.Context, id string) (model.Order, error)
}

type AvailabilityRepo interface {
	Get(ctx context.Context, id string) (model.RoomAvailability, error)
	GetAll(_ context.Context) ([]model.RoomAvailability, error)
	GetUnavailableDays(ctx context.Context, hotelID, roomID string, from, to time.Time) ([]time.Time, error)
	DecreaseQuotasByDate(ctx context.Context, hotelID, roomID string, from, to time.Time) ([]model.RoomAvailability, error)
	IncreaseQuotasByIds(ctx context.Context, hotelID, roomID string, ids []string)
	Create(ctx context.Context, av model.RoomAvailability)
}
