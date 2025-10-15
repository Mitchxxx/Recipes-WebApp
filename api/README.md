# Recipes API Service

RESTful API service for the Recipes WebApp, built with Go and the Gin framework.

## Overview

This is the core backend service that provides recipe management functionality including CRUD operations, user authentication, and caching capabilities. The API uses MongoDB for data persistence and Redis for caching frequently accessed recipes.

## Technology Stack

- **Go 1.23+**: Programming language
- **Gin**: Web framework with middleware support
- **MongoDB**: Document database for storing recipes and users
- **Redis**: In-memory cache for performance optimization
- **JWT**: JSON Web Tokens for authentication
- **Swagger**: API documentation

## API Endpoints

### Authentication
- `POST /signin` - Authenticate user and receive JWT token
- `POST /refresh` - Refresh an existing JWT token

### Recipes
- `GET /recipes` - List all recipes (cached in Redis)
- `POST /recipes` - Create a new recipe
- `GET /recipes/:id` - Get a specific recipe by ID
- `PUT /recipes/:id` - Update an existing recipe
- `DELETE /recipes/:id` - Delete a recipe

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `MONGO_URI` | MongoDB connection string | Yes | - |
| `MONGO_DATABASE` | Database name | Yes | `demo` |
| `REDIS_ADDRESS` | Redis server address | Yes | `localhost:6379` |
| `JWT_SECRET` | Secret key for signing JWT tokens | Yes | - |
| `X_API_KEY` | API key for authentication middleware | No | - |

## Local Development

### Prerequisites
- Go 1.23 or higher
- MongoDB 4.4 or higher
- Redis 6.2 or higher

### Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set environment variables**
   ```bash
   export MONGO_URI="mongodb://admin:password@localhost:27017/admin?authSource=admin"
   export MONGO_DATABASE="demo"
   export REDIS_ADDRESS="localhost:6379"
   export JWT_SECRET="your-secret-key"
   ```

3. **Run the service**
   ```bash
   go run main.go
   ```

The API will start on port `8080`.

## Docker Deployment

### Build the Docker image
```bash
docker build -t recipes-api .
```

### Run with Docker Compose
```bash
docker-compose up -d
```

This starts:
- MongoDB (port 27017)
- Redis (port 6379)
- API service (3 instances for load balancing)
- Nginx reverse proxy (port 80)

## Project Structure

```
api/
├── handlers/          # HTTP request handlers
│   ├── auth.go       # Authentication handlers
│   └── recipes.go    # Recipe CRUD handlers
├── models/           # Data models
│   ├── recipe.go     # Recipe model
│   └── user.go       # User model
├── main.go           # Application entry point
├── Dockerfile        # Container image definition
├── docker-compose.yml # Service orchestration
├── nginx.conf        # Nginx configuration
├── swagger.json      # API documentation
└── go.mod            # Go module dependencies
```

## Features

### Caching
Recipe listings are cached in Redis to improve performance. Cache is invalidated when:
- A new recipe is created
- An existing recipe is updated
- A recipe is deleted

### CORS Support
The API includes CORS middleware configured with default settings to allow cross-origin requests from the frontend.

### Error Handling
All endpoints return appropriate HTTP status codes:
- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Load Balancing
The Docker Compose configuration includes Nginx as a reverse proxy that load balances requests across 3 API instances.

## API Documentation

Full Swagger documentation is available in `swagger.json`. The API follows OpenAPI 2.0 specification.

**API Contact**: Mitchel Egboko <megboko@gmail.com>

## Testing

Run tests with:
```bash
go test ./...
```

## Notes

- The authentication middleware is currently commented out in the main routes
- The API uses graceful shutdown to handle termination signals
- Server timeouts are configured: 45 seconds read/write timeout
- MongoDB connection includes health checks via ping operations
