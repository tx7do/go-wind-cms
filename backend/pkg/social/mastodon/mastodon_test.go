package mastodon

import (
	"context"
	"encoding/json"
	"io"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-wind-cms/pkg/social"
)

func creds(server string) social.Credentials {
	return social.Credentials{"server": server, "access_token": "token-abc"}
}

func TestValidateCredentials_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/accounts/verify_credentials" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token-abc" {
			t.Errorf("missing bearer, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(accountView{ID: "1", Username: "alice"})
	}))
	defer srv.Close()

	c := NewClient()
	if err := c.ValidateCredentials(context.Background(), creds(srv.URL)); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_BadToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "The access token is invalid"})
	}))
	defer srv.Close()

	err := NewClient().ValidateCredentials(context.Background(), creds(srv.URL))
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateCredentials_MissingToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	}))
	defer srv.Close()

	err := NewClient().ValidateCredentials(context.Background(), social.Credentials{"server": srv.URL})
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for missing token, got %v", err)
	}
}

func TestValidateCredentials_MissingServer(t *testing.T) {
	err := NewClient().ValidateCredentials(context.Background(), social.Credentials{"access_token": "x"})
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for missing server, got %v", err)
	}
}

func TestPublish_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses" || r.Method != http.MethodPost {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		body := string(raw)
		if !strings.Contains(body, "status=hello") || !strings.Contains(body, "visibility=public") {
			t.Errorf("unexpected form body: %s", body)
		}
		_ = json.NewEncoder(w).Encode(statusView{ID: "109", URL: "https://mastodon.social/@alice/109"})
	}))
	defer srv.Close()

	res, err := NewClient().Publish(context.Background(), creds(srv.URL), &social.PublishInput{Text: "hello"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "109" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
}

func TestPublish_TooLong(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	}))
	defer srv.Close()

	long := strings.Repeat("好", 501)
	_, err := NewClient().Publish(context.Background(), creds(srv.URL), &social.PublishInput{Text: long})
	if !errors.Is(err, social.ErrContentTooLong) {
		t.Fatalf("expected ErrContentTooLong, got %v", err)
	}
}

func TestPublish_RateLimited(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Rate limited"})
	}))
	defer srv.Close()

	_, err := NewClient().Publish(context.Background(), creds(srv.URL), &social.PublishInput{Text: "hi"})
	if !errors.Is(err, social.ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestFetchStats_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/statuses/109" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(statusView{
			ID: "109", FavouritesCount: 7, RepliesCount: 2, ReblogsCount: 4,
		})
	}))
	defer srv.Close()

	stats, err := NewClient().FetchStats(context.Background(), creds(srv.URL), "109", "")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if stats.Likes != 7 || stats.Replies != 2 || stats.Reposts != 4 {
		t.Errorf("unexpected stats: %+v", stats)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("mastodon not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
