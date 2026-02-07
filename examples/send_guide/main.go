package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== 主动发送消息指南 ===\n")

	fmt.Println("💡 主动发送消息的最佳实践:")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("方式 1: Webhook 自定义机器人 (✅ 强烈推荐)")
	fmt.Println("  优势: 配置简单，无需复杂认证，稳定可靠")
	fmt.Println("  获取: 群设置 -> 智能群助手 -> 添加机器人 -> 自定义")
	fmt.Println()
	fmt.Println("方式 2: Stream 模式")
	fmt.Println("  优势: 支持双向通信，可以接收和发送消息")
	fmt.Println("  适用: 需要交互的场景")
	fmt.Println()
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()

	// ==================== Webhook 方式 ====================
	fmt.Println("【推荐】使用 Webhook 发送消息")
	fmt.Println()

	fmt.Println("步骤 1: 获取 Webhook URL")
	fmt.Println("  1. 打开钉钉群聊")
	fmt.Println("  2. 点击群设置 -> 智能群助手")
	fmt.Println("  3. 添加机器人 -> 自定义")
	fmt.Println("  4. 复制 Webhook URL")
	fmt.Println()

	fmt.Println("步骤 2: 使用 SDK 发送消息")
	fmt.Println()

	// ==================== 示例代码 ====================
	fmt.Println("完整示例代码:")
	fmt.Println("```go")
	fmt.Println(`package main

import (
    "fmt"
    "time"
    "github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
    // 创建客户端
    credential := client.Credential{
        ClientID:     "any",
        ClientSecret: "any",
    }
    dingClient := client.NewDingTalkClient(credential)
    
    // 你的 Webhook URL
    webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=你的token"
    
    // 1. 文本消息
    textMsg := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]interface{}{
            "content": "Hello, 钉钉！",
        },
    }
    dingClient.SendWebhookMessage(webhookURL, textMsg)
    
    // 2. Markdown 消息
    markdownMsg := map[string]interface{}{
        "msgtype": "markdown",
        "markdown": map[string]interface{}{
            "title": "通知标题",
            "text": "### 重要通知\n\n**内容**: 这是测试消息",
        },
    }
    dingClient.SendWebhookMessage(webhookURL, markdownMsg)
    
    // 3. 链接消息
    linkMsg := map[string]interface{}{
        "msgtype": "link",
        "link": map[string]string{
            "title":      "查看详情",
            "text":       "点击查看",
            "messageUrl": "https://www.dingtalk.com",
            "picUrl":     "https://example.com/image.png",
        },
    }
    dingClient.SendWebhookMessage(webhookURL, linkMsg)
}`)
	fmt.Println("```")
	fmt.Println()

	// ==================== Stream 模式说明 ====================
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("【高级】Stream 模式 - 支持消息接收和发送")
	fmt.Println()
	fmt.Println("Stream 模式可以:")
	fmt.Println("  ✅ 接收用户发送的消息")
	fmt.Println("  ✅ 主动回复消息")
	fmt.Println("  ✅ 处理各种事件")
	fmt.Println()
	fmt.Println("参考: examples/stream_v2/main.go")
	fmt.Println()

	// ==================== 关于 chatId 的说明 ====================
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("📝 关于你提供的群聊信息:")
	fmt.Println()
	fmt.Println("  chatId: chat52fb673c7b0c7722facfe07d6b48dbb6")
	fmt.Println("  openConversationId: cid1+dPH/0LUVUSBFDIcYjYSA==")
	fmt.Println()
	fmt.Println("说明:")
	fmt.Println("  - chatId 是钉钉内部的会话标识")
	fmt.Println("  - openConversationId 是开放平台的标准 ID")
	fmt.Println("  - 主动发送推荐使用 Webhook 方式(无需这些 ID)")
	fmt.Println("  - 如需在消息回调中使用，可以直接用 openConversationId")
	fmt.Println()

	// ==================== 总结 ====================
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("📚 完整文档和示例:")
	fmt.Println("  - docs/ACTIVE_SEND_GUIDE.md - 主动发送消息完整指南")
	fmt.Println("  - docs/STREAM_V2_GUIDE.md - Stream 模式使用指南")
	fmt.Println("  - examples/webhook/main.go - Webhook 完整示例")
	fmt.Println("  - examples/stream_v2/main.go - Stream 模式示例")
	fmt.Println()
	fmt.Println("🚀 快速开始:")
	fmt.Println("  1. 获取 Webhook URL (推荐)")
	fmt.Println("  2. 运行 examples/webhook/main.go 测试")
	fmt.Println("  3. 查看群聊收到的消息")
	fmt.Println()
}
