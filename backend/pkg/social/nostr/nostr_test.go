package nostr

import (
	"context"
	"errors"
	"strings"
	"testing"

	stdnostr "github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"

	"go-wind-cms/pkg/social"
)

// 测试私钥（仅用于单测的随机 hex）
const testSK = "6b911fd37cdf5c81d4c0adb1ab7fa822ed253ab0ad9aa18d77257c88b29b718e"

func creds() social.Credentials {
	return social.Credentials{
		"private_key": testSK,
		"relays":      "wss://relay-a.example, wss://relay-b.example",
	}
}

func TestValidateCredentials_HexKey(t *testing.T) {
	if err := NewClient().ValidateCredentials(context.Background(), creds()); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_NsecKey(t *testing.T) {
	nsec, err := nip19.EncodePrivateKey(testSK)
	if err != nil {
		t.Fatalf("encode nsec: %v", err)
	}
	bad := social.Credentials{"private_key": nsec, "relays": "wss://x"}
	if err := NewClient().ValidateCredentials(context.Background(), bad); err != nil {
		t.Fatalf("expected nil error for nsec key, got %v", err)
	}
}

func TestValidateCredentials_BadKey(t *testing.T) {
	bad := social.Credentials{"private_key": "not-a-key", "relays": "wss://x"}
	if err := NewClient().ValidateCredentials(context.Background(), bad); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if err := NewClient().ValidateCredentials(context.Background(), social.Credentials{}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatal("expected ErrInvalidCredentials for empty key")
	}
}

func TestPublish_SuccessWhenAnyRelayAccepts(t *testing.T) {
	var signed stdnostr.Event
	c := NewClient(WithPublisher(func(ctx context.Context, urls []string, evt stdnostr.Event) []publishOutcome {
		signed = evt
		return []publishOutcome{
			{URL: urls[0], Err: errors.New("blocked")},
			{URL: urls[1], Err: nil}, // 一个成功即成功
		}
	}))

	res, err := c.Publish(context.Background(), creds(), &social.PublishInput{Text: "hello nostr"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID == "" {
		t.Fatal("expected event id as external id")
	}
	if signed.ID == "" || signed.Sig == "" || signed.Kind != 1 || signed.Content != "hello nostr" {
		t.Errorf("event not properly signed: %+v", signed)
	}
}

func TestPublish_AllRelaysReject(t *testing.T) {
	c := NewClient(WithPublisher(func(ctx context.Context, urls []string, evt stdnostr.Event) []publishOutcome {
		var out []publishOutcome
		for _, u := range urls {
			out = append(out, publishOutcome{URL: u, Err: errors.New("rejected")})
		}
		return out
	}))

	_, err := c.Publish(context.Background(), creds(), &social.PublishInput{Text: "hello"})
	if err == nil || !strings.Contains(err.Error(), "all relays rejected") {
		t.Fatalf("expected all-relays-rejected error, got %v", err)
	}
}

func TestPublish_NoRelayConfigured(t *testing.T) {
	bad := social.Credentials{"private_key": testSK}
	if _, err := NewClient().Publish(context.Background(), bad, &social.PublishInput{Text: "x"}); err == nil ||
		!strings.Contains(err.Error(), "no relay") {
		t.Fatalf("expected no-relay error, got %v", err)
	}
}

func TestParseRelays(t *testing.T) {
	got := parseRelays("wss://a.example, nosuch.example\nwss://b.example;;ws://c.example")
	want := []string{"wss://a.example", "wss://nosuch.example", "wss://b.example", "ws://c.example"}
	if len(got) != len(want) {
		t.Fatalf("expected %d relays, got %v", len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("relay[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestFetchStats_Unsupported(t *testing.T) {
	if _, err := NewClient().FetchStats(context.Background(), creds(), "eventid", ""); !errors.Is(err, social.ErrStatsUnsupported) {
		t.Fatalf("expected ErrStatsUnsupported, got %v", err)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("nostr not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
