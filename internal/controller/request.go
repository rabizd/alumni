package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

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

// dbError answers 500 when the database fails. The client only sees "database
// error"; the real error goes to the log, where the developer can read it.
func dbError(w http.ResponseWriter, err error) {
	log.Printf("database error: %v", err)
	http.Error(w, "database error", http.StatusInternalServerError)
}
