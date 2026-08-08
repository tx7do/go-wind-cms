package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/loginpolicy"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

type LoginPolicyRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper          *mapper.CopierMapper[authenticationV1.LoginPolicy, ent.LoginPolicy]
	typeConverter   *mapper.EnumTypeConverter[authenticationV1.LoginPolicy_Type, loginpolicy.Type]
	methodConverter *mapper.EnumTypeConverter[authenticationV1.LoginPolicy_Method, loginpolicy.Method]

	repository *entCrud.Repository[
		ent.LoginPolicyQuery, ent.LoginPolicySelect,
		ent.LoginPolicyCreate, ent.LoginPolicyCreateBulk,
		ent.LoginPolicyUpdate, ent.LoginPolicyUpdateOne,
		ent.LoginPolicyDelete,
		predicate.LoginPolicy,
		authenticationV1.LoginPolicy, ent.LoginPolicy,
	]
}

func NewLoginPolicyRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *LoginPolicyRepo {
	repo := &LoginPolicyRepo{
		log:       ctx.NewLoggerHelper("login-policy/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[authenticationV1.LoginPolicy, ent.LoginPolicy](),
		typeConverter: mapper.NewEnumTypeConverter[authenticationV1.LoginPolicy_Type, loginpolicy.Type](
			authenticationV1.LoginPolicy_Type_name, authenticationV1.LoginPolicy_Type_value,
		),
		methodConverter: mapper.NewEnumTypeConverter[authenticationV1.LoginPolicy_Method, loginpolicy.Method](
			authenticationV1.LoginPolicy_Method_name, authenticationV1.LoginPolicy_Method_value,
		),
	}

	repo.init()

	return repo
}

func (r *LoginPolicyRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.LoginPolicyQuery, ent.LoginPolicySelect,
		ent.LoginPolicyCreate, ent.LoginPolicyCreateBulk,
		ent.LoginPolicyUpdate, ent.LoginPolicyUpdateOne,
		ent.LoginPolicyDelete,
		predicate.LoginPolicy,
		authenticationV1.LoginPolicy, ent.LoginPolicy,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.typeConverter.NewConverterPair())
	r.mapper.AppendConverters(r.methodConverter.NewConverterPair())
}

func (r *LoginPolicyRepo) Count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().LoginPolicy.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
		return 0, authenticationV1.ErrorInternalServerError("query count failed")
	}

	return count, nil
}

func (r *LoginPolicyRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*authenticationV1.ListLoginPolicyResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().LoginPolicy.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &authenticationV1.ListLoginPolicyResponse{Total: 0, Items: nil}, nil
	}

	return &authenticationV1.ListLoginPolicyResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *LoginPolicyRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().LoginPolicy.Query().
		Where(loginpolicy.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query exist failed: %s", err.Error())
		return false, authenticationV1.ErrorInternalServerError("query exist failed")
	}
	return exist, nil
}

func (r *LoginPolicyRepo) Get(ctx context.Context, req *authenticationV1.GetLoginPolicyRequest) (*authenticationV1.LoginPolicy, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().LoginPolicy.Query()

	var whereCond []func(s *sql.Selector)
	switch req.QueryBy.(type) {
	default:
	case *authenticationV1.GetLoginPolicyRequest_Id:
		whereCond = append(whereCond, loginpolicy.IDEQ(req.GetId()))
	}

	dto, err := r.repository.Get(ctx, builder, req.GetViewMask(), whereCond...)
	if err != nil {
		return nil, err
	}

	return dto, err
}

func (r *LoginPolicyRepo) Create(ctx context.Context, req *authenticationV1.CreateLoginPolicyRequest) error {
	if req == nil || req.Data == nil {
		return authenticationV1.ErrorBadRequest("invalid request")
	}

	builder := r.entClient.Client().LoginPolicy.Create().
		SetNillableTenantID(req.Data.TenantId).
		SetNillableTargetID(req.Data.TargetId).
		SetNillableType(r.typeConverter.ToEntity(req.Data.Type)).
		SetNillableMethod(r.methodConverter.ToEntity(req.Data.Method)).
		SetNillableValue(req.Data.Value).
		SetNillableReason(req.Data.Reason).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("insert admin login restriction failed: %s", err.Error())
		return authenticationV1.ErrorInternalServerError("insert admin login restriction failed")
	}

	return nil
}

func (r *LoginPolicyRepo) Update(ctx context.Context, req *authenticationV1.UpdateLoginPolicyRequest) error {
	if req == nil || req.Data == nil {
		return authenticationV1.ErrorBadRequest("invalid request")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.IsExist(ctx, req.GetId())
		if err != nil {
			return err
		}
		if !exist {
			createReq := &authenticationV1.CreateLoginPolicyRequest{Data: req.Data}
			createReq.Data.CreatedBy = createReq.Data.UpdatedBy
			createReq.Data.UpdatedBy = nil
			return r.Create(ctx, createReq)
		}
	}

	builder := r.entClient.Client().Debug().LoginPolicy.Update()
	// 租户作用域：仅更新本租户登录策略，避免跨租户改他人策略（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(loginpolicy.TenantIDEQ(tid))
	}
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	err := r.repository.UpdateX(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *authenticationV1.LoginPolicy) {
			builder.
				SetNillableTargetID(req.Data.TargetId).
				SetNillableType(r.typeConverter.ToEntity(req.Data.Type)).
				SetNillableMethod(r.methodConverter.ToEntity(req.Data.Method)).
				SetNillableValue(req.Data.Value).
				SetNillableReason(req.Data.Reason).
				SetUpdatedAt(time.Now())

			// updated_by 强制取调用者身份，保证审计归属真实，忽略客户端值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(loginpolicy.FieldID, req.GetId()))
		},
	)

	return err
}

func (r *LoginPolicyRepo) Delete(ctx context.Context, req *authenticationV1.DeleteLoginPolicyRequest) error {
	if req == nil {
		return authenticationV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Debug().LoginPolicy.Delete()
	// 租户作用域：仅删除本租户登录策略，避免跨租户删他人策略（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(loginpolicy.TenantIDEQ(tid))
	}
	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		s.Where(sql.EQ(loginpolicy.FieldID, req.GetId()))
	})
	if err != nil {
		r.log.Errorf("delete internal message categories failed: %s", err.Error())
		return authenticationV1.ErrorInternalServerError("delete admin login restriction failed")
	}

	return nil
}
