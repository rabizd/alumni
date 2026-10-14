package model

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Alumni is one graduate. The json tags decide the field names in the API:
// Go uses Name, the JSON uses "name".
type Alumni struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Department     string `json:"department"`
	GraduationYear int    `json:"graduationYear"`
}

// Validate trims the text fields and reports a missing name.
func (a *Alumni) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	a.Department = strings.TrimSpace(a.Department)
	if a.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// AlumniStore reads and writes the alumni table in PostgreSQL. The table
// comes from migrations/000003_create_alumni.up.sql, the rows from the seeder.
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

// Find returns the graduate with the given id. The bool reports whether it existed.
func (s *AlumniStore) Find(id int) (Alumni, bool, error) {
	var a Alumni
	err := db.QueryRow(context.Background(),
		"SELECT id, name, department, graduation_year FROM alumni WHERE id = $1", id).
		Scan(&a.ID, &a.Name, &a.Department, &a.GraduationYear)
	return foundAlumni(a, err)
}

// Replace overwrites every field of the graduate with the given id and keeps the id.
func (s *AlumniStore) Replace(id int, a Alumni) (Alumni, bool, error) {
	err := db.QueryRow(context.Background(),
		"UPDATE alumni SET name = $1, department = $2, graduation_year = $3 WHERE id = $4 RETURNING id",
		a.Name, a.Department, a.GraduationYear, id).Scan(&a.ID)
	return foundAlumni(a, err)
}

// Remove deletes the graduate with the given id, reporting whether it existed.
func (s *AlumniStore) Remove(id int) (bool, error) {
	tag, err := db.Exec(context.Background(), "DELETE FROM alumni WHERE id = $1", id)
	return tag.RowsAffected() > 0, err
}

// foundAlumni turns "no row" into found=false instead of an error.
func foundAlumni(a Alumni, err error) (Alumni, bool, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return Alumni{}, false, nil
	}
	return a, err == nil, err
}
