package model

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// Announcement is one post on the announcement board. Like the users, the
// announcements live in memory for now and are gone when the server restarts.
type Announcement struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
	// CreatedAt is set by the server when the announcement is added.
	CreatedAt time.Time `json:"createdAt"`
}

// Validate trims the text fields and reports the first required one that is missing.
func (a *Announcement) Validate() error {
	a.Title = strings.TrimSpace(a.Title)
	a.Body = strings.TrimSpace(a.Body)
	a.Author = strings.TrimSpace(a.Author)
	switch {
	case a.Title == "":
		return errors.New("title is required")
	case a.Body == "":
		return errors.New("body is required")
	case a.Author == "":
		return errors.New("author is required")
	}
	return nil
}

type AnnouncementStore struct {
	mu     sync.Mutex
	nextID int
	items  []Announcement
}

var Announcements = &AnnouncementStore{
	nextID: 3,
	items: []Announcement{
		{ID: 1, Title: "Mezunlar buluşması", Author: "Mezunlar Ofisi",
			Body:      "Bu yılki mezunlar buluşması 15 Kasım'da Beyazıt kampüsünde. Kayıt için profilinizin güncel olması yeterli.",
			CreatedAt: time.Date(2026, 10, 1, 10, 0, 0, 0, time.Local)},
		{ID: 2, Title: "Kariyer günleri başvuruları açıldı", Author: "Kariyer Merkezi",
			Body:      "Şirketinizle kariyer günlerine katılmak ve öğrencilerle tanışmak isteyen mezunlarımız 31 Ekim'e kadar başvurabilir.",
			CreatedAt: time.Date(2026, 10, 5, 14, 30, 0, 0, time.Local)},
	},
}

// List returns every announcement, newest first.
func (s *AnnouncementStore) List() []Announcement {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Announcement, len(s.items))
	for i, a := range s.items {
		out[len(s.items)-1-i] = a
	}
	return out
}

// Add gives the announcement the next id and the current time, and stores it.
func (s *AnnouncementStore) Add(a Announcement) Announcement {
	s.mu.Lock()
	defer s.mu.Unlock()
	a.ID = s.nextID
	s.nextID++
	a.CreatedAt = time.Now().Truncate(time.Second)
	s.items = append(s.items, a)
	return a
}

// Find returns the announcement with the given id. The bool reports whether it existed.
func (s *AnnouncementStore) Find(id int) (Announcement, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.items {
		if a.ID == id {
			return a, true
		}
	}
	return Announcement{}, false
}

// Replace overwrites the title, body and author of the announcement with the
// given id. The id and the creation time stay as they were.
func (s *AnnouncementStore) Replace(id int, a Announcement) (Announcement, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existing := range s.items {
		if existing.ID == id {
			a.ID = id
			a.CreatedAt = existing.CreatedAt
			s.items[i] = a
			return a, true
		}
	}
	return Announcement{}, false
}

// Remove deletes the announcement with the given id, reporting whether it existed.
func (s *AnnouncementStore) Remove(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, a := range s.items {
		if a.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}
