.PHONY: run build test docker-up docker-down migrate

run:
	go run main.go

build:
	go build -o bin/main main.go

test:
	go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

migrate:
	go run main.go migrate

install-deps:
	go mod tidy
	go mod download

dev:
	air # Para hot reload (instale: go install github.com/cosmtrek/air@latest)
