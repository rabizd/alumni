// Package routes maps every URL to the controller that answers it.
// web.go holds the routes for people (HTML pages), api.go the routes for
// programs (JSON). cmd/api/main.go registers both.
package routes

import (
	"net/http"

	"github.com/rabizd/alumni/internal/controller"
	"github.com/rabizd/alumni/internal/view"
)

// Web registers the HTML routes: the pages, and the /users, /announcements,
// /graduates, /job-postings and /applications screens.
func Web(mux *http.ServeMux) {
	// Lecture exercises and plain pages.
	mux.HandleFunc("GET /{$}", controller.Root)
	mux.HandleFunc("GET /main", controller.Main)
	mux.HandleFunc("GET /about", controller.About)
	mux.HandleFunc("GET /hello", controller.Hello)
	mux.HandleFunc("GET /hello/{name}", controller.HelloName)
	mux.HandleFunc("GET /sum/{number1}/{number2}", controller.Sum)
	mux.HandleFunc("GET /temporary", controller.Temporary)

	// Images the pages use, e.g. the campus photos behind the main page.
	mux.Handle("GET /static/", view.Static)

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

	// Graduates (the alumni table) as HTML pages. The JSON is at /alumni.
	graduate := controller.GraduateController{}
	mux.HandleFunc("GET /graduates", graduate.Index)
	mux.HandleFunc("GET /graduates/new", graduate.Create)
	mux.HandleFunc("POST /graduates", graduate.Store)
	mux.HandleFunc("GET /graduates/{id}", graduate.Show)
	mux.HandleFunc("GET /graduates/{id}/edit", graduate.Edit)
	mux.HandleFunc("POST /graduates/{id}", graduate.Update)
	mux.HandleFunc("POST /graduates/{id}/delete", graduate.Destroy)

	// Job postings as HTML pages. A posting's page lists its applications.
	jobPosting := controller.JobPostingController{}
	mux.HandleFunc("GET /job-postings", jobPosting.Index)
	mux.HandleFunc("GET /job-postings/new", jobPosting.Create)
	mux.HandleFunc("POST /job-postings", jobPosting.Store)
	mux.HandleFunc("GET /job-postings/{id}", jobPosting.Show)
	mux.HandleFunc("GET /job-postings/{id}/edit", jobPosting.Edit)
	mux.HandleFunc("POST /job-postings/{id}", jobPosting.Update)
	mux.HandleFunc("POST /job-postings/{id}/delete", jobPosting.Destroy)

	// Applications as HTML. Applying is nested under the posting; the status
	// form and delete work on the application's own id.
	application := controller.ApplicationController{}
	mux.HandleFunc("GET /applications", application.Index)
	mux.HandleFunc("POST /job-postings/{id}/applications", application.Store)
	mux.HandleFunc("POST /applications/{id}", application.Update)
	mux.HandleFunc("POST /applications/{id}/delete", application.Destroy)
}
