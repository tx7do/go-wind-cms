package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-wind-cms/pkg/social"
)

func creds(server string) social.Credentials {
	return social.Credentials{"bot_token": "123:ABC", "chat_id": "@mychannel", "server": server}
}

func TestValidateCredentials_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/bot123:ABC/getMe") {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(apiResponse{OK: true, Result: json.RawMessage(`{"id":1,"username":"bot"}`)})
	}))
	defer srv.Close()

	c := NewClient(WithEndpoint(srv.URL))
	if err := c.ValidateCredentials(context.Background(), creds(srv.URL)); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_BadToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(apiResponse{OK: false, ErrorCode: 401, Description: "Unauthorized"})
	}))
	defer srv.Close()

	err := NewClient(WithEndpoint(srv.URL)).ValidateCredentials(context.Background(), creds(srv.URL))
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateCredentials_MissingToken(t *testing.T) {
	err := NewClient().ValidateCredentials(context.Background(), social.Credentials{})
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestPublish_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/sendMessage") || r.Method != http.MethodPost {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		body := string(raw)
		if !strings.Contains(body, "chat_id=%40mychannel") || !strings.Contains(body, "text=hello") {
			t.Errorf("unexpected form body: %s", body)
		}
		_ = json.NewEncoder(w).Encode(apiResponse{OK: true, Result: json.RawMessage(`{"message_id":42,"chat":{"id":-100,"username":"mychannel"}}`)})
	}))
	defer srv.Close()

	res, err := NewClient(WithEndpoint(srv.URL)).Publish(context.Background(), creds(srv.URL),
		&social.PublishInput{Text: "hello"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "42" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
	if res.ExternalURI != "https://t.me/mychannel/42" {
		t.Errorf("unexpected external uri: %s", res.ExternalURI)
	}
}

func TestPublish_TooLong(t *testing.T) {
	long := strings.Repeat("字", maxMessageRunes+1)
	_, err := NewClient().Publish(context.Background(), creds("x"),
		&social.PublishInput{Text: long})
	if !errors.Is(err, social.ErrContentTooLong) {
		t.Fatalf("expected ErrContentTooLong, got %v", err)
	}
}

func TestFetchStats_Unsupported(t *testing.T) {
	_, err := NewClient().FetchStats(context.Background(), creds("x"), "42", "")
	if !errors.Is(err, social.ErrStatsUnsupported) {
		t.Fatalf("expected ErrStatsUnsupported, got %v", err)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("telegram not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
