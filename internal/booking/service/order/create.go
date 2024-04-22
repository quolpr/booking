package order

import (
	"context"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/repository"
	"github.com/quolpr/booking/internal/booking/service"
	"github.com/quolpr/booking/internal/booking/validator"

	"github.com/quolpr/booking/internal/app/transaction"
)

var _ service.OrderCreator = (*Creator)(nil)

type Creator struct {
	availRepo  repository.AvailabilityRepo
	ordersRepo repository.OrdersRepo
	validator  validator.OrderCreationValidator
}

func NewCreator(
	availRepo repository.AvailabilityRepo,
	ordersRepo repository.OrdersRepo,
	validator validator.OrderCreationValidator,
) *Creator {
	return &Creator{
		availRepo:  availRepo,
		ordersRepo: ordersRepo,
		validator:  validator,
	}
}

func (s *Creator) Create(ctx context.Context, order model.Order) (createdOrder model.Order, err error) {
	err = s.validator.Validate(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	err = transaction.InTransaction(ctx, func(ctx context.Context) error {
		_, err := s.availRepo.DecreaseQuotasByDate(
			ctx, order.HotelID, order.RoomID,
			order.From, order.To,
		)
		if err != nil {
			return err
		}

		createdOrder, err = s.ordersRepo.Create(ctx, order)
		if err != nil {
			return err
		}

		return nil
	})

	return
}
