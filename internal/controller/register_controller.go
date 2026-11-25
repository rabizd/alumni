package controller

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// RegisterController serves /register, the sign-up page.
type RegisterController struct{}

// registerForm is what register.html needs. The password is never sent back
// to the page, so after an error it has to be typed again.
type registerForm struct {
	User          model.User
	TermsAccepted bool
	Error         string
}

const minPasswordLength = 8

// Show: GET /register -> an empty sign-up form
func (RegisterController) Show(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, r, http.StatusOK, "register.html", registerForm{})
}

// Store: POST /register -> check the form, create the account, go to its page
func (RegisterController) Store(w http.ResponseWriter, r *http.Request) {
	form := registerForm{User: userFromForm(r), TermsAccepted: r.FormValue("terms") == "on"}
	password := r.FormValue("password")

	fail := func(status int, msg string) {
		form.Error = msg
		view.HTML(w, r, status, "register.html", form)
	}

	if err := form.User.Validate(); err != nil {
		fail(http.StatusBadRequest, err.Error())
		return
	}
	switch {
	case utf8.RuneCountInString(password) < minPasswordLength:
		fail(http.StatusBadRequest, fmt.Sprintf("password must be at least %d characters", minPasswordLength))
		return
	case password != r.FormValue("passwordConfirm"):
		fail(http.StatusBadRequest, "the two passwords do not match")
		return
	case !form.TermsAccepted:
		fail(http.StatusBadRequest, "accept the terms of use to register")
		return
	}

	taken, err := model.Users.EmailTaken(form.User.Email)
	if err != nil {
		dbError(w, err)
		return
	}
	if taken {
		fail(http.StatusConflict, "an account with this email already exists")
		return
	}

	// bcrypt adds a random salt and is slow on purpose, so a stolen hash is
	// hard to turn back into the password.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fail(http.StatusBadRequest, "password is too long")
		return
	}

	u, err := model.Users.Register(form.User, string(hash))
	if err != nil {
		dbError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/users/%d", u.ID), http.StatusSeeOther)
}
