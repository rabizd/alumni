package middleware

import (
	"net/http"
	"strconv"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// forbiddenPage is what forbidden.html needs.
type forbiddenPage struct {
	Need, Have string
}

// RequireRole lets a request through only when someone is signed in with at
// least the given role (see the role table in DDD.md). A guest is sent to the
// sign-in page; a signed-in user whose role is too low gets 403 Forbidden.
// It runs after Session, which has already put the user in the context.
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := model.CurrentUser(r.Context())
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if !u.HasRole(role) {
			view.HTML(w, r, http.StatusForbidden, "forbidden.html", forbiddenPage{Need: role, Have: u.Role})
			return
		}
		next(w, r)
	}
}

// RequireSelfOrRole is RequireRole for /users/{id} pages, with one exception:
// a user may always reach their own page, whatever their role.
func RequireSelfOrRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u, ok := model.CurrentUser(r.Context()); ok && strconv.Itoa(u.ID) == r.PathValue("id") {
			next(w, r)
			return
		}
		RequireRole(role, next)(w, r)
	}
}
