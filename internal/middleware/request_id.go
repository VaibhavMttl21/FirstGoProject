
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Private context key type to avoid collisions with other packages.
type requestIDKey struct{}

// RequestIDMiddleware generates a unique ID for every incoming request.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Generate 16 cryptographically secure random bytes.
		randomBytes := make([]byte, 16)

		if _, err := rand.Read(randomBytes); err != nil {
			http.Error(w, "Unable to generate request ID", http.StatusInternalServerError)
			return
		}

		// Convert the random bytes into a hexadecimal string.
		requestID := hex.EncodeToString(randomBytes)

		// Store the request ID in the request context.
		ctx := context.WithValue(
			r.Context(),
			requestIDKey{},
			requestID,
		)

		// Create a new request with the updated context.
		r = r.WithContext(ctx)

		// Expose the request ID in the HTTP response header.
		w.Header().Set("X-Request-ID", requestID)

		// Pass control to the next handler.
		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		return ""
	}

	return requestID
}
