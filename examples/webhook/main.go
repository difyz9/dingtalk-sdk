package main

import (
	"fmt"
	"log"

	"github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
	fmt.Println("=== 钉钉 Webhook 机器人消息发送示例 ===\n")

	// Webhook URL（自定义机器人）
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=cc444c66b477c4a83014535b461dc40b02d7ab7a45b4b1ea235b17e158c8a644"

	fmt.Println("✅ 使用 Webhook 方式发送消息\n")
	fmt.Println("📝 Webhook URL:", webhookURL[:60]+"...\n")

	// 1. 发送文本消息
	fmt.Println("1. 发送文本消息...")
	textMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "Hello from DingTalk SDK Webhook! 🤖\n\n这是通过自定义机器人发送的消息。",
		},
	}

	err := client.SendWebhookMessage(webhookURL, textMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ 文本消息发送成功\n")
	}

	// 2. 发送 Markdown 消息
	fmt.Println("2. 发送 Markdown 消息...")
	markdownMsg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "系统通知",
			"text": `### 📢 Webhook 机器人测试

**功能特点**:
- ✅ 无需 OAuth 认证
- ✅ 直接使用 Webhook URL
- ✅ 简单快速

#### 使用场景
1. 告警通知
2. 日志推送
3. 状态监控

> 更新时间: 2026-02-07

---
**Powered by DingTalk SDK**
`,
		},
	}

	err = client.SendWebhookMessage(webhookURL, markdownMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ Markdown 消息发送成功\n")
	}

	// 3. 发送链接消息
	fmt.Println("3. 发送链接消息...")
	linkMsg := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      "钉钉开放平台",
			"text":       "了解更多自定义机器人的使用方法",
			"messageUrl": "https://open.dingtalk.com/document/robots/custom-robot-access",
			"picUrl":     "https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png",
		},
	}

	err = client.SendWebhookMessage(webhookURL, linkMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ 链接消息发送成功\n")
	}

	// 4. 发送 @某人的消息
	fmt.Println("4. 发送 @所有人的消息...")
	atMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "重要通知：请所有人查看！",
		},
		"at": map[string]interface{}{
			"isAtAll": true, // @所有人
			// 或者 @指定人:
			// "atMobiles": []string{"138xxxxxxxx"},
			// "atUserIds": []string{"user123"},
		},
	}

	err = client.SendWebhookMessage(webhookURL, atMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ @所有人消息发送成功\n")
	}

	fmt.Println("=== 示例完成 ===\n")
	fmt.Println("💡 Webhook 机器人 vs 企业内部机器人:")
	fmt.Println("\n  📌 Webhook 机器人（自定义机器人）:")
	fmt.Println("    • 优点: 简单，无需 OAuth 认证")
	fmt.Println("    • 缺点: 功能受限，只能发送消息")
	fmt.Println("    • 适用: 告警、通知、日志推送")
	fmt.Println("\n  📌 企业内部机器人:")
	fmt.Println("    • 优点: 功能强大，可接收消息、管理群等")
	fmt.Println("    • 缺点: 需要 OAuth 认证")
	fmt.Println("    • 适用: 复杂的交互场景")
	fmt.Println("\n📚 参考文档:")
	fmt.Println("  https://open.dingtalk.com/document/robots/custom-robot-access")
}
