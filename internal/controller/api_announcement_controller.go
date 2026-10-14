package controller

import (
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// ApiAnnouncementController answers /api/announcements with JSON, for
// programs. AnnouncementController serves the same announcements as HTML pages.
type ApiAnnouncementController struct{}

// Index: GET /api/announcements -> every announcement, newest first
func (ApiAnnouncementController) Index(w http.ResponseWriter, r *http.Request) {
	list, err := model.Announcements.List()
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusOK, list)
}

// Show: GET /api/announcements/{id} -> one announcement, or 404
func (ApiAnnouncementController) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a, found, err := model.Announcements.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no announcement with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, a)
}

// Store: POST /api/announcements -> create one; 201 with the id and time the server assigned
func (ApiAnnouncementController) Store(w http.ResponseWriter, r *http.Request) {
	var a model.Announcement
	if !decodeJSON(w, r, &a) {
		return
	}

	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := model.Announcements.Add(a)
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusCreated, created)
}

// Update: PUT /api/announcements/{id} -> replace title, body and author
func (ApiAnnouncementController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var a model.Announcement
	if !decodeJSON(w, r, &a) {
		return
	}

	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, found, err := model.Announcements.Replace(id, a)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no announcement with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, updated)
}

// Destroy: DELETE /api/announcements/{id} -> 204 with an empty body, or 404
func (ApiAnnouncementController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	found, err := model.Announcements.Remove(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no announcement with that id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
