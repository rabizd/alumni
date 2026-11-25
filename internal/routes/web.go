// Package routes maps every URL to the controller that answers it.
// web.go holds the routes for people (HTML pages), api.go the routes for
// programs (JSON). cmd/api/main.go registers both.
package routes

import (
	"net/http"

	"github.com/rabizd/alumni/internal/controller"
	"github.com/rabizd/alumni/internal/middleware"
	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// Web registers the HTML routes: the pages, and the /users, /announcements,
// /graduates, /job-postings and /applications screens.
//
// Who may open which page follows the role table in DDD.md: RequireRole
// needs at least that role, RequireSelfOrRole also lets users reach their own
// /users/{id} pages. Routes without either are open to everyone, guests too.
func Web(mux *http.ServeMux) {
	// Lecture exercises and plain pages.
	mux.HandleFunc("GET /{$}", controller.Root)
	mux.HandleFunc("GET /main", controller.Main)
	mux.HandleFunc("GET /about", controller.About)
	mux.HandleFunc("GET /404", controller.NotFound)
	mux.HandleFunc("GET /hello", controller.Hello)
	mux.HandleFunc("GET /hello/{name}", controller.HelloName)
	mux.HandleFunc("GET /sum/{number1}/{number2}", controller.Sum)
	mux.HandleFunc("GET /temporary", controller.Temporary)

	// Images the pages use, e.g. the campus photos behind the main page.
	mux.Handle("GET /static/", view.Static)

	// Sign in and sign out. Signing out is a POST so a link or a prefetch
	// cannot sign anyone out by accident.
	auth := controller.AuthController{}
	mux.HandleFunc("GET /login", auth.ShowLogin)
	mux.HandleFunc("POST /login", auth.Login)
	mux.HandleFunc("POST /logout", auth.Logout)

	// The signed-in user's own page: profile, postings, sent and received applications.
	mux.HandleFunc("GET /me", middleware.RequireRole(model.RoleMember, controller.MeController{}.Show))

	// Admin panel: counts, roles, and delete buttons for users, postings and announcements.
	mux.HandleFunc("GET /admin", middleware.RequireRole(model.RoleAdmin, controller.AdminController{}.Index))

	// Sign-up page: creates a user with a password.
	register := controller.RegisterController{}
	mux.HandleFunc("GET /register", register.Show)
	mux.HandleFunc("POST /register", register.Store)

	// Users as HTML pages. HTML forms can only send GET and POST, so update
	// and delete are POSTs here.
	user := controller.UserController{}
	mux.HandleFunc("GET /users", middleware.RequireRole(model.RoleGraduate, user.Index))
	mux.HandleFunc("GET /users/new", middleware.RequireRole(model.RoleAdmin, user.Create))
	mux.HandleFunc("POST /users", middleware.RequireRole(model.RoleAdmin, user.Store))
	mux.HandleFunc("GET /users/{id}", middleware.RequireSelfOrRole(model.RoleGraduate, user.Show))
	mux.HandleFunc("GET /users/{id}/edit", middleware.RequireSelfOrRole(model.RoleAdmin, user.Edit))
	mux.HandleFunc("POST /users/{id}", middleware.RequireSelfOrRole(model.RoleAdmin, user.Update))
	mux.HandleFunc("POST /users/{id}/delete", middleware.RequireSelfOrRole(model.RoleAdmin, user.Destroy))
	mux.HandleFunc("POST /users/{id}/role", middleware.RequireRole(model.RoleAdmin, user.SetRole))

	// Announcements as HTML pages: the interface to manage them.
	announcement := controller.AnnouncementController{}
	mux.HandleFunc("GET /announcements", announcement.Index)
	mux.HandleFunc("GET /announcements/new", middleware.RequireRole(model.RoleAdmin, announcement.Create))
	mux.HandleFunc("POST /announcements", middleware.RequireRole(model.RoleAdmin, announcement.Store))
	mux.HandleFunc("GET /announcements/{id}", announcement.Show)
	mux.HandleFunc("GET /announcements/{id}/edit", middleware.RequireRole(model.RoleAdmin, announcement.Edit))
	mux.HandleFunc("POST /announcements/{id}", middleware.RequireRole(model.RoleAdmin, announcement.Update))
	mux.HandleFunc("POST /announcements/{id}/delete", middleware.RequireRole(model.RoleAdmin, announcement.Destroy))

	// Graduates (the alumni table) as HTML pages. The JSON is at /alumni.
	graduate := controller.GraduateController{}
	mux.HandleFunc("GET /graduates", middleware.RequireRole(model.RoleGraduate, graduate.Index))
	mux.HandleFunc("GET /graduates/new", middleware.RequireRole(model.RoleAdmin, graduate.Create))
	mux.HandleFunc("POST /graduates", middleware.RequireRole(model.RoleAdmin, graduate.Store))
	mux.HandleFunc("GET /graduates/{id}", middleware.RequireRole(model.RoleGraduate, graduate.Show))
	mux.HandleFunc("GET /graduates/{id}/edit", middleware.RequireRole(model.RoleAdmin, graduate.Edit))
	mux.HandleFunc("POST /graduates/{id}", middleware.RequireRole(model.RoleAdmin, graduate.Update))
	mux.HandleFunc("POST /graduates/{id}/delete", middleware.RequireRole(model.RoleAdmin, graduate.Destroy))

	// Job postings as HTML pages. A posting's page lists its applications.
	jobPosting := controller.JobPostingController{}
	mux.HandleFunc("GET /job-postings", middleware.RequireRole(model.RoleGraduate, jobPosting.Index))
	mux.HandleFunc("GET /job-postings/new", middleware.RequireRole(model.RoleGraduate, jobPosting.Create))
	mux.HandleFunc("POST /job-postings", middleware.RequireRole(model.RoleGraduate, jobPosting.Store))
	mux.HandleFunc("GET /job-postings/{id}", middleware.RequireRole(model.RoleGraduate, jobPosting.Show))
	mux.HandleFunc("GET /job-postings/{id}/edit", middleware.RequireRole(model.RoleGraduate, jobPosting.Edit))
	mux.HandleFunc("POST /job-postings/{id}", middleware.RequireRole(model.RoleGraduate, jobPosting.Update))
	mux.HandleFunc("POST /job-postings/{id}/delete", middleware.RequireRole(model.RoleGraduate, jobPosting.Destroy))

	// Applications as HTML. Applying is nested under the posting; the status
	// form and delete work on the application's own id.
	application := controller.ApplicationController{}
	mux.HandleFunc("GET /applications", middleware.RequireRole(model.RoleGraduate, application.Index))
	mux.HandleFunc("POST /job-postings/{id}/applications", middleware.RequireRole(model.RoleGraduate, application.Store))
	mux.HandleFunc("POST /applications/{id}", middleware.RequireRole(model.RoleGraduate, application.Update))
	mux.HandleFunc("POST /applications/{id}/delete", middleware.RequireRole(model.RoleGraduate, application.Destroy))
}
