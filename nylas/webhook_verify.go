package nylas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidSignature = errors.New("invalid webhook signature")
var ErrOldTimestamp = errors.New("webhook timestamp too old")

// VerifyWebhookHMAC checks a signature header using HMAC SHA-256.
// Supports formats: "sha256=<hex>", raw hex, or base64.
// The payload hashed is the body bytes by default. If timestamp is provided, we hash timestamp + "." + body.
func VerifyWebhookHMAC(secret string, body []byte, signatureHeader string, timestamp *string, tolerance time.Duration) error {
	mac := hmac.New(sha256.New, []byte(secret))
	if timestamp != nil {
		mac.Write([]byte(*timestamp))
		mac.Write([]byte("."))
	}
	mac.Write(body)
	expected := mac.Sum(nil)

	// normalize provided signature
	sig := signatureHeader
	if strings.HasPrefix(strings.ToLower(sig), "sha256=") {
		sig = sig[len("sha256="):]
	}
	// try hex
	if bs, err := hex.DecodeString(sig); err == nil {
		if hmac.Equal(expected, bs) {
			return nil
		}
	}
	// try base64
	if bs, err := base64.StdEncoding.DecodeString(sig); err == nil {
		if hmac.Equal(expected, bs) {
			return nil
		}
	}
	return ErrInvalidSignature
}

// Verify timestamp staleness if a timestamp header is present.
func VerifyWebhookTimestamp(ts string, tolerance time.Duration) error {
	if ts == "" {
		return nil
	}
	// accept integer seconds epoch or RFC3339
	if secs, err := strconv.ParseInt(ts, 10, 64); err == nil {
		when := time.Unix(secs, 0)
		if time.Since(when) > tolerance {
			return ErrOldTimestamp
		}
		return nil
	}
	if when, err := time.Parse(time.RFC3339, ts); err == nil {
		if time.Since(when) > tolerance {
			return ErrOldTimestamp
		}
		return nil
	}
	return nil
}
