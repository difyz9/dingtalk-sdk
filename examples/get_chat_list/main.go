package main

import (
	"fmt"
	"log"

	"github.com/difyz9/dingtalk-sdk.git/client"
)

func main() {
	fmt.Println("=== 钉钉获取 OpenConversationId 示例 ===\n")

	// 创建钉钉客户端
	credential := client.Credential{
		// ClientID:     "your_client_id",     // 替换为你的 Client ID
		// ClientSecret: "your_client_secret", // 替换为你的 Client Secret
		        ClientID:     "dingd0xxxxxxxxxxxfd6x",     // 替换为你的 Client ID
        ClientSecret: "qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3", // 替换为你的 Client Secret
	}

	dingClient := client.NewDingTalkClient(credential)

	// 1. 获取 Access Token
	fmt.Println("1. 获取 Access Token...")
	token, err := dingClient.GetAccessToken()
	if err != nil {
		log.Fatalf("❌ 获取 Access Token 失败: %v", err)
	}
	fmt.Printf("✅ Access Token: %s\n\n", token)

	// 2. 通过 chatId 获取 OpenConversationId
	// 注意：chatId 通常在创建群时获得，或从群信息中获取
	fmt.Println("2. 获取 OpenConversationId...\n")
	
	// 示例 chatId（需要替换为实际值）
	chatID := "your_chat_id" // 替换为实际的群 chatId
	
	if chatID == "your_chat_id" {
		fmt.Println("⚠️  请先设置实际的 chatID")
		fmt.Println("\n💡 获取 chatID 的方法:")
		fmt.Println("  1. 创建群时会返回 chatId")
		fmt.Println("  2. 从群信息查询接口获取")
		fmt.Println("  3. 从消息回调的 conversationId 获取\n")
		fmt.Println("示例代码:")
		fmt.Println("  chatID := \"chatfaabe59a460527f5fb72fbbdfe3f061e\"")
		return
	}

	openConversationId, err := dingClient.GetOpenConversationId(chatID)
	if err != nil {
		log.Fatalf("❌ 获取 OpenConversationId 失败: %v\n\n可能的原因:\n  • chatId 不正确\n  • 应用未开通群基础信息读权限\n  • Client ID/Secret 错误\n", err)
	}

	fmt.Printf("✅ 成功获取 OpenConversationId\n\n")
	fmt.Println("=" + fmt.Sprintf("%60s", "="))
	fmt.Printf("ChatID:              %s\n", chatID)
	fmt.Printf("OpenConversationId:  %s\n", openConversationId)
	fmt.Println("=" + fmt.Sprintf("%60s", "="))

	fmt.Println("\n=== 完成 ===")
	fmt.Println("\n💡 使用说明:")
	fmt.Println("  • OpenConversationId 可用于发送群消息")
	fmt.Println("  • 可以在 examples/send_message 中使用这个 OpenConversationId")
	fmt.Println("  • 确保机器人已加入目标群聊")
	fmt.Println("\n📚 相关文档:")
	fmt.Println("  • API文档: https://open.dingtalk.com/document/development/obtain-group-openconversationid")
}
