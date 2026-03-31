package aliexpress

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"cross-border-admin/config"
	"cross-border-admin/models"
	"cross-border-admin/platform"
)

type AliExpressClient struct {
	config     config.PlatformConfig
	httpClient *http.Client
}

func NewAliExpressClient() *AliExpressClient {
	cfg, exists := config.GetPlatformConfig("aliexpress")
	if !exists {
		cfg = config.PlatformConfig{
			APIURL: "https://api-sg.aliexpress.com/sync",
		}
	}
	return &AliExpressClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AliExpressClient) GetName() string {
	return "aliexpress"
}

func (c *AliExpressClient) Auth() platform.AuthAPI {
	return &AliExpressAuthAPI{client: c}
}

func (c *AliExpressClient) Order() platform.OrderAPI {
	return &AliExpressOrderAPI{client: c}
}

func (c *AliExpressClient) Product() platform.ProductAPI {
	return &AliExpressProductAPI{client: c}
}

func (c *AliExpressClient) Inventory() platform.InventoryAPI {
	return &AliExpressInventoryAPI{client: c}
}

func (c *AliExpressClient) Logistics() platform.LogisticsAPI {
	return &AliExpressLogisticsAPI{client: c}
}

// generateSign 生成签名
func (c *AliExpressClient) generateSign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr strings.Builder
	signStr.WriteString(c.config.AppSecret)
	for _, k := range keys {
		signStr.WriteString(k)
		signStr.WriteString(params[k])
	}
	signStr.WriteString(c.config.AppSecret)

	h := hmac.New(md5.New, []byte(c.config.AppSecret))
	h.Write([]byte(signStr.String()))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// request 发送 HTTP 请求
func (c *AliExpressClient) request(ctx context.Context, method, apiMethod string, params map[string]string) ([]byte, error) {
	if params == nil {
		params = make(map[string]string)
	}

	params["app_key"] = c.config.AppKey
	params["method"] = apiMethod
	params["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	params["format"] = "json"
	params["v"] = "2.0"
	params["sign_method"] = "md5"

	sign := c.generateSign(params)
	params["sign"] = sign

	req, err := http.NewRequestWithContext(ctx, method, c.config.APIURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		c.logAPI(apiMethod, method, req.URL.String(), nil, 0, nil, duration, false, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	c.logAPI(apiMethod, method, req.URL.String(), nil, resp.StatusCode, respBody, duration, success, "")

	return respBody, nil
}

func (c *AliExpressClient) logAPI(apiName, method, url string, reqBody []byte, status int, respBody []byte, duration int64, success bool, errMsg string) {
	log := &models.APILog{
		Platform:       "aliexpress",
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
type AliExpressAuthAPI struct {
	client *AliExpressClient
}

func (a *AliExpressAuthAPI) GetAccessToken(ctx context.Context) (string, error) {
	return "", fmt.Errorf("未实现")
}

func (a *AliExpressAuthAPI) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return "", fmt.Errorf("未实现")
}

func (a *AliExpressAuthAPI) ValidateToken(ctx context.Context, token string) (bool, error) {
	return false, fmt.Errorf("未实现")
}

// Order API 实现
type AliExpressOrderAPI struct {
	client *AliExpressClient
}

func (o *AliExpressOrderAPI) GetOrderList(ctx context.Context, req *platform.OrderListRequest) (*platform.OrderListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *AliExpressOrderAPI) GetOrderDetail(ctx context.Context, orderID string) (*platform.OrderDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *AliExpressOrderAPI) GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time) (*platform.OrderListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (o *AliExpressOrderAPI) CancelOrder(ctx context.Context, orderID string, reason string) error {
	return fmt.Errorf("未实现")
}

func (o *AliExpressOrderAPI) ShipOrder(ctx context.Context, req *platform.ShipOrderRequest) error {
	return fmt.Errorf("未实现")
}

// Product API 实现
type AliExpressProductAPI struct {
	client *AliExpressClient
}

func (p *AliExpressProductAPI) GetProductList(ctx context.Context, req *platform.ProductListRequest) (*platform.ProductListResponse, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *AliExpressProductAPI) GetProductDetail(ctx context.Context, productID string) (*platform.ProductDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *AliExpressProductAPI) CreateProduct(ctx context.Context, req *platform.CreateProductRequest) (*platform.ProductDetail, error) {
	return nil, fmt.Errorf("未实现")
}

func (p *AliExpressProductAPI) UpdateProduct(ctx context.Context, productID string, req *platform.UpdateProductRequest) error {
	return fmt.Errorf("未实现")
}

func (p *AliExpressProductAPI) DeleteProduct(ctx context.Context, productID string) error {
	return fmt.Errorf("未实现")
}

func (p *AliExpressProductAPI) UpdateProductStatus(ctx context.Context, productID string, status string) error {
	return fmt.Errorf("未实现")
}

// Inventory API 实现
type AliExpressInventoryAPI struct {
	client *AliExpressClient
}

func (i *AliExpressInventoryAPI) GetInventory(ctx context.Context, sku string) (*platform.InventoryInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (i *AliExpressInventoryAPI) UpdateInventory(ctx context.Context, sku string, quantity int) error {
	return fmt.Errorf("未实现")
}

func (i *AliExpressInventoryAPI) BatchUpdateInventory(ctx context.Context, items []platform.InventoryUpdateItem) error {
	return fmt.Errorf("未实现")
}

// Logistics API 实现
type AliExpressLogisticsAPI struct {
	client *AliExpressClient
}

func (l *AliExpressLogisticsAPI) GetLogisticsInfo(ctx context.Context, orderID string) (*platform.LogisticsInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (l *AliExpressLogisticsAPI) GetTrackingInfo(ctx context.Context, trackingNumber string) (*platform.TrackingInfo, error) {
	return nil, fmt.Errorf("未实现")
}

func (l *AliExpressLogisticsAPI) GetShippingProviders(ctx context.Context) ([]platform.ShippingProvider, error) {
	return nil, fmt.Errorf("未实现")
}
