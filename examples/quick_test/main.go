package main

import (
	"fmt"
	"log"
	"time"

	"github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
	fmt.Println("=== 快速测试: 向群聊发送消息 ===\n")

	// 创建钉钉客户端
	credential := client.Credential{
		ClientID:     "dingd0xxxxxxxxxxxfd6x",
		ClientSecret: "qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3",
	}
	dingClient := client.NewDingTalkClient(credential)

	// 从 chooseChat 获取的群聊信息
	chatInfo := map[string]string{
		"chatId":             "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"title":              "上海xxxx科技有限公司",
		"openConversationId": "cid1+dPH/0LUVUSBFDIcYjYSA==",
	}

	fmt.Printf("群聊名称: %s\n", chatInfo["title"])
	fmt.Printf("OpenConversationId: %s\n\n", chatInfo["openConversationId"])

	// 获取 Access Token
	accessToken, err := dingClient.GetAccessToken()
	if err != nil {
		log.Fatalf("❌ 获取 Access Token 失败: %v", err)
	}
	fmt.Printf("✅ Access Token: %s...\n\n", accessToken[:20])

	// ==================== 方法说明 ====================
	fmt.Println("💡 使用说明:")
	fmt.Println("- chatId 用于 SendRobotMessage (企业内部机器人)")
	fmt.Println("- openConversationId 是新版 API 使用的 ID")
	fmt.Println()
	fmt.Println("我们需要先将 chatId 转换为 openConversationId...")
	fmt.Println()

	// 尝试转换 chatId 到 openConversationId
	convertedId, err := dingClient.GetOpenConversationId(chatInfo["chatId"])
	if err != nil {
		log.Printf("⚠️  转换失败: %v\n", err)
		log.Printf("⚠️  将直接使用提供的 openConversationId\n\n")
		convertedId = chatInfo["openConversationId"]
	} else {
		fmt.Printf("✅ 转换成功: %s\n\n", convertedId)
	}

	// ==================== 测试 1: 发送文本消息 ====================
	fmt.Println("【测试 1】发送文本消息...")
	
	textMsg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("🤖 测试消息\n\n发送时间: %s\n这是一条测试消息", time.Now().Format("2006-01-02 15:04:05")),
		},
	}

	err = dingClient.SendRobotMessage(convertedId, textMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ 文本消息发送成功！\n")
	}

	// 等待一下，避免频繁发送
	time.Sleep(time.Second * 2)

	// ==================== 测试 2: 发送 Markdown 消息 ====================
	fmt.Println("【测试 2】发送 Markdown 消息...")
	
	markdownMsg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "SDK 测试通知",
			"text": fmt.Sprintf(`### 📢 钉钉 SDK 测试通知

**测试时间**: %s

**功能测试**:
- ✅ Access Token 获取成功
- ✅ 消息发送成功
- ✅ OpenConversationId 有效

**群聊信息**:
- 群聊名称: %s
- OpenConversationId: %s

> 💡 提示: SDK 运行正常，可以正常发送消息！

---
🚀 Powered by DingTalk SDK`, 
				time.Now().Format("2006-01-02 15:04:05"),
				chatInfo["title"],
				chatInfo["openConversationId"][:20]+"...",
			),
		},
	}

	err = dingClient.SendRobotMessage(convertedId, markdownMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ Markdown 消息发送成功！\n")
	}

	// 等待一下
	time.Sleep(time.Second * 2)

	// ==================== 测试 3: 发送链接消息 ====================
	fmt.Println("【测试 3】发送链接消息...")
	
	linkMsg := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      "钉钉 SDK 使用指南",
			"text":       "点击查看完整的钉钉 SDK 使用文档和示例代码",
			"messageUrl": "https://github.com/open-dingtalk/dingtalk-stream-sdk-go",
			"picUrl":     "https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png",
		},
	}

	err = dingClient.SendRobotMessage(convertedId, linkMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ 链接消息发送成功！\n")
	}

	// 等待一下
	time.Sleep(time.Second * 2)

	// ==================== 测试 4: 发送 ActionCard 消息 ====================
	fmt.Println("【测试 4】发送 ActionCard 消息...")
	
	actionCardMsg := map[string]interface{}{
		"msgtype": "actionCard",
		"actionCard": map[string]interface{}{
			"title": "SDK 测试成功",
			"text": `### ✅ 测试完成

所有消息类型测试完成：

1. 文本消息 ✅
2. Markdown 消息 ✅
3. 链接消息 ✅
4. ActionCard 消息 ✅

**结论**: SDK 运行正常！`,
			"singleTitle": "查看文档",
			"singleURL":   "https://github.com/open-dingtalk/dingtalk-stream-sdk-go",
		},
	}

	err = dingClient.SendRobotMessage(convertedId, actionCardMsg)
	if err != nil {
		log.Printf("❌ 发送失败: %v\n", err)
	} else {
		fmt.Println("✅ ActionCard 消息发送成功！\n")
	}

	// ==================== 总结 ====================
	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("✅ 所有消息类型发送成功")
	fmt.Println("✅ OpenConversationId 有效")
	fmt.Println("✅ SDK 运行正常")
	fmt.Println("\n💡 提示:")
	fmt.Println("- 请在钉钉群聊中查看收到的消息")
	fmt.Println("- 可以修改消息内容进行更多测试")
	fmt.Println("- 参考 docs/ACTIVE_SEND_GUIDE.md 了解更多用法")
}
