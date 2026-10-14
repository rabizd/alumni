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
