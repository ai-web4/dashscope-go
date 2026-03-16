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
	Headers      map[string]string // 自定义请求头，会附加到每个请求中
	PollInterval time.Duration
}

type Client struct {
	apiKey       string
	baseURL      string
	httpClient   *http.Client
	headers      map[string]string
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
		headers:      cfg.Headers,
		pollInterval: pollInterval,
	}, nil
}

// SubmitAsync 提交异步任务到任意 DashScope API 路径（透传 raw JSON body）
// 返回 CreateTaskResponse，需后续轮询获取结果
func (c *Client) SubmitAsync(ctx context.Context, path string, body json.RawMessage, header ...map[string]string) (*CreateTaskResponse, error) {
	var resp CreateTaskResponse
	if err := c.doJSON(ctx, http.MethodPost, path, body, &resp, true, mergeHeaders(header)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CallRaw 调用任意 DashScope API 路径（同步，透传 raw JSON body）
// 返回原始 JSON 响应
func (c *Client) CallRaw(ctx context.Context, path string, body json.RawMessage, header ...map[string]string) (json.RawMessage, error) {
	var resp json.RawMessage
	if err := c.doJSON(ctx, http.MethodPost, path, body, &resp, false, mergeHeaders(header)); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) SubmitVideoGeneration(ctx context.Context, req GenerationRequest, header ...map[string]string) (*CreateTaskResponse, error) {
	return c.submitTask(ctx, ServiceVideoGeneration, req, mergeHeaders(header))
}

func (c *Client) SubmitImageToVideo(ctx context.Context, req GenerationRequest, header ...map[string]string) (*CreateTaskResponse, error) {
	return c.submitTask(ctx, ServiceImageToVideo, req, mergeHeaders(header))
}

func (c *Client) SubmitAndWaitVideoGeneration(ctx context.Context, req GenerationRequest, header ...map[string]string) (*TaskResponse, error) {
	created, err := c.SubmitVideoGeneration(ctx, req, mergeHeaders(header))
	if err != nil {
		return nil, err
	}
	return c.WaitTask(ctx, created.Output.TaskID)
}

func (c *Client) SubmitAndWaitImageToVideo(ctx context.Context, req GenerationRequest, header ...map[string]string) (*TaskResponse, error) {
	created, err := c.SubmitImageToVideo(ctx, req, mergeHeaders(header))
	if err != nil {
		return nil, err
	}
	return c.WaitTask(ctx, created.Output.TaskID)
}

// CallSync 发起同步 POST 请求（用于 detect 等同步接口）。
// path 为完整 API 路径，如 "/api/v1/services/aigc/image2video/video-synthesis"。
func (c *Client) CallSync(ctx context.Context, path string, reqBody any, respBody any) error {
	return c.doJSON(ctx, http.MethodPost, path, reqBody, respBody, false, nil)
}

func (c *Client) submitTask(ctx context.Context, service Service, reqBody GenerationRequest, header map[string]string) (*CreateTaskResponse, error) {
	var resp CreateTaskResponse
	if err := c.doJSON(ctx, http.MethodPost, service.Path(), reqBody, &resp, true, header); err != nil {
		return nil, err
	}
	return &resp, nil
}

// mergeHeaders 取 variadic header 参数的第一个值（用于可选参数模式）
func mergeHeaders(headers []map[string]string) map[string]string {
	if len(headers) > 0 {
		return headers[0]
	}
	return nil
}

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	requestBody any,
	responseBody any,
	async bool,
	header map[string]string,
) error {
	var reader io.Reader
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("序列化请求失败: %w", err)
		}
		reader = bytes.NewReader(data)
	}

	fmt.Println("doJSON", c.baseURL+path, requestBody)
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 优先级：调用方 header > 全局 c.headers > SDK 内部 header
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	for k, v := range header {
		req.Header.Set(k, v)
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
