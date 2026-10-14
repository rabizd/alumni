-- Accounts on the network, with two sample users.
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

INSERT INTO users (name, email, department, start_year, graduation_year, double_major, minor, advisor, prep_exempt) VALUES
    ('Ayşe Yılmaz', 'ayse@example.com', 'Bilgisayar Mühendisliği', 2017, 2021, 'Matematik', '', 'Prof. Dr. Ahmet Kaya', true),
    ('Mehmet Demir', 'mehmet@example.com', 'Hukuk', 2015, 2019, '', 'İşletme', '', false);
