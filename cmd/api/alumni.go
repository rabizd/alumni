package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

// store keeps the graduates in memory for now, so everything is lost when the
// server restarts. PostgreSQL replaces this later.
type store struct {
	mu     sync.Mutex
	nextID int
	items  []Alumni
}

var alumniStore = &store{
	nextID: 3,
	items: []Alumni{
		{ID: 1, Name: "Ayşe Yılmaz", Department: "Bilgisayar Mühendisliği", GraduationYear: 2021},
		{ID: 2, Name: "Mehmet Demir", Department: "Hukuk", GraduationYear: 2019},
	},
}

func (s *store) list() []Alumni {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Alumni(nil), s.items...)
}

func (s *store) add(a Alumni) Alumni {
	s.mu.Lock()
	defer s.mu.Unlock()
	a.ID = s.nextID
	s.nextID++
	s.items = append(s.items, a)
	return a
}

// GET /alumni -> every graduate, as a JSON array
func handleListAlumni(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, alumniStore.list())
}

// POST /alumni -> create one graduate from the JSON in the request body
func handleCreateAlumni(w http.ResponseWriter, r *http.Request) {
	var a Alumni

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // a typo in a field name is a mistake, not something to ignore
	if err := dec.Decode(&a); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	a.Name = strings.TrimSpace(a.Name)
	if a.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, alumniStore.add(a))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// The status line is already sent, so all we can do is record it.
		log.Printf("writing JSON response: %v", err)
	}
}
