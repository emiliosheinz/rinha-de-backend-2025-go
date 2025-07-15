# Rinha de Backend 2025

This repository contains my Go-based implementation for the [Rinha de Backend 2025 challenge](https://github.com/zanfranceschi/rinha-de-backend-2025). The solution focuses on maximizing financial profit by intelligently routing payment requests between two Payment Processor services (default and fallback), handling service instabilities, and ensuring data consistency via a payment summary endpoint.

## Architecture

The system is composed of the following services:

- **API Service (Go)**: Two instances (`api01` and `api02`) receive and enqueue payment requests asynchronously.
- **Asynq Worker**: A Redis-backed background worker pool processes enqueued payment tasks.
- **Health Manager**: A leader-elected component that periodically polls both processors' health endpoints and caches their status/latency in Redis.
- **PostgreSQL**: Persists payment records (`payments` table) for auditing and summary reporting.
- **Redis**: Used as a queue backend (Asynq) and to store cached health metrics and perform leader election.
- **NGINX Load Balancer**: Distributes incoming HTTP traffic across API instances on port `9999`.
- **Docker Compose**: Orchestrates all containers, including the external `payment-processor` network with default and fallback processors.

## Key Design Decisions

### Asynchronous Processing
Incoming `POST /payments` requests return immediately (HTTP 202) after enqueueing a task. Background workers handle the actual payment processing to improve throughput and resilience.

### Health-Driven Routing
A single leader polls `/payments/service-health` every 5 s (within the rate limit of one call per 5 s) and stores the metrics (`failing` flag and `minResponseTime`) in Redis. Routing logic:

- Always avoid routing to a processor marked as failing.
- Switch to the fallback if the default's latency exceeds 1.25× the fallback's minimum response time.

### Leader Election
Instances use a Redis-based leader election (SETNX with TTL renewal every 10 s) to ensure only one node performs health checks and updates the shared cache.

### Data Persistence & Summary
Successful payments are recorded in PostgreSQL with correlation ID, amount, timestamp, and chosen processor. The `GET /payments-summary` endpoint aggregates totals per processor, optionally filtered by time window (`from`/`to` in RFC3339Nano).

### Containerization & Resource Limits
- Multi-stage Dockerfile produces a minimal scratch binary for production.
- Docker Compose defines CPU (up to 1.5 cores total) and memory (350 MB) limits per service as per challenge requirements.
- Uses bridge networking and declares the external `payment-processor` network for the default/fallback services.

## API Endpoints

| Method | Path                | Description                                                        |
|--------|---------------------|--------------------------------------------------------------------|
| POST   | `/payments`         | Enqueue a new payment. Returns HTTP 202 if input is valid.         |
| GET    | `/payments-summary` | Retrieve aggregated payment summary. Supports optional `from`/`to` filters. |

## Getting Started

### 1. Start the Payment Processors
```bash
docker-compose -f payment-processors/docker-compose.yml up -d
```

### 2. Set Environment Variables
Create a `.env` file (Compose will auto-load it) through the following command:
```bash
cp .env.example .env
```

### 3. (Optional) Enable Live Reload & Monitoring
Copy the development override and use Air for hot-reloading and AsynqMon for queue monitoring:
```bash
cp docker-compose.override.example.yml docker-compose.override.yml
docker-compose up --build
```

### 4. Launch the Full Stack
```bash
docker-compose up -d --build
```

### 5. Verify Services
- API & workers: logs via `docker-compose logs -f api01 api02`
- NGINX LB: listen on `http://localhost:9999`
- AsynqMon (if enabled): `http://localhost:8080`

## Testing

Execute the official test suite to validate the implementation by running:
```bash
./test.sh
```

## Conclusion
This implementation tried to balances performance and resilience by combining asynchronous task processing, dynamic health-aware routing, and leader election to meet the challenge's scoring criteria for profit maximization and latency targets. The main idea behind it all was to implement something that could be used in a real-world scenario and not just a hacky solution to pass the tests.
