package routes

import (
	"net/http"

	"github.com/rabizd/alumni/internal/controller"
)

// API registers the JSON routes: /api/users, /api/announcements,
// /api/job-postings, /api/applications, the health check, the Swagger
// documentation, and /alumni.
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

	// Job postings as JSON. Direct routes reach a posting by its own id...
	apiJobPosting := controller.ApiJobPostingController{}
	mux.HandleFunc("GET /api/job-postings", apiJobPosting.Index)
	mux.HandleFunc("POST /api/job-postings", apiJobPosting.Store)
	mux.HandleFunc("GET /api/job-postings/{id}", apiJobPosting.Show)
	mux.HandleFunc("PUT /api/job-postings/{id}", apiJobPosting.Update)
	mux.HandleFunc("DELETE /api/job-postings/{id}", apiJobPosting.Destroy)
	// ...nested routes go through the user who posted it ({id} is the user's id).
	mux.HandleFunc("GET /api/users/{id}/job-postings", apiJobPosting.UserIndex)
	mux.HandleFunc("POST /api/users/{id}/job-postings", apiJobPosting.UserStore)

	// Applications as JSON. Direct routes reach an application by its own id...
	apiApplication := controller.ApiApplicationController{}
	mux.HandleFunc("GET /api/applications", apiApplication.Index)
	mux.HandleFunc("POST /api/applications", apiApplication.Store)
	mux.HandleFunc("GET /api/applications/{id}", apiApplication.Show)
	mux.HandleFunc("PUT /api/applications/{id}", apiApplication.Update)
	mux.HandleFunc("DELETE /api/applications/{id}", apiApplication.Destroy)
	// ...nested routes go through the posting ({id} is the posting's id) or the user ({id} is the user's id).
	mux.HandleFunc("GET /api/job-postings/{id}/applications", apiApplication.JobPostingIndex)
	mux.HandleFunc("POST /api/job-postings/{id}/applications", apiApplication.JobPostingStore)
	mux.HandleFunc("GET /api/users/{id}/applications", apiApplication.UserIndex)

	mux.HandleFunc("GET /alumni", controller.ListAlumni)
	mux.HandleFunc("POST /alumni", controller.CreateAlumni)
	mux.HandleFunc("GET /alumni/{id}", controller.ShowAlumni)
	mux.HandleFunc("PUT /alumni/{id}", controller.UpdateAlumni)
	mux.HandleFunc("DELETE /alumni/{id}", controller.DeleteAlumni)
}
