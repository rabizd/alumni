-- One row per signed-in browser. The token is the random value in the
-- "session" cookie; signing out deletes the row. Deleting a user deletes
-- their sessions too.
CREATE TABLE sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL
);
