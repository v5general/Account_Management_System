package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GenerateID 生成唯一ID
func GenerateID(prefix string) string {
	b := make([]byte, 8)
	rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}

// GetTimestamp 获取当前时间戳（毫秒）
func GetTimestamp() int64 {
	return time.Now().UnixMilli()
}

// Response 统一响应格式
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) Response {
	return Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: GetTimestamp(),
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string) Response {
	return Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: GetTimestamp(),
	}
}
