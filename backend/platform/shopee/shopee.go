package shopee

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cross-border-admin/config"
	"cross-border-admin/models"
	"cross-border-admin/platform"
)

type ShopeeClient struct {
	config     config.PlatformConfig
	httpClient *http.Client
}

func NewShopeeClient() *ShopeeClient {
	cfg, exists := config.GetPlatformConfig("shopee")
	if !exists {
		cfg = config.PlatformConfig{
			APIURL: "https://partner.shopeemobile.com",
		}
	}
	return &ShopeeClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ShopeeClient) GetName() string {
	return "shopee"
}

func (c *ShopeeClient) Auth() platform.AuthAPI {
	return &ShopeeAuthAPI{client: c}
}

func (c *ShopeeClient) Order() platform.OrderAPI {
	return &ShopeeOrderAPI{client: c}
}

func (c *ShopeeClient) Product() platform.ProductAPI {
	return &ShopeeProductAPI{client: c}
}

func (c *ShopeeClient) Inventory() platform.InventoryAPI {
	return &ShopeeInventoryAPI{client: c}
}

func (c *ShopeeClient) Logistics() platform.LogisticsAPI {
	return &ShopeeLogisticsAPI{client: c}
}

// generateSign 生成签名
func (c *ShopeeClient) generateSign(path string, timestamp int64) string {
	baseString := fmt.Sprintf("%s%s%s%d", c.config.PartnerID, path, c.config.PartnerKey, timestamp)
	h := hmac.New(sha256.New, []byte(c.config.PartnerKey))
	h.Write([]byte(baseString))
	return hex.EncodeToString(h.Sum(nil))
}

// request 发送 HTTP 请求
func (c *ShopeeClient) request(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	timestamp := time.Now().Unix()
	sign := c.generateSign(path, timestamp)

	url := fmt.Sprintf("%s%s?partner_id=%s&timestamp=%d&sign=%s&shop_id=%s",
		c.config.APIURL, path, c.config.PartnerID, timestamp, sign, c.config.ShopID)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = io.NopCloser(nil)
		_ = jsonData
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// 记录 API 日志
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		c.logAPI(path, method, url, nil, 0, nil, duration, false, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	c.logAPI(path, method, url, nil, resp.StatusCode, respBody, duration, success, "")

	return respBody, nil
}

// logAPI 记录 API 调用日志
func (c *ShopeeClient) logAPI(apiName, method, url string, reqBody []byte, status int, respBody []byte, duration int64, success bool, errMsg string) {
	log := &models.APILog{
		Platform:       "shopee",
		APIName:        apiName,
		Method:         method,
		URL:            url,
		RequestBody:    string(reqBody),
		ResponseStatus: status,
		ResponseBody:   string(respBody),
		Duration:       duration,
		Success:        success,
		ErrorMessage:   errMsg,
	}
	models.DB.Create(log)
}

// Auth API 实现
type ShopeeAuthAPI struct {
	client *ShopeeClient
}

func (a *ShopeeAuthAPI) GetAccessToken(ctx context.Context) (string, error) {
	// TODO: 实现 Shopee 获取访问令牌逻辑
	return "", fmt.Errorf("未实现")
}

func (a *ShopeeAuthAPI) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", fmt.Errorf("未实现")
}

func (a *ShopeeAuthAPI) ValidateToken(ctx context.Context, token string) (bool, error) {
	return false, fmt.Errorf("未实现")
}

// Order API 实现
type ShopeeOrderAPI struct {
	client *ShopeeClient
}

func (o *ShopeeOrderAPI) GetOrderList(ctx context.Context, req *platform.OrderListRequest) (*platform.OrderListResponse, error) {
	// TODO: 实现 Shopee 获取订单列表逻辑
	return nil, fmt.Errorf("未实现")
}

func (o *ShopeeOrderAPI) GetOrderDetail(ctx context.Context, orderID string) (*platform.OrderDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *ShopeeOrderAPI) GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time) (*platform.OrderListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *ShopeeOrderAPI) CancelOrder(ctx context.Context, orderID string, reason string) error {
	return fmt.Errorf("未实现")
}

func (o *ShopeeOrderAPI) ShipOrder(ctx context.Context, req *platform.ShipOrderRequest) error {
	return fmt.Errorf("未实现")
}

// Product API 实现
type ShopeeProductAPI struct {
	client *ShopeeClient
}

func (p *ShopeeProductAPI) GetProductList(ctx context.Context, req *platform.ProductListRequest) (*platform.ProductListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *ShopeeProductAPI) GetProductDetail(ctx context.Context, productID string) (*platform.ProductDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *ShopeeProductAPI) CreateProduct(ctx context.Context, req *platform.CreateProductRequest) (*platform.ProductDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *ShopeeProductAPI) UpdateProduct(ctx context.Context, productID string, req *platform.UpdateProductRequest) error {
	return fmt.Errorf("未实现")
}

func (p *ShopeeProductAPI) DeleteProduct(ctx context.Context, productID string) error {
	return fmt.Errorf("未实现")
}

func (p *ShopeeProductAPI) UpdateProductStatus(ctx context.Context, productID string, status string) error {
	return fmt.Errorf("未实现")
}

// Inventory API 实现
type ShopeeInventoryAPI struct {
	client *ShopeeClient
}

func (i *ShopeeInventoryAPI) GetInventory(ctx context.Context, sku string) (*platform.InventoryInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (i *ShopeeInventoryAPI) UpdateInventory(ctx context.Context, sku string, quantity int) error {
	return fmt.Errorf("未实现")
}

func (i *ShopeeInventoryAPI) BatchUpdateInventory(ctx context.Context, items []platform.InventoryUpdateItem) error {
	return fmt.Errorf("未实现")
}

// Logistics API 实现
type ShopeeLogisticsAPI struct {
	client *ShopeeClient
}

func (l *ShopeeLogisticsAPI) GetLogisticsInfo(ctx context.Context, orderID string) (*platform.LogisticsInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (l *ShopeeLogisticsAPI) GetTrackingInfo(ctx context.Context, trackingNumber string) (*platform.TrackingInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (l *ShopeeLogisticsAPI) GetShippingProviders(ctx context.Context) ([]platform.ShippingProvider, error) {
	return nil, fmt.Errorf("未实现")
}
