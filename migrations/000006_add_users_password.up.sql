-- Accounts made on /register have a password. Only its bcrypt hash is stored.
-- Rows that existed before (seeded users) get '' and cannot sign in.
ALTER TABLE users ADD COLUMN password_hash TEXT NOT NULL DEFAULT '';
