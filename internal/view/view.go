// Package view is the V in MVC: it turns data into what the client receives,
// an HTML page or a JSON body. It shows things; it never decides anything.
package view

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// The HTML lives in its own files; go:embed bakes them into the binary at
// build time, so the server stays a single executable with no loose files.
//
//go:embed templates/*.html
var templateFiles embed.FS

var templates = template.Must(template.ParseFS(templateFiles, "templates/*.html"))

// The OpenAPI document is baked in the same way, so the running server always
// serves the spec that was built alongside it.
//
//go:embed openapi.json
var openAPISpec []byte

// HTML renders one page from templates/.
func HTML(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, page, nil); err != nil {
		log.Printf("rendering %s: %v", page, err)
	}
}

// JSON writes v as the response body with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// The status line is already sent, so all we can do is record it.
		log.Printf("writing JSON response: %v", err)
	}
}

// OpenAPI writes the OpenAPI document as it is.
func OpenAPI(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(openAPISpec)
}
