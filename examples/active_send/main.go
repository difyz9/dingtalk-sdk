package main

import (
	"fmt"
	"log"

	"github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
	fmt.Println("=== 钉钉主动发送消息示例 ===\n")

	// ==================== 方式 1: 企业内部机器人 (推荐) ====================
	fmt.Println("【方式 1】使用企业内部机器人主动发送消息")
	fmt.Println("适用场景: 向企业内部群聊或单聊发送消息")
	fmt.Println("优势: 功能强大，支持各种消息类型\n")

	// 创建钉钉客户端
	credential := client.Credential{
		ClientID:     "dingd0xxxxxxxxxxxfd6x",
		ClientSecret: "qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3",
	}
	dingClient := client.NewDingTalkClient(credential)

	// 1.1 发送文本消息到群聊
	fmt.Println("1. 发送文本消息到群聊...")
	chatID := "your_chat_id" // 替换为实际的群 chatId

	textMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "🤖 这是一条主动发送的消息\n\n发送时间: 2026-02-07",
		},
		// 可选: @指定用户
		// "at": map[string]interface{}{
		// 	"atUserIds": []string{"user_id_1", "user_id_2"},
		// 	"isAtAll":   false,
		// },
	}

	err := dingClient.SendRobotMessage(chatID, textMsg)
	if err != nil {
		log.Printf("❌ 发送文本消息失败: %v\n", err)
	} else {
		fmt.Println("✅ 文本消息发送成功\n")
	}

	// 1.2 发送 Markdown 消息
	fmt.Println("2. 发送 Markdown 消息...")
	markdownMsg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "系统通知",
			"text": `### 📢 重要通知

**发送方式**: 企业内部机器人主动推送

**功能特性**:
- ✅ 支持富文本格式
- ✅ 支持 @用户
- ✅ 支持各种消息类型

> 💡 提示: 这是主动发送的消息，不需要用户触发

---
**时间**: 2026-02-07
**状态**: 🟢 正常`,
		},
	}

	err = dingClient.SendRobotMessage(chatID, markdownMsg)
	if err != nil {
		log.Printf("❌ 发送 Markdown 消息失败: %v\n", err)
	} else {
		fmt.Println("✅ Markdown 消息发送成功\n")
	}

	// 1.3 发送链接消息
	fmt.Println("3. 发送链接消息...")
	linkMsg := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      "查看详情",
			"text":       "点击查看完整内容",
			"messageUrl": "https://www.dingtalk.com",
			"picUrl":     "https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png",
		},
	}

	err = dingClient.SendRobotMessage(chatID, linkMsg)
	if err != nil {
		log.Printf("❌ 发送链接消息失败: %v\n", err)
	} else {
		fmt.Println("✅ 链接消息发送成功\n")
	}

	// 1.4 发送 ActionCard 消息
	fmt.Println("4. 发送 ActionCard 消息...")
	actionCardMsg := map[string]interface{}{
		"msgtype": "actionCard",
		"actionCard": map[string]interface{}{
			"title": "重要通知",
			"text": `### 系统维护通知

**维护时间**: 2026-02-08 02:00-04:00

**影响范围**: 部分功能暂时不可用

请提前做好准备！`,
			"singleTitle": "知道了",
			"singleURL":   "https://www.dingtalk.com",
		},
	}

	err = dingClient.SendRobotMessage(chatID, actionCardMsg)
	if err != nil {
		log.Printf("❌ 发送 ActionCard 消息失败: %v\n", err)
	} else {
		fmt.Println("✅ ActionCard 消息发送成功\n")
	}

	// ==================== 方式 2: Webhook 自定义机器人 ====================
	fmt.Println("\n【方式 2】使用 Webhook 自定义机器人发送消息")
	fmt.Println("适用场景: 向特定群聊发送通知消息")
	fmt.Println("优势: 配置简单，无需 OAuth 认证\n")

	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=cc444c66b477c4a83014535b461dc40b02d7ab7a45b4b1ea235b17e158c8a644"

	// 2.1 通过 Webhook 发送文本消息
	fmt.Println("1. 通过 Webhook 发送文本消息...")
	webhookTextMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": "📣 Webhook 主动推送: 系统运行正常",
		},
	}

	err = dingClient.SendWebhookMessage(webhookURL, webhookTextMsg)
	if err != nil {
		log.Printf("❌ Webhook 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ Webhook 消息发送成功\n")
	}

	// 2.2 通过 Webhook 发送 Markdown
	fmt.Println("2. 通过 Webhook 发送 Markdown...")
	webhookMarkdownMsg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": "数据报告",
			"text": `### 📊 每日数据报告

**日期**: 2026-02-07

| 指标 | 数值 |
|------|------|
| 新增用户 | 1,234 |
| 活跃用户 | 5,678 |
| 订单量 | 890 |

> ✅ 所有指标正常`,
		},
	}

	err = dingClient.SendWebhookMessage(webhookURL, webhookMarkdownMsg)
	if err != nil {
		log.Printf("❌ Webhook Markdown 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ Webhook Markdown 发送成功\n")
	}

	// ==================== 总结 ====================
	fmt.Println("\n=== 发送方式对比 ===")
	fmt.Println("┌─────────────────┬──────────────────┬──────────────────┐")
	fmt.Println("│ 特性            │ 企业内部机器人   │ Webhook 机器人   │")
	fmt.Println("├─────────────────┼──────────────────┼──────────────────┤")
	fmt.Println("│ 配置难度        │ 中等             │ 简单             │")
	fmt.Println("│ 需要 OAuth      │ 是               │ 否               │")
	fmt.Println("│ 功能丰富度      │ 高               │ 中等             │")
	fmt.Println("│ 可发送群聊/单聊 │ 是               │ 仅群聊           │")
	fmt.Println("│ @指定用户       │ 支持             │ 支持             │")
	fmt.Println("│ 适用场景        │ 企业内部通知     │ 简单群通知       │")
	fmt.Println("└─────────────────┴──────────────────┴──────────────────┘")

	fmt.Println("\n💡 提示:")
	fmt.Println("1. chatId 需要从消息回调中获取，或使用 GetOpenConversationId 转换")
	fmt.Println("2. Webhook URL 在群设置 -> 智能群助手 -> 自定义机器人中获取")
	fmt.Println("3. 两种方式可以混合使用，根据场景选择")
}
