# 钉钉消息发送使用指南

## 概述

本 SDK 提供两种消息发送方式：

1. **SessionWebhook 回复** - 用于响应用户消息
2. **主动推送消息** - 用于主动向群/用户发送消息

## 方式一：SessionWebhook 回复消息（推荐）

这是最常用的方式，适用于机器人接收到用户消息后的回复场景。

### 特点
- ✅ 无需额外权限
- ✅ 实现简单
- ✅ 适合对话式交互
- ❌ 只能在接收到消息后的20分钟内使用

### 代码示例

```go
package main

import (
    "fmt"
    "github.com/difyz9/dingtalk-sdk.git/message"
)

func main() {
    // 模拟从钉钉回调接收到的消息
    receiveMsg := message.ReceiveMsg{
        // SessionWebhook 由钉钉回调提供，20分钟内有效
        SessionWebhook:    "https://oapi.dingtalk.com/robot/sendBySession?session=xxx",
        SenderNick:        "张三",
        SenderStaffId:     "user123",
        ConversationType:  "2", // 1: 私聊, 2: 群聊
        ConversationTitle: "技术交流群",
        Text: message.Text{
            Content: "你好，机器人",
        },
    }

    // 1. 发送文本消息
    statusCode, err := receiveMsg.ReplyToDingtalk(
        string(message.TEXT),
        "你好！我是钉钉机器人 🤖\n\n我可以帮你做什么？",
    )
    if err != nil {
        fmt.Printf("发送失败: %v\n", err)
    } else {
        fmt.Printf("✅ 消息发送成功，状态码: %d\n", statusCode)
    }

    // 2. 发送 Markdown 消息
    markdownText := `### 📋 功能菜单

**我可以帮您：**

1. 📊 查询数据报表
2. 🔔 接收系统通知
3. 💬 智能问答
4. 🤖 自动化任务

> 发送对应序号即可使用相应功能
`
    
    receiveMsg.ReplyToDingtalk(string(message.MARKDOWN), markdownText)
}
```

### 获取 SessionWebhook

SessionWebhook 来自钉钉的事件回调，在你的 HTTP 服务器中接收：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/difyz9/dingtalk-sdk.git/message"
)

func main() {
    router := gin.Default()
    
    // 钉钉回调接口
    router.POST("/webhook", func(c *gin.Context) {
        var msg message.ReceiveMsg
        if err := c.BindJSON(&msg); err != nil {
            return
        }
        
        // msg.SessionWebhook 就是可以用来回复的 webhook
        // msg.Text.Content 是用户发送的消息内容
        
        // 处理消息并回复
        msg.ReplyToDingtalk(string(message.TEXT), "收到您的消息："+msg.Text.Content)
    })
    
    router.Run(":8080")
}
```

## 方式二：主动推送群消息

这种方式可以主动向群发送消息，无需等待用户触发。

### 特点
- ✅ 可以主动推送
- ✅ 适合定时任务、告警通知等
- ❌ 需要申请 `qyapi_chat_manage` 权限
- ❌ 需要知道群的 chatId

### 申请权限

1. 访问钉钉开放平台：https://open-dev.dingtalk.com/
2. 进入你的应用
3. 点击"权限管理"
4. 搜索并申请 `qyapi_chat_manage` 权限
5. 或者直接访问提示中的链接申请

### 代码示例

```go
package main

import (
    "fmt"
    "github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
    // 创建客户端
    credential := client.Credential{
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
    }
    
    dingClient := client.NewDingTalkClient(credential)
    
    // 获取群的 chatId（方法见下文）
    chatID := "chatxxxxxxxxxxxxxx"
    
    // 发送文本消息
    textMsg := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": "【系统通知】\n服务器负载过高，请注意！",
        },
    }
    
    err := dingClient.SendRobotMessage(chatID, textMsg)
    if err != nil {
        fmt.Printf("发送失败: %v\n", err)
    } else {
        fmt.Println("✅ 消息发送成功")
    }
    
    // 发送 Markdown 消息
    markdownMsg := map[string]interface{}{
        "msgtype": "markdown",
        "markdown": map[string]string{
            "title": "告警通知",
            "text": `### ⚠️ 服务器告警

**告警级别**: 🔴 严重

**告警时间**: 2026-02-07 20:00:00

**问题描述**:
- CPU 使用率: 95%
- 内存使用率: 88%
- 磁盘使用率: 92%

**建议操作**:
1. 检查异常进程
2. 清理临时文件
3. 扩容服务器资源

> 请立即处理！
`,
        },
    }
    
    dingClient.SendRobotMessage(chatID, markdownMsg)
}
```

## 获取群的 chatId

### 方法1: 通过机器人命令（推荐）

在群里 @ 机器人并发送"群ID"，机器人会回复群的 chatId（需要实现这个功能）

### 方法2: 通过 API 获取

```go
// TODO: 实现获取群列表的 API
// 参考文档: https://open.dingtalk.com/document/orgapp/create-a-group-session
```

### 方法3: 从回调中获取

当用户在群里 @ 机器人时，回调数据中的 `ConversationID` 就是 chatId：

```go
router.POST("/webhook", func(c *gin.Context) {
    var msg message.ReceiveMsg
    c.BindJSON(&msg)
    
    // msg.ConversationID 就是群的 chatId
    fmt.Println("群 ID:", msg.ConversationID)
})
```

## 消息类型对比

| 类型 | 用途 | 优点 | 缺点 |
|------|------|------|------|
| SessionWebhook | 回复用户消息 | 简单、无需额外权限 | 20分钟时效、被动触发 |
| 主动推送 | 定时通知、告警 | 可主动发送、无时效限制 | 需要权限、需要chatId |

## 完整示例：HTTP 服务器

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/difyz9/dingtalk-sdk.git/client"
    "github.com/difyz9/dingtalk-sdk.git/message"
)

var dingClient *client.DingTalkClient

func init() {
    credential := client.Credential{
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
    }
    dingClient = client.NewDingTalkClient(credential)
}

func main() {
    router := gin.Default()
    
    // 接收钉钉回调
    router.POST("/webhook", handleDingTalkCallback)
    
    // 主动推送消息的 API
    router.POST("/send", handleSendMessage)
    
    router.Run(":8080")
}

// 处理钉钉回调
func handleDingTalkCallback(c *gin.Context) {
    var msg message.ReceiveMsg
    if err := c.BindJSON(&msg); err != nil {
        return
    }
    
    // 记录群 ID
    if msg.ConversationType == "2" {
        fmt.Printf("群ID: %s, 群名: %s\n", msg.ConversationID, msg.ConversationTitle)
    }
    
    // 根据用户消息回复
    switch msg.Text.Content {
    case "帮助", "help":
        helpText := `### 🤖 机器人使用指南

**命令列表**:
- 帮助/help - 显示本帮助信息
- 群ID - 获取当前群的ID
- 状态 - 查询系统状态

> 更多功能开发中...`
        msg.ReplyToDingtalk(string(message.MARKDOWN), helpText)
        
    case "群ID":
        msg.ReplyToDingtalk(
            string(message.TEXT),
            fmt.Sprintf("当前群ID: %s", msg.ConversationID),
        )
        
    case "状态":
        msg.ReplyToDingtalk(
            string(message.TEXT),
            "✅ 系统运行正常\n🟢 所有服务状态良好",
        )
        
    default:
        msg.ReplyToDingtalk(
            string(message.TEXT),
            "收到消息："+msg.Text.Content+"\n\n发送 '帮助' 查看可用命令",
        )
    }
}

// 主动发送消息（需要权限）
func handleSendMessage(c *gin.Context) {
    type Request struct {
        ChatID  string `json:"chatId"`
        Message string `json:"message"`
    }
    
    var req Request
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    textMsg := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": req.Message,
        },
    }
    
    err := dingClient.SendRobotMessage(req.ChatID, textMsg)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"status": "success"})
}
```

## 测试主动发送

```bash
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "chatId": "chatxxxxxx",
    "message": "测试消息"
  }'
```

## 常见问题

### Q: SessionWebhook 过期了怎么办？
A: SessionWebhook 有20分钟时效，过期后无法使用。需要等待用户再次发送消息获取新的 webhook。

### Q: 如何获取 chatId？
A: 
1. 让用户在群里 @ 机器人，从回调中的 `ConversationID` 获取
2. 实现"群ID"命令，让机器人回复 chatId
3. 使用钉钉 API 获取群列表

### Q: 权限申请需要多久？
A: 通常几分钟到几小时不等，具体看钉钉审核速度。

### Q: 可以发送图片吗？
A: 可以，先使用 `UploadMedia` 上传图片获得 media_id，然后发送图片消息。

## 更多资源

- [钉钉开放平台文档](https://open.dingtalk.com/document/)
- [机器人开发指南](https://open.dingtalk.com/document/orgapp/robot-overview)
- [消息类型说明](https://open.dingtalk.com/document/orgapp/message-types-and-data-format)
