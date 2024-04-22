package order

import (
	"encoding/json"
	"errors"
	"net/http"

	jsonresp2 "github.com/quolpr/booking/internal/pkg/jsonresp"

	"github.com/quolpr/booking/internal/booking/model"
	"github.com/quolpr/booking/internal/booking/service"
	"github.com/quolpr/booking/internal/booking/validator"

	"github.com/quolpr/booking/internal/app/appctx"
)

type Handlers struct {
	orderCreator service.OrderCreator
}

func NewHandlers(orderCreator service.OrderCreator) *Handlers {
	return &Handlers{orderCreator: orderCreator}
}

func (h *Handlers) CreateOrder(r *http.Request) (jsonresp2.JSONResp, error) {
	logger, err := appctx.GetLogger(r.Context())
	if err != nil {
		return jsonresp2.JSONResp{}, err
	}

	var newOrder model.Order

	err = json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		return jsonresp2.JSONResp{}, &jsonresp2.JSONError{Type: "decode_failed", StatusCode: http.StatusBadRequest, Err: err}
	}

	result, err := h.orderCreator.Create(r.Context(), newOrder)
	var validationErr *validator.ValidationMsgError
	if err != nil && errors.As(err, &validationErr) {
		return jsonresp2.JSONResp{}, &jsonresp2.JSONError{
			Type:       "validation_failed",
			StatusCode: http.StatusBadRequest,
			Payload:    validationErr.Msgs,
			Err:        err,
		}
	} else if err != nil {
		logger.Warn("Failed to create order", "err", err)

		return jsonresp2.JSONResp{}, &jsonresp2.JSONError{
			Type:       "unknown",
			StatusCode: http.StatusBadRequest,
			Payload:    err.Error(),
			Err:        err,
		}
	}

	return jsonresp2.JSONResp{
		Body:       result,
		StatusCode: http.StatusCreated,
	}, nil
}
