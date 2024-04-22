package validator

import (
	"context"
	"errors"
	"fmt"

	"github.com/quolpr/booking/internal/booking/model"
)

var ErrValidation = errors.New("validation error")

type ValidationMsgError struct {
	Msgs []ValidationMsg
}

func (err *ValidationMsgError) Error() string {
	return fmt.Sprintf("Following validation errors happened: %+v", err.Msgs)
}

func (err *ValidationMsgError) Unwrap() error {
	return ErrValidation
}

type ValidationMsg struct {
	Field    string `json:"field"`
	ErrorTag string `json:"error_tag"`
	Payload  any    `json:"payload"`
}

type OrderCreationValidator interface {
	Validate(ctx context.Context, order model.Order) error
}
