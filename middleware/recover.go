package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type ErrorHandlerFunc = func(ctx context.Context, w http.ResponseWriter, err error)

func Recovery(sendError ErrorHandlerFunc) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				ctx := r.Context()
				err := recover()
				if err != nil {
					var panicErr error
					switch x := err.(type) {
					case string:
						panicErr = errors.New(x)
					case error:
						panicErr = x
					default:
						panicErr = fmt.Errorf("unknown panic: %v", x)
					}
					sendError(ctx, w, panicErr)
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}
