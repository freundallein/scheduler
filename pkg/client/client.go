package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	domain "github.com/freundallein/scheduler/pkg"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"time"
)

type rpcResponse struct {
	ID     string                 `json:"id"`
	Error  string                 `json:"error,omitempty"`
}

type Scheduler struct {
	address     string
	accessToken string
	httpcli     *http.Client
}

func NewScheduler(address string, timeout time.Duration, opts ...Option) *Scheduler {
	client := &http.Client{
		Timeout: timeout,
	}
	sched := &Scheduler{
		address: address,
		httpcli: client,
	}
	for _, opt := range opts {
		opt(sched)
	}
	return sched
}

type setResponse struct {
	rpcResponse
	Result map[string]interface{} `json:"result"`
}

func (s *Scheduler) Set(executeAt, deadline time.Time, payload map[string]interface{}) (*uuid.UUID, error) {
	taskID := uuid.New()
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Scheduler.Set",
		"id":      "1",
		"params": []map[string]interface{}{
			{
				"id":        taskID,
				"executeAt": executeAt,
				"deadline":  deadline,
				"payload":   payload,
			},
		},
	})
	if err != nil {
		return &taskID, err
	}
	url := fmt.Sprintf("http://%s/rpc/v0", s.address)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return &taskID, err
	}
	request.Header.Set("Content-Type", "application/json")
	if s.accessToken != "" {
		request.Header.Set("Auth", s.accessToken)
	}
	if err != nil {
		return &taskID, err
	}
	resp, err := s.httpcli.Do(request)
	if err != nil {
		return &taskID, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &taskID, err
	}
	var response setResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &taskID, err
	}
	if response.Error != "" {
		return &taskID, errors.New(response.Error)
	}
	return &taskID, nil
}

type getResponse struct {
	rpcResponse
	Result struct{
		Meta map[string]interface{} `json:"meta"`
		Task *domain.Task  `json:"task"`
	} `json:"result"`
}

func (s *Scheduler) Get(id uuid.UUID) (*domain.Task, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Scheduler.Get",
		"id":      "1",
		"params": []map[string]interface{}{
			{
				"id": id,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s/rpc/v0", s.address)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if s.accessToken != "" {
		request.Header.Set("Auth", s.accessToken)
	}
	if err != nil {
		return nil, err
	}
	resp, err := s.httpcli.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response getResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result.Task, nil
}

type Worker struct{}

func NewWorker() *Worker {
	return &Worker{}
}
