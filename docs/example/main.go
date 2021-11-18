package main

import (
	domain "github.com/freundallein/scheduler/pkg"
	"github.com/freundallein/scheduler/pkg/client"
	"github.com/freundallein/scheduler/pkg/utils"
	log "github.com/freundallein/scheduler/pkg/utils/logging"
	"github.com/google/uuid"
	"time"
)

const (
	logLevelKey    = "LOG_LEVEL"
	tokenKey       = "TOKEN"
	workerTokenKey = "WORKER_TOKEN"
)

func worker() {
	token := utils.GetEnv(workerTokenKey, "token")
	worker := client.NewWorker(
		"0.0.0.0:8000",
		time.Second,
		client.WithWorkerToken(token),
	)
	for {
		tasks, err := worker.Claim(2)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("worker_claim_error")
		}
		for idx, task := range tasks {
			payload := task.Payload
			log.WithFields(log.Fields{
				"uid":     task.ID,
				"payload": payload,
			}).Error("worker_processing_task")
			//err := worker.Fail(task.ID, *task.ClaimID, "nobody is at home")
			err := worker.Succeed(task.ID, *task.ClaimID, map[string]interface{}{
				"idx": idx,
			})
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("worker_fail_error")
			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	logLevel := utils.GetEnv(logLevelKey, "debug")
	log.Init("example", logLevel)
	log.Info("init_service")
	token := utils.GetEnv(tokenKey, "token")

	service := client.NewScheduler(
		"127.0.0.1:8000",
		time.Second,
		client.WithToken(token),
	)
	go worker()
	uids := map[uuid.UUID]struct{}{}
	for i := 0; i < 10; i++ {
		uid, err := service.Set(time.Now(), time.Now().Add(time.Hour), map[string]interface{}{
			"source": "example.com",
			"type":   "parse",
			"number": i,
		})
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("scheduler_set_failed")
			panic(err)
		}
		uids[*uid] = struct{}{}
		log.WithFields(log.Fields{
			"uid": uid,
		}).Info("task_was_set")
	}
	for len(uids) > 0 {
		for uid := range uids {
			task, err := service.Get(uid)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("scheduler_get_failed")
				panic(err)
			}
			if task.State == domain.StateSucceeded {
				log.WithFields(log.Fields{
					"uid":    uid,
					"state":  task.State,
					"result": task.Result,
				}).Info("task_was_succeeded")
				delete(uids, uid)
			} else {
				log.WithFields(log.Fields{
					"uid":   uid,
					"state": task.State,
				}).Info("task_is_processing")
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

}
