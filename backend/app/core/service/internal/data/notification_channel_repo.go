package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	pkgcrypto "go-wind-cms/pkg/crypto"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/notificationchannel"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	notifyV1 "go-wind-cms/api/gen/go/notify/service/v1"
)

// NotificationChannelRepo 通知渠道仓储。
// 凭据（SMTP 密码 / Webhook 签名密钥）AES-GCM 加密落库（ENCRYPTION_KEY 同源），
// 对外 DTO 永不携带凭据（仅 Has* 布尔标识），sender 经 GetDecrypted*/GetFirstEnabled* 拿解密配置。
type NotificationChannelRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper        *mapper.CopierMapper[notifyV1.NotificationChannel, ent.NotificationChannel]
	typeConverter *mapper.EnumTypeConverter[notifyV1.NotificationChannel_Type, notificationchannel.Type]
	tlsConverter  *mapper.EnumTypeConverter[notifyV1.NotificationChannel_TlsMode, notificationchannel.SMTPTLS]
	// styleConverter 按枚举名字字符串配对，所以 proto 的成员名必须与 ent 列值逐字相同
	//（CUSTOM / NONE / DINGTALK / FEISHU / WECOM）。给任何一侧加前缀都不会报错，只会让
	// 这一列在读路径上静默变成 nil——proto 侧的注释里钉着同一条约束。
	styleConverter *mapper.EnumTypeConverter[notifyV1.SignStyle, notificationchannel.WebhookSignStyle]
	// smsProviderConverter 同按名字配对（ALIYUN / TENCENT）
	smsProviderConverter *mapper.EnumTypeConverter[notifyV1.NotificationChannel_SmsProvider, notificationchannel.SmsProvider]

	repository *entCrud.Repository[
		ent.NotificationChannelQuery, ent.NotificationChannelSelect,
		ent.NotificationChannelCreate, ent.NotificationChannelCreateBulk,
		ent.NotificationChannelUpdate, ent.NotificationChannelUpdateOne,
		ent.NotificationChannelDelete,
		predicate.NotificationChannel,
		notifyV1.NotificationChannel, ent.NotificationChannel,
	]
}

func NewNotificationChannelRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *NotificationChannelRepo {
	repo := &NotificationChannelRepo{
		log:       ctx.NewLoggerHelper("notification-channel/repo/core-service"),
		entClient: entClient,
		mapper:    mapper.NewCopierMapper[notifyV1.NotificationChannel, ent.NotificationChannel](),
		typeConverter: mapper.NewEnumTypeConverter[notifyV1.NotificationChannel_Type, notificationchannel.Type](
			notifyV1.NotificationChannel_Type_name, notifyV1.NotificationChannel_Type_value,
		),
		tlsConverter: mapper.NewEnumTypeConverter[notifyV1.NotificationChannel_TlsMode, notificationchannel.SMTPTLS](
			notifyV1.NotificationChannel_TlsMode_name, notifyV1.NotificationChannel_TlsMode_value,
		),
		styleConverter: mapper.NewEnumTypeConverter[notifyV1.SignStyle, notificationchannel.WebhookSignStyle](
			notifyV1.SignStyle_name, notifyV1.SignStyle_value,
		),
		smsProviderConverter: mapper.NewEnumTypeConverter[notifyV1.NotificationChannel_SmsProvider, notificationchannel.SmsProvider](
			notifyV1.NotificationChannel_SmsProvider_name, notifyV1.NotificationChannel_SmsProvider_value,
		),
	}

	repo.repository = entCrud.NewRepository[
		ent.NotificationChannelQuery, ent.NotificationChannelSelect,
		ent.NotificationChannelCreate, ent.NotificationChannelCreateBulk,
		ent.NotificationChannelUpdate, ent.NotificationChannelUpdateOne,
		ent.NotificationChannelDelete,
		predicate.NotificationChannel,
		notifyV1.NotificationChannel, ent.NotificationChannel,
	](repo.mapper)

	repo.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	repo.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	repo.mapper.AppendConverters(repo.typeConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.tlsConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.styleConverter.NewConverterPair())
	repo.mapper.AppendConverters(repo.smsProviderConverter.NewConverterPair())

	return repo
}

// List 分页查询通知渠道（凭据不出现在 DTO，靠 Has* 标识）。
func (r *NotificationChannelRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*notifyV1.ListNotificationChannelResponse, error) {
	if req == nil {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "invalid parameter")
	}

	builder := r.entClient.Client().NotificationChannel.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		r.log.Errorf("list notification channels failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "list notification channels failed")
	}
	if ret == nil {
		return &notifyV1.ListNotificationChannelResponse{Total: 0, Items: nil}, nil
	}

	// 凭据有无标识与启用态回填（Sensitive 列与 status 枚举列不在 mapper 的可拷贝面里）
	r.fillSecretFlagsByIDs(ctx, ret.Items)
	r.fillEnabledByIDs(ctx, ret.Items)

	return &notifyV1.ListNotificationChannelResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

// fillSecretFlagsByIDs 为 DTO 列表填充"是否已配置凭据"标识（SMTP 密码 / Webhook 密钥）。
// 两列都是 Sensitive 且 DTO 根本没有对应字段（只有 has* 布尔），读路径不该把密文带出去。
func (r *NotificationChannelRepo) fillSecretFlagsByIDs(ctx context.Context, items []*notifyV1.NotificationChannel) {
	if len(items) == 0 {
		return
	}
	entities, err := r.entClient.Client().NotificationChannel.Query().
		Select(notificationchannel.FieldID, notificationchannel.FieldSMTPPassword, notificationchannel.FieldWebhookSecret).
		All(ctx)
	if err != nil {
		r.log.Errorf("query credential flags failed: %s", err.Error())
		return
	}
	hasSmtpPwd := make(map[uint32]bool, len(entities))
	hasWebhookSecret := make(map[uint32]bool, len(entities))
	for _, e := range entities {
		hasSmtpPwd[e.ID] = e.SMTPPassword != nil
		hasWebhookSecret[e.ID] = e.WebhookSecret != nil
	}
	for _, it := range items {
		it.HasPassword = trans.Ptr(hasSmtpPwd[it.GetId()])
		it.HasWebhookSecret = trans.Ptr(hasWebhookSecret[it.GetId()])
	}
}

// fillEnabledByIDs 为 DTO 列表回填 Enabled 标识（SwitchStatus 的 status 枚举列 → 布尔，
// copier 按字段名映射 Status≠Enabled，读视图不回填则恒呈停用态）。
func (r *NotificationChannelRepo) fillEnabledByIDs(ctx context.Context, items []*notifyV1.NotificationChannel) {
	if len(items) == 0 {
		return
	}
	entities, err := r.entClient.Client().NotificationChannel.Query().
		Select(notificationchannel.FieldID, notificationchannel.FieldStatus).
		All(ctx)
	if err != nil {
		r.log.Errorf("query channel statuses failed: %s", err.Error())
		return
	}
	statuses := make(map[uint32]notificationchannel.Status, len(entities))
	for _, e := range entities {
		if e.Status != nil {
			statuses[e.ID] = *e.Status
		}
	}
	for _, it := range items {
		s, ok := statuses[it.GetId()]
		it.Enabled = trans.Ptr(ok && s == notificationchannel.StatusOn)
	}
}

// Get 按 ID 查询渠道（脱敏：DTO 不含凭据，仅 Has* 标识）。
func (r *NotificationChannelRepo) Get(ctx context.Context, id uint32) (*notifyV1.NotificationChannel, error) {
	entity, err := r.entClient.Client().NotificationChannel.Get(ctx, id)
	if err != nil {
		r.log.Errorf("get notification channel [%d] failed: %s", id, err.Error())
		return nil, errors.NotFound("NOTIFICATION_CHANNEL", "notification channel not found")
	}
	dto := r.mapper.ToDTO(entity)
	dto.HasPassword = trans.Ptr(entity.SMTPPassword != nil)
	dto.HasWebhookSecret = trans.Ptr(entity.WebhookSecret != nil)
	dto.HasSmsAccessSecret = trans.Ptr(entity.SmsAccessSecret != nil)
	// Enabled 回填：写路径把 enabled 布尔落为 status 枚举列（statusFromEnabled），
	// copier 按字段名映射 Status≠Enabled，读视图不回填则恒呈停用态。
	dto.Enabled = trans.Ptr(entity.Status != nil && *entity.Status == notificationchannel.StatusOn)
	return dto, nil
}

// Create 创建渠道；凭据明文经 EncryptIfNeeded 加密后落库。
func (r *NotificationChannelRepo) Create(ctx context.Context, req *notifyV1.CreateNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "invalid request")
	}
	if req.Data.GetName() == "" {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "channel name is required")
	}

	builder := r.entClient.Client().NotificationChannel.Create().
		SetName(req.Data.GetName()).
		SetNillableType(r.typeConverter.ToEntity(req.Data.Type)).
		SetNillableSMTPHost(req.Data.SmtpHost).
		SetNillableSMTPPort(req.Data.SmtpPort).
		SetNillableSMTPUsername(req.Data.SmtpUsername).
		SetNillableSMTPFrom(req.Data.SmtpFrom).
		SetNillableSMTPFromName(req.Data.SmtpFromName).
		SetNillableSMTPTLS(r.tlsConverter.ToEntity(req.Data.SmtpTls)).
		SetNillableWebhookURL(req.Data.WebhookUrl).
		SetNillableWebhookSignStyle(r.styleConverter.ToEntity(req.Data.WebhookSignStyle)).
		SetNillableWebhookPayloadTemplate(req.Data.WebhookPayloadTemplate).
		// 扩展渠道配置整包（service 层已按 schema 校验并加密 sensitive 字段）
		SetNillableSettings(req.Data.Settings).
		// SMS 渠道七列（Type=SMS 专用）
		SetNillableSmsProvider(r.smsProviderConverter.ToEntity(req.Data.SmsProvider)).
		SetNillableSmsAccessKeyID(req.Data.SmsAccessKeyId).
		SetNillableSmsSignName(req.Data.SmsSignName).
		SetNillableSmsTemplateCode(req.Data.SmsTemplateCode).
		SetNillableSmsAppID(req.Data.SmsAppId).
		SetNillableSmsRegion(req.Data.SmsRegion).
		SetStatus(statusFromEnabled(req.Data.GetEnabled())).
		SetNillableRemark(req.Data.Remark)

	if req.GetPassword() != "" {
		encrypted, err := pkgcrypto.EncryptIfNeeded(req.GetPassword())
		if err != nil {
			r.log.Errorf("encrypt smtp password failed: %s", err.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "encrypt password failed")
		}
		builder.SetSMTPPassword(encrypted)
	}
	if req.GetWebhookSecret() != "" {
		secret, err := pkgcrypto.EncryptIfNeeded(req.GetWebhookSecret())
		if err != nil {
			r.log.Errorf("encrypt webhook secret failed: %s", err.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "encrypt webhook secret failed")
		}
		builder.SetWebhookSecret(secret)
	}
	if req.GetSmsAccessSecret() != "" {
		encrypted, encErr := pkgcrypto.EncryptIfNeeded(req.GetSmsAccessSecret())
		if encErr != nil {
			r.log.Errorf("encrypt sms access secret failed: %s", encErr.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "encrypt sms access secret failed")
		}
		builder.SetSmsAccessSecret(encrypted)
	}
	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetCreatedBy(callerUserID).SetUpdatedBy(callerUserID)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		r.log.Errorf("insert notification channel failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "insert notification channel failed")
	}

	return r.Get(ctx, created.ID)
}

// Update 更新渠道；凭据留空表示不修改已存值。
func (r *NotificationChannelRepo) Update(ctx context.Context, req *notifyV1.UpdateNotificationChannelRequest) (*notifyV1.NotificationChannel, error) {
	if req == nil || req.Data == nil {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "invalid request")
	}
	if req.GetId() == 0 {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "id is required")
	}

	exist, err := r.entClient.Client().NotificationChannel.Query().
		Where(notificationchannel.IDEQ(req.GetId())).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_CHANNEL", "notification channel not found")
		}
		r.log.Errorf("get notification channel [%d] failed: %s", req.GetId(), err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "update notification channel failed")
	}

	builder := r.entClient.Client().NotificationChannel.UpdateOneID(exist.ID)
	if req.Data.Name != nil {
		builder.SetName(req.Data.GetName())
	}
	if req.Data.Type != nil {
		builder.SetNillableType(r.typeConverter.ToEntity(req.Data.Type))
	}
	if req.Data.SmtpHost != nil {
		builder.SetNillableSMTPHost(req.Data.SmtpHost)
	}
	if req.Data.SmtpPort != nil {
		builder.SetNillableSMTPPort(req.Data.SmtpPort)
	}
	if req.Data.SmtpUsername != nil {
		builder.SetNillableSMTPUsername(req.Data.SmtpUsername)
	}
	if req.Data.SmtpFrom != nil {
		builder.SetNillableSMTPFrom(req.Data.SmtpFrom)
	}
	if req.Data.SmtpFromName != nil {
		builder.SetNillableSMTPFromName(req.Data.SmtpFromName)
	}
	if req.Data.SmtpTls != nil {
		builder.SetNillableSMTPTLS(r.tlsConverter.ToEntity(req.Data.SmtpTls))
	}
	if req.Data.WebhookUrl != nil {
		builder.SetNillableWebhookURL(req.Data.WebhookUrl)
	}
	if req.Data.WebhookSignStyle != nil {
		builder.SetNillableWebhookSignStyle(r.styleConverter.ToEntity(req.Data.WebhookSignStyle))
	}
	if req.Data.WebhookPayloadTemplate != nil {
		builder.SetNillableWebhookPayloadTemplate(req.Data.WebhookPayloadTemplate)
	}
	// 扩展渠道配置整包（service 层已按 schema 合并：sensitive 留空 = 保留已存密文）
	if req.Data.Settings != nil {
		builder.SetNillableSettings(req.Data.Settings)
	}
	if req.Data.Enabled != nil {
		builder.SetStatus(statusFromEnabled(req.Data.GetEnabled()))
	}
	if req.Data.Remark != nil {
		builder.SetNillableRemark(req.Data.Remark)
	}
	if req.GetPassword() != "" {
		encrypted, encErr := pkgcrypto.EncryptIfNeeded(req.GetPassword())
		if encErr != nil {
			r.log.Errorf("encrypt smtp password failed: %s", encErr.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "encrypt password failed")
		}
		builder.SetSMTPPassword(encrypted)
	}
	if req.GetWebhookSecret() != "" {
		secret, encErr := pkgcrypto.EncryptIfNeeded(req.GetWebhookSecret())
		if encErr != nil {
			r.log.Errorf("encrypt webhook secret failed: %s", encErr.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "encrypt webhook secret failed")
		}
		builder.SetWebhookSecret(secret)
	}
	if req.GetSmsAccessSecret() != "" {
		encrypted, encErr := pkgcrypto.EncryptIfNeeded(req.GetSmsAccessSecret())
		if encErr != nil {
			r.log.Errorf("encrypt sms access secret failed: %s", encErr.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "encrypt sms access secret failed")
		}
		builder.SetSmsAccessSecret(encrypted)
	}
	if callerUserID, hasUser := viewerUserIDFromContext(ctx); hasUser {
		builder.SetUpdatedBy(callerUserID)
	}

	if _, err := builder.Save(ctx); err != nil {
		r.log.Errorf("update notification channel [%d] failed: %s", req.GetId(), err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "update notification channel failed")
	}

	return r.Get(ctx, req.GetId())
}

func (r *NotificationChannelRepo) Delete(ctx context.Context, id uint32) error {
	if id == 0 {
		return errors.BadRequest("NOTIFICATION_CHANNEL", "id is required")
	}

	if err := r.entClient.Client().NotificationChannel.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return errors.NotFound("NOTIFICATION_CHANNEL", "notification channel not found")
		}
		r.log.Errorf("delete notification channel [%d] failed: %s", id, err.Error())
		return errors.InternalServer("NOTIFICATION_CHANNEL", "delete notification channel failed")
	}

	return nil
}

// SmtpAccount 发送邮件所需的解密后 SMTP 配置（仅服务层内部使用，禁止外传）。
//
// ID 是这条配置在 notification_channels 里的主键：渠道自选后还要靠它把
// "实际用了哪个账号"写回投递台账。
type SmtpAccount struct {
	ID       uint32
	Host     string
	Port     uint32
	Username string
	Password string
	From     string
	FromName string
	TlsMode  string
	Enabled  bool
}

// GetFirstEnabledEmailChannel 取第一个启用的 EMAIL 渠道（事件路由未显式指定渠道时用）。
func (r *NotificationChannelRepo) GetFirstEnabledEmailChannel(ctx context.Context) (*SmtpAccount, error) {
	entities, err := r.entClient.Client().NotificationChannel.Query().
		Where(
			notificationchannel.TypeEQ(notificationchannel.TypeEmail),
			notificationchannel.StatusEQ(notificationchannel.StatusOn),
		).
		Order(ent.Asc(notificationchannel.FieldID)).
		Limit(1).
		All(ctx)
	if err != nil {
		r.log.Errorf("query enabled email channel failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "query email channel failed")
	}
	if len(entities) == 0 {
		return nil, errors.NotFound("NOTIFICATION_CHANNEL", "no enabled email channel")
	}
	return r.smtpAccountOf(ctx, entities[0])
}

// GetDecryptedSmtpAccount 取指定渠道的解密 SMTP 配置；类型不符即报错，
// 免得一条 WEBHOOK 行被当 SMTP 账号使（反向守卫在 EmailSender.resolveAccount）。
func (r *NotificationChannelRepo) GetDecryptedSmtpAccount(ctx context.Context, id uint32) (*SmtpAccount, error) {
	entity, err := r.entClient.Client().NotificationChannel.Query().
		Where(notificationchannel.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_CHANNEL", "notification channel not found")
		}
		r.log.Errorf("get notification channel [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "get notification channel failed")
	}
	if entity.Type == nil || *entity.Type != notificationchannel.TypeEmail {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "channel is not an EMAIL channel")
	}

	return r.smtpAccountOf(ctx, entity)
}

// smtpAccountOf 实体 → SmtpAccount，密码按 DecryptIfNeeded 解出明文。
// 解密失败上抛而不是降级成空密码：带认证的 SMTP 缺密码必然被拒，
// 让它失败在"配置坏"这一步比失败在对端 535 更可排障。
func (r *NotificationChannelRepo) smtpAccountOf(ctx context.Context, entity *ent.NotificationChannel) (*SmtpAccount, error) {
	password := ""
	if entity.SMTPPassword != nil {
		decrypted, decErr := pkgcrypto.DecryptIfNeeded(*entity.SMTPPassword)
		if decErr != nil {
			r.log.Errorf("decrypt smtp password failed for channel [%d]: %s", entity.ID, decErr.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "decrypt password failed")
		}
		password = decrypted
	}

	return &SmtpAccount{
		ID:       entity.ID,
		Host:     derefString(entity.SMTPHost),
		Port:     derefUint32(entity.SMTPPort),
		Username: derefString(entity.SMTPUsername),
		Password: password,
		From:     derefString(entity.SMTPFrom),
		FromName: derefString(entity.SMTPFromName),
		// SMTPTLS 为可空列：存量行可能为 NULL，必须 nil 安全取值
		TlsMode: derefEnumPtr(entity.SMTPTLS),
		Enabled: entity.Status != nil && *entity.Status == notificationchannel.StatusOn,
	}, nil
}

// WebhookAccount 投递一条 WEBHOOK 通知所需的解密后配置（仅服务层内部使用，禁止外传）。
//
// 与 SmtpAccount 同形：ID 是"实际用了哪条渠道配置"这条事实的来源，
// URL 是回调地址，Secret 为空表示不签名（签名头不下发）。
type WebhookAccount struct {
	ID     uint32
	URL    string
	Secret string

	// SignStyle 风格列的原始值，不在 repo 里归一化：认不出的值由 sender 报错，
	// 比悄悄换成 CUSTOM 更可排障——台账里看得出自家库里此刻写着什么。
	SignStyle string

	// PayloadTemplate 载荷模板，空串表示"用该风格的内置默认形状"。
	PayloadTemplate string
}

// SmsAccount SMS 渠道账号（Type=SMS 的解密投影）。
type SmsAccount struct {
	ID uint32
	// Provider 服务商（proto SmsProvider 成员名：ALIYUN/TENCENT）
	Provider string
	// AccessKeyID / AccessSecret 服务商访问密钥对；Secret 为空即配置不全，sender 拦下。
	AccessKeyID  string
	AccessSecret string
	// SignName 服务商侧已备案的短信签名
	SignName string
	// TemplateCode 验证码模板（阿里云正文含 ${code}，腾讯云正文含 {1}）
	TemplateCode string
	// AppID 腾讯云 SmsSdkAppId（必填）；阿里云留空
	AppID string
	// Region 腾讯云地域（必填，空则 sender 兜默认）；阿里云留空
	Region string
}

// GetFirstEnabledSmsChannel 取第一个启用的 SMS 渠道（验证码事件路由未显式指定渠道时用）。
func (r *NotificationChannelRepo) GetFirstEnabledSmsChannel(ctx context.Context) (*SmsAccount, error) {
	entities, err := r.entClient.Client().NotificationChannel.Query().
		Where(
			notificationchannel.TypeEQ(notificationchannel.TypeSms),
			notificationchannel.StatusEQ(notificationchannel.StatusOn),
		).
		Order(ent.Asc(notificationchannel.FieldID)).
		Limit(1).
		All(ctx)
	if err != nil {
		r.log.Errorf("query enabled sms channel failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "query sms channel failed")
	}
	if len(entities) == 0 {
		return nil, errors.NotFound("NOTIFICATION_CHANNEL", "no enabled sms channel")
	}
	return r.smsAccountOf(ctx, entities[0])
}

// GetDecryptedSmsAccount 取指定渠道的解密 SMS 配置；类型不符即报错，
// 免得一条 EMAIL/WEBHOOK 行被当短信账号使（反向守卫在 SmsSender.pickAccount）。
func (r *NotificationChannelRepo) GetDecryptedSmsAccount(ctx context.Context, id uint32) (*SmsAccount, error) {
	entity, err := r.entClient.Client().NotificationChannel.Query().
		Where(notificationchannel.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_CHANNEL", "notification channel not found")
		}
		r.log.Errorf("get sms channel [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "get notification channel failed")
	}
	if entity.Type == nil || *entity.Type != notificationchannel.TypeSms {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "channel is not an SMS channel")
	}
	return r.smsAccountOf(ctx, entity)
}

// smsAccountOf 实体 → SmsAccount，密钥按 DecryptIfNeeded 解出明文。
func (r *NotificationChannelRepo) smsAccountOf(ctx context.Context, entity *ent.NotificationChannel) (*SmsAccount, error) {
	secret := ""
	if entity.SmsAccessSecret != nil {
		decrypted, decErr := pkgcrypto.DecryptIfNeeded(*entity.SmsAccessSecret)
		if decErr != nil {
			r.log.Errorf("decrypt sms access secret failed for channel [%d]: %s", entity.ID, decErr.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "decrypt sms access secret failed")
		}
		secret = decrypted
	}

	return &SmsAccount{
		ID:           entity.ID,
		Provider:     derefEnum(entity.SmsProvider),
		AccessKeyID:  derefString(entity.SmsAccessKeyID),
		AccessSecret: secret,
		SignName:     derefString(entity.SmsSignName),
		TemplateCode: derefString(entity.SmsTemplateCode),
		AppID:        derefString(entity.SmsAppID),
		Region:       derefString(entity.SmsRegion),
	}, nil
}

// GetFirstEnabledWebhookChannel 取第一个启用的 WEBHOOK 渠道（事件路由未显式指定渠道时用）。
func (r *NotificationChannelRepo) GetFirstEnabledWebhookChannel(ctx context.Context) (*WebhookAccount, error) {
	entities, err := r.entClient.Client().NotificationChannel.Query().
		Where(
			notificationchannel.TypeEQ(notificationchannel.TypeWebhook),
			notificationchannel.StatusEQ(notificationchannel.StatusOn),
		).
		Order(ent.Asc(notificationchannel.FieldID)).
		Limit(1).
		All(ctx)
	if err != nil {
		r.log.Errorf("query enabled webhook channel failed: %s", err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "query webhook channel failed")
	}
	if len(entities) == 0 {
		return nil, errors.NotFound("NOTIFICATION_CHANNEL", "no enabled webhook channel")
	}

	return r.webhookAccountOf(ctx, entities[0])
}

// GetDecryptedWebhookAccount 取指定渠道的解密 WEBHOOK 配置；类型不符即报错，
// 免得一条 SMTP 行被当 webhook 目标使（反向守卫在 WebhookSender.resolveAccount）。
func (r *NotificationChannelRepo) GetDecryptedWebhookAccount(ctx context.Context, id uint32) (*WebhookAccount, error) {
	entity, err := r.entClient.Client().NotificationChannel.Query().
		Where(notificationchannel.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.NotFound("NOTIFICATION_CHANNEL", "notification channel not found")
		}
		r.log.Errorf("get notification channel [%d] failed: %s", id, err.Error())
		return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "get notification channel failed")
	}
	if entity.Type == nil || *entity.Type != notificationchannel.TypeWebhook {
		return nil, errors.BadRequest("NOTIFICATION_CHANNEL", "channel is not a WEBHOOK channel")
	}

	return r.webhookAccountOf(ctx, entity)
}

// webhookAccountOf 实体 → WebhookAccount，密钥按 DecryptIfNeeded 解出明文。
//
// 解密失败上抛而不是降级成"不签名发出去"：对端配置了密钥却要收不到签名，
// 是让它把"验签失败"当成对端的故障，比这条投递直接失败更难排。
func (r *NotificationChannelRepo) webhookAccountOf(ctx context.Context, entity *ent.NotificationChannel) (*WebhookAccount, error) {
	secret := ""
	if entity.WebhookSecret != nil {
		decrypted, decErr := pkgcrypto.DecryptIfNeeded(*entity.WebhookSecret)
		if decErr != nil {
			r.log.Errorf("decrypt webhook secret failed for channel [%d]: %s", entity.ID, decErr.Error())
			return nil, errors.InternalServer("NOTIFICATION_CHANNEL", "decrypt webhook secret failed")
		}
		secret = decrypted
	}

	return &WebhookAccount{
		ID:     entity.ID,
		URL:    derefString(entity.WebhookURL),
		Secret: secret,

		// 两列都是 Optional().Nillable()：ent 迁移只补 DEFAULT 给新建行，存量行是 NULL，
		// 所以这里必须能容忍空值（风格空 → sender 按 CUSTOM 走，模板空 → 内置默认形状）。
		SignStyle:       derefEnumPtr(entity.WebhookSignStyle),
		PayloadTemplate: derefEnumPtr(entity.WebhookPayloadTemplate),
	}, nil
}

// statusFromEnabled DTO enabled 布尔 → ent status 枚举
// GetSettingsRaw 取扩展渠道的存储态配置原文（settings 为加密 JSON，schema 感知的解密在 channel 包）。
// 返回：类型名（proto Type 枚举名）、settings 存储态 JSON、启用状态。
func (r *NotificationChannelRepo) GetSettingsRaw(ctx context.Context, id uint32) (typeName, storeJSON string, enabled bool, err error) {
	entity, getErr := r.entClient.Client().NotificationChannel.Query().
		Where(notificationchannel.IDEQ(id)).
		Only(ctx)
	if getErr != nil {
		if ent.IsNotFound(getErr) {
			return "", "", false, errors.NotFound("NOTIFICATION_CHANNEL", "notification channel not found")
		}
		r.log.Errorf("get settings channel [%d] failed: %s", id, getErr.Error())
		return "", "", false, errors.InternalServer("NOTIFICATION_CHANNEL", "get notification channel failed")
	}
	if entity.Type == nil {
		return "", "", false, errors.BadRequest("NOTIFICATION_CHANNEL", "channel type is empty")
	}
	return string(*entity.Type), derefString(entity.Settings), entity.Status != nil && *entity.Status == notificationchannel.StatusOn, nil
}

// GetFirstEnabledSettingsRaw 自选策略：类型落在 types 清单内的第一个启用行。
func (r *NotificationChannelRepo) GetFirstEnabledSettingsRaw(ctx context.Context, types []string) (typeName, storeJSON string, err error) {
	if len(types) == 0 {
		return "", "", errors.NotFound("NOTIFICATION_CHANNEL", "no enabled settings channel")
	}
	typeValues := make([]notificationchannel.Type, 0, len(types))
	for _, t := range types {
		typeValues = append(typeValues, notificationchannel.Type(t))
	}
	entities, queryErr := r.entClient.Client().NotificationChannel.Query().
		Where(
			notificationchannel.TypeIn(typeValues...),
			notificationchannel.StatusEQ(notificationchannel.StatusOn),
		).
		Order(ent.Asc(notificationchannel.FieldID)).
		Limit(1).
		All(ctx)
	if queryErr != nil {
		r.log.Errorf("query enabled settings channel failed: %s", queryErr.Error())
		return "", "", errors.InternalServer("NOTIFICATION_CHANNEL", "query settings channel failed")
	}
	if len(entities) == 0 {
		return "", "", errors.NotFound("NOTIFICATION_CHANNEL", "no enabled settings channel")
	}
	if entities[0].Type == nil {
		return "", "", errors.BadRequest("NOTIFICATION_CHANNEL", "channel type is empty")
	}
	return string(*entities[0].Type), derefString(entities[0].Settings), nil
}

func derefEnum[T ~string](p *T) string {
	if p == nil {
		return ""
	}
	return string(*p)
}

func statusFromEnabled(enabled bool) notificationchannel.Status {
	if enabled {
		return notificationchannel.StatusOn
	}
	return notificationchannel.StatusOff
}

// derefUint32 nil 安全解引用（ent 数值字段为 nillable 指针）。
func derefUint32(v *uint32) uint32 {
	if v == nil {
		return 0
	}
	return *v
}

// derefEnumPtr nil 安全解引用字符串型枚举指针（ent 可空枚举列读回 *T）。
func derefEnumPtr[T ~string](p *T) string {
	if p == nil {
		return ""
	}
	return string(*p)
}
