# scheduler

## Description

Provides ability to manage background tasks with delay.

You set a task with any payload and can poll it state until completion.

In the same time, worker will claim this task, try to process it and push results back to the service.

## API
Service implements JSON-RPC 2.0 and contains 2 types of API:
- scheduler (public) used for setting a task;
- worker (private) for processing a task.

[Do you want to know more?](https://github.com/freundallein/scheduler/blob/master/docs/api_v0.md)

## Design
![Alt text](https://github.com/freundallein/scheduler/blob/master/docs/design.png)

## SLA

...

## Deployment
Service uses environment variables for configuration.

Most important parameter is `DB_DSN` - address of a postgres instance.

You can see a full list of parameters in `Makefile`.

### Docker
To get image run  
`docker pull ghcr.io/freundallein/scheduler:latest`


[Docker-compose example](https://github.com/freundallein/scheduler/blob/master/docker/docker-compose.yml)

### Build 

You can build scheduler via `Makefile`:
```
make build
```
Binary file will be delivered to `./bin/scheduler`.

## Example
For proper work you will need a client and a worker.

You can run example with client and worker via `Makefile`:
```
make run
make example
```
Ensure, that you change variables in `Makefile`, especially `DB_DSN`

[Example code](https://github.com/freundallein/scheduler/blob/master/docs/example/main.go)

## Load test

...

## TODO
- [ ] extend unittests
- [ ] add grafana dashboard
- [ ] add prolong operation for worker
- [ ] benchmarks
- [ ] fill README.md