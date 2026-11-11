package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

// errorBody is the JSON every guarded error gets.
type errorBody struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
	Path   string `json:"path"`
}

// guardWriter holds back a response whose status matches catch, so a
// guardian can rewrite it; any other status goes straight through.
type guardWriter struct {
	http.ResponseWriter
	catch  func(code int) bool
	caught bool
	status int
	body   bytes.Buffer
}

func (g *guardWriter) WriteHeader(code int) {
	if g.catch(code) {
		g.caught, g.status = true, code
		return
	}
	g.ResponseWriter.WriteHeader(code)
}

func (g *guardWriter) Write(b []byte) (int, error) {
	if g.caught {
		return g.body.Write(b)
	}
	return g.ResponseWriter.Write(b)
}

// send writes the held-back error as JSON.
func (g *guardWriter) send(r *http.Request, reason string) {
	g.ResponseWriter.Header().Set("Content-Type", "application/json")
	g.ResponseWriter.WriteHeader(g.status)
	json.NewEncoder(g.ResponseWriter).Encode(errorBody{
		Status: g.status,
		Error:  http.StatusText(g.status),
		Reason: reason,
		Path:   r.URL.Path,
	})
}

// NotFoundGuardian turns every 404 into the same JSON with a reason. A
// controller's own message ("no user with that id") becomes the reason; when
// no route matched at all, it writes the reason itself.
func NotFoundGuardian(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g := &guardWriter{ResponseWriter: w, catch: func(code int) bool { return code == http.StatusNotFound }}
		next.ServeHTTP(g, r)
		if !g.caught {
			return
		}

		reason := strings.TrimSpace(g.body.String())
		if reason == "" || reason == "404 page not found" { // Go's answer when no route matches
			reason = "no route matches " + r.Method + " " + r.URL.Path
		}
		g.send(r, reason)
	})
}

// ServerErrorGuardian turns every 5xx into the same JSON with a reason, and
// catches a panic so the client gets a 500 instead of a dropped connection.
// The panic's details go to the terminal, never to the client.
func ServerErrorGuardian(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g := &guardWriter{ResponseWriter: w, catch: func(code int) bool { return code >= 500 }}
		defer func() {
			if p := recover(); p != nil {
				log.Printf("panic on %s %s: %v\n%s", r.Method, r.URL.Path, p, debug.Stack())
				g.caught, g.status = true, http.StatusInternalServerError
				g.body.Reset()
				g.body.WriteString("unexpected server error; the details are in the server log")
			}
			if g.caught {
				g.send(r, strings.TrimSpace(g.body.String()))
			}
		}()
		next.ServeHTTP(g, r)
	})
}
