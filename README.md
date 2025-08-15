# PostgreSQL Metrics Writer

A simple Go application that writes to PostgreSQL and exposes Prometheus metrics.

## Features

- Connects to PostgreSQL database
- Writes a row every second with ~10% simulated error rate
- Exposes Prometheus metrics on port 8080
- Automatic database table creation
- Environment variable configuration

## Quick Start

```bash
# Build the Docker image
docker build -t app .

# Load image to kind
kind load docker-image .

# Deploy to Kubernetes
kubectl apply -f k8s/postgresql.yaml
kubectl apply -f k8s/app.yaml
```

## Metrics

- `postgres_rows_written_total` - Counter of successful writes
- `postgres_write_errors_total` - Counter of failed writes

## Environment Variables

- `DB_HOST` (default: localhost)
- `DB_PORT` (default: 5432)
- `DB_USER` (default: devuser)
- `DB_PASSWORD` (default: devpass)
- `DB_NAME` (default: devdb)
