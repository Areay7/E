package handlers

import (
	"net/http"
	"time"

	"cross-border-admin/models"

	"github.com/gin-gonic/gin"
)

type DashboardSummary struct {
	// 核心指标
	TodayOrders     int64              `json:"todayOrders"`
	TodaySales      float64            `json:"todaySales"`
	PendingShipment int64              `json:"pendingShipment"`
	OpenOrders      int64              `json:"openOrders"`
	LowStockCount   int64              `json:"lowStockCount"`

	// 对比数据
	YesterdayOrders int64              `json:"yesterdayOrders"`
	YesterdaySales  float64            `json:"yesterdaySales"`

	// 平台统计
	PlatformStats   []PlatformStat     `json:"platformStats"`

	// 趋势数据
	SalesTrend      []DailySales       `json:"salesTrend"`
	OrderTrend      []DailyOrders      `json:"orderTrend"`

	// 排行榜
	TopProducts     []ProductSales     `json:"topProducts"`
	TopCountries    []CountrySales     `json:"topCountries"`

	// 状态分布
	OrderStatusDist []OrderStatusCount `json:"orderStatusDist"`

	// 实时数据
	RecentOrders    []RecentOrder      `json:"recentOrders"`

	// 系统健康
	SystemHealth    SystemHealth       `json:"systemHealth"`
}

type PlatformStat struct {
	Platform       string  `json:"platform"`
	OrderCount     int64   `json:"orderCount"`
	TotalAmount    float64 `json:"totalAmount"`
	PendingCount   int64   `json:"pendingCount"`
	ShippedCount   int64   `json:"shippedCount"`
	CompletedCount int64   `json:"completedCount"`
	GrowthRate     float64 `json:"growthRate"`
}

type DailySales struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Count  int64   `json:"count"`
}

type DailyOrders struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type ProductSales struct {
	ProductName string  `json:"productName"`
	Platform    string  `json:"platform"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"`
}

type CountrySales struct {
	Country string  `json:"country"`
	Count   int64   `json:"count"`
	Amount  float64 `json:"amount"`
}

type OrderStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type RecentOrder struct {
	OrderID     string    `json:"orderId"`
	Platform    string    `json:"platform"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Country     string    `json:"country"`
	OrderTime   time.Time `json:"orderTime"`
}

type SystemHealth struct {
	TotalProducts   int64   `json:"totalProducts"`
	TotalInventory  int64   `json:"totalInventory"`
	APISuccessRate  float64 `json:"apiSuccessRate"`
	LastSyncTime    string  `json:"lastSyncTime"`
}

func GetDashboardSummary(c *gin.Context) {
	var summary DashboardSummary

	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	sevenDaysAgo := today.AddDate(0, 0, -7)
	thirtyDaysAgo := today.AddDate(0, 0, -30)

	// 今日订单数
	models.DB.Model(&models.Order{}).
		Where("order_time >= ?", today).
		Count(&summary.TodayOrders)

	// 昨日订单数
	models.DB.Model(&models.Order{}).
		Where("order_time >= ? AND order_time < ?", yesterday, today).
		Count(&summary.YesterdayOrders)

	// 今日销售额
	models.DB.Model(&models.Order{}).
		Where("order_time >= ?", today).
		Select("COALESCE(SUM(actual_amount), 0)").
		Scan(&summary.TodaySales)

	// 昨日销售额
	models.DB.Model(&models.Order{}).
		Where("order_time >= ? AND order_time < ?", yesterday, today).
		Select("COALESCE(SUM(actual_amount), 0)").
		Scan(&summary.YesterdaySales)

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

	// 各平台统计（近30天）
	type PlatformStatRaw struct {
		Platform       string
		OrderCount     int64
		TotalAmount    float64
		PendingCount   int64
		ShippedCount   int64
		CompletedCount int64
	}
	var platformStatsRaw []PlatformStatRaw
	models.DB.Model(&models.Order{}).
		Select(`
			platform,
			COUNT(*) as order_count,
			COALESCE(SUM(actual_amount), 0) as total_amount,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending_count,
			SUM(CASE WHEN shipping_status = 'shipped' THEN 1 ELSE 0 END) as shipped_count,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count
		`).
		Where("order_time >= ?", thirtyDaysAgo).
		Group("platform").
		Scan(&platformStatsRaw)

	summary.PlatformStats = make([]PlatformStat, len(platformStatsRaw))
	for i, raw := range platformStatsRaw {
		summary.PlatformStats[i] = PlatformStat{
			Platform:       raw.Platform,
			OrderCount:     raw.OrderCount,
			TotalAmount:    raw.TotalAmount,
			PendingCount:   raw.PendingCount,
			ShippedCount:   raw.ShippedCount,
			CompletedCount: raw.CompletedCount,
			GrowthRate:     0, // 可以后续计算增长率
		}
	}

	// 近7天销售趋势
	models.DB.Model(&models.Order{}).
		Select("DATE(order_time) as date, COALESCE(SUM(actual_amount), 0) as amount, COUNT(*) as count").
		Where("order_time >= ?", sevenDaysAgo).
		Group("DATE(order_time)").
		Order("date ASC").
		Scan(&summary.SalesTrend)

	// 近7天订单趋势
	models.DB.Model(&models.Order{}).
		Select("DATE(order_time) as date, COUNT(*) as count").
		Where("order_time >= ?", sevenDaysAgo).
		Group("DATE(order_time)").
		Order("date ASC").
		Scan(&summary.OrderTrend)

	// 热销商品TOP10（近30天）
	models.DB.Table("order_items").
		Select("product_name, platform, SUM(quantity) as quantity, SUM(total_price) as amount").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.order_time >= ?", thirtyDaysAgo).
		Group("product_name, platform").
		Order("quantity DESC").
		Limit(10).
		Scan(&summary.TopProducts)

	// 热门国家TOP10（近30天）
	models.DB.Model(&models.Order{}).
		Select("country, COUNT(*) as count, COALESCE(SUM(actual_amount), 0) as amount").
		Where("order_time >= ? AND country != ''", thirtyDaysAgo).
		Group("country").
		Order("count DESC").
		Limit(10).
		Scan(&summary.TopCountries)

	// 订单状态分布
	models.DB.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&summary.OrderStatusDist)

	// 最近10条订单
	var recentOrders []models.Order
	models.DB.Model(&models.Order{}).
		Select("order_id, platform, actual_amount, status, country, order_time").
		Order("order_time DESC").
		Limit(10).
		Find(&recentOrders)

	summary.RecentOrders = make([]RecentOrder, len(recentOrders))
	for i, order := range recentOrders {
		summary.RecentOrders[i] = RecentOrder{
			OrderID:   order.OrderID,
			Platform:  order.Platform,
			Amount:    order.ActualAmount,
			Status:    order.Status,
			Country:   order.Country,
			OrderTime: order.OrderTime,
		}
	}

	// 系统健康状态
	var totalProducts, totalInventory int64
	models.DB.Model(&models.Product{}).Count(&totalProducts)
	models.DB.Model(&models.Inventory{}).
		Select("COALESCE(SUM(stock), 0)").
		Scan(&totalInventory)

	var apiStats struct {
		Total   int64
		Success int64
	}
	models.DB.Model(&models.APILog{}).
		Where("created_at >= ?", today.AddDate(0, 0, -1)).
		Count(&apiStats.Total)
	models.DB.Model(&models.APILog{}).
		Where("created_at >= ? AND success = ?", today.AddDate(0, 0, -1), true).
		Count(&apiStats.Success)

	successRate := 0.0
	if apiStats.Total > 0 {
		successRate = float64(apiStats.Success) / float64(apiStats.Total) * 100
	}

	var lastSync models.SyncTask
	models.DB.Model(&models.SyncTask{}).
		Order("created_at DESC").
		First(&lastSync)

	summary.SystemHealth = SystemHealth{
		TotalProducts:  totalProducts,
		TotalInventory: totalInventory,
		APISuccessRate: successRate,
		LastSyncTime:   lastSync.CreatedAt.Format("2006-01-02 15:04:05"),
	}

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
