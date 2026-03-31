package handlers

import (
	"net/http"
	"time"

	"cross-border-admin/models"

	"github.com/gin-gonic/gin"
)

type DashboardSummary struct {
	TodayOrders     int64              `json:"todayOrders"`
	TodaySales      float64            `json:"todaySales"`
	PendingShipment int64              `json:"pendingShipment"`
	OpenOrders      int64              `json:"openOrders"`
	LowStockCount   int64              `json:"lowStockCount"`
	PlatformStats   []PlatformStat     `json:"platformStats"`
	SalesTrend      []DailySales       `json:"salesTrend"`
	TopProducts     []ProductSales     `json:"topProducts"`
	OrderStatusDist []OrderStatusCount `json:"orderStatusDist"`
}

type PlatformStat struct {
	Platform    string  `json:"platform"`
	OrderCount  int64   `json:"orderCount"`
	TotalAmount float64 `json:"totalAmount"`
}

type DailySales struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Count  int64   `json:"count"`
}

type ProductSales struct {
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"`
}

type OrderStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

func GetDashboardSummary(c *gin.Context) {
	var summary DashboardSummary

	today := time.Now().Truncate(24 * time.Hour)
	sevenDaysAgo := today.AddDate(0, 0, -7)

	// 今日订单数
	models.DB.Model(&models.Order{}).
		Where("order_time >= ?", today).
		Count(&summary.TodayOrders)

	// 今日销售额
	models.DB.Model(&models.Order{}).
		Where("order_time >= ?", today).
		Select("COALESCE(SUM(actual_amount), 0)").
		Scan(&summary.TodaySales)

	// 待发货订单数
	models.DB.Model(&models.Order{}).
		Where("shipping_status = ?", "pending").
		Count(&summary.PendingShipment)

	// 进行中订单数
	models.DB.Model(&models.Order{}).
		Where("status IN ?", []string{"pending", "processing", "paid"}).
		Count(&summary.OpenOrders)

	// 低库存商品数（库存<10）
	models.DB.Model(&models.Inventory{}).
		Where("available_stock < ?", 10).
		Count(&summary.LowStockCount)

	// 各平台统计
	models.DB.Model(&models.Order{}).
		Select("platform, COUNT(*) as order_count, COALESCE(SUM(actual_amount), 0) as total_amount").
		Where("order_time >= ?", sevenDaysAgo).
		Group("platform").
		Scan(&summary.PlatformStats)

	// 近7天销售趋势
	models.DB.Model(&models.Order{}).
		Select("DATE(order_time) as date, COALESCE(SUM(actual_amount), 0) as amount, COUNT(*) as count").
		Where("order_time >= ?", sevenDaysAgo).
		Group("DATE(order_time)").
		Order("date ASC").
		Scan(&summary.SalesTrend)

	// 热销商品TOP10
	models.DB.Table("order_items").
		Select("product_name, SUM(quantity) as quantity, SUM(total_price) as amount").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.order_time >= ?", sevenDaysAgo).
		Group("product_name").
		Order("quantity DESC").
		Limit(10).
		Scan(&summary.TopProducts)

	// 订单状态分布
	models.DB.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&summary.OrderStatusDist)

	c.JSON(http.StatusOK, summary)
}

// GetSalesReport 获取销售报表
func GetSalesReport(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	platform := c.Query("platform")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供开始和结束日期"})
		return
	}

	db := models.DB.Model(&models.Order{}).
		Where("order_time BETWEEN ? AND ?", startDate, endDate)

	if platform != "" {
		db = db.Where("platform = ?", platform)
	}

	var report struct {
		TotalOrders  int64   `json:"totalOrders"`
		TotalAmount  float64 `json:"totalAmount"`
		AvgAmount    float64 `json:"avgAmount"`
		PaidOrders   int64   `json:"paidOrders"`
		ShippedCount int64   `json:"shippedCount"`
	}

	db.Count(&report.TotalOrders)
	db.Select("COALESCE(SUM(actual_amount), 0)").Scan(&report.TotalAmount)
	db.Select("COALESCE(AVG(actual_amount), 0)").Scan(&report.AvgAmount)
	db.Where("payment_status = ?", "paid").Count(&report.PaidOrders)
	db.Where("shipping_status = ?", "shipped").Count(&report.ShippedCount)

	c.JSON(http.StatusOK, report)
}
