<div align="center">

# 🎓 Alumni — Istanbul University Alumni Network

**A social network for Istanbul University graduates — part directory, part feed, entirely closed to everyone else.**

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
- [The Social Layer](#the-social-layer)
- [Planned Features](#planned-features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Data Model](#data-model)
- [Prior Art](#prior-art)
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

## Data Model

The initial schema sketch. It will change as the code is written.

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

**Phase 3 — the social layer**
- [ ] Posts and a newest-first feed
- [ ] Follows, reactions, and comments
- [ ] Faculty / department circles and class-year pages
- [ ] Job & internship board
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

## License

Distributed under the [MIT](LICENSE) license.

## Contact

**Project owner:** [@rabizd](https://github.com/rabizd)

**Project link:** [https://github.com/rabizd/alumni](https://github.com/rabizd/alumni)

<div align="center">
<sub>Built for the graduates of Istanbul University 🎓</sub>
</div>
