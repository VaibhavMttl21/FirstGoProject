package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		start := time.Now()

	log.Printf(
		"Request started | Method: %s | Path: %s",
		r.Method,
		r.URL.Path,
	)
	recorder := &statusRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	next.ServeHTTP(recorder, r)
//we intercept the method call, add our own behavior, 
// and then delegate the actual response-writing work to Go's original writer.
	duration := time.Since(start)

		log.Printf(
		"Request completed | Method: %s | Path: %s | Status: %d | Duration: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
	)
	
	})
}