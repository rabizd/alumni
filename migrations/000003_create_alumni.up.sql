-- The simple graduate directory, with two sample graduates.
CREATE TABLE alumni (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    department      TEXT NOT NULL DEFAULT '',
    graduation_year INTEGER NOT NULL DEFAULT 0
);

INSERT INTO alumni (name, department, graduation_year) VALUES
    ('Ayşe Yılmaz', 'Bilgisayar Mühendisliği', 2021),
    ('Mehmet Demir', 'Hukuk', 2019);
