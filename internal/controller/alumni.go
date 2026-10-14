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
		http.Error(w, "database error", http.StatusInternalServerError)
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
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	view.JSON(w, http.StatusCreated, created)
}
