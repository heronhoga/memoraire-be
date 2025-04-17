package utils

import "net/http"

func WithMiddleware(h http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.Handler {
	handler := http.Handler(h)
	for _, mw := range mws {
		handler = mw(handler)
	}
	return handler
}
