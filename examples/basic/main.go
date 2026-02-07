package main

import (
	"fmt"
	"log"

	"github.com/dingtalk-sdk/client"
)

func main() {
	// 创建钉钉客户端凭证
	credential := client.Credential{
		ClientID:     "your_client_id",     // 替换为你的 Client ID
		ClientSecret: "your_client_secret", // 替换为你的 Client Secret
	}

	// 创建钉钉客户端
	dingClient := client.NewDingTalkClient(credential)

	// 获取 Access Token
	token, err := dingClient.GetAccessToken()
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	fmt.Printf("Access Token: %s\n", token)

	// 上传图片示例（需要替换为实际的图片数据）
	// imageData, _ := ioutil.ReadFile("example.png")
	// result, err := dingClient.UploadMedia(
	// 	imageData,
	// 	"example.png",
	// 	client.MediaTypeImage,
	// 	client.MimeTypeImagePng,
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to upload media: %v", err)
	// }
	// fmt.Printf("Media ID: %s\n", result.MediaID)
}
