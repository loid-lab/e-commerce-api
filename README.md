# E-Commerce API

This is a **Go-based e-commerce backend API** built with the Gin framework and GORM for ORM. It is designed to be backend-only, RESTful, and optionally Stripe-integrated. You can use it for web or mobile storefronts as a secure and modular API service.

---

## 🔧 Tech Stack

- **Go** (Gin + GORM)
- **PostgreSQL** (Neon or local)
- **JWT Auth**
- **Stripe Payments**
- **Zod** (for shared schema validation)

---

## 📦 Features

- ✅ User authentication (signup, login, JWT)
- 🛒 Cart and cart items per user
- 📦 Product, category, and order management
- 💳 Stripe Checkout integration
- 🔐 Auth middleware (with claims)
- ✅ Zod schema validation (extra layer on frontend/backend if needed)

---

## 🗂 Folder Structure

```
e-commerce-api/
├── controllers/
├── middleware/
├── models/
├── initializers/
├── validators/
├── main.go
├── go.mod
├── .env.example
└── README.md
```

---

## 📌 Setup Instructions

### 1. Clone the repo

```bash
git clone https://github.com/loid-lab/e-commerce-api.git
cd e-commerce-api
```

### 2. Setup environment

```bash
cp .env.example .env
```

Update `.env` with your configuration values:

- For **Neon (cloud Postgres)**:

  ```
  DB_URL=postgres://user:password@ep-xxx.neon.tech:5432/dbname
  SECRET=your_jwt_secret
  STRIPE_SECRET_KEY=sk_test_...
  REDIS_URL=redis://redis:6379
  ```

- For **local Postgres** (optional, see Docker Compose below):

  ```
  DB_URL=postgres://postgres:example@localhost:6379/ecommercedb?sslmode=disable
  SECRET=your_jwt_secret
  STRIPE_SECRET_KEY=sk_test_...
  REDIS_URL=redis://localhost:6379
  ```

---

### 3. Running with Docker Compose

#### Using Neon (cloud DB)

Your `docker-compose.yaml` should look like:

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_URL: ${DB_URL}
      SECRET: ${SECRET}
      STRIPE_SECRET_KEY: ${STRIPE_SECRET_KEY}
      REDIS_URL: ${REDIS_URL}
    restart: unless-stopped
```

Simply run:

```bash
docker-compose up --build
```

#### Using local Postgres and Redis

Uncomment or add the following `db` and `redis` services in your `docker-compose.yaml`:

```yaml
version: '3.8'

services:
  api:
    build: .
    depends_on:
      - db
      - redis
    ports:
      - "8080:8080"
    environment:
      DB_URL: postgres://postgres:example@db:5432/ecommercedb?sslmode=disable
      SECRET: ${SECRET}
      STRIPE_SECRET_KEY: ${STRIPE_SECRET_KEY}
      REDIS_URL: redis://redis:6379
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: example
      POSTGRES_DB: ecommercedb
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    restart: always
    ports:
      - "6379:6379"

volumes:
  pgdata:
```

Run with:

```bash
docker-compose up --build
```

---

## 🔐 Authentication

- `POST /auth/signup` — Register new user  
- `POST /auth/login` — Login and receive JWT  
- Authenticated routes require `Authorization: Bearer <token>` header

---

## 💳 Payments (Stripe)

- `POST /user/orders/:id/pay` — Initiate Stripe checkout session  
- Payments stored in DB with status `pending`/`paid`

---

## 🛠 Zod Validation

Zod schemas (stored in `validators/`) are optional but recommended for:

- Frontend-to-backend shared validation  
- Backend request payload validation  

---

## 📘 License

MIT — feel free to use, modify, or contribute.
