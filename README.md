# scheduler

## Description

Provides ability to manage background tasks with delay.

You set a task with any payload and can poll it state until completion.

In the same time worker will claim this task, try to process it and push results back to the service.

## API
Service implements JSON-RPC 2.0 and contains 2 types of API:
- scheduler (public) used for setting a task;
- worker (private) for processing a task.

[Do you want to know more?](https://github.com/freundallein/scheduler/docs/api_v0.md)
## Design

## SLA

## Deployment

## Load test

## TODO
- [ ] fill README.md
- [ ] extend unittests
- [ ] lint all comments
- [ ] add metrics
- [ ] add grafana dashboard
- [ ] add prolong operation for worker
- [ ] benchmarks