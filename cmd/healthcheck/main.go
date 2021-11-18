package main

import (
	"fmt"
	"github.com/freundallein/scheduler/pkg/utils"
	"net/http"
	"os"
)

const (
	opsPortKey = "OPS_PORT"
)

func main() {
	opsPort := utils.GetEnv(opsPortKey, "8001")
	_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/ops/healthcheck", opsPort))
	if err != nil {
		os.Exit(1)
	}
}
