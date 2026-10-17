package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/notificationtemplate"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationTemplateRepo 通知模板仓储：可复用的标题/正文占位模板，平台全局，
// code 全局唯一（SendDirect 与测试渲染的引用锚）。
type NotificationTemplateRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[notifyV1.NotificationTemplate, ent.NotificationTemplate]

	repository *entCrud.Repository[
		ent.NotificationTemplateQuery, ent.NotificationTemplateSelect,
		ent.NotificationTemplateCreate, ent.NotificationTemplateCreateBulk,
		ent.NotificationTemplateUpdate, ent.NotificationTemplateUpdateOne,
		ent.NotificationTemplateDelete,
		predicate.NotificationTemplate,
		notifyV1.NotificationTemplate, ent.NotificationTemplate,
	]
}

func NewNotificationTemplateRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *NotificationTemplateRepo {
	repo := &NotificationTemplateRepo{
		log:       ctx.NewLoggerHelper("notification-template/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[notifyV1.NotificationTemplate, ent.NotificationTemplate](),
	}

	repo.repository = entCrud.NewRepository[
		ent.NotificationTemplateQuery, ent.NotificationTemplateSelect,
		ent.NotificationTemplateCreate, ent.NotificationTemplateCreateBulk,
		ent.NotificationTemplateUpdate, ent.NotificationTemplateUpdateOne,
		ent.NotificationTemplateDelete,
		predicate.NotificationTemplate,
		notifyV1.NotificationTemplate, ent.NotificationTemplate,
	](repo.mapper)

	repo.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	repo.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	return repo
}

func (r *NotificationTemplateRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationTemplateResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("NOTIFICATION_TEMPLATE", "invalid parameter")
	}

	builder := r.entClient.Client().NotificationTemplate.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		r.log.Errorf("list notification templates failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_TEMPLATE", "list notification templates failed")
	}
	if ret == nil {
		return &notifyV1.ListNotificationTemplateResponse{Total: 0, Items: nil}, nil
	}

	return &notifyV1.ListNotificationTemplateResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *NotificationTemplateRepo) Get(ctx context.Context, id uint32) (*notifyV1.NotificationTemplate, error) {
	if id == 0 {
		return nil, errors.BadRequest("NOTIFICATION_TEMPLATE", "id is required")
	}

	entity, err := r.entClient.Client().NotificationTemplate.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_TEMPLATE", "notification template not found")
		}
		r.log.Errorf("get notification template [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_TEMPLATE", "get notification template failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// GetByCode 按编码取模板；没有该行时返回 (nil, nil)——
// 与 GetByEventType 同一区分："没配这个模板"是配置事实，"查库失败"是基础设施故障，
// 两者压成同一种 error 会把 SendDirect 的排障方向引反。
func (r *NotificationTemplateRepo) GetByCode(ctx context.Context, code string) (*notifyV1.NotificationTemplate, error) {
	if code == "" {
		return nil, nil
	}

	entity, err := r.entClient.Client().NotificationTemplate.Query().
		Where(notificationtemplate.CodeEQ(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		r.log.Errorf("query notification template by code [%s] failed: %s", code, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_TEMPLATE", "query notification template failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// Create 建一条模板；code 全局唯一，撞码在写入前显式拒掉（唯一索引兜底）。
// created_at 显式写入，理由同规则表（mixin 列没有 DB 默认值）。
func (r *NotificationTemplateRepo) Create(ctx context.Context, req *notifyV1.NotificationTemplate, operatorID uint32) (uint32, error) {
	if req == nil {
		return 0, errors.BadRequest("NOTIFICATION_TEMPLATE", "invalid request")
	}
	if req.GetCode() == "" {
		return 0, errors.BadRequest("NOTIFICATION_TEMPLATE", "code is required")
	}

	exist, err := r.entClient.Client().NotificationTemplate.Query().
		Where(notificationtemplate.CodeEQ(req.GetCode())).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("check duplicated template code failed: %s", err.Error())
		return 0, errors.InternalServer("NOTIFICATION_TEMPLATE", "check notification template failed")
	}
	if exist {
		return 0, errors.BadRequest("NOTIFICATION_TEMPLATE", "template code "+req.GetCode()+" already exists")
	}

	builder := r.entClient.Client().NotificationTemplate.Create().
		SetNillableName(req.Name).
		SetNillableCode(req.Code).
		SetNillableTitleTemplate(req.TitleTemplate).
		SetNillableContentTemplate(req.ContentTemplate).
		SetNillableIsEnabled(req.IsEnabled).
		SetNillableRemark(req.Remark).
		SetCreatedAt(time.Now())

	if operatorID != 0 {
		builder.SetCreatedBy(operatorID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert notification template failed: %s", err.Error())
		return 0, errors.InternalServer("NOTIFICATION_TEMPLATE", "insert notification template failed")
	}

	return created.ID, nil
}

func (r *NotificationTemplateRepo) Update(ctx context.Context, req *notifyV1.UpdateNotificationTemplateRequest, operatorID uint32) (*notifyV1.NotificationTemplate, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("NOTIFICATION_TEMPLATE", "invalid request")
	}
	if req.GetId() == 0 {
		return nil, errors.BadRequest("NOTIFICATION_TEMPLATE", "id is required")
	}

	exist, err := r.entClient.Client().NotificationTemplate.Query().
		Where(notificationtemplate.IDEQ(req.GetId())).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_TEMPLATE", "notification template not found")
		}
		r.log.Errorf("get notification template [%d] failed: %s", req.GetId(), err.Error())
		return nil, errors.InternalServer("NOTIFICATION_TEMPLATE", "update notification template failed")
	}

	// code 若在更新里出现且与当前行不同，查重（唯一索引兜底并发窗口）
	if req.Data.Code != nil && exist.Code != nil && *exist.Code != req.Data.GetCode() {
		codeExist, err := r.entClient.Client().NotificationTemplate.Query().
			Where(notificationtemplate.CodeEQ(req.Data.GetCode())).
			Exist(ctx)
		if err != nil {
			r.log.Errorf("check duplicated template code failed: %s", err.Error())
			return nil, errors.InternalServer("NOTIFICATION_TEMPLATE", "check notification template failed")
		}
		if codeExist {
			return nil, errors.BadRequest("NOTIFICATION_TEMPLATE", "template code "+req.Data.GetCode()+" already exists")
		}
	}

	builder := r.entClient.Client().NotificationTemplate.UpdateOneID(exist.ID)
	if req.Data.Name != nil {
		builder.SetNillableName(req.Data.Name)
	}
	if req.Data.Code != nil {
		builder.SetNillableCode(req.Data.Code)
	}
	if req.Data.TitleTemplate != nil {
		builder.SetNillableTitleTemplate(req.Data.TitleTemplate)
	}
	if req.Data.ContentTemplate != nil {
		builder.SetNillableContentTemplate(req.Data.ContentTemplate)
	}
	if req.Data.IsEnabled != nil {
		builder.SetNillableIsEnabled(req.Data.IsEnabled)
	}
	if req.Data.Remark != nil {
		builder.SetNillableRemark(req.Data.Remark)
	}
	if operatorID != 0 {
		builder.SetUpdatedBy(operatorID)
	}
	builder.SetUpdatedAt(time.Now())

	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("update notification template [%d] failed: %s", req.GetId(), err.Error())
		return nil, errors.InternalServer("NOTIFICATION_TEMPLATE", "update notification template failed")
	}

	return r.Get(ctx, req.GetId())
}

func (r *NotificationTemplateRepo) Delete(ctx context.Context, id uint32) error {
	if id == 0 {
		return errors.BadRequest("NOTIFICATION_TEMPLATE", "id is required")
	}

	if err := r.entClient.Client().NotificationTemplate.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("NOTIFICATION_TEMPLATE", "notification template not found")
		}
		r.log.Errorf("delete notification template [%d] failed: %s", id, err.Error())
		return errors.InternalServer("NOTIFICATION_TEMPLATE", "delete notification template failed")
	}

	return nil
}
