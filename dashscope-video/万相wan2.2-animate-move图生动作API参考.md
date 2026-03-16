万相wan2.2-animate-move图生动作API参考-大模型服务平台百炼(Model Studio)-阿里云帮助中心

[官方文档](https://help.aliyun.com/)

输入文档关键字查找

万相-图生动作模型，可基于人物图片和参考视频，生成人物动作视频。

功能简述：将视频角色的动作/表情迁移到图片角色中，赋予图片角色动态表现力。

适用场景：复刻舞蹈、复刻高难度肢体动作、复刻影视剧表演表情及肢体动作细节，低成本动捕替代。

## 模型效果

万相-图生动作模型wan2.2-animate-move提供标准模式`wan-std` 和专业模式`wan-pro` 两种服务模式，不同模式在效果和计费上存在差异，请参见 [模型调用计费](https://help.aliyun.com/zh/model-studio/model-pricing#1f708fcd62rdc)。

人物图片

参考视频

输出视频（标准模式`wan-std`）

输出视频（专业模式`wan-pro`）

## 前提条件

您需要已 [获取API Key](https://help.aliyun.com/zh/model-studio/get-api-key) 并 [配置API Key到环境变量](https://help.aliyun.com/zh/model-studio/configure-api-key-through-environment-variables)。

重要

北京和新加坡地域拥有独立的 API Key 与请求地址，不可混用，跨地域调用将导致鉴权失败或服务报错。

## HTTP调用

由于视频生成任务耗时较长，API采用异步调用。整个流程包含 “创建任务 -> 轮询获取” 两个核心步骤，具体如下：

### 步骤1：创建任务获取任务ID

北京地域：`POST https://dashscope.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis`

新加坡地域：`POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis`

说明

创建成功后，使用接口返回的`task_id` 查询结果，task_id 有效期为 24 小时。请勿重复创建任务，轮询获取即可。

新手指引请参见 [Postman](https://help.aliyun.com/zh/model-studio/first-call-to-image-and-video-api)。

#### 请求参数

## 图生动作

以下为北京地域 base_url，若使用新加坡地域的模型，需将 base_url 替换为：`https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis`

```curl
curl --location 'https://dashscope.aliyuncs.com/api/v1/services/aigc/image2video/video-synthesis' --header 'X-DashScope-Async: enable' --header "Authorization: Bearer $DASHSCOPE_API_KEY" --header 'Content-Type: application/json' --data '{
    "model": "wan2.2-animate-move",
    "input": {
        "image_url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250919/adsyrp/move_input_image.jpeg",
        "video_url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250919/kaakcn/move_input_video.mp4",
        "watermark": true
    },
    "parameters": {
        "mode": "wan-std"
    }
  }'
```

[内容已按页面 Markdown 落盘，正文同抓取结果一致]

联系我们：4008013260
