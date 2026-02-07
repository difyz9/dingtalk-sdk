package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MsgType 消息类型
type MsgType string

const (
	TEXT     MsgType = "text"
	MARKDOWN MsgType = "markdown"
)

// ReceiveMsg 接收的消息体
type ReceiveMsg struct {
	ConversationID string `json:"conversationId"`
	AtUsers        []struct {
		DingtalkID string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatbotUserID             string  `json:"chatbotUserId"`
	MsgID                     string  `json:"msgId"`
	SenderNick                string  `json:"senderNick"`
	IsAdmin                   bool    `json:"isAdmin"`
	SenderStaffId             string  `json:"senderStaffId"`
	SessionWebhookExpiredTime int64   `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64   `json:"createAt"`
	ConversationType          string  `json:"conversationType"`
	SenderID                  string  `json:"senderId"`
	ConversationTitle         string  `json:"conversationTitle"`
	IsInAtList                bool    `json:"isInAtList"`
	SessionWebhook            string  `json:"sessionWebhook"`
	Text                      Text    `json:"text"`
	RobotCode                 string  `json:"robotCode"`
	Msgtype                   MsgType `json:"msgtype"`
}

// TextMessage 文本消息
type TextMessage struct {
	MsgType MsgType `json:"msgtype"`
	At      *At     `json:"at"`
	Text    *Text   `json:"text"`
}

// Text 消息内容
type Text struct {
	Content string `json:"content"`
}

// MarkDownMessage Markdown 消息
type MarkDownMessage struct {
	MsgType  MsgType   `json:"msgtype"`
	At       *At       `json:"at"`
	MarkDown *MarkDown `json:"markdown"`
}

// MarkDown Markdown 消息内容
type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// At @内容
type At struct {
	AtUserIds []string `json:"atUserIds"`
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

// GetSenderIdentifier 获取用户标识，兼容当 SenderStaffId 字段为空的场景
func (r ReceiveMsg) GetSenderIdentifier() (uid string) {
	if r.SenderStaffId != "" {
		uid = r.SenderStaffId
	} else {
		uid = r.SenderID
	}
	return uid
}

// GetChatTitle 获取聊天的群名字，如果是私聊，则命名为 昵称_私聊
func (r ReceiveMsg) GetChatTitle() (chatType string) {
	if r.ConversationType == "1" {
		chatType = r.SenderNick + "_私聊"
	} else {
		chatType = r.ConversationTitle
	}
	return chatType
}

// ReplyToDingtalk 发消息给钉钉
func (r ReceiveMsg) ReplyToDingtalk(msgType, msg string) (statuscode int, err error) {
	atUser := r.SenderStaffId
	if atUser == "" {
		msg = fmt.Sprintf("%s\n\n@%s", msg, r.SenderNick)
	}
	var msgtmp interface{}
	switch msgType {
	case string(TEXT):
		msgtmp = &TextMessage{Text: &Text{Content: msg}, MsgType: TEXT, At: &At{AtUserIds: []string{atUser}}}
	case string(MARKDOWN):
		if atUser != "" && r.ConversationType != "1" {
			msg = fmt.Sprintf("%s\n\n@%s", msg, atUser)
		}
		msgtmp = &MarkDownMessage{MsgType: MARKDOWN, At: &At{AtUserIds: []string{atUser}}, MarkDown: &MarkDown{Title: "Markdown Msg", Text: msg}}
	default:
		msgtmp = &TextMessage{Text: &Text{Content: msg}, MsgType: TEXT, At: &At{AtUserIds: []string{atUser}}}
	}

	data, err := json.Marshal(msgtmp)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", r.SessionWebhook, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
