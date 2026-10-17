package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/emptypb"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"

	"go-wind-cms/app/core/service/internal/data"
)

// 静音时段默认窗口：22:00 → 次日 08:00（分钟数自当日 00:00 起）。
// 与 ent schema 的列默认值保持一致；用户从未配置过偏好时按此返回。
const (
	defaultQuietStartMinute = 1320
	defaultQuietEndMinute   = 480

	quietMinuteUpperBound = 24*60 - 1 // 1439
)

// NotificationPreferenceService 用户通知偏好（自助服务）：
// 读写的都是当前操作人自己的偏好，无管理面、无平台守卫。
type NotificationPreferenceService struct {
	notifyV1.UnimplementedNotificationPreferenceServiceServer

	prefRepo     *data.NotificationPreferenceRepo
	categoryRepo *data.InternalMessageCategoryRepo
	log          *log.Helper
}

func NewNotificationPreferenceService(
	ctx *bootstrap.Context,
	prefRepo *data.NotificationPreferenceRepo,
	categoryRepo *data.InternalMessageCategoryRepo,
) *NotificationPreferenceService {
	return &NotificationPreferenceService{
		prefRepo:     prefRepo,
		categoryRepo: categoryRepo,
		log:          ctx.NewLoggerHelper("notification-preference/service/core-service"),
	}
}

// GetMyNotificationPreference 返回当前用户的偏好；从未配置过时返回默认值
//（未启用静音、无退订），前端据此渲染表单，落库发生在第一次保存。
func (s *NotificationPreferenceService) GetMyNotificationPreference(ctx context.Context, _ *emptypb.Empty) (*notifyV1.NotificationPreference, error) {
	userID, hasUser := viewerUserIDFromContext(ctx)
	if !hasUser {
		return nil, errors.BadRequest("NOTIFICATION_PREFERENCE", "operator identity is required")
	}

	pref, err := s.prefRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if pref == nil {
		return &notifyV1.NotificationPreference{
			UserId:           trans.Ptr(userID),
			QuietEnabled:     trans.Ptr(false),
			QuietStartMinute: trans.Ptr(int32(defaultQuietStartMinute)),
			QuietEndMinute:   trans.Ptr(int32(defaultQuietEndMinute)),
		}, nil
	}

	return pref, nil
}

// UpdateMyNotificationPreference 保存当前用户的偏好。user_id 从操作人钉定，
// 客户端传什么都不影响写入对象。
func (s *NotificationPreferenceService) UpdateMyNotificationPreference(ctx context.Context, req *notifyV1.UpdateNotificationPreferenceRequest) (*notifyV1.NotificationPreference, error) {
	userID, hasUser := viewerUserIDFromContext(ctx)
	if !hasUser {
		return nil, errors.BadRequest("NOTIFICATION_PREFERENCE", "operator identity is required")
	}

	if req == nil {
		return nil, errors.BadRequest("NOTIFICATION_PREFERENCE", "invalid parameter")
	}
	if err := validateQuietWindow(req.GetQuietEnabled(), req.GetQuietStartMinute(), req.GetQuietEndMinute()); err != nil {
		return nil, err
	}
	if len(req.GetMutedCategoryIds()) > 100 {
		return nil, errors.BadRequest("NOTIFICATION_PREFERENCE", "too many muted categories")
	}

	pref, err := s.prefRepo.Upsert(ctx, userID, req, userID)
	if err != nil {
		return nil, err
	}

	return pref, nil
}

// ListMyNotificationCategories 当前用户可退订的启用分类清单。
// 分类是租户实体，viewer 租户谓词天然把清单筛成本租户。
// v1 无全员广播，本清单暂无消费方（偏好组件不渲染退订），先随域一并落地。
func (s *NotificationPreferenceService) ListMyNotificationCategories(ctx context.Context, _ *emptypb.Empty) (*notifyV1.ListMyNotificationCategoriesResponse, error) {
	categories, err := s.categoryRepo.List(ctx, &paginationV1.PagingRequest{
		NoPaging: trans.Ptr(true),
	})
	if err != nil {
		return nil, err
	}

	items := make([]*notifyV1.NotificationCategoryItem, 0, len(categories.GetItems()))
	for _, c := range categories.GetItems() {
		if !c.GetIsEnabled() {
			continue
		}
		items = append(items, &notifyV1.NotificationCategoryItem{
			Id:   c.Id,
			Name: c.Name,
		})
	}

	return &notifyV1.ListMyNotificationCategoriesResponse{Items: items}, nil
}

// validateQuietWindow 校验静音窗口：分钟数在 0..1439 且启用时 start != end
//（start == end 的窗口宽为零，等价于不静音，按配置错误拒绝而不是默默不生效）。
// start > end 是合法的跨零点窗口（如 22:00→08:00）。
func validateQuietWindow(enabled bool, startMinute, endMinute int32) error {
	for _, m := range []int32{startMinute, endMinute} {
		if m < 0 || m > quietMinuteUpperBound {
			return errors.BadRequest("NOTIFICATION_PREFERENCE", "quiet minutes must be within 0..1439")
		}
	}
	if enabled && startMinute == endMinute {
		return errors.BadRequest("NOTIFICATION_PREFERENCE", "quiet start and end must differ")
	}
	return nil
}
