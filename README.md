# Technical Coding Challenge

REST service that supports the management of a device database.

## Requirements

- Docker
- Go (v1.18 or higher - required for running unit tests)

## Setup

1. Run the service using `docker-compose up --build`.
2. To run the unit tests, use `go test ./... -v `.

## Technologies Used

- Golang
- PostgreSQL
- Docker

## Endpoints

- You can check and try out every endpoint with Swagger. With the service running, it is accessible via [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)