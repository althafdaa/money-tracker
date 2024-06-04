# Project Money Tracker

Golang API Application to track money transactions

## Docs

- [Bruno HTTP Client](https://docs.usebruno.com/)
  - import the collection from `docs/money-tracker-bruno.json`
- [Database](docs/database/db.dbml)
- [API Endpoints](docs/api-endpoints.md)

## Technology Stack

- Go v1.21.1
- Postgres SQL
- Docker
- Fiber
- GORM
- Google Cloud Build and Cloud Run

## Development

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Makefile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```
