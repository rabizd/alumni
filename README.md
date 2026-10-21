<div align="center">

# 🎓 Alumni — Istanbul University Alumni Network

**A social network for Istanbul University graduates — part directory, part feed, entirely closed to everyone else.**

[![Status](https://img.shields.io/badge/status-early%20development-orange)](https://github.com/rabizd/alumni)
[![Go](https://img.shields.io/badge/backend-Go-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/database-PostgreSQL-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/cache-Redis-DC382D?logo=redis&logoColor=white)](https://redis.io)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

> ⚠️ **Early development.** A first slice of the API runs — see [What Works Today](#what-works-today). Everything below marked *planned* is a design decision, not a shipped feature.

## Table of Contents

- [About](#about)
- [Why This Project](#why-this-project)
- [The Social Layer](#the-social-layer)
- [Planned Features](#planned-features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [What Works Today](#what-works-today)
- [Data Model](#data-model)
- [Data Design Documentation](DDD.md) ↗
- [Prior Art](#prior-art)
- [Open Design Questions](#open-design-questions)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Project Structure (MVC)](#project-structure-mvc)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About

**Alumni** is a web platform built for the graduates of **Istanbul University**. Graduation should not be the end of a community — it should be the beginning of a much larger one. Yet in practice, the moment students walk out of the gates of Beyazıt, the network they spent years building quietly dissolves into scattered group chats and outdated spreadsheets.

Alumni brings that network back into one place. It collects graduate profiles, career histories, and the institutions they work at, then makes that information **searchable, discoverable, and genuinely useful** — whether you are a recent graduate looking for a mentor, a senior professional hiring from your own faculty, or the university itself trying to stay in touch with the people it educated.

But a list of names is not a community. On top of the directory sits a feed where graduates post, reply, and follow each other — see [The Social Layer](#the-social-layer).

## Why This Project

- **The network already exists — it is just invisible.** Thousands of Istanbul University graduates work across the country and the world. Alumni makes that reality legible.
- **Careers move fast; contact lists do not.** Profiles that graduates maintain themselves stay alive without anyone curating a spreadsheet.
- **Mentorship needs a starting point.** Finding "someone from my department who does what I want to do" should take one search, not six months of luck.
- **Institutional memory has value.** A university that knows where its graduates went can advise the students who have not left yet.

## The Social Layer

Alumni is not only a searchable directory. A directory is something you visit twice — once when you sign up, once when you need a favour. What keeps a network alive is people having somewhere to *talk*.

So the platform has a social side, shaped for graduates rather than for the general public:

- **A shared feed.** Graduates post news, job openings, questions, and announcements. A senior developer posting "we are hiring two juniors, my department first" is worth more here than on any public job board, because everyone reading it went to the same school.
- **Faculty and department circles.** The Faculty of Science feed and the Law feed are different conversations. Graduates follow the circles they belong to, plus anyone whose career they want to keep up with.
- **Reactions, comments, and follows.** The ordinary social vocabulary — because it is the vocabulary people already know, and because a post with fifteen replies is how a job opening actually reaches someone.
- **Class-year nostalgia.** "Who else graduated in 2019?" is a legitimate feature, not a joke. Shared years and shared classrooms are the strongest reason two strangers on this platform will talk to each other.

The difference from a public social network is the door: **everyone in the feed is a verified Istanbul University graduate.** That single constraint is what makes the conversation worth having — no bots, no strangers, no noise. It also means the feed must be built on top of verification, not before it.

## Planned Features

**Core — the first version is not useful without these**

- 📇 **Alumni profiles** — faculty, department, graduation year, city, and current position
- 🔐 **Sign-up, sign-in, and sessions** — JWT access tokens, with sessions tracked in Redis
- ✅ **Graduate verification** — a profile is only trustworthy if the person really graduated (see [Open Design Questions](#open-design-questions))
- 🔍 **Search & filtering** — by name, faculty, department, graduation year, city, or employer
- 📝 **Career timeline** — positions and promotions added over time, so a profile ages well
- 🙋 **"Open to helping" flag** — searchable, and specific: CV review, mock interview, a short career chat
- 🛡️ **Privacy controls** — every graduate decides which fields are public, which are visible to other verified alumni, and which stay private

**Social — what makes people come back after signing up**

- 📰 **Shared feed** — posts, announcements, and questions from verified graduates
- 🏛️ **Faculty & department circles** — follow the conversations you actually belong to
- 👥 **Follows, reactions, and comments** — the ordinary social vocabulary, inside a closed network
- 💼 **Job & internship board** — openings posted by graduates, for graduates
- ✉️ **Direct messaging** — reach another graduate without exposing private contact details
- 🎓 **Class-year pages** — find the people you sat next to

**Later — valuable, but only once the core and the feed work**

- 🏢 **Company & institution directory** — see who works where, at a glance
- 🤝 **Mentorship matching** — connect students and juniors with graduates in their field
- 📅 **Events & reunions** — announcements for faculty meetups and alumni gatherings
- 📊 **Admin dashboard** — review verification requests, moderate reported posts, and read platform statistics

## Tech Stack

| Layer | Technology | Why |
| --- | --- | --- |
| **Backend** | [Go](https://go.dev) | Fast, statically typed, and comfortable with many concurrent requests on a small server |
| **Database** | [PostgreSQL](https://www.postgresql.org) | Relational integrity for graduates, departments, and employment records; full-text search built in, so no separate search engine is needed early on |
| **Cache & sessions** | [Redis](https://redis.io) | Session storage, cached search results, and rate limiting — all read constantly and changed rarely |
| **API** | REST (JSON) | Simple, predictable, and easy for any frontend to consume |
| **Containerization** | Docker & Docker Compose | One command brings Postgres and Redis up locally, identically on every machine |

The frontend is deliberately undecided. The API is being designed first so that whatever consumes it — a web app, a mobile app, or both — is free to change later.

## Architecture

```
          ┌──────────────┐
          │    Client    │
          └──────┬───────┘
                 │ HTTPS / JSON
          ┌──────▼───────┐
          │    Go API    │   auth · profiles · search · feed · messaging
          └───┬──────┬───┘
              │      │
   ┌──────────▼─┐  ┌─▼───────────┐
   │ PostgreSQL │  │    Redis    │
   │  source of │  │  sessions,  │
   │    truth   │  │ cache, rate │
   └────────────┘  └─────────────┘
```

PostgreSQL is the single source of truth: every profile, employment record, post, and verification decision lives there and survives a restart. Redis holds only data that can be rebuilt — sessions, cached search results, assembled feed pages, rate-limit counters — so losing the cache degrades speed, never correctness.

The feed is the one place where that split really earns itself. Building a graduate's timeline means reading from everyone they follow, on every page load, for every user at once — exactly the query that gets expensive first. Caching assembled pages in Redis keeps the feed fast, and because the posts themselves still live in Postgres, a flushed cache costs a few slow requests rather than a lost conversation.

## What Works Today

The API runs on `http://localhost:8080` (`docker compose up --build`, see
[Getting Started](#getting-started)). Users, announcements and alumni are stored in PostgreSQL,
so they survive a restart and `docker compose down`.

| Route | Response |
| --- | --- |
| `GET /` | the HTML landing page |
| `GET /api/health` | `{"status":"ok"}` — is the server up? |
| `GET /api/users` | every user, as a JSON array |
| `POST /api/users` | creates a user; `201` with the assigned id |
| `GET /api/users/{id}` | one user, or `404` |
| `DELETE /api/users/{id}` | removes a user; `204` with an empty body |
| `PUT /api/users/{id}` | replaces the whole user; every field required |
| `PATCH /api/users/{id}` | changes only the fields the body mentions |
| `GET /api/swagger` | Swagger UI, rendered from the OpenAPI document |
| `GET /api/swagger.json` | the OpenAPI document itself |
| `GET /main` | HTML page listing every route |
| `GET /about` | HTML about page (placeholder content) |
| `GET /users` | HTML page listing every user, with edit and delete buttons and a form that posts to `POST /users` |
| `GET /users/new` | HTML form for a new user |
| `POST /users` | creates a user from the form, then redirects (`303`) to `/users` |
| `GET /users/{id}` | HTML page for one user, or `404` |
| `GET /users/{id}/edit` | HTML form filled in with the user's current values |
| `POST /users/{id}` | saves the edit form, then redirects to `/users/{id}` |
| `POST /users/{id}/delete` | removes the user, then redirects to `/users` |
| `GET /announcements` | HTML page listing every announcement (newest first), with a form that posts to `POST /announcements` |
| `GET /announcements/new`, `GET /announcements/{id}`, `GET /announcements/{id}/edit` | empty form, one announcement, edit form |
| `POST /announcements`, `POST /announcements/{id}`, `POST /announcements/{id}/delete` | create, save an edit, delete; each redirects (`303`) |
| `GET /api/announcements`, `GET /api/announcements/{id}` | every announcement, or one, as JSON |
| `POST /api/announcements` | creates an announcement; `201` with the id and `createdAt` the server assigned |
| `PUT /api/announcements/{id}`, `DELETE /api/announcements/{id}` | replaces title, body and author; deletes (`204`) |
| `GET /alumni` | every graduate, as a JSON array |
| `POST /alumni` | creates a graduate; `201` with the assigned id |
| `GET /alumni/{id}`, `PUT /alumni/{id}`, `DELETE /alumni/{id}` | one graduate; replace name, department and year; delete (`204`) |
| `GET /graduates` | HTML page listing every graduate, with edit and delete buttons and a create form |
| `GET /graduates/new`, `GET /graduates/{id}`, `GET /graduates/{id}/edit` | empty form, one graduate, edit form |
| `POST /graduates`, `POST /graduates/{id}`, `POST /graduates/{id}/delete` | create, save an edit, delete; each redirects (`303`) |
| `GET /hello`, `GET /hello/{name}` | greeting, from the lecture exercises |
| `GET /sum/{number1}/{number2}` | the sum, or `400` on non-numeric input |
| `GET /temporary` | `307` temporary redirect to `/main` |

Routing uses the Go 1.22 standard-library `ServeMux` — no third-party router. HTML lives in
`internal/view/templates/`, embedded into the binary with `go:embed` and rendered through
`html/template`; styling is [Pico.css](https://picocss.com), so the markup stays plain HTML.

`requests.http` at the repository root fires every endpoint, including the error cases, from
the VS Code REST Client extension.

### Keep the API documentation current

`internal/view/openapi.json` describes the API in the [OpenAPI 3](https://swagger.io/specification/)
format, and `GET /api/swagger` renders it as a Swagger UI page you can send requests from.

**Every pull request that changes an endpoint must update `openapi.json` in the same commit.**
That means a new route, a removed one, a renamed field, a different status code — anything a
caller would notice. Documentation that is updated "later" is documentation that quietly starts
lying, and a wrong API description is worse than none: it is believed. Treat the spec as part of
the endpoint, not as a chore that follows it.

A quick way to check yourself before opening a PR: open `/api/swagger`, press **Try it out** on
each endpoint you touched, and confirm the real response matches what the page promises.

## Data Model

> 📄 The tables that exist **today** (columns, types, rules, and why the data survives
> `docker compose down`) are documented in **[DDD.md — Data Design Documentation](DDD.md)**.

The initial schema sketch below is where the platform is heading. It will change as the code is written.

| Table | Holds | Notes |
| --- | --- | --- |
| `users` | login credentials, email, role | one row per account |
| `alumni_profiles` | faculty, department, graduation year, city, bio | one-to-one with `users` |
| `faculties` / `departments` | the university's own structure | reference data, seeded once |
| `employments` | company, title, start/end date | many per profile — this is the career timeline |
| `verifications` | proof submitted, reviewer, decision | an audit trail, not a single boolean |
| `posts` | body, author, optional faculty/department circle | the feed; a job opening is a post with a type |
| `comments` | body, author, parent post | threaded replies come later, flat first |
| `reactions` | who reacted to what, and how | one row per (user, post) pair |
| `follows` | follower → followed (a person or a circle) | what a graduate's feed is assembled from |
| `help_offers` | what a graduate is willing to do: CV review, mock interview, a coffee | small, explicit asks are the ones people say yes to |
| `reports` | reported post, reporter, reason | moderation needs a queue, not ad-hoc deletes |

## Prior Art

Almost every serious university already runs something in this space. What they do — and where they stop — is the clearest specification available.

| Institution | What it does well | Where it stops |
| --- | --- | --- |
| **Istanbul University** — [Mezun Bilgi Sistemi](https://mezun.istanbul.edu.tr) & [Mezun Doğrulama Sistemi](https://dogrulama.istanbul.edu.tr) | Official graduate records, profile and CV fields, and a real diploma-verification service backed by YÖKSİS | A records system, not a community: graduates have no reason to return after filling the form once |
| **METU** ([mezun.metu.edu.tr](https://mezun.metu.edu.tr/en/)) | The alumni card is a genuine reason to register — library, cafeteria, sports facilities, campus access, partner discounts. Mentoring and speed-networking events | Networking largely happens at physical events, not continuously online |
| **Boğaziçi** ([mbs.boun.edu.tr](https://mbs.boun.edu.tr/)) | A lifelong `@alumni.bogazici.edu.tr` address that doubles as the login, plus LinkedIn sign-in; digital alumni card | Closed portal; little public evidence of an ongoing feed |
| **İTÜ** ([mezun.itu.edu.tr](https://mezun.itu.edu.tr/)) | Alumni card, İTÜ Day, e-bulletin, and links to the alumni foundations and associations | Events and newsletters — one-way communication, no job board or mentoring in the portal |
| **Stanford** ([alumni.stanford.edu](https://alumni.stanford.edu/help/directory/)) | The strongest directory design: filters by class-year range, area of study, region, industry, skills, employer, and *whether the person is open to helping* — plus per-profile privacy that hides you from search entirely | Closed to non-alumni by design |
| **Harvard** ([alumni.harvard.edu](https://alumni.harvard.edu/community/alumni-services)) | 50+ Shared Interest Groups organised around purpose rather than class or faculty; "flash mentoring" — one mock interview or one CV review, not a six-month commitment | Heavy institutional infrastructure behind it |
| **MIT** ([Infinite Connection](https://alum.mit.edu/about/benefits-and-offerings/infinite-connection)) | Email-for-life at `@alum.mit.edu`, an alumni job board, and a directory searchable by region, industry, course, and living group — one account for everything | — |

**What this project takes from them**

1. **Give people a reason to register that is not altruism.** Every Turkish system above is built around the alumni card. A benefit you can hold is what converts a graduate into a user; the social features only matter afterwards.
2. **Verification is the foundation, not a feature.** MIT and Boğaziçi turn it into a lifelong email address, so the credential and the benefit are the same object.
3. **"Open to helping" belongs in the schema.** Stanford filters on it, Harvard sizes the ask down to a single CV review. Cheap for the mentor, decisive for the student.
4. **Interest groups beat class years alone.** Harvard's SIGs exist because graduates have more in common than a graduation date.
5. **The gap worth filling.** Turkish alumni systems are records and event announcements — one-way. None of them is a place where graduates talk to each other daily. That is exactly where this project aims.

## Open Design Questions

Honest unknowns, written down so they are decided deliberately rather than by accident:

- **How is "is this person really a graduate?" answered?** Student addresses (`@ogr.iu.edu.tr`) stop working after graduation, so email-domain checks alone cannot work for the people this platform is for. The realistic options, cheapest first: a graduation document from **e-Devlet**, whose verification code anyone can re-check; the university's own [Mezun Doğrulama Sistemi](https://dogrulama.istanbul.edu.tr), which validates a diploma against YÖKSİS records (manual web form today — no public API, so a human reviewer sits in the loop until the university offers an integration); a lifelong `@alumni.istanbul.edu.tr` address, which would be the best answer but needs the university to issue it; or vouching by already-verified graduates, which bootstraps quickly and decays quickly.
- **Who may see contact details?** Open to every verified graduate, or only after both sides accept a connection?
- **What happens to a profile its owner abandons?** Stale career data is worse than no career data.
- **How is the feed ordered?** Newest-first is honest and trivial to build. Any ranking beyond that needs a reason, and ranking a small network too aggressively just hides most of it.
- **Who moderates the feed?** A closed, verified network needs far less moderation than a public one — but "far less" is not "none", and one unhandled report is enough to make people stop posting.
- **KVKK / personal data.** Personal data of real people is involved; consent, retention, and deletion need to be designed in, not bolted on.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- [Go](https://go.dev/dl/) 1.22 or newer — only to run the app outside Docker

### Run it

```bash
git clone https://github.com/rabizd/alumni.git
cd alumni
cp .env.example .env                                  # settings: ports, database user and name
echo "pick-any-password" > secrets/db_password.txt   # the database password
docker compose up --build -d
docker compose run --rm seed                          # fill the empty tables with sample data
```

Then open <http://localhost:8080/main>. Compose starts `db` (PostgreSQL), runs `migrate` to
create the tables, then starts `app` (the Go server). The tables start empty; `seed` fills them
(see [Factories and seeding](#factories-and-seeding)).

`.env` and `secrets/db_password.txt` must exist before the first `docker compose up`, otherwise
Compose stops with an error that names what is missing. Neither file is committed. The password is set when the database is created for the first time;
changing the file later does not change it (run `docker compose down -v` to start over).

To run the Go code without Docker while the database container is up:
`docker compose up -d db`, then `go run ./cmd/api`.

### The database survives `docker compose down`

PostgreSQL keeps its files in the named volume `db-data`. `docker compose down` deletes the
containers but not the volume, so the data is still there on the next `up`:

1. Create a record, for example `POST /api/users` (see `requests.http`)
2. `docker compose down`
3. `docker compose up -d`
4. `GET /api/users` — the record is still there

Only `docker compose down -v` deletes the volume, and with it every record. On the next start
the migrations run again from the first one and bring back the empty tables; run `seed` to fill them.

### Migrations

The database schema is built by **migrations**: numbered SQL files in `migrations/`, applied in
order by [golang-migrate](https://github.com/golang-migrate/migrate). Each one has two halves:

| File | Does |
| --- | --- |
| `000001_create_users.up.sql` | applies the change (`CREATE TABLE users …`) |
| `000001_create_users.down.sql` | undoes it (`DROP TABLE users`) |

On every `docker compose up`, the `migrate` service runs first, applies only the migrations the
database has not seen yet, and exits; the app starts after it succeeds. The database remembers
the last applied number in its `schema_migrations` table, so an `up` with nothing new prints
`no change` and the data is untouched.

**Changing the schema** never needs `down -v`: add the next pair of files, for example

```
migrations/000004_add_city_to_users.up.sql     ALTER TABLE users ADD COLUMN city TEXT NOT NULL DEFAULT '';
migrations/000004_add_city_to_users.down.sql   ALTER TABLE users DROP COLUMN city;
```

and run `docker compose up -d` (or `docker compose run --rm migrate up`). The existing rows stay.

| Command | What it does |
| --- | --- |
| `docker compose run --rm migrate up` | apply every new migration |
| `docker compose run --rm migrate down 1` | undo the last migration |
| `docker compose run --rm migrate version` | print the number of the last applied migration |

> ⚠️ Never edit a migration that has already run — the database will not run it again. Add a
> new one instead.

Migrations hold only the **structure**. The rows come from the seeder.

### Factories and seeding

| Part | File | What it does |
| --- | --- | --- |
| **Factory** | `internal/seed/factory.go` | makes one random but valid record: `User()`, `Announcement()`, `Alumni()`, `JobPosting(userID)`. Names, departments and years are picked from Turkish lists; every user passes `User.Validate()` (start year in range, graduation after start, never in the future) and gets a unique email |
| **Seeder** | `internal/seed/seeder.go` | puts the fixed samples (Ayşe Yılmaz, Mehmet Demir, two announcements, two job postings) into any **empty** table, then adds the random records. Job postings come last, each by a random user already in the database. Each record goes through `Validate()` and the normal store, like an API request |
| **Command** | `cmd/seed/main.go` | runs the seeder; the `seed` service in `docker-compose.yml` runs it in Docker |

| Command | What it does |
| --- | --- |
| `docker compose run --rm seed` | samples (if the table is empty) + 10 users, 5 announcements, 10 alumni, 10 job postings |
| `docker compose run --rm seed -users 50 -alumni 0 -jobs 20` | choose how many random records of each kind |
| `docker compose run --rm seed -fresh` | empty every table first (ids start from 1 again), then seed |
| `docker compose run --rm seed -seed 42` | a fixed random seed: the same number always makes the same records |

The `seed` service has the `tools` profile, so `docker compose up` does not start it; it only
runs when you ask. Running it twice adds more random records but never a second copy of the
samples.

### Looking inside the database

**In the browser — Adminer.** `docker compose up` also starts Adminer, a small database viewer.
Open <http://localhost:8081> and log in:

| Field | Value |
| --- | --- |
| System | `PostgreSQL` |
| Server | `db` |
| Username | `DB_USER` from `.env` (`alumni`) |
| Password | the contents of `secrets/db_password.txt` |
| Database | `DB_NAME` from `.env` (`alumni`) |

Click a table, then **Select data**, to see its rows. Edit a graduate at `/graduates`, refresh
Adminer, and the change is there. Adminer listens on `127.0.0.1` only, so other computers on the
network cannot reach it.

**In the terminal — psql.** `docker compose exec db psql -U alumni -d alumni`, then SQL such as
`SELECT * FROM alumni;`. `\dt` lists the tables, `\q` quits.

**In the log — every change.** The `db` service runs Postgres with `log_statement=mod`, so each
`INSERT`, `UPDATE` and `DELETE` is written to its log together with the values, for example
`DELETE FROM alumni WHERE id = $1` / `Parameters: $1 = '15'`. Watch it with
`docker compose logs -f db` (or the **Logs** tab of `alumni-db-1` in Docker Desktop). Reads
(`SELECT`) are not logged, so the log stays short.

## Configuration

Settings live in a **`.env`** file in the project root. Docker Compose reads it automatically
and puts the values into `docker-compose.yml` wherever it says `${...}`.

| File | Committed? | What it is |
| --- | --- | --- |
| `.env.example` | ✅ yes | the template: every variable, with safe example values and a comment |
| `.env` | ❌ no, in `.gitignore` | your own copy (`cp .env.example .env`), with your values |

Each machine has its own `.env`, so changing a port on your laptop never ends up in someone
else's checkout. When a new variable is added, it goes into `.env.example` too.

| Variable (in `.env`) | Description | Example |
| --- | --- | --- |
| `APP_PORT` | Port the API listens on, inside the container and on your machine | `8080` |
| `DB_USER` | PostgreSQL user | `alumni` |
| `DB_NAME` | PostgreSQL database name | `alumni` |
| `DB_PORT` | Port PostgreSQL is reachable on from your machine | `5432` |
| `ADMINER_PORT` | Port of Adminer, the database viewer (optional; `8081` if unset) | `8081` |

Compose also passes two fixed values to the app: `DB_HOST=db` (the database container's name)
and `DB_PASSWORD_FILE=/run/secrets/db_password`. When you `go run ./cmd/api` outside Docker, the
app falls back to `localhost`, `alumni` and `secrets/db_password.txt`.

The password is a **Docker Compose secret**, not an environment variable. Compose mounts
`secrets/db_password.txt` into both containers as `/run/secrets/db_password`; Postgres reads it
through `POSTGRES_PASSWORD_FILE` and the app through `DB_PASSWORD_FILE`. Environment variables
show up in `docker inspect` and in logs; a secret file does not.

> 🔒 Never commit `.env` or `secrets/db_password.txt`. `.gitignore` keeps both out of git, and
> `.dockerignore` keeps both out of the image. Only `.env.example` belongs in version control.

*Planned:* `REDIS_URL` for sessions and `JWT_SECRET` for signing tokens.

## Project Structure (MVC)

The code follows the **Model–View–Controller** pattern, and each layer has its own folder
under `internal/`:

- **Model** (`internal/model`) — the data and the rules that keep it valid: what a user or a graduate *is*, how it is checked, and how it is stored, found, changed, and removed. It knows nothing about HTTP.
- **View** (`internal/view`) — what the client receives: an HTML page or a JSON body. A view shows data and decides nothing.
- **Controller** (`internal/controller`) — one function per route. It reads the request (URL, path values, JSON body), checks it, asks the model for data, and hands the result to a view.

Routes sit outside the three layers, in `internal/routes`: `web.go` maps the HTML routes and
`api.go` the JSON routes to their controllers. `cmd/api/main.go` registers both and starts the server.

### How one request moves through the app

```
  Client
    │  GET /api/users/1
    ▼
┌───────────────────────────────────────┐
│ Router      internal/routes/api.go    │  ServeMux matches method + path
└───────────────────┬───────────────────┘
                    ▼
┌───────────────────────────────────────┐
│ Controller  internal/controller/      │  ApiUserController.Show: read {id},
│             api_user_controller.go    │  400 if it is not a number
└───────────────────┬───────────────────┘
                    ▼
┌───────────────────────────────────────┐
│ Model       internal/model/user.go    │  model.Users.Find(1) → User, found?
└───────────────────┬───────────────────┘
                    ▼
┌───────────────────────────────────────┐
│ View        internal/view/view.go     │  view.JSON(...) → {"id":1,"name":...}
│                                       │  view.HTML(...) → templates/*.html
└───────────────────┬───────────────────┘
                    ▼
                 Client
```

The arrows only point one way. A controller imports the model and the view; the model and the
view never import a controller, and they never import each other.

### Directories and files

```
alumni/
├── cmd/
│   ├── api/
│   │   └── main.go             # entry point: registers the routes, reads APP_PORT, starts the server
│   └── seed/
│       └── main.go             # seed command: -users, -announcements, -alumni, -fresh, -seed
│
├── internal/                   # the application, split into the three MVC layers
│   ├── routes/                 # ROUTES: which URL goes to which controller
│   │   ├── web.go              #   Web(): pages, UserController at /users, AnnouncementController at /announcements  (HTML)
│   │   └── api.go              #   API(): ApiUserController, ApiAnnouncementController, health, swagger, /alumni  (JSON)
│   │
│   ├── model/                  # MODEL
│   │   ├── db.go               #   Connect(): opens the PostgreSQL pool, reads the password from the secret file
│   │   ├── user.go             #   User (+ department, years, ÇAP, yandal, advisor, prep exemption) + Validate(), UserPatch, UserStore (List/Add/Find/Remove/Replace), Apply()
│   │   ├── announcement.go     #   Announcement + Validate(), AnnouncementStore (List/Add/Find/Replace/Remove)
│   │   ├── alumni.go           #   Alumni + Validate(), AlumniStore (List/Add/Find/Replace/Remove)
│   │   └── health.go           #   Health, the {"status":"ok"} shape
│   │
│   ├── seed/                   # FACTORIES + SEEDER
│   │   ├── factory.go          #   random valid User / Announcement / Alumni
│   │   └── seeder.go           #   Run(): samples into empty tables, then random records
│   │
│   ├── view/                   # VIEW
│   │   ├── view.go             #   HTML() renders a template with data, JSON() writes a JSON body, OpenAPI() serves the spec
│   │   ├── openapi.json        #   the API description (embedded with go:embed)
│   │   └── templates/          #   HTML pages (embedded with go:embed)
│   │       ├── main.html       #     landing page, lists every route
│   │       ├── about.html      #     about page
│   │       ├── swagger.html    #     Swagger UI, loads /api/swagger.json
│   │       ├── layout.html     #     shared top and bottom of every page, and the form fields of each resource
│   │       ├── announcements.html #  announcement list + create form
│   │       ├── announcement.html #   one announcement
│   │       ├── announcement_form.html # create and edit form
│   │       ├── graduates.html  #     graduate list + create form
│   │       ├── graduate.html   #     one graduate
│   │       ├── graduate_form.html # create and edit form
│   │       ├── users.html      #     user list
│   │       ├── user.html       #     one user
│   │       └── user_form.html  #     create and edit form
│   │
│   └── controller/             # CONTROLLER
│       ├── pages.go            #   Root, Main, About, Hello, HelloName, Sum, Temporary
│       ├── api_user_controller.go # ApiUserController: JSON at /api/users  (Index, Show, Store, Update, Patch, Destroy)
│       ├── user_controller.go  #   UserController: HTML at /users  (Index, Show, Create, Store, Edit, Update, Destroy)
│       ├── api_announcement_controller.go # ApiAnnouncementController: JSON at /api/announcements  (Index, Show, Store, Update, Destroy)
│       ├── announcement_controller.go     # AnnouncementController: HTML at /announcements  (Index, Show, Create, Store, Edit, Update, Destroy)
│       ├── alumni.go           #   JSON at /alumni: ListAlumni, CreateAlumni, ShowAlumni, UpdateAlumni, DeleteAlumni
│       ├── graduate_controller.go # GraduateController: HTML at /graduates  (Index, Show, Create, Store, Edit, Update, Destroy)
│       ├── api.go              #   Health, Swagger, SwaggerSpec
│       └── request.go          #   shared helpers: pathID (reads {id}), decodeJSON (strict body parsing), dbError (logs, answers 500)
│
├── migrations/                 # numbered schema changes, NNNNNN_name.up.sql + .down.sql, applied in order
├── secrets/                    # db_password.txt lives here (not committed)
│
├── requests.http               # every endpoint, for the VS Code REST Client (manual testing)
├── Dockerfile                  # two-stage build: compile in golang, run in a small alpine image
├── docker-compose.yml          # app on port 8080, migrate (runs migrations/), PostgreSQL with the db-data volume and the password secret
├── .dockerignore               # files kept out of the Docker build
├── .env.example                # template for .env: copy it, then change the values
├── .gitignore                  # keeps .env and secrets/ out of git
├── go.mod                      # Go module: github.com/rabizd/alumni, Go 1.22, pgx (PostgreSQL driver)
├── README.md
├── DDD.md                      # Data Design Documentation: tables, columns, rules
└── LICENSE
```

Go treats a folder named `internal/` specially: only code inside this module may import it.
The layers are the app's own parts, not a library for other projects.

### The same thing, layer by layer

| Layer | Folder | What it contains | What it must not do |
| --- | --- | --- | --- |
| **Model** | `internal/model` | Structs (`User`, `UserPatch`, `Alumni`, `Health`), their `Validate()` rules (trimming spaces, required fields, the PATCH checks), and stores that run the SQL queries. The database assigns ids; a request body never does. | Touch `http.Request`, write a response, or know about HTML |
| **View** | `internal/view` | `HTML()` for pages, `JSON()` for API bodies, `OpenAPI()` for the spec, plus the templates and `openapi.json` | Check input or read and change the stores |
| **Controller** | `internal/controller` | One function per route; the two user controllers are structs whose methods are the CRUD actions. It returns `400` on bad input, `404` on a missing id and `500` when the database fails (the real error goes to the log), and otherwise calls the model and then a view. | Hold data itself, or build HTML or JSON by hand |

### The User model: CRUD on PostgreSQL

`internal/model/user.go` is the User model. `model.Users` reads and writes the `users` table;
`migrations/000001_create_users.up.sql` creates it and the seeder fills it. Each CRUD operation is one method, and each
also returns an `error` when the database fails (the controllers answer `500`):

| CRUD | Method on `model.Users` | What it does | Route that uses it |
| --- | --- | --- | --- |
| **Create** | `Add(u User)` | `INSERT`; the database gives the next id | `POST /api/users` |
| **Read** | `List()` | every user, by id | `GET /api/users` |
| **Read** | `Find(id int)` | one user; `false` if the id does not exist | `GET /api/users/{id}` |
| **Update** | `Replace(id int, u User)` | `UPDATE` of every field, keeps the id | `PUT /api/users/{id}`, and `PATCH` after `Find` + `Apply` |
| **Delete** | `Remove(id int)` | `DELETE`; `false` if the id does not exist | `DELETE /api/users/{id}` |

### Two controllers for the same model

The routes that reach them are defined in `internal/routes`: `web.go` sends `/users` to
`UserController`, and `api.go` sends `/api/users` to `ApiUserController`. Both controllers are
in the Swagger document at `/api/swagger`, each under its own tag, and every operation's
`operationId` is the controller method that answers it (for example `ApiUserController.Show`).

The users are served two ways, by two controllers in `internal/controller`. Both call the same
`model.Users` methods; only the input and the output differ.

| CRUD | `ApiUserController` (JSON, for programs) | `UserController` (HTML, for people) |
| --- | --- | --- |
| **Create** | `Store` — `POST /api/users` → `201` + JSON | `Create` — `GET /users/new` shows the form; `Store` — `POST /users` saves it |
| **Read** | `Index` — `GET /api/users`; `Show` — `GET /api/users/{id}` | `Index` — `GET /users`; `Show` — `GET /users/{id}` |
| **Update** | `Update` — `PUT /api/users/{id}`; `Patch` — `PATCH /api/users/{id}` | `Edit` — `GET /users/{id}/edit` shows the form; `Update` — `POST /users/{id}` saves it |
| **Delete** | `Destroy` — `DELETE /api/users/{id}` → `204` | `Destroy` — `POST /users/{id}/delete` |

`ApiUserController` reads a JSON body and answers with JSON. `UserController` reads an HTML
form and answers with a page, or with a `303` redirect after a successful POST so that reloading
the page does not submit the form twice. HTML forms can only send `GET` and `POST`, which is why
the HTML side uses `POST` where the API uses `PUT`, `PATCH` and `DELETE`.

Before Create and Update, the controller calls `User.Validate()` or `UserPatch.Validate()`.
These trim spaces and enforce the rules below.

| Field | JSON / form name | Rule |
| --- | --- | --- |
| Name, e-mail | `name`, `email` | required |
| Department (okuduğu bölüm) | `department` | required |
| Start year (başlangıç yılı) | `startYear` | required, between 1900 and this year |
| Graduation year (mezuniyet yılı) | `graduationYear` | required, not before the start year, not in the future |
| Double major (ÇAP) | `doubleMajor` | optional, empty means none |
| Minor (yandal) | `minor` | optional, empty means none |
| Advisor (danışman) | `advisor` | optional |
| English prep exemption (İngilizce hazırlık muafiyeti) | `prepExempt` | `true` / `false`; in the HTML form a checkbox |

A PATCH is applied first and the patched user is validated as a whole, so changing only
`graduationYear` to a year before the existing `startYear` is still rejected, and nothing is saved.

### Announcements: the same pattern again

`internal/model/announcement.go` is the Announcement model. It has a title, a body, an author,
and a `createdAt` time that the database sets. `model.Announcements` reads and writes the
`announcements` table, which `migrations/000002_create_announcements.up.sql` creates and the seeder fills. `List()` returns the newest first.

| CRUD | Model (`model.Announcements`) | `ApiAnnouncementController` (JSON) | `AnnouncementController` (HTML, the management interface) |
| --- | --- | --- | --- |
| **Create** | `Add` | `Store` — `POST /api/announcements` → `201` | `Create` — `GET /announcements/new`; `Store` — `POST /announcements` |
| **Read** | `List`, `Find` | `Index` — `GET /api/announcements`; `Show` — `GET /api/announcements/{id}` | `Index` — `GET /announcements`; `Show` — `GET /announcements/{id}` |
| **Update** | `Replace` (keeps the id and `createdAt`) | `Update` — `PUT /api/announcements/{id}` | `Edit` — `GET /announcements/{id}/edit`; `Update` — `POST /announcements/{id}` |
| **Delete** | `Remove` | `Destroy` — `DELETE /api/announcements/{id}` → `204` | `Destroy` — `POST /announcements/{id}/delete` |

Title, body and author are required; `Announcement.Validate()` trims them and rejects any that
are empty. The management interface is at `/announcements`, and the navigation bar on every
page links to `/users`, `/announcements` and `/graduates`.

### Graduates: the alumni table, with both interfaces

The `alumni` table has the same full CRUD. Its JSON stays at `/alumni`, the address from the
lecture exercise, and its HTML management interface is at `/graduates`.

| CRUD | Model (`model.Graduates`) | JSON (`internal/controller/alumni.go`) | `GraduateController` (HTML) |
| --- | --- | --- | --- |
| **Create** | `Add` | `CreateAlumni` — `POST /alumni` → `201` | `Create` — `GET /graduates/new`; `Store` — `POST /graduates` |
| **Read** | `List`, `Find` | `ListAlumni` — `GET /alumni`; `ShowAlumni` — `GET /alumni/{id}` | `Index` — `GET /graduates`; `Show` — `GET /graduates/{id}` |
| **Update** | `Replace` | `UpdateAlumni` — `PUT /alumni/{id}` | `Edit` — `GET /graduates/{id}/edit`; `Update` — `POST /graduates/{id}` |
| **Delete** | `Remove` | `DeleteAlumni` — `DELETE /alumni/{id}` → `204` | `Destroy` — `POST /graduates/{id}/delete` |

Only the name is required; department and graduation year may be empty.

### When the database fails

Every controller sends database errors through `dbError` in `request.go`. The client gets
`500 database error` and nothing more, because the details are not its business; the real
error is written to the app's log, for example
`database error: failed to connect to … lookup db: no such host` when the `db` container is
stopped. Read it with `docker compose logs app`.

### What is still on the way

- **No Redis yet.** Sessions and caching arrive with sign-in.

## Roadmap

**Phase 1 — foundations**
- [x] Go module and an HTTP server with routing
- [x] HTML landing / about pages, and `GET /api/health` as a JSON health check
- [x] `GET` and `POST /alumni` against an in-memory store
- [x] `/api/users` with `GET`, `POST`, `PUT`, `PATCH` and `DELETE`, still in memory
- [x] OpenAPI document and a Swagger UI page at `/api/swagger`
- [x] Code split into `internal/model`, `internal/view`, and `internal/controller` (MVC)
- [x] Announcements: in-memory model, `AnnouncementController` + `ApiAnnouncementController`, and a management interface at `/announcements`
- [x] `docker-compose.yml` with PostgreSQL, a named volume and the password as a secret
- [x] Users, announcements and alumni moved from memory onto PostgreSQL
- [x] Database migrations with golang-migrate (`migrations/`, run by the `migrate` service)
- [ ] Redis in `docker-compose.yml`

**Phase 2 — a usable core**
- [ ] Sign-up, sign-in, and Redis-backed sessions
- [ ] Alumni profile create / read / update
- [ ] Graduate verification flow
- [ ] Search and filtering

**Phase 3 — the social layer**
- [ ] Posts and a newest-first feed
- [ ] Follows, reactions, and comments
- [ ] Faculty / department circles and class-year pages
- [x] Job & internship board: postings, applications (`/job-postings`, `/applications`)
- [ ] "Ağımızdan yüzler" on the main page: a strip of graduates' profile photos, like Stanford's alumni site. Only graduates who opt in, so it needs sign-in and photo upload first
- [ ] Direct messaging
- [ ] Reporting and moderation queue

**Phase 4 — the network effects**
- [ ] Company / institution directory
- [ ] Mentorship matching
- [ ] Events and reunions
- [ ] Admin dashboard

## Contributing

Contributions are welcome — this project grows faster with more hands.

1. **Fork** this repository.
2. Create a branch: `git checkout -b feature/your-feature`
3. Commit your changes: `git commit -m "Add your feature"`
4. Push the branch: `git push origin feature/your-feature`
5. Open a **Pull Request**.

Before step 5, make sure `internal/view/openapi.json` reflects any endpoint you added or changed —
see [Keep the API documentation current](#keep-the-api-documentation-current).

## License

Distributed under the [MIT](LICENSE) license.

## Contact

**Project owner:** [@rabizd](https://github.com/rabizd)

**Project link:** [https://github.com/rabizd/alumni](https://github.com/rabizd/alumni)

<div align="center">
<sub>Built for the graduates of Istanbul University 🎓</sub>
</div>
