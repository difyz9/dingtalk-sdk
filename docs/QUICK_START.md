# 快速开始 - 5分钟实现钉钉消息发送

## 目标

5分钟内实现一个能接收和回复钉钉消息的机器人。

## 前置条件

1. ✅ Go 1.19 或更高版本
2. ✅ 钉钉开发者账号
3. ✅ 已创建的钉钉企业内部应用或机器人

## 步骤

### 1. 获取钉钉凭证（2分钟）

#### 方式 A: 企业内部应用
1. 登录 [钉钉开发者平台](https://open-dev.dingtalk.com/)
2. 创建应用 → 企业内部应用
3. 记录 **AppKey** (ClientID) 和 **AppSecret** (ClientSecret)

#### 方式 B: 机器人
1. 钉钉群 → 群设置 → 智能群助手 → 添加机器人
2. 选择"自定义"机器人
3. 配置 HTTP 回调地址（如 `https://yourdomain.com/webhook`）
4. 记录 Webhook 地址

### 2. 安装 SDK（30秒）

```bash
go get github.com/difyz9/dingtalk-sdk.git
```

### 3. 编写代码（2分钟）

创建 `main.go`:

```go
package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/difyz9/dingtalk-sdk.git/client"
    "github.com/difyz9/dingtalk-sdk.git/message"
)

func main() {
    // 1. 创建钉钉客户端
    credential := client.Credential{
        ClientID:     "dingxxxxxx",        // 替换为你的 AppKey
        ClientSecret: "your_app_secret",   // 替换为你的 AppSecret
    }
    dingClient := client.NewDingTalkClient(credential)
    
    // 2. 验证 Token（可选）
    token, _ := dingClient.GetAccessToken()
    fmt.Println("✅ Access Token:", token)
    
    // 3. 启动 HTTP 服务接收钉钉回调
    router := gin.Default()
    
    router.POST("/webhook", func(c *gin.Context) {
        var msg message.ReceiveMsg
        if err := c.BindJSON(&msg); err != nil {
            return
        }
        
        fmt.Printf("收到消息: %s\n", msg.Text.Content)
        
        // 4. 回复消息
        switch msg.Text.Content {
        case "hi", "hello", "你好":
            msg.ReplyToDingtalk(string(message.TEXT), "你好！我是钉钉机器人 🤖")
            
        case "help", "帮助":
            helpText := `### 🤖 机器人帮助

**可用命令**:
- hi/hello/你好 - 打招呼
- help/帮助 - 显示此帮助
- time - 获取当前时间
- status - 查看系统状态`
            msg.ReplyToDingtalk(string(message.MARKDOWN), helpText)
            
        default:
            msg.ReplyToDingtalk(string(message.TEXT), "收到："+msg.Text.Content)
        }
    })
    
    fmt.Println("🚀 服务器启动在 :8080")
    router.Run(":8080")
}
```

### 4. 运行测试（30秒）

```bash
go run main.go
```

---

## Stream 模式接入（推荐）

使用官方 Stream SDK (`dingtalk-stream-sdk-go`) 可以更简单地接收钉钉消息和事件，无需配置公网 Webhook 地址。

### 1. 安装 Stream SDK

```bash
go get github.com/open-dingtalk/dingtalk-stream-sdk-go
```

### 2. 使用 Builder 模式创建客户端

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
)

func OnEventReceived(ctx context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
    eventHeader := event.NewEventHeaderFromDataFrame(df)
    fmt.Printf("收到事件: %s\n", eventHeader.EventType)
    return event.NewSuccessResponse()
}

func main() {
    // 配置日志
    logger.SetLogger(logger.NewStdTestLoggerWithDebug())
    
    // 创建客户端
    cli := client.NewStreamClient(
        client.WithAppCredential(client.NewAppCredentialConfig(
            "your_client_id",
            "your_client_secret",
        )),
    )
    
    // 监听所有事件
    cli.RegisterAllEventRouter(OnEventReceived)
    
    // 启动客户端
    err := cli.Start(context.Background())
    if err != nil {
        panic(err)
    }
    
    defer cli.Close()
    fmt.Println("✅ Stream 客户端已启动")
    select {} // 阻塞主线程
}
```

### 3. Stream 模式优势

✅ **无需公网 IP** - 不需要配置 Webhook 回调地址  
✅ **自动重连** - 内置断线重连机制  
✅ **实时推送** - WebSocket 长连接，低延迟  
✅ **更安全** - 不暴露公网端点  

### 4. 完整示例

参考项目中的 [examples/stream_v2/main.go](../examples/stream_v2/main.go)

---

详细文档请查看: [Stream V2 使用指南](STREAM_V2_GUIDE.md)

输出:
```
✅ Access Token: 605e241440c43d8f924417e64fc25fb2
🚀 服务器启动在 :8080
```

### 5. 配置钉钉回调（1分钟）

#### 本地开发（使用 ngrok）
```bash
# 安装 ngrok
brew install ngrok  # macOS
# 或从 https://ngrok.com/ 下载

# 启动隧道
ngrok http 8080
```

复制 ngrok 提供的 HTTPS 地址（如 `https://abc123.ngrok.io`）

#### 配置机器人
1. 钉钉开放平台 → 应用开发 → 机器人配置
2. 设置 **消息接收地址**: `https://abc123.ngrok.io/webhook`
3. 保存配置

### 6. 测试（30秒）

1. 在钉钉群里 @ 你的机器人
2. 发送 "hello"
3. 机器人应该回复 "你好！我是钉钉机器人 🤖"

## 成功！🎉

你的第一个钉钉机器人已经运行起来了！

## 下一步

### 添加更多功能

#### 1. 定时推送消息

```go
// 每小时推送一次
go func() {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        msg := map[string]interface{}{
            "msgtype": "text",
            "text": map[string]string{
                "content": "定时提醒：该休息一下了！",
            },
        }
        dingClient.SendRobotMessage(chatID, msg)
    }
}()
```

#### 2. 集成外部 API

```go
case "天气":
    weather := getWeatherFromAPI() // 调用天气 API
    msg.ReplyToDingtalk(string(message.TEXT), weather)
```

#### 3. 数据库存储

```go
// 记录用户消息
db.Save(msg.SenderNick, msg.Text.Content)
```

## 常见问题

### Q: 为什么机器人没有回复？

**检查清单**:
- [ ] HTTP 服务是否运行？
- [ ] ngrok 隧道是否正常？
- [ ] 回调地址配置是否正确？
- [ ] 代码是否有错误输出？

**调试方法**:
```go
// 添加日志
fmt.Printf("收到回调: %+v\n", msg)
```

### Q: 如何查看详细的错误信息？

```go
statusCode, err := msg.ReplyToDingtalk(string(message.TEXT), "test")
if err != nil {
    fmt.Printf("发送失败: %v, 状态码: %d\n", err, statusCode)
}
```

### Q: SessionWebhook 是什么？

SessionWebhook 是钉钉回调提供的临时 URL，用于在20分钟内向用户回复消息。

### Q: 生产环境部署？

```bash
# 1. 编译
go build -o dingtalk-bot main.go

# 2. 部署到服务器
scp dingtalk-bot user@server:/app/

# 3. 使用 systemd 管理
sudo systemctl start dingtalk-bot
```

## 完整示例项目

查看 `examples/` 目录获取更多示例：

- `examples/basic/` - 基础功能
- `examples/message/` - 消息发送
- `examples/send_message/` - 完整实战示例

## 更多资源

- 📖 [完整文档](./MESSAGE_GUIDE.md)
- 📖 [API 参考](./API_REFERENCE.md)
- 💬 [钉钉开放平台](https://open.dingtalk.com/)
- 🐛 [问题反馈](https://github.com/difyz9/dingtalk-sdk/issues)

## 祝你使用愉快！🚀
