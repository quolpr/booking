package order

import (
	"context"
	"testing"
	"time"

	"github.com/quolpr/booking/internal/transaction"
	"github.com/quolpr/booking/internal/transaction/inmem"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/repository"
	"github.com/quolpr/booking/internal/booking/repository/availability"
	"github.com/quolpr/booking/internal/booking/repository/order"
	"github.com/quolpr/booking/internal/booking/validator"
	orderValidator "github.com/quolpr/booking/internal/booking/validator/order"

	"github.com/stretchr/testify/assert"

	"github.com/quolpr/booking/pkg/days"
)

type TestingCase struct {
	name string

	ctx       context.Context
	availRepo repository.AvailabilityRepo
	orderRepo repository.OrdersRepo
	svc       *Creator
}

func setupInmemCase() TestingCase {
	ctx := context.Background()
	ctx = transaction.WithTransactionManager(ctx, &inmem.TransactionManager{})

	availRepo := availability.NewRepo(ctx)
	orderRepo := order.NewRepo()

	return TestingCase{
		name:      "InMem",
		ctx:       ctx,
		availRepo: availRepo,
		orderRepo: orderRepo,
		svc:       NewCreator(availRepo, orderRepo, orderValidator.NewCreateValidator(availRepo)),
	}
}

type Case struct {
	name string
	fn   func() TestingCase
}

func TestCreate(t *testing.T) {
	t.Parallel()

	cases := []Case{{"InMem", setupInmemCase}}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			t.Run("Works", func(t *testing.T) {
				testCreateWorks(t, c.fn)
			})

			// TODO: check that rollback works correctly
		})
	}
}

func testCreateWorks(t *testing.T, setupFunc func() TestingCase) {
	setup := setupFunc()
	ctx := setup.ctx

	avail := []model.RoomAvailability{
		{ID: "123", HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 1), Quota: 2},
		{ID: "456", HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 2), Quota: 1},
	}
	for _, av := range avail {
		_, err := setup.availRepo.Create(ctx, av)

		assert.Nil(t, err)
	}

	_, err := setup.svc.Create(setup.ctx, model.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "test@test.com",
		From:      days.Date(2024, 1, 1),
		To:        days.Date(2024, 1, 2),
	})

	assert.Nil(t, err)

	// Try to book same date
	_, err = setup.svc.Create(setup.ctx, model.Order{
		HotelID:   "reddison",
		RoomID:    "lux",
		UserEmail: "test2@test.com",
		From:      days.Date(2024, 1, 1),
		To:        days.Date(2024, 1, 2),
	})

	assert.NotNil(t, err)

	var validationErr *validator.ValidationMsgError
	assert.ErrorAs(t, err, &validationErr)

	assert.Equal(t, err,
		&validator.ValidationMsgError{Msgs: []validator.ValidationMsg{
			{Field: "from", ErrorTag: "unavailable", Payload: []time.Time{
				days.Date(2024, 1, 2),
			}},
		}},
	)
	res, err := setup.availRepo.GetAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res,
		[]model.RoomAvailability{
			{ID: "123", HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 1), Quota: 1},
			{ID: "456", HotelID: "reddison", RoomID: "lux", Date: days.Date(2024, 1, 2), Quota: 0},
		},
	)
}
