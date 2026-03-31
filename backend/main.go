package main

import (
	"fmt"
	"log"
	"net/http"

	"cross-border-admin/config"
	"cross-border-admin/models"
	"cross-border-admin/platform"
	"cross-border-admin/platform/aliexpress"
	"cross-border-admin/platform/shopee"
	"cross-border-admin/platform/tiktok"
	"cross-border-admin/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("配置文件加载失败: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 注册平台客户端
	platform.Register("shopee", shopee.NewShopeeClient())
	platform.Register("aliexpress", aliexpress.NewAliExpressClient())
	platform.Register("tiktok", tiktok.NewTikTokClient())

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建 Gin 引擎
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.CORS.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: config.AppConfig.CORS.AllowCredentials,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务器
	port := config.AppConfig.Server.Port
	log.Printf("服务器启动在 :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
