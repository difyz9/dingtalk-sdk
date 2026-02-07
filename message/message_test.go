package message

import (
	"testing"
)

func TestGetSenderIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		msg      ReceiveMsg
		expected string
	}{
		{
			name: "With SenderStaffId",
			msg: ReceiveMsg{
				SenderStaffId: "staff123",
				SenderID:      "sender456",
			},
			expected: "staff123",
		},
		{
			name: "Without SenderStaffId",
			msg: ReceiveMsg{
				SenderStaffId: "",
				SenderID:      "sender456",
			},
			expected: "sender456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.msg.GetSenderIdentifier()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetChatTitle(t *testing.T) {
	tests := []struct {
		name     string
		msg      ReceiveMsg
		expected string
	}{
		{
			name: "Private chat",
			msg: ReceiveMsg{
				ConversationType:  "1",
				SenderNick:        "张三",
				ConversationTitle: "测试群",
			},
			expected: "张三_私聊",
		},
		{
			name: "Group chat",
			msg: ReceiveMsg{
				ConversationType:  "2",
				SenderNick:        "张三",
				ConversationTitle: "测试群",
			},
			expected: "测试群",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.msg.GetChatTitle()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
