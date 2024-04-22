package book

import (
	"context"

	"github.com/quolpr/booking/internal/booking"
	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/transaction"
)

type Creator struct {
	booking booking.Domain
	// Notifier
}

func NewCreator(booking booking.Domain) *Creator {
	return &Creator{booking: booking}
}

func (c *Creator) Create(ctx context.Context) (res model.Order, err error) {
	err = transaction.InTransaction(ctx, func(ctx context.Context) error {
		res, err = c.booking.CreateOrder(ctx, model.Order{})

		// notify

		return err
	})

	return
}
