package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// UserController serves /users as HTML pages, for people using a browser.
// HTML forms can only send GET and POST, so update and delete are POSTs here,
// where ApiUserController uses PUT, PATCH and DELETE.
type UserController struct{}

// userForm is what user_form.html needs to draw the create and edit forms.
type userForm struct {
	Title  string
	Action string
	User   model.User
	Error  string
}

// userList is what users.html needs: the users, and an empty user for the
// create form under the list.
type userList struct {
	Users []model.User
	New   model.User
}

// Index: GET /users -> the list page
func (UserController) Index(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "users.html", userList{Users: model.Users.List()})
}

// Show: GET /users/{id} -> one user's page, or 404
func (UserController) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	u, found := model.Users.Find(id)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	view.HTML(w, http.StatusOK, "user.html", u)
}

// Create: GET /users/new -> an empty form
func (UserController) Create(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "user_form.html", userForm{Title: "Yeni kullanıcı", Action: "/users"})
}

// Store: POST /users -> create from the form, then go to the list
func (UserController) Store(w http.ResponseWriter, r *http.Request) {
	u := userFromForm(r)
	if err := u.Validate(); err != nil {
		// Show the form again with what was typed, so nothing has to be re-entered.
		view.HTML(w, http.StatusBadRequest, "user_form.html",
			userForm{Title: "Yeni kullanıcı", Action: "/users", User: u, Error: err.Error()})
		return
	}

	model.Users.Add(u)
	// 303 See Other: after a POST, send the browser to a page it can safely reload.
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// Edit: GET /users/{id}/edit -> the form, filled in with the current values
func (UserController) Edit(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	u, found := model.Users.Find(id)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	view.HTML(w, http.StatusOK, "user_form.html",
		userForm{Title: "Kullanıcıyı düzenle", Action: fmt.Sprintf("/users/%d", id), User: u})
}

// Update: POST /users/{id} -> save the edit form, then go to the user's page
func (UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	u := userFromForm(r)
	u.ID = id
	if err := u.Validate(); err != nil {
		view.HTML(w, http.StatusBadRequest, "user_form.html",
			userForm{Title: "Kullanıcıyı düzenle", Action: fmt.Sprintf("/users/%d", id), User: u, Error: err.Error()})
		return
	}

	if _, found := model.Users.Replace(id, u); !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", id), http.StatusSeeOther)
}

// Destroy: POST /users/{id}/delete -> remove the user, then go to the list
func (UserController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if !model.Users.Remove(id) {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// userFromForm reads the create and edit forms. A year that is empty or not a
// number becomes 0, which User.Validate reports as missing. An unticked
// checkbox is not sent at all, so PrepExempt is true only when it is ticked.
func userFromForm(r *http.Request) model.User {
	startYear, _ := strconv.Atoi(r.FormValue("startYear"))
	graduationYear, _ := strconv.Atoi(r.FormValue("graduationYear"))
	return model.User{
		Name:           r.FormValue("name"),
		Email:          r.FormValue("email"),
		Department:     r.FormValue("department"),
		StartYear:      startYear,
		GraduationYear: graduationYear,
		DoubleMajor:    r.FormValue("doubleMajor"),
		Minor:          r.FormValue("minor"),
		Advisor:        r.FormValue("advisor"),
		PrepExempt:     r.FormValue("prepExempt") == "on",
	}
}
