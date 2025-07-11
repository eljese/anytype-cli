package daemon

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/anyproto/anytype-cli/core/config"
)

const (
	defaultTimeout = 5 * time.Second
)

// SendTaskStart sends a start request for a given task.
func SendTaskStart(task string, params map[string]string) (*TaskResponse, error) {
	reqData := TaskRequest{Task: task, Params: params}
	b, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Post(config.DaemonHTTPURL+"/task/start", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var taskResp TaskResponse
	err = json.Unmarshal(body, &taskResp)
	if err != nil {
		return nil, err
	}
	if taskResp.Status == "error" {
		return &taskResp, errors.New(taskResp.Error)
	}
	return &taskResp, nil
}

// SendTaskStop sends a stop request for a given task.
func SendTaskStop(task string, params map[string]string) (*TaskResponse, error) {
	reqData := TaskRequest{Task: task, Params: params}
	b, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Post(config.DaemonHTTPURL+"/task/stop", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var taskResp TaskResponse
	err = json.Unmarshal(body, &taskResp)
	if err != nil {
		return nil, err
	}
	if taskResp.Status == "error" {
		return &taskResp, errors.New(taskResp.Error)
	}
	return &taskResp, nil

}

// SendTaskStatus sends a status request for a given task.
func SendTaskStatus(task string) (*TaskResponse, error) {
	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Get(config.DaemonHTTPURL + "/task/status?task=" + task)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var taskResp TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return nil, err
	}
	if taskResp.Status == "error" {
		return &taskResp, errors.New(taskResp.Error)
	}
	return &taskResp, nil
}
