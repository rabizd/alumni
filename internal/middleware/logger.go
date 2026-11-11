// Package middleware holds code that runs around every request, before and
// after the controller.
package middleware

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// statusRecorder remembers the status code the controller sends, because
// http.ResponseWriter does not let us read it back.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Logger prints one line per request to the terminal:
//
//	10:20:30  POST    /api/users  201  18.0ms
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// 200 is what Go sends when the controller never calls WriteHeader.
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		ms := float64(time.Since(start).Microseconds()) / 1000
		line := fmt.Sprintf("%s  %-7s %s  %d  %.1fms\n",
			start.Format("15:04:05"), r.Method, r.URL.Path, rec.status, ms)
		fmt.Print(line)
		if rec.status >= 400 {
			// The file keeps lines for days, so it also needs the date.
			fmt.Fprint(ErrorLog, start.Format("2006/01/02 ")+line)
		}
	})
}

// ErrorLog receives the line of every request that ends in 4xx or 5xx.
// main points it at logs/error.log.
var ErrorLog io.Writer = io.Discard
