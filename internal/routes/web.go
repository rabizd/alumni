// Package routes maps every URL to the controller that answers it.
// web.go holds the routes for people (HTML pages), api.go the routes for
// programs (JSON). cmd/api/main.go registers both.
package routes

import (
	"net/http"

	"github.com/rabizd/alumni/internal/controller"
)

// Web registers the HTML routes: the pages, and the /users and /announcements screens.
func Web(mux *http.ServeMux) {
	// Lecture exercises and plain pages.
	mux.HandleFunc("GET /{$}", controller.Root)
	mux.HandleFunc("GET /main", controller.Main)
	mux.HandleFunc("GET /about", controller.About)
	mux.HandleFunc("GET /hello", controller.Hello)
	mux.HandleFunc("GET /hello/{name}", controller.HelloName)
	mux.HandleFunc("GET /sum/{number1}/{number2}", controller.Sum)
	mux.HandleFunc("GET /temporary", controller.Temporary)

	// Users as HTML pages. HTML forms can only send GET and POST, so update
	// and delete are POSTs here.
	user := controller.UserController{}
	mux.HandleFunc("GET /users", user.Index)
	mux.HandleFunc("GET /users/new", user.Create)
	mux.HandleFunc("POST /users", user.Store)
	mux.HandleFunc("GET /users/{id}", user.Show)
	mux.HandleFunc("GET /users/{id}/edit", user.Edit)
	mux.HandleFunc("POST /users/{id}", user.Update)
	mux.HandleFunc("POST /users/{id}/delete", user.Destroy)

	// Announcements as HTML pages: the interface to manage them.
	announcement := controller.AnnouncementController{}
	mux.HandleFunc("GET /announcements", announcement.Index)
	mux.HandleFunc("GET /announcements/new", announcement.Create)
	mux.HandleFunc("POST /announcements", announcement.Store)
	mux.HandleFunc("GET /announcements/{id}", announcement.Show)
	mux.HandleFunc("GET /announcements/{id}/edit", announcement.Edit)
	mux.HandleFunc("POST /announcements/{id}", announcement.Update)
	mux.HandleFunc("POST /announcements/{id}/delete", announcement.Destroy)
}
