package dashscope

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (c *Client) GetTask(ctx context.Context, taskID string) (*TaskResponse, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, fmt.Errorf("taskID 不能为空")
	}

	var resp TaskResponse
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/tasks/"+taskID, nil, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) WaitTask(ctx context.Context, taskID string) (*TaskResponse, error) {
	return c.WaitTaskWithInterval(ctx, taskID, c.pollInterval)
}

func (c *Client) WaitTaskWithInterval(ctx context.Context, taskID string, interval time.Duration) (*TaskResponse, error) {
	if interval <= 0 {
		interval = defaultPollInterval
	}

	for {
		resp, err := c.GetTask(ctx, taskID)
		if err != nil {
			return nil, err
		}

		if !resp.IsDone() {
			if err := sleepContext(ctx, interval); err != nil {
				return nil, err
			}
			continue
		}

		if resp.IsSuccess() {
			return resp, nil
		}

		return resp, &TaskFailedError{
			TaskID:     resp.Output.TaskID,
			TaskStatus: resp.Output.TaskStatus,
			Code:       resp.Output.Code,
			Message:    resp.Output.Message,
		}
	}
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
