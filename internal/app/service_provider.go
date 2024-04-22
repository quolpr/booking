package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/quolpr/booking/internal/booking"
	orderHandlers "github.com/quolpr/booking/internal/httpapi/order/v1"
	"github.com/quolpr/booking/internal/transaction"
	"github.com/quolpr/booking/internal/transaction/inmem"

	"github.com/quolpr/booking/internal/booking/repository/availability"
	orderRepo "github.com/quolpr/booking/internal/booking/repository/order"
	orderSvc "github.com/quolpr/booking/internal/booking/service/order"
	orderValidator "github.com/quolpr/booking/internal/booking/validator/order"
)

type serviceProvider struct {
	Logger    *slog.Logger
	TrManager transaction.TransactionManager

	OrdersHandler *orderHandlers.Handlers

	Booking *booking.Domain
}

func newServiceProvider(ctx context.Context, _ *config) *serviceProvider {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ordersRepo := orderRepo.NewRepo()
	availabilityRepo := availability.NewRepo(ctx)
	createOrderValidator := orderValidator.NewCreateValidator(availabilityRepo)
	createOrderService := orderSvc.NewCreator(availabilityRepo, ordersRepo, createOrderValidator)

	bookingDomain := booking.NewDomain(ordersRepo, availabilityRepo, createOrderService)

	return &serviceProvider{
		Logger:        logger,
		TrManager:     &inmem.TransactionManager{},
		OrdersHandler: orderHandlers.NewHandlers(bookingDomain),
		Booking:       bookingDomain,
	}
}
