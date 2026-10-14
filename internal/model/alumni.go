package model

import (
	"context"
	"errors"
	"strings"
)

// Alumni is one graduate. The json tags decide the field names in the API:
// Go uses Name, the JSON uses "name".
type Alumni struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Department     string `json:"department"`
	GraduationYear int    `json:"graduationYear"`
}

// Validate trims the name and reports it if it is missing.
func (a *Alumni) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	if a.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// AlumniStore reads and writes the alumni table in PostgreSQL. The table and
// the two starting graduates come from db/init.sql.
type AlumniStore struct{}

var Graduates = &AlumniStore{}

func (s *AlumniStore) List() ([]Alumni, error) {
	rows, err := db.Query(context.Background(),
		"SELECT id, name, department, graduation_year FROM alumni ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Alumni{}
	for rows.Next() {
		var a Alumni
		if err := rows.Scan(&a.ID, &a.Name, &a.Department, &a.GraduationYear); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Add stores the graduate. The database gives it the next id.
func (s *AlumniStore) Add(a Alumni) (Alumni, error) {
	err := db.QueryRow(context.Background(),
		"INSERT INTO alumni (name, department, graduation_year) VALUES ($1, $2, $3) RETURNING id",
		a.Name, a.Department, a.GraduationYear).Scan(&a.ID)
	return a, err
}
