package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/quolpr/booking/internal/pkg/transaction"
	"github.com/quolpr/booking/internal/pkg/transaction/inmem"

	orderHandlers "github.com/quolpr/booking/internal/booking/httpapi/v1/order"
	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/repository"
	"github.com/quolpr/booking/internal/booking/repository/availability"
	orderRepo "github.com/quolpr/booking/internal/booking/repository/order"
	"github.com/quolpr/booking/internal/booking/service"
	orderSvc "github.com/quolpr/booking/internal/booking/service/order"
	"github.com/quolpr/booking/internal/booking/validator"
	orderValidator "github.com/quolpr/booking/internal/booking/validator/order"

	"github.com/quolpr/booking/pkg/days"
)

type serviceProvider struct {
	logger    *slog.Logger
	trManager transaction.TransactionManager

	ordersHandler *orderHandlers.Handlers

	ordersRepo       repository.OrdersRepo
	availabilityRepo repository.AvailabilityRepo

	createOrderService service.OrderCreator

	creationOrderValidator validator.OrderCreationValidator
}

func newServiceProvider(ctx context.Context, _ *config) *serviceProvider {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ordersRepo := orderRepo.NewRepo()
	availabilityRepo := availability.NewRepo(ctx)
	createOrderValidator := orderValidator.NewCreateValidator(availabilityRepo)
	createOrderService := orderSvc.NewCreator(availabilityRepo, ordersRepo, createOrderValidator)

	avail := []model.RoomAvailability{
		{HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 1), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 2), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 3), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 4), Quota: 1},
		{HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 5), Quota: 0},
	}

	for _, av := range avail {
		availabilityRepo.Create(ctx, av)
	}

	return &serviceProvider{
		logger:                 logger,
		trManager:              &inmem.TransactionManager{},
		ordersHandler:          orderHandlers.NewHandlers(createOrderService),
		ordersRepo:             ordersRepo,
		availabilityRepo:       availabilityRepo,
		createOrderService:     createOrderService,
		creationOrderValidator: createOrderValidator,
	}
}
