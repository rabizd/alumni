// Package model is the M in MVC: the data the app works with and the rules
// that keep it valid. It knows nothing about HTTP, HTML, or JSON responses.
package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// User is one account on the network.
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

// Apply returns u with every field the patch mentions replaced.
func (u User) Apply(p UserPatch) User {
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

// UserStore reads and writes the users table in PostgreSQL, so users survive
// a restart. The table comes from migrations/000001_create_users.up.sql, the
// sample rows from the seeder.
type UserStore struct{}

var Users = &UserStore{}

const userColumns = "id, name, email, department, start_year, graduation_year, double_major, minor, advisor, prep_exempt"

func scanUser(row interface{ Scan(...any) error }) (User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Department, &u.StartYear, &u.GraduationYear,
		&u.DoubleMajor, &u.Minor, &u.Advisor, &u.PrepExempt)
	return u, err
}

func (s *UserStore) List() ([]User, error) {
	rows, err := db.Query(context.Background(), "SELECT "+userColumns+" FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// Add stores the user. The database gives it the next id.
func (s *UserStore) Add(u User) (User, error) {
	return scanUser(db.QueryRow(context.Background(),
		`INSERT INTO users (name, email, department, start_year, graduation_year, double_major, minor, advisor, prep_exempt)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING `+userColumns,
		u.Name, u.Email, u.Department, u.StartYear, u.GraduationYear, u.DoubleMajor, u.Minor, u.Advisor, u.PrepExempt))
}

// Find returns the user with the given id. The bool reports whether it existed.
func (s *UserStore) Find(id int) (User, bool, error) {
	return foundUser(scanUser(db.QueryRow(context.Background(),
		"SELECT "+userColumns+" FROM users WHERE id = $1", id)))
}

// Remove deletes the user with the given id, reporting whether it existed.
func (s *UserStore) Remove(id int) (bool, error) {
	tag, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	return tag.RowsAffected() > 0, err
}

// Replace overwrites every field of the user with the given id. The bool
// reports whether such a user existed.
func (s *UserStore) Replace(id int, u User) (User, bool, error) {
	return foundUser(scanUser(db.QueryRow(context.Background(),
		`UPDATE users SET name = $1, email = $2, department = $3, start_year = $4, graduation_year = $5,
		 double_major = $6, minor = $7, advisor = $8, prep_exempt = $9
		 WHERE id = $10 RETURNING `+userColumns,
		u.Name, u.Email, u.Department, u.StartYear, u.GraduationYear, u.DoubleMajor, u.Minor, u.Advisor, u.PrepExempt, id)))
}

// foundUser turns "no row" into found=false instead of an error.
func foundUser(u User, err error) (User, bool, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, false, nil
	}
	return u, err == nil, err
}
