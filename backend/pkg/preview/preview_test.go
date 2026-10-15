package preview

import (
	"strings"
	"testing"
	"time"
)

func TestMintVerifyRoundtrip(t *testing.T) {
	now := time.Now()
	token, exp, err := Mint(KindPost, 42, time.Minute, now)
	if err != nil {
		t.Fatalf("mint failed: %v", err)
	}
	if !exp.After(now) {
		t.Fatalf("expiry should be in the future")
	}

	id, err := Verify(token, KindPost, now.Add(30*time.Second))
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if id != 42 {
		t.Fatalf("id = %d, want 42", id)
	}
}

func TestVerifyExpired(t *testing.T) {
	now := time.Now()
	token, _, err := Mint(KindPage, 7, time.Minute, now)
	if err != nil {
		t.Fatalf("mint failed: %v", err)
	}
	if _, err = Verify(token, KindPage, now.Add(2*time.Minute)); err == nil {
		t.Fatalf("expired token should fail verification")
	}
}

func TestVerifyKindMismatch(t *testing.T) {
	now := time.Now()
	token, _, err := Mint(KindPost, 7, time.Minute, now)
	if err != nil {
		t.Fatalf("mint failed: %v", err)
	}
	if _, err = Verify(token, KindPage, now); err == nil {
		t.Fatalf("kind mismatch should fail verification")
	}
}

func TestVerifyTampered(t *testing.T) {
	now := time.Now()
	token, _, err := Mint(KindPost, 7, time.Minute, now)
	if err != nil {
		t.Fatalf("mint failed: %v", err)
	}
	tampered := token[:len(token)-2] + "xx"
	if _, err = Verify(tampered, KindPost, now); err == nil {
		t.Fatalf("tampered token should fail verification")
	}
	if _, err = Verify(strings.ToUpper(token), KindPost, now); err == nil {
		t.Fatalf("re-encoded token should fail verification")
	}
}

func TestMintZeroID(t *testing.T) {
	if _, _, err := Mint(KindPost, 0, time.Minute, time.Now()); err == nil {
		t.Fatalf("zero id should be rejected")
	}
}
