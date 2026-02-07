package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
	fmt.Println("=== 钉钉消息发送实战示例 ===\n")

	// ========== 方式一：使用 Webhook 机器人（推荐，最简单） ==========
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=cc444c66b477c4a83014535b461dc40b02d7ab7a45b4b1ea235b17e158c8a644"
	
	fmt.Println("📌 方式一：使用 Webhook 机器人发送消息")
	fmt.Println("优点：简单快速，无需 OAuth 认证\n")

	// 发送文本消息
	textMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "Hello from DingTalk SDK! 🤖\n\n当前时间: 2026-02-07",
		},
	}

	err := client.SendWebhookMessage(webhookURL, textMsg)
	if err != nil {
		log.Printf("❌ Webhook 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ Webhook 消息发送成功！\n")
	}

	// ========== 方式二：使用企业内部机器人（功能更强大） ==========
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📌 方式二：使用企业内部机器人（需要 OAuth 认证）")
	fmt.Println(strings.Repeat("=", 50) + "\n")

	// 创建钉钉客户端（使用您的凭证）
	credential := client.Credential{
		//  ClientID:     "your_client_id",     // 替换为你的 Client ID
		//  ClientSecret: "your_client_secret", // 替换为你的 Client Secret

		        ClientID:     "dingd0xxxxxxxxxxxfd6x",     // 替换为你的 Client ID
        ClientSecret: "qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3", // 替换为你的 Client Secret

	}

	// dingd0xxxxxxxxxxxfd6x

	dingClient := client.NewDingTalkClient(credential)

	// 1. 获取 Access Token（验证凭证是否正确）
	fmt.Println("1. 获取 Access Token...")
	token, err := dingClient.GetAccessToken()
	if err != nil {
		log.Fatalf("❌ 获取 Access Token 失败: %v", err)
	}
	fmt.Printf("✅ Access Token: %s\n\n", token)

	// https://open.dingtalk.com/tools/explorer/jsapi?id=10303
	// 2. 获取群聊 ChatID 说明
	fmt.Println("2. 如何获取群聊 ChatID")
	fmt.Println("=" + fmt.Sprintf("%50s", "="))
	fmt.Println("\n📋 获取 ChatID 的方法:")
	fmt.Println("\n  方法一：通过消息回调获取（推荐）")
	fmt.Println("    当机器人接收到群消息时，钉钉会在回调数据中提供 ConversationID")
	fmt.Println("    这个 ConversationID 就是发送消息时需要的 ChatID")
	fmt.Println("    示例：在 Stream 模式下，从 chatbot.BotCallbackDataModel.ConversationId 获取")
	fmt.Println("\n  方法二：在群里发送 '群ID' 命令")
	fmt.Println("    在群聊中向机器人发送 '群ID'，程序会在日志中输出该群的 ConversationID")
	fmt.Println("    日志格式：企业内部机器人 在『群名』群的ConversationID为: cid...")
	fmt.Println("\n  方法三：查看日志")
	fmt.Println("    机器人接收消息时，会自动记录 ConversationID 到日志")
	fmt.Println("\n=" + fmt.Sprintf("%50s", "="))
	
	// 使用示例 ChatID（需要替换为实际值）
	var chatID string = "cid1+dPH/0LUVUSBFDIcYjYSA==" // 从上述方法获取

	// 3. 发送群消息示例
	// 注意：需要先获取群的 chatId
	// 获取 chatId 的方法：
	//   - 使用上面的 GetChatList() 方法获取
	//   - 在群里让机器人发送"群ID"命令
	//   - 或通过钉钉开放平台 API 获取
	
	if chatID == "" || chatID == "your_chat_id" {
		fmt.Println("\n⚠️  未设置 chatID，跳过发送消息示例")
		fmt.Println("💡 请使用上面的 GetChatList() 获取群聊ID，或手动设置 chatID 变量\n")
	} else {
		fmt.Printf("3. 发送文本消息到群 (ChatID: %s)...\n", chatID)
	textMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "大家好！这是来自钉钉 SDK 的测试消息 🤖\n\n当前时间: 2026-02-07",
		},
	}

	err = dingClient.SendRobotMessage(chatID, textMsg)
	if err != nil {
		log.Printf("⚠️  发送群消息失败: %v\n", err)
		log.Println("\n💡 可能的原因:")
		log.Println("  • chatId 不正确或群不存在")
		log.Println("  • 机器人未加入该群")
		log.Println("  • Client ID/Secret 权限不足")
	} else {
	fmt.Println("4✅ 文本消息发送成功\n")
	}

	// 3. 发送 Markdown 消息
	fmt.Println("3. 发送 Markdown 消息到群...")
	markdownMsg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "系统通知",
			"text": `### 📢 系统状态报告
**运行状态**: 🟢 正常

#### 今日数据统计
| 指标 | 数值 | 趋势 |
|------|------|------|
| 在线用户 | 1,234 | ⬆️ +5% |
| 活跃会话 | 567 | ⬆️ +12% |
| 错误率 | 0.1% | ⬇️ -2% |

#### 最新更新
- ✅ 新增消息发送功能
- ✅ 优化 Token 缓存机制
- ✅ 支持 Markdown 格式

> 数据更新时间: 2026-02-07 10:00:00
> 
> 如有问题请联系管理员

---
**Powered by DingTalk SDK v1.0**
`,
		},
	}

	err = dingClient.SendRobotMessage(chatID, markdownMsg)
	if err != nil {
		log.Printf("⚠️  发送 Markdown 消息失败: %v\n", err)
	} else {
		fmt.Println("✅ Markdown 消息发送成功\n")
	}

	// 5. 发送链接消息
	fmt.Println("5. 发送链接消息到群...")
	linkMsg := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      "钉钉开放平台文档",
			"text":       "查看更多钉钉机器人开发文档，了解如何使用各种消息类型和功能。",
			"messageUrl": "https://open.dingtalk.com/document/",
			"picUrl":     "https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png",
		},
	}

	err = dingClient.SendRobotMessage(chatID, linkMsg)
	if err != nil {
		log.Printf("⚠️  发送链接消息失败: %v\n", err)
	} else {
		fmt.Println("✅ 链接消息发送成功\n")
	}
	} // 结束 chatID 检查

	fmt.Println("=== 示例完成 ===\n")
	fmt.Println("💡 使用说明:")
	fmt.Println("  1. 已集成自动获取群聊列表功能")
	fmt.Println("  2. 如需手动指定 chatID，请修改代码中的 chatID 变量")
	fmt.Println("  3. 确保机器人已经加入到目标群")
	fmt.Println("  4. 确保机器人有发送消息的权限")
	fmt.Println("\n📚 获取 chatID 的方法:")
	fmt.Println("  • ✅ 使用 GetChatList() 方法自动获取（推荐）")
	fmt.Println("  • 在群里向机器人发送 '群ID' 命令")
	fmt.Println("  • 使用钉钉开放平台的群管理 API")
	fmt.Println("  • 查看机器人接收消息时的 conversationId 字段")
}
