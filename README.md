Tracky — Telegram Mini App for Package Tracking

Tracky is a lightweight Telegram Mini App that allows users to track shipments directly within Telegram. It uses AfterShip for tracking data and is built with Go and a web frontend.

🔧 Stack

Backend: Go (Gin), PostgreSQL

Frontend: WebApp for Telegram (HTML/CSS/JS)

Integrations: AfterShip API

Migrations: golang-migrate

📁 Project Structure

tracky/
├── cmd/server/         # Entry point for the Go server
├── internal/           # Business logic, handlers, integrations
├── migrations/         # SQL migration files
└── ui/webapp/          # Frontend for Telegram Mini App

🚀 Getting Started

Clone the repo:

git clone github.com/tansystems/tracky.git

cd tracky

Copy and configure environment variables:

cp .env.example .env

Start the backend:

go run cmd/server/main.go

🛠️ Database Migrations

Install migrate CLI:

brew install golang-migrate

Run migrations:

migrate -database "$DB_DSN" -path ./migrations up

Make sure $DB_DSN is set, e.g.:

export DB_DSN="postgres://user:password@localhost:5432/tracky?sslmode=disable"

✨ Features

Track packages from multiple couriers via AfterShip
Clean Telegram WebApp interface
Lightweight Go backend
REST API for frontend integration








