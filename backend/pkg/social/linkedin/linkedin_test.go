package linkedin

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
		if r.URL.Path != "/v2/userinfo" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer TOKEN" {
			t.Errorf("missing bearer")
		}
		_ = json.NewEncoder(w).Encode(userInfoView{Sub: "abc123", Name: "Alice"})
	}))
	defer srv.Close()

	if err := NewClient(WithEndpoint(srv.URL)).ValidateCredentials(context.Background(), creds(srv.URL)); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_Expired(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "expired access token"})
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
		switch r.URL.Path {
		case "/v2/userinfo":
			_ = json.NewEncoder(w).Encode(userInfoView{Sub: "abc123"})
		case "/rest/posts":
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.Header.Get("X-Restli-Protocol") != "2.0.0" || r.Header.Get("LinkedIn-Version") == "" {
				t.Errorf("missing required headers")
			}
			raw, _ := io.ReadAll(r.Body)
			var req createPostRequest
			_ = json.Unmarshal(raw, &req)
			if req.Author != "urn:li:person:abc123" || req.Visibility != "PUBLIC" || req.Commentary != "hello" {
				t.Errorf("unexpected post payload: %+v", req)
			}
			w.Header().Set("X-RestLi-Id", "7789")
			w.WriteHeader(http.StatusCreated)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	res, err := NewClient(WithEndpoint(srv.URL)).Publish(context.Background(), creds(srv.URL), &social.PublishInput{Text: "hello"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "7789" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
	if res.ExternalURI != "https://www.linkedin.com/feed/update/urn:li:activity:7789" {
		t.Errorf("unexpected external uri: %s", res.ExternalURI)
	}
}

func TestPublish_TooLong(t *testing.T) {
	long := strings.Repeat("字", maxCommentaryRunes+1)
	_, err := NewClient().Publish(context.Background(), creds("x"), &social.PublishInput{Text: long})
	if !errors.Is(err, social.ErrContentTooLong) {
		t.Fatalf("expected ErrContentTooLong, got %v", err)
	}
}

func TestFetchStats_Unsupported(t *testing.T) {
	if _, err := NewClient().FetchStats(context.Background(), creds("x"), "7789", ""); !errors.Is(err, social.ErrStatsUnsupported) {
		t.Fatalf("expected ErrStatsUnsupported, got %v", err)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("linkedin not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
