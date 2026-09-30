package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// userPatch is the body of a PATCH. Every field is a pointer so the handler can
// tell "the client sent an empty name" (pointer to "") from "the client did not
// mention name at all" (nil), which is the whole difference between PATCH
// and PUT.
type userPatch struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type userStore struct {
	mu     sync.Mutex
	nextID int
	items  []User
}

var users = &userStore{
	nextID: 3,
	items: []User{
		{ID: 1, Name: "Ayşe Yılmaz", Email: "ayse@example.com"},
		{ID: 2, Name: "Mehmet Demir", Email: "mehmet@example.com"},
	},
}

func (s *userStore) list() []User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]User(nil), s.items...)
}

func (s *userStore) add(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextID
	s.nextID++
	s.items = append(s.items, u)
	return u
}

// find returns the user with the given id. The bool reports whether it existed.
func (s *userStore) find(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.items {
		if u.ID == id {
			return u, true
		}
	}
	return User{}, false
}

// remove deletes the user with the given id, reporting whether it existed.
func (s *userStore) remove(id int) bool {
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

// replace overwrites every field of the user with the given id. The bool
// reports whether such a user existed.
func (s *userStore) replace(id int, u User) (User, bool) {
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

// patch changes only the fields the client actually sent.
func (s *userStore) patch(id int, p userPatch) (User, bool) {
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

// GET /api/users -> every user, as a JSON array
func handleListUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, users.list())
}

// GET /api/users/{id} -> one user, or 404
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	u, found := users.find(id)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, u)
}

// DELETE /api/users/{id} -> 204 with an empty body, or 404
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	if !users.remove(id) {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	// 204 No Content: it worked, and there is nothing left to send back.
	w.WriteHeader(http.StatusNoContent)
}

// POST /api/users -> create a user; 201 with the id the server assigned
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if !decodeJSON(w, r, &u) {
		return
	}

	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	if msg, ok := validateUser(u); !ok {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, users.add(u))
}

// PUT /api/users/{id} -> replace the whole user, so every field is required
func handleReplaceUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var u User
	if !decodeJSON(w, r, &u) {
		return
	}

	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	if msg, ok := validateUser(u); !ok {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	updated, found := users.replace(id, u)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// PATCH /api/users/{id} -> change only the fields the body mentions
func handlePatchUser(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}

	var p userPatch
	if !decodeJSON(w, r, &p) {
		return
	}

	// An empty body would silently do nothing, which is more likely a mistake
	// than an intention.
	if p.Name == nil && p.Email == nil {
		http.Error(w, "send at least one of name or email", http.StatusBadRequest)
		return
	}

	if p.Name != nil {
		trimmed := strings.TrimSpace(*p.Name)
		if trimmed == "" {
			http.Error(w, "name cannot be empty", http.StatusBadRequest)
			return
		}
		p.Name = &trimmed
	}
	if p.Email != nil {
		trimmed := strings.TrimSpace(*p.Email)
		if trimmed == "" {
			http.Error(w, "email cannot be empty", http.StatusBadRequest)
			return
		}
		p.Email = &trimmed
	}

	updated, found := users.patch(id, p)
	if !found {
		http.Error(w, "no user with that id", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// pathID reads {id} from the URL. It answers the request itself and reports
// false when the id is not a number.
func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id must be a whole number", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

// decodeJSON reads the request body into dst, answering with 400 and reporting
// false if the body is not the JSON we expect.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // a typo in a field name is a mistake, not something to ignore
	if err := dec.Decode(dst); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

func validateUser(u User) (string, bool) {
	switch {
	case u.Name == "":
		return "name is required", false
	case u.Email == "":
		return "email is required", false
	}
	return "", true
}
