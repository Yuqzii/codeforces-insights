all: build run

ARGS := ""

build:
	docker compose build

run:
	docker compose up -d

stop:
	docker compose down

fetch: 
	docker compose run --build --rm fetcher $(ARGS)

lint:
	cd backend && golangci-lint run
	cd frontend && npm run lint

format:
	cd backend && gofmt -d -s -l
