-- Job and internship openings posted by graduates. Each one belongs to the user who posted it.
CREATE TABLE job_postings (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    company         TEXT NOT NULL,
    location        TEXT NOT NULL DEFAULT '',
    employment_type TEXT NOT NULL,
    description     TEXT NOT NULL,
    apply_url       TEXT NOT NULL DEFAULT '',
    deadline        DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT date_trunc('second', now())
);
