package telegram

import (
	"testing"
)

func TestParseInitData(t *testing.T) {
	// 测试用例：模拟 Telegram Web App 初始化数据
	initData := `query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A99281932%2C%22first_name%22%3A%22Andrew%22%2C%22last_name%22%3A%22Rogue%22%2C%22username%22%3A%22rogue%22%2C%22language_code%22%3A%22en%22%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2`

	data, err := ParseInitData(initData)
	if err != nil {
		t.Fatalf("ParseInitData failed: %v", err)
	}

	if data.User == nil {
		t.Fatal("User is nil")
	}

	if data.User.ID != 99281932 {
		t.Errorf("Expected user ID 99281932, got %d", data.User.ID)
	}

	if data.User.FirstName != "Andrew" {
		t.Errorf("Expected first name 'Andrew', got '%s'", data.User.FirstName)
	}

	if data.User.Username != "rogue" {
		t.Errorf("Expected username 'rogue', got '%s'", data.User.Username)
	}
}

func TestGetUserID(t *testing.T) {
	// 测试用例：从初始化数据中获取用户 ID
	initData := `user=%7B%22id%22%3A123456789%2C%22first_name%22%3A%22Test%22%7D`

	userID, err := GetUserID(initData)
	if err != nil {
		t.Fatalf("GetUserID failed: %v", err)
	}

	if userID != 123456789 {
		t.Errorf("Expected user ID 123456789, got %d", userID)
	}
}

func TestGetUserID_NoUser(t *testing.T) {
	// 测试用例：没有用户信息
	initData := `query_id=test123`

	userID, err := GetUserID(initData)
	if err != nil {
		t.Fatalf("GetUserID failed: %v", err)
	}

	if userID != 0 {
		t.Errorf("Expected user ID 0, got %d", userID)
	}
}

func TestGetUserID_InvalidData(t *testing.T) {
	// 测试用例：无效的初始化数据
	initData := `invalid%%%data`

	_, err := GetUserID(initData)
	if err == nil {
		t.Error("Expected error for invalid data, got nil")
	}
}
