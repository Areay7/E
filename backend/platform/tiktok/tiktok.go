package tiktok

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

type TikTokClient struct {
	config     config.PlatformConfig
	httpClient *http.Client
}

func NewTikTokClient() *TikTokClient {
	cfg, exists := config.GetPlatformConfig("tiktok")
	if !exists {
		cfg = config.PlatformConfig{
			APIURL: "https://open-api.tiktokglobalshop.com",
		}
	}
	return &TikTokClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *TikTokClient) GetName() string {
	return "tiktok"
}

func (c *TikTokClient) Auth() platform.AuthAPI {
	return &TikTokAuthAPI{client: c}
}

func (c *TikTokClient) Order() platform.OrderAPI {
	return &TikTokOrderAPI{client: c}
}

func (c *TikTokClient) Product() platform.ProductAPI {
	return &TikTokProductAPI{client: c}
}

func (c *TikTokClient) Inventory() platform.InventoryAPI {
	return &TikTokInventoryAPI{client: c}
}

func (c *TikTokClient) Logistics() platform.LogisticsAPI {
	return &TikTokLogisticsAPI{client: c}
}

// generateSign 生成签名
func (c *TikTokClient) generateSign(path string, timestamp int64, body []byte) string {
	message := fmt.Sprintf("%s%d%s", path, timestamp, string(body))
	h := hmac.New(sha256.New, []byte(c.config.AppSecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// request 发送 HTTP 请求
func (c *TikTokClient) request(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	timestamp := time.Now().Unix()

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	sign := c.generateSign(path, timestamp, reqBody)
	url := fmt.Sprintf("%s%s", c.config.APIURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-tts-access-token", c.config.AppKey)
	req.Header.Set("x-tts-timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("x-tts-signature", sign)

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		c.logAPI(path, method, url, reqBody, 0, nil, duration, false, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	c.logAPI(path, method, url, reqBody, resp.StatusCode, respBody, duration, success, "")

	return respBody, nil
}

func (c *TikTokClient) logAPI(apiName, method, url string, reqBody []byte, status int, respBody []byte, duration int64, success bool, errMsg string) {
	log := &models.APILog{
		Platform:       "tiktok",
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
type TikTokAuthAPI struct {
	client *TikTokClient
}

func (a *TikTokAuthAPI) GetAccessToken(ctx context.Context) (string, error) {
	return "", fmt.Errorf("未实现")
}

func (a *TikTokAuthAPI) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", fmt.Errorf("未实现")
}

func (a *TikTokAuthAPI) ValidateToken(ctx context.Context, token string) (bool, error) {
	return false, fmt.Errorf("未实现")
}

// Order API 实现
type TikTokOrderAPI struct {
	client *TikTokClient
}

func (o *TikTokOrderAPI) GetOrderList(ctx context.Context, req *platform.OrderListRequest) (*platform.OrderListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *TikTokOrderAPI) GetOrderDetail(ctx context.Context, orderID string) (*platform.OrderDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *TikTokOrderAPI) GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time) (*platform.OrderListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *TikTokOrderAPI) CancelOrder(ctx context.Context, orderID string, reason string) error {
	return fmt.Errorf("未实现")
}

func (o *TikTokOrderAPI) ShipOrder(ctx context.Context, req *platform.ShipOrderRequest) error {
	return fmt.Errorf("未实现")
}

// Product API 实现
type TikTokProductAPI struct {
	client *TikTokClient
}

func (p *TikTokProductAPI) GetProductList(ctx context.Context, req *platform.ProductListRequest) (*platform.ProductListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *TikTokProductAPI) GetProductDetail(ctx context.Context, productID string) (*platform.ProductDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *TikTokProductAPI) CreateProduct(ctx context.Context, req *platform.CreateProductRequest) (*platform.ProductDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *TikTokProductAPI) UpdateProduct(ctx context.Context, productID string, req *platform.UpdateProductRequest) error {
	return fmt.Errorf("未实现")
}

func (p *TikTokProductAPI) DeleteProduct(ctx context.Context, productID string) error {
	return fmt.Errorf("未实现")
}

func (p *TikTokProductAPI) UpdateProductStatus(ctx context.Context, productID string, status string) error {
	return fmt.Errorf("未实现")
}

// Inventory API 实现
type TikTokInventoryAPI struct {
	client *TikTokClient
}

func (i *TikTokInventoryAPI) GetInventory(ctx context.Context, sku string) (*platform.InventoryInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (i *TikTokInventoryAPI) UpdateInventory(ctx context.Context, sku string, quantity int) error {
	return fmt.Errorf("未实现")
}

func (i *TikTokInventoryAPI) BatchUpdateInventory(ctx context.Context, items []platform.InventoryUpdateItem) error {
	return fmt.Errorf("未实现")
}

// Logistics API 实现
type TikTokLogisticsAPI struct {
	client *TikTokClient
}

func (l *TikTokLogisticsAPI) GetLogisticsInfo(ctx context.Context, orderID string) (*platform.LogisticsInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (l *TikTokLogisticsAPI) GetTrackingInfo(ctx context.Context, trackingNumber string) (*platform.TrackingInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (l *TikTokLogisticsAPI) GetShippingProviders(ctx context.Context) ([]platform.ShippingProvider, error) {
	return nil, fmt.Errorf("未实现")
}
