<div align="center">

# 🎓 Alumni — Istanbul University Alumni Network

**A modern platform that keeps Istanbul University graduates connected — to each other, and to their university.**

[![Status](https://img.shields.io/badge/status-in%20development-orange)](https://github.com/rabizd/alumni)
[![Go](https://img.shields.io/badge/backend-Go-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/database-PostgreSQL-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/cache-Redis-DC382D?logo=redis&logoColor=white)](https://redis.io)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

> ⚠️ **Work in progress.** This repository is under active development. The feature set and architecture are still taking shape, and breaking changes should be expected.

## Table of Contents

- [About](#about)
- [Why This Project](#why-this-project)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About

**Alumni** is a web platform built for the graduates of **Istanbul University**. Graduation should not be the end of a community — it should be the beginning of a much larger one. Yet in practice, the moment students walk out of the gates of Beyazıt, the network they spent years building quietly dissolves into scattered group chats and outdated spreadsheets.

Alumni brings that network back into one place. It collects graduate profiles, career histories, and the institutions they work at, then makes that information **searchable, discoverable, and genuinely useful** — whether you are a recent graduate looking for a mentor, a senior professional hiring from your own faculty, or the university itself trying to stay in touch with the people it educated.

## Why This Project

- **The network already exists — it is just invisible.** Thousands of Istanbul University graduates work across the country and the world. Alumni makes that reality legible.
- **Careers move fast; contact lists do not.** Self-updating profiles keep the data alive without anyone maintaining a spreadsheet.
- **Mentorship needs a starting point.** Finding "someone from my department who does what I want to do" should take one search, not six months of luck.
- **Institutional memory has value.** A university that knows where its graduates went can advise the students who have not left yet.

## Features

Legend: ✅ done · 🚧 in progress · 📋 planned

- 📋 📇 **Alumni profiles** — contact details, faculty, department, and graduation year
- 📋 🏢 **Company & institution directory** — see who works where, at a glance
- 📋 🔍 **Rich search & filtering** — by name, faculty, department, graduation year, city, or employer
- 📋 📝 **Career timeline** — graduates add positions and promotions as their careers progress
- 📋 🔐 **Authentication & authorization** — secure sign-up and sign-in with JWT-based sessions
- 📋 ✉️ **Direct messaging** — reach another graduate without exposing private contact details
- 📋 🤝 **Mentorship matching** — connect students and juniors with graduates in their field
- 📋 📅 **Events & reunions** — announcements for faculty meetups and alumni gatherings
- 📋 🛡️ **Privacy controls** — every graduate decides which fields are public
- 📋 📊 **Admin dashboard** — verification of graduate records and platform statistics

## Tech Stack

| Layer | Technology | Why |
| --- | --- | --- |
| **Backend** | [Go](https://go.dev) | Fast, statically typed, and excellent at handling many concurrent requests with a small memory footprint |
| **Database** | [PostgreSQL](https://www.postgresql.org) | Relational integrity for graduates, departments, and employment records; powerful full-text search built in |
| **Cache & sessions** | [Redis](https://redis.io) | Sub-millisecond lookups for sessions, hot search results, and rate limiting |
| **API** | REST (JSON) | Simple, predictable, and easy for any frontend to consume |
| **Containerization** | Docker & Docker Compose | One command brings the whole stack up, identically on every machine |

## Architecture

```
          ┌──────────────┐
          │    Client    │
          └──────┬───────┘
                 │ HTTPS / JSON
          ┌──────▼───────┐
          │    Go API    │   auth · profiles · search · messaging
          └───┬──────┬───┘
              │      │
   ┌──────────▼─┐  ┌─▼───────────┐
   │ PostgreSQL │  │    Redis    │
   │  source of │  │  sessions,  │
   │    truth   │  │ cache, rate │
   └────────────┘  └─────────────┘
```

PostgreSQL is the single source of truth. Redis sits in front of it for anything read often and changed rarely — session tokens, popular search results, and request rate limits — so the database never becomes the bottleneck.

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or newer
- [PostgreSQL](https://www.postgresql.org/download/) 15 or newer
- [Redis](https://redis.io/download/) 7 or newer
- [Docker](https://www.docker.com/) (optional, but the quickest path)

### Installation

```bash
# Clone the repository
git clone https://github.com/rabizd/alumni.git
cd alumni

# Copy the example environment file and fill in your values
cp .env.example .env

# Download Go dependencies
go mod download

# Run the API
go run ./cmd/api
```

### With Docker Compose

```bash
docker compose up --build
```

This starts the API, PostgreSQL, and Redis together. The API is then available at `http://localhost:8080`.

## Configuration

Configuration is read from environment variables (see `.env.example`):

| Variable | Description | Example |
| --- | --- | --- |
| `APP_PORT` | Port the API listens on | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/alumni?sslmode=disable` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379/0` |
| `JWT_SECRET` | Secret used to sign access tokens | *(a long random string)* |

> 🔒 Never commit your `.env` file. Only `.env.example` belongs in version control.

## Usage

1. **Create an account** with your university email address.
2. **Complete your profile** — faculty, department, graduation year, and current employer.
3. **Search the network** by department, company, city, or graduation year.
4. **Connect** — send a message, offer mentorship, or join an upcoming reunion.

## Project Structure

> This section will grow as the code does.

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
├── docs/               # Documentation
├── docker-compose.yml
├── README.md
└── LICENSE
```

## Roadmap

- [ ] Database schema and migrations
- [ ] Authentication and session management
- [ ] Alumni profile CRUD
- [ ] Search and filtering
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
