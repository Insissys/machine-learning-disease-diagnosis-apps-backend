run:
	@go run cmd/server/main.go

tidy:
	@go mod tidy

migrate:
	@go run cmd/migration/main.go

config:
	@go run cmd/config/main.go