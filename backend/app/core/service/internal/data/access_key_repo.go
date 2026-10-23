package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"entgo.io/ent/dialect/sql"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	appViewer "go-wind-cms/pkg/entgo/viewer"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/accesskey"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	accesskeyV1 "go-wind-cms/api/gen/go/access_key/service/v1"
)

// AccessKeyRepo OpenAPI 访问凭证（AK/SK）仓储。
// 读路径租户隔离由 ent TenantPrivacy 查询策略注入；变更路径（Update/Delete）
// TenantPrivacy 不覆盖，必须显式加租户谓词（见 Update/Delete）；
// 令牌交换（免鉴权流程）的查询经 SystemViewer 旁路租户 scope（见 *BySystem 方法）。
type AccessKeyRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
	mapper    *mapper.CopierMapper[accesskeyV1.AccessKey, ent.AccessKey]
	// 状态枚举转换（ent SwitchStatus ON/OFF <-> proto AccessKey_Status）
	statusConverter *mapper.EnumTypeConverter[accesskeyV1.AccessKey_Status, accesskey.Status]

	repository *entCrud.Repository[
		ent.AccessKeyQuery, ent.AccessKeySelect,
		ent.AccessKeyCreate, ent.AccessKeyCreateBulk,
		ent.AccessKeyUpdate, ent.AccessKeyUpdateOne,
		ent.AccessKeyDelete,
		predicate.AccessKey,
		accesskeyV1.AccessKey, ent.AccessKey,
	]
}

func NewAccessKeyRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *AccessKeyRepo {
	r := &AccessKeyRepo{
		log:             ctx.NewLoggerHelper("access-key/repo/core-service"),
		entClient:       entClient,
		mapper:          mapper.NewCopierMapper[accesskeyV1.AccessKey, ent.AccessKey](),
		statusConverter: mapper.NewEnumTypeConverter[accesskeyV1.AccessKey_Status, accesskey.Status](accesskeyV1.AccessKey_Status_name, accesskeyV1.AccessKey_Status_value),
	}
	r.init()
	return r
}

func (r *AccessKeyRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.AccessKeyQuery, ent.AccessKeySelect,
		ent.AccessKeyCreate, ent.AccessKeyCreateBulk,
		ent.AccessKeyUpdate, ent.AccessKeyUpdateOne,
		ent.AccessKeyDelete,
		predicate.AccessKey,
		accesskeyV1.AccessKey, ent.AccessKey,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
}

func (r *AccessKeyRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*accesskeyV1.ListAccessKeyResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("ACCESS_KEY", "invalid parameter")
	}

	builder := r.entClient.Client().AccessKey.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &accesskeyV1.ListAccessKeyResponse{Total: 0, Items: nil}, nil
	}

	return &accesskeyV1.ListAccessKeyResponse{Total: ret.Total, Items: ret.Items}, nil
}

func (r *AccessKeyRepo) Get(ctx context.Context, req *accesskeyV1.GetAccessKeyRequest) (*accesskeyV1.AccessKey, error) {
	if req == nil {
		return nil, errors.BadRequest("ACCESS_KEY", "invalid parameter")
	}

	builder := r.entClient.Client().AccessKey.Query()

	var whereCond []func(s *sql.Selector)
	switch req.QueryBy.(type) {
	default:
	case *accesskeyV1.GetAccessKeyRequest_Id:
		whereCond = append(whereCond, accesskey.IDEQ(req.GetId()))
	case *accesskeyV1.GetAccessKeyRequest_AccessKey:
		whereCond = append(whereCond, accesskey.AccessKeyEQ(req.GetAccessKey()))
	}

	return r.repository.Get(ctx, builder, req.GetViewMask(), whereCond...)
}

// Create 创建凭证。AK 与 secret 摘要由 service 层生成后传入
// （proto DTO 不携带 secret 摘要，不能走通用 DTO 映射落库）。
// 返回落库后的实体供回显（id/AK/时间戳）。
func (r *AccessKeyRepo) Create(ctx context.Context, data *accesskeyV1.AccessKey, accessKey, secretHash string, operatorID uint32) (*ent.AccessKey, error) {
	builder := r.entClient.Client().AccessKey.Create().
		SetNillableName(data.Name).
		SetAccessKey(accessKey).
		SetSecretHash(secretHash).
		SetNillableExpiresAt(tsOf(data.ExpiresAt)).
		SetCreatedAt(time.Now())

	// 未指定状态时默认启用（与 ent schema 的 ON 默认一致）
	if data.Status == nil || *data.Status == accesskeyV1.AccessKey_ON {
		builder.SetStatus(accesskey.StatusOn)
	} else {
		builder.SetStatus(accesskey.StatusOff)
	}

	if data.TenantId != nil {
		builder.SetTenantID(*data.TenantId)
	}
	if operatorID != 0 {
		builder.SetCreatedBy(operatorID)
	}

	entity, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert access key failed: %s", err.Error())
		return nil, errors.InternalServer("ACCESS_KEY", "insert access key failed")
	}
	return entity, nil
}

// Update 更新凭证。可编辑字段白名单只有 name/status/expires_at——AK、
// secret 摘要与租户归属按设计不可变（上游 e852daba 的教训：黑名单放行
// 会把 DTO 上不存在的 secret_hash 一起写，掩码管线直接炸）。
// mask 为空时更新全部可编辑字段（上游 nil 掩码同语义）。
func (r *AccessKeyRepo) Update(ctx context.Context, req *accesskeyV1.UpdateAccessKeyRequest, operatorID uint32) error {
	if req == nil || req.Data == nil {
		return errors.BadRequest("ACCESS_KEY", "invalid parameter")
	}
	if req.GetId() == 0 {
		return errors.BadRequest("ACCESS_KEY", "id is required")
	}

	builder := r.entClient.Client().AccessKey.Update().Where(accesskey.IDEQ(req.GetId()))

	// 租户作用域：Update 是 ent 变更路径，TenantPrivacy 不覆盖，必须显式加租户过滤，
	// 否则 WHERE 仅含 id，跨租户可改他人 AccessKey（与其余 *_repo 的 Update 写法对齐）。
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(accesskey.TenantIDEQ(tid))
	}

	applyField := func(path string) bool {
		switch path {
		case "name":
			builder.SetNillableName(req.Data.Name)
			return true
		case "status":
			if req.Data.Status != nil {
				if stEnt := r.statusConverter.ToEntity(req.Data.Status); stEnt != nil {
					builder.SetStatus(*stEnt)
				}
			}
			return true
		case "expiresAt", "expires_at":
			builder.SetNillableExpiresAt(tsOf(req.Data.ExpiresAt))
			return true
		}
		return false
	}

	if req.UpdateMask != nil {
		for _, p := range req.UpdateMask.Paths {
			applyField(p)
		}
	} else {
		for _, p := range []string{"name", "status", "expiresAt"} {
			applyField(p)
		}
	}

	builder.SetUpdatedAt(time.Now())
	if operatorID != 0 {
		builder.SetUpdatedBy(operatorID)
	}

	// 白名单没命中任何字段时只落操作人/时间戳（上游"纯不可变掩码=仅盖章空操作"同语义）

	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("update access key failed (id=%d): %s", req.GetId(), err.Error())
		return errors.InternalServer("ACCESS_KEY", "update access key failed")
	}
	return nil
}

func (r *AccessKeyRepo) Delete(ctx context.Context, req *accesskeyV1.DeleteAccessKeyRequest) error {
	if req == nil {
		return errors.BadRequest("ACCESS_KEY", "invalid parameter")
	}

	// 租户作用域：Delete 是 ent 变更路径，TenantPrivacy 不覆盖，必须显式加租户过滤，
	// 否则仅按 id 删除，跨租户可删他人 AccessKey。
	builder := r.entClient.Client().AccessKey.Delete().
		Where(accesskey.IDEQ(req.GetId()))
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(accesskey.TenantIDEQ(tid))
	}

	if _, err := builder.Exec(ctx); err != nil {
		r.log.Errorf("delete access key failed (id=%d): %s", req.GetId(), err.Error())
		return errors.InternalServer("ACCESS_KEY", "delete access key failed")
	}
	return nil
}

// UpdateSecretHash 重置密钥摘要（轮换）。旧的已签发机器令牌在自身过期前仍有效。
func (r *AccessKeyRepo) UpdateSecretHash(ctx context.Context, id uint32, secretHash string) error {
	err := r.entClient.Client().AccessKey.UpdateOneID(id).
		SetSecretHash(secretHash).
		SetUpdatedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("update access key secret hash failed (id=%d): %s", id, err.Error())
		return errors.InternalServer("ACCESS_KEY", "update secret failed")
	}
	return nil
}

// GetByAccessKeyBySystem 按访问键查询凭证（系统旁路视图，绕过租户 scope）——
// 仅供令牌交换（免鉴权流程）使用：交换发生时请求尚无任何 viewer。
func (r *AccessKeyRepo) GetByAccessKeyBySystem(ctx context.Context, accessKey string) (*ent.AccessKey, error) {
	svCtx := appViewer.NewSystemViewerContext(ctx)
	entity, err := r.entClient.Client().AccessKey.Query().
		Where(accesskey.AccessKeyEQ(accessKey)).
		Only(svCtx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// TouchLastUsedBySystem 刷新最近使用时间（系统旁路；尽力而为，失败不影响交换）。
func (r *AccessKeyRepo) TouchLastUsedBySystem(ctx context.Context, id uint32) {
	svCtx := appViewer.NewSystemViewerContext(ctx)
	err := r.entClient.Client().AccessKey.UpdateOneID(id).
		SetLastUsedAt(time.Now()).
		Exec(svCtx)
	if err != nil {
		r.log.Warnf("touch access key [%d] last_used_at failed: %s", id, err.Error())
	}
}
