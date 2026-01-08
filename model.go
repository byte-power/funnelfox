package funnelfox

import (
	"encoding/json"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999999"

// Response FunnelFox Billing API 响应结构
type Response struct {
	Data   json.RawMessage `json:"data,omitempty"`
	Status string          `json:"status"`
	ReqID  string          `json:"req_id,omitempty"`
	Error  []APIError      `json:"error,omitempty"`
}

// APIError FunnelFox API 错误结构
type APIError struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

// ===== Payment Management =====

// RefundRequest 退款请求
type RefundRequest struct {
	ExternalID string  `json:"external_id"`
	OrderID    string  `json:"order_id"`         // 订单ID
	Reason     *string `json:"reason,omitempty"` // 退款原因（可选）
	Comment    *string `json:"comment,omitempty"`
	Amount     *string `json:"amount,omitempty"`      // 退款金额（可选，用于部分退款）
	SoftRefund *bool   `json:"soft_refund,omitempty"` // 是否软退款（可选）
}

// PaymentsHistoryRequest 支付历史请求
type PaymentsHistoryRequest struct {
	ExternalID string `json:"external_id"`
}

// OneClickPurchaseRequest 一键购买请求
type OneClickPurchaseRequest struct {
	ExternalID     string         `json:"external_id"`               // 用户外部ID
	PPIdent        string         `json:"pp_ident"`                  // 价格点标识
	ClientMetadata map[string]any `json:"client_metadata,omitempty"` // 客户端元数据（可选）
}

// OneClickPurchaseResponse 一键购买响应
type OneClickPurchaseResponse PaymentResult

// rawPayment 原始支付信息（用于解析）
type rawPayment struct {
	Amount        string   `json:"amount"`
	CreatedAt     string   `json:"created_at"`
	Currency      Currency `json:"currency"`
	Last4         string   `json:"last4"`
	Network       string   `json:"network"`
	OneoffID      *string  `json:"oneoff_id"`
	OrderID       string   `json:"order_id"`
	PaymentMethod string   `json:"payment_method"`
	Refunded      string   `json:"refunded"`
	SubsID        string   `json:"subs_id"`
}

// Payment 支付信息
type Payment struct {
	Amount        string     `json:"amount"`
	CreatedAt     *time.Time `json:"created_at"`
	Currency      Currency   `json:"currency"`
	Last4         string     `json:"last4"`
	Network       string     `json:"network"`
	OneoffID      *string    `json:"oneoff_id"`
	OrderID       string     `json:"order_id"`
	PaymentMethod string     `json:"payment_method"`
	Refunded      string     `json:"refunded"`
	SubsID        string     `json:"subs_id"`
}

// rawPaymentsHistoryResponse 原始支付历史响应
type rawPaymentsHistoryResponse struct {
	Payments []rawPayment `json:"payments"`
}

// PaymentsHistoryResponse 支付历史响应
type PaymentsHistoryResponse struct {
	Payments []Payment `json:"payments"`
}

func (raw rawPaymentsHistoryResponse) toPaymentsHistoryResponse() *PaymentsHistoryResponse {
	var res PaymentsHistoryResponse
	for _, rawPayment := range raw.Payments {
		createdAt := parseTimePointer(rawPayment.CreatedAt)
		res.Payments = append(res.Payments, Payment{
			Amount:        rawPayment.Amount,
			CreatedAt:     createdAt,
			Currency:      rawPayment.Currency,
			Last4:         rawPayment.Last4,
			Network:       rawPayment.Network,
			OneoffID:      rawPayment.OneoffID,
			OrderID:       rawPayment.OrderID,
			PaymentMethod: rawPayment.PaymentMethod,
			Refunded:      rawPayment.Refunded,
			SubsID:        rawPayment.SubsID,
		})
	}
	return &res
}

// TransactionReportRequest 获取所有交易请求
type TransactionReportRequest struct {
	LastTransactionDate string  `json:"last_transaction_date"` // 最后交易日期（RFC3339 格式）
	SubsID              *string `json:"subs_id,omitempty"`     // 订阅ID（可选）
	OrderID             *string `json:"order_id,omitempty"`    // 订单ID（可选）
	OneoffID            *string `json:"oneoff_id,omitempty"`   // 一次性购买ID（可选）
	Limit               *int    `json:"limit,omitempty"`       // 限制返回数量（可选，默认100，范围1-500）
}

// rawTransaction 原始交易信息（用于解析）
type rawTransaction struct {
	OrderID                             string  `json:"order_id"`
	PPID                                int     `json:"pp_id"`
	PPVersion                           int     `json:"pp_version"`
	Status                              string  `json:"status"`
	IsCIT                               bool    `json:"is_cit"`
	IntegrationType                     string  `json:"integration_type"`
	Region                              string  `json:"region"`
	Amount                              string  `json:"amount"`
	CurrencyCode                        string  `json:"currency_code"`
	TrxID                               string  `json:"trx_id"`
	PSP                                 string  `json:"psp"`
	PSPTransactionID                    string  `json:"psp_transaction_id"`
	AmountUSD                           string  `json:"amount_usd"`
	IsFallback                          bool    `json:"is_fallback"`
	PSPMerchantID                       string  `json:"psp_merchant_id"`
	PSPTransactionType                  string  `json:"psp_transaction_type"`
	PSPStatus                           string  `json:"psp_status"`
	PSPDate                             string  `json:"psp_date"`
	PSPCardTokenType                    *string `json:"psp_card_token_type"`
	PSPReasonMessage                    *string `json:"psp_reason_message"`
	PSPReasonType                       *string `json:"psp_reason_type"`
	PSPReasonCode                       *string `json:"psp_reason_code"`
	PSPReasonDeclineType                *string `json:"psp_reason_decline_type"`
	TrxCreatedAt                        string  `json:"trx_created_at"`
	MetaSubsID                          *string `json:"meta_subs_id"`
	MetaSubsIteration                   *string `json:"meta_subs_iteration"`
	MetaRetryStep                       *string `json:"meta_retry_step"`
	MetaOneoffID                        *string `json:"meta_oneoff_id"`
	MetaClientFFProjectID               *string `json:"meta_client_ff_project_id"`
	MetaClientFFSessionID               *string `json:"meta_client_ff_session_id"`
	MetaClientFFPriceID                 *string `json:"meta_client_ff_price_id"`
	PMType                              string  `json:"pm_type"`
	ThreeDSChallengeIssued              *bool   `json:"threeds_challenge_issued"`
	ThreeDSProtocolVersion              *string `json:"threeds_protocol_version"`
	ThreeDSResponseCode                 *string `json:"threeds_response_code"`
	ThreeDSReasonCode                   *string `json:"threeds_reason_code"`
	ThreeDSReasonText                   *string `json:"threeds_reason_text"`
	AuthorizationType                   string  `json:"authorization_type"`
	IsVaulted                           bool    `json:"is_vaulted"`
	PMDataBinAccountFundingType         *string `json:"pm_data_bin_account_funding_type"`
	PMDataBinAccountNumberType          *string `json:"pm_data_bin_account_number_type"`
	PMDataBinIssuerName                 *string `json:"pm_data_bin_issuer_name"`
	PMDataBinIssuerCountryCode          *string `json:"pm_data_bin_issuer_country_code"`
	PMDataBinIssuerCurrencyCode         *string `json:"pm_data_bin_issuer_currency_code"`
	Network                             *string `json:"network"`
	PMDataBinPrepaidReloadableIndicator *string `json:"pm_data_bin_prepaid_reloadable_indicator"`
	PMDataBinProductCode                *string `json:"pm_data_bin_product_code"`
	PMDataBinProductName                *string `json:"pm_data_bin_product_name"`
	PMDataBinProductUsageType           *string `json:"pm_data_bin_product_usage_type"`
	PMDataBinRegionalRestriction        *string `json:"pm_data_bin_regional_restriction"`
	PMDataExpirationDate                *string `json:"pm_data_expiration_date"`
	PMDataFirst6                        *string `json:"pm_data_first6"`
	PMDataLast4                         *string `json:"pm_data_last4"`
	PMDataIsNetworkTokenized            *bool   `json:"pm_data_is_network_tokenized"`
}

// Transaction 交易信息
type Transaction struct {
	OrderID                             string     `json:"order_id"`
	PPID                                int        `json:"pp_id"`
	PPVersion                           int        `json:"pp_version"`
	Status                              string     `json:"status"`
	IsCIT                               bool       `json:"is_cit"`
	IntegrationType                     string     `json:"integration_type"`
	Region                              string     `json:"region"`
	Amount                              string     `json:"amount"`
	CurrencyCode                        string     `json:"currency_code"`
	TrxID                               string     `json:"trx_id"`
	PSP                                 string     `json:"psp"`
	PSPTransactionID                    string     `json:"psp_transaction_id"`
	AmountUSD                           string     `json:"amount_usd"`
	IsFallback                          bool       `json:"is_fallback"`
	PSPMerchantID                       string     `json:"psp_merchant_id"`
	PSPTransactionType                  string     `json:"psp_transaction_type"`
	PSPStatus                           string     `json:"psp_status"`
	PSPDate                             *time.Time `json:"psp_date"`
	PSPCardTokenType                    *string    `json:"psp_card_token_type"`
	PSPReasonMessage                    *string    `json:"psp_reason_message"`
	PSPReasonType                       *string    `json:"psp_reason_type"`
	PSPReasonCode                       *string    `json:"psp_reason_code"`
	PSPReasonDeclineType                *string    `json:"psp_reason_decline_type"`
	TrxCreatedAt                        *time.Time `json:"trx_created_at"`
	MetaSubsID                          *string    `json:"meta_subs_id"`
	MetaSubsIteration                   *string    `json:"meta_subs_iteration"`
	MetaRetryStep                       *string    `json:"meta_retry_step"`
	MetaOneoffID                        *string    `json:"meta_oneoff_id"`
	MetaClientFFProjectID               *string    `json:"meta_client_ff_project_id"`
	MetaClientFFSessionID               *string    `json:"meta_client_ff_session_id"`
	MetaClientFFPriceID                 *string    `json:"meta_client_ff_price_id"`
	PMType                              string     `json:"pm_type"`
	ThreeDSChallengeIssued              *bool      `json:"threeds_challenge_issued"`
	ThreeDSProtocolVersion              *string    `json:"threeds_protocol_version"`
	ThreeDSResponseCode                 *string    `json:"threeds_response_code"`
	ThreeDSReasonCode                   *string    `json:"threeds_reason_code"`
	ThreeDSReasonText                   *string    `json:"threeds_reason_text"`
	AuthorizationType                   string     `json:"authorization_type"`
	IsVaulted                           bool       `json:"is_vaulted"`
	PMDataBinAccountFundingType         *string    `json:"pm_data_bin_account_funding_type"`
	PMDataBinAccountNumberType          *string    `json:"pm_data_bin_account_number_type"`
	PMDataBinIssuerName                 *string    `json:"pm_data_bin_issuer_name"`
	PMDataBinIssuerCountryCode          *string    `json:"pm_data_bin_issuer_country_code"`
	PMDataBinIssuerCurrencyCode         *string    `json:"pm_data_bin_issuer_currency_code"`
	Network                             *string    `json:"network"`
	PMDataBinPrepaidReloadableIndicator *string    `json:"pm_data_bin_prepaid_reloadable_indicator"`
	PMDataBinProductCode                *string    `json:"pm_data_bin_product_code"`
	PMDataBinProductName                *string    `json:"pm_data_bin_product_name"`
	PMDataBinProductUsageType           *string    `json:"pm_data_bin_product_usage_type"`
	PMDataBinRegionalRestriction        *string    `json:"pm_data_bin_regional_restriction"`
	PMDataExpirationDate                *string    `json:"pm_data_expiration_date"`
	PMDataFirst6                        *string    `json:"pm_data_first6"`
	PMDataLast4                         *string    `json:"pm_data_last4"`
	PMDataIsNetworkTokenized            *bool      `json:"pm_data_is_network_tokenized"`
}

// rawTransactionReportResponse 原始交易报告响应
type rawTransactionReportResponse struct {
	Transactions []rawTransaction `json:"transactions"`
}

// TransactionReportResponse 交易报告响应
type TransactionReportResponse struct {
	Transactions []Transaction `json:"transactions"`
}

func (raw rawTransactionReportResponse) toTransactionReportResponse() *TransactionReportResponse {
	var res TransactionReportResponse
	for _, rawTx := range raw.Transactions {
		pspDate := parseTimePointer(rawTx.PSPDate)
		trxCreatedAt := parseTimePointer(rawTx.TrxCreatedAt)
		res.Transactions = append(res.Transactions, Transaction{
			OrderID:                             rawTx.OrderID,
			PPID:                                rawTx.PPID,
			PPVersion:                           rawTx.PPVersion,
			Status:                              rawTx.Status,
			IsCIT:                               rawTx.IsCIT,
			IntegrationType:                     rawTx.IntegrationType,
			Region:                              rawTx.Region,
			Amount:                              rawTx.Amount,
			CurrencyCode:                        rawTx.CurrencyCode,
			TrxID:                               rawTx.TrxID,
			PSP:                                 rawTx.PSP,
			PSPTransactionID:                    rawTx.PSPTransactionID,
			AmountUSD:                           rawTx.AmountUSD,
			IsFallback:                          rawTx.IsFallback,
			PSPMerchantID:                       rawTx.PSPMerchantID,
			PSPTransactionType:                  rawTx.PSPTransactionType,
			PSPStatus:                           rawTx.PSPStatus,
			PSPDate:                             pspDate,
			PSPCardTokenType:                    rawTx.PSPCardTokenType,
			PSPReasonMessage:                    rawTx.PSPReasonMessage,
			PSPReasonType:                       rawTx.PSPReasonType,
			PSPReasonCode:                       rawTx.PSPReasonCode,
			PSPReasonDeclineType:                rawTx.PSPReasonDeclineType,
			TrxCreatedAt:                        trxCreatedAt,
			MetaSubsID:                          rawTx.MetaSubsID,
			MetaSubsIteration:                   rawTx.MetaSubsIteration,
			MetaRetryStep:                       rawTx.MetaRetryStep,
			MetaOneoffID:                        rawTx.MetaOneoffID,
			MetaClientFFProjectID:               rawTx.MetaClientFFProjectID,
			MetaClientFFSessionID:               rawTx.MetaClientFFSessionID,
			MetaClientFFPriceID:                 rawTx.MetaClientFFPriceID,
			PMType:                              rawTx.PMType,
			ThreeDSChallengeIssued:              rawTx.ThreeDSChallengeIssued,
			ThreeDSProtocolVersion:              rawTx.ThreeDSProtocolVersion,
			ThreeDSResponseCode:                 rawTx.ThreeDSResponseCode,
			ThreeDSReasonCode:                   rawTx.ThreeDSReasonCode,
			ThreeDSReasonText:                   rawTx.ThreeDSReasonText,
			AuthorizationType:                   rawTx.AuthorizationType,
			IsVaulted:                           rawTx.IsVaulted,
			PMDataBinAccountFundingType:         rawTx.PMDataBinAccountFundingType,
			PMDataBinAccountNumberType:          rawTx.PMDataBinAccountNumberType,
			PMDataBinIssuerName:                 rawTx.PMDataBinIssuerName,
			PMDataBinIssuerCountryCode:          rawTx.PMDataBinIssuerCountryCode,
			PMDataBinIssuerCurrencyCode:         rawTx.PMDataBinIssuerCurrencyCode,
			Network:                             rawTx.Network,
			PMDataBinPrepaidReloadableIndicator: rawTx.PMDataBinPrepaidReloadableIndicator,
			PMDataBinProductCode:                rawTx.PMDataBinProductCode,
			PMDataBinProductName:                rawTx.PMDataBinProductName,
			PMDataBinProductUsageType:           rawTx.PMDataBinProductUsageType,
			PMDataBinRegionalRestriction:        rawTx.PMDataBinRegionalRestriction,
			PMDataExpirationDate:                rawTx.PMDataExpirationDate,
			PMDataFirst6:                        rawTx.PMDataFirst6,
			PMDataLast4:                         rawTx.PMDataLast4,
			PMDataIsNetworkTokenized:            rawTx.PMDataIsNetworkTokenized,
		})
	}
	return &res
}

// ===== Subscription Management =====

// EnableAutoRenewRequest 启用自动续费请求
type EnableAutoRenewRequest struct {
	ExternalID string  `json:"external_id"`
	SubsID     string  `json:"subs_id"`
	Reason     *string `json:"reason,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

// DisableAutoRenewRequest 禁用自动续费请求
type DisableAutoRenewRequest struct {
	ExternalID string  `json:"external_id"`
	SubsID     string  `json:"subs_id"`
	Reason     *string `json:"reason,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

type MigrationStrategy string

const (
	MigrationStrategyDelayedStart MigrationStrategy = "delayed_start"
	MigrationStrategyPriceProrate MigrationStrategy = "price_prorate"
)

// SubscriptionMigrationRequest 订阅迁移请求
type SubscriptionMigrationRequest struct {
	ExternalID string            `json:"external_id"`
	SubsID     string            `json:"subs_id"`
	PPIdent    string            `json:"pp_ident"`
	Reason     *string           `json:"reason,omitempty"`
	Comment    *string           `json:"comment,omitempty"`
	Strategy   MigrationStrategy `json:"strategy"`
	DryRun     *bool             `json:"dry_run"`
	StrictMode *bool             `json:"strict_mode"`
}

type CheckoutStatus string

const (
	CheckoutStatusProcessing CheckoutStatus = "processing"
	CheckoutStatusSucceeded  CheckoutStatus = "succeeded"
	CheckoutStatusFailed     CheckoutStatus = "failed"
	CheckoutStatusCancelled  CheckoutStatus = "cancelled"
)

// PaymentResult 支付结果
type PaymentResult struct {
	ActionRequiredToken  string         `json:"action_required_token"`
	CheckoutStatus       CheckoutStatus `json:"checkout_status"`
	FailedMessageForUser string         `json:"failed_message_for_user"`
	OrderID              *string        `json:"order_id"`
}

// SubscriptionMigrationResponse 订阅迁移响应
type SubscriptionMigrationResponse struct {
	MigrationStrategy MigrationStrategy `json:"migration_strategy"`
	PaymentResult     PaymentResult     `json:"payment_result"`
	ChargedAmount     any               `json:"charged_amount"`
	SubsID            *string           `json:"subs_id"`
	OneoffID          *string           `json:"oneoff_id"`
}

// DiscountRequest 折扣请求
type DiscountRequest struct {
	ExternalID        string  `json:"external_id"` // 订阅的外部ID
	SubsID            string  `json:"subs_id"`
	Percent           int     `json:"percent"`
	Reason            *string `json:"reason,omitempty"`
	Comment           *string `json:"comment,omitempty"`
	CountOfIterations *int    `json:"count_of_iterations"`
}

// SubscriptionDeferRequest 延迟订阅请求
type SubscriptionDeferRequest struct {
	ExternalID string  `json:"external_id"`
	SubsID     string  `json:"subs_id"`
	DeferTill  string  `json:"defer_till"`
	Reason     *string `json:"reason,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

// SubscriptionPauseRequest 暂停订阅请求
type SubscriptionPauseRequest struct {
	ExternalID string  `json:"external_id"`
	SubsID     string  `json:"subs_id"`
	PauseTill  string  `json:"pause_till"`
	Reason     *string `json:"reason,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

// SubscriptionResumeRequest 恢复订阅请求
type SubscriptionResumeRequest struct {
	ExternalID string  `json:"external_id"` // 订阅的外部ID
	SubsID     string  `json:"subs_id"`
	Reason     *string `json:"reason,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

// ===== PricePoints =====

// PricePointsListRequest 价格点列表请求
type PricePointsListRequest struct {
	Ident *string `json:"ident,omitempty"` // 可选：按标识过滤
}

type PricePointCreateRequest struct {
	Ident                        string             `json:"ident,omitempty"`
	IntroType                    IntroType          `json:"intro_type,omitempty"`
	CurrencyCode                 string             `json:"currency_code,omitempty"`
	LifetimePrice                float64            `json:"lifetime_price,omitempty"`
	IntroFreeTrialPeriod         int                `json:"intro_free_trial_period,omitempty"`
	IntroFreeTrialPeriodDuration PeriodDurationUnit `json:"intro_free_trial_period_duration,omitempty"`
	IntroPaidTrialPrice          float64            `json:"intro_paid_trial_price,omitempty"`
	IntroPaidTrialPeriod         int                `json:"intro_paid_trial_period,omitempty"`
	IntroPaidTrialPeriodDuration PeriodDurationUnit `json:"intro_paid_trial_period_duration,omitempty"`
	NextPrice                    float64            `json:"next_price,omitempty"`
	NextPeriod                   int                `json:"next_period,omitempty"`
	NextPeriodDuration           PeriodDurationUnit `json:"next_period_duration,omitempty"`
	Descriptor                   string             `json:"descriptor,omitempty"`
	Features                     []string           `json:"features,omitempty"`
}

type PricePointCreateResponse struct{}

type FeatureType string

const (
	FeatureTypeTimebased  = "timebased"
	FeatureTypeLifetime   = "lifetime"
	FeatureTypeConsumable = "consumable"
)

type FeatureCreateRequest struct {
	Ident       string      `json:"ident,omitempty"`
	FeatureType FeatureType `json:"feature_type"`
}

type FeatureCreateResponse struct{}

// Currency 货币信息
type Currency struct {
	Code       string `json:"code"`        // 货币代码（如 USD）
	MinorUnits int    `json:"minor_units"` // 小数位数
	Title      string `json:"title"`       // 货币名称
	Symbol     string `json:"symbol"`      // 货币符号
}

// Feature 功能特性
type Feature struct {
	Ident string `json:"ident"`
}

type IntroType string

const (
	IntroTypeNoIntro   IntroType = "no_intro"
	IntroTypeFreeTrial IntroType = "free_trial"
	IntroTypePaidTrial IntroType = "paid_trial"
)

type PeriodDurationUnit string

const (
	PeriodDurationUnitMinutes PeriodDurationUnit = "minutes"
	PeriodDurationUnitDays    PeriodDurationUnit = "days"
	PeriodDurationUnitWeeks   PeriodDurationUnit = "weeks"
	PeriodDurationUnitMonths  PeriodDurationUnit = "months"
	PeriodDurationUnitYears   PeriodDurationUnit = "years"
)

// PricePoint 价格点
type PricePoint struct {
	Ident                        string              `json:"ident"`
	Currency                     Currency            `json:"currency"`
	IntroType                    IntroType           `json:"intro_type"`
	Features                     []Feature           `json:"features"`
	LifetimePrice                *string             `json:"lifetime_price"`
	IntroFreeTrialPeriod         *int                `json:"intro_free_trial_period"`
	IntroFreeTrialPeriodDuration *PeriodDurationUnit `json:"intro_free_trial_period_duration"`
	IntroPaidTrialPrice          *string             `json:"intro_paid_trial_price"`
	IntroPaidTrialPeriod         *int                `json:"intro_paid_trial_period"`
	IntroPaidTrialPeriodDuration *PeriodDurationUnit `json:"intro_paid_trial_period_duration"`
	NextPrice                    *string             `json:"next_price"`
	NextPeriod                   *int                `json:"next_period"`
	NextPeriodDuration           *PeriodDurationUnit `json:"next_period_duration"`
}

// PricePointsListResponse 价格点列表响应
type PricePointsListResponse struct {
	PricePoints []PricePoint `json:"price_points"`
}

// ===== Information =====

// MyAssetsRequest 获取用户资产请求
type MyAssetsRequest struct {
	ExternalID string `json:"external_id"` // 用户外部ID
}

type SubscriptionField struct {
	SubsID               string         `json:"subs_id"`
	IsActive             bool           `json:"is_active"`
	PricePoint           PricePoint     `json:"price_point"`
	Status               []string       `json:"status"`
	AvailableActions     []string       `json:"available_actions"`
	InitialOrderMetadata map[string]any `json:"initial_order_metadata"`
	Iteration            int            `json:"iteration"`
}

// rawSubscription 订阅信息
type rawSubscription struct {
	SubscriptionField     `json:",inline"`
	StartedAt             string `json:"started_at"`
	CurrentPeriodStartsAt string `json:"current_period_starts_at"`
	CurrentPeriodEndsAt   string `json:"current_period_ends_at"`
	NextCheckAt           string `json:"next_check_at"`
}

type Subscription struct {
	SubscriptionField     `json:",inline"`
	StartedAt             *time.Time `json:"started_at"`
	CurrentPeriodStartsAt *time.Time `json:"current_period_starts_at"`
	CurrentPeriodEndsAt   *time.Time `json:"current_period_ends_at"`
	NextCheckAt           *time.Time `json:"next_check_at"`
}

type OneoffField struct {
	OneoffID             string         `json:"oneoff_id"` // 订单ID
	OrderID              string         `json:"order_id"`
	IsActive             bool           `json:"is_active"`
	PricePoint           PricePoint     `json:"price_point"`
	InitialOrderMetadata map[string]any `json:"initial_order_metadata"`
}

// rawOneOffPurchase 一次性购买
type rawOneOffPurchase struct {
	OneoffField `json:",inline"`
	StartedAt   string `json:"started_at"`
	RevokedAt   string `json:"revoked_at"`
}

type OneOffPurchase struct {
	OneoffField `json:",inline"`
	StartedAt   *time.Time `json:"started_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

// rawMyAssetsResponse 用户资产响应
type rawMyAssetsResponse struct {
	Subscriptions   []rawSubscription   `json:"subscriptions"`
	OneOffPurchases []rawOneOffPurchase `json:"oneoffs"`
}

type MyAssetsResponse struct {
	Subscriptions   []Subscription   `json:"subscriptions"`
	OneOffPurchases []OneOffPurchase `json:"oneoffs"`
}

func parseTimePointer(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(timeFormat, s)
	if err == nil {
		return &t
	}
	bs, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	var t1 time.Time
	err = json.Unmarshal(bs, &t1)
	return &t1
}

func (raw rawMyAssetsResponse) toMyAssetsResponse() *MyAssetsResponse {
	var res MyAssetsResponse
	for _, rawSub := range raw.Subscriptions {
		startedAt := parseTimePointer(rawSub.StartedAt)
		currStart := parseTimePointer(rawSub.CurrentPeriodStartsAt)
		currEnd := parseTimePointer(rawSub.CurrentPeriodEndsAt)
		nextCheck := parseTimePointer(rawSub.NextCheckAt)
		res.Subscriptions = append(res.Subscriptions, Subscription{
			SubscriptionField:     rawSub.SubscriptionField,
			StartedAt:             startedAt,
			CurrentPeriodStartsAt: currStart,
			CurrentPeriodEndsAt:   currEnd,
			NextCheckAt:           nextCheck,
		})
	}

	for _, rawOneoff := range raw.OneOffPurchases {
		startedAt := parseTimePointer(rawOneoff.StartedAt)
		revokedAt := parseTimePointer(rawOneoff.RevokedAt)
		res.OneOffPurchases = append(res.OneOffPurchases, OneOffPurchase{
			OneoffField: rawOneoff.OneoffField,
			StartedAt:   startedAt,
			RevokedAt:   revokedAt,
		})
	}
	return &res
}

type OrderStatus string

const (
	OrderStatusAuthorized OrderStatus = "authorized"
	OrderStatusSettled    OrderStatus = "settled"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type OrderField struct {
	OrderID              string         `json:"order_id"`
	Amount               string         `json:"amount"`
	CurrencyCode         string         `json:"currency_code"`
	ExternalID           string         `json:"external_id"`
	SubsID               *string        `json:"subs_id"`
	UserUUID             string         `json:"user_uuid"`
	OneoffID             *string        `json:"oneoff_id"`
	InitialOrderMetadata map[string]any `json:"initial_order_metadata"`
	Status               OrderStatus    `json:"status"`
	DeclineReason        *any           `json:"decline_reason"`
	RetryStep            *any           `json:"retry_step"`
	PSP                  *string        `json:"psp"`
}

type rawOrder struct {
	OrderField `json:",inline"`
	CreatedAt  string `json:"created_at"`
}

// Order 订单信息
type Order struct {
	OrderField `json:",inline"`
	CreatedAt  *time.Time `json:"created_at"`
}

type EventType string

const (
	EventTypeSubscription EventType = "subscription"
	EventTypeOrder        EventType = "order"
	EventTypeRefund       EventType = "refund"
	EventTypeOneoff       EventType = "oneoff"
)

type EventSubtype string

const (
	EventSubtypeSubscriptionStartingTrial                 EventSubtype = "starting_trial"
	EventSubtypeSubscriptionConversion                    EventSubtype = "convertion"
	EventSubtypeSubscriptionRenewing                      EventSubtype = "renewing"
	EventSubtypeSubscriptionUnsubscription                EventSubtype = "unsubscription"
	EventSubtypeSubscriptionPausing                       EventSubtype = "pausing"
	EventSubtypeSubscriptionDeferring                     EventSubtype = "deferring"
	EventSubtypeSubscriptionResuming                      EventSubtype = "resuming"
	EventSubtypeSubscriptionRecoveringAutorenew           EventSubtype = "recovering_autorenew"
	EventSubtypeSubscriptionExpiration                    EventSubtype = "expiration"
	EventSubtypeSubscriptionUnknown                       EventSubtype = "unknown"
	EventSubtypeSubscriptionStartGrace                    EventSubtype = "start_grace"
	EventSubtypeSubscriptionStartRetry                    EventSubtype = "start_retry"
	EventSubtypeSubscriptionFinishGrace                   EventSubtype = "finish_grace"
	EventSubtypeSubscriptionRecovering                    EventSubtype = "recovering"
	EventSubtypeSubscriptionPlanningPostponedSubscription EventSubtype = "planning_postponed_subscription"

	EventSubtypeOrderSettled  EventSubtype = "settled"
	EventSubtypeOrderDeclined EventSubtype = "declined"

	EventSubtypeOneoffGranted EventSubtype = "granted"
	EventSubtypeOneoffRevoked EventSubtype = "revoked"

	EventSubtypeRefundSettled EventSubtype = "settled"
)

// Event 通用事件结构（泛型）
type Event struct {
	EventID        string       `json:"event_id"`
	EventTimestamp time.Time    `json:"event_timestamp"`
	EventType      EventType    `json:"event_type"`
	Subtype        EventSubtype `json:"subtype"`
	ExternalID     *string      `json:"external_id,omitempty"`
	IsLivemode     *bool        `json:"is_livemode,omitempty"`
	*RefundInfo    `json:",inline,omitempty"`
	User           struct {
		Email      string `json:"email"`
		ExternalID string `json:"external_id"`
	} `json:"user"`
	Subscription *Subscription   `json:"subscription"`
	Order        *Order          `json:"order"`
	Oneoff       *OneOffPurchase `json:"oneoff"`
}

type RefundInfoField struct {
	AmountRefunded string `json:"amount_refunded"`
	OrderID        string `json:"order_id"`
	TrxID          string `json:"trx_id"`
	CurrencyCode   string `json:"currency_code"`
}

type RawRefundInfo struct {
	RefundInfoField `json:",inline"`
	CreatedAt       string `json:"created_at"`
}

type RefundInfo struct {
	RefundInfoField `json:",inline"`
	CreatedAt       *time.Time `json:"created_at"`
}

type rawEvent struct {
	EventID        string       `json:"event_id"`
	EventTimestamp string       `json:"event_timestamp"`
	EventType      EventType    `json:"event_type"`
	Subtype        EventSubtype `json:"subtype"`
	ExternalID     *string      `json:"external_id,omitempty"`
	IsLivemode     *bool        `json:"is_livemode,omitempty"`
	*RawRefundInfo `json:",inline,omitempty"`
	User           struct {
		Email      string `json:"email"`
		ExternalID string `json:"external_id"`
	} `json:"user"`
	Subscription *rawSubscription   `json:"subscription"`
	Order        *rawOrder          `json:"order"`
	Oneoff       *rawOneOffPurchase `json:"oneoff"`
}

// parseEventTimestamp 解析事件时间戳
func parseEventTimestamp(s string) (time.Time, error) {
	t, err := time.Parse(timeFormat, s)
	if err == nil {
		return t, nil
	}
	bs, err := json.Marshal(s)
	if err != nil {
		return time.Time{}, err
	}
	var t1 time.Time
	err = json.Unmarshal(bs, &t1)
	return t1, err
}

func ParseEvent(data []byte) (*Event, error) {
	var raw rawEvent
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	event := Event{
		EventID:    raw.EventID,
		EventType:  raw.EventType,
		Subtype:    raw.Subtype,
		ExternalID: raw.ExternalID,
		IsLivemode: raw.IsLivemode,
		User:       raw.User,
	}

	eventTimestamp, err := parseEventTimestamp(raw.EventTimestamp)
	if err != nil {
		return nil, err
	}
	event.EventTimestamp = eventTimestamp

	// 转换 rawSubscription 为 Subscription
	if raw.Subscription != nil {
		event.Subscription = &Subscription{
			SubscriptionField:     raw.Subscription.SubscriptionField,
			StartedAt:             parseTimePointer(raw.Subscription.StartedAt),
			CurrentPeriodStartsAt: parseTimePointer(raw.Subscription.CurrentPeriodStartsAt),
			CurrentPeriodEndsAt:   parseTimePointer(raw.Subscription.CurrentPeriodEndsAt),
			NextCheckAt:           parseTimePointer(raw.Subscription.NextCheckAt),
		}
	}
	if raw.Order != nil {
		event.Order = &Order{
			OrderField: raw.Order.OrderField,
			CreatedAt:  parseTimePointer(raw.Order.CreatedAt),
		}
	}
	if raw.Oneoff != nil {
		event.Oneoff = &OneOffPurchase{
			OneoffField: raw.Oneoff.OneoffField,
			StartedAt:   parseTimePointer(raw.Oneoff.StartedAt),
			RevokedAt:   parseTimePointer(raw.Oneoff.RevokedAt),
		}
	}
	if raw.RawRefundInfo != nil {
		event.RefundInfo = &RefundInfo{
			RefundInfoField: raw.RawRefundInfo.RefundInfoField,
			CreatedAt:       parseTimePointer(raw.RawRefundInfo.CreatedAt),
		}
	}

	return &event, nil
}
