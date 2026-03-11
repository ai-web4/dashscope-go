package dashscope

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL      = "https://dashscope.aliyuncs.com"
	defaultPollInterval = 15 * time.Second
	asyncHeaderValue    = "enable"
)

type Region string

const (
	RegionBeijing   Region = "beijing"
	RegionSingapore Region = "singapore"
	RegionVirginia  Region = "virginia"
)

func (r Region) BaseURL() string {
	switch r {
	case RegionSingapore:
		return "https://dashscope-intl.aliyuncs.com"
	case RegionVirginia:
		return "https://dashscope-us.aliyuncs.com"
	case RegionBeijing, "":
		return defaultBaseURL
	default:
		return defaultBaseURL
	}
}

type Config struct {
	APIKey       string
	BaseURL      string
	Region       Region
	HTTPClient   *http.Client
	PollInterval time.Duration
}

type Client struct {
	apiKey       string
	baseURL      string
	httpClient   *http.Client
	pollInterval time.Duration
}

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("api key 不能为空")
	}

	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = cfg.Region.BaseURL()
	}
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 120 * time.Second}
	}

	pollInterval := cfg.PollInterval
	if pollInterval <= 0 {
		pollInterval = defaultPollInterval
	}

	return &Client{
		apiKey:       cfg.APIKey,
		baseURL:      strings.TrimRight(baseURL, "/"),
		httpClient:   httpClient,
		pollInterval: pollInterval,
	}, nil
}

func (c *Client) SubmitVideoGeneration(ctx context.Context, req GenerationRequest) (*CreateTaskResponse, error) {
	return c.submitTask(ctx, ServiceVideoGeneration, req)
}

func (c *Client) SubmitImageToVideo(ctx context.Context, req GenerationRequest) (*CreateTaskResponse, error) {
	return c.submitTask(ctx, ServiceImageToVideo, req)
}

func (c *Client) SubmitAndWaitVideoGeneration(ctx context.Context, req GenerationRequest) (*TaskResponse, error) {
	created, err := c.SubmitVideoGeneration(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.WaitTask(ctx, created.Output.TaskID)
}

func (c *Client) SubmitAndWaitImageToVideo(ctx context.Context, req GenerationRequest) (*TaskResponse, error) {
	created, err := c.SubmitImageToVideo(ctx, req)
	if err != nil {
		return nil, err
	}
	return c.WaitTask(ctx, created.Output.TaskID)
}

// CallSync 发起同步 POST 请求（用于 detect 等同步接口）。
// path 为完整 API 路径，如 "/api/v1/services/aigc/image2video/video-synthesis"。
func (c *Client) CallSync(ctx context.Context, path string, reqBody any, respBody any) error {
	return c.doJSON(ctx, http.MethodPost, path, reqBody, respBody, false)
}

func (c *Client) submitTask(ctx context.Context, service Service, reqBody GenerationRequest) (*CreateTaskResponse, error) {
	var resp CreateTaskResponse
	if err := c.doJSON(ctx, http.MethodPost, service.Path(), reqBody, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	requestBody any,
	responseBody any,
	async bool,
) error {
	var reader io.Reader
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("序列化请求失败: %w", err)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if async {
		req.Header.Set("X-DashScope-Async", asyncHeaderValue)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求 DashScope 失败: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return parseAPIError(resp.StatusCode, raw)
	}

	if responseBody == nil || len(raw) == 0 {
		return nil
	}

	if err := json.Unmarshal(raw, responseBody); err != nil {
		return fmt.Errorf("解析响应失败: %w; body=%s", err, string(raw))
	}

	return nil
}
