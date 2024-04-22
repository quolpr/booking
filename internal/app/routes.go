package app

import (
	"net/http"

	"github.com/quolpr/booking/internal/app/jsonresp"
	"github.com/quolpr/booking/internal/app/middleware"
)

func newRoutes(serviceProvider *serviceProvider) *http.ServeMux {
	mux := http.NewServeMux()

	applyMiddleware := func(h func(r *http.Request) (jsonresp.JSONResp, error)) http.HandlerFunc {
		return middleware.Logger(
			middleware.PanicRecover(
				middleware.ErrorHandler(h),
			), serviceProvider.logger,
		)
	}

	mux.HandleFunc("POST /v1/orders", applyMiddleware(serviceProvider.ordersHandler.CreateOrder))

	return mux
}
