package controller

import (
	"net/http"
	"strings"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// GET /alumni -> every graduate, as a JSON array
func ListAlumni(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, model.Graduates.List())
}

// POST /alumni -> create one graduate from the JSON in the request body
func CreateAlumni(w http.ResponseWriter, r *http.Request) {
	var a model.Alumni
	if !decodeJSON(w, r, &a) {
		return
	}

	a.Name = strings.TrimSpace(a.Name)
	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	view.JSON(w, http.StatusCreated, model.Graduates.Add(a))
}
