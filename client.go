package funnelfox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var defaultFunnelFoxHTTPClient = &http.Client{
	Timeout: 2 * time.Minute,
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	orgID      string
	secretKey  string
	logger     Logger
}

// NewClient 创建新的 FunnelFox 客户端
// orgID: 组织ID
// secretKey: 密钥，部分 API 需要
// logger: 日志记录器，可以为 nil（使用 NopLogger）
func NewClient(orgID, secretKey string, logger Logger) *Client {
	baseURL := fmt.Sprintf("https://billing.funnelfox.com/%s/v1", orgID)
	if logger == nil {
		logger = &NopLogger{}
	}
	return &Client{
		httpClient: defaultFunnelFoxHTTPClient,
		baseURL:    baseURL,
		orgID:      orgID,
		secretKey:  secretKey,
		logger:     logger,
	}
}

// NewClientWithHTTPClient 使用自定义 HTTP 客户端创建 FunnelFox 客户端
func NewClientWithHTTPClient(orgID, secretKey string, httpClient *http.Client, logger Logger) *Client {
	baseURL := fmt.Sprintf("https://billing.funnelfox.com/%s/v1", orgID)
	if httpClient == nil {
		httpClient = defaultFunnelFoxHTTPClient
	}
	if logger == nil {
		logger = &NopLogger{}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		orgID:      orgID,
		secretKey:  secretKey,
		logger:     logger,
	}
}

// doRequest 执行 HTTP POST 请求
func (c *Client) doRequest(endpoint string, requestBody, response any, withSecretKey bool) *Error {
	url := c.baseURL + endpoint
	var bodyReader io.Reader
	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return WrapError(err, "failed to marshal request body")
		}
		bodyReader = bytes.NewReader(bodyBytes)
		c.logger.Debug("funnelfox_request",
			String("url", url),
			String("body", string(bodyBytes)))
	}

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return WrapError(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if withSecretKey {
		req.Header.Set("ff-secret-key", c.secretKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("funnelfox_request_error",
			String("url", url),
			ErrorField(err))
		return WrapError(err, "request failed")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("funnelfox_read_response_error",
			String("url", url),
			ErrorField(err))
		return WrapError(err, "failed to read response")
	}

	c.logger.Debug("funnelfox_response",
		String("url", url),
		Number("status_code", resp.StatusCode),
		String("body", string(respBody)))

	// 解析响应
	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return WrapError(err, "failed to unmarshal response")
	}

	// 检查响应状态
	if apiResp.Status == "error" {
		var errMsg string
		errMsgs := make([]string, 0, len(apiResp.Error))
		for _, err := range apiResp.Error {
			errMsgs = append(errMsgs, fmt.Sprintf("%s_%s", err.Type, err.Msg))
		}
		if len(errMsgs) > 0 {
			errMsg = strings.Join(errMsgs, " ")
		} else {
			errMsg = "API error"
		}
		c.logger.Error("funnelfox_api_error",
			String("url", url),
			String("req_id", apiResp.ReqID),
			String("error", errMsg))
		return WrapErrorf(nil, "FunnelFox API error: %s (req_id: %s)", errMsg, apiResp.ReqID)
	}

	// 如果响应状态码不是 200-299，返回错误
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return WrapErrorf(nil, "HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应数据
	if response != nil && len(apiResp.Data) > 0 {
		if err := json.Unmarshal(apiResp.Data, response); err != nil {
			return WrapError(err, "failed to unmarshal response data")
		}
	}

	return nil
}

// ===== Payment Management =====

// Refund 退款订单（全额或部分，可选软退款）
func (c *Client) Refund(req RefundRequest) *Error {
	if err := c.doRequest("/payment/refund", req, nil, true); err != nil {
		return err
	}
	return nil
}

// OneClickPurchase 执行一键购买（使用用户保存的支付方式）
func (c *Client) OneClickPurchase(req OneClickPurchaseRequest) (*OneClickPurchaseResponse, *Error) {
	var resp OneClickPurchaseResponse
	if err := c.doRequest("/checkout/one_click", req, &resp, false); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ===== Subscription Management =====

// EnableAutoRenew 启用自动续费
func (c *Client) EnableAutoRenew(req EnableAutoRenewRequest) *Error {
	return c.doRequest("/subscription/enable_autorenew", req, nil, true)
}

// DisableAutoRenew 禁用自动续费
func (c *Client) DisableAutoRenew(req DisableAutoRenewRequest) *Error {
	return c.doRequest("/subscription/disable_autorenew", req, nil, true)
}

// SubscriptionMigration 迁移订阅到另一个价格点
func (c *Client) SubscriptionMigration(req SubscriptionMigrationRequest) (*SubscriptionMigrationResponse, *Error) {
	var resp SubscriptionMigrationResponse
	if err := c.doRequest("/subscription/migration", req, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ApplyDiscount 应用百分比折扣
func (c *Client) ApplyDiscount(req DiscountRequest) *Error {
	return c.doRequest("/discount", req, nil, true)
}

// DeferSubscription 延迟订阅的下次扣费时间
func (c *Client) DeferSubscription(req SubscriptionDeferRequest) *Error {
	return c.doRequest("/subscription/defer", req, nil, true)
}

// PauseSubscription 暂停订阅
func (c *Client) PauseSubscription(req SubscriptionPauseRequest) *Error {
	return c.doRequest("/subscription/pause", req, nil, true)
}

// ResumeSubscription 恢复订阅
func (c *Client) ResumeSubscription(req SubscriptionResumeRequest) *Error {
	return c.doRequest("/subscription/resume", req, nil, true)
}

// ===== PricePoints =====

// ListPricePoints 列出价格点（可按 ident 过滤）
func (c *Client) ListPricePoints(req PricePointsListRequest) (*PricePointsListResponse, *Error) {
	var resp PricePointsListResponse
	if err := c.doRequest("/price_points", req, &resp, false); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateFeature 创建 feature
func (c *Client) CreateFeature(req FeatureCreateRequest) (*FeatureCreateResponse, *Error) {
	var resp FeatureCreateResponse
	if err := c.doRequest("/feature/create", req, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreatePricePoint 创建 price point
func (c *Client) CreatePricePoint(req PricePointCreateRequest) (*PricePointCreateResponse, *Error) {
	var resp PricePointCreateResponse
	if err := c.doRequest("/pp/create", req, &resp, true); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ===== Information =====

// GetMyAssets 获取用户资产（订阅和一次性购买）
func (c *Client) GetMyAssets(req MyAssetsRequest) (*MyAssetsResponse, *Error) {
	var resp rawMyAssetsResponse
	if err := c.doRequest("/my_assets", req, &resp, false); err != nil {
		return nil, err
	}
	return resp.toMyAssetsResponse(), nil
}

// GetPaymentsHistory 获取支付历史
func (c *Client) GetPaymentsHistory(req PaymentsHistoryRequest) (*PaymentsHistoryResponse, *Error) {
	var resp rawPaymentsHistoryResponse
	if err := c.doRequest("/payments_history", req, &resp, true); err != nil {
		return nil, err
	}
	return resp.toPaymentsHistoryResponse(), nil
}
