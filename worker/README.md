# Worker App

## Description

This is the worker service written in Go. It consumes vote messages from a Kafka topic and stores them in a PostgreSQL database for later retrieval by the result service.

## Prerequisites

- Go 1.24.2
- Kafka cluster
- PostgreSQL database
- Environment variables for database and Kafka configuration

## Installation

1. Clone the repository.
2. Navigate to the worker folder.
3. Run `go mod tidy` to download dependencies.

## Running the App

Set the following environment variables:

- POSTGRES_HOST: PostgreSQL host
- POSTGRES_PORT: PostgreSQL port
- POSTGRES_USER: Database username
- POSTGRES_PASSWORD: Database password
- POSTGRES_DB: Database name

Optional flags (can be set via command line or defaults):

- --brokerList: Kafka brokers (default: kafka:9092)
- --topic: Kafka topic (default: votes)
- --messageCountStart: Starting message count

Then run:

```
go run main.go
```

Or build the binary:

```
go build -o worker main.go
./worker
```

## Features

- Consumes vote messages from Kafka topic asynchronously
- Stores votes in a PostgreSQL table with unique IDs
- Handles database connection and table creation
- Supports graceful shutdown on interrupt signals