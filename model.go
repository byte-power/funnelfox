package funnelfox

import (
	"encoding/json"
	"fmt"
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
	PaymentResult PaymentResult `json:"payment_result"`
	ChargedAmount any           `json:"charged_amount"`
	SubsID        *string       `json:"subs_id"`
	OneoffID      *string       `json:"oneoff_id"`
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
	t, err := time.Parse(timeFormat, s)
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

// ===== Events =====

// Order 订单信息
type Order struct {
	OrderID              string         `json:"order_id"`
	Amount               string         `json:"amount"`
	CurrencyCode         string         `json:"currency_code"`
	ExternalID           string         `json:"external_id"`
	SubsID               string         `json:"subs_id"`
	UserUUID             string         `json:"user_uuid"`
	OneoffID             *string        `json:"oneoff_id"`
	InitialOrderMetadata map[string]any `json:"initial_order_metadata"`
}

// EventPayload 事件载荷接口
type EventPayload interface {
	SubscriptionEventPayload | OrderEventPayload
}

// SubscriptionEventPayload 订阅事件载荷
type SubscriptionEventPayload struct {
	Subscription Subscription `json:"subscription"`
}

// OrderEventPayload 订单事件载荷
type OrderEventPayload struct {
	Order Order `json:"order"`
}

type EventType string

const (
	EventTypeSubscription EventType = "subscription"
	EventTypeOrder        EventType = "order"
)

type EventSubtype string

const (
	EventSubtypeSubscriptionStartingTrial       EventSubtype = "starting_trial"
	EventSubtypeSubscriptionConversion          EventSubtype = "conversion"
	EventSubtypeSubscriptionRenewing            EventSubtype = "renewing"
	EventSubtypeSubscriptionUnsubscribed        EventSubtype = "unsubscribed"
	EventSubtypeSubscriptionPausing             EventSubtype = "pausing"
	EventSubtypeSubscriptionDeferring           EventSubtype = "deferring"
	EventSubtypeSubscriptionResuming            EventSubtype = "resuming"
	EventSubtypeSubscriptionRecoveringAutorenew EventSubtype = "recovering_autorenew"
	EventSubtypeSubscriptionExpiration          EventSubtype = "expiration"
	EventSubtypeSubscriptionUnknown             EventSubtype = "unknown"
	EventSubtypeSubscriptionStartGrace          EventSubtype = "start_grace"
	EventSubtypeSubscriptionStartRetry          EventSubtype = "start_retry"
	EventSubtypeSubscriptionFinishGrace         EventSubtype = "finish_grace"
	EventSubtypeSubscriptionRecovering          EventSubtype = "recovering"

	EventSubtypeOrderSettled  EventSubtype = "settled"
	EventSubtypeOrderDeclined EventSubtype = "declined"

	EventSubtypeOneoffGranted EventSubtype = "granted"
	EventSubtypeOneoffRevoked EventSubtype = "revoked"

	EventSubtypeRefund EventSubtype = "refund"
)

// Event 通用事件结构（泛型）
type Event[T EventPayload] struct {
	EventID        string       `json:"event_id"`
	EventTimestamp time.Time    `json:"event_timestamp"`
	EventType      EventType    `json:"event_type"`
	Subtype        EventSubtype `json:"subtype"`
	ExternalID     *string      `json:"external_id,omitempty"`
	IsLivemode     *bool        `json:"is_livemode,omitempty"`
	Payload        T            `json:"-"` // Not serialized directly, inlined via custom MarshalJSON
}

// SubscriptionEvent 订阅事件
type SubscriptionEvent Event[SubscriptionEventPayload]

// MarshalJSON inlines the payload fields into the event JSON
func (e SubscriptionEvent) MarshalJSON() ([]byte, error) {
	aux := &struct {
		EventID        string       `json:"event_id"`
		EventTimestamp time.Time    `json:"event_timestamp"`
		EventType      EventType    `json:"event_type"`
		Subtype        EventSubtype `json:"subtype"`
		ExternalID     *string      `json:"external_id,omitempty"`
		IsLivemode     *bool        `json:"is_livemode,omitempty"`
		Subscription   Subscription `json:"subscription"`
	}{
		EventID:        e.EventID,
		EventTimestamp: e.EventTimestamp,
		EventType:      e.EventType,
		Subtype:        e.Subtype,
		ExternalID:     e.ExternalID,
		IsLivemode:     e.IsLivemode,
		Subscription:   e.Payload.Subscription,
	}
	return json.Marshal(aux)
}

// OrderEvent 订单事件
type OrderEvent Event[OrderEventPayload]

// MarshalJSON inlines the payload fields into the event JSON
func (e OrderEvent) MarshalJSON() ([]byte, error) {
	aux := &struct {
		EventID        string       `json:"event_id"`
		EventTimestamp time.Time    `json:"event_timestamp"`
		EventType      EventType    `json:"event_type"`
		Subtype        EventSubtype `json:"subtype"`
		ExternalID     *string      `json:"external_id,omitempty"`
		IsLivemode     *bool        `json:"is_livemode,omitempty"`
		Order          Order        `json:"order"`
	}{
		EventID:        e.EventID,
		EventTimestamp: e.EventTimestamp,
		EventType:      e.EventType,
		Subtype:        e.Subtype,
		ExternalID:     e.ExternalID,
		IsLivemode:     e.IsLivemode,
		Order:          e.Payload.Order,
	}
	return json.Marshal(aux)
}

// rawEventBase 原始事件基础结构（用于解析）
type rawEventBase struct {
	EventID        string          `json:"event_id"`
	EventTimestamp string          `json:"event_timestamp"`
	EventType      string          `json:"event_type"`
	Subtype        EventSubtype    `json:"subtype"`
	ExternalID     *string         `json:"external_id,omitempty"`
	IsLivemode     *bool           `json:"is_livemode,omitempty"`
	Payload        json.RawMessage `json:"-"` // Will be extracted based on event_type
}

// rawSubscriptionEvent 原始订阅事件（用于解析）
type rawSubscriptionEvent struct {
	rawEventBase
	Subscription rawSubscription `json:"subscription"`
}

// rawOrderEvent 原始订单事件（用于解析）
type rawOrderEvent struct {
	rawEventBase
	Order struct {
		Amount               string         `json:"amount"`
		CurrencyCode         string         `json:"currency_code"`
		ExternalID           string         `json:"external_id"`
		InitialOrderMetadata map[string]any `json:"initial_order_metadata"`
		OneoffID             *string        `json:"oneoff_id"`
		OrderID              string         `json:"order_id"`
		SubsID               string         `json:"subs_id"`
		UserUUID             string         `json:"user_uuid"`
	} `json:"order"`
}

// parseEventTimestamp 解析事件时间戳
func parseEventTimestamp(s string) (time.Time, error) {
	return time.Parse(timeFormat, s)
}

// parseSubscriptionEvent 解析订阅事件
func parseSubscriptionEvent(data []byte) (*SubscriptionEvent, error) {
	var raw rawSubscriptionEvent
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	eventTimestamp, err := parseEventTimestamp(raw.EventTimestamp)
	if err != nil {
		return nil, err
	}

	// 转换 rawSubscription 为 Subscription
	startedAt := parseTimePointer(raw.Subscription.StartedAt)
	currStart := parseTimePointer(raw.Subscription.CurrentPeriodStartsAt)
	currEnd := parseTimePointer(raw.Subscription.CurrentPeriodEndsAt)
	nextCheck := parseTimePointer(raw.Subscription.NextCheckAt)

	subscription := Subscription{
		subscriptionField:     raw.Subscription.subscriptionField,
		StartedAt:             startedAt,
		CurrentPeriodStartsAt: currStart,
		CurrentPeriodEndsAt:   currEnd,
		NextCheckAt:           nextCheck,
	}

	event := SubscriptionEvent{
		EventID:        raw.EventID,
		EventTimestamp: eventTimestamp,
		EventType:      EventType(raw.EventType),
		Subtype:        raw.Subtype,
		ExternalID:     raw.ExternalID,
		IsLivemode:     raw.IsLivemode,
		Payload: SubscriptionEventPayload{
			Subscription: subscription,
		},
	}
	return &event, nil
}

// parseOrderEvent 解析订单事件
func parseOrderEvent(data []byte) (*OrderEvent, error) {
	var raw rawOrderEvent
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	eventTimestamp, err := parseEventTimestamp(raw.EventTimestamp)
	if err != nil {
		return nil, err
	}

	order := Order{
		OrderID:              raw.Order.OrderID,
		Amount:               raw.Order.Amount,
		CurrencyCode:         raw.Order.CurrencyCode,
		ExternalID:           raw.Order.ExternalID,
		SubsID:               raw.Order.SubsID,
		UserUUID:             raw.Order.UserUUID,
		OneoffID:             raw.Order.OneoffID,
		InitialOrderMetadata: raw.Order.InitialOrderMetadata,
	}

	event := OrderEvent{
		EventID:        raw.EventID,
		EventTimestamp: eventTimestamp,
		EventType:      EventType(raw.EventType),
		Subtype:        raw.Subtype,
		ExternalID:     raw.ExternalID,
		IsLivemode:     raw.IsLivemode,
		Payload: OrderEventPayload{
			Order: order,
		},
	}
	return &event, nil
}

// ParseEvent 泛型函数，解析 []byte 为指定类型的事件对象
// 示例：event, err := ParseEvent[SubscriptionEvent](data)
func ParseEvent[T SubscriptionEvent | OrderEvent](data []byte) (*T, error) {
	var base rawEventBase
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	var result any
	var err error

	switch base.EventType {
	case "subscription":
		var subEvent *SubscriptionEvent
		subEvent, err = parseSubscriptionEvent(data)
		if err != nil {
			return nil, err
		}
		result = subEvent
	case "order":
		var orderEvent *OrderEvent
		orderEvent, err = parseOrderEvent(data)
		if err != nil {
			return nil, err
		}
		result = orderEvent
	default:
		return nil, fmt.Errorf("unknown event_type: %s", base.EventType)
	}

	// 类型断言并验证
	if typed, ok := result.(*T); ok {
		return typed, nil
	}

	// 类型不匹配的情况
	var zeroPtr *T
	return zeroPtr, fmt.Errorf("event_type %s does not match requested type", base.EventType)
}

// ParseEventAuto 自动推断事件类型并解析
// 返回类型为 any，需要使用类型断言
func ParseEventAuto(data []byte) (any, error) {
	var base rawEventBase
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.EventType {
	case "subscription":
		return parseSubscriptionEvent(data)
	case "order":
		return parseOrderEvent(data)
	default:
		return nil, fmt.Errorf("unknown event_type: %s", base.EventType)
	}
}
