package wordpress

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-wind-cms/pkg/social"
)

func creds(server string) social.Credentials {
	return social.Credentials{"server": server, "username": "alice", "app_password": "abcd efgh"}
}

func TestValidateCredentials_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wp-json/wp/v2/users/me" || !strings.Contains(r.URL.RawQuery, "context=edit") {
			t.Errorf("unexpected %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		want := "Basic " + base64.StdEncoding.EncodeToString([]byte("alice:abcd efgh"))
		if got := r.Header.Get("Authorization"); got != want {
			t.Errorf("unexpected auth header %q", got)
		}
		_ = json.NewEncoder(w).Encode(userView{ID: 1, Name: "Alice"})
	}))
	defer srv.Close()

	if err := NewClient().ValidateCredentials(context.Background(), creds(srv.URL)); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "incorrect_password", "message": "密码不正确"})
	}))
	defer srv.Close()

	err := NewClient().ValidateCredentials(context.Background(), creds(srv.URL))
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateCredentials_MissingFields(t *testing.T) {
	if err := NewClient().ValidateCredentials(context.Background(), social.Credentials{"server": "https://x.com"}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatal("expected ErrInvalidCredentials for missing username/app_password")
	}
	if err := NewClient().ValidateCredentials(context.Background(), social.Credentials{"username": "a"}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatal("expected ErrInvalidCredentials for missing server")
	}
}

func TestPublish_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wp-json/wp/v2/posts" || r.Method != http.MethodPost {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var req createPostRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Title != "My Title" || req.Status != "publish" {
			t.Errorf("unexpected post payload: %+v", req)
		}
		if !strings.Contains(req.Content, "<p>body text</p>") {
			t.Errorf("expected html content, got %q", req.Content)
		}
		_ = json.NewEncoder(w).Encode(postView{ID: 7, Link: "https://example.com/my-title/"})
	}))
	defer srv.Close()

	res, err := NewClient().Publish(context.Background(), creds(srv.URL),
		&social.PublishInput{Title: "My Title", Text: "body text"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "7" || res.ExternalURI != "https://example.com/my-title/" {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestPublish_EmptyContent(t *testing.T) {
	if _, err := NewClient().Publish(context.Background(), creds("https://x.com"),
		&social.PublishInput{Text: "  "}); err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestFetchStats_Unsupported(t *testing.T) {
	if _, err := NewClient().FetchStats(context.Background(), creds("https://x.com"), "7", ""); !errors.Is(err, social.ErrStatsUnsupported) {
		t.Fatalf("expected ErrStatsUnsupported, got %v", err)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("wordpress not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
