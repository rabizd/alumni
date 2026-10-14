-- Posts on the announcement board. The sample rows come from the seeder (cmd/seed).
CREATE TABLE announcements (
    id         SERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    body       TEXT NOT NULL,
    author     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT date_trunc('second', now())
);
