package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/rabizd/alumni/internal/middleware"
	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// AuthController serves /login and /logout.
type AuthController struct{}

// loginForm is what login.html needs. The password is never sent back.
type loginForm struct {
	Email string
	Error string
}

// ShowLogin: GET /login -> the sign-in form
func (AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, r, http.StatusOK, "login.html", loginForm{})
}

// Login: POST /login -> check email and password, start a session, go home
func (AuthController) Login(w http.ResponseWriter, r *http.Request) {
	form := loginForm{Email: r.FormValue("email")}

	u, hash, found, err := model.Users.FindForLogin(form.Email)
	if err != nil {
		dbError(w, err)
		return
	}
	// The same answer for "no such email" and "wrong password", so the form
	// does not tell a stranger which emails have an account. Users made
	// before /register existed have no password ("") and cannot sign in.
	if !found || hash == "" || bcrypt.CompareHashAndPassword([]byte(hash), []byte(r.FormValue("password"))) != nil {
		form.Error = "email or password is wrong"
		view.HTML(w, r, http.StatusUnauthorized, "login.html", form)
		return
	}

	token, err := model.Sessions.Create(u.ID)
	if err != nil {
		dbError(w, err)
		return
	}
	// HttpOnly: page scripts cannot read the cookie. SameSite=Lax: other
	// sites cannot send it along with a form they post here.
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(model.SessionLength.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

// Logout: POST /logout -> end the session and forget the cookie, go home
func (AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(middleware.SessionCookie); err == nil {
		if err := model.Sessions.Delete(c.Value); err != nil {
			dbError(w, err)
			return
		}
	}
	// MaxAge -1 tells the browser to delete the cookie now.
	http.SetCookie(w, &http.Cookie{Name: middleware.SessionCookie, Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
