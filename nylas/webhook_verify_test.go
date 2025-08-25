package nylas

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestVerifyWebhookHMAC_Hex_NoTimestamp(t *testing.T) {
	secret := "s3cr3t"
	body := []byte(`{"ok":true}`)

	m := hmac.New(sha256.New, []byte(secret))
	m.Write(body)
	sum := m.Sum(nil)
	sig := hex.EncodeToString(sum)

	if err := VerifyWebhookHMAC(secret, body, sig, nil, 0); err != nil {
		t.Fatalf("want ok, got err=%v", err)
	}
}

func TestVerifyWebhookHMAC_Base64_WithTimestamp(t *testing.T) {
	secret := "s3cr3t"
	body := []byte("abc")
	ts := time.Now().UTC().Format(time.RFC3339)

	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte(ts))
	m.Write([]byte("."))
	m.Write(body)
	sig := base64.StdEncoding.EncodeToString(m.Sum(nil))

	if err := VerifyWebhookHMAC(secret, body, sig, &ts, 0); err != nil {
		t.Fatalf("want ok, got err=%v", err)
	}
	// Header form "sha256=<hex>" should also be accepted
	m2 := hmac.New(sha256.New, []byte(secret))
	m2.Write([]byte(ts))
	m2.Write([]byte("."))
	m2.Write(body)
	hexSig := "sha256=" + hex.EncodeToString(m2.Sum(nil))
	if err := VerifyWebhookHMAC(secret, body, hexSig, &ts, 0); err != nil {
		t.Fatalf("want ok, got err=%v", err)
	}
}

func TestVerifyWebhookHMAC_BadSig(t *testing.T) {
	secret := "s3cr3t"
	body := []byte("x")
	if err := VerifyWebhookHMAC(secret, body, "deadbeef", nil, 0); err == nil {
		t.Fatalf("want error for bad signature")
	}
}

func TestVerifyWebhookTimestamp(t *testing.T) {
	// Empty: allowed
	if err := VerifyWebhookTimestamp("", time.Minute); err != nil {
		t.Fatalf("empty timestamp should be ok, got %v", err)
	}

	now := time.Now().UTC()
	epoch := now.Unix()
	if err := VerifyWebhookTimestamp(strconv64(epoch), time.Minute); err != nil {
		t.Fatalf("epoch timestamp ok, got %v", err)
	}
	if err := VerifyWebhookTimestamp(now.Format(time.RFC3339), time.Minute); err != nil {
		t.Fatalf("rfc3339 timestamp ok, got %v", err)
	}

	// Old (beyond tolerance)
	old := now.Add(-2 * time.Hour)
	if err := VerifyWebhookTimestamp(strconv64(old.Unix()), time.Minute); err == nil {
		t.Fatalf("want ErrOldTimestamp for old epoch")
	}
	if err := VerifyWebhookTimestamp(old.Format(time.RFC3339), time.Minute); err == nil {
		t.Fatalf("want ErrOldTimestamp for old rfc3339")
	}
}

func TestVerifyWebhookTimestamp_InvalidFormat(t *testing.T) {
	if err := VerifyWebhookTimestamp("not-a-time", 0); err != nil {
		t.Fatalf("invalid format should be ignored (nil), got %v", err)
	}
}

// small helper avoids extra imports
func strconv64(v int64) string { return string(bytes.Trim([]byte(fmtInt(v)), "\x00")) }
func fmtInt(v int64) []byte    { return []byte(fmt.Sprintf("%d", v)) }
