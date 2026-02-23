package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"account-management-system/backend/config"
)

// UploadFile 上传文件
func UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 获取文件扩展名
	ext := filepath.Ext(file.Filename)
	fileType := "image"
	if ext == ".pdf" {
		fileType = "pdf"
	}

	// 生成文件名
	filename := fmt.Sprintf("%s/%s%s", config.AppConfig.OSS.UploadPath, GenerateID("file"), ext)

	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建目标文件
	dst, err := os.Create(filename)
	if err != nil {
		return "", "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", "", fmt.Errorf("保存文件失败: %w", err)
	}

	return filename, fileType, nil
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// GetFileType 获取文件类型
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "image"
	case ".pdf":
		return "pdf"
	default:
		return "unknown"
	}
}

// ValidateFileType 验证文件类型
func ValidateFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"}
	for _, allowed := range allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}
