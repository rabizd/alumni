-- What a user may do on the network. Every existing and new account starts as
-- 'member' (not verified yet); 'graduate' and 'admin' are given later.
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'member'
    CHECK (role IN ('member', 'graduate', 'admin'));
