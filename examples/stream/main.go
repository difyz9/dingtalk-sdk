package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dingtalk-sdk/client"
	"github.com/dingtalk-sdk/stream"
	"github.com/google/uuid"
)

func main() {
	// 创建钉钉客户端
	credential := client.Credential{
		ClientID:     "your_client_id",     // 替换为你的 Client ID
		ClientSecret: "your_client_secret", // 替换为你的 Client Secret
	}

	dingClient := client.NewDingTalkClient(credential)

	// 获取 Access Token
	accessToken, err := dingClient.GetAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// 创建流式卡片客户端
	streamClient, err := stream.NewStreamCardClient()
	if err != nil {
		log.Fatalf("Failed to create stream card client: %v", err)
	}

	// 生成唯一的追踪 ID
	trackID := uuid.New().String()

	// 创建并投放卡片
	fmt.Println("创建并投放流式卡片...")
	cardReq := &stream.CreateAndDeliverCardRequest{
		CardTemplateID:   "your_card_template_id", // 替换为你的卡片模板 ID
		OutTrackID:       trackID,
		OpenSpaceID:      "dtv1.card//IM_ROBOT.user_id", // 替换为实际的 OpenSpaceID
		ConversationType: "1",                            // 1: 私聊, 2: 群聊
		RobotCode:        "your_robot_code",              // 替换为你的机器人 code
		CardData: map[string]string{
			"content": "正在处理中...",
		},
	}

	err = streamClient.CreateAndDeliverCard(accessToken, cardReq)
	if err != nil {
		log.Fatalf("Failed to create and deliver card: %v", err)
	}
	fmt.Println("卡片创建成功！")

	// 模拟流式更新
	contents := []string{
		"第一步：分析问题...",
		"第二步：查找资料...",
		"第三步：生成答案...",
		"完成！这是最终答案。",
	}

	for i, content := range contents {
		time.Sleep(2 * time.Second) // 模拟处理时间

		isFinalize := i == len(contents)-1 // 最后一次更新标记为完成
		fmt.Printf("更新卡片内容 (%d/%d): %s\n", i+1, len(contents), content)

		updateReq := &stream.StreamingUpdateRequest{
			OutTrackID: trackID,
			Key:        "content",
			Content:    content,
			IsFull:     true,
			IsFinalize: isFinalize,
		}

		err = streamClient.StreamingUpdate(accessToken, updateReq)
		if err != nil {
			log.Printf("Failed to update stream card: %v", err)
			continue
		}
	}

	fmt.Println("流式卡片更新完成！")
}
