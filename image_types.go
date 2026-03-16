package dashscope

// -------- 图像生成服务端点 --------

type ImageService string

const (
	// V1版（wanx-v1）：仅异步
	ServiceTextToImage ImageService = "text2image"
	// V2版（wan2.2/2.5/2.6）：异步
	ServiceImageGeneration ImageService = "image-generation"
	// V2.6 同步接口
	ServiceMultimodalGeneration ImageService = "multimodal-generation"
)

func (s ImageService) Path() string {
	switch s {
	case ServiceTextToImage:
		return "/api/v1/services/aigc/text2image/image-synthesis"
	case ServiceImageGeneration:
		return "/api/v1/services/aigc/image-generation/generation"
	case ServiceMultimodalGeneration:
		return "/api/v1/services/aigc/multimodal-generation/generation"
	default:
		return "/api/v1/services/aigc/image-generation/generation"
	}
}

// -------- 图像模型名称 --------

const (
	// 文生图 V2（推荐）
	ModelWan26T2I        Model = "wan2.6-t2i"
	ModelWan25T2IPreview Model = "wan2.5-t2i-preview"
	ModelWan22T2IFlash   Model = "wan2.2-t2i-flash"
	ModelWan22T2IPlus    Model = "wan2.2-t2i-plus"
	ModelWanx21T2ITurbo  Model = "wanx2.1-t2i-turbo"
	ModelWanx21T2IPlus   Model = "wanx2.1-t2i-plus"
	ModelWanx20T2ITurbo  Model = "wanx2.0-t2i-turbo"

	// 图像生成与编辑（multimodal-generation 同步接口）
	ModelWan26Image Model = "wan2.6-image"

	// 文生图 V1（旧版）
	ModelWanxV1 Model = "wanx-v1"

	// 图像编辑 wan2.5
	ModelWan25ImageEdit Model = "wan2.5-image-edit"
)

// -------- 图像风格（V1版 wanx-v1 专用） --------

type ImageStyle string

const (
	ImageStyleAuto         ImageStyle = "<auto>"              // 自动
	ImageStylePhotography  ImageStyle = "<photography>"       // 摄影
	ImageStylePortrait     ImageStyle = "<portrait>"          // 人像写真
	ImageStyle3DCartoon    ImageStyle = "<3d cartoon>"        // 3D卡通
	ImageStyleAnime        ImageStyle = "<anime>"             // 动画
	ImageStyleOilPainting  ImageStyle = "<oil painting>"      // 油画
	ImageStyleWatercolor   ImageStyle = "<watercolor>"        // 水彩
	ImageStyleSketch       ImageStyle = "<sketch>"            // 素描
	ImageStyleChinesePaint ImageStyle = "<chinese painting>"  // 中国画
	ImageStyleFlat         ImageStyle = "<flat illustration>" // 扁平插画
)

// RefMode 参考图生成模式（V1版专用）
type RefMode string

const (
	RefModeRepaint RefMode = "repaint" // 基于参考图内容
	RefModeRefOnly RefMode = "refonly" // 基于参考图风格
)

// -------- V1版文生图请求（wanx-v1） --------

// ImageSynthesisRequest V1版文生图请求
// 端点: /api/v1/services/aigc/text2image/image-synthesis
type ImageSynthesisRequest struct {
	Model      Model                    `json:"model"`
	Input      ImageSynthesisInput      `json:"input"`
	Parameters ImageSynthesisParameters `json:"parameters,omitempty"`
}

// ImageSynthesisInput V1版输入
type ImageSynthesisInput struct {
	Prompt         string `json:"prompt"`                    // 正向提示词
	NegativePrompt string `json:"negative_prompt,omitempty"` // 反向提示词
	RefImg         string `json:"ref_img,omitempty"`         // 参考图 URL
}

// ImageSynthesisParameters V1版参数
type ImageSynthesisParameters struct {
	Style       ImageStyle `json:"style,omitempty"`        // 图像风格
	Size        string     `json:"size,omitempty"`         // 分辨率，如 1024*1024
	N           int        `json:"n,omitempty"`            // 生成数量 1~4
	Seed        *int64     `json:"seed,omitempty"`         // 随机种子
	RefStrength *float64   `json:"ref_strength,omitempty"` // 参考图相似度 [0.0, 1.0]
	RefMode     RefMode    `json:"ref_mode,omitempty"`     // 参考图模式
}

// -------- V2版文生图/图像编辑请求（messages 格式） --------

// ImageGenerationRequest V2版图像生成请求（wan2.2~2.6）
// 异步端点: /api/v1/services/aigc/image-generation/generation
// 同步端点: /api/v1/services/aigc/multimodal-generation/generation（仅 wan2.6）
type ImageGenerationRequest struct {
	Model      Model                     `json:"model"`
	Input      ImageGenerationInput      `json:"input"`
	Parameters ImageGenerationParameters `json:"parameters,omitempty"`
}

// ImageGenerationInput V2版输入（messages 格式）
type ImageGenerationInput struct {
	Messages []ImageMessage `json:"messages"`
}

// ImageMessage 图像消息
type ImageMessage struct {
	Role    string         `json:"role"`    // 固定为 "user"
	Content []ImageContent `json:"content"` // 文本和图片内容
}

// ImageContent 消息内容项
type ImageContent struct {
	Text  string `json:"text,omitempty"`  // 文本提示词
	Image string `json:"image,omitempty"` // 图片 URL 或 base64
	Type  string `json:"type,omitempty"`  // 输出时包含：text / image
}

// ImageGenerationParameters V2版参数
type ImageGenerationParameters struct {
	NegativePrompt   string `json:"negative_prompt,omitempty"`   // 反向提示词
	Size             string `json:"size,omitempty"`              // 分辨率
	N                int    `json:"n,omitempty"`                 // 生成数量 1~4
	PromptExtend     *bool  `json:"prompt_extend,omitempty"`     // Prompt 智能改写
	Watermark        *bool  `json:"watermark,omitempty"`         // 是否添加水印
	Seed             *int64 `json:"seed,omitempty"`              // 随机种子
	EnableInterleave *bool  `json:"enable_interleave,omitempty"` // 图文混排模式
	MaxImages        int    `json:"max_images,omitempty"`        // 图文混排模式下最大图片数
	Stream           *bool  `json:"stream,omitempty"`            // 流式输出（图文混排需要）
}

// -------- 图像生成响应 --------

// ImageGenerationResponse V2版同步响应 / 异步查询响应
type ImageGenerationResponse struct {
	RequestID string      `json:"request_id,omitempty"`
	Output    ImageOutput `json:"output"`
	Usage     *ImageUsage `json:"usage,omitempty"`
	Code      string      `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
}

// ImageOutput 图像输出
type ImageOutput struct {
	// 异步任务字段
	TaskID     string     `json:"task_id,omitempty"`
	TaskStatus TaskStatus `json:"task_status,omitempty"`

	// 时间戳
	SubmitTime    string `json:"submit_time,omitempty"`
	ScheduledTime string `json:"scheduled_time,omitempty"`
	EndTime       string `json:"end_time,omitempty"`

	// V2版输出（choices 格式）
	Choices  []ImageChoice `json:"choices,omitempty"`
	Finished *bool         `json:"finished,omitempty"`

	// V1版输出（results 格式）
	Results     []ImageResult `json:"results,omitempty"`
	TaskMetrics *TaskMetrics  `json:"task_metrics,omitempty"`

	// 错误信息
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ImageChoice V2版输出内容
type ImageChoice struct {
	FinishReason string       `json:"finish_reason,omitempty"`
	Message      ImageMessage `json:"message,omitempty"`
}

// ImageResult V1版输出内容
type ImageResult struct {
	URL     string `json:"url,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// TaskMetrics V1版任务统计
type TaskMetrics struct {
	Total     int `json:"TOTAL"`
	Succeeded int `json:"SUCCEEDED"`
	Failed    int `json:"FAILED"`
}

// ImageUsage 图像用量统计
type ImageUsage struct {
	ImageCount   int    `json:"image_count,omitempty"`
	Size         string `json:"size,omitempty"`
	InputTokens  int    `json:"input_tokens,omitempty"`
	OutputTokens int    `json:"output_tokens,omitempty"`
	TotalTokens  int    `json:"total_tokens,omitempty"`
}

// -------- 辅助方法 --------

// IsDone 判断异步任务是否完成
func (r ImageGenerationResponse) IsDone() bool {
	switch r.Output.TaskStatus {
	case TaskStatusSucceeded, TaskStatusFailed, TaskStatusCanceled, TaskStatusUnknown:
		return true
	default:
		return false
	}
}

// IsSuccess 判断异步任务是否成功
func (r ImageGenerationResponse) IsSuccess() bool {
	return r.Output.TaskStatus == TaskStatusSucceeded
}

// GetImageURLs 获取所有生成的图片 URL
func (r ImageGenerationResponse) GetImageURLs() []string {
	var urls []string
	// V2 格式（choices）
	for _, choice := range r.Output.Choices {
		for _, c := range choice.Message.Content {
			if c.Image != "" {
				urls = append(urls, c.Image)
			}
		}
	}
	// V1 格式（results）
	for _, result := range r.Output.Results {
		if result.URL != "" {
			urls = append(urls, result.URL)
		}
	}
	return urls
}

// GetFirstImageURL 获取第一张生成图片的 URL
func (r ImageGenerationResponse) GetFirstImageURL() string {
	urls := r.GetImageURLs()
	if len(urls) > 0 {
		return urls[0]
	}
	return ""
}
