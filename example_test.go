package dashscope_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ai-web4/dashscope-go"
)

// 首尾帧生视频
func ExampleClient_SubmitAndWaitImageToVideo() {
	client, _ := dashscope.NewClient(dashscope.Config{
		APIKey: "sk-xxx",
		Region: dashscope.RegionBeijing,
	})

	resp, err := client.SubmitAndWaitImageToVideo(context.Background(), dashscope.GenerationRequest{
		Model: "wan2.2-kf2v-flash",
		Input: dashscope.FirstLastFrameInput{
			Prompt:        "一只小猫从窗边跳到沙发上",
			FirstFrameURL: "https://example.com/first.png",
			LastFrameURL:  "https://example.com/last.png",
		},
		Parameters: dashscope.GenerationParameters{
			Resolution: "720P",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetVideoURL())
}

// 文生视频
func ExampleClient_SubmitAndWaitVideoGeneration_textToVideo() {
	client, _ := dashscope.NewClient(dashscope.Config{
		APIKey: "sk-xxx",
		Region: dashscope.RegionBeijing,
	})

	resp, err := client.SubmitAndWaitVideoGeneration(context.Background(), dashscope.GenerationRequest{
		Model: "wan2.6-t2v",
		Input: dashscope.TextToVideoInput{
			Prompt: "夕阳下，一只金毛犬在海滩上奔跑",
		},
		Parameters: dashscope.GenerationParameters{
			Size:     "1280*720",
			Duration: 5,
			ShotType: "single",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetVideoURL())
}

// 参考生视频
func ExampleClient_SubmitAndWaitVideoGeneration_referenceVideo() {
	client, _ := dashscope.NewClient(dashscope.Config{
		APIKey: "sk-xxx",
		Region: dashscope.RegionBeijing,
	})

	resp, err := client.SubmitAndWaitVideoGeneration(context.Background(), dashscope.GenerationRequest{
		Model: "wan2.6-r2v-flash",
		Input: dashscope.ReferenceVideoInput{
			Prompt: "character1在沙发上开心地看电影",
			ReferenceURLs: []string{
				"https://example.com/character.mp4",
			},
		},
		Parameters: dashscope.GenerationParameters{
			Size:     "1280*720",
			Duration: 5,
			ShotType: "single",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetVideoURL())
}

// VACE 视频重绘
func ExampleClient_SubmitAndWaitVideoGeneration_vaceRepainting() {
	client, _ := dashscope.NewClient(dashscope.Config{
		APIKey: "sk-xxx",
		Region: dashscope.RegionBeijing,
	})

	strength := 0.8
	resp, err := client.SubmitAndWaitVideoGeneration(context.Background(), dashscope.GenerationRequest{
		Model: "wanx2.1-vace-plus",
		Input: dashscope.VACEInput{
			Prompt:   "一位穿红色裙子的女孩在跳舞",
			Function: "video_repainting",
			VideoURL: "https://example.com/dance.mp4",
		},
		Parameters: dashscope.VACEParameters{
			ControlCondition: "posebodyface",
			Strength:         &strength,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetVideoURL())
}

// 图生动作
func ExampleClient_SubmitAndWaitImageToVideo_animateMove() {
	client, _ := dashscope.NewClient(dashscope.Config{
		APIKey: "sk-xxx",
		Region: dashscope.RegionBeijing,
	})

	resp, err := client.SubmitAndWaitImageToVideo(context.Background(), dashscope.GenerationRequest{
		Model: "wan2.2-animate-move",
		Input: dashscope.AnimateMoveInput{
			ImageURL: "https://example.com/person.jpg",
			VideoURL: "https://example.com/dance_ref.mp4",
		},
		Parameters: dashscope.AnimateParameters{
			Mode: "wan-std",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetVideoURL())
}

// 同步 detect 接口调用示例
func ExampleClient_CallSync_detect() {
	client, _ := dashscope.NewClient(dashscope.Config{
		APIKey: "sk-xxx",
		Region: dashscope.RegionBeijing,
	})

	var resp dashscope.DetectResponse
	err := client.CallSync(context.Background(),
		"/api/v1/services/aigc/image2video/video-synthesis",
		dashscope.DetectRequest{
			Model: "emo-detect-v1",
			Input: dashscope.DetectInput{
				ImageURL: "https://example.com/portrait.jpg",
			},
		},
		&resp,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("detect output: %s\n", string(resp.Output))
}

// 手动提交 + 自定义轮询间隔
func ExampleClient_WaitTaskWithInterval() {
	client, _ := dashscope.NewClient(dashscope.Config{
		APIKey: "sk-xxx",
		Region: dashscope.RegionBeijing,
	})

	created, err := client.SubmitVideoGeneration(context.Background(), dashscope.GenerationRequest{
		Model: "wan2.6-t2v",
		Input: dashscope.TextToVideoInput{
			Prompt: "星空下的湖面倒映着月光",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.WaitTaskWithInterval(context.Background(), created.Output.TaskID, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetVideoURL())
}
