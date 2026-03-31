package models

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	OrderID         string         `gorm:"type:varchar(100);uniqueIndex:idx_platform_order;not null" json:"orderId"`
	Platform        string         `gorm:"index;uniqueIndex:idx_platform_order;not null" json:"platform"`
	ShopID          string         `gorm:"index" json:"shopId"`
	OrderSN         string         `gorm:"index" json:"orderSn"`
	Status          string         `gorm:"index" json:"status"`
	PaymentStatus   string         `json:"paymentStatus"`
	ShippingStatus  string         `json:"shippingStatus"`

	// 买家信息
	BuyerUsername   string         `json:"buyerUsername"`
	BuyerEmail      string         `json:"buyerEmail"`

	// 收货地址
	RecipientName   string         `json:"recipientName"`
	RecipientPhone  string         `json:"recipientPhone"`
	ShippingAddress string         `json:"shippingAddress"`
	Country         string         `gorm:"index" json:"country"`
	Province        string         `json:"province"`
	City            string         `json:"city"`
	District        string         `json:"district"`
	Zipcode         string         `json:"zipcode"`

	// 金额信息
	TotalAmount     float64        `json:"totalAmount"`
	Currency        string         `json:"currency"`
	ShippingFee     float64        `json:"shippingFee"`
	DiscountAmount  float64        `json:"discountAmount"`
	ActualAmount    float64        `json:"actualAmount"`

	// 物流信息
	TrackingNumber  string         `gorm:"index" json:"trackingNumber"`
	ShippingCarrier string         `json:"shippingCarrier"`
	ShippedAt       *time.Time     `json:"shippedAt"`
	DeliveredAt     *time.Time     `json:"deliveredAt"`

	// 备注
	BuyerMessage    string         `json:"buyerMessage"`
	SellerNote      string         `json:"sellerNote"`

	// 平台原始数据
	RawData         string         `gorm:"type:text" json:"rawData"`

	// 时间戳
	OrderTime       time.Time      `gorm:"index" json:"orderTime"`
	PaymentTime     *time.Time     `json:"paymentTime"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderItem 订单明细
type OrderItem struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	OrderID         uint           `gorm:"index;not null" json:"orderId"`

	Platform        string         `json:"platform"`
	ItemID          string         `json:"itemId"`
	ProductID       string         `gorm:"index" json:"productId"`
	VariationID     string         `json:"variationId"`

	ProductName     string         `json:"productName"`
	VariationName   string         `json:"variationName"`
	SKU             string         `gorm:"index" json:"sku"`

	Quantity        int            `json:"quantity"`
	UnitPrice       float64        `json:"unitPrice"`
	TotalPrice      float64        `json:"totalPrice"`
	DiscountAmount  float64        `json:"discountAmount"`

	ImageURL        string         `json:"imageUrl"`

	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// Product 商品模型
type Product struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Platform        string         `gorm:"index;uniqueIndex:idx_platform_product;not null" json:"platform"`
	ProductID       string         `gorm:"uniqueIndex:idx_platform_product;not null" json:"productId"`
	ShopID          string         `gorm:"index" json:"shopId"`

	Name            string         `json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	Category        string         `gorm:"index" json:"category"`
	Brand           string         `json:"brand"`

	Status          string         `gorm:"index" json:"status"`
	Price           float64        `json:"price"`
	OriginalPrice   float64        `json:"originalPrice"`
	Currency        string         `json:"currency"`

	Stock           int            `json:"stock"`
	SKU             string         `gorm:"index" json:"sku"`

	MainImage       string         `json:"mainImage"`
	Images          string         `gorm:"type:text" json:"images"`

	Weight          float64        `json:"weight"`
	Length          float64        `json:"length"`
	Width           float64        `json:"width"`
	Height          float64        `json:"height"`

	SoldCount       int            `json:"soldCount"`
	ViewCount       int            `json:"viewCount"`
	Rating          float64        `json:"rating"`
	ReviewCount     int            `json:"reviewCount"`

	RawData         string         `gorm:"type:text" json:"rawData"`

	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProductVariation 商品变体
type ProductVariation struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	ProductID       uint           `gorm:"index;not null" json:"productId"`

	Platform        string         `json:"platform"`
	VariationID     string         `gorm:"index" json:"variationId"`

	Name            string         `json:"name"`
	SKU             string         `gorm:"index" json:"sku"`

	Price           float64        `json:"price"`
	Stock           int            `json:"stock"`

	Attributes      string         `gorm:"type:text" json:"attributes"`
	ImageURL        string         `json:"imageUrl"`

	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// Inventory 库存记录
type Inventory struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Platform        string         `gorm:"index;uniqueIndex:idx_platform_sku;not null" json:"platform"`
	SKU             string         `gorm:"uniqueIndex:idx_platform_sku;not null" json:"sku"`
	ProductID       string         `gorm:"index" json:"productId"`
	VariationID     string         `json:"variationId"`

	Stock           int            `json:"stock"`
	ReservedStock   int            `json:"reservedStock"`
	AvailableStock  int            `json:"availableStock"`

	WarehouseCode   string         `gorm:"index" json:"warehouseCode"`
	Location        string         `json:"location"`

	LastSyncAt      *time.Time     `json:"lastSyncAt"`

	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// InventoryLog 库存变动日志
type InventoryLog struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	InventoryID     uint           `gorm:"index;not null" json:"inventoryId"`

	Type            string         `gorm:"index" json:"type"`
	Quantity        int            `json:"quantity"`
	BeforeStock     int            `json:"beforeStock"`
	AfterStock      int            `json:"afterStock"`

	Reason          string         `json:"reason"`
	RelatedOrderID  string         `gorm:"index" json:"relatedOrderId"`
	OperatorID      uint           `json:"operatorId"`

	CreatedAt       time.Time      `json:"createdAt"`
}

// Logistics 物流信息
type Logistics struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	OrderID         uint           `gorm:"index;not null" json:"orderId"`

	Platform        string         `json:"platform"`
	TrackingNumber  string         `gorm:"index" json:"trackingNumber"`
	Carrier         string         `json:"carrier"`

	Status          string         `gorm:"index" json:"status"`
	CurrentLocation string         `json:"currentLocation"`

	ShippedAt       *time.Time     `json:"shippedAt"`
	InTransitAt     *time.Time     `json:"inTransitAt"`
	DeliveredAt     *time.Time     `json:"deliveredAt"`

	EstimatedDelivery *time.Time   `json:"estimatedDelivery"`

	TrackingEvents  string         `gorm:"type:text" json:"trackingEvents"`

	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// PlatformConfig 平台配置
type PlatformConfig struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Platform        string         `gorm:"uniqueIndex;not null" json:"platform"`
	ShopID          string         `json:"shopId"`
	ShopName        string         `json:"shopName"`

	Enabled         bool           `gorm:"default:true" json:"enabled"`

	AppKey          string         `json:"appKey"`
	AppSecret       string         `json:"-"`
	PartnerID       string         `json:"partnerId"`
	PartnerKey      string         `json:"-"`

	AccessToken     string         `json:"-"`
	RefreshToken    string         `json:"-"`
	TokenExpiresAt  *time.Time     `json:"tokenExpiresAt"`

	APIURL          string         `json:"apiUrl"`

	SyncEnabled     bool           `gorm:"default:true" json:"syncEnabled"`
	LastSyncAt      *time.Time     `json:"lastSyncAt"`

	Settings        string         `gorm:"type:text" json:"settings"`

	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// APILog API调用日志
type APILog struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Platform        string         `gorm:"index" json:"platform"`
	APIName         string         `gorm:"index" json:"apiName"`
	Method          string         `json:"method"`
	URL             string         `json:"url"`

	RequestHeaders  string         `gorm:"type:text" json:"requestHeaders"`
	RequestBody     string         `gorm:"type:text" json:"requestBody"`

	ResponseStatus  int            `json:"responseStatus"`
	ResponseHeaders string         `gorm:"type:text" json:"responseHeaders"`
	ResponseBody    string         `gorm:"type:text" json:"responseBody"`

	Duration        int64          `json:"duration"`
	Success         bool           `gorm:"index" json:"success"`
	ErrorMessage    string         `json:"errorMessage"`

	CreatedAt       time.Time      `gorm:"index" json:"createdAt"`
}

// SyncTask 同步任务
type SyncTask struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Platform        string         `gorm:"index" json:"platform"`
	TaskType        string         `gorm:"index" json:"taskType"`
	Status          string         `gorm:"index" json:"status"`

	StartTime       *time.Time     `json:"startTime"`
	EndTime         *time.Time     `json:"endTime"`

	TotalCount      int            `json:"totalCount"`
	SuccessCount    int            `json:"successCount"`
	FailCount       int            `json:"failCount"`

	ErrorMessage    string         `gorm:"type:text" json:"errorMessage"`

	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}
