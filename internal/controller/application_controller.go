package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// ApplicationController serves applications as HTML. Applying happens on a
// posting's page (POST /job-postings/{id}/applications); /applications lists
// every application, where its status can be changed or it can be deleted.
type ApplicationController struct{}

// applicationRow is an application with the names the pages show.
type applicationRow struct {
	model.Application
	Applicant    string
	PostingTitle string
}

// applicationList is what applications.html needs.
type applicationList struct {
	Applications []applicationRow
	Statuses     []string
}

// Index: GET /applications -> every application, newest first
func (ApplicationController) Index(w http.ResponseWriter, r *http.Request) {
	apps, err := model.Applications.List()
	if err != nil {
		dbError(w, err)
		return
	}
	names, _, err := userNames()
	if err != nil {
		dbError(w, err)
		return
	}
	postings, err := model.JobPostings.List()
	if err != nil {
		dbError(w, err)
		return
	}
	titles := make(map[int]string, len(postings))
	for _, p := range postings {
		titles[p.ID] = p.Title + " · " + p.Company
	}

	rows := make([]applicationRow, len(apps))
	for i, a := range apps {
		rows[i] = applicationRow{Application: a, Applicant: names[a.UserID], PostingTitle: titles[a.JobPostingID]}
	}
	view.HTML(w, http.StatusOK, "applications.html", applicationList{Applications: rows, Statuses: model.ApplicationStatuses})
}

// Store: POST /job-postings/{id}/applications -> apply to that posting, then
// go back to it; on an error its page is shown again with the message
func (ApplicationController) Store(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	userID, _ := strconv.Atoi(r.FormValue("userId"))
	a := model.Application{JobPostingID: id, UserID: userID, CoverLetter: r.FormValue("coverLetter")}
	if err := a.Validate(); err != nil {
		JobPostingController{}.renderPage(w, id, http.StatusBadRequest, a, err.Error())
		return
	}

	_, err := model.Applications.Add(a)
	switch {
	case errors.Is(err, model.ErrJobPostingNotFound):
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	case errors.Is(err, model.ErrAlreadyApplied):
		JobPostingController{}.renderPage(w, id, http.StatusConflict, a, "Bu kullanıcı bu ilana zaten başvurmuş.")
		return
	case errors.Is(err, model.ErrUserNotFound):
		JobPostingController{}.renderPage(w, id, http.StatusBadRequest, a, err.Error())
		return
	case err != nil:
		dbError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/job-postings/%d", id), http.StatusSeeOther)
}

// Update: POST /applications/{id} -> change only the status, then go back to
// the page the form was on (the "back" field), or to /applications
func (ApplicationController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a, found, err := model.Applications.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no application with that id", http.StatusNotFound)
		return
	}

	a.Status = r.FormValue("status")
	if err := a.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, _, err := model.Applications.Replace(id, a); err != nil {
		dbError(w, err)
		return
	}
	http.Redirect(w, r, backTo(r), http.StatusSeeOther)
}

// Destroy: POST /applications/{id}/delete -> remove it, then go back
func (ApplicationController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	found, err := model.Applications.Remove(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no application with that id", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, backTo(r), http.StatusSeeOther)
}

// backTo is where to go after a status change or delete: the page the form
// was on, sent in the "back" field. Only a local path is accepted, so the
// field cannot send the browser to another site.
func backTo(r *http.Request) string {
	back := r.FormValue("back")
	if len(back) > 1 && back[0] == '/' && back[1] != '/' && back[1] != '\\' {
		return back
	}
	return "/applications"
}
