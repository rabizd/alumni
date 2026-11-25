package model

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// SessionLength is how long a sign-in lasts before the user has to sign in again.
const SessionLength = 7 * 24 * time.Hour

// SessionStore reads and writes the sessions table: which browser is signed
// in as which user. The table comes from migrations/000008_create_sessions.up.sql.
type SessionStore struct{}

var Sessions = &SessionStore{}

// Create starts a session for the user and returns its token, a random value
// nobody can guess, which the browser keeps in a cookie.
func (s *SessionStore) Create(userID int) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	_, err := db.Exec(context.Background(),
		"INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)",
		token, userID, time.Now().Add(SessionLength))
	return token, err
}

// User returns the user a token belongs to. The bool is false when the token
// is unknown or expired.
func (s *SessionStore) User(token string) (User, bool, error) {
	return foundUser(scanUser(db.QueryRow(context.Background(),
		`SELECT `+userColumns+` FROM users
		 WHERE id = (SELECT user_id FROM sessions WHERE token = $1 AND expires_at > now())`, token)))
}

// Delete ends the session with the given token.
func (s *SessionStore) Delete(token string) error {
	_, err := db.Exec(context.Background(), "DELETE FROM sessions WHERE token = $1", token)
	return err
}

// FindForLogin returns the user with this email (ignoring case) and their
// password hash. The bool is false when no user has that email.
func (s *UserStore) FindForLogin(email string) (User, string, bool, error) {
	var hash string
	row := db.QueryRow(context.Background(),
		"SELECT "+userColumns+", password_hash FROM users WHERE lower(email) = lower($1) ORDER BY id LIMIT 1", email)
	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Department, &u.StartYear, &u.GraduationYear,
		&u.DoubleMajor, &u.Minor, &u.Advisor, &u.PrepExempt, &u.Role, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", false, nil
	}
	return u, hash, err == nil, err
}

type currentUserKey struct{}

// WithCurrentUser returns a context that carries the signed-in user.
func WithCurrentUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, currentUserKey{}, u)
}

// CurrentUser returns the signed-in user the context carries, if any.
func CurrentUser(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(currentUserKey{}).(User)
	return u, ok
}
