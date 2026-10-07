package model

import (
	"errors"
	"sync"
)

// Alumni is one graduate. The json tags decide the field names in the API:
// Go uses Name, the JSON uses "name".
type Alumni struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Department     string `json:"department"`
	GraduationYear int    `json:"graduationYear"`
}

// Validate reports the first required field that is missing.
func (a Alumni) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// AlumniStore keeps the graduates in memory for now, so everything is lost
// when the server restarts. PostgreSQL replaces this later.
type AlumniStore struct {
	mu     sync.Mutex
	nextID int
	items  []Alumni
}

var Graduates = &AlumniStore{
	nextID: 3,
	items: []Alumni{
		{ID: 1, Name: "Ayşe Yılmaz", Department: "Bilgisayar Mühendisliği", GraduationYear: 2021},
		{ID: 2, Name: "Mehmet Demir", Department: "Hukuk", GraduationYear: 2019},
	},
}

func (s *AlumniStore) List() []Alumni {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Alumni(nil), s.items...)
}

func (s *AlumniStore) Add(a Alumni) Alumni {
	s.mu.Lock()
	defer s.mu.Unlock()
	a.ID = s.nextID
	s.nextID++
	s.items = append(s.items, a)
	return a
}
