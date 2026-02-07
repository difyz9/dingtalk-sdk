# 主动发送消息指南

本指南介绍如何使用钉钉 SDK 主动向用户发送消息。

## 目录

- [概述](#概述)
- [方式 1: 企业内部机器人](#方式-1-企业内部机器人)
- [方式 2: Webhook 自定义机器人](#方式-2-webhook-自定义机器人)
- [获取 ChatID](#获取-chatid)
- [消息类型](#消息类型)
- [最佳实践](#最佳实践)

## 概述

钉钉支持两种主动发送消息的方式：

| 方式 | 优势 | 劣势 | 适用场景 |
|------|------|------|----------|
| 企业内部机器人 | 功能强大，支持群聊/单聊 | 需要 OAuth 认证 | 企业内部系统通知 |
| Webhook 自定义机器人 | 配置简单，无需认证 | 仅支持群聊 | 简单的群通知 |

## 方式 1: 企业内部机器人

### 1.1 准备工作

1. 登录 [钉钉开发者平台](https://open-dev.dingtalk.com/)
2. 创建企业内部应用
3. 获取 `ClientID` 和 `ClientSecret`
4. 配置机器人权限：`企业内部群消息机器人`

### 1.2 基础示例

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
    
    // 准备消息
    chatID := "your_chat_id" // 群聊或单聊的 ID
    
    message := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": "你好！这是一条主动发送的消息",
        },
    }
    
    // 发送消息
    err := dingClient.SendRobotMessage(chatID, message)
    if err != nil {
        panic(err)
    }
}
```

### 1.3 支持的消息类型

#### 文本消息

```go
textMsg := map[string]interface{}{
    "msgtype": "text",
    "text": map[string]string{
        "content": "消息内容",
    },
    // 可选: @指定用户
    "at": map[string]interface{}{
        "atUserIds": []string{"user_id_1", "user_id_2"},
        "isAtAll":   false, // true 表示 @所有人
    },
}

dingClient.SendRobotMessage(chatID, textMsg)
```

#### Markdown 消息

```go
markdownMsg := map[string]interface{}{
    "msgtype": "markdown",
    "markdown": map[string]string{
        "title": "消息标题",
        "text": `### 标题

**粗体文本**

- 列表项 1
- 列表项 2

> 引用内容`,
    },
}

dingClient.SendRobotMessage(chatID, markdownMsg)
```

#### 链接消息

```go
linkMsg := map[string]interface{}{
    "msgtype": "link",
    "link": map[string]string{
        "title":      "链接标题",
        "text":       "链接描述",
        "messageUrl": "https://www.example.com",
        "picUrl":     "https://example.com/image.png",
    },
}

dingClient.SendRobotMessage(chatID, linkMsg)
```

#### ActionCard 消息

```go
actionCardMsg := map[string]interface{}{
    "msgtype": "actionCard",
    "actionCard": map[string]interface{}{
        "title": "卡片标题",
        "text":  "卡片内容",
        "singleTitle": "按钮文字",
        "singleURL":   "https://www.example.com",
    },
}

dingClient.SendRobotMessage(chatID, actionCardMsg)
```

## 方式 2: Webhook 自定义机器人

### 2.1 准备工作

1. 进入钉钉群聊
2. 群设置 → 智能群助手 → 添加机器人 → 自定义
3. 配置机器人名称、头像等
4. 获取 Webhook URL

### 2.2 基础示例

```go
package main

import (
    "github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
    // 创建客户端（即使使用 Webhook，仍需要客户端实例）
    credential := client.Credential{
        ClientID:     "any",
        ClientSecret: "any",
    }
    dingClient := client.NewDingTalkClient(credential)
    
    // Webhook URL
    webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=xxx"
    
    // 准备消息
    message := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]interface{}{
            "content": "Webhook 发送的消息",
        },
    }
    
    // 发送消息
    err := dingClient.SendWebhookMessage(webhookURL, message)
    if err != nil {
        panic(err)
    }
}
```

### 2.3 @指定用户

```go
// @指定用户（需要用户的手机号）
message := map[string]interface{}{
    "msgtype": "text",
    "text": map[string]interface{}{
        "content": "@张三 请查看",
    },
    "at": map[string]interface{}{
        "atMobiles": []string{"13800138000"},
        "isAtAll":   false,
    },
}

dingClient.SendWebhookMessage(webhookURL, message)
```

## 获取 ChatID

ChatID 是发送消息的关键参数，有以下几种获取方式：

### 方法 1: 从消息回调中获取

当用户给机器人发消息时，钉钉会回调你的服务器，回调数据中包含 `conversationId`：

```go
type ReceiveMsg struct {
    ConversationId   string `json:"conversationId"`   // 这就是 chatId
    ChatbotUserId    string `json:"chatbotUserId"`
    SenderStaffId    string `json:"senderStaffId"`
    // ... 其他字段
}
```

### 方法 2: 使用 GetOpenConversationId 转换

如果你有原始的 chatId，可以转换为 OpenConversationId：

```go
// chatId: 群聊的 chatId（从其他渠道获得）
openConversationId, err := dingClient.GetOpenConversationId("chatId")
if err != nil {
    panic(err)
}

// 使用转换后的 ID 发送消息
message := map[string]interface{}{
    "msgtype": "text",
    "text": map[string]string{
        "content": "消息内容",
    },
}
dingClient.SendRobotMessage(openConversationId, message)
```

参考示例: [examples/get_chat_list/main.go](../examples/get_chat_list/main.go)

### 方法 3: 固定群聊

对于固定的群聊，可以：

1. 让机器人加入群聊
2. 在群内发一条消息
3. 从回调中获取并保存 `conversationId`
4. 后续使用保存的 ID 发送消息

## 消息类型

### 文本消息

```go
map[string]interface{}{
    "msgtype": "text",
    "text": map[string]string{
        "content": "纯文本内容",
    },
}
```

### Markdown 消息

支持标准 Markdown 语法：

```go
map[string]interface{}{
    "msgtype": "markdown",
    "markdown": map[string]string{
        "title": "标题",
        "text": `
### 一级标题
#### 二级标题

**粗体** *斜体* 

- 无序列表
1. 有序列表

> 引用

[链接](https://www.dingtalk.com)

![图片](https://example.com/image.png)
`,
    },
}
```

### 链接消息

```go
map[string]interface{}{
    "msgtype": "link",
    "link": map[string]string{
        "title":      "链接标题",
        "text":       "链接描述文字",
        "messageUrl": "https://www.example.com",
        "picUrl":     "https://example.com/image.png", // 可选
    },
}
```

### ActionCard 消息

单按钮：

```go
map[string]interface{}{
    "msgtype": "actionCard",
    "actionCard": map[string]interface{}{
        "title": "卡片标题",
        "text":  "卡片内容（支持 Markdown）",
        "singleTitle": "查看详情",
        "singleURL":   "https://www.example.com",
    },
}
```

多按钮：

```go
map[string]interface{}{
    "msgtype": "actionCard",
    "actionCard": map[string]interface{}{
        "title": "卡片标题",
        "text":  "卡片内容",
        "btns": []map[string]string{
            {
                "title":     "同意",
                "actionURL": "https://www.example.com/approve",
            },
            {
                "title":     "拒绝",
                "actionURL": "https://www.example.com/reject",
            },
        },
    },
}
```

## 最佳实践

### 1. 合理选择发送方式

```go
// 企业内部通知、需要精确控制 -> 使用企业内部机器人
func sendToEmployee(dingClient *client.DingTalkClient, chatID string) {
    message := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": "您有新的审批待处理",
        },
    }
    dingClient.SendRobotMessage(chatID, message)
}

// 简单的群通知 -> 使用 Webhook
func sendGroupNotification(dingClient *client.DingTalkClient, webhookURL string) {
    message := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]interface{}{
            "content": "系统维护通知",
        },
    }
    dingClient.SendWebhookMessage(webhookURL, message)
}
```

### 2. 错误处理

```go
err := dingClient.SendRobotMessage(chatID, message)
if err != nil {
    log.Printf("发送消息失败: %v", err)
    // 可以实现重试机制
    retryCount := 3
    for i := 0; i < retryCount; i++ {
        time.Sleep(time.Second * 2)
        err = dingClient.SendRobotMessage(chatID, message)
        if err == nil {
            break
        }
    }
}
```

### 3. 批量发送

```go
// 向多个群发送相同消息
func sendToMultipleChats(dingClient *client.DingTalkClient, chatIDs []string, message map[string]interface{}) {
    for _, chatID := range chatIDs {
        err := dingClient.SendRobotMessage(chatID, message)
        if err != nil {
            log.Printf("发送到 %s 失败: %v", chatID, err)
            continue
        }
        
        // 避免频繁发送，加入延迟
        time.Sleep(time.Millisecond * 100)
    }
}
```

### 4. 定时任务发送

```go
import (
    "time"
    "github.com/robfig/cron/v3"
)

func setupDailyReport(dingClient *client.DingTalkClient, chatID string) {
    c := cron.New()
    
    // 每天 9:00 发送日报
    c.AddFunc("0 9 * * *", func() {
        message := map[string]interface{}{
            "msgtype": "markdown",
            "markdown": map[string]string{
                "title": "每日数据报告",
                "text": generateDailyReport(),
            },
        }
        
        dingClient.SendRobotMessage(chatID, message)
    })
    
    c.Start()
}

func generateDailyReport() string {
    return fmt.Sprintf(`### 📊 每日数据报告

**日期**: %s

| 指标 | 数值 |
|------|------|
| 新增用户 | 1,234 |
| 活跃用户 | 5,678 |

> ✅ 所有指标正常`, time.Now().Format("2006-01-02"))
}
```

### 5. 消息模板化

```go
type MessageTemplate struct {
    Type    string
    Title   string
    Content string
}

func (t *MessageTemplate) Build() map[string]interface{} {
    switch t.Type {
    case "alert":
        return map[string]interface{}{
            "msgtype": "markdown",
            "markdown": map[string]string{
                "title": "⚠️ " + t.Title,
                "text": fmt.Sprintf(`### ⚠️ %s

%s

**时间**: %s`, t.Title, t.Content, time.Now().Format("15:04:05")),
            },
        }
    case "success":
        return map[string]interface{}{
            "msgtype": "markdown",
            "markdown": map[string]string{
                "title": "✅ " + t.Title,
                "text": fmt.Sprintf(`### ✅ %s

%s`, t.Title, t.Content),
            },
        }
    default:
        return map[string]interface{}{
            "msgtype": "text",
            "text": map[string]string{
                "content": t.Content,
            },
        }
    }
}

// 使用示例
func sendAlert(dingClient *client.DingTalkClient, chatID string) {
    template := MessageTemplate{
        Type:    "alert",
        Title:   "系统告警",
        Content: "服务器 CPU 使用率超过 80%",
    }
    
    dingClient.SendRobotMessage(chatID, template.Build())
}
```

## 完整示例

查看完整的可运行示例：

- [examples/active_send/main.go](../examples/active_send/main.go) - 主动发送消息完整示例
- [examples/send_message/main.go](../examples/send_message/main.go) - 双模式发送示例
- [examples/webhook/main.go](../examples/webhook/main.go) - Webhook 发送示例

## 常见问题

### 1. 如何获取 chatId？

最简单的方法是让用户先给机器人发一条消息，从回调中获取 `conversationId`。

### 2. 为什么发送失败？

常见原因：
- chatId 不正确
- 机器人未加入群聊
- 权限配置不正确
- Access Token 过期

### 3. 可以发送给单个用户吗？

可以。使用企业内部机器人方式，chatId 设置为用户的单聊 ID。

### 4. 发送频率有限制吗？

有。钉钉对消息发送有频率限制，建议：
- 同一群聊：间隔 > 100ms
- 批量发送：加入适当延迟

### 5. Webhook 和企业内部机器人可以混用吗？

可以。根据不同场景选择合适的方式。

## 参考资料

- [钉钉开放平台](https://open.dingtalk.com/)
- [机器人消息类型](https://open.dingtalk.com/document/orgapp/message-types-and-data-format)
- [企业内部机器人](https://open.dingtalk.com/document/orgapp/enterprise-internal-robot)
