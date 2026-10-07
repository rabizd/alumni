package routes

import (
	"net/http"

	"github.com/rabizd/alumni/internal/controller"
)

// API registers the JSON routes: /api/users, /api/announcements, the health
// check, the Swagger documentation, and /alumni.
func API(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", controller.Health)
	mux.HandleFunc("GET /api/swagger", controller.Swagger)
	mux.HandleFunc("GET /api/swagger.json", controller.SwaggerSpec)

	// Users as JSON.
	apiUser := controller.ApiUserController{}
	mux.HandleFunc("GET /api/users", apiUser.Index)
	mux.HandleFunc("POST /api/users", apiUser.Store)
	mux.HandleFunc("GET /api/users/{id}", apiUser.Show)
	mux.HandleFunc("PUT /api/users/{id}", apiUser.Update)
	mux.HandleFunc("PATCH /api/users/{id}", apiUser.Patch)
	mux.HandleFunc("DELETE /api/users/{id}", apiUser.Destroy)

	// Announcements as JSON.
	apiAnnouncement := controller.ApiAnnouncementController{}
	mux.HandleFunc("GET /api/announcements", apiAnnouncement.Index)
	mux.HandleFunc("POST /api/announcements", apiAnnouncement.Store)
	mux.HandleFunc("GET /api/announcements/{id}", apiAnnouncement.Show)
	mux.HandleFunc("PUT /api/announcements/{id}", apiAnnouncement.Update)
	mux.HandleFunc("DELETE /api/announcements/{id}", apiAnnouncement.Destroy)

	mux.HandleFunc("GET /alumni", controller.ListAlumni)
	mux.HandleFunc("POST /alumni", controller.CreateAlumni)
}
