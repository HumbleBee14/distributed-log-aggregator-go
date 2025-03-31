# Distributed Log Aggregator (Golang)

A high-performance log aggregation service for microservices written in Go.

## Overview

This distributed log aggregation service provides a RESTful API for collecting logs from multiple microservices and querying them efficiently. The service is designed to be horizontally scalable, thread-safe, and suitable for high-throughput production environments.

## Features

- **Log Ingestion API**: Submit logs from any microservice
- **Query API**: Retrieve logs by service name and time range
- **Efficient Storage**: Uses Redis with automatic expiration (TTL)
- **Robust Configuration**: Multiple configuration sources with clear precedence
- **Production-Ready Logging**: Configurable log levels and outputs
- **Graceful Shutdown**: Clean termination with connection handling

## Architecture

The service uses a layered architecture with clean separation of concerns:

- **API Layer**: RESTful endpoints built with Go's standard HTTP package
- **Storage Layer**: Redis implementation for high-performance in-memory storage
- **Configuration**: Environment-based configuration with multiple sources
- **Logging**: Structured logging with level-based filtering

## Requirements

- Go 1.19+  (Using 1.24.0)
- Redis server (or Redis Cloud account)

## Installation

Clone the repository:

```bash
git clone https://github.com/HumbleBee14/distributed-log-aggregator-go.git
cd distributed-log-aggregator-go
```

Install dependencies:

```bash
go mod tidy
```

## Configuration

Create a `.env` file in the project root:

You can refer the .env.example or use the bash script to generate the env file (setup_env.sh)

Configuration has the following precedence:
1. Environment variables (`.env` file)
2. Default values

## Running the Service

Start the service with:

```bash
go run cmd/server/main.go
```

For production deployment:

```bash
go build -o log-aggregator cmd/server/main.go
./log-aggregator
```

## API Usage

### Log Ingestion

```bash
# Submit a log entry
curl -X POST http://localhost:8080/logs \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "auth-service",
    "timestamp": "2025-03-17T10:15:00Z",
    "message": "User login successful"
  }'
```

### Log Querying

```bash
# Query logs by service and time range
curl "http://localhost:8080/logs?service=auth-service&start=2025-03-17T10:00:00Z&end=2025-03-17T10:30:00Z"
```

## Scaling

For high-volume environments, consider:

1. Running multiple instances behind a load balancer
2. Configuring Redis in cluster mode
3. Setting up log rotation for persistent storage

## License
<!-- MIT -->

## Author

👤 **DINESH**  
🔗 [GitHub](https://github.com/HumbleBee14/distributed-log-aggregator-go)
