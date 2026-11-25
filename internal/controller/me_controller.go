package controller

import (
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// MeController serves /me, the signed-in user's own page. Everything on it
// is read with the user's id from the session, never from the URL, so
// nobody can open another person's /me.
type MeController struct{}

// mePage is what me.html needs.
type mePage struct {
	User     model.User
	Postings []model.JobPosting
	// Sent are the user's own applications; Received are the applications
	// to the user's postings.
	Sent, Received []applicationRow
	Statuses       []string
}

// Show: GET /me -> profile, own postings, sent and received applications
func (MeController) Show(w http.ResponseWriter, r *http.Request) {
	me, _ := model.CurrentUser(r.Context()) // the route needs sign-in, so me is set

	postings, err := model.JobPostings.ListByUser(me.ID)
	if err != nil {
		dbError(w, err)
		return
	}
	sent, err := model.Applications.ListByUser(me.ID)
	if err != nil {
		dbError(w, err)
		return
	}
	allPostings, err := model.JobPostings.List()
	if err != nil {
		dbError(w, err)
		return
	}
	names, _, err := userNames()
	if err != nil {
		dbError(w, err)
		return
	}
	titles := make(map[int]string, len(allPostings))
	for _, p := range allPostings {
		titles[p.ID] = p.Title
	}

	page := mePage{User: me, Postings: postings, Statuses: model.ApplicationStatuses}
	for _, a := range sent {
		page.Sent = append(page.Sent, applicationRow{Application: a, Applicant: me.Name, PostingTitle: titles[a.JobPostingID]})
	}
	for _, p := range postings {
		received, err := model.Applications.ListByJobPosting(p.ID)
		if err != nil {
			dbError(w, err)
			return
		}
		for _, a := range received {
			page.Received = append(page.Received, applicationRow{Application: a, Applicant: names[a.UserID], PostingTitle: p.Title})
		}
	}
	view.HTML(w, r, http.StatusOK, "me.html", page)
}
