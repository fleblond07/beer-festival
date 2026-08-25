# 🍺 Festivals de Bière - Beer Festival Tracker

Beers, for everyone. That's the idea behind the website, register and share every beer festival or any type of event beer-related in France.

# Technical details:
Feel free to run a copy of this website, below instruction on features and how to run it!

## Features

- Next festival highlight
- Interactive Leaflet map
- Festival and brewery browsing
- Admin login and festival creation
- Go backend API
- Vue frontend
- Dockerized Postgres persistence

## Quick Start

### Prerequisites

- Node.js 22.x or higher
- npm 10.x or higher
- Go 1.21.1 or higher for local backend development
- Docker and Docker Compose for the full stack
- `pg_dump` and `psql` when transferring data from Supabase

### Configure environment

Copy `.env.example` to `.env` and fill in the required values before starting anything:

```bash
cp .env.example .env
```

- `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` - required, used to initialize Postgres
- `DATABASE_URL` - required, Postgres connection string for the backend
- `JWT_SECRET` - required, HMAC secret used to sign bearer tokens
- `ADMIN_EMAIL`, `ADMIN_PASSWORD` - optional; set both to have Docker create/update an admin user on first init

`docker-compose` reads `.env` from the project root automatically. When running the backend outside Docker (`make backend`), export `DATABASE_URL`, `JWT_SECRET`, and optionally `PORT` in your shell instead.

### Installation

```bash
make install
```

### Start Postgres

The backend needs a reachable Postgres instance, including for local (non-Docker) backend development:

```bash
docker-compose up -d postgres
```

### Run the app

```bash
make dev

# Or start them separately:
make backend
make frontend
```

### Docker Compose (full stack)

With `.env` configured as above:

```bash
docker-compose up --build

docker-compose down
```

This starts:

- Postgres on `localhost:5432`
- Backend on `http://localhost:1337`
- Frontend on `http://localhost:1338`

The admin user is created only when the Postgres volume is initialized and `ADMIN_EMAIL`/`ADMIN_PASSWORD` are both set. For an existing volume, insert or update the `users` table manually, or run the transfer script with `ADMIN_EMAIL` and `ADMIN_PASSWORD` set.

## Transfer Supabase Data

The app now uses its own Docker Postgres database. To copy the existing Supabase public table data into it:

```bash
docker-compose up -d postgres

SUPABASE_DATABASE_URL='postgres://postgres:<password>@db.<project-ref>.supabase.co:5432/postgres?sslmode=require' \
DATABASE_URL='postgres://beer_festival:beer_festival@localhost:5432/beer_festival?sslmode=disable' \
ADMIN_EMAIL=admin@example.com \
ADMIN_PASSWORD='change-me' \
bash scripts/transfer_supabase_data.sh
```

The script truncates and reloads these public tables: `festivals`, `breweries`, and `festivals_breweries`. Supabase Auth passwords cannot be exported directly, so create local admin accounts with `ADMIN_EMAIL` and `ADMIN_PASSWORD`.

## Build the frontend

```bash
make build
```

This builds only the frontend (`frontend/dist`); the Go backend is run with `go run .` or built separately with `go build`.

## Testing

```bash
make test
```

## Project Structure

```text
beer-festival/
├── db/                # Postgres initialization SQL
├── frontend/          # Vue.js frontend application
├── backend/           # Go backend API
├── scripts/           # Data transfer utilities
└── Makefile           # Build and development commands
```

## Backend API

- `GET /api/festivals` - Returns all festivals
- `GET /api/breweries` - Returns all breweries
- `GET /api/festivals/{id}/breweries` - Returns breweries for a festival
- `POST /api/auth/login` - Logs in with local Postgres users
- `GET /api/auth/verify` - Verifies a bearer token
- `POST /api/festivals/create` - Creates a festival; requires bearer token
- `GET /health` - Health check endpoint

## Configuration

- `PORT` - Backend port
- `ALLOWED_ORIGINS` - CORS origins, comma-separated or `*`
- `DATABASE_URL` - Postgres connection string
- `JWT_SECRET` - HMAC secret used to sign bearer tokens
- `ADMIN_EMAIL` - Optional Docker init or transfer-script admin email
- `ADMIN_PASSWORD` - Optional Docker init or transfer-script admin password
