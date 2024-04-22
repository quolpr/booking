package order

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/quolpr/booking/internal/app/appctx"
	"github.com/quolpr/booking/internal/app/jsonresp"
	"github.com/quolpr/booking/internal/model"
	"github.com/quolpr/booking/internal/service"
	"github.com/quolpr/booking/internal/validator"
)

type Handlers struct {
	orderCreator service.OrderCreator
}

func NewHandlers(orderCreator service.OrderCreator) *Handlers {
	return &Handlers{orderCreator: orderCreator}
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

	result, err := h.orderCreator.Create(r.Context(), newOrder)
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
	}, errors.New("not implemented")
}
