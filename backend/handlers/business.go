package handlers

import (
	"context"
	"net/http"
	"strconv"

	"cross-border-admin/models"
	"cross-border-admin/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
	}
}

// GetProductList 获取商品列表
func (h *ProductHandler) GetProductList(c *gin.Context) {
	platform := c.Query("platform")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	products, total, err := h.productService.GetProductList(platform, page, pageSize, keyword, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  products,
		"total": total,
	})
}

// GetProductDetail 获取商品详情
func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var product models.Product
	if err := models.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	var variations []models.ProductVariation
	models.DB.Where("product_id = ?", id).Find(&variations)

	c.JSON(http.StatusOK, gin.H{
		"product":    product,
		"variations": variations,
	})
}

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{
		inventoryService: service.NewInventoryService(),
	}
}

// GetInventoryList 获取库存列表
func (h *InventoryHandler) GetInventoryList(c *gin.Context) {
	platform := c.Query("platform")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	inventories, total, err := h.inventoryService.GetInventoryList(platform, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  inventories,
		"total": total,
	})
}

// UpdateInventory 更新库存
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	var req struct {
		Platform string `json:"platform" binding:"required"`
		SKU      string `json:"sku" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := h.inventoryService.UpdateInventory(ctx, req.Platform, req.SKU, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

type OrderHandlerV2 struct {
	orderService *service.OrderService
}

func NewOrderHandlerV2() *OrderHandlerV2 {
	return &OrderHandlerV2{
		orderService: service.NewOrderService(),
	}
}

// GetOrderList 获取订单列表
func (h *OrderHandlerV2) GetOrderList(c *gin.Context) {
	platform := c.Query("platform")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	orders, total, err := h.orderService.GetOrderList(platform, page, pageSize, keyword, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  orders,
		"total": total,
	})
}

// GetOrderDetail 获取订单详情
func (h *OrderHandlerV2) GetOrderDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	order, items, err := h.orderService.GetOrderDetail(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"items": items,
	})
}

// ShipOrder 发货
func (h *OrderHandlerV2) ShipOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		TrackingNumber string `json:"trackingNumber" binding:"required"`
		Carrier        string `json:"carrier" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := h.orderService.ShipOrder(ctx, uint(id), req.TrackingNumber, req.Carrier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发货失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发货成功"})
}
