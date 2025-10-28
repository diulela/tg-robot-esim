package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// TelegramWebAppMiddleware Telegram Web App 身份验证中间件
func TelegramWebAppMiddleware(botToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 获取初始化数据
			initData := r.Header.Get("X-Telegram-Init-Data")
			if initData == "" {
				// 开发模式：从查询参数获取
				initData = r.URL.Query().Get("init_data")
			}

			// 如果没有初始化数据，允许通过（开发模式）
			if initData == "" {
				next.ServeHTTP(w, r)
				return
			}

			// 验证初始化数据
			if !validateTelegramWebAppData(initData, botToken) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// validateTelegramWebAppData 验证 Telegram Web App 初始化数据
func validateTelegramWebAppData(initData string, botToken string) bool {
	// 解析初始化数据
	values, err := url.ParseQuery(initData)
	if err != nil {
		return false
	}

	// 获取 hash
	hash := values.Get("hash")
	if hash == "" {
		return false
	}

	// 移除 hash 参数
	values.Del("hash")

	// 构建数据检查字符串
	var dataCheckArr []string
	for key, vals := range values {
		if len(vals) > 0 {
			dataCheckArr = append(dataCheckArr, key+"="+vals[0])
		}
	}
	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// 计算密钥
	secretKey := sha256.Sum256([]byte(botToken))

	// 计算 HMAC
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// 比较 hash
	return calculatedHash == hash
}

// CORSMiddleware CORS 中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Telegram-Init-Data")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
