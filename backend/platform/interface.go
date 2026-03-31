package platform

import (
	"context"
	"time"
)

// Platform 平台接口定义
type Platform interface {
	// GetName 获取平台名称
	GetName() string

	// Auth 认证相关
	Auth() AuthAPI

	// Order 订单相关
	Order() OrderAPI

	// Product 商品相关
	Product() ProductAPI

	// Inventory 库存相关
	Inventory() InventoryAPI

	// Logistics 物流相关
	Logistics() LogisticsAPI
}

// AuthAPI 认证接口
type AuthAPI interface {
	// GetAccessToken 获取访问令牌
	GetAccessToken(ctx context.Context) (string, error)

	// RefreshToken 刷新令牌
	RefreshToken(ctx context.Context, refreshToken string) (string, error)

	// ValidateToken 验证令牌
	ValidateToken(ctx context.Context, token string) (bool, error)
}

// OrderAPI 订单接口
type OrderAPI interface {
	// GetOrderList 获取订单列表
	GetOrderList(ctx context.Context, req *OrderListRequest) (*OrderListResponse, error)

	// GetOrderDetail 获取订单详情
	GetOrderDetail(ctx context.Context, orderID string) (*OrderDetail, error)

	// GetOrdersByTimeRange 按时间范围获取订单
	GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time) (*OrderListResponse, error)

	// CancelOrder 取消订单
	CancelOrder(ctx context.Context, orderID string, reason string) error

	// ShipOrder 发货
	ShipOrder(ctx context.Context, req *ShipOrderRequest) error
}

// ProductAPI 商品接口
type ProductAPI interface {
	// GetProductList 获取商品列表
	GetProductList(ctx context.Context, req *ProductListRequest) (*ProductListResponse, error)

	// GetProductDetail 获取商品详情
	GetProductDetail(ctx context.Context, productID string) (*ProductDetail, error)

	// CreateProduct 创建商品
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductDetail, error)

	// UpdateProduct 更新商品
	UpdateProduct(ctx context.Context, productID string, req *UpdateProductRequest) error

	// DeleteProduct 删除商品
	DeleteProduct(ctx context.Context, productID string) error

	// UpdateProductStatus 更新商品状态
	UpdateProductStatus(ctx context.Context, productID string, status string) error
}

// InventoryAPI 库存接口
type InventoryAPI interface {
	// GetInventory 获取库存
	GetInventory(ctx context.Context, sku string) (*InventoryInfo, error)

	// UpdateInventory 更新库存
	UpdateInventory(ctx context.Context, sku string, quantity int) error

	// BatchUpdateInventory 批量更新库存
	BatchUpdateInventory(ctx context.Context, items []InventoryUpdateItem) error
}

// LogisticsAPI 物流接口
type LogisticsAPI interface {
	// GetLogisticsInfo 获取物流信息
	GetLogisticsInfo(ctx context.Context, orderID string) (*LogisticsInfo, error)

	// GetTrackingInfo 获取物流跟踪信息
	GetTrackingInfo(ctx context.Context, trackingNumber string) (*TrackingInfo, error)

	// GetShippingProviders 获取物流服务商列表
	GetShippingProviders(ctx context.Context) ([]ShippingProvider, error)
}

// 请求和响应结构体

type OrderListRequest struct {
	Page         int
	PageSize     int
	Status       string
	StartTime    *time.Time
	EndTime      *time.Time
	UpdatedAfter *time.Time
}

type OrderListResponse struct {
	Orders     []OrderDetail
	Total      int
	HasMore    bool
	NextCursor string
}

type OrderDetail struct {
	OrderID         string
	OrderSN         string
	Status          string
	PaymentStatus   string
	ShippingStatus  string
	BuyerUsername   string
	BuyerEmail      string
	RecipientName   string
	RecipientPhone  string
	ShippingAddress string
	Country         string
	TotalAmount     float64
	Currency        string
	ShippingFee     float64
	DiscountAmount  float64
	ActualAmount    float64
	Items           []OrderItem
	TrackingNumber  string
	ShippingCarrier string
	OrderTime       time.Time
	PaymentTime     *time.Time
	ShippedAt       *time.Time
	RawData         interface{}
}

type OrderItem struct {
	ItemID         string
	ProductID      string
	VariationID    string
	ProductName    string
	VariationName  string
	SKU            string
	Quantity       int
	UnitPrice      float64
	TotalPrice     float64
	DiscountAmount float64
	ImageURL       string
}

type ShipOrderRequest struct {
	OrderID         string
	TrackingNumber  string
	ShippingCarrier string
	ShipTime        time.Time
}

type ProductListRequest struct {
	Page       int
	PageSize   int
	Status     string
	CategoryID string
	Keyword    string
}

type ProductListResponse struct {
	Products   []ProductDetail
	Total      int
	HasMore    bool
	NextCursor string
}

type ProductDetail struct {
	ProductID      string
	Name           string
	Description    string
	Category       string
	Brand          string
	Status         string
	Price          float64
	OriginalPrice  float64
	Currency       string
	Stock          int
	SKU            string
	MainImage      string
	Images         []string
	Variations     []ProductVariation
	Weight         float64
	Dimensions     Dimensions
	SoldCount      int
	ViewCount      int
	Rating         float64
	ReviewCount    int
	RawData        interface{}
}

type ProductVariation struct {
	VariationID string
	Name        string
	SKU         string
	Price       float64
	Stock       int
	Attributes  map[string]string
	ImageURL    string
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
}

type CreateProductRequest struct {
	Name          string
	Description   string
	Category      string
	Brand         string
	Price         float64
	OriginalPrice float64
	Stock         int
	SKU           string
	Images        []string
	Variations    []ProductVariation
	Weight        float64
	Dimensions    Dimensions
}

type UpdateProductRequest struct {
	Name          *string
	Description   *string
	Price         *float64
	OriginalPrice *float64
	Stock         *int
	Images        []string
}

type InventoryInfo struct {
	SKU            string
	Stock          int
	ReservedStock  int
	AvailableStock int
	WarehouseCode  string
	Location       string
}

type InventoryUpdateItem struct {
	SKU      string
	Quantity int
}

type LogisticsInfo struct {
	OrderID         string
	TrackingNumber  string
	Carrier         string
	Status          string
	CurrentLocation string
	ShippedAt       *time.Time
	EstimatedDelivery *time.Time
	TrackingEvents  []TrackingEvent
}

type TrackingInfo struct {
	TrackingNumber string
	Carrier        string
	Status         string
	Events         []TrackingEvent
}

type TrackingEvent struct {
	Time        time.Time
	Status      string
	Location    string
	Description string
}

type ShippingProvider struct {
	Code string
	Name string
}
