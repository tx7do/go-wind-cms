package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	"go-wind-cms/app/core/service/internal/data"
)

// editLockTTL 锁 TTL：心跳间隔（前端 60s）的两倍余量，进程死亡后自动释放。
const editLockTTL = 120 * time.Second

// validEditLockTargetTypes 允许加锁的目标类型（可扩展）。
var validEditLockTargetTypes = map[string]bool{"post": true, "page": true}

// luaEditLockAcquire 获取/续期锁：本人持锁则刷新 TTL（心跳语义），
// 他人持锁原样返回现值供调用方判归属；空闲则上锁。
var luaEditLockAcquire = redis.NewScript(`
local cur = redis.call('GET', KEYS[1])
if cur then
  local ok, d = pcall(cjson.decode, cur)
  if ok and d ~= nil and tostring(d.user_id) == ARGV[1] then
    redis.call('SET', KEYS[1], ARGV[2], 'EX', tonumber(ARGV[3]))
    return ARGV[2]
  end
  return cur
end
redis.call('SET', KEYS[1], ARGV[2], 'EX', tonumber(ARGV[3]))
return ARGV[2]
`)

// luaEditLockRelease 释放锁：仅持锁者可删；未持锁视为无事发生。
var luaEditLockRelease = redis.NewScript(`
local cur = redis.call('GET', KEYS[1])
if not cur then
  return 1
end
local ok, d = pcall(cjson.decode, cur)
if ok and d ~= nil and tostring(d.user_id) == ARGV[1] then
  redis.call('DEL', KEYS[1])
end
return 1
`)

// editLockValue Redis 锁值。
type editLockValue struct {
	UserID     uint32 `json:"user_id"`
	UserName   string `json:"user_name"`
	AcquiredAt string `json:"acquired_at"`
}

// EditLockService 编辑锁（协同编辑 v1）：Redis 短 TTL 软锁，只做编辑感知
// 与冲突提示，不做强制写拦截（保存冲突由内容版本化兜底）。
type EditLockService struct {
	contentV1.UnimplementedEditLockServiceServer

	rdb      *redis.Client
	userRepo data.UserRepo
	log      *log.Helper
}

func NewEditLockService(ctx *bootstrap.Context, rdb *redis.Client, userRepo data.UserRepo) *EditLockService {
	return &EditLockService{
		log:      ctx.NewLoggerHelper("edit-lock/service/core-service"),
		rdb:      rdb,
		userRepo: userRepo,
	}
}

// Acquire 获取/续期编辑锁（重复调用即心跳）。
func (s *EditLockService) Acquire(ctx context.Context, req *contentV1.AcquireEditLockRequest) (*contentV1.EditLock, error) {
	if req == nil || !validEditLockTargetTypes[req.GetTargetType()] || req.GetTargetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	uid, hasUser := viewerUserIDFromContext(ctx)
	if !hasUser {
		return nil, contentV1.ErrorForbidden("editor lock requires an authenticated user")
	}

	userName := ""
	if user, err := s.userRepo.Get(ctx, &identityV1.GetUserRequest{
		QueryBy: &identityV1.GetUserRequest_Id{Id: uid},
	}); err != nil || user == nil {
		s.log.Warnf("edit lock: resolve user %d name failed: %v", uid, err)
	} else {
		userName = firstNonEmpty(user.GetNickname(), user.GetUsername())
	}

	value, err := json.Marshal(editLockValue{
		UserID:     uid,
		UserName:   userName,
		AcquiredAt: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return nil, contentV1.ErrorInternalServerError("marshal lock value failed")
	}

	key := editLockKey(req.GetTargetType(), req.GetTargetId())
	raw, err := luaEditLockAcquire.Run(ctx, s.rdb, []string{key},
		fmt.Sprintf("%d", uid), string(value), fmt.Sprintf("%d", int(editLockTTL.Seconds())),
	).Text()
	if err != nil {
		s.log.Errorf("edit lock acquire failed (%s): %v", key, err)
		return nil, contentV1.ErrorInternalServerError("acquire edit lock failed")
	}

	var held editLockValue
	if err := json.Unmarshal([]byte(raw), &held); err != nil {
		s.log.Errorf("edit lock value decode failed (%s): %v", key, err)
		return nil, contentV1.ErrorInternalServerError("decode edit lock failed")
	}

	at, _ := time.Parse(time.RFC3339, held.AcquiredAt)
	return &contentV1.EditLock{
		Acquired:   trans.Ptr(held.UserID == uid),
		HolderId:   &held.UserID,
		HolderName: &held.UserName,
		AcquiredAt: timestamppb.New(at),
	}, nil
}

// Release 释放编辑锁（仅持锁者可释放）。
func (s *EditLockService) Release(ctx context.Context, req *contentV1.ReleaseEditLockRequest) (*emptypb.Empty, error) {
	if req == nil || !validEditLockTargetTypes[req.GetTargetType()] || req.GetTargetId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}
	uid, hasUser := viewerUserIDFromContext(ctx)
	if !hasUser {
		return nil, contentV1.ErrorForbidden("editor lock requires an authenticated user")
	}
	if _, err := luaEditLockRelease.Run(ctx, s.rdb,
		[]string{editLockKey(req.GetTargetType(), req.GetTargetId())},
		fmt.Sprintf("%d", uid),
	).Result(); err != nil {
		s.log.Errorf("edit lock release failed (%s/%d): %v", req.GetTargetType(), req.GetTargetId(), err)
		return nil, contentV1.ErrorInternalServerError("release edit lock failed")
	}
	return &emptypb.Empty{}, nil
}

func editLockKey(targetType string, targetID uint32) string {
	return fmt.Sprintf("gowind-cms:editlock:%s:%d", targetType, targetID)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
