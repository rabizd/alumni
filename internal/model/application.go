package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Application is a user's application to a job posting. It belongs to the
// posting (JobPostingID) and to the user who applied (UserID); through it,
// users and job postings are many-to-many.
type Application struct {
	ID           int `json:"id"`
	JobPostingID int `json:"jobPostingId"`
	UserID       int `json:"userId"`
	// CoverLetter is optional; "" means none.
	CoverLetter string `json:"coverLetter"`
	// Status is one of the ApplicationStatuses; "" becomes "pending".
	Status string `json:"status"`
	// CreatedAt is set by the database when the application is added.
	CreatedAt time.Time `json:"createdAt"`
}

// ApplicationStatuses are the values Status accepts.
var ApplicationStatuses = []string{"pending", "accepted", "rejected"}

var (
	// ErrJobPostingNotFound is returned when an application names a posting that does not exist.
	ErrJobPostingNotFound = errors.New("jobPostingId does not match an existing job posting")
	// ErrAlreadyApplied is returned when the user has already applied to that posting.
	ErrAlreadyApplied = errors.New("this user has already applied to this job posting")
)

// Validate trims the cover letter, fills in the default status and reports
// the first rule the application breaks. Whether the posting and the user
// exist, and whether this is a second application, is checked by the database.
func (a *Application) Validate() error {
	a.CoverLetter = strings.TrimSpace(a.CoverLetter)
	a.Status = strings.TrimSpace(a.Status)
	if a.Status == "" {
		a.Status = "pending"
	}

	switch {
	case a.JobPostingID == 0:
		return errors.New("jobPostingId is required")
	case a.UserID == 0:
		return errors.New("userId is required")
	case !isApplicationStatus(a.Status):
		return errors.New("status must be pending, accepted or rejected")
	}
	return nil
}

func isApplicationStatus(s string) bool {
	for _, st := range ApplicationStatuses {
		if s == st {
			return true
		}
	}
	return false
}

// JobPosting returns the posting a was made to (belongs to).
func (a Application) JobPosting() (JobPosting, bool, error) {
	return JobPostings.Find(a.JobPostingID)
}

// Applicant returns the user who made a (belongs to).
func (a Application) Applicant() (User, bool, error) {
	return Users.Find(a.UserID)
}

// Applications returns the applications p has received (has many).
func (p JobPosting) Applications() ([]Application, error) {
	return Applications.ListByJobPosting(p.ID)
}

// Applications returns the applications u has made (has many).
func (u User) Applications() ([]Application, error) {
	return Applications.ListByUser(u.ID)
}

// ApplicationStore reads and writes the applications table in PostgreSQL. The
// table comes from migrations/000005_create_applications.up.sql.
type ApplicationStore struct{}

var Applications = &ApplicationStore{}

const applicationColumns = "id, job_posting_id, user_id, cover_letter, status, created_at"

func scanApplication(row interface{ Scan(...any) error }) (Application, error) {
	var a Application
	err := row.Scan(&a.ID, &a.JobPostingID, &a.UserID, &a.CoverLetter, &a.Status, &a.CreatedAt)
	return a, err
}

// List returns every application, newest first.
func (s *ApplicationStore) List() ([]Application, error) {
	return s.query("SELECT " + applicationColumns + " FROM applications ORDER BY created_at DESC, id DESC")
}

// ListByJobPosting returns the applications to one posting, newest first.
func (s *ApplicationStore) ListByJobPosting(jobPostingID int) ([]Application, error) {
	return s.query("SELECT "+applicationColumns+" FROM applications WHERE job_posting_id = $1 ORDER BY created_at DESC, id DESC", jobPostingID)
}

// ListByUser returns the applications one user has made, newest first.
func (s *ApplicationStore) ListByUser(userID int) ([]Application, error) {
	return s.query("SELECT "+applicationColumns+" FROM applications WHERE user_id = $1 ORDER BY created_at DESC, id DESC", userID)
}

func (s *ApplicationStore) query(sql string, args ...any) ([]Application, error) {
	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Application{}
	for rows.Next() {
		a, err := scanApplication(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Add stores the application. The database gives it the next id and the
// current time; see applicationError for the errors it can return.
func (s *ApplicationStore) Add(a Application) (Application, error) {
	a, err := scanApplication(db.QueryRow(context.Background(),
		`INSERT INTO applications (job_posting_id, user_id, cover_letter, status)
		 VALUES ($1, $2, $3, $4) RETURNING `+applicationColumns,
		a.JobPostingID, a.UserID, a.CoverLetter, a.Status))
	return a, applicationError(err)
}

// Find returns the application with the given id. The bool reports whether it existed.
func (s *ApplicationStore) Find(id int) (Application, bool, error) {
	return foundApplication(scanApplication(db.QueryRow(context.Background(),
		"SELECT "+applicationColumns+" FROM applications WHERE id = $1", id)))
}

// Replace overwrites every field of the application with the given id except
// the id and the creation time. The bool reports whether such an application existed.
func (s *ApplicationStore) Replace(id int, a Application) (Application, bool, error) {
	a, ok, err := foundApplication(scanApplication(db.QueryRow(context.Background(),
		`UPDATE applications SET job_posting_id = $1, user_id = $2, cover_letter = $3, status = $4
		 WHERE id = $5 RETURNING `+applicationColumns,
		a.JobPostingID, a.UserID, a.CoverLetter, a.Status, id)))
	return a, ok, applicationError(err)
}

// Remove deletes the application with the given id, reporting whether it existed.
func (s *ApplicationStore) Remove(id int) (bool, error) {
	tag, err := db.Exec(context.Background(), "DELETE FROM applications WHERE id = $1", id)
	return tag.RowsAffected() > 0, err
}

// foundApplication turns "no row" into found=false instead of an error.
func foundApplication(a Application, err error) (Application, bool, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return Application{}, false, nil
	}
	return a, err == nil, err
}

// applicationError turns the database's constraint errors into ones a
// controller can answer with 400 or 409 instead of 500: a foreign key
// violation (23503) names the missing posting or user, a unique violation
// (23505) means a second application.
func applicationError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch {
	case pgErr.Code == "23503" && pgErr.ConstraintName == "applications_job_posting_id_fkey":
		return ErrJobPostingNotFound
	case pgErr.Code == "23503" && pgErr.ConstraintName == "applications_user_id_fkey":
		return ErrUserNotFound
	case pgErr.Code == "23505":
		return ErrAlreadyApplied
	}
	return err
}
