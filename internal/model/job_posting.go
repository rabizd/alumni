package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// JobPosting is a job or internship opening posted by a graduate. It belongs
// to the user who posted it (UserID); a user can have many.
type JobPosting struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userId"`
	Title    string `json:"title"`
	Company  string `json:"company"`
	Location string `json:"location"`
	// EmploymentType is one of the EmploymentTypes.
	EmploymentType string `json:"employmentType"`
	Description    string `json:"description"`
	// ApplyURL is optional; "" means contact the poster.
	ApplyURL string `json:"applyUrl"`
	// Deadline is the last day to apply, as "2006-01-02". nil means open until filled.
	Deadline *string `json:"deadline"`
	// CreatedAt is set by the database when the posting is added.
	CreatedAt time.Time `json:"createdAt"`
}

// EmploymentTypes are the values EmploymentType accepts.
var EmploymentTypes = []string{"full_time", "part_time", "internship"}

// ErrUserNotFound is returned when a posting names a user that does not exist.
var ErrUserNotFound = errors.New("userId does not match an existing user")

const dateLayout = "2006-01-02"

// Validate trims the text fields and reports the first rule the posting breaks.
// Whether UserID is a real user is checked by the database (the foreign key).
func (p *JobPosting) Validate() error {
	p.Title = strings.TrimSpace(p.Title)
	p.Company = strings.TrimSpace(p.Company)
	p.Location = strings.TrimSpace(p.Location)
	p.EmploymentType = strings.TrimSpace(p.EmploymentType)
	p.Description = strings.TrimSpace(p.Description)
	p.ApplyURL = strings.TrimSpace(p.ApplyURL)

	switch {
	case p.UserID == 0:
		return errors.New("userId is required")
	case p.Title == "":
		return errors.New("title is required")
	case p.Company == "":
		return errors.New("company is required")
	case p.EmploymentType == "":
		return errors.New("employmentType is required")
	case !isEmploymentType(p.EmploymentType):
		return errors.New("employmentType must be full_time, part_time or internship")
	case p.Description == "":
		return errors.New("description is required")
	}

	if p.Deadline != nil {
		if _, err := time.Parse(dateLayout, *p.Deadline); err != nil {
			return errors.New("deadline must be a date like 2026-12-31")
		}
		// Dates in this layout sort as strings, so no time zone arithmetic is needed.
		if *p.Deadline < time.Now().Format(dateLayout) {
			return errors.New("deadline cannot be in the past")
		}
	}
	return nil
}

func isEmploymentType(s string) bool {
	for _, t := range EmploymentTypes {
		if s == t {
			return true
		}
	}
	return false
}

// Poster returns the user who posted p (the "belongs to" side of the relation).
func (p JobPosting) Poster() (User, bool, error) {
	return Users.Find(p.UserID)
}

// JobPostings returns the postings u has posted (the "has many" side of the relation).
func (u User) JobPostings() ([]JobPosting, error) {
	return JobPostings.ListByUser(u.ID)
}

// JobPostingStore reads and writes the job_postings table in PostgreSQL. The
// table comes from migrations/000004_create_job_postings.up.sql.
type JobPostingStore struct{}

var JobPostings = &JobPostingStore{}

// deadline::text gives "2006-01-02", so the column scans straight into *string.
const jobPostingColumns = "id, user_id, title, company, location, employment_type, description, apply_url, deadline::text, created_at"

func scanJobPosting(row interface{ Scan(...any) error }) (JobPosting, error) {
	var p JobPosting
	err := row.Scan(&p.ID, &p.UserID, &p.Title, &p.Company, &p.Location, &p.EmploymentType,
		&p.Description, &p.ApplyURL, &p.Deadline, &p.CreatedAt)
	return p, err
}

// List returns every posting, newest first.
func (s *JobPostingStore) List() ([]JobPosting, error) {
	return s.query("SELECT " + jobPostingColumns + " FROM job_postings ORDER BY created_at DESC, id DESC")
}

// ListByUser returns the postings of one user, newest first.
func (s *JobPostingStore) ListByUser(userID int) ([]JobPosting, error) {
	return s.query("SELECT "+jobPostingColumns+" FROM job_postings WHERE user_id = $1 ORDER BY created_at DESC, id DESC", userID)
}

func (s *JobPostingStore) query(sql string, args ...any) ([]JobPosting, error) {
	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []JobPosting{}
	for rows.Next() {
		p, err := scanJobPosting(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Add stores the posting. The database gives it the next id and the current
// time; ErrUserNotFound means UserID is not a user.
func (s *JobPostingStore) Add(p JobPosting) (JobPosting, error) {
	p, err := scanJobPosting(db.QueryRow(context.Background(),
		`INSERT INTO job_postings (user_id, title, company, location, employment_type, description, apply_url, deadline)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING `+jobPostingColumns,
		p.UserID, p.Title, p.Company, p.Location, p.EmploymentType, p.Description, p.ApplyURL, p.Deadline))
	return p, userNotFound(err)
}

// Find returns the posting with the given id. The bool reports whether it existed.
func (s *JobPostingStore) Find(id int) (JobPosting, bool, error) {
	return foundJobPosting(scanJobPosting(db.QueryRow(context.Background(),
		"SELECT "+jobPostingColumns+" FROM job_postings WHERE id = $1", id)))
}

// Replace overwrites every field of the posting with the given id except the
// id and the creation time. The bool reports whether such a posting existed.
func (s *JobPostingStore) Replace(id int, p JobPosting) (JobPosting, bool, error) {
	p, ok, err := foundJobPosting(scanJobPosting(db.QueryRow(context.Background(),
		`UPDATE job_postings SET user_id = $1, title = $2, company = $3, location = $4, employment_type = $5,
		 description = $6, apply_url = $7, deadline = $8
		 WHERE id = $9 RETURNING `+jobPostingColumns,
		p.UserID, p.Title, p.Company, p.Location, p.EmploymentType, p.Description, p.ApplyURL, p.Deadline, id)))
	return p, ok, userNotFound(err)
}

// Remove deletes the posting with the given id, reporting whether it existed.
func (s *JobPostingStore) Remove(id int) (bool, error) {
	tag, err := db.Exec(context.Background(), "DELETE FROM job_postings WHERE id = $1", id)
	return tag.RowsAffected() > 0, err
}

// foundJobPosting turns "no row" into found=false instead of an error.
func foundJobPosting(p JobPosting, err error) (JobPosting, bool, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return JobPosting{}, false, nil
	}
	return p, err == nil, err
}

// userNotFound turns a foreign key violation (Postgres code 23503) into
// ErrUserNotFound, so a controller can answer 400 instead of 500.
func userNotFound(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrUserNotFound
	}
	return err
}
