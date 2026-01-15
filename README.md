# GopherFeed

A production-ready Go backend implementing a social feed API with authentication, caching, async messaging, and comprehensive observability.

![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)

## Features

- **RESTful API** with versioned routes (`/v1`) and Swagger documentation
- **Posts & Comments** - Full CRUD operations with ownership validation
- **User Feed** - Filterable, sortable, paginated user feed
- **Authentication** - JWT tokens with HTTP-only cookies, Basic Auth support
- **Authorization** - Role-based access control (user/moderator/admin)
- **Email Verification** - Async email delivery via RabbitMQ + Mailtrap
- **Redis Caching** - Optional caching layer for user data
- **Rate Limiting** - Fixed-window algorithm per IP
- **Structured Logging** - Production-ready JSON logs via Zap
- **Graceful Shutdown** - Signal handling for clean termination
- **Server Metrics** - Customizable runtime stats
- **CI/CD** - GitHub Actions for linting and testing

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Router | chi v5 |
| Database | PostgreSQL |
| Cache | Redis 6.2 |
| Message Queue | RabbitMQ 3 |
| Email | Mailtrap SMTP |
| Auth | JWT (golang-jwt/jwt/v5) |
| Docs | Swagger (swaggo) |
| Logging | uber-go/zap |
| Migrations | goose |

## Project Structure

```
├── cmd/
│   ├── api/           # API server (handlers, middleware, routes)
│   ├── migrate/       # Database migrations and seed scripts
│   └── worker/        # RabbitMQ email consumer
├── internal/
│   ├── auth/          # JWT authenticator
│   ├── db/            # Database connection
│   ├── email/         # Mailtrap client + templates
│   ├── env/           # Environment config loader
│   ├── mq/            # RabbitMQ publisher/consumer
│   ├── ratelimiter/   # Fixed-window rate limiter
│   ├── store/         # Data access layer
│   │   └── cache/     # Redis caching layer
│   └── utils/         # Shared utilities
├── web/               # Email templates
├── docs/              # Generated Swagger files
├── scripts/           # DB init scripts
└── compose.yml        # Docker services
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [goose](https://github.com/pressly/goose) (migrations)
- [Air](https://github.com/air-verse/air) (hot reload, optional)

### Quick Start

1. **Start infrastructure:**
   ```bash
   docker compose up -d
   ```

2. **Configure environment** (create `.envrc` or export directly):
   ```bash
   export DB_URL="postgres://user:password@localhost:5432/gopherfeed?sslmode=disable"
   export RABBITMQ_URL="amqp://guest:guest@localhost:5672/"
   export REDIS_ADDR="localhost:6379"
   export JWT_SECRET="your-secret-key"
   export MAIL_FROM_EMAIL="noreply@gopherfeed.com"
   export MAILTRAP_API_KEY="your-mailtrap-key"
   ```

3. **Run migrations:**
   ```bash
   make migrate-up
   ```

4. **Seed database (optional):**
   ```bash
   make seed
   ```

5. **Start the API:**
   ```bash
   go run cmd/api/main.go
   ```

6. **Start the worker (separate terminal):**
   ```bash
   go run cmd/worker/main.go
   ```

### Hot Reload Development

```bash
# API server
air -c .air.toml

# Worker (separate terminal)
air -c .air-worker.toml
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | `8080` |
| `DB_URL` | PostgreSQL connection string | Required |
| `DB_MAX_OPEN_CONNS` | Max open DB connections | `30` |
| `DB_MAX_IDLE_CONNS` | Max idle DB connections | `30` |
| `DB_MAX_IDLE_TIME` | Max idle connection time | `15m` |
| `REDIS_ADDR` | Redis address | `""` (disabled) |
| `REDIS_ENABLED` | Enable Redis caching | `false` |
| `RABBITMQ_URL` | RabbitMQ connection string | Required |
| `JWT_SECRET` | JWT signing secret | Required |
| `JWT_EXPIRY` | Token expiration duration | `72h` |
| `JWT_ISSUER` | Token issuer claim | `gopherfeed` |
| `JWT_AUDIENCE` | Token audience claim | `gopherfeed` |
| `MAILTRAP_API_KEY` | Mailtrap API key | Required |
| `MAIL_FROM_EMAIL` | Sender email address | Required |
| `RATE_LIMIT_ENABLED` | Enable rate limiting | `true` |
| `RATE_LIMIT_RPS` | Requests per second | `20` |
| `RATE_LIMIT_BURST` | Burst limit | `40` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `""` |

## API Endpoints

Base URL: `http://localhost:8080/v1`

### Public
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/auth/user` | Register user |
| PUT | `/auth/user/activate/{token}` | Activate account |
| POST | `/auth/token` | Login (get JWT) |

### Protected (JWT Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/posts` | List posts (paginated) |
| POST | `/posts` | Create post |
| GET | `/posts/{id}` | Get post by ID |
| PATCH | `/posts/{id}` | Update post (owner only) |
| DELETE | `/posts/{id}` | Delete post (owner/admin) |
| POST | `/posts/{id}/comments` | Add comment |
| GET | `/users/{id}` | Get user profile |
| PUT | `/users/{id}/follow` | Follow user |
| PUT | `/users/{id}/unfollow` | Unfollow user |
| GET | `/users/feed` | Get personalized feed |

### Swagger Documentation

Visit `http://localhost:8080/v1/swagger/index.html` after starting the API.

Regenerate docs:
```bash
make gen-docs
```

## Key Features

### Authentication & Authorization

- JWT tokens stored in HTTP-only cookies for CSRF protection
- Basic Auth alternative for simple integrations
- Role-based permissions: `user`, `moderator`, `admin`
- Email verification required for account activation

### Async Email

Uses RabbitMQ for reliable, non-blocking email delivery:
1. API publishes email task to queue
2. Worker consumes and sends via Mailtrap
3. Embedded HTML templates in `web/` directory

### Caching

Optional Redis caching for user lookups. Enable with:
```bash
export REDIS_ENABLED=true
export REDIS_ADDR=localhost:6379
```

### Rate Limiting

Fixed-window algorithm limits requests per IP. Configure via:
```bash
export RATE_LIMIT_RPS=20
export RATE_LIMIT_BURST=40
```

## Development

### Run Tests

```bash
make test
```

### Database Migrations

```bash
# Create new migration
make migration add_new_table

# Apply migrations
make migrate-up

# Rollback last migration
make migrate-down
```

### Docker Services

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

Management UIs:
- RabbitMQ: http://localhost:15672 (guest/guest)
- Redis Commander: http://localhost:8081

## License

MIT License - see [LICENSE](LICENSE)
