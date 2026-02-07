# 钉钉 SDK 使用指南

本文档详细说明如何使用本 SDK 向钉钉群聊发送消息，包含三种已验证可用的方式。

## 📋 目录

- [前置准备](#前置准备)
- [方式一：Webhook 自定义机器人（推荐）](#方式一webhook-自定义机器人推荐)
- [方式二：Stream V2 模式（官方推荐）](#方式二stream-v2-模式官方推荐)
- [方式三：阿里云官方 SDK](#方式三阿里云官方-sdk)
- [测试成功案例](#测试成功案例)
- [常见问题](#常见问题)

---

## 前置准备

### 1. 获取应用凭证

在钉钉开放平台创建应用后，你将获得：
- **ClientID（AppKey）**: 应用唯一标识
- **ClientSecret（AppSecret）**: 应用密钥

示例：
```
ClientID: dingd0xxxxxxxxxxxfd6x
ClientSecret: qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3
```

### 2. 获取群聊信息

你需要知道目标群的会话 ID：
- **chatId**: 如 `chat52fb673c7b0c7722facfe07d6b48dbb6`
- **openConversationId**: 如 `cid1+dPH/0LUVUSBFDIcYjYSA==`

获取方法：
- 参考钉钉开放平台文档：https://open.dingtalk.com/tools/explorer/jsapi?id=10303
- 或通过 API：`/v1.0/im/conversations/users/{userId}/chatIds` 转换

---

## 方式一：Webhook 自定义机器人（推荐）

### ✅ 优势

- 🟢 **最简单** - 无需复杂配置
- 🟢 **无需公网 IP**
- 🟢 **配置快速** - 5分钟即可完成
- 🟢 **稳定可靠** - 官方标准接口

### 📝 配置步骤

#### 1. 创建自定义机器人

1. 打开钉钉群聊
2. 点击右上角 `···` → 群设置
3. 选择 `智能群助手` → `添加机器人` → `自定义`
4. 填写机器人名称和描述
5. 安全设置：选择 `加签` 或 `自定义关键词`
6. 点击 `完成`，复制 **Webhook URL**

Webhook URL 格式：
```
https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxxxxx
```

#### 2. 发送消息代码

示例代码：[examples/webhook/main.go](../examples/webhook/main.go)

```go
package main

import (
    "fmt"
    "github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
    // 你的 Webhook URL
    webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=your_token"
    
    // 1. 发送文本消息
    textMsg := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]interface{}{
            "content": "📢 测试消息：Webhook 方式发送成功！",
        },
    }
    err := client.SendWebhookMessage(webhookURL, textMsg)
    if err != nil {
        fmt.Printf("❌ 发送失败: %v\n", err)
    } else {
        fmt.Println("✅ 文本消息发送成功")
    }
    
    // 2. 发送 Markdown 消息
    markdownMsg := map[string]interface{}{
        "msgtype": "markdown",
        "markdown": map[string]interface{}{
            "title": "系统通知",
            "text": "### 📊 数据报告\n\n- **状态**: ✅ 正常\n- **时间**: 2026-02-07\n\n> 所有服务运行正常",
        },
    }
    err = client.SendWebhookMessage(webhookURL, markdownMsg)
    if err != nil {
        fmt.Printf("❌ 发送失败: %v\n", err)
    } else {
        fmt.Println("✅ Markdown 消息发送成功")
    }
}
```

#### 3. 运行测试

```bash
go run examples/webhook/main.go
```

### 🎯 适用场景

- 系统告警通知
- 定时任务报告
- 监控数据推送
- 简单的单向消息发送

---

## 方式二：Stream V2 模式（官方推荐）

### ✅ 优势

- 🟢 **无需公网 IP** - 使用 WebSocket 长连接
- 🟢 **支持双向通信** - 可接收和发送消息
- 🟢 **实时响应** - 用户 @机器人 立即回复
- 🟢 **官方推荐** - 钉钉官方最新推荐方式

### 📝 配置步骤

#### 1. 在钉钉开放平台配置

1. 进入你的应用管理页面
2. 开通 `机器人能力`
3. 配置机器人信息（名称、头像、描述）
4. 将机器人添加到测试群聊

#### 2. 接收消息并回复

示例代码：[examples/stream_v2/main.go](../examples/stream_v2/main.go)

```go
package main

import (
    "context"
    "fmt"
    "github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
    streamclient "github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
)

func OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
    fmt.Printf("收到消息: %s\n", data.Text.Content)
    
    // 通过 SessionWebhook 回复消息
    replyMsg := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]interface{}{
            "content": "收到你的消息: " + data.Text.Content,
        },
    }
    
    // 发送回复
    sendMessageViaWebhook(data.SessionWebhook, replyMsg)
    
    return []byte(`{}`), nil
}

func main() {
    clientID := "your_client_id"
    clientSecret := "your_client_secret"
    
    cli := streamclient.NewStreamClient(
        streamclient.WithAppCredential(
            streamclient.NewAppCredentialConfig(clientID, clientSecret),
        ),
    )
    
    cli.RegisterChatBotCallbackRouter(OnChatBotMessageReceived)
    
    cli.Start(context.Background())
    defer cli.Close()
    
    fmt.Println("✅ Stream 客户端已启动，等待接收消息...")
    select {} // 阻塞保持连接
}
```

#### 3. 运行测试

```bash
go run examples/stream_v2/main.go
```

然后在群聊中 @机器人 发送消息，机器人会自动回复。

### 🎯 测试成功案例

**测试群**: 银河护卫队科技有限公司  
**会话 ID**: `cidGCUBTzi5e6/D2Drgx6UHT2cAEyncMJx6pMZePDxhb2k=`

测试结果：
```
📩 收到第 1 条消息:
  发送人: 蜘蛛侠
  内容: 1
  会话 ID: cidGCUBTzi5e6/D2Drgx6UHT2cAEyncMJx6pMZePDxhb2k=
  消息类型: text
  → 回复: 文本消息
  ✅ 发送成功
```

支持的命令：
- `1` 或 `文本` → 测试文本消息
- `2` 或 `markdown` → 测试 Markdown 消息
- `3` 或 `链接` → 测试 Link 消息
- `4` 或 `卡片` → 测试 ActionCard 消息
- `help` 或 `帮助` → 查看帮助信息

### 🎯 适用场景

- 智能客服机器人
- 交互式问答系统
- 任务管理助手
- 需要实时响应的场景

---

## 方式三：阿里云官方 SDK

### ✅ 优势

- 🟢 **官方支持** - 阿里云官方维护
- 🟢 **功能完整** - 支持所有钉钉 API
- 🟢 **文档齐全** - 完整的 API 文档
- 🟢 **已验证可用** - 实测成功发送

### 📝 配置步骤

#### 1. 安装依赖

```bash
go get github.com/alibabacloud-go/dingtalk
go get github.com/alibabacloud-go/tea
go get github.com/alibabacloud-go/darabonba-openapi/v2
```

#### 2. 发送消息代码

示例代码：[examples/alicloud_sdk/main.go](../examples/alicloud_sdk/main.go)

```go
package main

import (
    "fmt"
    "github.com/difyz9/dingtalk-sdk.git/client"
    dingtalkrobot_1_0 "github.com/alibabacloud-go/dingtalk/robot_1_0"
    openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
    "github.com/alibabacloud-go/tea/tea"
)

func CreateClient() (*dingtalkrobot_1_0.Client, error) {
    config := &openapi.Config{}
    config.Protocol = tea.String("https")
    config.RegionId = tea.String("central")
    return dingtalkrobot_1_0.NewClient(config)
}

func main() {
    // 步骤 1: 获取 AccessToken
    credential := client.Credential{
        ClientID:     "your_client_id",
        ClientSecret: "your_client_secret",
    }
    dingClient := client.NewDingTalkClient(credential)
    accessToken, _ := dingClient.GetAccessToken()
    
    // 步骤 2: 创建阿里云 SDK 客户端
    aliClient, _ := CreateClient()
    
    // 步骤 3: 发送消息
    headers := &dingtalkrobot_1_0.OrgGroupSendHeaders{}
    headers.XAcsDingtalkAccessToken = tea.String(accessToken)
    
    request := &dingtalkrobot_1_0.OrgGroupSendRequest{
        MsgParam:           tea.String("{\"content\":\"📢 测试消息\"}"),
        MsgKey:             tea.String("sampleText"),
        OpenConversationId: tea.String("cid1+dPH/0LUVUSBFDIcYjYSA=="),
        RobotCode:          tea.String("your_client_id"),
    }
    
    result, err := aliClient.OrgGroupSendWithOptions(request, headers, nil)
    if err != nil {
        fmt.Printf("❌ 发送失败: %v\n", err)
    } else {
        fmt.Printf("✅ 发送成功: %v\n", result)
    }
}
```

#### 3. 运行测试

```bash
go run examples/alicloud_sdk/main.go
```

### 🎯 测试成功案例

**测试群**: 银河护卫队科技有限公司  
**OpenConversationId**: `cid1+dPH/0LUVUSBFDIcYjYSA==`

测试结果：
```
=== 阿里云官方 SDK 发送消息测试 ===

【步骤 1】获取 AccessToken...
✅ AccessToken: 605e241440c43d8f9244...

【步骤 2】创建阿里云 SDK 客户端...
✅ 客户端创建成功

【步骤 3】尝试使用 OrgGroupSend API 发送消息...
✅ 发送成功！
响应: {
   "statusCode": 200,
   "body": {
      "processQueryKey": "h2Jh2kbkPlnUZ6w3PBSaHaZXM/uYDtWB1UaA6Ihttow="
   }
}
```

### 🎯 适用场景

- 需要使用完整钉钉 API 的场景
- 企业级应用开发
- 需要官方技术支持
- 复杂的钉钉集成需求

---

## 测试成功案例

### 案例 1: Webhook 方式发送消息

**测试信息**:
- 方式: Webhook 自定义机器人
- 群聊: 银河护卫队科技有限公司
- 消息类型: 文本、Markdown、Link、ActionCard

**测试结果**: ✅ 全部成功

### 案例 2: Stream V2 接收并回复消息

**测试信息**:
- 方式: Stream V2 模式
- 群聊: 银河护卫队科技有限公司
- 会话 ID: `cidGCUBTzi5e6/D2Drgx6UHT2cAEyncMJx6pMZePDxhb2k=`
- 交互: 用户发送 "1"、"help"、"5555" 等命令

**测试结果**: ✅ 成功接收消息并自动回复

测试日志：
```
📩 收到第 1 条消息:
  发送人: 蜘蛛侠
  内容: 好
  → 回复: 默认智能应答
  ✅ 发送成功

📩 收到第 2 条消息:
  发送人: 蜘蛛侠
  内容: 1
  → 回复: 文本消息
  ✅ 发送成功

📩 收到第 4 条消息:
  发送人: 蜘蛛侠
  内容: help
  → 回复: 帮助信息
  ✅ 发送成功
```

### 案例 3: 阿里云官方 SDK 发送消息

**测试信息**:
- 方式: 阿里云官方 SDK (OrgGroupSend API)
- 群聊: 银河护卫队科技有限公司
- OpenConversationId: `cid1+dPH/0LUVUSBFDIcYjYSA==`
- ClientID: `dingd0xxxxxxxxxxxfd6x`

**测试结果**: ✅ 成功发送

响应：
```json
{
  "statusCode": 200,
  "body": {
    "processQueryKey": "h2Jh2kbkPlnUZ6w3PBSaHaZXM/uYDtWB1UaA6Ihttow="
  }
}
```

---

## 三种方式对比

| 特性 | Webhook | Stream V2 | 阿里云 SDK |
|------|---------|-----------|-----------|
| **配置难度** | 🟢 最简单 | 🟡 中等 | 🟡 中等 |
| **公网要求** | ❌ 不需要 | ❌ 不需要 | ❌ 不需要 |
| **双向通信** | ❌ 单向 | ✅ 支持 | ✅ 支持 |
| **实时性** | 🟡 HTTP | 🟢 WebSocket | 🟡 HTTP |
| **官方支持** | ✅ 官方 | ⭐ 官方推荐 | ✅ 阿里云官方 |
| **测试状态** | ✅ 成功 | ✅ 成功 | ✅ 成功 |
| **推荐指数** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

### 选择建议

1. **简单通知场景** → 使用 **Webhook** （最快速）
2. **需要交互** → 使用 **Stream V2** （实时响应）
3. **企业级应用** → 使用 **阿里云 SDK** （功能完整）

---

## 常见问题

### Q1: 如何获取 chatId？

**方法 1**: 通过钉钉 JSAPI
```javascript
// 在钉钉内打开网页
dd.biz.chat.pickConversation({
  onSuccess: function(result) {
    console.log('chatId:', result.chatId);
  }
});
```

**方法 2**: 通过 API 转换
```go
openConvId := dingClient.GetOpenConversationId("chatId")
```

参考：https://open.dingtalk.com/tools/explorer/jsapi?id=10303

### Q2: AccessToken 有效期多久？

AccessToken 有效期为 **2 小时**。本 SDK 已自动实现缓存和刷新机制。

### Q3: Webhook 安全设置如何配置？

推荐使用 **加签** 方式：

1. 创建机器人时选择"加签"
2. 获得密钥（secret）
3. 在发送消息时计算签名

```go
// 计算签名的示例代码
timestamp := time.Now().UnixMilli()
sign := calculateSign(timestamp, secret)
webhookURL := fmt.Sprintf("%s&timestamp=%d&sign=%s", baseURL, timestamp, sign)
```

### Q4: Stream 模式如何保持长连接？

Stream 模式使用 WebSocket 长连接，SDK 会自动处理断线重连。只需在 main 函数最后保持阻塞：

```go
select {} // 阻塞主线程，保持连接
```

### Q5: 三种方式都需要 AccessToken 吗？

- **Webhook**: ❌ 不需要（直接使用 Webhook URL）
- **Stream V2**: ✅ 需要（SDK 自动获取）
- **阿里云 SDK**: ✅ 需要（手动获取并传递）

---

## 完整示例代码

所有示例代码位于 [examples/](../examples/) 目录：

- `examples/webhook/` - Webhook 完整示例
- `examples/stream_v2/` - Stream V2 完整示例
- `examples/alicloud_sdk/` - 阿里云 SDK 完整示例
- `examples/active_send/` - 主动发送消息综合示例
- `examples/send_guide/` - 发送消息使用指南

## 相关文档

- [快速开始](QUICK_START.md) - 5分钟入门
- [Stream V2 指南](STREAM_V2_GUIDE.md) - Stream 模式详细文档
- [主动发送消息指南](ACTIVE_SEND_GUIDE.md) - 主动发送完整指南
- [API 文档](API.md) - 完整 API 参考

---

## 技术支持

如有问题，请参考：
1. 钉钉开放平台文档：https://open.dingtalk.com
2. 本项目 Issues：提交问题和建议
3. 示例代码：[examples/](../examples/) 目录

最后更新：2026-02-07
