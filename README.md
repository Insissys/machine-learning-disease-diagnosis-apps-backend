# Gin Backend – Disease Diagnosis App

This is the backend service built with Gin (Golang) for handling API requests, user management, medical records, and communicating with a FastAPI-based ML prediction service.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.21+
- Docker (for containerization)
- MySQL (or your DB of choice)

## Installation

Install dependencies:

```bash
# Gin (web framework)
go get -u github.com/gin-gonic/gin

# GORM (ORM) + MySQL driver
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql

# Resty (HTTP client)
go get -u github.com/go-resty/resty/v2

# YAML v3 (for config)
go get -u gopkg.in/yaml.v3
```

## Run with Docker

To build and run the FastAPI ML service in a Docker container:

### 1. Build the Docker Image

```bash
docker build -t gin-backend .
```

### 2. Run the Container

```bash
docker run -p 8080:8080 gin-backend
```
