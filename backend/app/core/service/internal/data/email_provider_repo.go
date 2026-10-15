package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	entCrud "github.com/tx7do/go-crud/entgo"

	pkgcrypto "go-wind-cms/pkg/crypto"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/emailprovider"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// EmailProviderRepo 邮件出站 SMTP 配置仓储。
// 密码 AES-GCM 加密落库（ENCRYPTION_KEY 全局 Encryptor，与社媒凭据同源）；
// 对外 DTO 永不携带密码（写专用，读取置空），worker 经 GetActive 拿解密后的配置。
type EmailProviderRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewEmailProviderRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *EmailProviderRepo {
	return &EmailProviderRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("email-provider/repo/core-service"),
	}
}

// DecryptedConfig 解密后的 SMTP 连接配置（worker/测试发送专用，不出服务层）。
type DecryptedConfig struct {
	Host     string
	Port     uint32
	Username string
	Password string
	From     string
	TLSMode  string
	Enabled  bool
	Name     string
}

// List 配置列表（密码恒为空）。
func (r *EmailProviderRepo) List(ctx context.Context) ([]*notifyV1.EmailProvider, int, error) {
	conds := r.tenantConds(ctx)
	rows, err := r.entClient.Client().EmailProvider.Query().Where(conds...).
		Order(ent.Desc(emailprovider.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("list email providers failed: %s", err.Error())
		return nil, 0, err
	}
	items := make([]*notifyV1.EmailProvider, 0, len(rows))
	for _, row := range rows {
		items = append(items, r.toDTO(row))
	}
	return items, len(items), nil
}

// Get 配置详情（密码恒为空）。
func (r *EmailProviderRepo) Get(ctx context.Context, id uint32) (*notifyV1.EmailProvider, error) {
	conds := append(r.tenantConds(ctx), emailprovider.IDEQ(id))
	row, err := r.entClient.Client().EmailProvider.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("EMAIL_PROVIDER", "email provider not found")
		}
		r.log.Errorf("get email provider failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// GetActive 第一个启用项（解密密码；worker 专用，无租户过滤——平台级基础设施配置）。
func (r *EmailProviderRepo) GetActive(ctx context.Context) (*DecryptedConfig, error) {
	row, err := r.entClient.Client().EmailProvider.Query().
		Where(emailprovider.EnabledEQ(true)).
		Order(ent.Desc(emailprovider.FieldID)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	password := ""
	if row.Password != nil && *row.Password != "" {
		plain, dErr := pkgcrypto.DecryptIfNeeded(*row.Password)
		if dErr != nil {
			r.log.Errorf("decrypt email provider password failed (id=%d): %s", row.ID, dErr.Error())
			return nil, dErr
		}
		password = plain
	}
	name := row.Name
	tlsMode := "starttls"
	if row.TLSMode != nil {
		tlsMode = *row.TLSMode
	}
	port := uint32(587)
	if row.Port != nil {
		port = *row.Port
	}
	username := ""
	if row.Username != nil {
		username = *row.Username
	}
	return &DecryptedConfig{
		Host:     row.Host,
		Port:     port,
		Username: username,
		Password: password,
		From:     row.From,
		TLSMode:  tlsMode,
		Enabled:  row.Enabled != nil && *row.Enabled,
		Name:     name,
	}, nil
}

// Create 新建（密码加密落库）。
func (r *EmailProviderRepo) Create(ctx context.Context, dto *notifyV1.EmailProvider) (*notifyV1.EmailProvider, error) {
	builder := r.entClient.Client().EmailProvider.Create().
		SetName(dto.GetName()).
		SetHost(dto.GetHost()).
		SetFrom(dto.GetFrom()).
		SetEnabled(dto.GetEnabled())
	if dto.GetPort() > 0 {
		builder.SetPort(dto.GetPort())
	}
	if dto.GetUsername() != "" {
		builder.SetUsername(dto.GetUsername())
	}
	if dto.GetPassword() != "" {
		enc, encErr := pkgcrypto.EncryptIfNeeded(dto.GetPassword())
		if encErr != nil {
			r.log.Errorf("encrypt email provider password failed: %s", encErr.Error())
			return nil, encErr
		}
		builder.SetPassword(enc)
	}
	if dto.GetTlsMode() != "" {
		builder.SetTLSMode(dto.GetTlsMode())
	}
	if dto.GetDescription() != "" {
		builder.SetDescription(dto.GetDescription())
	}
	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetCreatedBy(callerUserID).SetUpdatedBy(callerUserID)
	}
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		builder.SetTenantID(tid)
	}

	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("create email provider failed: %s", err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Update 按 mask 更新（密码仅在传入非空时覆写）。
func (r *EmailProviderRepo) Update(ctx context.Context, id uint32, dto *notifyV1.EmailProvider) (*notifyV1.EmailProvider, error) {
	conds := append(r.tenantConds(ctx), emailprovider.IDEQ(id))
	exist, err := r.entClient.Client().EmailProvider.Query().Where(conds...).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("EMAIL_PROVIDER", "email provider not found")
		}
		return nil, err
	}

	builder := r.entClient.Client().EmailProvider.UpdateOneID(exist.ID)
	if dto.Name != nil {
		builder.SetName(dto.GetName())
	}
	if dto.Host != nil {
		builder.SetHost(dto.GetHost())
	}
	if dto.Port != nil {
		builder.SetPort(dto.GetPort())
	}
	if dto.Username != nil {
		builder.SetUsername(dto.GetUsername())
	}
	if dto.Password != nil && dto.GetPassword() != "" {
		enc, encErr := pkgcrypto.EncryptIfNeeded(dto.GetPassword())
		if encErr != nil {
			r.log.Errorf("encrypt email provider password failed (id=%d): %s", id, encErr.Error())
			return nil, encErr
		}
		builder.SetPassword(enc)
	}
	if dto.From != nil {
		builder.SetFrom(dto.GetFrom())
	}
	if dto.TlsMode != nil {
		builder.SetTLSMode(dto.GetTlsMode())
	}
	if dto.Enabled != nil {
		builder.SetEnabled(dto.GetEnabled())
	}
	if dto.Description != nil {
		if dto.GetDescription() == "" {
			builder.ClearDescription()
		} else {
			builder.SetDescription(dto.GetDescription())
		}
	}
	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetUpdatedBy(callerUserID)
	}

	row, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("update email provider failed (id=%d): %s", id, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Delete 删除。
func (r *EmailProviderRepo) Delete(ctx context.Context, id uint32) error {
	conds := append(r.tenantConds(ctx), emailprovider.IDEQ(id))
	count, err := r.entClient.Client().EmailProvider.Delete().Where(conds...).Exec(ctx)
	if err != nil {
		r.log.Errorf("delete email provider failed (id=%d): %s", id, err.Error())
		return err
	}
	if count == 0 {
		return errors.NotFound("EMAIL_PROVIDER", "email provider not found")
	}
	return nil
}

func (r *EmailProviderRepo) tenantConds(ctx context.Context) []predicate.EmailProvider {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant {
		return nil
	}
	return []predicate.EmailProvider{emailprovider.TenantIDEQ(tid)}
}

// toDTO 密码恒为空（写专用，接口永不回显）。
func (r *EmailProviderRepo) toDTO(row *ent.EmailProvider) *notifyV1.EmailProvider {
	enabled := true
	if row.Enabled != nil {
		enabled = *row.Enabled
	}
	port := uint32(587)
	if row.Port != nil {
		port = *row.Port
	}
	dto := &notifyV1.EmailProvider{
		Id:          &row.ID,
		Name:        &row.Name,
		Host:        &row.Host,
		Port:        &port,
		Username:    row.Username,
		Password:    nil,
		From:        &row.From,
		TlsMode:     row.TLSMode,
		Enabled:     &enabled,
		Description: row.Description,
		CreatedBy:   row.CreatedBy,
		UpdatedBy:   row.UpdatedBy,
	}
	if row.CreatedAt != nil {
		dto.CreatedAt = timestamppb.New(*row.CreatedAt)
	}
	if row.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*row.UpdatedAt)
	}
	return dto
}
