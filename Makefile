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

CONTAINER_NAME=codeforces-insights-server-1

debug:
	docker run -it --rm \
		--pid=container:$(CONTAINER_NAME) \
		--net=container:$(CONTAINER_NAME) \
		--cap-add=SYS_PTRACE \
		golang:alpine \
		sh -c "apk add --no-cache delve htop && sh"

PROFILE_TIME=30
profile:
	@echo "Starting $(PROFILE_TIME)-second CPU profiling, wait until \"Fetching profile over HTTP from ...\" appears."
	@docker run --rm \
		--pid=container:$(CONTAINER_NAME) \
		--net=container:$(CONTAINER_NAME) \
		-v $(PWD):/out:z \
		golang:alpine \
		sh -c "go tool pprof -proto http://localhost:6060/debug/pprof/profile?seconds=$(PROFILE_TIME) > /out/cpu.pprof"
	@echo "Profile captured to ./cpu.pprof, opening flamegraph..."
	go tool pprof -http=:8000 cpu.pprof
