package wechat

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

func creds(server string) social.Credentials {
	return social.Credentials{
		"server": server, "app_id": "wx123", "app_secret": "sec",
		"thumb_media_id": "THUMB_MEDIA",
	}
}

// newScriptedServer 按调用顺序返回：token → draft/add → freepublish/submit →
// freepublish/get × N。publishStatus 控制轮询看到的结果态。
func newScriptedServer(t *testing.T, publishStatus int) (*Client, *[]string) {
	var calls []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.URL.Path)
		switch r.URL.Path {
		case "/cgi-bin/token":
			if r.URL.Query().Get("appid") != "wx123" || r.URL.Query().Get("secret") != "sec" {
				_ = json.NewEncoder(w).Encode(tokenResponse{wxEnvelope: wxEnvelope{ErrCode: 40125, ErrMsg: "invalid appsecret"}})
				return
			}
			_ = json.NewEncoder(w).Encode(tokenResponse{AccessToken: "AT"})
		case "/cgi-bin/draft/add":
			var req draftAddRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if len(req.Articles) != 1 || req.Articles[0].ThumbMediaID != "THUMB_MEDIA" {
				_ = json.NewEncoder(w).Encode(wxEnvelope{ErrCode: 40007, ErrMsg: "invalid media_id"})
				return
			}
			_ = json.NewEncoder(w).Encode(draftAddResponse{MediaID: "DRAFT1"})
		case "/cgi-bin/freepublish/submit":
			_ = json.NewEncoder(w).Encode(freepublishSubmitResponse{PublishID: "PUB1"})
		case "/cgi-bin/freepublish/get":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"publish_status": publishStatus,
				"article_id": map[string]any{
					"item": []map[string]any{{"idx": 1, "article_url": "https://mp.weixin.qq.com/s/abc"}},
				},
			})
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}))
	t.Cleanup(srv.Close)
	c := NewClient(WithEndpoint(srv.URL), WithPollInterval(1_000_000)) // 1ms
	return c, &calls
}

func TestValidateCredentials_OK(t *testing.T) {
	c, _ := newScriptedServer(t, 0)
	if err := c.ValidateCredentials(context.Background(), creds("x")); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateCredentials_BadSecret(t *testing.T) {
	c, _ := newScriptedServer(t, 0)
	bad := creds("x")
	bad["app_secret"] = "wrong"
	err := c.ValidateCredentials(context.Background(), bad)
	if !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestPublish_FullFlow(t *testing.T) {
	c, calls := newScriptedServer(t, 0)

	res, err := c.Publish(context.Background(), creds("x"),
		&social.PublishInput{Title: "标题", Text: "正文内容"})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res.ExternalID != "PUB1" {
		t.Errorf("unexpected external id: %s", res.ExternalID)
	}
	if res.ExternalURI != "https://mp.weixin.qq.com/s/abc" {
		t.Errorf("unexpected external uri: %s", res.ExternalURI)
	}

	// 调用顺序：token → draft/add → submit → get
	want := []string{"/cgi-bin/token", "/cgi-bin/draft/add", "/cgi-bin/freepublish/submit", "/cgi-bin/freepublish/get"}
	if len(*calls) < len(want) {
		t.Fatalf("expected at least %d calls, got %v", len(want), *calls)
	}
	for i, p := range want {
		if (*calls)[i] != p {
			t.Errorf("call[%d] = %s, want %s", i, (*calls)[i], p)
		}
	}
}

func TestPublish_MissingThumbMediaID(t *testing.T) {
	c, _ := newScriptedServer(t, 0)
	bad := creds("x")
	delete(bad, "thumb_media_id")
	_, err := c.Publish(context.Background(), bad, &social.PublishInput{Title: "t", Text: "b"})
	if err == nil || !strings.Contains(err.Error(), "thumb_media_id") {
		t.Fatalf("expected thumb_media_id guidance error, got %v", err)
	}
}

func TestValidateCredentials_MissingFields(t *testing.T) {
	if err := NewClient().ValidateCredentials(context.Background(), social.Credentials{"app_id": "x"}); !errors.Is(err, social.ErrInvalidCredentials) {
		t.Fatal("expected ErrInvalidCredentials for missing app_secret")
	}
}

func TestFetchStats_Unsupported(t *testing.T) {
	if _, err := NewClient().FetchStats(context.Background(), creds("x"), "PUB1", ""); !errors.Is(err, social.ErrStatsUnsupported) {
		t.Fatalf("expected ErrStatsUnsupported, got %v", err)
	}
}

func TestRegistryRegistration(t *testing.T) {
	p, err := social.GetProvider(Platform)
	if err != nil {
		t.Fatalf("wechat not registered via init(): %v", err)
	}
	if _, ok := p.(*Client); !ok {
		t.Fatalf("registered provider is not *Client")
	}
}
