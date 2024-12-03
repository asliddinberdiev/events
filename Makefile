-include .env
  
DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)"

tidy:
	@go mod tidy
	@go mod vendor

run:
	@go run cmd/main.go

migration:
	@migrate create -ext sql -dir ./migrations -seq $(name)

up:
	@migrate -path migrations -database "$(DB_URL)" -verbose up

down:
	@migrate -path migrations -database "$(DB_URL)" -verbose down

compose-up:
	@docker-compose up -d

compose-down:
	@docker-compose down
