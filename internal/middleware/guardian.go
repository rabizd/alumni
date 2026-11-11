package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

// notFoundBody is the JSON every 404 gets.
type notFoundBody struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
	Path   string `json:"path"`
}

// guardWriter holds back a 404 so Guardian can rewrite it; any other status
// goes straight through.
type guardWriter struct {
	http.ResponseWriter
	notFound bool
	body     bytes.Buffer
}

func (g *guardWriter) WriteHeader(code int) {
	if code == http.StatusNotFound {
		g.notFound = true
		return
	}
	g.ResponseWriter.WriteHeader(code)
}

func (g *guardWriter) Write(b []byte) (int, error) {
	if g.notFound {
		return g.body.Write(b)
	}
	return g.ResponseWriter.Write(b)
}

// Guardian turns every 404 into the same JSON with a reason. A controller's
// own message ("no user with that id") becomes the reason; when no route
// matched at all, Guardian writes the reason itself.
func Guardian(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g := &guardWriter{ResponseWriter: w}
		next.ServeHTTP(g, r)
		if !g.notFound {
			return
		}

		reason := strings.TrimSpace(g.body.String())
		if reason == "" || reason == "404 page not found" { // Go's answer when no route matches
			reason = "no route matches " + r.Method + " " + r.URL.Path
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(notFoundBody{
			Status: http.StatusNotFound,
			Error:  "Not Found",
			Reason: reason,
			Path:   r.URL.Path,
		})
	})
}
