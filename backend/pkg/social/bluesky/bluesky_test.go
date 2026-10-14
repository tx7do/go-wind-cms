package bluesky

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-wind-cms/pkg/social"
)

// newTestServer 起一个假的 PDS：按路径分发 XRPC 响应。
func newTestServer(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewClient(WithEndpoint(srv.URL))
}

func creds() social.Credentials {
	return social.Credentials{"handle": "alice.test", "app_password": "abcd-efgh-ijkl"}
}

// sessionHandler 返回统一的 createSession 成功响应。
func sessionHandler(w http.ResponseWriter) {
	_ = json.NewEncoder(w).Encode(createSessionResponse{
		Did: "did:plc:test1234", Handle: "alice.test", AccessJwt: "jwt-token",
	})
}

func TestValidateCredentials_OK(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/xrpc/com.atproto.server.createSession" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var req createSessionRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Identifier != "alice.test" || req.Password != "abcd-efgh-ijkl" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(xrpcError{Error: "AuthenticationNeeded", Message: "invalid"})
			return
		}
		sessionHandler(w)
	})

	if err := c.ValidateCredentials(context.Background(), creds()); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_BadPassword(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(xrpcError{Error: "AuthenticationNeeded", Message: "Invalid identifier or password"})
	})

	err := c.ValidateCredentials(context.Background(), creds())
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestValidateCredentials_MissingFields(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	})
	if err := c.ValidateCredentials(context.Background(), social.Credentials{"handle": "x"}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for missing app_password, got %v", err)
	}
}

func TestPublish_OK(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/xrpc/com.atproto.server.createSession":
			sessionHandler(w)
		case "/xrpc/com.atproto.repo.createRecord":
			if got := r.Header.Get("Authorization"); got != "Bearer jwt-token" {
				t.Errorf("missing bearer, got %q", got)
			}
			var req createRecordRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.Repo != "did:plc:test1234" || req.Collection != postCollection {
				t.Errorf("unexpected createRecord body: %+v", req)
			}
			if req.Record.Text != "hello world" {
				t.Errorf("unexpected text: %q", req.Record.Text)
			}
			_ = json.NewEncoder(w).Encode(createRecordResponse{
				Uri: "at://did:plc:test1234/app.bsky.feed.post/3kxyz",
				Cid: "bafytest",
			})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})

	res, err := c.Publish(context.Background(), creds(), &social.PublishInput{Text: "hello world"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "at://did:plc:test1234/app.bsky.feed.post/3kxyz" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
	if res.ExternalURI != "https://bsky.app/profile/alice.test/post/3kxyz" {
		t.Errorf("unexpected external uri: %s", res.ExternalURI)
	}
}

func TestPublish_TooLong(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	})
	long := strings.Repeat("好", 301)
	_, err := c.Publish(context.Background(), creds(), &social.PublishInput{Text: long})
	if !errors.Is(err, social.ErrContentTooLong) {
		t.Fatalf("expected ErrContentTooLong, got %v", err)
	}
}

func TestPublish_EmptyText(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not reach server")
	})
	if _, err := c.Publish(context.Background(), creds(), &social.PublishInput{Text: "  "}); err == nil {
		t.Fatal("expected error for empty text")
	}
}

func TestPublish_RateLimited(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/xrpc/com.atproto.server.createSession" {
			sessionHandler(w)
			return
		}
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(xrpcError{Error: "RateLimitExceeded", Message: "slow down"})
	})

	_, err := c.Publish(context.Background(), creds(), &social.PublishInput{Text: "hi"})
	if !errors.Is(err, social.ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestFetchStats_OK(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/xrpc/com.atproto.server.createSession":
			sessionHandler(w)
		case "/xrpc/app.bsky.feed.getPosts":
			uri := r.URL.Query().Get("uris")
			if uri != "at://did:plc:test1234/app.bsky.feed.post/3kxyz" {
				t.Errorf("unexpected uris param: %q", uri)
			}
			_ = json.NewEncoder(w).Encode(getPostsResponse{Posts: []postView{{
				Uri: uri, LikeCount: 12, ReplyCount: 3, RepostCount: 5,
			}}})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})

	stats, err := c.FetchStats(context.Background(), creds(),
		"at://did:plc:test1234/app.bsky.feed.post/3kxyz", "")
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if stats.Likes != 12 || stats.Replies != 3 || stats.Reposts != 5 {
		t.Errorf("unexpected stats: %+v", stats)
	}
}

func TestFetchStats_NotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/xrpc/com.atproto.server.createSession" {
			sessionHandler(w)
			return
		}
		_ = json.NewEncoder(w).Encode(getPostsResponse{Posts: nil})
	})

	if _, err := c.FetchStats(context.Background(), creds(), "at://gone", ""); !errors.Is(err, social.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestWebURL(t *testing.T) {
	got := webURL("alice.test", "at://did:plc:x/app.bsky.feed.post/3k")
	if got != "https://bsky.app/profile/alice.test/post/3k" {
		t.Errorf("unexpected web url: %s", got)
	}
	// 非 at:// URI 原样返回
	if got := webURL("", "https://example.com"); got != "https://example.com" {
		t.Errorf("expected passthrough, got %s", got)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("bluesky not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
