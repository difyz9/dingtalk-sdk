// 使用阿里云官方 SDK 发送消息测试
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/difyz9/dingtalk-sdk.git/client"
	dingtalkrobot_1_0 "github.com/alibabacloud-go/dingtalk/robot_1_0"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

/**
 * 使用 Token 初始化账号Client
 * @return Client
 * @throws Exception
 */
func CreateClient() (_result *dingtalkrobot_1_0.Client, _err error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	_result = &dingtalkrobot_1_0.Client{}
	_result, _err = dingtalkrobot_1_0.NewClient(config)
	return _result, _err
}

func _main(args []*string) (_err error) {
	// 配置凭证
	clientID := "dingd0xxxxxxxxxxxfd6x"
	clientSecret := "qbxr1T5_deG9UPxcu1-Ek_xxxxxxxxxxx_KpA0OjLCUBb6wnOLN3"
	openConversationId := "cid1+dPH/0LUVUSBFDIcYjYSA=="
	
	fmt.Println("=== 阿里云官方 SDK 发送消息测试 ===")
	fmt.Println()
	
	// 步骤 1: 获取 AccessToken
	fmt.Println("【步骤 1】获取 AccessToken...")
	credential := client.Credential{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	dingClient := client.NewDingTalkClient(credential)
	
	accessToken, err := dingClient.GetAccessToken()
	if err != nil {
		fmt.Printf("❌ 获取 AccessToken 失败: %v\n", err)
		return err
	}
	fmt.Printf("✅ AccessToken: %s...\n", accessToken[:20])
	fmt.Println()
	
	// 步骤 2: 创建阿里云 SDK 客户端
	fmt.Println("【步骤 2】创建阿里云 SDK 客户端...")
	aliClient, _err := CreateClient()
	if _err != nil {
		fmt.Printf("❌ 创建客户端失败: %v\n", _err)
		return _err
	}
	fmt.Println("✅ 客户端创建成功")
	fmt.Println()

	// 步骤 3: 发送消息（注意：此 API 需要 RobotCode，但我们没有）
	fmt.Println("【步骤 3】尝试使用 OrgGroupSend API 发送消息...")
	fmt.Println("⚠️  注意: OrgGroupSend API 需要 RobotCode（机器人编码）")
	fmt.Println("   这个参数通常在企业内部应用中配置")
	fmt.Println()
	
	orgGroupSendHeaders := &dingtalkrobot_1_0.OrgGroupSendHeaders{}
	orgGroupSendHeaders.XAcsDingtalkAccessToken = tea.String(accessToken)
	orgGroupSendRequest := &dingtalkrobot_1_0.OrgGroupSendRequest{
		MsgParam:           tea.String("{\"content\":\"📢 通过阿里云官方SDK发送的测试消息\"}"),
		MsgKey:             tea.String("sampleText"),
		OpenConversationId: tea.String(openConversationId),
		RobotCode:          tea.String("dingd0xxxxxxxxxxxfd6x"), // 尝试使用 ClientID
		CoolAppCode:        tea.String(""),
	}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		result, _err := aliClient.OrgGroupSendWithOptions(orgGroupSendRequest, orgGroupSendHeaders, &util.RuntimeOptions{})
		if _err != nil {
			return _err
		}

		fmt.Println("✅ 发送成功！")
		fmt.Printf("响应: %v\n", result)
		return nil
	}()

	if tryErr != nil {
		var sdkErr = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			sdkErr = _t
		} else {
			sdkErr.Message = tea.String(tryErr.Error())
		}
		
		fmt.Println("❌ 发送失败:")
		if !tea.BoolValue(util.Empty(sdkErr.Code)) && !tea.BoolValue(util.Empty(sdkErr.Message)) {
			fmt.Printf("  错误代码: %s\n", tea.StringValue(sdkErr.Code))
			fmt.Printf("  错误信息: %s\n", tea.StringValue(sdkErr.Message))
		} else {
			fmt.Printf("  %v\n", tryErr)
		}
		
		fmt.Println()
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println()
		fmt.Println("📝 结论:")
		fmt.Println("  阿里云官方 SDK 的 OrgGroupSend API 需要以下条件:")
		fmt.Println("  1. ✅ AccessToken - 已成功获取")
		fmt.Println("  2. ❌ RobotCode - 需要在钉钉开放平台配置机器人")
		fmt.Println("  3. ✅ OpenConversationId - 已有")
		fmt.Println()
		fmt.Println("💡 推荐方案:")
		fmt.Println("  → 使用 Webhook 自定义机器人（最简单）")
		fmt.Println("    参考: examples/webhook/main.go")
		fmt.Println()
		fmt.Println("  → 使用 Stream 模式（支持交互）")
		fmt.Println("    参考: examples/stream_v2/main.go")
	}
	
	return nil
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		fmt.Printf("执行出错: %v\n", err)
	}
}
