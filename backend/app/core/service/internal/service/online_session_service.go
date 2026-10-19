package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	onlineSessionV1 "go-wind-cms/api/gen/go/online_session/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// OnlineSessionService 在线会话视图（基于 Redis 令牌对元数据）与强制下线能力。
// 会话记录在令牌签发时写入、随 refresh token 过期或吊销时删除，
// 本服务只做读取与转发吊销，不持有独立存储。
type OnlineSessionService struct {
	onlineSessionV1.UnimplementedOnlineSessionServiceServer

	log           *log.Helper
	authenticator *data.Authenticator
}

// NewOnlineSessionService 构造（authenticator 持有会话元数据的读写与吊销）。
func NewOnlineSessionService(ctx *bootstrap.Context, authenticator *data.Authenticator) *OnlineSessionService {
	return &OnlineSessionService{
		log:           ctx.NewLoggerHelper("online-session/service/core-service"),
		authenticator: authenticator,
	}
}

// entriesToProto 扫描条目 → proto 列表（关键词过滤 + 登录时间倒序 + 标记当前会话）。
func entriesToProto(entries []data.SessionEntry, keyword, currentJti string) []*onlineSessionV1.OnlineSession {
	kw := strings.ToLower(strings.TrimSpace(keyword))

	items := make([]*onlineSessionV1.OnlineSession, 0, len(entries))
	for _, entry := range entries {
		meta := entry.Meta
		if meta == nil {
			continue
		}
		if kw != "" &&
			!strings.Contains(strings.ToLower(meta.Username), kw) &&
			!strings.Contains(strings.ToLower(meta.Ip), kw) {
			continue
		}

		loginAt := time.Unix(meta.LoginAt, 0)
		if meta.LoginAt <= 0 {
			loginAt = time.Now()
		}

		items = append(items, &onlineSessionV1.OnlineSession{
			Jti:        trans.Ptr(entry.Jti),
			UserId:     trans.Ptr(entry.UserId),
			Username:   trans.Ptr(meta.Username),
			TenantId:   trans.Ptr(meta.TenantId),
			ClientType: trans.Ptr(entry.ClientType),
			IpAddress:  trans.Ptr(meta.Ip),
			UserAgent:  trans.Ptr(meta.UserAgent),
			LoginAt:    timestamppb.New(loginAt),
			Current:    trans.Ptr(currentJti != "" && currentJti == entry.Jti),
		})
	}

	// 登录时间倒序：最新登录的会话排前面
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetLoginAt().AsTime().After(items[j].GetLoginAt().AsTime())
	})
	return items
}

// ListOnlineSession 全量在线会话列表（平台管理员）。
func (s *OnlineSessionService) List(ctx context.Context, req *onlineSessionV1.ListOnlineSessionRequest) (*onlineSessionV1.ListOnlineSessionResponse, error) {
	if err := requireNotificationPlatformContext(ctx, "online sessions"); err != nil {
		return nil, err
	}

	entries, err := s.authenticator.ListSessionEntries(ctx)
	if err != nil {
		s.log.Errorf("list session entries failed: %s", err.Error())
		return nil, errors.InternalServer("ONLINE_SESSION", "list online sessions failed")
	}

	items := entriesToProto(entries, req.GetKeyword(), "")
	total := uint64(len(items))

	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	start := (page - 1) * pageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	return &onlineSessionV1.ListOnlineSessionResponse{
		Items: items[start:end],
		Total: total,
	}, nil
}

// ForceLogoutSession 强制下线指定会话（平台管理员）。
func (s *OnlineSessionService) ForceLogout(ctx context.Context, req *onlineSessionV1.ForceLogoutSessionRequest) (*onlineSessionV1.ForceLogoutSessionResponse, error) {
	if err := requireNotificationPlatformContext(ctx, "online sessions"); err != nil {
		return nil, err
	}
	if req.GetUserId() == 0 || req.GetJti() == "" {
		return nil, errors.BadRequest("ONLINE_SESSION", "userId and jti are required")
	}

	if err := s.authenticator.RevokeTokenByJti(ctx, req.ClientType, req.GetUserId(), req.GetJti()); err != nil {
		s.log.Errorf("force logout session [%d/%s] failed: %s", req.GetUserId(), req.GetJti(), err.Error())
		return nil, errors.InternalServer("ONLINE_SESSION", "force logout failed")
	}

	s.log.Infof("session [%d/%s] force-logged-out by operator", req.GetUserId(), req.GetJti())
	return &onlineSessionV1.ForceLogoutSessionResponse{}, nil
}

// ListMy 当前登录用户的在线会话列表。
func (s *OnlineSessionService) ListMy(ctx context.Context, req *onlineSessionV1.ListMyOnlineSessionRequest) (*onlineSessionV1.ListOnlineSessionResponse, error) {
	if req.GetUserId() == 0 {
		return nil, errors.BadRequest("ONLINE_SESSION", "userId is required")
	}

	entries, err := s.authenticator.ListSessionEntries(ctx)
	if err != nil {
		s.log.Errorf("list my session entries failed: %s", err.Error())
		return nil, errors.InternalServer("ONLINE_SESSION", "list online sessions failed")
	}

	mine := make([]data.SessionEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.UserId == req.GetUserId() {
			mine = append(mine, entry)
		}
	}

	return &onlineSessionV1.ListOnlineSessionResponse{
		Items: entriesToProto(mine, "", req.GetCurrentJti()),
		Total: uint64(len(mine)),
	}, nil
}

// RevokeMy 当前用户强制下线自己的指定会话。
func (s *OnlineSessionService) RevokeMy(ctx context.Context, req *onlineSessionV1.RevokeMyOnlineSessionRequest) (*onlineSessionV1.RevokeMyOnlineSessionResponse, error) {
	if req.GetUserId() == 0 || req.GetJti() == "" {
		return nil, errors.BadRequest("ONLINE_SESSION", "userId and jti are required")
	}

	if err := s.authenticator.RevokeTokenByJti(ctx, req.ClientType, req.GetUserId(), req.GetJti()); err != nil {
		s.log.Errorf("revoke my session [%d/%s] failed: %s", req.GetUserId(), req.GetJti(), err.Error())
		return nil, errors.InternalServer("ONLINE_SESSION", "revoke session failed")
	}

	return &onlineSessionV1.RevokeMyOnlineSessionResponse{}, nil
}
