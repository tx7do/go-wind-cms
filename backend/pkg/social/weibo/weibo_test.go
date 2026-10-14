package weibo

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
	return social.Credentials{"server": server, "access_token": "TOKEN"}
}

func TestValidateCredentials_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/account/rate_limit_status.json" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("access_token") != "TOKEN" {
			t.Errorf("missing access_token query")
		}
		_ = json.NewEncoder(w).Encode(rateLimitView{UserID: 123})
	}))
	defer srv.Close()

	if err := NewClient(WithEndpoint(srv.URL)).ValidateCredentials(context.Background(), creds(srv.URL)); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_Expired(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "applications over quota", "error_code": 21327})
	}))
	defer srv.Close()

	err := NewClient().ValidateCredentials(context.Background(), creds(srv.URL))
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateCredentials_MissingToken(t *testing.T) {
	if err := NewClient().ValidateCredentials(context.Background(), social.Credentials{}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatal("expected ErrInvalidCredentials for missing token")
	}
}

func TestPublish_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/statuses/update.json" || r.Method != http.MethodPost {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(raw), "status=hello") {
			t.Errorf("unexpected form body: %s", string(raw))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 5001, "user": map[string]any{"id": 123}})
	}))
	defer srv.Close()

	res, err := NewClient(WithEndpoint(srv.URL)).Publish(context.Background(), creds(srv.URL), &social.PublishInput{Text: "hello"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "5001" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
	if res.ExternalURI != "" {
		t.Errorf("weibo url should be empty, got %s", res.ExternalURI)
	}
}

func TestPublish_TooLong(t *testing.T) {
	long := strings.Repeat("博", maxStatusRunes+1)
	_, err := NewClient().Publish(context.Background(), creds("x"), &social.PublishInput{Text: long})
	if !errors.Is(err, social.ErrContentTooLong) {
		t.Fatalf("expected ErrContentTooLong, got %v", err)
	}
}

func TestFetchStats_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/statuses/show.json" || r.URL.Query().Get("id") != "5001" {
			t.Errorf("unexpected %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(statusView{ID: 5001, AttitudesCount: 6, CommentsCount: 2, RepostsCount: 1})
	}))
	defer srv.Close()

	stats, err := NewClient(WithEndpoint(srv.URL)).FetchStats(context.Background(), creds(srv.URL), "5001", "")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if stats.Likes != 6 || stats.Replies != 2 || stats.Reposts != 1 {
		t.Errorf("unexpected stats: %+v", stats)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("weibo not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
