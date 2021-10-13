export LOG_LEVEL=debug
export PORT=8000
export OPS_PORT=8001
export PG_DSN=postgres://scheduler:scheduler@192.168.64.6:5432/scheduler

run:
	go run cmd/main.go

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -o ./bin/scheduler ./cmd/

build_docker:
	docker build --tag=freundallein/scheduler:latest --file=./docker/Dockerfile .

