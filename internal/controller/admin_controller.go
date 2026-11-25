package controller

import (
	"errors"
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// AdminController serves /admin, the panel only admins can open (the route
// is guarded by middleware.RequireRole). The panel's forms post to the
// existing routes with back=/admin, so after a change the admin lands here again.
type AdminController struct{}

// adminPage is what admin.html needs.
type adminPage struct {
	RoleCounts    map[string]int
	Users         []model.User
	Postings      []model.JobPosting
	Announcements []model.Announcement
	Applications  int
	Roles         []string
}

// Index: GET /admin -> counts, every user with a role form, postings and announcements
func (AdminController) Index(w http.ResponseWriter, r *http.Request) {
	users, err1 := model.Users.List()
	postings, err2 := model.JobPostings.List()
	announcements, err3 := model.Announcements.List()
	applications, err4 := model.Applications.List()
	if err := errors.Join(err1, err2, err3, err4); err != nil {
		dbError(w, err)
		return
	}

	counts := map[string]int{}
	for _, u := range users {
		counts[u.Role]++
	}
	view.HTML(w, r, http.StatusOK, "admin.html", adminPage{
		RoleCounts: counts, Users: users, Postings: postings, Announcements: announcements,
		Applications: len(applications), Roles: model.Roles,
	})
}
