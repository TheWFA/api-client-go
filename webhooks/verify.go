package webhooks

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ErrorKind distinguishes the two ways constructing a webhook event can
// fail.
type ErrorKind int

const (
	// ErrorKindSignature means the delivery failed signature or timestamp
	// verification.
	ErrorKindSignature ErrorKind = iota
	// ErrorKindPayload means the delivery's headers or body were malformed
	// or unrecognized.
	ErrorKindPayload
)

// Error is returned by VerifySignature, ExtractHeaders, ParsePayload,
// ConstructEvent and ConstructEventFromHeaders. Check Kind, or use
// IsSignatureError / IsPayloadError, to tell what went wrong.
type Error struct {
	Kind    ErrorKind
	Message string
}

func (e *Error) Error() string { return e.Message }

func errorKindIs(err error, kind ErrorKind) bool {
	var webhookErr *Error
	return errors.As(err, &webhookErr) && webhookErr.Kind == kind
}

// IsSignatureError reports whether err is an *Error caused by a signature or
// timestamp verification failure.
func IsSignatureError(err error) bool { return errorKindIs(err, ErrorKindSignature) }

// IsPayloadError reports whether err is an *Error caused by a malformed or
// unrecognized delivery.
func IsPayloadError(err error) bool { return errorKindIs(err, ErrorKindPayload) }

const (
	signaturePrefix         = "ed25519="
	defaultToleranceSeconds = 300
	headerSignature         = "X-WFA-Signature"
	headerTimestamp         = "X-WFA-Timestamp"
	headerEventType         = "X-WFA-Event-Type"
	headerDeliveryID        = "X-WFA-Delivery-Id"
)

// Headers holds the four X-WFA-* delivery headers sent with every webhook
// request.
type Headers struct {
	Signature  string
	Timestamp  string
	EventType  string
	DeliveryID string
}

// ExtractHeaders pulls the four X-WFA-* delivery headers off an inbound
// webhook request.
func ExtractHeaders(header http.Header) (Headers, error) {
	h := Headers{
		Signature:  header.Get(headerSignature),
		Timestamp:  header.Get(headerTimestamp),
		EventType:  header.Get(headerEventType),
		DeliveryID: header.Get(headerDeliveryID),
	}

	var missing []string

	if h.Signature == "" {
		missing = append(missing, headerSignature)
	}

	if h.Timestamp == "" {
		missing = append(missing, headerTimestamp)
	}

	if h.EventType == "" {
		missing = append(missing, headerEventType)
	}

	if h.DeliveryID == "" {
		missing = append(missing, headerDeliveryID)
	}

	if len(missing) > 0 {
		return Headers{}, &Error{
			Kind:    ErrorKindPayload,
			Message: fmt.Sprintf("wfa: missing required webhook header(s): %s", strings.Join(missing, ", ")),
		}
	}

	return h, nil
}

// VerifySignatureParams holds the inputs to VerifySignature.
type VerifySignatureParams struct {
	// Body is the raw request body, exactly as received — must not be
	// re-serialized JSON.
	Body []byte
	// Timestamp is the value of the X-WFA-Timestamp header (ms since epoch,
	// as a string).
	Timestamp string
	// Signature is the value of the X-WFA-Signature header, e.g.
	// "ed25519=<base64>".
	Signature string
	// PublicKey is the subscription's Ed25519 public key, SPKI PEM format.
	PublicKey string
	// Tolerance is the maximum allowed age of the timestamp. Zero (the
	// default) uses 5 minutes.
	Tolerance time.Duration
}

func parseEd25519PublicKeyPEM(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("not a valid PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	edPub, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("key is not an Ed25519 public key")
	}

	return edPub, nil
}

// VerifySignature verifies a webhook delivery's Ed25519 signature and
// timestamp freshness. It returns an *Error with Kind == ErrorKindSignature
// if either check fails.
func VerifySignature(params VerifySignatureParams) error {
	tolerance := params.Tolerance
	if tolerance <= 0 {
		tolerance = defaultToleranceSeconds * time.Second
	}

	if !strings.HasPrefix(params.Signature, signaturePrefix) {
		return &Error{
			Kind:    ErrorKindSignature,
			Message: fmt.Sprintf("wfa: signature must be prefixed with %q", signaturePrefix),
		}
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(params.Signature, signaturePrefix))
	if err != nil {
		return &Error{
			Kind:    ErrorKindSignature,
			Message: fmt.Sprintf("wfa: could not decode signature: %v", err),
		}
	}

	publicKey, err := parseEd25519PublicKeyPEM(params.PublicKey)
	if err != nil {
		return &Error{
			Kind:    ErrorKindSignature,
			Message: fmt.Sprintf("wfa: could not parse public key: %v", err),
		}
	}

	message := append([]byte(params.Timestamp+"."), params.Body...)

	if !ed25519.Verify(publicKey, message, signatureBytes) {
		return &Error{
			Kind:    ErrorKindSignature,
			Message: "wfa: signature does not match expected value",
		}
	}

	timestampMs, err := strconv.ParseInt(params.Timestamp, 10, 64)
	if err != nil {
		return &Error{
			Kind:    ErrorKindSignature,
			Message: "wfa: timestamp header is not a valid number",
		}
	}

	age := time.Since(time.UnixMilli(timestampMs))
	if age < 0 {
		age = -age
	}

	if age > tolerance {
		return &Error{
			Kind:    ErrorKindSignature,
			Message: "wfa: timestamp is outside the allowed tolerance",
		}
	}

	return nil
}

// ConstructEventOptions configures ConstructEvent and
// ConstructEventFromHeaders.
type ConstructEventOptions struct {
	// Tolerance is the maximum allowed age of the delivery's timestamp. Zero
	// (the default) uses 5 minutes.
	Tolerance time.Duration
}

// ConstructEventFromHeaders verifies and parses a webhook delivery from its
// raw body and headers. Use this when an *http.Request isn't available —
// e.g. a framework that has already read the raw body and holds headers as
// an http.Header.
func ConstructEventFromHeaders(body []byte, header http.Header, publicKey string, opts *ConstructEventOptions) (Event, error) {
	headers, err := ExtractHeaders(header)
	if err != nil {
		return nil, err
	}

	var tolerance time.Duration
	if opts != nil {
		tolerance = opts.Tolerance
	}

	if err := VerifySignature(VerifySignatureParams{
		Body:      body,
		Timestamp: headers.Timestamp,
		Signature: headers.Signature,
		PublicKey: publicKey,
		Tolerance: tolerance,
	}); err != nil {
		return nil, err
	}

	return ParsePayload(body)
}

// ConstructEvent verifies and parses a webhook delivery directly from an
// inbound *http.Request. It reads and closes the request body.
func ConstructEvent(r *http.Request, publicKey string, opts *ConstructEventOptions) (Event, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("wfa: reading webhook request body: %w", err)
	}

	return ConstructEventFromHeaders(body, r.Header, publicKey, opts)
}
