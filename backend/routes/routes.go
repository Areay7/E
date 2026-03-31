package routes

import (
	"cross-border-admin/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 认证相关路由（无需token）
	authHandler := handlers.NewAuthHandler()
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/captcha", authHandler.GetCaptcha)
	}

	// 需要认证的API
	api := r.Group("/api/v1")
	api.Use(handlers.AuthMiddleware())
	{
		// 用户信息
		api.GET("/user/current", authHandler.GetCurrentUser)

		// 数据看板
		api.GET("/dashboard/summary", handlers.GetDashboardSummary)
		api.GET("/dashboard/sales-report", handlers.GetSalesReport)

		// 平台管理
		platformHandler := handlers.NewPlatformHandler()
		platform := api.Group("/platforms")
		{
			platform.GET("", platformHandler.GetAllPlatformConfigs)
			platform.GET("/:platform/config", platformHandler.GetPlatformConfig)
			platform.PUT("/:platform/config", platformHandler.UpdatePlatformConfig)
			platform.POST("/:platform/sync/orders", platformHandler.SyncOrders)
			platform.POST("/:platform/sync/products", platformHandler.SyncProducts)
		}

		// API 日志
		api.GET("/api-logs", platformHandler.GetAPILogs)

		// 同步任务
		api.GET("/sync-tasks", platformHandler.GetSyncTasks)

		// 订单管理
		orderHandler := handlers.NewOrderHandlerV2()
		orders := api.Group("/orders")
		{
			orders.GET("", orderHandler.GetOrderList)
			orders.GET("/:id", orderHandler.GetOrderDetail)
			orders.POST("/:id/ship", orderHandler.ShipOrder)
		}

		// 商品管理
		productHandler := handlers.NewProductHandler()
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetProductList)
			products.GET("/:id", productHandler.GetProductDetail)
		}

		// 库存管理
		inventoryHandler := handlers.NewInventoryHandler()
		inventory := api.Group("/inventory")
		{
			inventory.GET("", inventoryHandler.GetInventoryList)
			inventory.PUT("", inventoryHandler.UpdateInventory)
		}
	}
}
