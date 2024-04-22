package availability

import (
	"context"
	"sync"
	"time"

	"github.com/quolpr/booking/internal/transaction"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/repository"

	"github.com/quolpr/booking/pkg/days"
	"github.com/quolpr/booking/pkg/randomstr"
)

var _ repository.AvailabilityRepo = (*Repo)(nil)

type AvailabilityError struct {
	unavailableDays []time.Time
}

func (e *AvailabilityError) Error() string { return "Some dates are not available" }

type availableKey struct {
	hotelID string
	roomID  string
	date    time.Time
}

type Repo struct {
	availability      map[string]model.RoomAvailability
	availabilityIndex map[availableKey]string

	// Actually it's better to use partial locks, but
	// I pick global lock as an easiest solution for now
	lock sync.RWMutex
}

func NewRepo(ctx context.Context) *Repo {
	r := &Repo{
		availabilityIndex: make(map[availableKey]string),
		availability:      make(map[string]model.RoomAvailability),
	}

	return r
}

func (r *Repo) Create(_ context.Context, av model.RoomAvailability) (model.RoomAvailability, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	av.ID = randomstr.Generate(20)

	r.availability[av.ID] = av
	r.availabilityIndex[availableKey{
		hotelID: av.HotelID,
		roomID:  av.RoomID,
		date:    av.Date,
	}] = av.ID

	return av, nil
}

func (r *Repo) Get(ctx context.Context, id string) (model.RoomAvailability, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.unsafeGet(ctx, id)
}

func (r *Repo) GetAll(_ context.Context) ([]model.RoomAvailability, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	availabilities := make([]model.RoomAvailability, 0, len(r.availability))

	for _, av := range r.availability {
		availabilities = append(availabilities, av)
	}

	return availabilities, nil
}

func (r *Repo) unsafeGet(_ context.Context, id string) (model.RoomAvailability, error) {
	av, ok := r.availability[id]
	if !ok {
		return model.RoomAvailability{}, repository.ErrNotFound
	}

	return av, nil
}

func (r *Repo) unsafeGetByAvailableKey(ctx context.Context, key availableKey) (model.RoomAvailability, error) {
	id, ok := r.availabilityIndex[key]
	if !ok {
		return model.RoomAvailability{}, repository.ErrNotFound
	}

	return r.unsafeGet(ctx, id)
}

func (r *Repo) unsafeGetUnavailableDays(
	ctx context.Context,
	hotelID, roomID string,
	from, to time.Time,
) ([]time.Time, error) {
	daysToBook, err := days.DaysBetween(from, to)
	if err != nil {
		return nil, err
	}

	unavailableDays := make([]time.Time, 0, len(daysToBook)/2)

	for _, dayToBook := range daysToBook {
		model, err := r.unsafeGetByAvailableKey(ctx, availableKey{hotelID: hotelID, roomID: roomID, date: dayToBook})

		if err != nil || model.Quota <= 0 {
			unavailableDays = append(unavailableDays, dayToBook)
		}
	}

	return unavailableDays, nil
}

func (r *Repo) GetUnavailableDays(
	ctx context.Context,
	hotelID, roomID string,
	from, to time.Time,
) ([]time.Time, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.unsafeGetUnavailableDays(ctx, hotelID, roomID, from, to)
}

func (r *Repo) DecreaseQuotasByDate(
	ctx context.Context,
	hotelID, roomID string,
	from, to time.Time,
) ([]model.RoomAvailability, error) {
	daysToBook, err := days.DaysBetween(from, to)
	if err != nil {
		return nil, err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	unavailableDays, err := r.unsafeGetUnavailableDays(ctx, hotelID, roomID, from, to)
	if err != nil {
		return nil, err
	}

	if len(unavailableDays) > 0 {
		return nil, &AvailabilityError{unavailableDays: unavailableDays}
	}

	result := make([]model.RoomAvailability, 0, len(daysToBook))

	for _, dayToBook := range daysToBook {
		model, err := r.unsafeGetByAvailableKey(ctx, availableKey{hotelID: hotelID, roomID: roomID, date: dayToBook})
		if err != nil {
			panic("This should not happen because we checked days in unsafeGetUnavailableDays before")
		}

		result = append(result, model)

		model.Quota--

		r.availability[model.ID] = model
	}

	transaction.RegisterTrCompensation(ctx, func(ctx context.Context) error {
		ids := make([]string, 0, len(result))

		for _, av := range result {
			ids = append(ids, av.ID)
		}

		r.IncreaseQuotasByIds(ctx, hotelID, roomID, ids)

		return nil
	})

	return result, nil
}

func (r *Repo) IncreaseQuotasByIds(ctx context.Context, hotelID, roomID string, ids []string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, id := range ids {
		model, err := r.unsafeGet(ctx, id)
		if err != nil {
			continue
		}

		model.Quota++

		r.availability[model.ID] = model
	}
}
