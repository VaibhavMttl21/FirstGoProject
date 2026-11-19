package middleware

import "net/http"

func Chain(
	handler http.Handler,
	middlewares ...func(http.Handler) http.Handler,
) http.Handler {

	// Start with the original handler.
	finalHandler := handler

	// Apply middleware in reverse order.
	// But remember: each middleware wraps the handler passed to it.
	for i := len(middlewares) - 1; i >= 0; i-- {
		finalHandler = middlewares[i](finalHandler)
	}

	return finalHandler
}



// The loop works like this:
// Iteration	Middleware applied	Result
// 1	Logging	Logging(handler)
// 2	Recovery	Recovery(Logging(handler))
// 3	Request ID	RequestID(Recovery(Logging(handler)))

// The final execution flow is:
// Request ID
// Recovery
// Logging
// Handler

// If we iterated forward, Logging would execute first and Request ID last—which is not what we want.