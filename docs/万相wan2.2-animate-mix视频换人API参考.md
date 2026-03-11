万相wan2.2-animate-mix视频换人API参考-大模型服务平台百炼(Model Studio)-阿里云帮助中心

[官方文档](https://help.aliyun.com/)

输入文档关键字查找

万相-视频换人模型能够依据人物图片和参考视频，将视频中的主角替换为图片中的角色，同时保留原视频的场景、光照和色调，实现无缝换人。

核心功能： 在不改变原始视频的动作、表情及环境的条件下，将视频中的角色替换为指定图片中的人物。

适用场景：适用于二次创作、影视后期制作等需要进行角色替换的场景。

## 效果示例

万相-视频换人模型wan2.2-animate-mix提供标准模式`wan-std` 和专业模式`wan-pro` 两种服务模式，不同模式在效果和计费上存在差异，详情参见计费与限流。

## HTTP调用

您需要已 [获取API Key](https://help.aliyun.com/zh/model-studio/get-api-key) 并 [配置API Key到环境变量](https://help.aliyun.com/zh/model-studio/configure-api-key-through-environment-variables)。

### 步骤1：创建任务获取任务ID

北京地域：`POST https://dashscope.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis`

新加坡地域：`POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis`

```curl
curl --location 'https://dashscope.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis' --header 'X-DashScope-Async: enable' --header "Authorization: Bearer $DASHSCOPE_API_KEY" --header 'Content-Type: application/json' --data '{
    "model": "wan2.2-animate-mix",
    "input": {
        "image_url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250919/bhkfor/mix_input_image.jpeg",
        "video_url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250919/wqefue/mix_input_video.mp4",
        "watermark": true
    },
    "parameters": {
        "mode": "wan-std"
    }
  }'
```

[内容已按页面 Markdown 落盘，正文同抓取结果一致]

联系我们：4008013260
