# Stream V2 模式使用指南

## 简介

Stream V2 是钉钉官方推荐的 Stream SDK (`dingtalk-stream-sdk-go`) 使用方式，采用构建器模式 (Builder Pattern) 创建客户端，比传统 Webhook 模式更简单、更可靠。

## 核心优势

| 特性 | Stream V2 模式 | Webhook 模式 |
|------|---------------|--------------|
| 公网要求 | ❌ 不需要 | ✅ 需要公网 IP |
| 配置难度 | 🟢 简单 | 🟡 中等 |
| 实时性 | 🟢 WebSocket 长连接 | 🟡 HTTP 轮询 |
| 自动重连 | ✅ 内置 | ❌ 需自行实现 |
| 安全性 | 🟢 高 | 🟡 中等 |

## 快速开始

### 1. 安装依赖

```bash
go get github.com/open-dingtalk/dingtalk-stream-sdk-go
```

### 2. 基础示例

```go
package main

import (
    "context"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
)

func OnEventReceived(ctx context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
    eventHeader := event.NewEventHeaderFromDataFrame(df)
    println("收到事件:", eventHeader.EventType, df.Data)
    
    // 返回成功响应
    return event.NewSuccessResponse()
    // 返回失败响应(稍后重试)
    // return event.NewLaterResponse()
}

func main() {
    logger.SetLogger(logger.NewStdTestLoggerWithDebug())
    
    cli := client.NewStreamClient(
        client.WithAppCredential(client.NewAppCredentialConfig(
            "your_client_id",
            "your_client_secret",
        )),
    )
    
    cli.RegisterAllEventRouter(OnEventReceived)
    
    err := cli.Start(context.Background())
    if err != nil {
        panic(err)
    }
    
    defer cli.Close()
    select {} // 阻塞主线程
}
```

## 完整功能示例

### 1. 监听机器人消息

```go
import (
    "context"
    "fmt"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
)

func OnChatBotMessage(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
    fmt.Printf("收到消息: %s\n", data.Text.Content)
    
    // 回复消息
    reply := fmt.Sprintf(`{"msgtype":"text","text":{"content":"收到: %s"}}`, data.Text.Content)
    return []byte(reply), nil
}

func main() {
    logger.SetLogger(logger.NewStdTestLoggerWithDebug())
    
    cli := client.NewStreamClient(
        client.WithAppCredential(client.NewAppCredentialConfig(
            "your_client_id",
            "your_client_secret",
        )),
    )
    
    cli.RegisterChatBotCallbackRouter(OnChatBotMessage)
    
    cli.Start(context.Background())
    defer cli.Close()
    select {}
}
```

### 2. 监听特定事件类型

```go
import (
    "context"
    "fmt"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
)

// 自定义事件处理器
func OnSpecificEvent(ctx context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
    eventHeader := event.NewEventHeaderFromDataFrame(df)
    
    fmt.Printf("事件类型: %s\n", eventHeader.EventType)
    fmt.Printf("事件数据: %s\n", df.Data)
    
    // 返回成功响应
    return event.NewSuccessResponse()
}

func main() {
    logger.SetLogger(logger.NewStdTestLoggerWithDebug())
    
    cli := client.NewStreamClient(
        client.WithAppCredential(client.NewAppCredentialConfig(
            "your_client_id",
            "your_client_secret",
        )),
    )
    
    cli.RegisterEventRouter("user_add_org", OnSpecificEvent) // 监听用户加入企业事件
    
    cli.Start(context.Background())
    defer cli.Close()
    select {}
}
```

### 3. 同时监听多种类型

```go
func main() {
    logger.SetLogger(logger.NewStdTestLoggerWithDebug())
    
    cli := client.NewStreamClient(
        client.WithAppCredential(client.NewAppCredentialConfig(
            "your_client_id",
            "your_client_secret",
        )),
    )
    
    // 监听机器人消息
    cli.RegisterChatBotCallbackRouter(OnChatBotMessage)
    
    // 监听所有事件
    cli.RegisterAllEventRouter(OnEventReceived)
    
    // 监听互动卡片回调
    cli.RegisterCardCallbackRouter(OnCardCallback)
    
    // 监听 AI 插件消息
    cli.RegisterPluginCallbackRouter(OnPluginMessage)
    
    cli.Start(context.Background())
    defer cli.Close()
    select {}
}
```

## API 参考

### 核心方法

| 方法 | 说明 | 参数 |
|------|------|------|
| `NewStreamClient(options...)` | 创建 Stream 客户端 | `...ClientOption` |
| `WithAppCredential()` | 设置认证凭证 | `*AppCredentialConfig` |
| `RegisterAllEventRouter()` | 监听所有事件 | `handler.IFrameHandler` |
| `RegisterEventRouter()` | 监听特定事件类型 | `topic string, handler` |
| `RegisterChatBotCallbackRouter()` | 监听机器人消息 | `chatbot.IChatBotMessageHandler` |
| `RegisterCardCallbackRouter()` | 监听互动卡片回调 | `card.ICardCallbackHandler` |
| `RegisterPluginCallbackRouter()` | 监听 AI 插件消息 | `plugin.IPluginMessageHandler` |
| `Start()` | 启动客户端 | `context.Context` |
| `Close()` | 关闭客户端 | - |

### 事件响应

```go
// 处理成功
return event.NewSuccessResponse()

// 处理失败,稍后重试
return event.NewLaterResponse()
```

## 配置选项

### 自定义日志

```go
import "github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"

// 设置日志级别
logger.SetLogger(logger.NewStdTestLoggerWithDebug())
```

### 自动重连

Stream 客户端默认启用自动重连，断线后会自动重新连接。

## 常见问题

### 1. 如何获取 ClientId 和 ClientSecret?

登录 [钉钉开发者平台](https://open-dev.dingtalk.com/) → 应用开发 → 企业内部应用 → 应用信息

### 2. EventStatusSuccess 和 EventStatusLater 的区别?

- `EventStatusSuccess`: 告诉钉钉服务器消息已成功处理
- `EventStatusLater`: 告诉钉钉服务器处理失败，请稍后重试

### 3. 如何在本地测试?

Stream 模式不需要公网 IP，可以直接在本地运行测试。

### 4. 如何回复机器人消息?

在 `RegisterChatBotCallbackHandler` 的回调函数中返回 JSON 格式的消息:

```go
reply := `{"msgtype":"text","text":{"content":"回复内容"}}`
return []byte(reply), nil
```

## 完整示例

查看项目中的完整示例代码：

- [examples/stream_v2/main.go](../examples/stream_v2/main.go) - Stream V2 基础示例
- [官方示例](https://github.com/open-dingtalk/dingtalk-stream-sdk-go/tree/main/example) - 官方完整示例

## 参考资料

- [官方 Stream SDK](https://github.com/open-dingtalk/dingtalk-stream-sdk-go)
- [Stream 模式文档](https://opensource.dingtalk.com/developerpedia/docs/learn/stream/overview)
- [钉钉开放平台](https://open.dingtalk.com/)
