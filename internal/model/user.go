// Package model is the M in MVC: the data the app works with and the rules
// that keep it valid. It knows nothing about HTTP, HTML, or JSON responses.
package model

import (
	"errors"
	"strings"
	"sync"
)

// User is one account on the network. No database yet: the users live in
// memory, so they are gone when the server restarts.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validate trims the fields and reports the first required one that is missing.
func (u *User) Validate() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	switch {
	case u.Name == "":
		return errors.New("name is required")
	case u.Email == "":
		return errors.New("email is required")
	}
	return nil
}

// UserPatch is the body of a PATCH. Every field is a pointer so the handler can
// tell "the client sent an empty name" (pointer to "") from "the client did not
// mention name at all" (nil), which is the whole difference between PATCH
// and PUT.
type UserPatch struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// Validate trims the fields that were sent and rejects a patch that would
// change nothing or leave a field empty.
func (p *UserPatch) Validate() error {
	// An empty body would silently do nothing, which is more likely a mistake
	// than an intention.
	if p.Name == nil && p.Email == nil {
		return errors.New("send at least one of name or email")
	}
	if p.Name != nil {
		trimmed := strings.TrimSpace(*p.Name)
		if trimmed == "" {
			return errors.New("name cannot be empty")
		}
		p.Name = &trimmed
	}
	if p.Email != nil {
		trimmed := strings.TrimSpace(*p.Email)
		if trimmed == "" {
			return errors.New("email cannot be empty")
		}
		p.Email = &trimmed
	}
	return nil
}

type UserStore struct {
	mu     sync.Mutex
	nextID int
	items  []User
}

var Users = &UserStore{
	nextID: 3,
	items: []User{
		{ID: 1, Name: "Ayşe Yılmaz", Email: "ayse@example.com"},
		{ID: 2, Name: "Mehmet Demir", Email: "mehmet@example.com"},
	},
}

func (s *UserStore) List() []User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]User(nil), s.items...)
}

func (s *UserStore) Add(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextID
	s.nextID++
	s.items = append(s.items, u)
	return u
}

// Find returns the user with the given id. The bool reports whether it existed.
func (s *UserStore) Find(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.items {
		if u.ID == id {
			return u, true
		}
	}
	return User{}, false
}

// Remove deletes the user with the given id, reporting whether it existed.
func (s *UserStore) Remove(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.items {
		if u.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

// Replace overwrites every field of the user with the given id. The bool
// reports whether such a user existed.
func (s *UserStore) Replace(id int, u User) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existing := range s.items {
		if existing.ID == id {
			u.ID = id // the id belongs to the server, not to the request body
			s.items[i] = u
			return u, true
		}
	}
	return User{}, false
}

// Patch changes only the fields the client actually sent.
func (s *UserStore) Patch(id int, p UserPatch) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existing := range s.items {
		if existing.ID == id {
			if p.Name != nil {
				existing.Name = *p.Name
			}
			if p.Email != nil {
				existing.Email = *p.Email
			}
			s.items[i] = existing
			return existing, true
		}
	}
	return User{}, false
}
