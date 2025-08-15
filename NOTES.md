# My Setup
I’m running this in a kind cluster on WSL, with Docker Desktop on Windows 11 and kind running in Debian within WSL.

# Running the project
- `docker build -t app .`
- If running under `kind`: `kind load docker-image app`
- `kubectl apply -k k8s/`

# Adjustments
- When I ran kubectl apply -f postgres.yaml, it failed due to a cgroup version mismatch, likely caused by a WSL2 kernel issue. I resolved the problem by commenting out the resource specifications in postgres.yaml.
- I'm deploying my custom app to the same namespace as postgres, I typically wouldn't do that

# Notes
- Default environment variables are being set in the go code. This is not ideal by any means, under normal circumstances I would store this in a secret manager like `vault` or `kubernetes secrets`.
- The `app.yaml` has the `env` set; it isn't actually necessary, this would be pulled from the secret manager
- I added a `kustomization.yaml` just because I would do this under normal circumstances

## Prometheus Metrics
- Metric names are based on the Prometheus documentation: https://prometheus.io/docs/practices/naming/
- `postgres_rows_written_total` defines how many rows have been written to the database
- `postgres_write_errors_total` defines how many errors have occurred while writing to the database

## Database Connection
- Implemented retry logic, there are 30 attempts every 2 seconds while waiting for containers to startup
- Database will automatic bootstrap; if the table doesn't exist, it will be created, otherwise it will continue

## Error Simulation
- There is roughly a 10% error rate; this is simplified with `rand.Intn(100) < 10` logic
- Used simulated errors such as `database connection failed`, `database timeout`, `database unavailable`, `database error`

## Container
- Set `imagePullPolicy` to `Never` because this is required for local Docker images
- App should wait for the database readiness rather than using init containers
- Metrics server should not block the main loop, it runs as a `goroutine`.
- Based on my research, it is best practice to include the `go.mod` and `go.sum` files in the `Dockerfile` to prevent version drift
