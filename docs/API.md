# 钉钉 SDK API 文档

## 目录

- [客户端管理](#客户端管理)
- [认证相关](#认证相关)
- [群聊管理](#群聊管理)
- [消息发送](#消息发送)
- [媒体上传](#媒体上传)

## 客户端管理

### NewDingTalkClient

创建钉钉客户端实例。

```go
func NewDingTalkClient(credential Credential) *DingTalkClient
```

**参数:**
- `credential` - 钉钉应用凭证

**返回:**
- `*DingTalkClient` - 钉钉客户端实例

**示例:**

```go
credential := client.Credential{
    ClientID:     "dingxxxxxxxx",
    ClientSecret: "your-secret-key",
}
dingClient := client.NewDingTalkClient(credential)
```

## 认证相关

### GetAccessToken

获取 AccessToken，自动处理缓存和刷新。

```go
func (c *DingTalkClient) GetAccessToken() (string, error)
```

**返回:**
- `string` - AccessToken
- `error` - 错误信息

**特性:**
- ✅ 自动缓存 Token
- ✅ 过期自动刷新
- ✅ 线程安全
- ✅ 预留 60 秒有效期避免临界点错误

**示例:**

```go
token, err := dingClient.GetAccessToken()
if err != nil {
    log.Fatal(err)
}
fmt.Println("Token:", token)
```

## 群聊管理

### GetChatList

获取所有可访问的群聊列表。

```go
func (c *DingTalkClient) GetChatList() (*ChatListResult, error)
```

**返回:**
- `*ChatListResult` - 群聊列表结果
- `error` - 错误信息

**ChatListResult 结构:**

```go
type ChatListResult struct {
    ErrorCode    int        `json:"errcode"`
    ErrorMessage string     `json:"errmsg"`
    ChatList     []ChatInfo `json:"chatlist"`
}
```

**ChatInfo 结构:**

```go
type ChatInfo struct {
    ChatID          string   `json:"chatid"`          // 群聊 ID
    Name            string   `json:"name"`            // 群名称
    Owner           string   `json:"owner"`           // 群主 userId
    UseridList      []string `json:"useridlist"`      // 群成员列表
    Icon            string   `json:"icon"`            // 群头像
    ConversationTag int      `json:"conversationtag"` // 0=单聊，1=群聊，2=企业群
}
```

**示例:**

```go
chatList, err := dingClient.GetChatList()
if err != nil {
    log.Fatal(err)
}

for i, chat := range chatList.ChatList {
    fmt.Printf("%d. %s (ID: %s, 成员: %d人)\n", 
        i+1, chat.Name, chat.ChatID, len(chat.UseridList))
}
```

**权限要求:**
- 应用需要开通群管理权限
- 只能获取机器人已加入的群聊

## 消息发送

### SendRobotMessage

发送机器人消息到指定群聊。

```go
func (c *DingTalkClient) SendRobotMessage(chatID string, message interface{}) error
```

**参数:**
- `chatID` - 群聊 ID（通过 GetChatList 获取）
- `message` - 消息内容（支持多种格式）

**支持的消息类型:**

#### 1. 文本消息

```go
textMsg := map[string]interface{}{
    "msgtype": "text",
    "text": map[string]string{
        "content": "Hello, World!",
    },
}
dingClient.SendRobotMessage(chatID, textMsg)
```

#### 2. Markdown 消息

```go
markdownMsg := map[string]interface{}{
    "msgtype": "markdown",
    "markdown": map[string]string{
        "title": "通知标题",
        "text": `### 📢 重要通知

**状态**: 🟢 正常

| 指标 | 数值 |
|------|------|
| CPU  | 45%  |
| 内存 | 78%  |

> 更新时间: 2026-02-07
`,
    },
}
dingClient.SendRobotMessage(chatID, markdownMsg)
```

#### 3. 链接消息

```go
linkMsg := map[string]interface{}{
    "msgtype": "link",
    "link": map[string]string{
        "title":      "钉钉开放平台",
        "text":       "查看更多开发文档和 API 接口",
        "messageUrl": "https://open.dingtalk.com",
        "picUrl":     "https://example.com/image.png",
    },
}
dingClient.SendRobotMessage(chatID, linkMsg)
```

#### 4. ActionCard 消息

```go
actionCardMsg := map[string]interface{}{
    "msgtype": "actionCard",
    "actionCard": map[string]interface{}{
        "title": "乔布斯的演讲",
        "text": "![screenshot](https://gw.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png) \n\n ### 乔布斯的演讲 \n\n Stay Hungry, Stay Foolish",
        "btnOrientation": "0",
        "singleTitle": "阅读全文",
        "singleURL": "https://www.dingtalk.com/",
    },
}
dingClient.SendRobotMessage(chatID, actionCardMsg)
```

**示例:**

```go
// 先获取群聊列表
chatList, _ := dingClient.GetChatList()
if len(chatList.ChatList) > 0 {
    chatID := chatList.ChatList[0].ChatID
    
    // 发送消息
    msg := map[string]interface{}{
        "msgtype": "text",
        "text": map[string]string{
            "content": "这是一条测试消息",
        },
    }
    
    err := dingClient.SendRobotMessage(chatID, msg)
    if err != nil {
        log.Fatal(err)
    }
}
```

**权限要求:**
- 机器人必须已加入目标群聊
- 应用需要有发送消息权限

## 媒体上传

### UploadMedia

上传媒体文件（图片、语音、视频、文件）。

```go
func (c *DingTalkClient) UploadMedia(content []byte, filename, mediaType, mimeType string) (*MediaUploadResult, error)
```

**参数:**
- `content` - 文件内容（字节数组）
- `filename` - 文件名
- `mediaType` - 媒体类型（image/voice/video/file）
- `mimeType` - MIME 类型（如 image/png）

**媒体类型常量:**

```go
const (
    MediaTypeImage string = "image"  // 图片
    MediaTypeVoice string = "voice"  // 语音
    MediaTypeVideo string = "video"  // 视频
    MediaTypeFile  string = "file"   // 文件
)

const (
    MimeTypeImagePng string = "image/png"  // PNG 图片
)
```

**MediaUploadResult 结构:**

```go
type MediaUploadResult struct {
    ErrorCode    int64  `json:"errcode"`
    ErrorMessage string `json:"errmsg"`
    MediaID      string `json:"media_id"`  // 媒体文件 ID
    CreatedAt    int64  `json:"created_at"`
    Type         string `json:"type"`
}
```

**示例:**

```go
import "os"

// 读取图片文件
imageData, err := os.ReadFile("example.png")
if err != nil {
    log.Fatal(err)
}

// 上传图片
result, err := dingClient.UploadMedia(
    imageData,
    "example.png",
    client.MediaTypeImage,
    client.MimeTypeImagePng,
)

if err != nil {
    log.Fatal(err)
}

fmt.Printf("上传成功! Media ID: %s\n", result.MediaID)
```

**限制:**
- 图片大小不超过 2MB
- 支持的图片格式：JPG、PNG
- 视频大小不超过 10MB

## 错误处理

所有 API 调用都会返回 error，建议进行错误检查：

```go
chatList, err := dingClient.GetChatList()
if err != nil {
    log.Printf("获取群聊列表失败: %v", err)
    // 根据错误类型进行处理
    return
}
```

常见错误：

| 错误码 | 说明 | 解决方案 |
|-------|------|---------|
| 40014 | 不合法的 access_token | 检查 ClientID 和 ClientSecret |
| 60011 | 设置已被管理员禁用 | 联系管理员开通权限 |
| 60020 | 机器人不在群里 | 将机器人添加到群聊 |

## 完整示例

查看 [examples](../examples) 目录获取更多完整示例：

- `examples/basic` - 基础使用示例
- `examples/get_chat_list` - 获取群聊列表
- `examples/send_message` - 发送各种类型消息
- `examples/message` - 消息处理示例
- `examples/stream` - Stream 模式示例
