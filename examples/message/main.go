package main

import (
	"fmt"
	"log"

	"github.com/difyz9/dingtalk-sdk.git/client"
	"github.com/difyz9/dingtalk-sdk.git/message"
)

func main() {
	fmt.Println("=== 钉钉消息发送示例 ===")

	// 方式1: 使用 SessionWebhook 回复消息（适用于接收到用户消息后的回复）
	fmt.Println("\n--- 方式1: 使用 SessionWebhook 回复消息 ---")
	//demoSessionWebhookReply()

	// 方式2: 使用企业内部机器人发送群消息（需要 chatId）
	fmt.Println("\n--- 方式2: 使用企业内部机器人发送群消息 ---")
	demoRobotGroupMessage()
}

// 方式1: 通过 SessionWebhook 回复消息
// 这种方式用于接收到钉钉回调后,使用回调中的 SessionWebhook 进行回复
func demoSessionWebhookReply() {
	// 模拟接收到的消息（实际使用中，这些信息来自钉钉的回调）
	receiveMsg := message.ReceiveMsg{
		// 注意: SessionWebhook 在实际使用中来自钉钉回调,这里仅作演示
		// 实际的 webhook 格式: https://oapi.dingtalk.com/robot/sendBySession?session=xxx
		SessionWebhook:    "https://oapi.dingtalk.com/robot/sendBySession?session=your_session_token",
		SenderNick:        "张三",
		SenderStaffId:     "user123",
		ConversationType:  "2", // 1: 私聊, 2: 群聊
		ConversationTitle: "技术交流群",
		Text: message.Text{
			Content: "你好，机器人",
		},
	}

	// 1. 发送文本消息
	fmt.Println("\n1. 发送文本消息...")
	statusCode, err := receiveMsg.ReplyToDingtalk(
		string(message.TEXT),
		"你好！这是一条文本消息。\n\n当前时间: 2026-02-07",
	)
	if err != nil {
		log.Printf("发送文本消息失败: %v", err)
	} else {
		fmt.Printf("✅ 文本消息发送成功, HTTP状态码: %d\n", statusCode)
	}

	// 2. 发送 Markdown 消息
	fmt.Println("\n2. 发送 Markdown 消息...")
	markdownText := `**欢迎使用钉钉 SDK** 🎉

> 这是一个功能强大的 Go 语言钉钉 SDK

### 主要特性：
- ✅ 消息发送（文本、Markdown）
- ✅ 流式卡片支持
- ✅ 媒体文件上传
- ✅ Access Token 自动管理

### 使用方法：
1. 创建客户端
2. 获取 Access Token
3. 发送消息或上传文件

---
**当前版本**: v1.0.0
`
	statusCode, err = receiveMsg.ReplyToDingtalk(string(message.MARKDOWN), markdownText)
	if err != nil {
		log.Printf("发送 Markdown 消息失败: %v", err)
	} else {
		fmt.Printf("✅ Markdown 消息发送成功, HTTP状态码: %d\n", statusCode)
	}

	// 3. 获取发送者信息
	fmt.Println("\n3. 获取消息元信息...")
	senderID := receiveMsg.GetSenderIdentifier()
	fmt.Printf("发送者标识: %s\n", senderID)

	chatTitle := receiveMsg.GetChatTitle()
	fmt.Printf("聊天标题: %s\n", chatTitle)
}

// 方式2: 使用企业内部机器人发送群消息
// 这种方式需要知道群的 chatId，可以主动向群发送消息
func demoRobotGroupMessage() {
	// 创建钉钉客户端
	credential := client.Credential{
		ClientID:     "dingd0xxxxxxxxxxxfd6x",     // 替换为你的 Client ID (AppKey)
		ClientSecret: "qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3", // 替换为你的 Client Secret
	}

	dingClient := client.NewDingTalkClient(credential)

	// 获取 Access Token（用于验证）
	token, err := dingClient.GetAccessToken()
	if err != nil {
		log.Printf("获取 Access Token 失败: %v", err)
		return
	}
	fmt.Printf("✅ Access Token: %s\n", token)

	// chatId 是群会话的ID，可以通过钉钉开放平台的接口获取
	// 或者在群里发送"群ID"命令让机器人回复
	// chatID := "chat4f4ed5da91cc6500c640ed463645a8d3" // 替换为实际的群 chatId
	chatID:= "cid1+dPH/0LUVUSBFDIcYjYSA==" // 替换为实际的群 chatId --- IGNORE ---

	// 发送文本消息到群
	fmt.Println("\n发送文本消息到群...")
	textMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "大家好！这是来自机器人的群消息 🤖",
		},
	}

	err = dingClient.SendRobotMessage(chatID, textMsg)
	if err != nil {
		log.Printf("发送群消息失败: %v\n", err)
		log.Println("提示: 请确保:")
		log.Println("  1. Client ID 和 Client Secret 正确")
		log.Println("  2. chatId 是有效的群会话 ID")
		log.Println("  3. 机器人已经加入到该群")
	} else {
		fmt.Println("✅ 群消息发送成功")
	}

	// 发送 Markdown 消息到群
	fmt.Println("\n发送 Markdown 消息到群...")
	markdownMsg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "每日报告",
			"text": `### 今日数据统计 📊

**系统状态**: 🟢 正常运行

| 指标 | 数值 |
|------|------|
| 活跃用户 | 1,234 |
| 新增用户 | 56 |
| 错误率 | 0.1% |

> 数据更新时间: 2026-02-07 10:00:00
`,
		},
	}

	err = dingClient.SendRobotMessage(chatID, markdownMsg)
	if err != nil {
		log.Printf("发送 Markdown 群消息失败: %v\n", err)
	} else {
		fmt.Println("✅ Markdown 群消息发送成功")
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 使用提示:")
	fmt.Println("  • 方式1适用于响应用户消息的场景")
	fmt.Println("  • 方式2适用于主动推送消息的场景")
	fmt.Println("  • 实际使用时请替换示例中的凭证和ID")
}
