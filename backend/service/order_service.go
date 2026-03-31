package service

import (
	"context"
	"time"

	"cross-border-admin/models"
	"cross-border-admin/platform"
)

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

// SyncOrders 同步订单
func (s *OrderService) SyncOrders(ctx context.Context, platformName string, startTime, endTime time.Time) error {
	// 获取平台实例
	p, err := platform.Get(platformName)
	if err != nil {
		return err
	}

	// 创建同步任务
	task := &models.SyncTask{
		Platform: platformName,
		TaskType: "order",
		Status:   "running",
	}
	now := time.Now()
	task.StartTime = &now
	if err := models.DB.Create(task).Error; err != nil {
		return err
	}

	// 获取订单列表
	resp, err := p.Order().GetOrdersByTimeRange(ctx, startTime, endTime)
	if err != nil {
		task.Status = "failed"
		task.ErrorMessage = err.Error()
		endTime := time.Now()
		task.EndTime = &endTime
		models.DB.Save(task)
		return err
	}

	// 保存订单
	task.TotalCount = len(resp.Orders)
	for _, orderDetail := range resp.Orders {
		if err := s.saveOrder(platformName, &orderDetail); err != nil {
			task.FailCount++
		} else {
			task.SuccessCount++
		}
	}

	// 更新任务状态
	task.Status = "completed"
	endTime = time.Now()
	task.EndTime = &endTime
	models.DB.Save(task)

	return nil
}

// saveOrder 保存订单到数据库
func (s *OrderService) saveOrder(platformName string, detail *platform.OrderDetail) error {
	order := &models.Order{
		OrderID:         detail.OrderID,
		Platform:        platformName,
		OrderSN:         detail.OrderSN,
		Status:          detail.Status,
		PaymentStatus:   detail.PaymentStatus,
		ShippingStatus:  detail.ShippingStatus,
		BuyerUsername:   detail.BuyerUsername,
		BuyerEmail:      detail.BuyerEmail,
		RecipientName:   detail.RecipientName,
		RecipientPhone:  detail.RecipientPhone,
		ShippingAddress: detail.ShippingAddress,
		Country:         detail.Country,
		TotalAmount:     detail.TotalAmount,
		Currency:        detail.Currency,
		ShippingFee:     detail.ShippingFee,
		DiscountAmount:  detail.DiscountAmount,
		ActualAmount:    detail.ActualAmount,
		TrackingNumber:  detail.TrackingNumber,
		ShippingCarrier: detail.ShippingCarrier,
		OrderTime:       detail.OrderTime,
		PaymentTime:     detail.PaymentTime,
		ShippedAt:       detail.ShippedAt,
	}

	// 检查订单是否已存在
	var existing models.Order
	result := models.DB.Where("platform = ? AND order_id = ?", platformName, detail.OrderID).First(&existing)

	if result.Error == nil {
		// 更新现有订单
		order.ID = existing.ID
		return models.DB.Save(order).Error
	}

	// 创建新订单
	if err := models.DB.Create(order).Error; err != nil {
		return err
	}

	// 保存订单明细
	for _, item := range detail.Items {
		orderItem := &models.OrderItem{
			OrderID:        order.ID,
			Platform:       platformName,
			ItemID:         item.ItemID,
			ProductID:      item.ProductID,
			VariationID:    item.VariationID,
			ProductName:    item.ProductName,
			VariationName:  item.VariationName,
			SKU:            item.SKU,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			TotalPrice:     item.TotalPrice,
			DiscountAmount: item.DiscountAmount,
			ImageURL:       item.ImageURL,
		}
		if err := models.DB.Create(orderItem).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetOrderList 获取订单列表
func (s *OrderService) GetOrderList(platformName string, page, pageSize int, keyword, status string) ([]models.Order, int64, error) {
	db := models.DB.Model(&models.Order{}).Where("platform = ?", platformName)

	if keyword != "" {
		db = db.Where("order_id LIKE ? OR buyer_username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	offset := (page - 1) * pageSize
	if err := db.Order("order_time DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrderDetail 获取订单详情
func (s *OrderService) GetOrderDetail(orderID uint) (*models.Order, []models.OrderItem, error) {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		return nil, nil, err
	}

	var items []models.OrderItem
	if err := models.DB.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, nil, err
	}

	return &order, items, nil
}

// ShipOrder 发货
func (s *OrderService) ShipOrder(ctx context.Context, orderID uint, trackingNumber, carrier string) error {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		return err
	}

	// 调用平台 API 发货
	p, err := platform.Get(order.Platform)
	if err != nil {
		return err
	}

	req := &platform.ShipOrderRequest{
		OrderID:         order.OrderID,
		TrackingNumber:  trackingNumber,
		ShippingCarrier: carrier,
		ShipTime:        time.Now(),
	}

	if err := p.Order().ShipOrder(ctx, req); err != nil {
		return err
	}

	// 更新本地订单状态
	now := time.Now()
	order.TrackingNumber = trackingNumber
	order.ShippingCarrier = carrier
	order.ShippedAt = &now
	order.ShippingStatus = "shipped"

	return models.DB.Save(&order).Error
}
