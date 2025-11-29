.PHONY: tests
tests:
	go test -v ./... -race -timeout 300s -cover

.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -o docs

lint:
	golangci-lint run

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" -verbose down

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

up:
	docker compose -f Docker-compose.yml up -d --build

down:
	docker compose down

restart:
	make stop
	make run