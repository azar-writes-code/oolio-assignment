# Oolio Product Service

A robust, high-performance microservice for managing products and processing orders with coupon validation. Built with Go, PostgreSQL, and BadgerDB.

## 🚀 Features

- **Product Management**: CRUD operations for products with support for bulk creation.
- **Order Processing**: Transactional order placement with atomic stock management.
- **Coupon Validation**: High-speed coupon lookup using BadgerDB (Key-Value store).
- **Observability**: Integrated Prometheus metrics and structured logging (slog).
- **Architecture**: Domain-driven design with clean separation of layers (Controllers, Services, Repositories).
- **Developer Experience**: Hot-reloading with Air and automated SQL code generation with SQLC.

## 🛠️ Technology Stack

- **Language**: Go 1.26+
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/) (via [pgx](https://github.com/jackc/pgx))
- **Key-Value Store**: [BadgerDB](https://github.com/dgraph-io/badger) (for Coupons)
- **SQL Tooling**: [SQLC](https://sqlc.dev/) (Type-safe SQL generator)
- **Monitoring**: Prometheus & Loki

---

## 🏃 Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose
- [Air](https://github.com/cosmtrek/air) (for live reload)
- [SQLC](https://sqlc.dev/docs/installing/) (for regenerating DB code)

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd product-service
   ```

2. **Setup environment variables**:
   ```bash
   cp .env.example .env
   # Update .env with your local database credentials if needed
   ```

3. **Spin up infrastructure**:
   ```bash
   make up
   ```

4. **Run migrations**:
   ```bash
   make run-migrations
   ```

5. **Start the application**:
   ```bash
   # Using hot reload (recommended for development)
   air

   # OR using Makefile (starts infra + migrations + air)
   make start
   ```

---

## 📂 Project Structure

```text
├── cmd/                # Entry points (App, DB Migrator, Coupon Verifier)
├── config/             # Configuration management (Viper/Env)
├── infra/              # Infrastructure (Docker Compose, Prometheus config)
├── pkg/
│   ├── app/            # Application lifecycle & bootstrap
│   ├── server/
│   │   └── rest/       # REST API Implementation
│   │       ├── controllers/ # HTTP Handlers
│   │       ├── services/    # Business Logic
│   │       ├── models/      # Domain Entities
│   │       ├── dtos/        # Data Access Objects (SQLC, Badger)
│   │       └── routes/      # Gin Route definitions
│   └── telemetry/      # Metrics & Logging setup
└── sqlc.yaml           # SQLC Configuration
```

---

## 📡 API Endpoints

### Products
- `POST /api/v1/products/create-product`: Create a single product.
- `POST /api/v1/products/create-many-products`: Bulk create products.
- `GET /api/v1/products/`: List products (paginated).
- `GET /api/v1/products/:id`: Get product details.

### Orders
- `POST /api/v1/order/`: Place a new order (Requires Auth).
- `POST /api/v1/order/validate-coupon`: Manually validate a promo code.

### System
- `GET /api/v1/health/health`: System health check.
- `GET /metrics`: Prometheus metrics.

---

## 🔧 Development Tasks (Makefile)

| Command | Description |
| :--- | :--- |
| `make up` | Start Docker infrastructure (Postgres, Prometheus, Loki). |
| `make down` | Stop infrastructure and remove volumes. |
| `make start` | Full developer setup (Infra + Migrations + Air). |
| `make clean-db` | Reset database and regenerate SQLC code. |
| `make run-migrations` | Run pending database migrations. |
| `make generate-sqlc` | Regenerate type-safe Go code from SQL queries. |
| `make logs` | Stream Docker container logs. |

---

## 🧪 Testing

Run the test suite:
```bash
go test ./... -v
```

---

## 📝 License

Distributed under a Proprietary License. Commercial use of this project, in whole or in part, is strictly prohibited without prior written permission from the copyright holder. See `LICENSE` for more information.
