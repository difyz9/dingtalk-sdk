package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	streamclient "github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
)

var messageCount = 0

// 通过 SessionWebhook 发送消息到钉钉
func sendMessageViaWebhook(webhookURL string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("响应状态码: %d", resp.StatusCode)
	}
	
	return nil
}

// 处理机器人收到的消息
func OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	messageCount++
	
	// 打印收到的消息
	fmt.Printf("\n📩 收到第 %d 条消息:\n", messageCount)
	fmt.Printf("  发送人: %s\n", data.SenderNick)
	fmt.Printf("  内容: %s\n", data.Text.Content)
	fmt.Printf("  会话 ID: %s\n", data.ConversationId)
	fmt.Printf("  消息类型: %s\n", data.Msgtype)
	fmt.Printf("  SessionWebhook: %s\n", data.SessionWebhook)
	
	// 根据不同的消息内容发送不同类型的回复
	userMsg := data.Text.Content
	var replyMsg interface{}
	var msgType string
	
	// 演示1: 文本消息回复
	if userMsg == "1" || userMsg == "文本" {
		msgType = "文本消息"
		replyMsg = map[string]interface{}{
			"msgtype": "text",
			"text": map[string]interface{}{
				"content": "✅ 这是一条文本消息回复\n\n当前时间: " + time.Now().Format("2006-01-02 15:04:05"),
			},
		}
	} else if userMsg == "2" || userMsg == "markdown" {
		// 演示2: Markdown 消息回复
		msgType = "Markdown 消息"
		replyMsg = map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]interface{}{
				"title": "Markdown 测试",
				"text":  fmt.Sprintf("### 📊 Stream 模式测试\n\n- **消息序号**: %d\n- **发送人**: %s\n- **时间**: %s\n\n> 这是一条 Markdown 格式的消息", messageCount, data.SenderNick, time.Now().Format("15:04:05")),
			},
		}
	} else if userMsg == "3" || userMsg == "link" || userMsg == "链接" {
		// 演示3: Link 消息回复
		msgType = "Link 消息"
		replyMsg = map[string]interface{}{
			"msgtype": "link",
			"link": map[string]interface{}{
				"title":      "点击查看详情",
				"text":       "这是一条链接消息，点击可以跳转到钉钉官网",
				"messageUrl": "https://www.dingtalk.com",
				"picUrl":     "https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png",
			},
		}
	} else if userMsg == "4" || userMsg == "actioncard" || userMsg == "卡片" {
		// 演示4: ActionCard 消息回复
		msgType = "ActionCard 消息"
		replyMsg = map[string]interface{}{
			"msgtype": "actionCard",
			"actionCard": map[string]interface{}{
				"title":       "任务提醒",
				"text":        "### 📋 您有新的任务待处理\n\n**任务名称**: 测试任务\n**截止时间**: " + time.Now().Add(24*time.Hour).Format("2006-01-02") + "\n\n请及时查看并处理",
				"singleTitle": "查看详情",
				"singleURL":   "https://www.dingtalk.com",
			},
		}
	} else if userMsg == "help" || userMsg == "帮助" || userMsg == "?" {
		// 演示5: 帮助菜单
		msgType = "帮助信息"
		replyMsg = map[string]interface{}{
			"msgtype": "markdown",
			"markdown": map[string]interface{}{
				"title": "使用帮助",
				"text":  "### 🤖 机器人使用指南\n\n发送以下命令测试不同消息类型:\n\n- **1** 或 **文本** - 文本消息\n- **2** 或 **markdown** - Markdown 消息\n- **3** 或 **链接** - Link 消息\n- **4** 或 **卡片** - ActionCard 消息\n- **help** 或 **帮助** - 查看此帮助\n\n---\n\n💡 提示: @我发送消息即可获得回复",
			},
		}
	} else {
		// 默认回复: 智能应答
		msgType = "默认智能应答"
		replyMsg = map[string]interface{}{
			"msgtype": "text",
			"text": map[string]interface{}{
				"content": fmt.Sprintf("收到你的消息: %s\n\n💡 发送 'help' 查看可用命令", userMsg),
			},
		}
	}
	
	fmt.Printf("  → 回复: %s\n", msgType)
	
	// 通过 SessionWebhook 发送消息
	err := sendMessageViaWebhook(data.SessionWebhook, replyMsg)
	if err != nil {
		fmt.Printf("  ❌ 发送失败: %v\n", err)
	} else {
		fmt.Printf("  ✅ 发送成功\n")
	}
	
	// 返回空响应（因为已经通过 webhook 发送了）
	return []byte(`{}`), nil
}

// 处理所有事件
func OnEventReceived(ctx context.Context, df *payload.DataFrame) (*payload.DataFrameResponse, error) {
	eventHeader := event.NewEventHeaderFromDataFrame(df)
	
	fmt.Printf("\n🔔 收到事件 - 类型: %s, ID: %s\n", 
		eventHeader.EventType, eventHeader.EventId)
	
	// 返回成功响应
	return event.NewSuccessResponse()
}

func main() {
	// 配置日志
	logger.SetLogger(logger.NewStdTestLoggerWithDebug())
	
	// 从环境变量或配置文件读取凭证
	clientID := "dingd0xxxxxxxxxxxfd6x"
	clientSecret := "qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3"
	
	fmt.Println("=== 钉钉 Stream 模式 - 消息接收与回复示例 ===")
	fmt.Println()
	fmt.Println("💡 功能说明:")
	fmt.Println("  - Stream 模式可以接收用户发送的消息")
	fmt.Println("  - 自动回复不同类型的消息")
	fmt.Println("  - 支持文本、Markdown、Link、ActionCard 等多种消息类型")
	fmt.Println()
	fmt.Println("📝 使用方法:")
	fmt.Println("  1. 在钉钉群中 @机器人")
	fmt.Println("  2. 发送以下命令:")
	fmt.Println("     - '1' 或 '文本' → 测试文本消息")
	fmt.Println("     - '2' 或 'markdown' → 测试 Markdown 消息")
	fmt.Println("     - '3' 或 '链接' → 测试 Link 消息")
	fmt.Println("     - '4' 或 '卡片' → 测试 ActionCard 消息")
	fmt.Println("     - 'help' 或 '帮助' → 查看帮助信息")
	fmt.Println()
	fmt.Println("🔌 正在连接 Stream 服务器...")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	
	// 使用 NewStreamClient + options 模式创建客户端
	cli := streamclient.NewStreamClient(
		streamclient.WithAppCredential(streamclient.NewAppCredentialConfig(clientID, clientSecret)),
	)
	
	// 注册聊天机器人消息处理器
	cli.RegisterChatBotCallbackRouter(OnChatBotMessageReceived)
	
	// 注册所有事件处理器
	cli.RegisterAllEventRouter(OnEventReceived)
	
	// 启动 Stream 客户端
	err := cli.Start(context.Background())
	if err != nil {
		log.Fatalf("❌ 启动 Stream 客户端失败: %v", err)
		return
	}
	
	defer cli.Close()
	
	fmt.Println("✅ Stream 客户端已启动成功！")
	fmt.Println("💬 等待接收消息...（在群聊中 @机器人 发送消息测试）")
	fmt.Println()
	
	// 阻塞主线程，保持连接
	select {}
}

