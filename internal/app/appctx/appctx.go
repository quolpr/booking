package appctx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type contextKey string

const (
	loggerKey = contextKey("logger")
	reqIDKey  = contextKey("reqId")
)

var ErrNotFound = errors.New("key in ctx not found")

func WithRequestID(c context.Context, reqID string) context.Context {
	return context.WithValue(c, reqIDKey, reqID)
}

func WithLogger(c context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(c, loggerKey, logger)
}

func GetLogger(ctx context.Context) (*slog.Logger, error) {
	if res := ctx.Value(loggerKey); res != nil {
		logger, ok := res.(*slog.Logger)

		if !ok {
			panic(fmt.Sprintf("Tried to get logger from ctx, but it has different type: %#v", res))
		}

		return logger, nil
	} else {
		return nil, fmt.Errorf("logger not found: %w", ErrNotFound)
	}
}

func GetRequestID(ctx context.Context) (string, error) {
	if res := ctx.Value(reqIDKey); res != nil {
		reqID, ok := res.(string)

		if !ok {
			panic(fmt.Sprintf("Tried to get request id from ctx, but it has different type: %#v", res))
		}

		return reqID, nil
	} else {
		return "", fmt.Errorf("request id not found: %w", ErrNotFound)
	}
}
