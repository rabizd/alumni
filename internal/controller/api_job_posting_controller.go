package controller

import (
	"errors"
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// ApiJobPostingController answers job postings with JSON. A posting is reached
// directly at /api/job-postings/{id}, or through the user who posted it at
// /api/users/{id}/job-postings (nested). Only listing and creating are nested:
// once a posting has its own id, that id is enough to find it.
type ApiJobPostingController struct{}

// Index: GET /api/job-postings -> every posting, newest first
func (ApiJobPostingController) Index(w http.ResponseWriter, r *http.Request) {
	list, err := model.JobPostings.List()
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusOK, list)
}

// Show: GET /api/job-postings/{id} -> one posting, or 404
func (ApiJobPostingController) Show(w http.ResponseWriter, r *http.Request) {
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

	view.JSON(w, http.StatusOK, p)
}

// Store: POST /api/job-postings -> create one for the userId in the body; 201
func (ApiJobPostingController) Store(w http.ResponseWriter, r *http.Request) {
	var p model.JobPosting
	if !decodeJSON(w, r, &p) {
		return
	}

	if err := p.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := model.JobPostings.Add(p)
	if errors.Is(err, model.ErrUserNotFound) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusCreated, created)
}

// Update: PUT /api/job-postings/{id} -> replace every field but id and createdAt
func (ApiJobPostingController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var p model.JobPosting
	if !decodeJSON(w, r, &p) {
		return
	}

	if err := p.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, found, err := model.JobPostings.Replace(id, p)
	if errors.Is(err, model.ErrUserNotFound) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, updated)
}

// Destroy: DELETE /api/job-postings/{id} -> 204 with an empty body, or 404
func (ApiJobPostingController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	found, err := model.JobPostings.Remove(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UserIndex: GET /api/users/{id}/job-postings -> that user's postings, or 404 if there is no such user
func (ApiJobPostingController) UserIndex(w http.ResponseWriter, r *http.Request) {
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

	list, err := u.JobPostings()
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusOK, list)
}

// UserStore: POST /api/users/{id}/job-postings -> create one for the user in
// the URL; a userId in the body is overwritten by it
func (ApiJobPostingController) UserStore(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var p model.JobPosting
	if !decodeJSON(w, r, &p) {
		return
	}
	p.UserID = id

	if err := p.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := model.JobPostings.Add(p)
	if errors.Is(err, model.ErrUserNotFound) {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusCreated, created)
}
