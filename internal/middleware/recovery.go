package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if recovered := recover(); recovered != nil {

				requestID := GetRequestID(r.Context())

				log.Printf(
					"Panic recovered | RequestID: %s | Method: %s | Path: %s | Error: %v\n%s",
					requestID,
					r.Method,
					r.URL.Path,
					recovered,
					debug.Stack(),
				)

				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}



// Request ID middleware
// Recovery middleware — defer registered
// Logging middleware
// Handler — panic occurs!
// Recovery catches panic and writes HTTP 500
// Logging middleware completes

// Recovery prevents an unhandled panic from propagating, but it cannot undo an HTTP response that has already been sent.