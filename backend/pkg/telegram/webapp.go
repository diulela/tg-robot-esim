package telegram

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// WebAppUser Telegram Web App 用户信息
type WebAppUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
	IsPremium    bool   `json:"is_premium,omitempty"`
	PhotoURL     string `json:"photo_url,omitempty"`
}

// WebAppInitData Telegram Web App 初始化数据
type WebAppInitData struct {
	QueryID      string      `json:"query_id,omitempty"`
	User         *WebAppUser `json:"user,omitempty"`
	Receiver     *WebAppUser `json:"receiver,omitempty"`
	Chat         interface{} `json:"chat,omitempty"`
	ChatType     string      `json:"chat_type,omitempty"`
	ChatInstance string      `json:"chat_instance,omitempty"`
	StartParam   string      `json:"start_param,omitempty"`
	CanSendAfter int         `json:"can_send_after,omitempty"`
	AuthDate     int64       `json:"auth_date,omitempty"`
	Hash         string      `json:"hash,omitempty"`
}

// ParseInitData 解析 Telegram Web App 初始化数据
func ParseInitData(initData string) (*WebAppInitData, error) {

	fmt.Println("========解析tg web app data=======", initData)

	// 解析 URL 查询参数
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	data := &WebAppInitData{
		QueryID:      values.Get("query_id"),
		ChatType:     values.Get("chat_type"),
		ChatInstance: values.Get("chat_instance"),
		StartParam:   values.Get("start_param"),
		Hash:         values.Get("hash"),
	}

	// 解析用户信息
	if userStr := values.Get("user"); userStr != "" {
		var user WebAppUser
		if err := json.Unmarshal([]byte(userStr), &user); err == nil {
			data.User = &user
		}
	}

	// 解析接收者信息
	if receiverStr := values.Get("receiver"); receiverStr != "" {
		var receiver WebAppUser
		if err := json.Unmarshal([]byte(receiverStr), &receiver); err == nil {
			data.Receiver = &receiver
		}
	}

	// 解析聊天信息
	if chatStr := values.Get("chat"); chatStr != "" {
		var chat interface{}
		if err := json.Unmarshal([]byte(chatStr), &chat); err == nil {
			data.Chat = chat
		}
	}

	return data, nil
}

// GetUserID 从初始化数据中获取用户 ID
func GetUserID(initData string) (int64, error) {
	data, err := ParseInitData(initData)
	if err != nil {
		return 0, err
	}

	if data.User != nil {
		return data.User.ID, nil
	}

	return 0, nil
}
