package main

import "net/http"

// health is the answer to GET /api/health. It is a struct rather than a map so
// the field order in the JSON is fixed.
type health struct {
	Status string `json:"status"`
}

// GET /api/health -> {"status":"ok"}
//
// Reaching this handler at all is the real check: it means the process is up
// and the router is wired correctly.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, health{Status: "ok"})
}
