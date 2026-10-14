-- Accounts on the network. The sample rows come from the seeder (cmd/seed).
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    email           TEXT NOT NULL,
    department      TEXT NOT NULL,
    start_year      INTEGER NOT NULL,
    graduation_year INTEGER NOT NULL,
    double_major    TEXT NOT NULL DEFAULT '',
    minor           TEXT NOT NULL DEFAULT '',
    advisor         TEXT NOT NULL DEFAULT '',
    prep_exempt     BOOLEAN NOT NULL DEFAULT false
);
