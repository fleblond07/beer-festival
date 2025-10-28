# 🍺 Festivals de Bière - Beer Festival Tracker

Beers, for everyone. That's the idea behind the website, register and share every beer festival or any type of event beer-related in France.

# Technical details:
Feel free to run a copy of this website, below instruction on features and how to run it!

## ✨ Features

- **Next Festival Highlight**: Prominently displays the upcoming festival
- **Interactive Map**: Explore all festivals on an interactive Leaflet map
- **Festival Cards**: Browse all festivals with detailed information
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile
- **Clean Architecture**: Component-based structure with TypeScript
- **Full Test Coverage**: ~95% test coverage with unit and E2E tests
- **Go Backend API**: RESTful API serving festival data with CORS support
- **Mock Data**: Pre-populated with example French beer festivals

## 🚀 Quick Start

### Prerequisites

- Node.js 22.x or higher
- npm 10.x or higher
- Go 1.x or higher (for backend)

### Installation

```bash
make install
```
```bash
# Start both backend and frontend concurrently
make dev

# Or start them separately:
make backend  # Starts backend on http://localhost:8080
make frontend # Starts frontend on http://localhost:5173
```

### Build for Production

```bash
make build
```

### Docker


#### Using Docker Compose

```bash
docker-compose up --build

docker-compose down
```

This starts:
- **Backend** on http://localhost:8080
- **Frontend** on http://localhost:3000


#### Building Individual Containers for Deployment

```bash
docker build -t beer-festival-backend ./backend
docker build -t beer-festival-frontend ./frontend

docker run -d -p 8080:8080 --name backend beer-festival-backend
docker run -d -p 3000:3000 --name frontend beer-festival-frontend
```

## 🧪 Testing

```bash
make test
```

## 📁 Project Structure

This is a monorepo with separate frontend and backend directories:

```
beer-festival/
├── frontend/          # Vue.js frontend application
│   ├── src/          # Frontend source code
│   ├── tests/        # Frontend tests
│   └── Dockerfile    # Frontend container configuration
├── backend/          # Go backend API
│   ├── main.go       # Backend entry point
│   ├── handlers.go   # API handlers
│   └── Dockerfile    # Backend container configuration
└── Makefile          # Build and development commands
```

## 🔌 Backend API

The Go backend provides a RESTful API for festival data:

### Endpoints

- `GET /api/v1/festivals` - Returns all festivals
- `GET /health` - Health check endpoint

### Configuration

The backend supports the following environment variables:

- `PORT` - Server port (default: `8080`)
- `FESTIVALS_FILE` - Path to festivals JSON file (default: `./festivals.json`)

## 🚧 Future Enhancements

- [ ] Festival detail pages
- [ ] Filter and search functionality
- [ ] Event calendar integration
- [ ] Social sharing features
- [ ] Multi-language support
- [ ] Database integration for persistent storage

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is a demonstration application. Feel free to use it as a template for your own projects.

## 🙏 Acknowledgments

- Mock festival data is fictional
- Map tiles from OpenStreetMap
- Icons from Heroicons
- Built with ❤️ for craft beer enthusiasts
