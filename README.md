# Recipes WebApp

A full-stack web application for discovering and managing recipes, built with a microservices architecture.

## 🎯 Overview

Recipes WebApp is a distributed application that allows users to browse, search, and manage recipes. The application features a React-based frontend, a RESTful API backend, and several supporting microservices including RSS feed parsing capabilities.

## 🏗️ Architecture

The application consists of multiple microservices:

- **API Service** (`/api`): RESTful API built with Go (Gin framework)
- **Frontend** (`/recipes-web`): React-based single-page application
- **Dashboard** (`/dashboard`): Server-side rendered web interface for recipe viewing
- **Consumer** (`/consumer`): RabbitMQ consumer for processing RSS feed entries
- **Producer** (`/producer`): RabbitMQ producer for queuing RSS feed URLs
- **RSS Parser** (`/rss-parser`): Service for parsing recipe RSS feeds
- **Go Assets** (`/go-assets`): Static file server using Go

## 🛠️ Technology Stack

### Backend
- **Go 1.23+**: Main programming language
- **Gin**: Web framework
- **MongoDB**: Primary database
- **Redis**: Caching layer
- **RabbitMQ**: Message queue for async processing
- **JWT**: Authentication
- **Swagger**: API documentation

### Frontend
- **React 19**: UI library
- **Auth0**: Authentication provider
- **Bootstrap 5**: CSS framework
- **Create React App**: Build tooling

### Infrastructure
- **Docker**: Containerization
- **Docker Compose**: Service orchestration
- **Nginx**: Reverse proxy and load balancing

## 📋 Prerequisites

- **Go**: Version 1.23 or higher
- **Node.js**: Version 14 or higher
- **npm**: Version 6 or higher
- **Docker**: Version 20.10 or higher (for containerized deployment)
- **Docker Compose**: Version 1.29 or higher
- **MongoDB**: Version 4.4 or higher (if running without Docker)
- **Redis**: Version 6.2 or higher (if running without Docker)
- **RabbitMQ**: Version 3.9 or higher (for consumer/producer services)

## 🚀 Getting Started

### Option 1: Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/Mitchxxx/Recipes-WebApp.git
   cd Recipes-WebApp
   ```

2. **Run the API service with Docker Compose**
   ```bash
   cd api
   docker-compose up -d
   ```
   This will start:
   - MongoDB on port 27017
   - Redis on port 6379
   - API service (scaled to 3 instances)
   - Nginx load balancer on port 80

3. **Run the frontend**
   ```bash
   cd recipes-web
   docker-compose up -d
   ```

### Option 2: Local Development

#### Backend API

1. **Set up environment variables**
   
   Create a `.env` file in the `api` directory or set the following environment variables:
   ```bash
   MONGO_URI=mongodb://admin:password@localhost:27017/admin?authSource=admin
   MONGO_DATABASE=demo
   REDIS_ADDRESS=localhost:6379
   JWT_SECRET=your-secret-key
   X_API_KEY=your-api-key
   ```

2. **Install dependencies and run**
   ```bash
   cd api
   go mod download
   go run main.go
   ```
   The API will be available at `http://localhost:8080`

#### Frontend

1. **Install dependencies**
   ```bash
   cd recipes-web
   npm install
   ```

2. **Configure Auth0**
   
   Update the Auth0 configuration in `recipes-web/src/index.js`:
   ```javascript
   domain="your-auth0-domain"
   clientId="your-auth0-client-id"
   ```

3. **Run the development server**
   ```bash
   npm start
   ```
   The app will open at `http://localhost:3000`

#### Dashboard

```bash
cd dashboard
go mod download
go run main.go
```

#### Consumer/Producer (Optional)

For RSS feed processing:

1. **Set up RabbitMQ**
   ```bash
   docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:3-management
   ```

2. **Set environment variables**
   ```bash
   export RABBITMQ_URI=amqp://guest:guest@localhost:5672/
   export RABBITMQ_QUEUE=recipes
   export MONGO_URI=mongodb://admin:password@localhost:27017/admin?authSource=admin
   export MONGO_DATABASE=demo
   ```

3. **Run producer**
   ```bash
   cd producer
   go run main.go
   ```

4. **Run consumer**
   ```bash
   cd consumer
   go run main.go
   ```

## 📚 API Documentation

The API includes Swagger documentation. Key endpoints:

- `GET /recipes` - List all recipes
- `POST /recipes` - Create a new recipe
- `GET /recipes/:id` - Get a specific recipe
- `PUT /recipes/:id` - Update a recipe
- `DELETE /recipes/:id` - Delete a recipe
- `POST /signin` - User authentication
- `POST /refresh` - Refresh JWT token

Full API documentation can be found in `/api/swagger.json` or by accessing the Swagger UI (if configured).

**Contact**: Mitchel Egboko <megboko@gmail.com>

## 🔐 Environment Variables

### API Service
| Variable | Description | Default |
|----------|-------------|---------|
| `MONGO_URI` | MongoDB connection string | Required |
| `MONGO_DATABASE` | MongoDB database name | `demo` |
| `REDIS_ADDRESS` | Redis server address | `localhost:6379` |
| `JWT_SECRET` | Secret key for JWT tokens | Required |
| `X_API_KEY` | API authentication key | Required |

### Consumer/Producer
| Variable | Description | Default |
|----------|-------------|---------|
| `RABBITMQ_URI` | RabbitMQ connection string | Required |
| `RABBITMQ_QUEUE` | RabbitMQ queue name | Required |
| `MONGO_URI` | MongoDB connection string | Required |
| `MONGO_DATABASE` | MongoDB database name | Required |

## 🏃 Running Tests

### Backend
```bash
cd api
go test ./...
```

### Frontend
```bash
cd recipes-web
npm test
```

## 🐳 Docker Build

### Build API Image
```bash
cd api
docker build -t recipes-api .
```

### Build Frontend Image
```bash
cd recipes-web
docker build -t recipes-web .
```

### Build Dashboard Image
```bash
cd dashboard
docker build -t recipes-dashboard .
```

## 📁 Project Structure

```
Recipes-WebApp/
├── api/                    # Go REST API service
│   ├── handlers/          # Request handlers
│   ├── models/            # Data models
│   ├── main.go           # API entry point
│   ├── Dockerfile        # API container definition
│   └── docker-compose.yml # API service orchestration
├── recipes-web/           # React frontend
│   ├── src/              # React source files
│   ├── public/           # Static assets
│   ├── Dockerfile        # Frontend container definition
│   └── package.json      # Frontend dependencies
├── dashboard/            # Server-side rendered dashboard
│   ├── template/         # HTML templates
│   ├── assets/           # Static assets
│   └── main.go          # Dashboard entry point
├── consumer/             # RabbitMQ consumer
│   └── main.go          # Consumer entry point
├── producer/             # RabbitMQ producer
│   └── main.go          # Producer entry point
├── rss-parser/           # RSS feed parser
│   └── main.go          # Parser entry point
└── go-assets/            # Static file server
    └── main.go          # Asset server entry point
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is part of the Building Distributed Applications in Gin learning materials.

## 🔗 Related Resources

- [Gin Web Framework](https://gin-gonic.com/)
- [React Documentation](https://reactjs.org/)
- [MongoDB Documentation](https://docs.mongodb.com/)
- [Auth0 Documentation](https://auth0.com/docs)
- [Building Distributed Applications in Gin](https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin)

## ⚠️ Notes

- The API uses CORS (Cross-Origin Resource Sharing) configured with default settings
- Authentication is currently disabled in the main routes (commented out in main.go)
- Redis is used for caching recipe listings
- The application supports horizontal scaling via Docker Compose (API service scaled to 3 instances by default)
