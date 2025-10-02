package middlewarex

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// responseWriter wrapper to capture response body & status
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, &bytes.Buffer{}}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

// Table Logger Middleware
func RequestResponseLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Read request body (clone, so handler still gets it)
		var reqBody bytes.Buffer
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			reqBody.Write(bodyBytes)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore body
		}

		// Wrap response
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		// Print as table row
		fmt.Printf("\n%-8s | %-30s | %-3d | %-10s | Req: %-40s | Res: %-40s\n",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(start).String(),
			truncate(reqBody.String(), 40),
			truncate(lrw.body.String(), 40),
		)

	})
}

// Helper: truncate long strings
func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
