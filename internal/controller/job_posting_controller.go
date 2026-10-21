package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// JobPostingController serves /job-postings as HTML pages. A posting's page
// also lists the applications it received and has a form to apply, so the
// relation between postings and applications is visible in the browser.
type JobPostingController struct{}

// jobPostingRow is a posting with the name of the user who posted it.
type jobPostingRow struct {
	model.JobPosting
	Poster string
}

// jobPostingPage is what job_posting.html needs.
type jobPostingPage struct {
	Posting      model.JobPosting
	Poster       model.User
	Applications []applicationRow
	Users        []model.User
	Statuses     []string
	Apply        model.Application
	Error        string
}

// jobPostingForm is what job_posting_form.html needs.
type jobPostingForm struct {
	Title   string
	Action  string
	Posting model.JobPosting
	Users   []model.User
	Types   []string
	Error   string
}

// Index: GET /job-postings -> the list page
func (JobPostingController) Index(w http.ResponseWriter, r *http.Request) {
	list, err := model.JobPostings.List()
	if err != nil {
		dbError(w, err)
		return
	}
	names, _, err := userNames()
	if err != nil {
		dbError(w, err)
		return
	}

	rows := make([]jobPostingRow, len(list))
	for i, p := range list {
		rows[i] = jobPostingRow{JobPosting: p, Poster: names[p.UserID]}
	}
	view.HTML(w, http.StatusOK, "job_postings.html", rows)
}

// Show: GET /job-postings/{id} -> one posting with its applications and an apply form, or 404
func (c JobPostingController) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	c.renderPage(w, id, http.StatusOK, model.Application{}, "")
}

// renderPage shows a posting's page. Apply and errMsg refill the apply form
// when an application could not be saved.
func (JobPostingController) renderPage(w http.ResponseWriter, id, status int, apply model.Application, errMsg string) {
	p, found, err := model.JobPostings.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}
	names, users, err := userNames()
	if err != nil {
		dbError(w, err)
		return
	}
	apps, err := p.Applications()
	if err != nil {
		dbError(w, err)
		return
	}

	rows := make([]applicationRow, len(apps))
	for i, a := range apps {
		rows[i] = applicationRow{Application: a, Applicant: names[a.UserID], PostingTitle: p.Title}
	}
	poster, _, err := p.Poster()
	if err != nil {
		dbError(w, err)
		return
	}
	view.HTML(w, status, "job_posting.html", jobPostingPage{
		Posting: p, Poster: poster, Applications: rows, Users: users,
		Statuses: model.ApplicationStatuses, Apply: apply, Error: errMsg,
	})
}

// Create: GET /job-postings/new -> an empty form
func (c JobPostingController) Create(w http.ResponseWriter, r *http.Request) {
	c.renderForm(w, http.StatusOK, "Yeni ilan", "/job-postings", model.JobPosting{}, "")
}

// Store: POST /job-postings -> create from the form, then go to the new posting
func (c JobPostingController) Store(w http.ResponseWriter, r *http.Request) {
	p := jobPostingFromForm(r)
	if err := p.Validate(); err != nil {
		c.renderForm(w, http.StatusBadRequest, "Yeni ilan", "/job-postings", p, err.Error())
		return
	}

	created, err := model.JobPostings.Add(p)
	if errors.Is(err, model.ErrUserNotFound) {
		c.renderForm(w, http.StatusBadRequest, "Yeni ilan", "/job-postings", p, err.Error())
		return
	}
	if err != nil {
		dbError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/job-postings/%d", created.ID), http.StatusSeeOther)
}

// Edit: GET /job-postings/{id}/edit -> the form, filled in with the current values
func (c JobPostingController) Edit(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	p, found, err := model.JobPostings.Find(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}

	c.renderForm(w, http.StatusOK, "İlanı düzenle", fmt.Sprintf("/job-postings/%d", id), p, "")
}

// Update: POST /job-postings/{id} -> save the edit form, then go to the posting
func (c JobPostingController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	p := jobPostingFromForm(r)
	p.ID = id
	action := fmt.Sprintf("/job-postings/%d", id)
	if err := p.Validate(); err != nil {
		c.renderForm(w, http.StatusBadRequest, "İlanı düzenle", action, p, err.Error())
		return
	}

	_, found, err := model.JobPostings.Replace(id, p)
	if errors.Is(err, model.ErrUserNotFound) {
		c.renderForm(w, http.StatusBadRequest, "İlanı düzenle", action, p, err.Error())
		return
	}
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, action, http.StatusSeeOther)
}

// Destroy: POST /job-postings/{id}/delete -> remove it (and its applications), then go to the list
func (JobPostingController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	found, err := model.JobPostings.Remove(id)
	if err != nil {
		dbError(w, err)
		return
	}
	if !found {
		http.Error(w, "no job posting with that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/job-postings", http.StatusSeeOther)
}

func (JobPostingController) renderForm(w http.ResponseWriter, status int, title, action string, p model.JobPosting, errMsg string) {
	_, users, err := userNames()
	if err != nil {
		dbError(w, err)
		return
	}
	view.HTML(w, status, "job_posting_form.html", jobPostingForm{
		Title: title, Action: action, Posting: p, Users: users, Types: model.EmploymentTypes, Error: errMsg,
	})
}

func jobPostingFromForm(r *http.Request) model.JobPosting {
	userID, _ := strconv.Atoi(r.FormValue("userId"))
	p := model.JobPosting{
		UserID:         userID,
		Title:          r.FormValue("title"),
		Company:        r.FormValue("company"),
		Location:       r.FormValue("location"),
		EmploymentType: r.FormValue("employmentType"),
		Description:    r.FormValue("description"),
		ApplyURL:       r.FormValue("applyUrl"),
	}
	// an empty date field means "no deadline"
	if d := strings.TrimSpace(r.FormValue("deadline")); d != "" {
		p.Deadline = &d
	}
	return p
}

// userNames returns every user, and their names by id, for the pages that
// show who posted or applied and the forms that let you pick a user.
func userNames() (map[int]string, []model.User, error) {
	users, err := model.Users.List()
	if err != nil {
		return nil, nil, err
	}
	names := make(map[int]string, len(users))
	for _, u := range users {
		names[u.ID] = u.Name
	}
	return names, users, nil
}
