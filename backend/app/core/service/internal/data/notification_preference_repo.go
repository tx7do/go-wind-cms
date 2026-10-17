package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/notificationpreference"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationPreferenceRepo 用户通知偏好仓储。
// schema 刻意无租户 mixin（偏好跟人不跟租户），读锚是 user_id；全链 fail-open：
// 这里返回 (nil, nil) 时由调用方按"未配置"处理（不静音不退订）。
type NotificationPreferenceRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewNotificationPreferenceRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *NotificationPreferenceRepo {
	return &NotificationPreferenceRepo{
		log:       ctx.NewLoggerHelper("notification-preference/repo/core-service"),
		entClient: entClient,
	}
}

func toNotificationPreferenceDTO(entity *ent.NotificationPreference) *notifyV1.NotificationPreference {
	if entity == nil {
		return nil
	}
	return &notifyV1.NotificationPreference{
		UserId:           trans.Ptr(entity.UserID),
		QuietEnabled:     trans.Ptr(entity.QuietEnabled),
		QuietStartMinute: trans.Ptr(entity.QuietStartMinute),
		QuietEndMinute:   trans.Ptr(entity.QuietEndMinute),
		MutedCategoryIds: entity.MutedCategoryIds,
		UpdatedAt:        timeutil.TimeToTimestamppb(entity.UpdatedAt),
	}
}

// GetByUserID 读某用户的偏好；从未配置过时返回 (nil, nil)，由调用方决定默认值。
func (r *NotificationPreferenceRepo) GetByUserID(ctx context.Context, userID uint32) (*notifyV1.NotificationPreference, error) {
	entity, err := r.entClient.Client().NotificationPreference.Query().
		Where(notificationpreference.UserIDEQ(userID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		r.log.Errorf("query notification preference by user [%d] failed: %s", userID, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_PREFERENCE", "query notification preference failed")
	}

	return toNotificationPreferenceDTO(entity), nil
}

// ListByUserIDs 批量读一组用户的偏好（投递链逐收件人判定静音时用），缺偏好的用户不在返回值里。
func (r *NotificationPreferenceRepo) ListByUserIDs(ctx context.Context, userIDs []uint32) (map[uint32]*notifyV1.NotificationPreference, error) {
	if len(userIDs) == 0 {
		return map[uint32]*notifyV1.NotificationPreference{}, nil
	}

	entities, err := r.entClient.Client().NotificationPreference.Query().
		Where(notificationpreference.UserIDIn(userIDs...)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query notification preferences by users failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_PREFERENCE", "query notification preferences failed")
	}

	result := make(map[uint32]*notifyV1.NotificationPreference, len(entities))
	for _, entity := range entities {
		result[entity.UserID] = toNotificationPreferenceDTO(entity)
	}
	return result, nil
}

// Upsert 按用户保存偏好（存在则更新，不存在则创建）；userID 由调用方从操作人钉定。
func (r *NotificationPreferenceRepo) Upsert(ctx context.Context, userID uint32, pref *notifyV1.UpdateNotificationPreferenceRequest, operatorUserID uint32) (*notifyV1.NotificationPreference, error) {
	if pref == nil {
		return nil, errors.BadRequest("NOTIFICATION_PREFERENCE", "invalid parameter")
	}

	existing, err := r.entClient.Client().NotificationPreference.Query().
		Where(notificationpreference.UserIDEQ(userID)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		r.log.Errorf("query notification preference for upsert failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_PREFERENCE", "save notification preference failed")
	}

	muted := make([]uint32, 0, len(pref.GetMutedCategoryIds()))
	muted = append(muted, pref.GetMutedCategoryIds()...)

	if existing == nil {
		entity, err := r.entClient.Client().NotificationPreference.Create().
			SetUserID(userID).
			SetQuietEnabled(pref.GetQuietEnabled()).
			SetQuietStartMinute(pref.GetQuietStartMinute()).
			SetQuietEndMinute(pref.GetQuietEndMinute()).
			SetMutedCategoryIds(muted).
			SetNillableCreatedBy(operatorUserIDOrNil(operatorUserID)).
			Save(ctx)
		if err != nil {
			r.log.Errorf("create notification preference for user [%d] failed: %s", userID, err.Error())
			return nil, errors.InternalServer("NOTIFICATION_PREFERENCE", "save notification preference failed")
		}
		return toNotificationPreferenceDTO(entity), nil
	}

	entity, err := existing.Update().
		SetQuietEnabled(pref.GetQuietEnabled()).
		SetQuietStartMinute(pref.GetQuietStartMinute()).
		SetQuietEndMinute(pref.GetQuietEndMinute()).
		SetMutedCategoryIds(muted).
		SetNillableUpdatedBy(operatorUserIDOrNil(operatorUserID)).
		Save(ctx)
	if err != nil {
		r.log.Errorf("update notification preference for user [%d] failed: %s", userID, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_PREFERENCE", "save notification preference failed")
	}

	return toNotificationPreferenceDTO(entity), nil
}

func operatorUserIDOrNil(userID uint32) *uint32 {
	if userID == 0 {
		return nil
	}
	return trans.Ptr(userID)
}
