package funnelfox

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestTransactionReportPreservesMetaClientAndRawMessage(t *testing.T) {
	var raw rawTransactionReportResponse
	payload := []byte(`{
		"transactions": [{
			"order_id": "ord_1",
			"meta_client": {"source": "ios"},
			"psp_date": "",
			"trx_created_at": "2026-01-02T03:04:05.000000"
		}]
	}`)

	if err := json.Unmarshal(payload, &raw); err != nil {
		t.Fatalf("unmarshal raw transaction report: %v", err)
	}

	resp := raw.toTransactionReportResponse()
	if len(resp.Transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(resp.Transactions))
	}

	tx := resp.Transactions[0]
	if tx.MetaClient["source"] != "ios" {
		t.Fatalf("expected meta_client source to be copied, got %#v", tx.MetaClient)
	}
	assertRawMessageContains(t, tx.RawMessage(), []byte(`"meta_client"`))
	assertRawMessageClone(t, tx.RawMessage, []byte(`"meta_client"`))
}

func TestParseEventPreservesRawMessages(t *testing.T) {
	payload := []byte(`{
		"event_id": "evt_1",
		"event_timestamp": "2026-01-02T03:04:05.000000",
		"event_type": "subscription",
		"subtype": "renewing",
		"user": {
			"email": "user@example.com",
			"external_id": "ext_1"
		},
		"subscription": {
			"subs_id": "sub_1",
			"started_at": "2026-01-02T03:04:05.000000",
			"current_period_starts_at": "",
			"current_period_ends_at": "",
			"next_check_at": ""
		},
		"order": {
			"order_id": "ord_1",
			"amount": "10.00",
			"currency_code": "USD",
			"external_id": "ext_1",
			"user_uuid": "user_1",
			"status": "settled",
			"created_at": "2026-01-02T03:04:05.000000"
		},
		"oneoff": {
			"oneoff_id": "one_1",
			"order_id": "ord_2",
			"is_active": true,
			"started_at": "2026-01-02T03:04:05.000000",
			"revoked_at": ""
		}
	}`)

	event, err := ParseEvent(payload)
	if err != nil {
		t.Fatalf("parse event: %v", err)
	}
	if event.Subscription == nil {
		t.Fatal("expected subscription")
	}
	if event.Order == nil {
		t.Fatal("expected order")
	}
	if event.Oneoff == nil {
		t.Fatal("expected oneoff")
	}

	assertRawMessageContains(t, event.Subscription.RawMessage(), []byte(`"subs_id"`))
	assertRawMessageContains(t, event.Order.RawMessage(), []byte(`"order_id"`))
	assertRawMessageContains(t, event.Oneoff.RawMessage(), []byte(`"oneoff_id"`))
	assertRawMessageClone(t, event.Subscription.RawMessage, []byte(`"subs_id"`))
	assertRawMessageClone(t, event.Order.RawMessage, []byte(`"order_id"`))
	assertRawMessageClone(t, event.Oneoff.RawMessage, []byte(`"oneoff_id"`))

	expectedCreatedAt, err := time.Parse(timeFormat, "2026-01-02T03:04:05.000000")
	if err != nil {
		t.Fatalf("parse expected created_at: %v", err)
	}
	if event.Order.CreatedAt == nil {
		t.Fatal("expected order created_at to be set")
	}
	if !event.Order.CreatedAt.Equal(expectedCreatedAt) {
		t.Fatalf("expected order created_at %v, got %v", expectedCreatedAt, *event.Order.CreatedAt)
	}
}

func TestOrderUnmarshalJSONCreatedAt(t *testing.T) {
	payload := []byte(`{
		"order_id": "ord_1",
		"amount": "10.00",
		"currency_code": "USD",
		"external_id": "ext_1",
		"user_uuid": "user_1",
		"status": "settled",
		"created_at": "2026-01-02T03:04:05.000000"
	}`)

	var order Order
	if err := json.Unmarshal(payload, &order); err != nil {
		t.Fatalf("unmarshal order: %v", err)
	}

	expectedCreatedAt, err := time.Parse(timeFormat, "2026-01-02T03:04:05.000000")
	if err != nil {
		t.Fatalf("parse expected created_at: %v", err)
	}
	if order.CreatedAt == nil {
		t.Fatal("expected order created_at to be set")
	}
	if !order.CreatedAt.Equal(expectedCreatedAt) {
		t.Fatalf("expected order created_at %v, got %v", expectedCreatedAt, *order.CreatedAt)
	}
	assertRawMessageContains(t, order.RawMessage(), []byte(`"created_at"`))
}

func assertRawMessageContains(t *testing.T, raw json.RawMessage, fragment []byte) {
	t.Helper()
	if !json.Valid(raw) {
		t.Fatalf("expected valid raw JSON, got %q", string(raw))
	}
	if !bytes.Contains(raw, fragment) {
		t.Fatalf("expected raw JSON to contain %q, got %q", string(fragment), string(raw))
	}
}

func assertRawMessageClone(t *testing.T, rawMessage func() json.RawMessage, fragment []byte) {
	t.Helper()
	raw := rawMessage()
	if len(raw) == 0 {
		t.Fatal("expected raw JSON")
	}
	raw[0] = 'x'
	assertRawMessageContains(t, rawMessage(), fragment)
}
