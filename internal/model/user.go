// Package model is the M in MVC: the data the app works with and the rules
// that keep it valid. It knows nothing about HTTP, HTML, or JSON responses.
package model

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// User is one account on the network. No database yet: the users live in
// memory, so they are gone when the server restarts.
type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Department     string `json:"department"`
	StartYear      int    `json:"startYear"`
	GraduationYear int    `json:"graduationYear"`
	// DoubleMajor (ÇAP), Minor (yandal) and Advisor (danışman) are optional;
	// "" means none.
	DoubleMajor string `json:"doubleMajor"`
	Minor       string `json:"minor"`
	Advisor     string `json:"advisor"`
	// PrepExempt is true when the user was exempt from the English prep year.
	PrepExempt bool `json:"prepExempt"`
}

// firstYear is the earliest start year the form accepts.
const firstYear = 1900

// Validate trims the text fields and reports the first rule the user breaks.
func (u *User) Validate() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.Department = strings.TrimSpace(u.Department)
	u.DoubleMajor = strings.TrimSpace(u.DoubleMajor)
	u.Minor = strings.TrimSpace(u.Minor)
	u.Advisor = strings.TrimSpace(u.Advisor)

	thisYear := time.Now().Year()
	switch {
	case u.Name == "":
		return errors.New("name is required")
	case u.Email == "":
		return errors.New("email is required")
	case u.Department == "":
		return errors.New("department is required")
	case u.StartYear == 0:
		return errors.New("startYear is required")
	case u.StartYear < firstYear || u.StartYear > thisYear:
		return fmt.Errorf("startYear must be between %d and %d", firstYear, thisYear)
	case u.GraduationYear == 0:
		return errors.New("graduationYear is required")
	case u.GraduationYear < u.StartYear:
		return errors.New("graduationYear cannot be before startYear")
	case u.GraduationYear > thisYear:
		return errors.New("graduationYear cannot be in the future")
	}
	return nil
}

// UserPatch is the body of a PATCH. Every field is a pointer so the handler can
// tell "the client sent an empty name" (pointer to "") from "the client did not
// mention name at all" (nil), which is the whole difference between PATCH
// and PUT.
type UserPatch struct {
	Name           *string `json:"name"`
	Email          *string `json:"email"`
	Department     *string `json:"department"`
	StartYear      *int    `json:"startYear"`
	GraduationYear *int    `json:"graduationYear"`
	DoubleMajor    *string `json:"doubleMajor"`
	Minor          *string `json:"minor"`
	Advisor        *string `json:"advisor"`
	PrepExempt     *bool   `json:"prepExempt"`
}

// Validate rejects a patch that would change nothing. The field rules are
// checked by User.Validate on the patched user, because some of them (a
// graduation year before the start year) need both old and new values.
func (p *UserPatch) Validate() error {
	// An empty body would silently do nothing, which is more likely a mistake
	// than an intention.
	if *p == (UserPatch{}) {
		return errors.New("send at least one field to change")
	}
	return nil
}

// apply returns u with every field the patch mentions replaced.
func (u User) apply(p UserPatch) User {
	if p.Name != nil {
		u.Name = *p.Name
	}
	if p.Email != nil {
		u.Email = *p.Email
	}
	if p.Department != nil {
		u.Department = *p.Department
	}
	if p.StartYear != nil {
		u.StartYear = *p.StartYear
	}
	if p.GraduationYear != nil {
		u.GraduationYear = *p.GraduationYear
	}
	if p.DoubleMajor != nil {
		u.DoubleMajor = *p.DoubleMajor
	}
	if p.Minor != nil {
		u.Minor = *p.Minor
	}
	if p.Advisor != nil {
		u.Advisor = *p.Advisor
	}
	if p.PrepExempt != nil {
		u.PrepExempt = *p.PrepExempt
	}
	return u
}

type UserStore struct {
	mu     sync.Mutex
	nextID int
	items  []User
}

var Users = &UserStore{
	nextID: 3,
	items: []User{
		{ID: 1, Name: "Ayşe Yılmaz", Email: "ayse@example.com", Department: "Bilgisayar Mühendisliği",
			StartYear: 2017, GraduationYear: 2021, DoubleMajor: "Matematik", Advisor: "Prof. Dr. Ahmet Kaya", PrepExempt: true},
		{ID: 2, Name: "Mehmet Demir", Email: "mehmet@example.com", Department: "Hukuk",
			StartYear: 2015, GraduationYear: 2019, Minor: "İşletme"},
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

// Patch changes only the fields the client actually sent. The bool reports
// whether the user existed; the error reports a patched user that would break
// a rule, in which case nothing is saved.
func (s *UserStore) Patch(id int, p UserPatch) (User, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existing := range s.items {
		if existing.ID == id {
			patched := existing.apply(p)
			if err := patched.Validate(); err != nil {
				return User{}, true, err
			}
			s.items[i] = patched
			return patched, true, nil
		}
	}
	return User{}, false, nil
}
