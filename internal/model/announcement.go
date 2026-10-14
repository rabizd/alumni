package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Announcement is one post on the announcement board.
type Announcement struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
	// CreatedAt is set by the database when the announcement is added.
	CreatedAt time.Time `json:"createdAt"`
}

// Validate trims the text fields and reports the first required one that is missing.
func (a *Announcement) Validate() error {
	a.Title = strings.TrimSpace(a.Title)
	a.Body = strings.TrimSpace(a.Body)
	a.Author = strings.TrimSpace(a.Author)
	switch {
	case a.Title == "":
		return errors.New("title is required")
	case a.Body == "":
		return errors.New("body is required")
	case a.Author == "":
		return errors.New("author is required")
	}
	return nil
}

// AnnouncementStore reads and writes the announcements table in PostgreSQL,
// so announcements survive a restart. The table comes from
// migrations/000002_create_announcements.up.sql, the sample rows from the seeder.
type AnnouncementStore struct{}

var Announcements = &AnnouncementStore{}

const announcementColumns = "id, title, body, author, created_at"

func scanAnnouncement(row interface{ Scan(...any) error }) (Announcement, error) {
	var a Announcement
	err := row.Scan(&a.ID, &a.Title, &a.Body, &a.Author, &a.CreatedAt)
	return a, err
}

// List returns every announcement, newest first.
func (s *AnnouncementStore) List() ([]Announcement, error) {
	rows, err := db.Query(context.Background(),
		"SELECT "+announcementColumns+" FROM announcements ORDER BY created_at DESC, id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Announcement{}
	for rows.Next() {
		a, err := scanAnnouncement(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Add stores the announcement. The database gives it the next id and the current time.
func (s *AnnouncementStore) Add(a Announcement) (Announcement, error) {
	return scanAnnouncement(db.QueryRow(context.Background(),
		"INSERT INTO announcements (title, body, author) VALUES ($1, $2, $3) RETURNING "+announcementColumns,
		a.Title, a.Body, a.Author))
}

// Find returns the announcement with the given id. The bool reports whether it existed.
func (s *AnnouncementStore) Find(id int) (Announcement, bool, error) {
	return found(scanAnnouncement(db.QueryRow(context.Background(),
		"SELECT "+announcementColumns+" FROM announcements WHERE id = $1", id)))
}

// Replace overwrites the title, body and author of the announcement with the
// given id. The id and the creation time stay as they were.
func (s *AnnouncementStore) Replace(id int, a Announcement) (Announcement, bool, error) {
	return found(scanAnnouncement(db.QueryRow(context.Background(),
		"UPDATE announcements SET title = $1, body = $2, author = $3 WHERE id = $4 RETURNING "+announcementColumns,
		a.Title, a.Body, a.Author, id)))
}

// Remove deletes the announcement with the given id, reporting whether it existed.
func (s *AnnouncementStore) Remove(id int) (bool, error) {
	tag, err := db.Exec(context.Background(), "DELETE FROM announcements WHERE id = $1", id)
	return tag.RowsAffected() > 0, err
}

// found turns "no row" into found=false instead of an error.
func found(a Announcement, err error) (Announcement, bool, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return Announcement{}, false, nil
	}
	return a, err == nil, err
}
