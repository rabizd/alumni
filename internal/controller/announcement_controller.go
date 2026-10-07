package controller

import (
	"fmt"
	"net/http"

	"github.com/rabizd/alumni/internal/model"
	"github.com/rabizd/alumni/internal/view"
)

// AnnouncementController serves /announcements as HTML pages: the interface
// for managing announcements in a browser. Like UserController, it uses POST
// for update and delete because HTML forms can only send GET and POST.
type AnnouncementController struct{}

// announcementList is what announcements.html needs: the announcements, and
// an empty one for the create form under the list.
type announcementList struct {
	Announcements []model.Announcement
	New           model.Announcement
}

// announcementForm is what announcement_form.html needs.
type announcementForm struct {
	Title        string
	Action       string
	Announcement model.Announcement
	Error        string
}

// Index: GET /announcements -> the list page, with a create form
func (AnnouncementController) Index(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "announcements.html", announcementList{Announcements: model.Announcements.List()})
}

// Show: GET /announcements/{id} -> one announcement's page, or 404
func (AnnouncementController) Show(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a, found := model.Announcements.Find(id)
	if !found {
		http.Error(w, "no announcement with that id", http.StatusNotFound)
		return
	}

	view.HTML(w, http.StatusOK, "announcement.html", a)
}

// Create: GET /announcements/new -> an empty form
func (AnnouncementController) Create(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, http.StatusOK, "announcement_form.html",
		announcementForm{Title: "Yeni duyuru", Action: "/announcements"})
}

// Store: POST /announcements -> create from the form, then go to the list
func (AnnouncementController) Store(w http.ResponseWriter, r *http.Request) {
	a := announcementFromForm(r)
	if err := a.Validate(); err != nil {
		view.HTML(w, http.StatusBadRequest, "announcement_form.html",
			announcementForm{Title: "Yeni duyuru", Action: "/announcements", Announcement: a, Error: err.Error()})
		return
	}

	model.Announcements.Add(a)
	http.Redirect(w, r, "/announcements", http.StatusSeeOther)
}

// Edit: GET /announcements/{id}/edit -> the form, filled in with the current values
func (AnnouncementController) Edit(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a, found := model.Announcements.Find(id)
	if !found {
		http.Error(w, "no announcement with that id", http.StatusNotFound)
		return
	}

	view.HTML(w, http.StatusOK, "announcement_form.html",
		announcementForm{Title: "Duyuruyu düzenle", Action: fmt.Sprintf("/announcements/%d", id), Announcement: a})
}

// Update: POST /announcements/{id} -> save the edit form, then go to the announcement
func (AnnouncementController) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	a := announcementFromForm(r)
	a.ID = id
	if err := a.Validate(); err != nil {
		view.HTML(w, http.StatusBadRequest, "announcement_form.html",
			announcementForm{Title: "Duyuruyu düzenle", Action: fmt.Sprintf("/announcements/%d", id), Announcement: a, Error: err.Error()})
		return
	}

	if _, found := model.Announcements.Replace(id, a); !found {
		http.Error(w, "no announcement with that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/announcements/%d", id), http.StatusSeeOther)
}

// Destroy: POST /announcements/{id}/delete -> remove it, then go to the list
func (AnnouncementController) Destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if !model.Announcements.Remove(id) {
		http.Error(w, "no announcement with that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/announcements", http.StatusSeeOther)
}

func announcementFromForm(r *http.Request) model.Announcement {
	return model.Announcement{
		Title:  r.FormValue("title"),
		Body:   r.FormValue("body"),
		Author: r.FormValue("author"),
	}
}
