package booking

import (
	"context"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/repository"
	"github.com/quolpr/booking/internal/booking/service"
)

// If I need notification, then notification context will be, and this will require this interface:
// type OrderCreatedNotifier interface {
// 	NotifyOrderCreated(ctx context.Context, order model.Order) error
// }

// OR

// Also I can put event to kafka that order was created,
// And then other bounded context will process it

// OR

// I could write cross-domain BookingService that will be calling both billing and booking

// ----

// Also, for this:

// - появятся скидки, промокоды, программы лояльности

// 1. I would create billing domain
// 2. I would create cross-domain service, that will be calling both billing and booking
// OR
// 1. Billing will be emitting Paid event, then booking will be listening to this event

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
