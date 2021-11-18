package main

import (
	"fmt"
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
	workerTokenKey = "WRK_TOKEN"
)

func worker() {
	//worker := client.NewWorker("0.0.0.0:8000", "token", time.Second)
}

func main() {
	logLevel := utils.GetEnv(logLevelKey, "debug")
	log.Init("example", logLevel)
	log.Info("init_service")
	token := utils.GetEnv(tokenKey, "token")
	//workerToken := utils.GetEnv(workerTokenKey, "token")

	service := client.NewScheduler(
		"0.0.0.0:8000",
		time.Second,
		client.WithToken(token),
	)
	go worker()
	uids := map[uuid.UUID]struct{}{}
	for i := 0; i < 2; i++ {
		uid, err := service.Set(time.Now(), time.Now().Add(time.Hour), map[string]interface{}{
			"source": "example.com",
			"type":   "parse",
			"number": i,
		})
		if err != nil {
			panic(err)
		}
		uids[*uid] = struct{}{}
	}
	for len(uids) > 0 {
		for uid := range uids {
			task, err := service.Get(uid)
			if err != nil {
				panic(err)
			}
			if task.State == domain.StateSucceeded {
				fmt.Println(uid, task.State, task.Result)
				delete(uids, uid)
			} else {
				fmt.Println(uid, task.State)
			}
			time.Sleep(time.Second)
		}
	}

}
