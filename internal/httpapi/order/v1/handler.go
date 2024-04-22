package order

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/quolpr/booking/internal/jsonresp"

	"github.com/quolpr/booking/internal/booking"
	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/validator"

	"github.com/quolpr/booking/internal/app/appctx"
)

type Handlers struct {
	booking *booking.Domain
}

func NewHandlers(d *booking.Domain) *Handlers {
	return &Handlers{booking: d}
}

func (h *Handlers) CreateOrder(r *http.Request) (jsonresp.JSONResp, error) {
	logger, err := appctx.GetLogger(r.Context())
	if err != nil {
		return jsonresp.JSONResp{}, err
	}

	var newOrder model.Order

	err = json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		return jsonresp.JSONResp{}, &jsonresp.JSONError{Type: "decode_failed", StatusCode: http.StatusBadRequest, Err: err}
	}

	result, err := h.booking.CreateOrder(r.Context(), newOrder)
	var validationErr *validator.ValidationMsgError
	if err != nil && errors.As(err, &validationErr) {
		return jsonresp.JSONResp{}, &jsonresp.JSONError{
			Type:       "validation_failed",
			StatusCode: http.StatusBadRequest,
			Payload:    validationErr.Msgs,
			Err:        err,
		}
	} else if err != nil {
		logger.Warn("Failed to create order", "err", err)

		return jsonresp.JSONResp{}, &jsonresp.JSONError{
			Type:       "unknown",
			StatusCode: http.StatusBadRequest,
			Payload:    err.Error(),
			Err:        err,
		}
	}

	return jsonresp.JSONResp{
		Body:       result,
		StatusCode: http.StatusCreated,
	}, nil
}
