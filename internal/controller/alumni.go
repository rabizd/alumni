package controller

import (
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// GET /alumni -> every graduate, as a JSON array
func ListAlumni(w http.ResponseWriter, r *http.Request) {
	list, err := model.Graduates.List()
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusOK, list)
}

// POST /alumni -> create one graduate from the JSON in the request body
func CreateAlumni(w http.ResponseWriter, r *http.Request) {
	var a model.Alumni
	if !decodeJSON(w, r, &a) {
		return
	}

	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := model.Graduates.Add(a)
	if err != nil {
		dbError(w, err)
		return
	}

	view.JSON(w, http.StatusCreated, created)
}

// GET /alumni/{id} -> one graduate, or 404
func ShowAlumni(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a, found, err := model.Graduates.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no graduate with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, a)
}

// PUT /alumni/{id} -> replace the graduate's name, department and graduation year
func UpdateAlumni(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var a model.Alumni
	if !decodeJSON(w, r, &a) {
		return
	}

	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, found, err := model.Graduates.Replace(id, a)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no graduate with that id", http.StatusNotFound)
		return
	}

	view.JSON(w, http.StatusOK, updated)
}

// DELETE /alumni/{id} -> 204 with an empty body, or 404
func DeleteAlumni(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	found, err := model.Graduates.Remove(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no graduate with that id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
