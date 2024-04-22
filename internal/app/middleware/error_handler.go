package middleware

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/quolpr/booking/internal/jsonresp"

	"github.com/quolpr/booking/internal/app/appctx"
)

type jsonWriter struct {
	writer http.ResponseWriter
	logger *slog.Logger
}

func (j *jsonWriter) writeJSON(code int, data any) {
	j.writer.Header().Add("Content-Type", "application/json")
	j.writer.WriteHeader(code)

	err := json.NewEncoder(j.writer).Encode(data)
	if err != nil {
		j.logger.Error("Failed to write response", "err", err.Error())
	}
}

func ErrorHandler(h func(r *http.Request) (jsonresp.JSONResp, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, err := appctx.GetLogger(r.Context())
		if err != nil {
			slog.Error("Failed to get logger from context", "err", err.Error())
			logger = slog.Default()
		}

		jsonWriter := &jsonWriter{writer: w, logger: logger}

		res, err := h(r)

		var jsonErr *jsonresp.JSONError

		if errors.As(err, &jsonErr) {
			logger.Info("Json error happened", "err", err.Error())

			jsonWriter.writeJSON(jsonErr.StatusCode, jsonErr)
			return
		} else if err != nil {
			logger.Error("Error happened", "err", err)

			jsonWriter.writeJSON(http.StatusInternalServerError, jsonresp.JSONError{Type: "internal_error"})
			return
		}

		jsonWriter.writeJSON(res.StatusCode, res.Body)
	}
}
