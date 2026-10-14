-- Postgres runs this once, when the db volume is empty. After that the data
-- in the volume is kept, so these two rows are not inserted again.
CREATE TABLE announcements (
    id         SERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    body       TEXT NOT NULL,
    author     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT date_trunc('second', now())
);

INSERT INTO announcements (title, body, author, created_at) VALUES
    ('Mezunlar buluşması',
     'Bu yılki mezunlar buluşması 15 Kasım''da Beyazıt kampüsünde. Kayıt için profilinizin güncel olması yeterli.',
     'Mezunlar Ofisi', '2026-10-01 10:00:00+03'),
    ('Kariyer günleri başvuruları açıldı',
     'Şirketinizle kariyer günlerine katılmak ve öğrencilerle tanışmak isteyen mezunlarımız 31 Ekim''e kadar başvurabilir.',
     'Kariyer Merkezi', '2026-10-05 14:30:00+03');

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

CREATE TABLE alumni (
    id              SERIAL PRIMARY KEY,
    name            TEXT NOT NULL,
    department      TEXT NOT NULL DEFAULT '',
    graduation_year INTEGER NOT NULL DEFAULT 0
);

INSERT INTO alumni (name, department, graduation_year) VALUES
    ('Ayşe Yılmaz', 'Bilgisayar Mühendisliği', 2021),
    ('Mehmet Demir', 'Hukuk', 2019);
