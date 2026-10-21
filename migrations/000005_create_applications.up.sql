-- Applications to job postings. Each one belongs to a posting and to the user who applied,
-- which makes users and job postings many-to-many.
CREATE TABLE applications (
    id             SERIAL PRIMARY KEY,
    job_posting_id INTEGER NOT NULL REFERENCES job_postings(id) ON DELETE CASCADE,
    user_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cover_letter   TEXT NOT NULL DEFAULT '',
    status         TEXT NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT date_trunc('second', now()),
    -- a user can apply to the same posting only once
    UNIQUE (job_posting_id, user_id)
);
