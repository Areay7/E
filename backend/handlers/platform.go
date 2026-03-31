package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"cross-border-admin/models"
	"cross-border-admin/service"

	"github.com/gin-gonic/gin"
)

type PlatformHandler struct {
	orderService     *service.OrderService
	productService   *service.ProductService
	inventoryService *service.InventoryService
}

func NewPlatformHandler() *PlatformHandler {
	return &PlatformHandler{
		orderService:     service.NewOrderService(),
		productService:   service.NewProductService(),
		inventoryService: service.NewInventoryService(),
	}
}

// GetPlatformConfig 获取平台配置
func (h *PlatformHandler) GetPlatformConfig(c *gin.Context) {
	platform := c.Param("platform")

	var config models.PlatformConfig
	if err := models.DB.Where("platform = ?", platform).First(&config).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "平台配置不存在"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdatePlatformConfig 更新平台配置
func (h *PlatformHandler) UpdatePlatformConfig(c *gin.Context) {
	platform := c.Param("platform")

	var req models.PlatformConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var config models.PlatformConfig
	result := models.DB.Where("platform = ?", platform).First(&config)

	if result.Error != nil {
		// 创建新配置
		req.Platform = platform
		if err := models.DB.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建配置失败"})
			return
		}
		c.JSON(http.StatusOK, req)
		return
	}

	// 更新配置
	config.ShopID = req.ShopID
	config.ShopName = req.ShopName
	config.Enabled = req.Enabled
	config.AppKey = req.AppKey
	if req.AppSecret != "" {
		config.AppSecret = req.AppSecret
	}
	config.PartnerID = req.PartnerID
	if req.PartnerKey != "" {
		config.PartnerKey = req.PartnerKey
	}
	config.APIURL = req.APIURL
	config.SyncEnabled = req.SyncEnabled

	if err := models.DB.Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GetAllPlatformConfigs 获取所有平台配置
func (h *PlatformHandler) GetAllPlatformConfigs(c *gin.Context) {
	var configs []models.PlatformConfig
	if err := models.DB.Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, configs)
}

// SyncOrders 同步订单
func (h *PlatformHandler) SyncOrders(c *gin.Context) {
	platform := c.Param("platform")

	// 异步执行同步任务
	go func() {
		ctx := context.Background()
		endTime := time.Now()
		startTime := endTime.Add(-24 * time.Hour)
		h.orderService.SyncOrders(ctx, platform, startTime, endTime)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "同步任务已启动"})
}

// SyncProducts 同步商品
func (h *PlatformHandler) SyncProducts(c *gin.Context) {
	platform := c.Param("platform")

	go func() {
		ctx := context.Background()
		h.productService.SyncProducts(ctx, platform)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "同步任务已启动"})
}

// GetAPILogs 获取API调用日志
func (h *PlatformHandler) GetAPILogs(c *gin.Context) {
	platform := c.Query("platform")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	db := models.DB.Model(&models.APILog{})
	if platform != "" {
		db = db.Where("platform = ?", platform)
	}

	var total int64
	db.Count(&total)

	var logs []models.APILog
	offset := (pageInt - 1) * pageSizeInt
	db.Order("created_at DESC").Offset(offset).Limit(pageSizeInt).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"list":  logs,
		"total": total,
	})
}

// GetSyncTasks 获取同步任务列表
func (h *PlatformHandler) GetSyncTasks(c *gin.Context) {
	platform := c.Query("platform")
	taskType := c.Query("taskType")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	db := models.DB.Model(&models.SyncTask{})
	if platform != "" {
		db = db.Where("platform = ?", platform)
	}
	if taskType != "" {
		db = db.Where("task_type = ?", taskType)
	}

	var total int64
	db.Count(&total)

	var tasks []models.SyncTask
	offset := (pageInt - 1) * pageSizeInt
	db.Order("created_at DESC").Offset(offset).Limit(pageSizeInt).Find(&tasks)

	c.JSON(http.StatusOK, gin.H{
		"list":  tasks,
		"total": total,
	})
}
