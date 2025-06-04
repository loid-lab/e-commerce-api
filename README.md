# E-Commerce API

This is a **Go-based e-commerce backend API** built with the Gin framework and GORM for ORM. It is designed to be backend-only, RESTful, and optionally Stripe-integrated. You can use it for web or mobile storefronts as a secure and modular API service.

---

## 🔧 Tech Stack

- **Go** (Gin + GORM)
- **PostgreSQL** (Neon or local)
- **JWT Auth**
- **Stripe Payments**
- **Zod** (for shared schema validation)
- **Redis** (for rate limiting + caching)

---

## 📦 Features

- ✅ User authentication (signup, login, JWT)
- 🛒 Cart and cart items per user
- 📦 Product, category, and order management
- 💳 Stripe Checkout integration
- 🔐 Auth middleware (with claims)
- ✅ Zod schema validation (extra layer on frontend/backend if needed)
- 🚦 Rate limiting middleware backed by Redis for enhanced security and scalability
- 🧠 Optional reCAPTCHA validation for signup/login to prevent bot activity
- 🚀 Redis caching integrated for product, category, order, and cart reads

---

## 🗂 Folder Structure

```
e-commerce-api/
├── controllers/     # HTTP handlers for routes
├── initializers/    # DB, Redis, and env config
├── middleware/      # JWT, rate limiting, etc.
├── models/          # GORM models
├── utils/           # Reusable utilities (cache, mail, recaptcha, etc.)
├── validators/      # Zod schemas (for optional validation)
├── main.go          # Application entry point
├── go.mod
├── .env.example     # Sample env vars
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
  MAILTRAP_HOST=smtp.mailtrap.io
  MAILTRAP_PORT=587
  MAILTRAP_USERNAME=your_username
  MAILTRAP_PASSWORD=your_password
  MAIL_FROM=your@email.com
  MAIL_TO=receiver@email.com
  ```

- For **local Postgres** (optional, see Docker Compose below):

  ```
  DB_URL=postgres://postgres:example@localhost:6379/ecommercedb?sslmode=disable
  SECRET=your_jwt_secret
  STRIPE_SECRET_KEY=sk_test_...
  REDIS_URL=redis://localhost:6379
  ... # same Mailtrap vars
  ```

> **Note:**  
> This project uses Redis for **rate limiting** and **caching**. Set `REDIS_URL` and ensure Redis is running locally or remotely.

---

### 3. Running with Docker Compose

#### Using Neon (cloud DB)

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
      MAILTRAP_HOST: ${MAILTRAP_HOST}
      MAILTRAP_PORT: ${MAILTRAP_PORT}
      MAILTRAP_USERNAME: ${MAILTRAP_USERNAME}
      MAILTRAP_PASSWORD: ${MAILTRAP_PASSWORD}
      MAIL_FROM: ${MAIL_FROM}
      MAIL_TO: ${MAIL_TO}
    restart: unless-stopped
```

```bash
docker-compose up --build
```

#### Using local Postgres + Redis

Add `db` and `redis` services in `docker-compose.yaml`:

```yaml
  db:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: example
      POSTGRES_DB: ecommercedb
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    restart: always
    ports:
      - "6379:6379"
```

---

## 🔐 Authentication

- `POST /auth/signup` — Register new user  
- `POST /auth/login` — Login and receive JWT  
- Authenticated routes require `Authorization: Bearer <token>`

---

### 🧠 reCAPTCHA Support

Signup/login support reCAPTCHA v2/v3. Send `recaptchaToken` in form payload.

---

## 💳 Payments (Stripe)

- `POST /user/orders/:id/pay` — Initiate Stripe checkout session  
- Stripe session stores payment reference

---

## 🧊 Redis Caching

Improves performance for heavy-read routes:

- 🔁 Products: `GET /products`, `GET /products/:id`
- 🔁 Orders: `GET /orders`, `GET /orders/:id`
- 🔁 Categories: `GET /categories`
- 🔁 Carts: `GET /cart`

Auto invalidation after `create/update/delete` where applicable.

---

## 🛠 Zod Validation

Used optionally for frontend/backend data agreement via `validators/`.

---

## 📘 License

MIT — feel free to use, modify, or contribute.
