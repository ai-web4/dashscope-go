package dashscope

import (
	"context"
	"net/http"
	"time"
)

// -------- V1版文生图（wanx-v1） --------

// SubmitTextToImage 提交 V1版文生图异步任务
// 端点: /api/v1/services/aigc/text2image/image-synthesis
func (c *Client) SubmitTextToImage(ctx context.Context, req ImageSynthesisRequest, header ...map[string]string) (*CreateTaskResponse, error) {
	var resp CreateTaskResponse
	if err := c.doJSON(ctx, http.MethodPost, ServiceTextToImage.Path(), req, &resp, true, mergeHeaders(header)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubmitAndWaitTextToImage 提交 V1版文生图并等待完成
func (c *Client) SubmitAndWaitTextToImage(ctx context.Context, req ImageSynthesisRequest, header ...map[string]string) (*TaskResponse, error) {
	created, err := c.SubmitTextToImage(ctx, req, mergeHeaders(header))
	if err != nil {
		return nil, err
	}
	return c.WaitTask(ctx, created.Output.TaskID)
}

// -------- V2版图像生成（wan2.2~2.6 异步） --------

// SubmitImageGeneration 提交 V2版图像生成异步任务
// 端点: /api/v1/services/aigc/image-generation/generation
func (c *Client) SubmitImageGeneration(ctx context.Context, req ImageGenerationRequest, header ...map[string]string) (*CreateTaskResponse, error) {
	var resp CreateTaskResponse
	if err := c.doJSON(ctx, http.MethodPost, ServiceImageGeneration.Path(), req, &resp, true, mergeHeaders(header)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubmitAndWaitImageGeneration 提交 V2版图像生成并等待完成
func (c *Client) SubmitAndWaitImageGeneration(ctx context.Context, req ImageGenerationRequest, header ...map[string]string) (*ImageGenerationResponse, error) {
	created, err := c.SubmitImageGeneration(ctx, req, mergeHeaders(header))
	if err != nil {
		return nil, err
	}
	return c.WaitImageTask(ctx, created.Output.TaskID)
}

// -------- V2.6 同步接口（wan2.6-t2i / wan2.6-image） --------

// GenerateImage 调用 V2.6 同步接口生成图像
// 端点: /api/v1/services/aigc/multimodal-generation/generation
func (c *Client) GenerateImage(ctx context.Context, req ImageGenerationRequest, header ...map[string]string) (*ImageGenerationResponse, error) {
	var resp ImageGenerationResponse
	if err := c.doJSON(ctx, http.MethodPost, ServiceMultimodalGeneration.Path(), req, &resp, false, mergeHeaders(header)); err != nil {
		return nil, err
	}
	return &resp, nil
}

// -------- 图像任务轮询 --------

// WaitImageTask 轮询等待图像生成任务完成（使用 ImageGenerationResponse 格式解析）
func (c *Client) WaitImageTask(ctx context.Context, taskID string) (*ImageGenerationResponse, error) {
	return c.WaitImageTaskWithInterval(ctx, taskID, c.pollInterval)
}

// WaitImageTaskWithInterval 自定义间隔轮询等待图像任务完成
func (c *Client) WaitImageTaskWithInterval(ctx context.Context, taskID string, interval time.Duration) (*ImageGenerationResponse, error) {
	for {
		var resp ImageGenerationResponse
		if err := c.doJSON(ctx, http.MethodGet, "/api/v1/tasks/"+taskID, nil, &resp, false, nil); err != nil {
			return nil, err
		}
		if resp.IsDone() {
			return &resp, nil
		}
		if err := sleepContext(ctx, interval); err != nil {
			return nil, err
		}
	}
}
