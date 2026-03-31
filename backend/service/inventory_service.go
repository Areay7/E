package service

import (
	"context"

	"cross-border-admin/models"
	"cross-border-admin/platform"
)

type InventoryService struct{}

func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

// UpdateInventory 更新库存
func (s *InventoryService) UpdateInventory(ctx context.Context, platformName, sku string, quantity int) error {
	// 更新平台库存
	p, err := platform.Get(platformName)
	if err != nil {
		return err
	}

	if err := p.Inventory().UpdateInventory(ctx, sku, quantity); err != nil {
		return err
	}

	// 更新本地库存
	var inventory models.Inventory
	result := models.DB.Where("platform = ? AND sku = ?", platformName, sku).First(&inventory)

	if result.Error != nil {
		// 创建新库存记录
		inventory = models.Inventory{
			Platform:       platformName,
			SKU:            sku,
			Stock:          quantity,
			AvailableStock: quantity,
		}
		return models.DB.Create(&inventory).Error
	}

	// 记录库存变动
	log := &models.InventoryLog{
		InventoryID: inventory.ID,
		Type:        "manual_update",
		Quantity:    quantity - inventory.Stock,
		BeforeStock: inventory.Stock,
		AfterStock:  quantity,
		Reason:      "手动更新",
	}
	models.DB.Create(log)

	// 更新库存
	inventory.Stock = quantity
	inventory.AvailableStock = quantity - inventory.ReservedStock
	return models.DB.Save(&inventory).Error
}

// GetInventory 获取库存
func (s *InventoryService) GetInventory(platformName, sku string) (*models.Inventory, error) {
	var inventory models.Inventory
	if err := models.DB.Where("platform = ? AND sku = ?", platformName, sku).First(&inventory).Error; err != nil {
		return nil, err
	}
	return &inventory, nil
}

// GetInventoryList 获取库存列表
func (s *InventoryService) GetInventoryList(platformName string, page, pageSize int) ([]models.Inventory, int64, error) {
	db := models.DB.Model(&models.Inventory{})

	if platformName != "" {
		db = db.Where("platform = ?", platformName)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var inventories []models.Inventory
	offset := (page - 1) * pageSize
	if err := db.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&inventories).Error; err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}
