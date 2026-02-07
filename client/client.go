package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	url2 "net/url"
	"sync"
	"time"
)

// OpenAPI doc: https://open.dingtalk.com/document/isvapp/upload-media-files
const (
	MediaTypeImage string = "image"
	MediaTypeVoice string = "voice"
	MediaTypeVideo string = "video"
	MediaTypeFile  string = "file"
)

const (
	MimeTypeImagePng string = "image/png"
)

// MediaUploadResult 媒体上传结果
type MediaUploadResult struct {
	ErrorCode    int64  `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
	MediaID      string `json:"media_id"`
	CreatedAt    int64  `json:"created_at"`
	Type         string `json:"type"`
}

// OAuthTokenResult OAuth Token 结果
type OAuthTokenResult struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// OpenConversationIdResult OpenConversationId 结果
type OpenConversationIdResult struct {
	OpenConversationId string `json:"openConversationId"`
}

// Credential 钉钉应用凭证
type Credential struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// DingTalkClientInterface 钉钉客户端接口
type DingTalkClientInterface interface {
	GetAccessToken() (string, error)
	UploadMedia(content []byte, filename, mediaType, mimeType string) (*MediaUploadResult, error)
	GetOpenConversationId(chatID string) (string, error)
}

// DingTalkClient 钉钉客户端
type DingTalkClient struct {
	Credential  Credential
	AccessToken string
	expireAt    int64
	mutex       sync.Mutex
}

// DingTalkClientManagerInterface 钉钉客户端管理器接口
type DingTalkClientManagerInterface interface {
	GetClientByOAuthClientID(clientId string) DingTalkClientInterface
}

// DingTalkClientManager 钉钉客户端管理器
type DingTalkClientManager struct {
	Credentials []Credential
	Clients     map[string]*DingTalkClient
	mutex       sync.Mutex
}

// NewDingTalkClient 创建钉钉客户端
func NewDingTalkClient(credential Credential) *DingTalkClient {
	return &DingTalkClient{
		Credential: credential,
	}
}

// NewDingTalkClientManager 创建钉钉客户端管理器
func NewDingTalkClientManager(credentials []Credential) *DingTalkClientManager {
	clients := make(map[string]*DingTalkClient)

	if credentials != nil {
		for _, credential := range credentials {
			clients[credential.ClientID] = NewDingTalkClient(credential)
		}
	}
	return &DingTalkClientManager{
		Credentials: credentials,
		Clients:     clients,
	}
}

// GetClientByOAuthClientID 通过 OAuth ClientID 获取客户端
func (m *DingTalkClientManager) GetClientByOAuthClientID(clientId string) DingTalkClientInterface {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if client, ok := m.Clients[clientId]; ok {
		return client
	}
	return nil
}

// GetAccessToken 获取 AccessToken（自动缓存）
func (c *DingTalkClient) GetAccessToken() (string, error) {
	accessToken := ""
	{
		// 先查询缓存
		c.mutex.Lock()
		now := time.Now().Unix()
		if c.expireAt > 0 && c.AccessToken != "" && (now+60) < c.expireAt {
			// 预留一分钟有效期避免在Token过期的临界点调用接口出现401错误
			accessToken = c.AccessToken
		}
		c.mutex.Unlock()
	}
	if accessToken != "" {
		return accessToken, nil
	}

	tokenResult, err := c.getAccessTokenFromDingTalk()
	if err != nil {
		return "", err
	}

	{
		// 更新缓存
		c.mutex.Lock()
		c.AccessToken = tokenResult.AccessToken
		c.expireAt = time.Now().Unix() + int64(tokenResult.ExpiresIn)
		c.mutex.Unlock()
	}
	return tokenResult.AccessToken, nil
}

// UploadMedia 上传媒体文件
func (c *DingTalkClient) UploadMedia(content []byte, filename, mediaType, mimeType string) (*MediaUploadResult, error) {
	// OpenAPI doc: https://open.dingtalk.com/document/isvapp/upload-media-files
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	if len(accessToken) == 0 {
		return nil, errors.New("empty access token")
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(content)
	if err != nil {
		return nil, err
	}
	if err = writer.WriteField("type", mediaType); err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request to upload the media file
	url := fmt.Sprintf("https://oapi.dingtalk.com/media/upload?access_token=%s", url2.QueryEscape(accessToken))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the HTTP request and parse the response
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse the response body as JSON and extract the media ID
	media := &MediaUploadResult{}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bodyBytes, media); err != nil {
		return nil, err
	}
	if media.ErrorCode != 0 {
		return nil, errors.New(media.ErrorMessage)
	}
	return media, nil
}

// SendRobotMessage 发送企业内部机器人消息
// 文档: https://open.dingtalk.com/document/orgapp/robot-sends-group-messages
func (c *DingTalkClient) SendRobotMessage(chatID string, message interface{}) error {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return err
	}

	// 构造请求参数
	params := map[string]interface{}{
		"chatId": chatID,
		"msg":    message,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://oapi.dingtalk.com/chat/send?access_token=%s", url2.QueryEscape(accessToken))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return err
	}

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("send robot message failed: %v", result["errmsg"])
	}
	return nil
}

// getAccessTokenFromDingTalk 从钉钉获取 AccessToken
func (c *DingTalkClient) getAccessTokenFromDingTalk() (*OAuthTokenResult, error) {
	// OpenAPI doc: https://open.dingtalk.com/document/orgapp/obtain-orgapp-token
	apiUrl := "https://oapi.dingtalk.com/gettoken"
	queryParams := url2.Values{}
	queryParams.Add("appkey", c.Credential.ClientID)
	queryParams.Add("appsecret", c.Credential.ClientSecret)

	// Create a new HTTP request to get the AccessToken
	req, err := http.NewRequest("GET", apiUrl+"?"+queryParams.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// Send the HTTP request and parse the response body as JSON
	client := http.Client{
		Timeout: time.Second * 60,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	tokenResult := &OAuthTokenResult{}
	err = json.Unmarshal(body, tokenResult)
	if err != nil {
		return nil, err
	}
	if tokenResult.ErrorCode != 0 {
		return nil, errors.New(tokenResult.ErrorMessage)
	}
	return tokenResult, nil
}

// GetOpenConversationId 通过 chatId 获取 OpenConversationId
// 文档: https://open.dingtalk.com/document/development/obtain-group-openconversationid
func (c *DingTalkClient) GetOpenConversationId(chatID string) (string, error) {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return "", err
	}

	// 使用新版 API 地址
	url := fmt.Sprintf("https://api.dingtalk.com/v1.0/im/chat/%s/convertToOpenConversationId", url2.PathEscape(chatID))
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	// 新版 API 使用 Header 传递 token
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// 检查 HTTP 状态码
	if res.StatusCode != 200 {
		return "", fmt.Errorf("API request failed with status %d: %s", res.StatusCode, string(bodyBytes))
	}

	result := &OpenConversationIdResult{}
	if err = json.Unmarshal(bodyBytes, result); err != nil {
		return "", err
	}

	return result.OpenConversationId, nil
}

// SendWebhookMessage 通过 Webhook URL 发送消息（自定义机器人）
// webhookURL: 完整的 webhook 地址，例如: https://oapi.dingtalk.com/robot/send?access_token=xxx
// message: 消息内容，支持 text/markdown/link 等格式
func SendWebhookMessage(webhookURL string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bodyBytes, &result); err != nil {
		return err
	}

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("send webhook message failed: %v", result["errmsg"])
	}
	return nil
}
