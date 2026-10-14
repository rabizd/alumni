package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// GraduateController serves /graduates as HTML pages: the interface for
// managing the alumni table in a browser. The JSON for the same table is at
// /alumni. Like the other HTML controllers, it uses POST for update and delete
// because HTML forms can only send GET and POST.
type GraduateController struct{}

// graduateList is what graduates.html needs: the graduates, and an empty one
// for the create form under the list.
type graduateList struct {
	Graduates []model.Alumni
	New       model.Alumni
}

// graduateForm is what graduate_form.html needs.
type graduateForm struct {
	Title    string
	Action   string
	Graduate model.Alumni
	Error    string
}

// Index: GET /graduates -> the list page, with a create form
func (GraduateController) Index(w http.ResponseWriter, r *http.Request) {
	list, err := model.Graduates.List()
	if err != nil {
		dbError(w, err)
		return
	}

	view.HTML(w, http.StatusOK, "graduates.html", graduateList{Graduates: list})
}

// Show: GET /graduates/{id} -> one graduate's page, or 404
func (GraduateController) Show(w http.ResponseWriter, r *http.Request) {
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

	view.HTML(w, http.StatusOK, "graduate.html", a)
}

// Create: GET /graduates/new -> an empty form
func (GraduateController) Create(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "graduate_form.html",
		graduateForm{Title: "Yeni mezun", Action: "/graduates"})
}

// Store: POST /graduates -> create from the form, then go to the list
func (GraduateController) Store(w http.ResponseWriter, r *http.Request) {
	a := graduateFromForm(r)
	if err := a.Validate(); err != nil {
		view.HTML(w, http.StatusBadRequest, "graduate_form.html",
			graduateForm{Title: "Yeni mezun", Action: "/graduates", Graduate: a, Error: err.Error()})
		return
	}

	if _, err := model.Graduates.Add(a); err != nil {
		dbError(w, err)
		return
	}
	http.Redirect(w, r, "/graduates", http.StatusSeeOther)
}

// Edit: GET /graduates/{id}/edit -> the form, filled in with the current values
func (GraduateController) Edit(w http.ResponseWriter, r *http.Request) {
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

	view.HTML(w, http.StatusOK, "graduate_form.html",
		graduateForm{Title: "Mezunu düzenle", Action: fmt.Sprintf("/graduates/%d", id), Graduate: a})
}

// Update: POST /graduates/{id} -> save the edit form, then go to the graduate
func (GraduateController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a := graduateFromForm(r)
	a.ID = id
	if err := a.Validate(); err != nil {
		view.HTML(w, http.StatusBadRequest, "graduate_form.html",
			graduateForm{Title: "Mezunu düzenle", Action: fmt.Sprintf("/graduates/%d", id), Graduate: a, Error: err.Error()})
		return
	}

	_, found, err := model.Graduates.Replace(id, a)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no graduate with that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/graduates/%d", id), http.StatusSeeOther)
}

// Destroy: POST /graduates/{id}/delete -> remove it, then go to the list
func (GraduateController) Destroy(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "/graduates", http.StatusSeeOther)
}

func graduateFromForm(r *http.Request) model.Alumni {
	graduationYear, _ := strconv.Atoi(r.FormValue("graduationYear"))
	return model.Alumni{
		Name:           r.FormValue("name"),
		Department:     r.FormValue("department"),
		GraduationYear: graduationYear,
	}
}
