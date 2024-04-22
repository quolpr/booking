package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/quolpr/booking/internal/pkg/transaction"
	"github.com/quolpr/booking/internal/pkg/transaction/inmem"

	orderHandlers "github.com/quolpr/booking/internal/booking/httpapi/order/v1"
	"github.com/quolpr/booking/internal/booking/repository"
	"github.com/quolpr/booking/internal/booking/repository/availability"
	orderRepo "github.com/quolpr/booking/internal/booking/repository/order"
	"github.com/quolpr/booking/internal/booking/service"
	orderSvc "github.com/quolpr/booking/internal/booking/service/order"
	"github.com/quolpr/booking/internal/booking/validator"
	orderValidator "github.com/quolpr/booking/internal/booking/validator/order"
)

type serviceProvider struct {
	Logger    *slog.Logger
	TrManager transaction.TransactionManager

	OrdersHandler *orderHandlers.Handlers

	OrdersRepo       repository.OrdersRepo
	AvailabilityRepo repository.AvailabilityRepo

	CreateOrderService service.OrderCreator

	CreationOrderValidator validator.OrderCreationValidator
}

func newServiceProvider(ctx context.Context, _ *config) *serviceProvider {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ordersRepo := orderRepo.NewRepo()
	availabilityRepo := availability.NewRepo(ctx)
	createOrderValidator := orderValidator.NewCreateValidator(availabilityRepo)
	createOrderService := orderSvc.NewCreator(availabilityRepo, ordersRepo, createOrderValidator)

	return &serviceProvider{
		Logger:                 logger,
		TrManager:              &inmem.TransactionManager{},
		OrdersHandler:          orderHandlers.NewHandlers(createOrderService),
		OrdersRepo:             ordersRepo,
		AvailabilityRepo:       availabilityRepo,
		CreateOrderService:     createOrderService,
		CreationOrderValidator: createOrderValidator,
	}
}
