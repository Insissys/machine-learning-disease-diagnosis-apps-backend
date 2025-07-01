run:
	@go run cmd/server/main.go

build:
	@go build -o server ./cmd/server/main.go

tidy:
	@go mod tidy

migrate:
	@go run cmd/migration/main.go

drop:
	@go run cmd/migration/drop/main.go

seeder:
	@go run cmd/migration/seeder/main.go

config:
	@go run cmd/config/main.go