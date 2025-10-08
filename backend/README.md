# Beer Festival Backend

Simple Go backend that serves festival data via a REST API.

## Environment Variables

- `PORT` - Server port (default: `8080`)
- `FESTIVALS_FILE` - Path to festivals JSON file (default: `./festivals.json`)
- `ALLOWED_ORIGINS` - CORS allowed origins (default: `*`)

## Running the Server

```bash
cd backend
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### GET /api/festivals

Returns a JSON array of all festivals.

**Response:**
```json
[
  {
    "id": "1",
    "name": "Paris Beer Week",
    "description": "La plus grande célébration de la bière...",
    "startDate": "2025-11-15",
    "endDate": "2025-11-22",
    "city": "Paris",
    "region": "Île-de-France",
    "location": {
      "latitude": 48.8566,
      "longitude": 2.3522
    },
    "image": "https://...",
    "website": "https://...",
    "breweryCount": 120
  }
]
```

## Running Tests

```bash
go test -v
```

## Configuration

Copy `.env.example` to `.env` and adjust values as needed:

```bash
cp .env.example .env
```
