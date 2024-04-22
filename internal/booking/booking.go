package booking

import (
	"context"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/repository"
	"github.com/quolpr/booking/internal/booking/service"
)

// This is Bounded Context

type Domain struct {
	ordersRepo   repository.OrdersRepo
	availRepo    repository.AvailabilityRepo
	orderCreator service.OrderCreator
}

func NewDomain(
	ordersRepo repository.OrdersRepo,
	availRepo repository.AvailabilityRepo,
	orderCreator service.OrderCreator,
) *Domain {
	return &Domain{
		ordersRepo:   ordersRepo,
		orderCreator: orderCreator,
		availRepo:    availRepo,
	}
}

func (d *Domain) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	return d.orderCreator.Create(ctx, order)
}

func (d *Domain) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	return d.ordersRepo.GetAll(ctx)
}

func (d *Domain) GetOrder(ctx context.Context, id string) (model.Order, error) {
	return d.ordersRepo.Get(ctx, id)
}

func (d *Domain) CreateAvailability(ctx context.Context, av model.RoomAvailability) (model.RoomAvailability, error) {
	return d.availRepo.Create(ctx, av)
}

func (d *Domain) GetAllAvailabilities(ctx context.Context) ([]model.RoomAvailability, error) {
	return d.availRepo.GetAll(ctx)
}
