<div align="center">

# 🎓 Alumni — Istanbul University Alumni Network

**A platform that keeps Istanbul University graduates connected — to each other, and to their university.**

[![Status](https://img.shields.io/badge/status-early%20development-orange)](https://github.com/rabizd/alumni)
[![Go](https://img.shields.io/badge/backend-Go-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/database-PostgreSQL-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/cache-Redis-DC382D?logo=redis&logoColor=white)](https://redis.io)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

> ⚠️ **Early development.** No application code has been written yet — this README describes the project that is being built. Everything below marked *planned* is a design decision, not a shipped feature.

## Table of Contents

- [About](#about)
- [Why This Project](#why-this-project)
- [Planned Features](#planned-features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Data Model](#data-model)
- [Open Design Questions](#open-design-questions)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Planned Project Structure](#planned-project-structure)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About

**Alumni** is a web platform built for the graduates of **Istanbul University**. Graduation should not be the end of a community — it should be the beginning of a much larger one. Yet in practice, the moment students walk out of the gates of Beyazıt, the network they spent years building quietly dissolves into scattered group chats and outdated spreadsheets.

Alumni brings that network back into one place. It collects graduate profiles, career histories, and the institutions they work at, then makes that information **searchable, discoverable, and genuinely useful** — whether you are a recent graduate looking for a mentor, a senior professional hiring from your own faculty, or the university itself trying to stay in touch with the people it educated.

## Why This Project

- **The network already exists — it is just invisible.** Thousands of Istanbul University graduates work across the country and the world. Alumni makes that reality legible.
- **Careers move fast; contact lists do not.** Profiles that graduates maintain themselves stay alive without anyone curating a spreadsheet.
- **Mentorship needs a starting point.** Finding "someone from my department who does what I want to do" should take one search, not six months of luck.
- **Institutional memory has value.** A university that knows where its graduates went can advise the students who have not left yet.

## Planned Features

**Core — the first version is not useful without these**

- 📇 **Alumni profiles** — faculty, department, graduation year, city, and current position
- 🔐 **Sign-up, sign-in, and sessions** — JWT access tokens, with sessions tracked in Redis
- ✅ **Graduate verification** — a profile is only trustworthy if the person really graduated (see [Open Design Questions](#open-design-questions))
- 🔍 **Search & filtering** — by name, faculty, department, graduation year, city, or employer
- 📝 **Career timeline** — positions and promotions added over time, so a profile ages well
- 🛡️ **Privacy controls** — every graduate decides which fields are public, which are visible to other verified alumni, and which stay private

**Later — valuable, but only once the core works**

- 🏢 **Company & institution directory** — see who works where, at a glance
- ✉️ **Direct messaging** — reach another graduate without exposing private contact details
- 🤝 **Mentorship matching** — connect students and juniors with graduates in their field
- 📅 **Events & reunions** — announcements for faculty meetups and alumni gatherings
- 📊 **Admin dashboard** — review verification requests and platform statistics

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
          │    Go API    │   auth · profiles · search · verification
          └───┬──────┬───┘
              │      │
   ┌──────────▼─┐  ┌─▼───────────┐
   │ PostgreSQL │  │    Redis    │
   │  source of │  │  sessions,  │
   │    truth   │  │ cache, rate │
   └────────────┘  └─────────────┘
```

PostgreSQL is the single source of truth: every profile, employment record, and verification decision lives there and survives a restart. Redis holds only data that can be rebuilt — sessions, cached search results, rate-limit counters — so losing the cache degrades speed, never correctness.

## Data Model

The initial schema sketch. It will change as the code is written.

| Table | Holds | Notes |
| --- | --- | --- |
| `users` | login credentials, email, role | one row per account |
| `alumni_profiles` | faculty, department, graduation year, city, bio | one-to-one with `users` |
| `faculties` / `departments` | the university's own structure | reference data, seeded once |
| `employments` | company, title, start/end date | many per profile — this is the career timeline |
| `verifications` | proof submitted, reviewer, decision | an audit trail, not a single boolean |

## Open Design Questions

Honest unknowns, written down so they are decided deliberately rather than by accident:

- **How is "is this person really a graduate?" answered?** Student email addresses (`@ogr.iu.edu.tr`) stop working after graduation, so email-domain verification alone cannot work for the people this platform is for. Candidate answers: diploma/transcript upload reviewed by an admin, a one-time invite code from the university, or vouching by already-verified graduates.
- **Who may see contact details?** Open to every verified graduate, or only after both sides accept a connection?
- **What happens to a profile its owner abandons?** Stale career data is worse than no career data.
- **KVKK / personal data.** Personal data of real people is involved; consent, retention, and deletion need to be designed in, not bolted on.

## Getting Started

> ⏳ The commands below describe the intended setup. They will not work until the API skeleton exists — see the [Roadmap](#roadmap).

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or newer
- [Docker](https://www.docker.com/) and Docker Compose — the easiest way to get PostgreSQL and Redis running
- Or, without Docker: [PostgreSQL](https://www.postgresql.org/download/) 15+ and [Redis](https://redis.io/download/) 7+ installed locally

### Installation

```bash
# Clone the repository
git clone https://github.com/rabizd/alumni.git
cd alumni

# Start PostgreSQL and Redis
docker compose up -d

# Copy the example environment file and fill in your values
cp .env.example .env

# Download Go dependencies and run the API
go mod download
go run ./cmd/api
```

The API will be available at `http://localhost:8080`.

## Configuration

Configuration is read from environment variables (see `.env.example`, once it exists):

| Variable | Description | Example |
| --- | --- | --- |
| `APP_PORT` | Port the API listens on | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/alumni?sslmode=disable` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379/0` |
| `JWT_SECRET` | Secret used to sign access tokens | *(a long random string)* |

> 🔒 Never commit your `.env` file. Only `.env.example` belongs in version control.

## Planned Project Structure

None of these directories exist yet; this is the layout the code is heading toward.

```
alumni/
├── cmd/
│   └── api/            # Application entry point
├── internal/
│   ├── handler/        # HTTP handlers (routing layer)
│   ├── service/        # Business logic
│   ├── repository/     # PostgreSQL and Redis access
│   └── model/          # Domain types
├── migrations/         # SQL schema migrations
├── docker-compose.yml  # PostgreSQL + Redis for local development
├── .env.example
├── README.md
└── LICENSE
```

## Roadmap

**Phase 1 — foundations**
- [ ] `docker-compose.yml` for PostgreSQL and Redis
- [ ] Go module, configuration loading, and a `/health` endpoint
- [ ] Database schema and migrations

**Phase 2 — a usable core**
- [ ] Sign-up, sign-in, and Redis-backed sessions
- [ ] Alumni profile create / read / update
- [ ] Graduate verification flow
- [ ] Search and filtering

**Phase 3 — the network effects**
- [ ] Company / institution directory
- [ ] Direct messaging
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

## License

Distributed under the [MIT](LICENSE) license.

## Contact

**Project owner:** [@rabizd](https://github.com/rabizd)

**Project link:** [https://github.com/rabizd/alumni](https://github.com/rabizd/alumni)

<div align="center">
<sub>Built for the graduates of Istanbul University 🎓</sub>
</div>
