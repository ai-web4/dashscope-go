package dashscope

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	StatusCode int    `json:"-"`
	RequestID  string `json:"request_id,omitempty"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Body       string `json:"-"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}

	if e.Code != "" || e.Message != "" {
		return fmt.Sprintf("dashscope api 错误: status=%d code=%s message=%s", e.StatusCode, e.Code, e.Message)
	}

	return fmt.Sprintf("dashscope api 错误: status=%d body=%s", e.StatusCode, e.Body)
}

type TaskFailedError struct {
	TaskID     string
	TaskStatus TaskStatus
	Code       string
	Message    string
}

func (e *TaskFailedError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf(
		"任务执行失败: task_id=%s status=%s code=%s message=%s",
		e.TaskID,
		e.TaskStatus,
		e.Code,
		e.Message,
	)
}

func parseAPIError(statusCode int, raw []byte) error {
	var apiErr APIError
	apiErr.StatusCode = statusCode
	apiErr.Body = string(raw)

	_ = json.Unmarshal(raw, &apiErr)
	return &apiErr
}
