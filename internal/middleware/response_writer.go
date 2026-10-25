package middleware

import "net/http"

type statusRecorder struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (sr *statusRecorder) WriteHeader(statusCode int) {
	// Only the first status code should be recorded.
	if sr.wroteHeader {
		return
	}

	sr.wroteHeader = true
	sr.statusCode = statusCode

	// Forward the status to the real HTTP response writer.
	sr.ResponseWriter.WriteHeader(statusCode)
}

func (sr *statusRecorder) Write(body []byte) (int, error) {
	// Writing a body without WriteHeader means status 200.
	if !sr.wroteHeader {
		sr.WriteHeader(http.StatusOK)
	}

	// Forward the body to the real HTTP response writer.
	return sr.ResponseWriter.Write(body)
}