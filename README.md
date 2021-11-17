# scheduler

## Description

Provides ability to manage background tasks with delay.

You set a task with any payload and can poll it state until comletion.
In that time, worker will claim this task, process it and push results back to the service.

## API

### Scheduler

### Worker

## Design

## SLA

## Deployment

## Load test

## TODO
- [ ] db scheme creation on startup
- [ ] separate worker and scheduler interfaces
- [ ] fill README.md
- [ ] add client library
- [ ] add worker example
- [ ] add supervisor routine to delete old tasks or to manage partitions
- [ ] add prolong operation for worker
- [ ] add metrics
- [ ] add grafana dashboard
