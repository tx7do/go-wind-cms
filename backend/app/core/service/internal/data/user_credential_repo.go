package data

import (
	"context"
	"encoding/json"

	"encoding/base64"
	"go-wind-cms/pkg/constants"
	passwordpolicy "go-wind-cms/pkg/password"
	"strconv"
	"time"

	pkgcrypto "go-wind-cms/pkg/crypto"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-crud/viewer"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/crypto"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/password"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/usercredential"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	identityV1 "go-wind-cms/api/gen/go/identity/service/v1"
)

type UserCredentialRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper                  *mapper.CopierMapper[authenticationV1.UserCredential, ent.UserCredential]
	statusConverter         *mapper.EnumTypeConverter[authenticationV1.UserCredential_Status, usercredential.Status]
	identityTypeConverter   *mapper.EnumTypeConverter[authenticationV1.UserCredential_IdentityType, usercredential.IdentityType]
	credentialTypeConverter *mapper.EnumTypeConverter[authenticationV1.UserCredential_CredentialType, usercredential.CredentialType]

	passwordCrypto password.Crypto
	configRepo     *ConfigRepo

	repository *entCrud.Repository[
		ent.UserCredentialQuery, ent.UserCredentialSelect,
		ent.UserCredentialCreate, ent.UserCredentialCreateBulk,
		ent.UserCredentialUpdate, ent.UserCredentialUpdateOne,
		ent.UserCredentialDelete,
		predicate.UserCredential,
		authenticationV1.UserCredential, ent.UserCredential,
	]
}

func NewUserCredentialRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client], passwordCrypto password.Crypto, configRepo *ConfigRepo) *UserCredentialRepo {
	repo := &UserCredentialRepo{
		log:                     ctx.NewLoggerHelper("user-credentials/repo/core-service"),
		entClient:               entClient,
		passwordCrypto:          passwordCrypto,
		configRepo:              configRepo,
		mapper:                  mapper.NewCopierMapper[authenticationV1.UserCredential, ent.UserCredential](),
		statusConverter:         mapper.NewEnumTypeConverter[authenticationV1.UserCredential_Status, usercredential.Status](authenticationV1.UserCredential_Status_name, authenticationV1.UserCredential_Status_value),
		identityTypeConverter:   mapper.NewEnumTypeConverter[authenticationV1.UserCredential_IdentityType, usercredential.IdentityType](authenticationV1.UserCredential_IdentityType_name, authenticationV1.UserCredential_IdentityType_value),
		credentialTypeConverter: mapper.NewEnumTypeConverter[authenticationV1.UserCredential_CredentialType, usercredential.CredentialType](authenticationV1.UserCredential_CredentialType_name, authenticationV1.UserCredential_CredentialType_value),
	}

	repo.init()

	return repo
}

func (r *UserCredentialRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.UserCredentialQuery, ent.UserCredentialSelect,
		ent.UserCredentialCreate, ent.UserCredentialCreateBulk,
		ent.UserCredentialUpdate, ent.UserCredentialUpdateOne,
		ent.UserCredentialDelete,
		predicate.UserCredential,
		authenticationV1.UserCredential, ent.UserCredential,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
	r.mapper.AppendConverters(r.identityTypeConverter.NewConverterPair())
	r.mapper.AppendConverters(r.credentialTypeConverter.NewConverterPair())
}

func (r *UserCredentialRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().UserCredential.Query().
		Where(usercredential.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query exist failed: %s", err.Error())
		return false, authenticationV1.ErrorInternalServerError("query exist failed")
	}
	return exist, nil
}

func (r *UserCredentialRepo) Count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().UserCredential.Query()
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

func (r *UserCredentialRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*authenticationV1.ListUserCredentialResponse, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().UserCredential.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &authenticationV1.ListUserCredentialResponse{Total: 0, Items: nil}, nil
	}

	return &authenticationV1.ListUserCredentialResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *UserCredentialRepo) Create(ctx context.Context, req *authenticationV1.CreateUserCredentialRequest) (err error) {
	if req == nil || req.Data == nil {
		return identityV1.ErrorBadRequest("invalid parameter")
	}

	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return identityV1.ErrorInternalServerError("start transaction failed")
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.Errorf("transaction rollback failed: %s", rollbackErr.Error())
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			r.log.Errorf("transaction commit failed: %s", commitErr.Error())
			err = identityV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	return r.CreateWithTx(ctx, tx, req.GetData())
}

func (r *UserCredentialRepo) CreateWithTx(ctx context.Context, tx *ent.Tx, data *authenticationV1.UserCredential) error {
	if data == nil {
		return authenticationV1.ErrorBadRequest("invalid request")
	}

	var err error

	if data.Credential != nil {
		var newCredential string
		newCredential, err = r.prepareCredential(r.credentialTypeConverter.ToEntity(data.CredentialType), data.GetCredential())
		if err != nil {
			r.log.Errorf("prepare new credential failed: %s", err.Error())
			return authenticationV1.ErrorBadRequest("prepare new credential failed")
		}
		data.Credential = trans.Ptr(newCredential)
	}

	builder := tx.UserCredential.Create()
	builder.
		SetUserID(data.GetUserId()).
		SetNillableTenantID(data.TenantId).
		SetNillableIdentityType(r.identityTypeConverter.ToEntity(data.IdentityType)).
		SetNillableIdentifier(data.Identifier).
		SetNillableCredentialType(r.credentialTypeConverter.ToEntity(data.CredentialType)).
		SetNillableCredential(data.Credential).
		SetNillableIsPrimary(data.IsPrimary).
		SetNillableStatus(r.statusConverter.ToEntity(data.Status)).
		SetNillableExtraInfo(data.ExtraInfo).
		SetNillableProvider(data.Provider).
		SetNillableProviderAccountID(data.ProviderAccountId).
		SetCreatedAt(time.Now())

	if err = builder.Exec(ctx); err != nil {
		r.log.Errorf("insert user credential failed: %s [%v]", err.Error(), data)
		return authenticationV1.ErrorInternalServerError("insert user credential failed")
	}

	return nil
}

func (r *UserCredentialRepo) Update(ctx context.Context, req *authenticationV1.UpdateUserCredentialRequest) error {
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
			err = r.Create(ctx, &authenticationV1.CreateUserCredentialRequest{Data: req.Data})
			return err
		}
	}

	var err error

	if req.Data.Credential != nil {
		var newCredential string
		newCredential, err = r.prepareCredential(r.credentialTypeConverter.ToEntity(req.Data.CredentialType), req.Data.GetCredential())
		if err != nil {
			r.log.Errorf("prepare new credential failed: %s", err.Error())
			return authenticationV1.ErrorBadRequest("prepare new credential failed")
		}
		req.Data.Credential = trans.Ptr(newCredential)
	}

	builder := r.entClient.Client().Debug().UserCredential.Update()
	err = r.repository.UpdateX(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *authenticationV1.UserCredential) {
			builder.
				SetNillableIdentityType(r.identityTypeConverter.ToEntity(req.Data.IdentityType)).
				SetNillableIdentifier(req.Data.Identifier).
				SetNillableCredentialType(r.credentialTypeConverter.ToEntity(req.Data.CredentialType)).
				SetNillableCredential(req.Data.Credential).
				SetNillableIsPrimary(req.Data.IsPrimary).
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableExtraInfo(req.Data.ExtraInfo).
				SetNillableProvider(req.Data.Provider).
				SetNillableProviderAccountID(req.Data.ProviderAccountId).
				SetUpdatedAt(time.Now())
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(usercredential.FieldID, req.GetId()))
		},
	)

	return err
}

func (r *UserCredentialRepo) Delete(ctx context.Context, id uint32) error {
	builder := r.entClient.Client().UserCredential.Delete()
	builder.Where(usercredential.IDEQ(id))
	if affected, err := builder.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return authenticationV1.ErrorNotFound("user credential not found")
		}

		r.log.Errorf("delete one data failed: %s", err.Error())

		return authenticationV1.ErrorInternalServerError("delete one data failed")
	} else {
		if affected == 0 {
			return authenticationV1.ErrorNotFound("user credential not found")
		} else {
			return nil
		}
	}
}

func (r *UserCredentialRepo) DeleteByUserId(ctx context.Context, userId uint32) error {
	builder := r.entClient.Client().UserCredential.Delete()
	builder.Where(usercredential.UserIDEQ(userId))
	if affected, err := builder.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return authenticationV1.ErrorNotFound("user credential not found")
		}

		r.log.Errorf("delete one data failed: %s", err.Error())

		return authenticationV1.ErrorInternalServerError("delete one data failed")
	} else {
		if affected == 0 {
			return authenticationV1.ErrorNotFound("user credential not found")
		} else {
			return nil
		}
	}
}

func (r *UserCredentialRepo) DeleteByIdentifier(ctx context.Context, identityType authenticationV1.UserCredential_IdentityType, identifier string) error {
	builder := r.entClient.Client().UserCredential.Delete()
	builder.Where(
		usercredential.IdentityTypeEQ(*r.identityTypeConverter.ToEntity(&identityType)),
		usercredential.IdentifierEQ(identifier),
	)
	// 租户范围过滤（平台管理员 tenant_id=0 不限定）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.Where(usercredential.TenantIDEQ(tid))
	}
	if affected, err := builder.Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return authenticationV1.ErrorNotFound("user credential not found")
		}

		r.log.Errorf("delete one data failed: %s", err.Error())

		return authenticationV1.ErrorInternalServerError("delete one data failed")
	} else {
		if affected == 0 {
			return authenticationV1.ErrorNotFound("user credential not found")
		} else {
			return nil
		}
	}
}

func (r *UserCredentialRepo) Get(ctx context.Context, req *authenticationV1.GetUserCredentialRequest) (*authenticationV1.UserCredential, error) {
	if req == nil {
		return nil, authenticationV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().UserCredential.Query()

	var whereCond []func(s *sql.Selector)
	switch req.QueryBy.(type) {
	default:
	case *authenticationV1.GetUserCredentialRequest_Id:
		whereCond = append(whereCond, usercredential.IDEQ(req.GetId()))
	}

	dto, err := r.repository.Get(ctx, builder, req.GetViewMask(), whereCond...)
	if err != nil {
		return nil, err
	}

	return dto, err
}

// maybeTenantFromViewer returns the viewer's tenant id if present and >0, plus whether a viewer context exists.
func maybeTenantFromViewer(ctx context.Context) (tenantID uint32, hasTenant bool) {
	vc, exist := viewer.FromContext(ctx)
	if !exist {
		return 0, false
	}
	tid := vc.TenantID()
	if tid == 0 {
		return 0, false
	}
	return uint32(tid), true
}

// viewerUserIDFromContext 从 viewer context 提取操作人的用户 ID。
// 用于强制按调用者身份过滤（如收件箱只能看到自己的消息）。
func viewerUserIDFromContext(ctx context.Context) (uint32, bool) {
	vc, exist := viewer.FromContext(ctx)
	if !exist || vc == nil {
		return 0, false
	}
	uid := vc.UserID()
	if uid == 0 {
		return 0, false
	}
	return uint32(uid), true
}

func (r *UserCredentialRepo) GetByIdentifier(ctx context.Context, req *authenticationV1.GetUserCredentialByIdentifierRequest) (*authenticationV1.UserCredential, error) {
	builder := r.entClient.Client().UserCredential.Query()

	var whereConds []predicate.UserCredential
	whereConds = append(whereConds,
		usercredential.IdentityTypeEQ(*r.identityTypeConverter.ToEntity(trans.Ptr(req.GetIdentityType()))),
		usercredential.IdentifierEQ(req.GetIdentifier()),
	)

	// 租户范围过滤（平台管理员 tenant_id=0 不限定）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		whereConds = append(whereConds, usercredential.TenantIDEQ(tid))
	}

	builder.Where(whereConds...)

	entity, err := builder.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, authenticationV1.ErrorNotFound("user credential not found")
		}

		r.log.Errorf("query one data failed: %s", err.Error())

		return nil, authenticationV1.ErrorInternalServerError("query data failed")
	}

	return r.mapper.ToDTO(entity), nil
}

// dummyPasswordHash 是一个合法的 bcrypt 哈希，用于在用户不存在时执行一次假校验，
// 让"用户不存在"与"密码错误"两条路径耗时一致，抹掉用户名枚举的计时侧信道。
// 明文不可逆，仅作为恒定耗时的目标。
const dummyPasswordHash = "$2a$10$1sbpKmhQDpXLHnDnEQ1nLe3oOnYyP2bUJyqHcX2T0Fq1qfyoXOrPm"

// performDummyVerify 执行一次假的密码校验，用于抹平计时差异。
// 返回值始终为 false，且从不抛错。
// ===== TOTP 两步验证（MFA）因子存取 =====
// 约定：TOTP secret 以 identity_type=USERID + identifier=用户ID + credential_type=TOTP
// 落 sys_user_credentials，secret 经全局加密器（ENCRYPTION_KEY）落库加密；一个用户
// 至多一条启用的 TOTP 因子（重复绑定覆盖旧因子）。

// totpCredentialPredicates TOTP 因子的公共查询谓词。
func totpCredentialPredicates(userID uint32) []predicate.UserCredential {
	return []predicate.UserCredential{
		usercredential.UserIDEQ(userID),
		usercredential.CredentialTypeEQ(usercredential.CredentialTypeTOTP),
		usercredential.StatusEQ(usercredential.StatusEnabled),
	}
}

// GetTOTPCredential 取用户启用的 TOTP secret（解密后返回）；未绑定返回 NotFound。
func (r *UserCredentialRepo) GetTOTPCredential(ctx context.Context, userID uint32) (string, error) {
	entity, err := r.entClient.Client().UserCredential.Query().
		Select(
			usercredential.FieldUserID,
			usercredential.FieldCredential,
		).
		Where(totpCredentialPredicates(userID)...).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", authenticationV1.ErrorNotFound("totp factor not found")
		}
		r.log.Errorf("query totp credential failed: %s", err.Error())
		return "", authenticationV1.ErrorServiceUnavailable("db error")
	}
	if entity.Credential == nil {
		return "", authenticationV1.ErrorNotFound("totp factor not found")
	}
	secret, err := pkgcrypto.DecryptIfNeeded(*entity.Credential)
	if err != nil {
		r.log.Errorf("decrypt totp secret failed: %s", err.Error())
		return "", authenticationV1.ErrorInternalServerError("decrypt totp secret failed")
	}
	return secret, nil
}

// HasTOTPCredential 用户是否已启用 TOTP 因子（登录链热路径，仅探活不取值）。
func (r *UserCredentialRepo) HasTOTPCredential(ctx context.Context, userID uint32) (bool, error) {
	exists, err := r.entClient.Client().UserCredential.Query().
		Where(totpCredentialPredicates(userID)...).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("probe totp credential failed: %s", err.Error())
		return false, authenticationV1.ErrorServiceUnavailable("db error")
	}
	return exists, nil
}

// SaveTOTPCredential 写入（或覆盖）用户的 TOTP 因子，secret 需已加密。
func (r *UserCredentialRepo) SaveTOTPCredential(ctx context.Context, userID, tenantID uint32, encryptedSecret string) error {
	client := r.entClient.Client()
	existing, err := client.UserCredential.Query().
		Select(usercredential.FieldID).
		Where(
			usercredential.UserIDEQ(userID),
			usercredential.CredentialTypeEQ(usercredential.CredentialTypeTOTP),
		).
		First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		r.log.Errorf("probe existing totp credential failed: %s", err.Error())
		return authenticationV1.ErrorServiceUnavailable("db error")
	}
	identityType := r.identityTypeConverter.ToEntity(trans.Ptr(authenticationV1.UserCredential_USERID))
	status := usercredential.StatusEnabled
	now := time.Now()
	if existing != nil {
		if _, uErr := client.UserCredential.UpdateOneID(existing.ID).
			SetCredential(encryptedSecret).
			SetStatus(status).
			SetNillableUpdatedAt(&now).
			Save(ctx); uErr != nil {
			r.log.Errorf("update totp credential failed: %s", uErr.Error())
			return authenticationV1.ErrorInternalServerError("save totp factor failed")
		}
		return nil
	}
	identifier := strconv.FormatUint(uint64(userID), 10)
	if _, cErr := client.UserCredential.Create().
		SetUserID(userID).
		SetTenantID(tenantID).
		SetNillableIdentityType(identityType).
		SetIdentifier(identifier).
		SetCredentialType(usercredential.CredentialTypeTOTP).
		SetCredential(encryptedSecret).
		SetStatus(status).
		SetNillableCreatedAt(&now).
		SetNillableUpdatedAt(&now).
		Save(ctx); cErr != nil {
		r.log.Errorf("create totp credential failed: %s", cErr.Error())
		return authenticationV1.ErrorInternalServerError("save totp factor failed")
	}
	return nil
}

// DeleteTOTPCredentials 清空用户全部 TOTP 因子（管理员救援重置）。
func (r *UserCredentialRepo) DeleteTOTPCredentials(ctx context.Context, userID uint32) (int, error) {
	affected, err := r.entClient.Client().UserCredential.Delete().
		Where(
			usercredential.UserIDEQ(userID),
			usercredential.CredentialTypeEQ(usercredential.CredentialTypeTOTP),
		).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("delete totp credentials failed: %s", err.Error())
		return 0, authenticationV1.ErrorInternalServerError("delete totp factors failed")
	}
	return affected, nil
}

func (r *UserCredentialRepo) performDummyVerify(plainCredential string) {
	_, _ = r.passwordCrypto.Verify(plainCredential, dummyPasswordHash)
}

// FindUserCredential 在指定租户范围内根据身份类型+标识符查询单条凭证并校验密码。
// tenantID 由调用方解析确定（登录时按 tenant_code 解析；留空视为平台 tenant 0）。
// (tenant_id, identity_type, identifier) 在该 tenant 内唯一，无需跨租户消歧。
// needDecrypt 为 true 时，plainCredential 应为 base64+AES 加密密文，会自动解密。
// 返回匹配到的用户 ID，调用方可据此通过 ID 查找用户。
func (r *UserCredentialRepo) FindUserCredential(ctx context.Context, tenantID uint32, identityType authenticationV1.UserCredential_IdentityType, identifier, plainCredential string, needDecrypt bool) (uint32, error) {
	if needDecrypt {
		bytesPass, err := base64.StdEncoding.DecodeString(plainCredential)
		if err != nil {
			r.log.Errorf("decode base64 credential failed: %s", err.Error())
			return 0, authenticationV1.ErrorBadRequest("invalid credential format")
		}
		decrypted, err := crypto.AesDecrypt(bytesPass, crypto.DefaultAESKey, nil)
		if err != nil {
			r.log.Errorf("decrypt credential failed: %s", err.Error())
			return 0, authenticationV1.ErrorBadRequest("decrypt credential failed")
		}
		plainCredential = string(decrypted)
	}

	entity, err := r.entClient.Client().UserCredential.Query().
		Select(
			usercredential.FieldUserID,
			usercredential.FieldCredentialType,
			usercredential.FieldCredential,
			usercredential.FieldStatus,
			usercredential.FieldUpdatedAt,
		).
		Where(
			usercredential.TenantIDEQ(tenantID),
			usercredential.IdentityTypeEQ(*r.identityTypeConverter.ToEntity(trans.Ptr(identityType))),
			usercredential.IdentifierEQ(identifier),
		).
		Only(ctx)
	if err != nil {
		// 恒定时间防护：用户不存在时也跑一次 bcrypt 校验，避免被计时攻击枚举用户名
		r.performDummyVerify(plainCredential)
		if ent.IsNotFound(err) {
			return 0, authenticationV1.ErrorUserNotFound("user not found")
		}
		r.log.Errorf("query credential failed: %s", err.Error())
		return 0, authenticationV1.ErrorServiceUnavailable("db error")
	}

	if entity.Status == nil || entity.CredentialType == nil || entity.Credential == nil || entity.UserID == nil {
		r.performDummyVerify(plainCredential)
		return 0, authenticationV1.ErrorUserNotFound("user not found")
	}
	if *entity.Status != usercredential.StatusEnabled {
		r.performDummyVerify(plainCredential)
		return 0, authenticationV1.ErrorUserNotFound("user not found")
	}

	if r.verifyCredential(entity.CredentialType, plainCredential, *entity.Credential) {
		// 等保口令有效期：超期拒绝登录，用户走重置/改密流程换新口令
		if expired, xErr := r.checkPasswordExpired(ctx, entity); xErr == nil && expired {
			r.log.Warnf("password expired for user [%d], login denied", *entity.UserID)
			return 0, authenticationV1.ErrorBadRequest("password expired, please reset your password")
		}
		return *entity.UserID, nil
	}

	return 0, authenticationV1.ErrorInvalidPassword("incorrect password")
}

// FindUserCredentialAcrossTenants 在全部租户范围内按身份类型+标识符定位凭证并校验密码。
// 仅当标识符在所有租户中唯一命中一条启用凭证时才成功，返回 (userID, credentialTenantID)；
// 命中 0 条返回 UserNotFound，命中多条（同名歧义）时要求调用方改传 tenant_code 精确登录。
// 供未携带租户上下文的 C 端登录回退使用，避免租户用户因无 tenant_code 而永远无法登录。
func (r *UserCredentialRepo) FindUserCredentialAcrossTenants(ctx context.Context, identityType authenticationV1.UserCredential_IdentityType, identifier, plainCredential string, needDecrypt bool) (uint32, uint32, error) {
	if needDecrypt {
		bytesPass, err := base64.StdEncoding.DecodeString(plainCredential)
		if err != nil {
			r.log.Errorf("decode base64 credential failed: %s", err.Error())
			return 0, 0, authenticationV1.ErrorBadRequest("invalid credential format")
		}
		decrypted, err := crypto.AesDecrypt(bytesPass, crypto.DefaultAESKey, nil)
		if err != nil {
			r.log.Errorf("decrypt credential failed: %s", err.Error())
			return 0, 0, authenticationV1.ErrorBadRequest("decrypt credential failed")
		}
		plainCredential = string(decrypted)
	}

	entities, err := r.entClient.Client().UserCredential.Query().
		Select(
			usercredential.FieldTenantID,
			usercredential.FieldUserID,
			usercredential.FieldCredentialType,
			usercredential.FieldCredential,
			usercredential.FieldStatus,
			usercredential.FieldUpdatedAt,
		).
		Where(
			usercredential.IdentityTypeEQ(*r.identityTypeConverter.ToEntity(trans.Ptr(identityType))),
			usercredential.IdentifierEQ(identifier),
		).
		Limit(2).
		All(ctx)
	if err != nil {
		r.performDummyVerify(plainCredential)
		r.log.Errorf("query credential across tenants failed: %s", err.Error())
		return 0, 0, authenticationV1.ErrorServiceUnavailable("db error")
	}

	var matched *ent.UserCredential
	ambiguous := false
	for i := range entities {
		e := entities[i]
		if e.Status == nil || *e.Status != usercredential.StatusEnabled ||
			e.CredentialType == nil || e.Credential == nil || e.UserID == nil || e.TenantID == nil {
			continue
		}
		if matched != nil {
			ambiguous = true
			break
		}
		matched = e
	}
	if ambiguous {
		// 同名标识符跨租户歧义：拒绝并要求 tenant_code，不做密码校验
		r.performDummyVerify(plainCredential)
		return 0, 0, authenticationV1.ErrorUserNotFound("ambiguous identifier, tenant code required")
	}
	if matched == nil {
		r.performDummyVerify(plainCredential)
		return 0, 0, authenticationV1.ErrorUserNotFound("user not found")
	}

	if r.verifyCredential(matched.CredentialType, plainCredential, *matched.Credential) {
		// 等保口令有效期（同 FindUserCredential）
		if expired, xErr := r.checkPasswordExpired(ctx, matched); xErr == nil && expired {
			r.log.Warnf("password expired for user [%d] (cross-tenant), login denied", *matched.UserID)
			return 0, 0, authenticationV1.ErrorBadRequest("password expired, please reset your password")
		}
		return *matched.UserID, *matched.TenantID, nil
	}
	return 0, 0, authenticationV1.ErrorInvalidPassword("incorrect password")
}

func (r *UserCredentialRepo) VerifyCredential(ctx context.Context, req *authenticationV1.VerifyCredentialRequest) (*authenticationV1.VerifyCredentialResponse, error) {
	// 该内部 RPC 请求未携带租户信息，按平台（tenant 0）范围校验。
	if _, err := r.FindUserCredential(ctx, 0, req.GetIdentityType(), req.GetIdentifier(), req.GetCredential(), req.GetNeedDecrypt()); err != nil {
		return nil, err
	}

	return &authenticationV1.VerifyCredentialResponse{
		Success: true,
	}, nil
}

func (r *UserCredentialRepo) verifyCredential(credentialType *usercredential.CredentialType, plainCredential, targetCredential string) bool {
	if credentialType == nil || plainCredential == "" {
		return false
	}

	switch *credentialType {
	case usercredential.CredentialTypePasswordHash:
		ok, err := r.passwordCrypto.Verify(plainCredential, targetCredential)
		if err != nil {
			r.log.Errorf("verify password failed: %s", err.Error())
			return false
		}
		return ok
	default:
		return plainCredential == targetCredential
	}
}

func (r *UserCredentialRepo) prepareCredential(credentialType *usercredential.CredentialType, plainCredential string) (string, error) {
	var newCredential string
	switch *credentialType {
	case usercredential.CredentialTypePasswordHash:
		var err error
		// 加密明文密码
		newCredential, err = r.passwordCrypto.Encrypt(plainCredential)
		if err != nil {
			r.log.Errorf("hash new password failed: %s", err.Error())
			return "", authenticationV1.ErrorBadRequest("hash new password failed")
		}

	default:
		newCredential = plainCredential
	}

	return newCredential, nil
}

// ChangeCredential 修改认证信息
func (r *UserCredentialRepo) ChangeCredential(ctx context.Context, req *authenticationV1.ChangeCredentialRequest) error {
	if req.GetNeedDecrypt() {
		// 解密旧密码
		bytesPass, err := base64.StdEncoding.DecodeString(req.GetOldCredential())
		if err != nil {
			return authenticationV1.ErrorBadRequest("invalid old credential format")
		}
		plainPassword, err := crypto.AesDecrypt(bytesPass, crypto.DefaultAESKey, nil)
		if err != nil {
			return authenticationV1.ErrorBadRequest("decrypt old credential failed")
		}
		req.OldCredential = string(plainPassword)

		// 解密新密码
		bytesPass, err = base64.StdEncoding.DecodeString(req.GetNewCredential())
		if err != nil {
			return authenticationV1.ErrorBadRequest("invalid new credential format")
		}
		plainPassword, err = crypto.AesDecrypt(bytesPass, crypto.DefaultAESKey, nil)
		if err != nil {
			return authenticationV1.ErrorBadRequest("decrypt new credential failed")
		}
		req.NewCredential = string(plainPassword)
	}

	// 租户范围 where 条件（平台管理员 tenant_id=0 不限定）
	tenantWhere := []predicate.UserCredential{
		usercredential.IdentityTypeEQ(*r.identityTypeConverter.ToEntity(trans.Ptr(req.GetIdentityType()))),
		usercredential.IdentifierEQ(req.GetIdentifier()),
	}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		tenantWhere = append(tenantWhere, usercredential.TenantIDEQ(tid))
	}

	entity, err := r.entClient.Client().UserCredential.
		Query().
		Select(
			usercredential.FieldCredentialType,
			usercredential.FieldCredential,
			usercredential.FieldExtraInfo,
		).
		Where(tenantWhere...).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return authenticationV1.ErrorNotFound("user credential not found")
		}
		r.log.Errorf("query one data failed: %s", err.Error())
		return authenticationV1.ErrorInternalServerError("query one data failed")
	}

	if entity.CredentialType == nil {
		return authenticationV1.ErrorNotFound("user credential not found")
	}

	// 验证旧认证信息
	if !r.verifyCredential(entity.CredentialType, req.GetOldCredential(), *entity.Credential) {
		return authenticationV1.ErrorBadRequest("invalid old password")
	}

	// 等保口令复杂度：新明文口令须长度达标且四类取三
	if err = passwordpolicy.ValidateComplexity(req.GetNewCredential(), r.passwordMinLen(ctx)); err != nil {
		return authenticationV1.ErrorBadRequest("%s", err.Error())
	}
	// 等保口令历史：新口令不得与最近 N 条历史口令相同
	if err = r.checkPasswordHistory(ctx, entity, req.GetNewCredential()); err != nil {
		return err
	}

	// 注意：这里必须用 GetNewCredential，原代码误写成 GetOldCredential 会导致
	// 改密码后库里哈希被重置成旧密码的哈希，新密码无法登录、旧密码永久有效。
	var newCredential string
	newCredential, err = r.prepareCredential(entity.CredentialType, req.GetNewCredential())
	if err != nil {
		r.log.Errorf("prepare new credential failed: %s", err.Error())
		return authenticationV1.ErrorBadRequest("prepare new credential failed")
	}

	if newCredential == "" {
		return authenticationV1.ErrorBadRequest("new credential cannot be empty")
	}

	// 等保口令历史：旧哈希追加进 extra_info（只保留最近 N 条）
	histLimit := r.passwordHistoryLimit(ctx)
	extraInfo := appendPasswordHistory(entity.ExtraInfo, *entity.Credential, histLimit)

	builder := r.entClient.Client().UserCredential.Update()
	builder.Where(tenantWhere...)
	builder.
		SetCredential(newCredential).
		SetNillableExtraInfo(extraInfo).
		SetUpdatedAt(time.Now())
	if err = builder.Exec(ctx); err != nil {
		r.log.Errorf("update one data failed: %s", err.Error())
		return authenticationV1.ErrorInternalServerError("update data failed")
	}

	return nil
}

// ResetCredential 修改认证信息
func (r *UserCredentialRepo) ResetCredential(ctx context.Context, req *authenticationV1.ResetCredentialRequest) error {
	if req.GetNeedDecrypt() {
		// 解密新密码
		bytesPass, err := base64.StdEncoding.DecodeString(req.GetNewCredential())
		if err != nil {
			return authenticationV1.ErrorBadRequest("invalid new credential format")
		}
		plainPassword, err := crypto.AesDecrypt(bytesPass, crypto.DefaultAESKey, nil)
		if err != nil {
			return authenticationV1.ErrorBadRequest("decrypt new credential failed")
		}
		req.NewCredential = string(plainPassword)
	}

	// 租户范围 where 条件（平台管理员 tenant_id=0 不限定）
	tenantWhere := []predicate.UserCredential{
		usercredential.IdentityTypeEQ(*r.identityTypeConverter.ToEntity(trans.Ptr(req.GetIdentityType()))),
		usercredential.IdentifierEQ(req.GetIdentifier()),
	}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		tenantWhere = append(tenantWhere, usercredential.TenantIDEQ(tid))
	}

	entity, err := r.entClient.Client().UserCredential.
		Query().
		Select(
			usercredential.FieldCredentialType,
			usercredential.FieldCredential,
			usercredential.FieldExtraInfo,
		).
		Where(tenantWhere...).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return authenticationV1.ErrorNotFound("user credential not found")
		}
		r.log.Errorf("query one data failed: %s", err.Error())
		return authenticationV1.ErrorInternalServerError("query one data failed")
	}

	if entity.CredentialType == nil {
		return authenticationV1.ErrorNotFound("user credential not found")
	}

	// 等保口令复杂度：重置码重置同样把关（ForgotPassword 链路的 NewCredential
	// 在 service 层已解密为明文）
	if err = passwordpolicy.ValidateComplexity(req.GetNewCredential(), r.passwordMinLen(ctx)); err != nil {
		return authenticationV1.ErrorBadRequest("%s", err.Error())
	}
	// 等保口令历史：新口令不得与最近 N 条相同
	if err = r.checkPasswordHistory(ctx, entity, req.GetNewCredential()); err != nil {
		return err
	}

	var newCredential string
	newCredential, err = r.prepareCredential(entity.CredentialType, req.GetNewCredential())
	if err != nil {
		r.log.Errorf("prepare new credential failed: %s", err.Error())
		return authenticationV1.ErrorBadRequest("prepare new credential failed")
	}

	if newCredential == "" {
		return authenticationV1.ErrorBadRequest("new credential cannot be empty")
	}

	// 等保口令历史：旧哈希追加（重置前的口令同样进入历史，防止绕回旧口令）
	histLimit := r.passwordHistoryLimit(ctx)
	extraInfo := appendPasswordHistory(entity.ExtraInfo, *entity.Credential, histLimit)

	builder := r.entClient.Client().UserCredential.Update()
	builder.Where(tenantWhere...)
	builder.
		SetCredential(newCredential).
		SetNillableExtraInfo(extraInfo).
		SetUpdatedAt(time.Now())
	if err = builder.Exec(ctx); err != nil {
		r.log.Errorf("update one data failed: %s", err.Error())
		return authenticationV1.ErrorInternalServerError("update data failed")
	}

	return nil
}

// ── 等保口令策略（阈值自 sys_config 平台参数，参数管理页可调；<=0 关闭对应检查）──
//
// 三条规则的落点：
//   - 复杂度（minLen + 四类取三）：改密/重置入口 prepareCredential 之前；
//   - 历史（最近 N 条 bcrypt 比对，存 extra_info.password_history）：同上，改密成功后追加；
//   - 有效期（updated_at 起算）：登录凭证校验成功后拒绝超期登录。

func (r *UserCredentialRepo) passwordMinLen(ctx context.Context) int {
	if r.configRepo == nil {
		return constants.DefaultPasswordMinLen
	}
	return r.configRepo.GetConfigInt(ctx, constants.ConfigKeyPasswordMinLen, constants.DefaultPasswordMinLen)
}

func (r *UserCredentialRepo) passwordMaxAgeDays(ctx context.Context) int {
	if r.configRepo == nil {
		return constants.DefaultPasswordMaxAgeDays
	}
	return r.configRepo.GetConfigInt(ctx, constants.ConfigKeyPasswordMaxAgeDays, constants.DefaultPasswordMaxAgeDays)
}

func (r *UserCredentialRepo) passwordHistoryLimit(ctx context.Context) int {
	if r.configRepo == nil {
		return constants.DefaultPasswordHistoryCount
	}
	return r.configRepo.GetConfigInt(ctx, constants.ConfigKeyPasswordHistoryCount, constants.DefaultPasswordHistoryCount)
}

// passwordHistoryKey 是 extra_info JSON 里历史口令哈希数组的键。
const passwordHistoryKey = "password_history"

// checkPasswordExpired 等保口令有效期：密码哈希行 updated_at 距今超过 maxAge 天
// 即超期（登录被拒，走重置/改密流程换新口令）。updated_at 为空视为不超期。
func (r *UserCredentialRepo) checkPasswordExpired(ctx context.Context, entity *ent.UserCredential) (bool, error) {
	maxAge := r.passwordMaxAgeDays(ctx)
	if maxAge <= 0 || entity == nil || entity.CredentialType == nil {
		return false, nil
	}
	if *entity.CredentialType != usercredential.CredentialTypePasswordHash {
		return false, nil
	}
	if entity.UpdatedAt == nil {
		return false, nil
	}
	return time.Since(*entity.UpdatedAt) > time.Duration(maxAge)*24*time.Hour, nil
}

// checkPasswordHistory 等保口令策略：新明文口令不得与历史哈希列表中任一条
// 相同（bcrypt 比对）。历史列表为空/解析失败时跳过（不阻塞主流程）。
func (r *UserCredentialRepo) checkPasswordHistory(ctx context.Context, entity *ent.UserCredential, newPlain string) error {
	limit := r.passwordHistoryLimit(ctx)
	if limit <= 0 || entity == nil || entity.CredentialType == nil {
		return nil
	}
	if *entity.CredentialType != usercredential.CredentialTypePasswordHash {
		return nil
	}
	var extra map[string]any
	if entity.ExtraInfo != nil && *entity.ExtraInfo != "" {
		if err := json.Unmarshal([]byte(*entity.ExtraInfo), &extra); err != nil {
			// extra_info 非历史格式（可能存了别的业务数据），跳过历史检查
			return nil
		}
	}
	arr, _ := extra[passwordHistoryKey].([]any)
	for _, item := range arr {
		hash, ok := item.(string)
		if !ok || hash == "" {
			continue
		}
		if ok, err := r.passwordCrypto.Verify(newPlain, hash); err == nil && ok {
			return authenticationV1.ErrorBadRequest("new password must differ from recent passwords")
		}
	}
	return nil
}

// appendPasswordHistory 把旧口令哈希追加进 extra_info 的历史列表，
// 只保留最近 limit 条；extra_info 里已有的其他键原样保留。
func appendPasswordHistory(extraInfo *string, oldHash string, limit int) *string {
	if limit <= 0 {
		return extraInfo
	}
	var extra map[string]any
	if extraInfo != nil && *extraInfo != "" {
		if err := json.Unmarshal([]byte(*extraInfo), &extra); err != nil || extra == nil {
			// 解析失败：历史让位，以新结构重建
			extra = make(map[string]any)
		}
	} else {
		extra = make(map[string]any)
	}

	arr, _ := extra[passwordHistoryKey].([]any)
	history := make([]string, 0, limit)
	for _, item := range arr {
		if s, ok := item.(string); ok && s != "" {
			history = append(history, s)
		}
	}
	if oldHash != "" {
		history = append(history, oldHash)
	}
	if len(history) > limit {
		history = history[len(history)-limit:]
	}

	strs := make([]any, 0, len(history))
	for _, h := range history {
		strs = append(strs, h)
	}
	extra[passwordHistoryKey] = strs

	b, err := json.Marshal(extra)
	if err != nil {
		return extraInfo
	}
	s := string(b)
	return &s
}
