-- The simple graduate directory. The sample rows come from the seeder (cmd/seed).
CREATE TABLE alumni (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    department      TEXT NOT NULL DEFAULT '',
    graduation_year INTEGER NOT NULL DEFAULT 0
);
