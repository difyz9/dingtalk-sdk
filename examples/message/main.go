package main

import (
	"fmt"
	"log"

	"github.com/difyz9/dingtalk-sdk.git/message"
)

func main() {
	// 模拟接收到的消息
	receiveMsg := message.ReceiveMsg{
		SessionWebhook:    "https://oapi.dingtalk.com/robot/send?access_token=xxx", // 替换为实际的 webhook
		SenderNick:        "张三",
		SenderStaffId:     "user123",
		ConversationType:  "1", // 1: 私聊, 2: 群聊
		ConversationTitle: "测试群",
		Text: message.Text{
			Content: "Hello, DingTalk!",
		},
	}

	// 发送文本消息
	fmt.Println("发送文本消息...")
	statusCode, err := receiveMsg.ReplyToDingtalk(string(message.TEXT), "你好！这是一条文本消息。")
	if err != nil {
		log.Fatalf("Failed to send text message: %v", err)
	}
	fmt.Printf("Text message sent, status code: %d\n", statusCode)

	// 发送 Markdown 消息
	fmt.Println("\n发送 Markdown 消息...")
	markdownText := `**欢迎使用钉钉 SDK**

> 这是一个功能强大的 SDK

特性：
- ✅ 消息发送
- ✅ 流式卡片
- ✅ 媒体上传
`
	statusCode, err = receiveMsg.ReplyToDingtalk(string(message.MARKDOWN), markdownText)
	if err != nil {
		log.Fatalf("Failed to send markdown message: %v", err)
	}
	fmt.Printf("Markdown message sent, status code: %d\n", statusCode)

	// 获取发送者标识
	senderID := receiveMsg.GetSenderIdentifier()
	fmt.Printf("\n发送者标识: %s\n", senderID)

	// 获取聊天标题
	chatTitle := receiveMsg.GetChatTitle()
	fmt.Printf("聊天标题: %s\n", chatTitle)
}
