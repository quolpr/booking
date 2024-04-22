package order

import (
	"context"
	"strings"
	"time"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/validator"
)

type unavailableDaysGetter interface {
	GetUnavailableDays(ctx context.Context, hotelID, roomID string, from, to time.Time) ([]time.Time, error)
}

type createValidator struct {
	unavailableDaysGetter unavailableDaysGetter
}

func NewCreateValidator(getter unavailableDaysGetter) *createValidator {
	return &createValidator{getter}
}

func (v *createValidator) Validate(ctx context.Context, order model.Order) error {
	msgs := make([]validator.ValidationMsg, 0)

	if len(order.HotelID) == 0 {
		msgs = append(msgs, validator.ValidationMsg{Field: "hotel_id", ErrorTag: "required"})
	}

	if len(order.RoomID) == 0 {
		msgs = append(msgs, validator.ValidationMsg{Field: "room_id", ErrorTag: "required"})
	}

	if len(order.UserEmail) == 0 {
		msgs = append(msgs, validator.ValidationMsg{Field: "email", ErrorTag: "required"})
	}

	if !strings.Contains(order.UserEmail, "@") {
		msgs = append(msgs, validator.ValidationMsg{Field: "email", ErrorTag: "invalid"})
	}

	if order.From.After(order.To) {
		msgs = append(msgs, validator.ValidationMsg{Field: "from", ErrorTag: "wrong_range"})
	}

	unavailableDays, err := v.unavailableDaysGetter.GetUnavailableDays(
		ctx,
		order.HotelID, order.RoomID,
		order.From, order.To,
	)
	if err != nil {
		return err
	}

	if len(unavailableDays) > 0 {
		msgs = append(msgs, validator.ValidationMsg{Field: "from", ErrorTag: "unavailable", Payload: unavailableDays})
	}

	if len(msgs) > 0 {
		return &validator.ValidationMsgError{Msgs: msgs}
	}

	return nil
}
