# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Setup and Run
```bash
# Initial setup
cp .env.example .env
docker-compose -f payment-processors/docker-compose.yml up -d

# Development with hot reload
cp docker-compose.override.example.yml docker-compose.override.yml
docker-compose up --build

# Production mode
docker-compose up -d --build

# Run official test suite
./test.sh
```

### Development Tools
- **Air**: Hot reload configured in `.air.toml`
- **AsynqMon**: Queue monitoring UI at http://localhost:8080 (dev mode only)

## High-Level Architecture

This is a Go-based payment processing system designed for the Rinha de Backend 2025 challenge. The architecture prioritizes throughput and resilience through asynchronous processing and intelligent routing.

### Core Components

1. **API Service** (`cmd/api/main.go`): HTTP server handling payment requests and summaries
   - Returns 202 Accepted immediately for payment requests
   - Queues payments for async processing

2. **Asynq Worker**: Redis-backed task processor
   - 32 concurrent workers processing payments
   - Implements retry logic with exponential backoff

3. **Health Manager** (`internal/health/`): Monitors payment processor availability
   - Leader-elected service using Redis SETNX
   - Polls processors every 5 seconds
   - Caches health data for all instances

4. **Payment Router** (`internal/payments/`): Intelligent routing logic
   - Avoids failing processors
   - Switches to fallback when default's latency > 1.25× fallback's minimum

### Data Flow

1. Payment request → HTTP 202 → Redis queue
2. Worker dequeues → Routes to healthy processor → Stores result in PostgreSQL
3. Health manager continuously updates processor status in Redis
4. Summary endpoint aggregates from PostgreSQL with time-based filtering

### Key Design Patterns

- **Leader Election**: Single instance performs health checks to respect rate limits
- **Circuit Breaker**: Processors marked as failing are avoided
- **Async Processing**: Decouples request handling from payment processing
- **Connection Pooling**: Reuses HTTP connections to payment processors

### Database Schema

PostgreSQL `payments` table:
- `correlation_id`: UUID primary key
- `amount`: NUMERIC(10,2) 
- `processed_at`: TIMESTAMPTZ (indexed)
- `processed_by`: payment_processor_type ENUM (indexed)

### Testing Strategy

The official test suite (`./test.sh`) uses K6 to simulate load and measure:
- Transaction success rate
- Response time distribution
- Overall profit calculation

Tests are configured in `rinha-test/` with scenarios for different load patterns.