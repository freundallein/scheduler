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
- [ ] add client library
- [ ] add worker example
- [ ] add supervisor routine to delete old tasks or to manage partitions
- [ ] fill README.md
- [ ] extend unittests
- [ ] add prolong operation for worker
- [ ] add metrics
- [ ] add grafana dashboard
