package data

import (
	"context"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	kratosMetadata "github.com/go-kratos/kratos/v2/metadata"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/go-utils/mapper"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/socialaccount"

	socialV1 "go-wind-cms/api/gen/go/social/service/v1"
)

// SocialAccountRepo 社交账号仓储。
// credentials 的加解密在 service 层完成，本仓储只负责持久化（字段原样存取）。
type SocialAccountRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[socialV1.SocialAccount, ent.SocialAccount]

	repository *entCrud.Repository[
		ent.SocialAccountQuery, ent.SocialAccountSelect,
		ent.SocialAccountCreate, ent.SocialAccountCreateBulk,
		ent.SocialAccountUpdate, ent.SocialAccountUpdateOne,
		ent.SocialAccountDelete,
		predicate.SocialAccount,
		socialV1.SocialAccount, ent.SocialAccount,
	]

	platformConverter *mapper.EnumTypeConverter[socialV1.Platform, socialaccount.Platform]
}

func NewSocialAccountRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SocialAccountRepo {
	repo := &SocialAccountRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("social-account/repo/core-service"),
		mapper:    mapper.NewCopierMapper[socialV1.SocialAccount, ent.SocialAccount](),
		platformConverter: mapper.NewEnumTypeConverter[socialV1.Platform, socialaccount.Platform](
			socialV1.Platform_name, socialV1.Platform_value,
		),
	}

	repo.init()

	return repo
}

func (r *SocialAccountRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.SocialAccountQuery, ent.SocialAccountSelect,
		ent.SocialAccountCreate, ent.SocialAccountCreateBulk,
		ent.SocialAccountUpdate, ent.SocialAccountUpdateOne,
		ent.SocialAccountDelete,
		predicate.SocialAccount,
		socialV1.SocialAccount, ent.SocialAccount,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	r.mapper.AppendConverters(r.platformConverter.NewConverterPair())
}

// List 租户过滤的账号列表（tid<=0 的平台视角不过滤，可见全部）。
func (r *SocialAccountRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*socialV1.ListSocialAccountsResponse, error) {
	if req == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SocialAccount.Query()
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(socialaccount.TenantIDEQ(tid))
	}

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &socialV1.ListSocialAccountsResponse{Total: 0, Items: nil}, nil
	}

	return &socialV1.ListSocialAccountsResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

// Get 按 ID 取账号（管理面，租户过滤）。
func (r *SocialAccountRepo) Get(ctx context.Context, req *socialV1.GetSocialAccountRequest) (*socialV1.SocialAccount, error) {
	if req == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SocialAccount.Query()

	var whereConds []func(*sql.Selector)
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		whereConds = append(whereConds, func(sels *sql.Selector) {
			sels.Where(sql.EQ(socialaccount.FieldTenantID, tid))
		})
	}

	dto, err := r.repository.Get(ctx, builder, req.GetViewMask(), whereConds...)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, socialV1.ErrorAccountNotFound("social account not found")
		}
		r.log.Errorf("query social account failed: %s", err.Error())
		return nil, socialV1.ErrorInternalServerError("query social account failed")
	}

	return dto, nil
}

// GetByID 内部读取（发布 worker 的 SystemViewer 上下文，跨租户）。
func (r *SocialAccountRepo) GetByID(ctx context.Context, id uint32) (*socialV1.SocialAccount, error) {
	entity, err := r.entClient.Client().SocialAccount.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, socialV1.ErrorAccountNotFound("social account not found")
		}
		r.log.Errorf("query social account by id failed: %s", err.Error())
		return nil, socialV1.ErrorInternalServerError("query social account failed")
	}
	return r.mapper.ToDTO(entity), nil
}

// Create 落库账号。credentials/hint 由 service 层加密后放进 req.Data 传入。
func (r *SocialAccountRepo) Create(ctx context.Context, req *socialV1.CreateSocialAccountRequest) (*socialV1.SocialAccount, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SocialAccount.Create().
		SetTenantID(socialTenantIDFromContext(ctx)).
		SetNillablePlatform(r.platformConverter.ToEntity(req.Data.Platform)).
		SetNillableName(req.Data.Name).
		SetNillableHandle(req.Data.Handle).
		SetNillableCredentials(req.Data.Credentials).
		SetNillableCredentialHint(req.Data.CredentialHint).
		SetNillableIsEnabled(req.Data.IsEnabled).
		SetNillableRemark(req.Data.Remark).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert social account failed: %s", err.Error())
		return nil, socialV1.ErrorInternalServerError("insert social account failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// Update 更新账号（租户守卫 + updated_by 由 viewer 推导）。
func (r *SocialAccountRepo) Update(ctx context.Context, req *socialV1.UpdateSocialAccountRequest) (*socialV1.SocialAccount, error) {
	if req == nil || req.Data == nil {
		return nil, socialV1.ErrorBadRequest("invalid parameter")
	}

	// 不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.entClient.Client().SocialAccount.Query().
			Where(socialaccount.IDEQ(req.GetId())).
			Exist(ctx)
		if err != nil {
			r.log.Errorf("query social account exist failed: %s", err.Error())
			return nil, socialV1.ErrorInternalServerError("query social account exist failed")
		}
		if !exist {
			req.Data.CreatedBy = req.Data.UpdatedBy
			req.Data.UpdatedBy = nil
			_, err = r.Create(ctx, &socialV1.CreateSocialAccountRequest{Data: req.Data})
			return nil, err
		}
	}

	tid, hasTenant := maybeTenantFromViewer(ctx)
	callerUserID, hasUser := viewerUserIDFromContext(ctx)

	builder := r.entClient.Client().SocialAccount.UpdateOneID(req.GetId())
	builder.Where(socialaccount.IDEQ(req.GetId()))
	if hasTenant {
		builder.Where(socialaccount.TenantIDEQ(tid))
	}

	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *socialV1.SocialAccount) {
			builder.
				SetNillablePlatform(r.platformConverter.ToEntity(req.Data.Platform)).
				SetNillableName(req.Data.Name).
				SetNillableHandle(req.Data.Handle).
				SetNillableCredentials(req.Data.Credentials).
				SetNillableCredentialHint(req.Data.CredentialHint).
				SetNillableIsEnabled(req.Data.IsEnabled).
				SetNillableRemark(req.Data.Remark).
				SetUpdatedAt(time.Now())

			// updated_by 强制由服务端 viewer context 推导，忽略客户端传入值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(socialaccount.FieldID, req.GetId()))
		},
	)

	return result, err
}

// Delete 删除账号（租户守卫）。分发表不级联删除，保留历史记录；
// service 层负责先把引用该账号的未发布分发置为失效。
func (r *SocialAccountRepo) Delete(ctx context.Context, req *socialV1.DeleteSocialAccountRequest) error {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	delBuilder := r.entClient.Client().SocialAccount.Delete()
	delBuilder.Where(socialaccount.IDEQ(req.GetId()))
	if hasTenant {
		delBuilder.Where(socialaccount.TenantIDEQ(tid))
	}
	_, err := delBuilder.Exec(ctx)
	if err != nil {
		r.log.Errorf("delete social account failed: %s", err.Error())
	}
	return err
}

// ListEnabledByPlatforms 收件箱同步用：指定平台上全部启用中的账号（SystemViewer 跨租户）。
func (r *SocialAccountRepo) ListEnabledByPlatforms(ctx context.Context, platforms []socialV1.Platform) ([]*socialV1.SocialAccount, error) {
	ents := make([]socialaccount.Platform, 0, len(platforms))
	for _, p := range platforms {
		if e := r.platformConverter.ToEntity(trans.Ptr(p)); e != nil {
			ents = append(ents, *e)
		}
	}
	if len(ents) == 0 {
		return nil, nil
	}

	rows, err := r.entClient.Client().SocialAccount.Query().
		Where(
			socialaccount.IsEnabled(true),
			socialaccount.PlatformIn(ents...),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("list enabled accounts failed: %s", err.Error())
		return nil, err
	}

	out := make([]*socialV1.SocialAccount, 0, len(rows))
	for _, row := range rows {
		out = append(out, r.mapper.ToDTO(row))
	}
	return out, nil
}

// socialTenantIDFromContext 解析写入侧租户归属：viewer 优先，
// 回落内网可信元数据 x-md-global-tenant-id（与评论链路口径一致），缺省 0（平台）。
func socialTenantIDFromContext(ctx context.Context) uint32 {
	if tid, ok := maybeTenantFromViewer(ctx); ok {
		return tid
	}
	if md, ok := kratosMetadata.FromServerContext(ctx); ok {
		if v := md.Get("x-md-global-tenant-id"); v != "" {
			if parsed, perr := strconv.ParseUint(v, 10, 32); perr == nil {
				return uint32(parsed)
			}
		}
	}
	return 0
}
