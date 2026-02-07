package stream

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalk "github.com/alibabacloud-go/dingtalk/card_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/google/uuid"
)

// StreamCardClient 流式卡片客户端
type StreamCardClient struct {
	client *dingtalk.Client
}

// NewStreamCardClient 创建流式卡片客户端
func NewStreamCardClient() (*StreamCardClient, error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	client, err := dingtalk.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &StreamCardClient{
		client: client,
	}, nil
}

// CreateAndDeliverCardRequest 创建并投放卡片请求
type CreateAndDeliverCardRequest struct {
	CardTemplateID   string
	OutTrackID       string
	ConversationID   string
	SenderStaffID    string
	RobotCode        string
	OpenSpaceID      string
	ConversationType string // "1" for private chat, "2" for group chat
	CardData         map[string]string
}

// CreateAndDeliverCard 创建并投放流式卡片
func (s *StreamCardClient) CreateAndDeliverCard(accessToken string, req *CreateAndDeliverCardRequest) error {
	headers := &dingtalk.CreateAndDeliverHeaders{
		XAcsDingtalkAccessToken: tea.String(accessToken),
	}

	cardData := &dingtalk.CreateAndDeliverRequestCardData{
		CardParamMap: make(map[string]*string),
	}
	for k, v := range req.CardData {
		cardData.CardParamMap[k] = tea.String(v)
	}

	createReq := &dingtalk.CreateAndDeliverRequest{
		CardTemplateId: tea.String(req.CardTemplateID),
		OutTrackId:     tea.String(req.OutTrackID),
		CardData:       cardData,
		CallbackType:   tea.String("STREAM"),
		UserIdType:     tea.Int32(1),
		ImGroupOpenSpaceModel: &dingtalk.CreateAndDeliverRequestImGroupOpenSpaceModel{
			SupportForward: tea.Bool(true),
		},
		ImRobotOpenSpaceModel: &dingtalk.CreateAndDeliverRequestImRobotOpenSpaceModel{
			SupportForward: tea.Bool(true),
		},
	}

	if req.OpenSpaceID != "" {
		createReq.SetOpenSpaceId(req.OpenSpaceID)
	}

	// Handle different conversation types with appropriate delivery models
	switch req.ConversationType {
	case "2": // Group chat
		if req.RobotCode != "" {
			createReq.SetImGroupOpenDeliverModel(
				&dingtalk.CreateAndDeliverRequestImGroupOpenDeliverModel{
					RobotCode: tea.String(req.RobotCode),
				})
		}
	case "1": // Private chat with robot
		// For private chat, use ImRobotOpenDeliverModel with SpaceType
		createReq.SetImRobotOpenDeliverModel(
			&dingtalk.CreateAndDeliverRequestImRobotOpenDeliverModel{
				SpaceType: tea.String("IM_GROUP"),
			})
	default:
		// Fallback to group model if conversation type is unknown
		if req.RobotCode != "" {
			createReq.SetImGroupOpenDeliverModel(
				&dingtalk.CreateAndDeliverRequestImGroupOpenDeliverModel{
					RobotCode: tea.String(req.RobotCode),
				})
		}
	}

	_, err := s.client.CreateAndDeliverWithOptions(createReq, headers, &util.RuntimeOptions{})
	return err
}

// StreamingUpdateRequest 流式更新请求
type StreamingUpdateRequest struct {
	OutTrackID string
	Key        string
	Content    string
	IsFull     bool
	IsFinalize bool
}

// StreamingUpdate 流式更新卡片内容
func (s *StreamCardClient) StreamingUpdate(accessToken string, req *StreamingUpdateRequest) error {
	headers := &dingtalk.StreamingUpdateHeaders{
		XAcsDingtalkAccessToken: tea.String(accessToken),
	}

	updateReq := &dingtalk.StreamingUpdateRequest{
		OutTrackId: tea.String(req.OutTrackID),
		Guid:       tea.String(uuid.New().String()),
		Key:        tea.String(req.Key),
		Content:    tea.String(req.Content),
		IsFull:     tea.Bool(req.IsFull),
		IsFinalize: tea.Bool(req.IsFinalize),
		IsError:    tea.Bool(false),
	}

	_, err := s.client.StreamingUpdateWithOptions(updateReq, headers, &util.RuntimeOptions{})
	return err
}

// UpdateAIStreamCard 更新AI流式卡片 (简化版本,不依赖卡片模板)
// 这个方法需要与 client 包集成，这里提供一个独立实现
func UpdateAIStreamCard(accessToken, trackID, content string, isFinalize bool) error {
	cardClient, err := NewStreamCardClient()
	if err != nil {
		return fmt.Errorf("failed to create stream card client: %w", err)
	}

	req := &StreamingUpdateRequest{
		OutTrackID: trackID,
		Key:        "content",
		Content:    content,
		IsFull:     true,
		IsFinalize: isFinalize,
	}

	return cardClient.StreamingUpdate(accessToken, req)
}
