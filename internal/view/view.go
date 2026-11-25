// Package view is the V in MVC: it turns data into what the client receives,
// an HTML page or a JSON body. It shows things; it never decides anything.
package view

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/rabizd/alumni/internal/model"
)

// The HTML lives in its own files; go:embed bakes them into the binary at
// build time, so the server stays a single executable with no loose files.
//
//go:embed templates/*.html
var templateFiles embed.FS

// currentUser and menu are placeholders here: HTML swaps in the signed-in
// user of each request and the menu links their role may see (menu.go).
var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"currentUser": func() *model.User { return nil },
	"menu":        func() []MenuItem { return nil },
}).ParseFS(templateFiles, "templates/*.html"))

// The OpenAPI document is baked in the same way, so the running server always
// serves the spec that was built alongside it.
//
//go:embed openapi.json
var openAPISpec []byte

// Images and other files the pages link to are baked in too, under static/.
//
//go:embed static
var staticFiles embed.FS

// Static serves the files in static/ at the same path, e.g.
// /static/img/campus-1.jpg.
var Static = http.FileServerFS(staticFiles)

// HTML renders one page from templates/ with the given status code. data is
// what the page shows; pages that show nothing take nil.
func HTML(w http.ResponseWriter, r *http.Request, status int, page string, data any) {
	var me *model.User
	if u, ok := model.CurrentUser(r.Context()); ok {
		me = &u
	}
	t, err := templates.Clone()
	if err != nil {
		log.Printf("rendering %s: %v", page, err)
		http.Error(w, "page could not be drawn", http.StatusInternalServerError)
		return
	}
	items := menuFor(me, r.URL.Path)
	t.Funcs(template.FuncMap{
		"currentUser": func() *model.User { return me },
		"menu":        func() []MenuItem { return items },
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := t.ExecuteTemplate(w, page, data); err != nil {
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
