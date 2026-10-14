package devto

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
	return social.Credentials{"server": server, "api_key": "key-abc"}
}

func TestValidateCredentials_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users/me" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if got := r.Header.Get("api-key"); got != "key-abc" {
			t.Errorf("missing api-key, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(userView{ID: 1, Username: "alice"})
	}))
	defer srv.Close()

	if err := NewClient().ValidateCredentials(context.Background(), creds(srv.URL)); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_BadKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
	}))
	defer srv.Close()

	err := NewClient().ValidateCredentials(context.Background(), creds(srv.URL))
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateCredentials_MissingKey(t *testing.T) {
	err := NewClient().ValidateCredentials(context.Background(), social.Credentials{})
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for missing api_key, got %v", err)
	}
}

func TestPublish_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/articles" || r.Method != http.MethodPost {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		var req createArticleRequest
		_ = json.Unmarshal(raw, &req)
		if req.Article.Title != "My Title" || req.Article.BodyMarkdown != "body text" || !req.Article.Published {
			t.Errorf("unexpected article payload: %+v", req.Article)
		}
		_ = json.NewEncoder(w).Encode(articleView{ID: 42, URL: "https://dev.to/alice/my-title"})
	}))
	defer srv.Close()

	res, err := NewClient().Publish(context.Background(), creds(srv.URL),
		&social.PublishInput{Title: "My Title", Text: "body text"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "42" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
}

func TestPublish_DeriveTitleFromFirstLine(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var req createArticleRequest
		_ = json.Unmarshal(raw, &req)
		if req.Article.Title != "Derived Title" {
			t.Errorf("expected derived title, got %q", req.Article.Title)
		}
		_ = json.NewEncoder(w).Encode(articleView{ID: 43})
	}))
	defer srv.Close()

	_, err := NewClient().Publish(context.Background(), creds(srv.URL),
		&social.PublishInput{Text: "## Derived Title\n\nbody"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestPublish_EmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	}))
	defer srv.Close()

	if _, err := NewClient().Publish(context.Background(), creds(srv.URL),
		&social.PublishInput{Text: "  "}); err == nil {
		t.Fatal("expected error for empty body")
	}
}

func TestFetchStats_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/articles/42" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(articleView{ID: 42, PositiveReactionsCount: 9, CommentsCount: 3})
	}))
	defer srv.Close()

	stats, err := NewClient().FetchStats(context.Background(), creds(srv.URL), "42", "")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if stats.Likes != 9 || stats.Replies != 3 {
		t.Errorf("unexpected stats: %+v", stats)
	}
}

func TestDeriveTitle_LongTruncated(t *testing.T) {
	long := strings.Repeat("题", maxTitleRunes+10)
	got := deriveTitle("", long)
	if runeCount := len([]rune(got)); runeCount != maxTitleRunes {
		t.Fatalf("expected title truncated to %d runes, got %d", maxTitleRunes, runeCount)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("devto not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
