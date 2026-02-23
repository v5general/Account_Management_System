package main

import (
	"account-management-system/backend/config"
	"account-management-system/backend/database"
	"account-management-system/backend/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.CloseDB()

	// 确保上传目录存在
	if err := os.MkdirAll(cfg.OSS.UploadPath, 0755); err != nil {
		log.Fatalf("创建上传目录失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 设置静态文件目录
	r.Static("/uploads", cfg.OSS.UploadPath)

	// 配置路由
	routes.SetupRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
