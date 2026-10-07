// Package controller is the C in MVC: one function per route. It reads the
// request, checks it, asks the model for data, and hands the result to a view.
package controller

import (
	"net/http"
	"strings"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// GET /api/users -> every user, as a JSON array
func ListUsers(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, model.Users.List())
}

// GET /api/users/{id} -> one user, or 404
func GetUser(w http.ResponseWriter, r *http.Request) {
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

// DELETE /api/users/{id} -> 204 with an empty body, or 404
func DeleteUser(w http.ResponseWriter, r *http.Request) {
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

// POST /api/users -> create a user; 201 with the id the server assigned
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	if !decodeJSON(w, r, &u) {
		return
	}

	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	if err := u.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	view.JSON(w, http.StatusCreated, model.Users.Add(u))
}

// PUT /api/users/{id} -> replace the whole user, so every field is required
func ReplaceUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var u model.User
	if !decodeJSON(w, r, &u) {
		return
	}

	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
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

// PATCH /api/users/{id} -> change only the fields the body mentions
func PatchUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var p model.UserPatch
	if !decodeJSON(w, r, &p) {
		return
	}

	// An empty body would silently do nothing, which is more likely a mistake
	// than an intention.
	if p.Name == nil && p.Email == nil {
		http.Error(w, "send at least one of name or email", http.StatusBadRequest)
		return
	}

	if p.Name != nil {
		trimmed := strings.TrimSpace(*p.Name)
		if trimmed == "" {
			http.Error(w, "name cannot be empty", http.StatusBadRequest)
			return
		}
		p.Name = &trimmed
	}
	if p.Email != nil {
		trimmed := strings.TrimSpace(*p.Email)
		if trimmed == "" {
			http.Error(w, "email cannot be empty", http.StatusBadRequest)
			return
		}
		p.Email = &trimmed
	}

	updated, found := model.Users.Patch(id, p)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, updated)
}
