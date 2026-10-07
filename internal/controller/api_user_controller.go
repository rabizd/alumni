// Package controller is the C in MVC: one function per route. It reads the
// request, checks it, asks the model for data, and hands the result to a view.
package controller

import (
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// ApiUserController answers /api/users with JSON, for programs rather than
// people. UserController serves the same users as HTML pages.
type ApiUserController struct{}

// Index: GET /api/users -> every user, as a JSON array
func (ApiUserController) Index(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, model.Users.List())
}

// Show: GET /api/users/{id} -> one user, or 404
func (ApiUserController) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	u, found := model.Users.Find(id)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, u)
}

// Store: POST /api/users -> create a user; 201 with the id the server assigned
func (ApiUserController) Store(w http.ResponseWriter, r *http.Request) {
	var u model.User
	if !decodeJSON(w, r, &u) {
		return
	}

	if err := u.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	view.JSON(w, http.StatusCreated, model.Users.Add(u))
}

// Update: PUT /api/users/{id} -> replace the whole user, so every field is required
func (ApiUserController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var u model.User
	if !decodeJSON(w, r, &u) {
		return
	}

	if err := u.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, found := model.Users.Replace(id, u)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, updated)
}

// Patch: PATCH /api/users/{id} -> change only the fields the body mentions
func (ApiUserController) Patch(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var p model.UserPatch
	if !decodeJSON(w, r, &p) {
		return
	}

	if err := p.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, found := model.Users.Patch(id, p)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, updated)
}

// Destroy: DELETE /api/users/{id} -> 204 with an empty body, or 404
func (ApiUserController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if !model.Users.Remove(id) {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	// 204 No Content: it worked, and there is nothing left to send back.
	w.WriteHeader(http.StatusNoContent)
}
