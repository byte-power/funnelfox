package funnelfox

import (
	"encoding/json"
	"time"
)

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
	Amount     *int    `json:"amount,omitempty"`      // 退款金额（可选，用于部分退款）
	SoftRefund *bool   `json:"soft_refund,omitempty"` // 是否软退款（可选）
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
	MigrationStrategyDelayedStart = "delayed_start"
	MigrationStrategyPriceProrate = "price_prorate"
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
	LifetimePrice                *string             `json:"lifetime_price,omitempty"`
	IntroFreeTrialPeriod         *int                `json:"intro_free_trial_period,omitempty"`
	IntroFreeTrialPeriodDuration *PeriodDurationUnit `json:"intro_free_trial_period_duration,omitempty"`
	IntroPaidTrialPrice          *string             `json:"intro_paid_trial_price,omitempty"`
	IntroPaidTrialPeriod         *int                `json:"intro_paid_trial_period,omitempty"`
	IntroPaidTrialPeriodDuration *PeriodDurationUnit `json:"intro_paid_trial_period_duration,omitempty"`
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

type subscriptionField struct {
	SubsID               string         `json:"subs_id"`
	IsActive             bool           `json:"is_active"`
	PricePoint           PricePoint     `json:"price_point"`
	Status               []string       `json:"status"`
	AvailableActions     []string       `json:"available_actions"`
	InitialOrderMetadata map[string]any `json:"initial_order_metadata"`
}

// rawSubscription 订阅信息
type rawSubscription struct {
	subscriptionField     `json:",inline"`
	StartedAt             string `json:"started_at"`
	CurrentPeriodStartsAt string `json:"current_period_starts_at"`
	CurrentPeriodEndsAt   string `json:"current_period_ends_at"`
	NextCheckAt           string `json:"next_check_at"`
}

type Subscription struct {
	subscriptionField     `json:",inline"`
	StartedAt             *time.Time `json:"started_at"`
	CurrentPeriodStartsAt *time.Time `json:"current_period_starts_at"`
	CurrentPeriodEndsAt   *time.Time `json:"current_period_ends_at"`
	NextCheckAt           *time.Time `json:"next_check_at"`
}

type oneoffField struct {
	OneoffID             string         `json:"oneoff_id"` // 订单ID
	IsActive             bool           `json:"is_active"`
	PricePoint           PricePoint     `json:"price_point"`
	InitialOrderMetadata map[string]any `json:"initial_order_metadata"`
}

// rawOneOffPurchase 一次性购买
type rawOneOffPurchase struct {
	oneoffField `json:",inline"`
	StartedAt   string `json:"started_at"`
	RevokedAt   string `json:"revoked_at"`
}

type OneOffPurchase struct {
	oneoffField `json:",inline"`
	StartedAt   *time.Time `json:"started_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

// rawMyAssetsResponse 用户资产响应
type rawMyAssetsResponse struct {
	Subscriptions   []rawSubscription   `json:"subscriptions"`
	OneOffPurchases []rawOneOffPurchase `json:"one_off_purchases"`
}

type MyAssetsResponse struct {
	Subscriptions   []Subscription   `json:"subscriptions"`
	OneOffPurchases []OneOffPurchase `json:"one_off_purchases"`
}

func parseTimePointer(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return nil
	}
	return &t
}

func (raw rawMyAssetsResponse) toMyAssetsResponse() *MyAssetsResponse {
	var res MyAssetsResponse
	for _, rawSub := range raw.Subscriptions {
		startedAt := parseTimePointer(rawSub.StartedAt)
		currStart := parseTimePointer(rawSub.CurrentPeriodStartsAt)
		currEnd := parseTimePointer(rawSub.CurrentPeriodEndsAt)
		nextCheck := parseTimePointer(rawSub.NextCheckAt)
		res.Subscriptions = append(res.Subscriptions, Subscription{
			subscriptionField:     rawSub.subscriptionField,
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
			oneoffField: rawOneoff.oneoffField,
			StartedAt:   startedAt,
			RevokedAt:   revokedAt,
		})
	}
	return &res
}
