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
	"strconv"
	"time"
)

type rpcResponse struct {
	// ID is current response ID.
	// Should be equals request ID.
	ID string `json:"id"`
	// Error shows error message
	Error string `json:"error,omitempty"`
}

// Scheduler implements client for a public interface domain.Scheduler.
type Scheduler struct {
	url         string
	accessToken string
	httpcli     *http.Client
}

// NewScheduler returns an instance of Scheduler.
func NewScheduler(address string, timeout time.Duration, opts ...SchedulerOption) *Scheduler {
	client := &http.Client{
		Timeout: timeout,
	}
	scheduler := &Scheduler{
		url:     fmt.Sprintf("http://%s/rpc/v0", address),
		httpcli: client,
	}
	for _, opt := range opts {
		opt(scheduler)
	}
	return scheduler
}

func (s *Scheduler) makeRequest(payload map[string]interface{}) (interface{}, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", s.url, bytes.NewBuffer(requestBody))
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
	return body, nil
}

type setResponse struct {
	rpcResponse
	Result map[string]interface{} `json:"result"`
}

// Set allows to enqueue task.
func (s *Scheduler) Set(executeAt, deadline time.Time, payload map[string]interface{}) (*uuid.UUID, error) {
	taskID := uuid.New()
	request := map[string]interface{}{
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
	}
	responseBody, err := s.makeRequest(request)
	var response setResponse
	err = json.Unmarshal(responseBody.([]byte), &response)
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
	Result struct {
		Meta map[string]interface{} `json:"meta"`
		Task *domain.Task           `json:"task"`
	} `json:"result"`
}

// Get allows to poll a task state.
func (s *Scheduler) Get(id uuid.UUID) (*domain.Task, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Scheduler.Get",
		"id":      "1",
		"params": []map[string]interface{}{
			{
				"id": id,
			},
		},
	}
	responseBody, err := s.makeRequest(request)
	var response getResponse
	err = json.Unmarshal(responseBody.([]byte), &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result.Task, nil
}

// Worker implements client for a private interface domain.Worker.
type Worker struct {
	url         string
	accessToken string
	httpcli     *http.Client
}

// NewWorker returns an instance of Worker.
func NewWorker(address string, timeout time.Duration, opts ...WorkerOption) *Worker {
	client := &http.Client{
		Timeout: timeout,
	}
	worker := &Worker{
		url:     fmt.Sprintf("http://%s/worker/v0", address),
		httpcli: client,
	}
	for _, opt := range opts {
		opt(worker)
	}
	return worker
}

func (w *Worker) makeRequest(payload map[string]interface{}) (interface{}, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", w.url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if w.accessToken != "" {
		request.Header.Set("Auth", w.accessToken)
	}
	if err != nil {
		return nil, err
	}
	resp, err := w.httpcli.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type claimResponse struct {
	rpcResponse
	Result struct {
		Count int            `json:"count"`
		Tasks []*domain.Task `json:"tasks"`
	} `json:"result"`
}

// Claim takes a list of tasks.
func (w *Worker) Claim(amount int) ([]*domain.Task, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Worker.Claim",
		"id":      "1",
		"params": []map[string]interface{}{
			{
				"amount": strconv.Itoa(amount),
			},
		},
	}
	responseBody, err := w.makeRequest(request)
	var response claimResponse
	err = json.Unmarshal(responseBody.([]byte), &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	return response.Result.Tasks, nil
}

type succeedResponse struct {
	rpcResponse
	Result struct {
		Message string `json:"message"`
	} `json:"result"`
}

// Succeed marks a task as done.
func (w *Worker) Succeed(id, claimID uuid.UUID, result map[string]interface{}) error {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Worker.Succeed",
		"id":      "1",
		"params": []map[string]interface{}{
			{
				"id":      id,
				"claimID": claimID,
				"result":  result,
			},
		},
	}
	responseBody, err := w.makeRequest(request)
	if err != nil {
		return err
	}
	var response succeedResponse
	err = json.Unmarshal(responseBody.([]byte), &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	if response.Result.Message != "success" {
		return errors.New("succeed op was unsuccessful")
	}
	return nil
}

type failResponse struct {
	rpcResponse
	Result struct {
		Message string `json:"message"`
	} `json:"result"`
}

// Fail marks a task as failed.
func (w *Worker) Fail(id, claimID uuid.UUID, reason string) error {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "Worker.Fail",
		"id":      "1",
		"params": []map[string]interface{}{
			{
				"id":      id,
				"claimID": claimID,
				"reason":  reason,
			},
		},
	}
	responseBody, err := w.makeRequest(request)
	if err != nil {
		return err
	}
	fmt.Println(string(responseBody.([]byte)))
	var response failResponse
	err = json.Unmarshal(responseBody.([]byte), &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	if response.Result.Message != "success" {
		return errors.New("fail op was unsuccessful")
	}
	return nil
}
