# DingTalk SDK

一个基于 Go 语言实现的钉钉机器人 SDK，提供完整的钉钉消息发送、接收、流式卡片等功能。

## 功能特性

- ✅ 钉钉机器人消息接收和发送
- ✅ 支持文本、Markdown 消息格式
- ✅ OAuth 2.0 认证和 AccessToken 管理
- ✅ 流式卡片创建和更新
- ✅ 媒体文件（图片、视频、文件）上传
- ✅ 自动 AccessToken 缓存和刷新

## 快速开始

### 安装

```bash
go get github.com/difyz9/dingtalk-sdk.git
```

### 基础使用

#### 1. 创建钉钉客户端

```go
package main

import (
    "github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
    // 创建客户端
    credential := client.Credential{
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
    }
    
    dingClient := client.NewDingTalkClient(credential)
    
    // 获取 Access Token
    token, err := dingClient.GetAccessToken()
    if err != nil {
        panic(err)
    }
    
    println("Access Token:", token)
}
```

#### 2. 发送消息

```go
package main

import (
    "github.com/difyz9/dingtalk-sdk.git/message"
)

func main() {
    // 接收到的消息
    receiveMsg := message.ReceiveMsg{
        SessionWebhook: "your_webhook_url",
        SenderNick:     "用户名",
        SenderStaffId:  "user_id",
    }
    
    // 发送文本消息
    receiveMsg.ReplyToDingtalk(string(message.TEXT), "Hello, DingTalk!")
    
    // 发送 Markdown 消息
    markdownText := "**这是粗体文本**\n\n> 这是引用"
    receiveMsg.ReplyToDingtalk(string(message.MARKDOWN), markdownText)
}
```

#### 3. 上传媒体文件

```go
package main

import (
    "io/ioutil"
    "github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
    credential := client.Credential{
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
    }
    
    dingClient := client.NewDingTalkClient(credential)
    
    // 读取图片文件
    imageData, _ := ioutil.ReadFile("image.png")
    
    // 上传图片
    result, err := dingClient.UploadMedia(
        imageData,
        "image.png",
        client.MediaTypeImage,
        client.MimeTypeImagePng,
    )
    
    if err != nil {
        panic(err)
    }
    
    println("Media ID:", result.MediaID)
}
```

#### 4. 流式卡片更新

```go
package main

import (
    "github.com/difyz9/dingtalk-sdk.git/stream"
    "github.com/google/uuid"
)

func main() {
    // 创建流式卡片客户端
    streamClient, err := stream.NewStreamCardClient()
    if err != nil {
        panic(err)
    }
    
    // 创建并投放卡片
    cardReq := &stream.CreateAndDeliverCardRequest{
        CardTemplateID:   "template_id",
        OutTrackID:       uuid.New().String(),
        OpenSpaceID:      "open_space_id",
        ConversationType: "2", // 群聊
        RobotCode:        "robot_code",
        CardData: map[string]string{
            "content": "初始内容",
        },
    }
    
    err = streamClient.CreateAndDeliverCard("access_token", cardReq)
    if err != nil {
        panic(err)
    }
    
    // 流式更新卡片内容
    updateReq := &stream.StreamingUpdateRequest{
        OutTrackID: cardReq.OutTrackID,
        Key:        "content",
        Content:    "更新后的内容",
        IsFull:     true,
        IsFinalize: true,
    }
    
    err = streamClient.StreamingUpdate("access_token", updateReq)
    if err != nil {
        panic(err)
    }
}
```

## 项目结构

```
dingtalk-sdk/
├── client/         # 钉钉客户端和认证
├── message/        # 消息接收和发送
├── stream/         # 流式卡片功能
├── examples/       # 使用示例
└── README.md
```

## API 文档

### Client 模块

- `NewDingTalkClient(credential Credential) *DingTalkClient` - 创建钉钉客户端
- `GetAccessToken() (string, error)` - 获取 AccessToken（自动缓存）
- `UploadMedia(content []byte, filename, mediaType, mimeType string) (*MediaUploadResult, error)` - 上传媒体文件

### Message 模块

- `ReplyToDingtalk(msgType, msg string) (int, error)` - 回复消息到钉钉
- `GetSenderIdentifier() string` - 获取发送者标识
- `GetChatTitle() string` - 获取聊天标题

### Stream 模块

- `NewStreamCardClient() (*StreamCardClient, error)` - 创建流式卡片客户端
- `CreateAndDeliverCard(accessToken string, req *CreateAndDeliverCardRequest) error` - 创建并投放卡片
- `StreamingUpdate(accessToken string, req *StreamingUpdateRequest) error` - 流式更新卡片

## 许可证

MIT License

## 鸣谢

本项目基于 [chatgpt-dingtalk](https://github.com/eryajf/chatgpt-dingtalk) 项目的钉钉模块改造而成。
