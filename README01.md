# DingTalk SDK
一个基于 Go 语言实现的钉钉机器人 SDK，提供完整的钉钉消息发送、接收、流式卡片等功能。

## 🎉 已验证可用的三种消息发送方式

| 方式 | 难度 | 推荐度 | 测试状态 |
|------|------|--------|---------|
| **Webhook 自定义机器人** | 🟢 最简单 | ⭐⭐⭐⭐⭐ | ✅ 已验证 |
| **Stream V2 模式** | 🟡 中等 | ⭐⭐⭐⭐ | ✅ 已验证 |
| **阿里云官方 SDK** | 🟡 中等 | ⭐⭐⭐⭐ | ✅ 已验证 |

**测试群**: 银河护卫队科技有限公司  
**测试日期**: 2026-02-07  
**测试结果**: 所有方式均成功发送消息 ✅

👉 **[完整使用指南](docs/USAGE_GUIDE.md)** - 包含详细配置步骤和测试案例

## 功能特性

- ✅ 钉钉机器人消息接收和发送
- ✅ 支持文本、Markdown、链接、ActionCard 等消息格式
- ✅ OAuth 2.0 认证和 AccessToken 管理
- ✅ **Stream V2 模式** - 官方推荐,无需公网 IP
- ✅ **阿里云官方 SDK 集成** - 已验证可用
- ✅ 流式卡片创建和更新
- ✅ 媒体文件（图片、视频、文件）上传
- ✅ Webhook 自定义机器人支持
- ✅ 自动 AccessToken 缓存和刷新

## 快速开始

### 安装

```bash
go get github.com/difyz9/dingtalk-sdk.git
```

### 方式一：Webhook（最简单，5分钟上手）

```go
package main

import "github.com/difyz9/dingtalk-sdk.git/client"

func main() {
    webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN"
    msg := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]interface{}{"content": "📢 测试消息"},
    }
    client.SendWebhookMessage(webhookURL, msg)
}
```

📖 **[完整使用指南](docs/USAGE_GUIDE.md)** - 详细配置步骤和成功案例

---

## 测试成功案例

### ✅ Stream V2 交互式回复
测试日志：
```
📩 收到消息: help
  → 回复: 帮助信息
  ✅ 发送成功
```

### ✅ 阿里云 SDK 发送
响应：
```json
{
  "statusCode": 200,
  "body": {"processQueryKey": "h2Jh2kbkPlnUZ6w3PBSaHaZXM/uYDtWB1UaA6Ihttow="}
}
```

---

## 示例程序

| 示例 | 说明 | 状态 |
|------|------|------|
| [webhook/](examples/webhook/) | Webhook 自定义机器人 | ✅ 已测试 |
| [stream_v2/](examples/stream_v2/) | Stream V2 交互式机器人 | ✅ 已测试 |
| [alicloud_sdk/](examples/alicloud_sdk/) | 阿里云官方 SDK | ✅ 已测试 |

## 文档

- 📖 **[完整使用指南](docs/USAGE_GUIDE.md)** - 详细配置和测试案例 ⭐
- 📖 [快速开始](docs/QUICK_START.md)
- ⭐ [Stream V2 指南](docs/STREAM_V2_GUIDE.md)
- 📝 [主动发送消息指南](docs/ACTIVE_SEND_GUIDE.md)
- 📚 [API 文档](docs/API.md)

## 许可证

MIT License