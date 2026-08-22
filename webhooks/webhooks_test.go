package webhooks_test

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/TheWFA/api-client-go/webhooks"
)

type testKeys struct {
	publicKeyPEM string
	privateKey   ed25519.PrivateKey
}

func generateTestKeys(t *testing.T) testKeys {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("MarshalPKIXPublicKey: %v", err)
	}

	block := &pem.Block{Type: "PUBLIC KEY", Bytes: der}

	return testKeys{
		publicKeyPEM: string(pem.EncodeToMemory(block)),
		privateKey:   priv,
	}
}

func (k testKeys) sign(timestamp, body string) string {
	sig := ed25519.Sign(k.privateKey, []byte(timestamp+"."+body))
	return "ed25519=" + base64.StdEncoding.EncodeToString(sig)
}

const pingBody = `{"detailType":"WebhookPing","occurredAt":"2026-08-17T12:00:00.000Z"}`

func TestVerifySignatureAccepts(t *testing.T) {
	keys := generateTestKeys(t)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := keys.sign(timestamp, pingBody)

	err := webhooks.VerifySignature(webhooks.VerifySignatureParams{
		Body:      []byte(pingBody),
		Timestamp: timestamp,
		Signature: signature,
		PublicKey: keys.publicKeyPEM,
	})
	if err != nil {
		t.Fatalf("VerifySignature: %v", err)
	}
}

func TestVerifySignatureRejectsMissingPrefix(t *testing.T) {
	keys := generateTestKeys(t)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	err := webhooks.VerifySignature(webhooks.VerifySignatureParams{
		Body:      []byte(pingBody),
		Timestamp: timestamp,
		Signature: "deadbeef",
		PublicKey: keys.publicKeyPEM,
	})
	if !webhooks.IsSignatureError(err) {
		t.Fatalf("expected a signature error, got %v", err)
	}
}

func TestVerifySignatureRejectsWrongKey(t *testing.T) {
	keys := generateTestKeys(t)
	otherKeys := generateTestKeys(t)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := keys.sign(timestamp, pingBody)

	err := webhooks.VerifySignature(webhooks.VerifySignatureParams{
		Body:      []byte(pingBody),
		Timestamp: timestamp,
		Signature: signature,
		PublicKey: otherKeys.publicKeyPEM,
	})
	if !webhooks.IsSignatureError(err) {
		t.Fatalf("expected a signature error, got %v", err)
	}
}

func TestVerifySignatureRejectsTamperedBody(t *testing.T) {
	keys := generateTestKeys(t)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := keys.sign(timestamp, pingBody)

	err := webhooks.VerifySignature(webhooks.VerifySignatureParams{
		Body:      []byte(strings.Replace(pingBody, "2026", "2099", 1)),
		Timestamp: timestamp,
		Signature: signature,
		PublicKey: keys.publicKeyPEM,
	})
	if !webhooks.IsSignatureError(err) {
		t.Fatalf("expected a signature error, got %v", err)
	}
}

func TestVerifySignatureRejectsStaleTimestamp(t *testing.T) {
	keys := generateTestKeys(t)
	staleTimestamp := strconv.FormatInt(time.Now().Add(-10*time.Minute).UnixMilli(), 10)
	signature := keys.sign(staleTimestamp, pingBody)

	err := webhooks.VerifySignature(webhooks.VerifySignatureParams{
		Body:      []byte(pingBody),
		Timestamp: staleTimestamp,
		Signature: signature,
		PublicKey: keys.publicKeyPEM,
	})
	if !webhooks.IsSignatureError(err) {
		t.Fatalf("expected a signature error, got %v", err)
	}
}

func TestVerifySignatureAcceptsStaleWithWidenedTolerance(t *testing.T) {
	keys := generateTestKeys(t)
	staleTimestamp := strconv.FormatInt(time.Now().Add(-10*time.Minute).UnixMilli(), 10)
	signature := keys.sign(staleTimestamp, pingBody)

	err := webhooks.VerifySignature(webhooks.VerifySignatureParams{
		Body:      []byte(pingBody),
		Timestamp: staleTimestamp,
		Signature: signature,
		PublicKey: keys.publicKeyPEM,
		Tolerance: 24 * time.Hour,
	})
	if err != nil {
		t.Fatalf("VerifySignature: %v", err)
	}
}

func TestParsePayloadPing(t *testing.T) {
	event, err := webhooks.ParsePayload([]byte(pingBody))
	if err != nil {
		t.Fatalf("ParsePayload: %v", err)
	}

	ping, ok := event.(webhooks.PingEvent)
	if !ok {
		t.Fatalf("event is %T, want PingEvent", event)
	}

	if ping.EventType() != webhooks.EventTypePing {
		t.Errorf("EventType() = %v, want %v", ping.EventType(), webhooks.EventTypePing)
	}
}

func TestParsePayloadGoalScoredWithResolvedRefs(t *testing.T) {
	body := `{
		"detailType": "GoalScored",
		"match": {
			"id": 1,
			"status": "second-half",
			"scheduledFor": "2026-08-17T10:00:00.000Z",
			"competition": {"id": 10, "name": "Premier League"},
			"homeTeam": {"id": 2, "name": "Home FC", "nickname": "Home", "badgeUrl": "https://x/h.png"},
			"awayTeam": {"id": 3, "name": "Away FC", "nickname": "Away", "badgeUrl": "https://x/a.png"},
			"score": {"home": 1, "away": 0, "homePenalty": 0, "awayPenalty": 0}
		},
		"team": {"id": 2, "name": "Home FC", "nickname": "Home", "badgeUrl": "https://x/h.png"},
		"scorer": {"id": 5, "name": "J. Smith"},
		"assister": null,
		"goalType": "goal",
		"isPenalty": false,
		"matchPeriod": "first-half",
		"matchTime": 12,
		"occurredAt": "2026-08-17T12:00:00.000Z"
	}`

	event, err := webhooks.ParsePayload([]byte(body))
	if err != nil {
		t.Fatalf("ParsePayload: %v", err)
	}

	goal, ok := event.(webhooks.GoalScoredEvent)
	if !ok {
		t.Fatalf("event is %T, want GoalScoredEvent", event)
	}

	if !goal.Match.Resolved() || goal.Match.Match.ID != 1 {
		t.Errorf("expected a resolved match, got %+v", goal.Match)
	}

	if goal.Scorer == nil || !goal.Scorer.Resolved() || goal.Scorer.Person.Name != "J. Smith" {
		t.Errorf("expected a resolved scorer, got %+v", goal.Scorer)
	}

	if goal.Assister != nil {
		t.Errorf("expected a nil assister, got %+v", goal.Assister)
	}
}

func TestParsePayloadFallsBackToBareID(t *testing.T) {
	body := `{
		"detailType": "CardIssued",
		"match": "1",
		"team": {"id": 2, "name": "Home FC", "nickname": "Home", "badgeUrl": "https://x/h.png"},
		"player": "5",
		"cardType": "yellow_card",
		"matchPeriod": "second-half",
		"matchTime": 60,
		"occurredAt": "2026-08-17T12:00:00.000Z"
	}`

	event, err := webhooks.ParsePayload([]byte(body))
	if err != nil {
		t.Fatalf("ParsePayload: %v", err)
	}

	card, ok := event.(webhooks.CardIssuedEvent)
	if !ok {
		t.Fatalf("event is %T, want CardIssuedEvent", event)
	}

	if card.Match.Resolved() || card.Match.ID != 1 {
		t.Errorf("expected an unresolved match with ID 1, got %+v", card.Match)
	}

	if card.Player.Resolved() || card.Player.ID != 5 {
		t.Errorf("expected an unresolved player with ID 5, got %+v", card.Player)
	}
}

func TestParsePayloadInvalidJSON(t *testing.T) {
	_, err := webhooks.ParsePayload([]byte("not json"))
	if !webhooks.IsPayloadError(err) {
		t.Fatalf("expected a payload error, got %v", err)
	}
}

func TestParsePayloadUnknownEventType(t *testing.T) {
	_, err := webhooks.ParsePayload([]byte(`{"detailType":"SomethingElse"}`))
	if !webhooks.IsPayloadError(err) {
		t.Fatalf("expected a payload error, got %v", err)
	}
}

func TestExtractHeaders(t *testing.T) {
	header := http.Header{}
	header.Set("X-WFA-Signature", "ed25519=abc")
	header.Set("X-WFA-Timestamp", "123")
	header.Set("X-WFA-Event-Type", "WebhookPing")
	header.Set("X-WFA-Delivery-Id", "delivery-1")

	headers, err := webhooks.ExtractHeaders(header)
	if err != nil {
		t.Fatalf("ExtractHeaders: %v", err)
	}

	if headers != (webhooks.Headers{
		Signature:  "ed25519=abc",
		Timestamp:  "123",
		EventType:  "WebhookPing",
		DeliveryID: "delivery-1",
	}) {
		t.Errorf("unexpected headers: %+v", headers)
	}
}

func TestExtractHeadersMissing(t *testing.T) {
	header := http.Header{}
	header.Set("X-WFA-Timestamp", "123")
	header.Set("X-WFA-Event-Type", "WebhookPing")
	header.Set("X-WFA-Delivery-Id", "delivery-1")

	_, err := webhooks.ExtractHeaders(header)
	if !webhooks.IsPayloadError(err) {
		t.Fatalf("expected a payload error, got %v", err)
	}
}

func TestConstructEventFromHeaders(t *testing.T) {
	keys := generateTestKeys(t)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := keys.sign(timestamp, pingBody)

	header := http.Header{}
	header.Set("X-WFA-Signature", signature)
	header.Set("X-WFA-Timestamp", timestamp)
	header.Set("X-WFA-Event-Type", "WebhookPing")
	header.Set("X-WFA-Delivery-Id", "delivery-1")

	event, err := webhooks.ConstructEventFromHeaders([]byte(pingBody), header, keys.publicKeyPEM, nil)
	if err != nil {
		t.Fatalf("ConstructEventFromHeaders: %v", err)
	}

	if event.EventType() != webhooks.EventTypePing {
		t.Errorf("EventType() = %v, want %v", event.EventType(), webhooks.EventTypePing)
	}
}

func TestConstructEventFromHeadersInvalidSignature(t *testing.T) {
	keys := generateTestKeys(t)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	header := http.Header{}
	header.Set("X-WFA-Signature", "ed25519=d3Jvbmc=")
	header.Set("X-WFA-Timestamp", timestamp)
	header.Set("X-WFA-Event-Type", "WebhookPing")
	header.Set("X-WFA-Delivery-Id", "delivery-1")

	_, err := webhooks.ConstructEventFromHeaders([]byte(pingBody), header, keys.publicKeyPEM, nil)
	if !webhooks.IsSignatureError(err) {
		t.Fatalf("expected a signature error, got %v", err)
	}
}

func TestConstructEventFromRequest(t *testing.T) {
	keys := generateTestKeys(t)
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := keys.sign(timestamp, pingBody)

	req := httptest.NewRequest(http.MethodPost, "https://example.com/webhooks/wfa", strings.NewReader(pingBody))
	req.Header.Set("X-WFA-Signature", signature)
	req.Header.Set("X-WFA-Timestamp", timestamp)
	req.Header.Set("X-WFA-Event-Type", "WebhookPing")
	req.Header.Set("X-WFA-Delivery-Id", "delivery-1")

	event, err := webhooks.ConstructEvent(req, keys.publicKeyPEM, nil)
	if err != nil {
		t.Fatalf("ConstructEvent: %v", err)
	}

	if event.EventType() != webhooks.EventTypePing {
		t.Errorf("EventType() = %v, want %v", event.EventType(), webhooks.EventTypePing)
	}
}
