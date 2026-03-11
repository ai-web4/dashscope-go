package dashscope

import "encoding/json"

// -------- 服务端点 --------

type Service string

const (
	ServiceVideoGeneration Service = "video-generation"
	ServiceImageToVideo    Service = "image2video"
)

func (s Service) Path() string {
	switch s {
	case ServiceImageToVideo:
		return "/api/v1/services/aigc/image2video/video-synthesis"
	case ServiceVideoGeneration, "":
		return "/api/v1/services/aigc/video-generation/video-synthesis"
	default:
		return "/api/v1/services/aigc/video-generation/video-synthesis"
	}
}

// -------- 枚举值 --------

// Size 视频分辨率（宽*高）
type Size string

const (
	// 480P 档位
	Size480P_832x480 Size = "832*480"  // 16:9
	Size480P_480x832 Size = "480*832"  // 9:16
	Size480P_624x624 Size = "624*624"  // 1:1
	// 720P 档位
	Size720P_1280x720 Size = "1280*720" // 16:9
	Size720P_720x1280 Size = "720*1280" // 9:16
	Size720P_960x960  Size = "960*960"  // 1:1
	Size720P_1088x832 Size = "1088*832" // 4:3
	Size720P_832x1088 Size = "832*1088" // 3:4
	// 1080P 档位
	Size1080P_1920x1080 Size = "1920*1080" // 16:9
	Size1080P_1080x1920 Size = "1080*1920" // 9:16
	Size1080P_1440x1440 Size = "1440*1440" // 1:1
	Size1080P_1632x1248 Size = "1632*1248" // 4:3
	Size1080P_1248x1632 Size = "1248*1632" // 3:4
)

// Resolution 视频分辨率档位
type Resolution string

const (
	Resolution480P  Resolution = "480P"
	Resolution720P  Resolution = "720P"
	Resolution1080P Resolution = "1080P"
)

// ShotType 镜头类型
type ShotType string

const (
	ShotTypeSingle ShotType = "single" // 单镜头（默认）
	ShotTypeMulti  ShotType = "multi"  // 多镜头
)

// VACEFunction VACE 视频编辑功能
type VACEFunction string

const (
	VACEFunctionImageReference   VACEFunction = "image_reference"   // 多图参考
	VACEFunctionVideoRepainting  VACEFunction = "video_repainting"  // 视频重绘
	VACEFunctionVideoEdit        VACEFunction = "video_edit"        // 局部编辑
	VACEFunctionVideoExtension   VACEFunction = "video_extension"   // 视频延展
	VACEFunctionVideoOutpainting VACEFunction = "video_outpainting" // 视频画面扩展
)

// ControlCondition 视频特征提取方式
type ControlCondition string

const (
	ControlConditionPoseBodyFace ControlCondition = "posebodyface" // 脸部表情+肢体动作
	ControlConditionPoseBody     ControlCondition = "posebody"     // 仅肢体动作
	ControlConditionDepth        ControlCondition = "depth"        // 构图和运动轮廓
	ControlConditionScribble     ControlCondition = "scribble"     // 线稿结构
)

// MaskType 掩码区域行为方式
type MaskType string

const (
	MaskTypeTracking MaskType = "tracking" // 动态跟随（默认）
	MaskTypeFixed    MaskType = "fixed"    // 固定位置
)

// ExpandMode 掩码区域形状
type ExpandMode string

const (
	ExpandModeHull     ExpandMode = "hull"     // 多边形（默认）
	ExpandModeBBox     ExpandMode = "bbox"     // 边界框
	ExpandModeOriginal ExpandMode = "original" // 保持原始形状
)

// RefRole 参考图像用途（VACE image_reference）
type RefRole string

const (
	RefRoleObject     RefRole = "obj" // 主体参考
	RefRoleBackground RefRole = "bg"  // 背景参考
)

// AnimateMode 图生动作/视频换人模式
type AnimateMode string

const (
	AnimateModeStd AnimateMode = "wan-std" // 标准模式
	AnimateModePro AnimateMode = "wan-pro" // 专业模式
)

// DigitalHumanStyle 数字人风格
type DigitalHumanStyle string

const (
	DigitalHumanStyleSpeech DigitalHumanStyle = "speech" // 说话风格
)

// -------- 通用请求 --------

type GenerationRequest struct {
	Model      string `json:"model"`
	Input      any    `json:"input,omitempty"`
	Parameters any    `json:"parameters,omitempty"`
}

// -------- Input 类型 --------

// 文生视频: wan2.6-t2v / wan2.5-t2v-preview 等
// 端点: video-generation
type TextToVideoInput struct {
	Prompt         string `json:"prompt,omitempty"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	AudioURL       string `json:"audio_url,omitempty"`
}

// 图生视频（基于首帧）: wan2.6-i2v / wanx2.1-i2v 等
// 端点: video-generation
type ImageToVideoInput struct {
	Prompt         string `json:"prompt,omitempty"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	ImgURL         string `json:"img_url,omitempty"`
	AudioURL       string `json:"audio_url,omitempty"`
	Template       string `json:"template,omitempty"`
}

// 首尾帧生视频: wan2.2-kf2v / wanx2.1-kf2v 等
// 端点: image2video
type FirstLastFrameInput struct {
	Prompt         string `json:"prompt,omitempty"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	FirstFrameURL  string `json:"first_frame_url,omitempty"`
	LastFrameURL   string `json:"last_frame_url,omitempty"`
	Template       string `json:"template,omitempty"`
}

// 参考生视频: wan2.6-r2v / wan2.6-r2v-flash
// 端点: video-generation
type ReferenceVideoInput struct {
	Prompt             string   `json:"prompt,omitempty"`
	NegativePrompt     string   `json:"negative_prompt,omitempty"`
	ReferenceURLs      []string `json:"reference_urls,omitempty"`
	ReferenceVideoURLs []string `json:"reference_video_urls,omitempty"` // 已废弃，兼容保留
}

// VACE 通用视频编辑: wanx2.1-vace-plus
// 端点: video-generation
// function 值: image_reference / video_repainting / video_edit / video_extension / video_outpainting
type VACEInput struct {
	Prompt       string   `json:"prompt,omitempty"`
	Function     VACEFunction `json:"function,omitempty"`
	RefImagesURL []string `json:"ref_images_url,omitempty"`
	VideoURL     string   `json:"video_url,omitempty"`
	ImageURL     string   `json:"image_url,omitempty"`
	// video_edit
	MaskImageURL string `json:"mask_image_url,omitempty"`
	MaskVideoURL string `json:"mask_video_url,omitempty"`
	MaskFrameID  *int   `json:"mask_frame_id,omitempty"`
	// video_extension
	FirstFrameURL string `json:"first_frame_url,omitempty"`
	LastFrameURL  string `json:"last_frame_url,omitempty"`
	FirstClipURL  string `json:"first_clip_url,omitempty"`
	LastClipURL   string `json:"last_clip_url,omitempty"`
}

// 图生动作: wan2.2-animate-move
// 端点: image2video
type AnimateMoveInput struct {
	ImageURL  string `json:"image_url,omitempty"`
	VideoURL  string `json:"video_url,omitempty"`
	Watermark *bool  `json:"watermark,omitempty"`
}

// 视频换人: wan2.2-animate-mix
// 端点: image2video
type AnimateMixInput struct {
	ImageURL  string `json:"image_url,omitempty"`
	VideoURL  string `json:"video_url,omitempty"`
	Watermark *bool  `json:"watermark,omitempty"`
}

// 数字人: wan2.2-s2v
// 端点: image2video
type DigitalHumanInput struct {
	ImageURL string `json:"image_url,omitempty"`
	AudioURL string `json:"audio_url,omitempty"`
}

// -------- Parameters 类型 --------

// GenerationParameters 通用参数（图生/文生/首尾帧/参考生视频）
type GenerationParameters struct {
	Resolution   Resolution `json:"resolution,omitempty"`
	Size         Size       `json:"size,omitempty"`
	Duration     int        `json:"duration,omitempty"`
	PromptExtend *bool      `json:"prompt_extend,omitempty"`
	ShotType     ShotType   `json:"shot_type,omitempty"`
	Audio        *bool      `json:"audio,omitempty"`
	Watermark    *bool      `json:"watermark,omitempty"`
	Seed         *int64     `json:"seed,omitempty"`
}

// AnimateParameters 图生动作 / 视频换人 专用参数
type AnimateParameters struct {
	Mode       AnimateMode `json:"mode,omitempty"`
	CheckImage *bool       `json:"check_image,omitempty"`
}

// DigitalHumanParameters 数字人 wan2.2-s2v 专用参数
type DigitalHumanParameters struct {
	Style DigitalHumanStyle `json:"style,omitempty"`
}

// VACEParameters VACE 通用视频编辑专用参数
type VACEParameters struct {
	// 通用
	Duration     int    `json:"duration,omitempty"`
	PromptExtend *bool  `json:"prompt_extend,omitempty"`
	Watermark    *bool  `json:"watermark,omitempty"`
	Seed         *int64 `json:"seed,omitempty"`
	Size         Size   `json:"size,omitempty"`
	// image_reference
	ObjOrBG []RefRole `json:"obj_or_bg,omitempty"`
	// video_repainting / video_edit / video_extension
	ControlCondition ControlCondition `json:"control_condition,omitempty"`
	Strength         *float64         `json:"strength,omitempty"`
	// video_edit
	MaskType    MaskType   `json:"mask_type,omitempty"`
	ExpandRatio *float64   `json:"expand_ratio,omitempty"`
	ExpandMode  ExpandMode `json:"expand_mode,omitempty"`
	// video_outpainting
	TopScale    *float64 `json:"top_scale,omitempty"`
	BottomScale *float64 `json:"bottom_scale,omitempty"`
	LeftScale   *float64 `json:"left_scale,omitempty"`
	RightScale  *float64 `json:"right_scale,omitempty"`
}

// -------- Detect 同步接口（AnimateAnyone / EMO / LivePortrait / Emoji 等） --------

// DetectRequest 图像检测请求（同步接口，无需异步头）
type DetectRequest struct {
	Model      string      `json:"model"`
	Input      DetectInput `json:"input"`
	Parameters any         `json:"parameters,omitempty"`
}

type DetectInput struct {
	ImageURL string `json:"image_url,omitempty"`
}

// DetectResponse 图像检测响应（不同模型返回字段不同，Data 兜底）
type DetectResponse struct {
	RequestID string          `json:"request_id,omitempty"`
	Output    json.RawMessage `json:"output,omitempty"`
	Usage     json.RawMessage `json:"usage,omitempty"`
	Code      string          `json:"code,omitempty"`
	Message   string          `json:"message,omitempty"`
}

// -------- 任务状态 --------

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusRunning   TaskStatus = "RUNNING"
	TaskStatusSucceeded TaskStatus = "SUCCEEDED"
	TaskStatusFailed    TaskStatus = "FAILED"
	TaskStatusCanceled  TaskStatus = "CANCELED"
	TaskStatusUnknown   TaskStatus = "UNKNOWN"
)

// -------- 创建任务响应 --------

type CreateTaskResponse struct {
	Output    CreateTaskOutput `json:"output"`
	RequestID string           `json:"request_id,omitempty"`
	Code      string           `json:"code,omitempty"`
	Message   string           `json:"message,omitempty"`
}

type CreateTaskOutput struct {
	TaskStatus TaskStatus `json:"task_status"`
	TaskID     string     `json:"task_id"`
}

// -------- 查询任务响应 --------

type TaskResponse struct {
	RequestID string     `json:"request_id,omitempty"`
	Output    TaskOutput `json:"output"`
	Usage     *Usage     `json:"usage,omitempty"`
}

type TaskOutput struct {
	TaskID        string       `json:"task_id,omitempty"`
	TaskStatus    TaskStatus   `json:"task_status,omitempty"`
	SubmitTime    string       `json:"submit_time,omitempty"`
	ScheduledTime string       `json:"scheduled_time,omitempty"`
	EndTime       string       `json:"end_time,omitempty"`
	OrigPrompt    string       `json:"orig_prompt,omitempty"`
	ActualPrompt  string       `json:"actual_prompt,omitempty"`
	VideoURL      string       `json:"video_url,omitempty"`
	Results       *TaskResults `json:"results,omitempty"`
	Code          string       `json:"code,omitempty"`
	Message       string       `json:"message,omitempty"`
}

type TaskResults struct {
	VideoURL string          `json:"video_url,omitempty"`
	Extra    json.RawMessage `json:"-"`
}

// Usage 任务用量统计（不同模型返回的字段组合不同，全部可选）
type Usage struct {
	Duration            float64 `json:"duration,omitempty"`
	InputVideoDuration  int     `json:"input_video_duration,omitempty"`
	OutputVideoDuration int     `json:"output_video_duration,omitempty"`
	VideoCount          int     `json:"video_count,omitempty"`
	VideoDuration       float64 `json:"video_duration,omitempty"`
	VideoRatio          string  `json:"video_ratio,omitempty"`
	SR                  int     `json:"SR,omitempty"`
	Size                string  `json:"size,omitempty"`
	Audio               *bool   `json:"audio,omitempty"`
}

// -------- 辅助方法 --------

func (r TaskResponse) IsDone() bool {
	switch r.Output.TaskStatus {
	case TaskStatusSucceeded, TaskStatusFailed, TaskStatusCanceled, TaskStatusUnknown:
		return true
	default:
		return false
	}
}

func (r TaskResponse) IsSuccess() bool {
	return r.Output.TaskStatus == TaskStatusSucceeded
}

func (r TaskResponse) GetVideoURL() string {
	if r.Output.VideoURL != "" {
		return r.Output.VideoURL
	}
	if r.Output.Results != nil && r.Output.Results.VideoURL != "" {
		return r.Output.Results.VideoURL
	}
	return ""
}
