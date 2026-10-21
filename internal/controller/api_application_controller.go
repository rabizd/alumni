package controller

import (
	"errors"
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// ApiApplicationController answers applications with JSON. An application is
// reached directly at /api/applications/{id}, or nested through the posting it
// was made to (/api/job-postings/{id}/applications) or the user who made it
// (/api/users/{id}/applications).
type ApiApplicationController struct{}

// Index: GET /api/applications -> every application, newest first
func (ApiApplicationController) Index(w http.ResponseWriter, r *http.Request) {
	list, err := model.Applications.List()
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusOK, list)
}

// Show: GET /api/applications/{id} -> one application, or 404
func (ApiApplicationController) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a, found, err := model.Applications.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no application with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, a)
}

// Store: POST /api/applications -> create one; 201
func (ApiApplicationController) Store(w http.ResponseWriter, r *http.Request) {
	var a model.Application
	if !decodeJSON(w, r, &a) {
		return
	}

	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := model.Applications.Add(a)
	if err != nil {
		applicationError(w, err)
		return
	}

	view.JSON(w, http.StatusCreated, created)
}

// Update: PUT /api/applications/{id} -> replace every field but id and createdAt
func (ApiApplicationController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var a model.Application
	if !decodeJSON(w, r, &a) {
		return
	}

	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, found, err := model.Applications.Replace(id, a)
	if err != nil {
		applicationError(w, err)
		return
	}
	if !found {
		http.Error(w, "no application with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, updated)
}

// Destroy: DELETE /api/applications/{id} -> 204 with an empty body, or 404
func (ApiApplicationController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	found, err := model.Applications.Remove(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no application with that id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// JobPostingIndex: GET /api/job-postings/{id}/applications -> the applications
// that posting received, or 404 if there is no such posting
func (ApiApplicationController) JobPostingIndex(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	p, found, err := model.JobPostings.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}

	list, err := p.Applications()
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusOK, list)
}

// JobPostingStore: POST /api/job-postings/{id}/applications -> apply to the
// posting in the URL; the applicant is the userId in the body
func (ApiApplicationController) JobPostingStore(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var a model.Application
	if !decodeJSON(w, r, &a) {
		return
	}
	a.JobPostingID = id

	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := model.Applications.Add(a)
	if errors.Is(err, model.ErrJobPostingNotFound) {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}
	if err != nil {
		applicationError(w, err)
		return
	}

	view.JSON(w, http.StatusCreated, created)
}

// UserIndex: GET /api/users/{id}/applications -> the applications that user
// made, or 404 if there is no such user
func (ApiApplicationController) UserIndex(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	u, found, err := model.Users.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	list, err := u.Applications()
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusOK, list)
}

// applicationError answers a failed Add or Replace: a missing posting or user
// in the body is 400, a second application to the same posting is 409, and
// anything else is a database error.
func applicationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrJobPostingNotFound), errors.Is(err, model.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, model.ErrAlreadyApplied):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		dbError(w, err)
	}
}
