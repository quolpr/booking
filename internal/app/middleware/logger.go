package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/quolpr/booking/internal/app/appctx"
	"github.com/quolpr/booking/pkg/randomstr"
)

func Logger(
	h func(w http.ResponseWriter, r *http.Request),
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger
		start := time.Now()

		reqID := randomstr.Generate(16)
		ctx := r.Context()

		logger = logger.With("path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr, "req_id", reqID)

		ctx = appctx.WithLogger(ctx, logger)
		ctx = appctx.WithRequestID(ctx, reqID)

		r = r.Clone(ctx)

		w.Header().Add("X-Request-Id", reqID)

		h(w, r)

		logger.InfoContext(
			ctx, "Done handling",
			"duration", time.Since(start),
			"happened_at", start,
		)
	}
}
