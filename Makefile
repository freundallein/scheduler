export LOG_LEVEL=debug
export PORT=8000
export OPS_PORT=8001
export DB_DSN=postgres://scheduler:scheduler@192.168.64.6:5432/scheduler
export TOKEN=token
export WORKER_TOKEN=workertoken
export STALE_HOURS=1

run:
	go run cmd/main.go

fmt:
	go fmt ./...

test:
	go test -v -race -count=1 -cover ./...

tidy:
	go mod tidy

build: build_scheduler build_healthcheck

build_scheduler:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -o ./bin/scheduler ./cmd/

build_healthcheck:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -o ./bin/healthcheck ./cmd/healthcheck

build_docker:
	docker build --tag=ghcr.io/freundallein/scheduler:latest --file=./docker/Dockerfile .

deliver:
	make build_docker
	docker push ghcr.io/freundallein/scheduler:latest

example:
	go run docs/example/main.go
