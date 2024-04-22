package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/quolpr/booking/internal/app/appctx"
)

func PanicRecover(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, err := appctx.GetLogger(r.Context())
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic happened", "err", err, "stack", string(debug.Stack()))

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h(w, r)
	}
}
