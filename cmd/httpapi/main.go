package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/quolpr/booking/internal/app"
	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/pkg/days"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	app := app.NewApp(ctx)

	logger := app.Logger()

	availabilityRepo := app.ServiceProvider.AvailabilityRepo

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

	err := app.ServeHTTP(ctx)
	if err != nil {
		logger.Error("Server closed", "error", err)

		defer os.Exit(1)

		return
	}

	logger.Info("Server closed")
}
