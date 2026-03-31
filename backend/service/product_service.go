package service

import (
	"context"
	"time"

	"cross-border-admin/models"
	"cross-border-admin/platform"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

// SyncProducts 同步商品
func (s *ProductService) SyncProducts(ctx context.Context, platformName string) error {
	p, err := platform.Get(platformName)
	if err != nil {
		return err
	}

	task := &models.SyncTask{
		Platform: platformName,
		TaskType: "product",
		Status:   "running",
	}
	now := time.Now()
	task.StartTime = &now
	models.DB.Create(task)

	page := 1
	pageSize := 100

	for {
		req := &platform.ProductListRequest{
			Page:     page,
			PageSize: pageSize,
		}

		resp, err := p.Product().GetProductList(ctx, req)
		if err != nil {
			task.Status = "failed"
			task.ErrorMessage = err.Error()
			endTime := time.Now()
			task.EndTime = &endTime
			models.DB.Save(task)
			return err
		}

		task.TotalCount += len(resp.Products)
		for _, productDetail := range resp.Products {
			if err := s.saveProduct(platformName, &productDetail); err != nil {
				task.FailCount++
			} else {
				task.SuccessCount++
			}
		}

		if !resp.HasMore {
			break
		}
		page++
	}

	task.Status = "completed"
	endTime := time.Now()
	task.EndTime = &endTime
	models.DB.Save(task)

	return nil
}

// saveProduct 保存商品
func (s *ProductService) saveProduct(platformName string, detail *platform.ProductDetail) error {
	product := &models.Product{
		Platform:      platformName,
		ProductID:     detail.ProductID,
		Name:          detail.Name,
		Description:   detail.Description,
		Category:      detail.Category,
		Brand:         detail.Brand,
		Status:        detail.Status,
		Price:         detail.Price,
		OriginalPrice: detail.OriginalPrice,
		Currency:      detail.Currency,
		Stock:         detail.Stock,
		SKU:           detail.SKU,
		MainImage:     detail.MainImage,
		Weight:        detail.Weight,
		Length:        detail.Dimensions.Length,
		Width:         detail.Dimensions.Width,
		Height:        detail.Dimensions.Height,
		SoldCount:     detail.SoldCount,
		ViewCount:     detail.ViewCount,
		Rating:        detail.Rating,
		ReviewCount:   detail.ReviewCount,
	}

	var existing models.Product
	result := models.DB.Where("platform = ? AND product_id = ?", platformName, detail.ProductID).First(&existing)

	if result.Error == nil {
		product.ID = existing.ID
		return models.DB.Save(product).Error
	}

	return models.DB.Create(product).Error
}

// GetProductList 获取商品列表
func (s *ProductService) GetProductList(platformName string, page, pageSize int, keyword, status string) ([]models.Product, int64, error) {
	db := models.DB.Model(&models.Product{}).Where("platform = ?", platformName)

	if keyword != "" {
		db = db.Where("name LIKE ? OR sku LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.Product
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
