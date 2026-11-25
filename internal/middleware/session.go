package middleware

import (
	"log"
	"net/http"

	"github.com/rabizd/alumni/internal/model"
)

// SessionCookie is the name of the cookie that holds the sign-in token.
const SessionCookie = "session"

// Session finds out who is signed in: when the request carries a valid
// session cookie, the user is put in the request's context, where
// model.CurrentUser finds them. Without one the request goes on as a guest.
func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(SessionCookie); err == nil {
			u, found, err := model.Sessions.User(c.Value)
			if err != nil {
				log.Printf("reading session: %v", err)
			} else if found {
				r = r.WithContext(model.WithCurrentUser(r.Context(), u))
			}
		}
		next.ServeHTTP(w, r)
	})
}
