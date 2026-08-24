package data

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/category"
	"go-wind-cms/app/core/service/internal/data/ent/contentmodel"
	"go-wind-cms/app/core/service/internal/data/ent/contentmodeltranslation"
	"go-wind-cms/app/core/service/internal/data/ent/fielddefinition"
	"go-wind-cms/app/core/service/internal/data/ent/fielddefinitiontranslation"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/schema"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

var (
	// relationEntityTypeWhitelist relation 字段允许引用的实体类型白名单。
	relationEntityTypeWhitelist = map[string]bool{
		"post":     true,
		"page":     true,
		"category": true,
		"tag":      true,
	}
	// mediaRefPrefix image/file 字段值的字符串化引用前缀。
	mediaRefPrefix = "media_asset:"
)

// ContentModelRepo 内容模型仓库：模型 CRUD 与 fields/translations 级联读写，
// 以及字段值校验（ValidateValues，供内容服务在写入 custom_fields 前调用）。
type ContentModelRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.ContentModel, ent.ContentModel]

	fieldMapper       *mapper.CopierMapper[contentV1.FieldDefinition, ent.FieldDefinition]
	modelTrMapper     *mapper.CopierMapper[contentV1.ContentModelTranslation, ent.ContentModelTranslation]
	fieldTrMapper     *mapper.CopierMapper[contentV1.FieldDefinitionTranslation, ent.FieldDefinitionTranslation]
	typeConverter     *mapper.EnumTypeConverter[contentV1.FieldDefinition_Type, fielddefinition.Type]
	relationConverter *mapper.CopierMapper[contentV1.RelationConfig, schema.RelationConfig]

	repository *entCrud.Repository[
		ent.ContentModelQuery, ent.ContentModelSelect,
		ent.ContentModelCreate, ent.ContentModelCreateBulk,
		ent.ContentModelUpdate, ent.ContentModelUpdateOne,
		ent.ContentModelDelete,
		predicate.ContentModel,
		contentV1.ContentModel, ent.ContentModel,
	]
}

func NewContentModelRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *ContentModelRepo {
	repo := &ContentModelRepo{
		log:        ctx.NewLoggerHelper("content-model/repo/core-service"),
		entClient:  entClient,
		mapper:     mapper.NewCopierMapper[contentV1.ContentModel, ent.ContentModel](),
		fieldMapper: mapper.NewCopierMapper[contentV1.FieldDefinition, ent.FieldDefinition](),
		modelTrMapper: mapper.NewCopierMapper[contentV1.ContentModelTranslation, ent.ContentModelTranslation](),
		fieldTrMapper: mapper.NewCopierMapper[contentV1.FieldDefinitionTranslation, ent.FieldDefinitionTranslation](),
		typeConverter: mapper.NewEnumTypeConverter[contentV1.FieldDefinition_Type, fielddefinition.Type](
			contentV1.FieldDefinition_Type_name, contentV1.FieldDefinition_Type_value,
		),
		relationConverter: mapper.NewCopierMapper[contentV1.RelationConfig, schema.RelationConfig](),
	}

	repo.init()

	return repo
}

func (r *ContentModelRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.ContentModelQuery, ent.ContentModelSelect,
		ent.ContentModelCreate, ent.ContentModelCreateBulk,
		ent.ContentModelUpdate, ent.ContentModelUpdateOne,
		ent.ContentModelDelete,
		predicate.ContentModel,
		contentV1.ContentModel, ent.ContentModel,
	](r.mapper)

	for _, m := range []*mapper.CopierMapper[contentV1.ContentModelTranslation, ent.ContentModelTranslation]{r.modelTrMapper} {
		m.AppendConverters(copierutil.NewTimeStringConverterPair())
		m.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	}
	for _, m := range []*mapper.CopierMapper[contentV1.FieldDefinitionTranslation, ent.FieldDefinitionTranslation]{r.fieldTrMapper} {
		m.AppendConverters(copierutil.NewTimeStringConverterPair())
		m.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
	}
}

func (r *ContentModelRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListContentModelResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().ContentModel.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListContentModelResponse{Total: 0, Items: nil}, nil
	}

	return &contentV1.ListContentModelResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *ContentModelRepo) Count(ctx context.Context, req *paginationV1.PagingRequest) (int, error) {
	builder := r.entClient.Client().ContentModel.Query()

	whereSelectors, _, err := r.repository.BuildListSelectorWithPaging(builder, req)
	if len(whereSelectors) != 0 {
		builder.Modify(whereSelectors...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
		return 0, contentV1.ErrorInternalServerError("query count failed")
	}

	return count, nil
}

func (r *ContentModelRepo) Get(ctx context.Context, req *contentV1.GetContentModelRequest) (*contentV1.ContentModel, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().ContentModel.Query()

	switch req.QueryBy.(type) {
	case *contentV1.GetContentModelRequest_Id:
		builder.Where(contentmodel.IDEQ(req.GetId()))
	case *contentV1.GetContentModelRequest_Code:
		builder.Where(contentmodel.CodeEQ(req.GetCode()))
	default:
		return nil, contentV1.ErrorBadRequest("invalid query_by value")
	}

	entity, err := builder.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorNotFound("content model not found")
		}
		r.log.Errorf("query content model failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query content model failed")
	}

	dto := r.mapper.ToDTO(entity)

	// 填充字段定义与翻译（Get 详情默认全量返回）
	if err := r.fillModelAssociations(ctx, dto); err != nil {
		return nil, err
	}

	return dto, nil
}

func (r *ContentModelRepo) Create(ctx context.Context, req *contentV1.CreateContentModelRequest) (*contentV1.ContentModel, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	var err error
	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("start transaction failed")
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
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	builder := tx.ContentModel.Create().
		SetNillableTenantID(req.Data.TenantId).
		SetNillableName(req.Data.Name).
		SetNillableCode(req.Data.Code).
		SetNillableDescription(req.Data.Description).
		SetNillableSortOrder(req.Data.SortOrder).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	var entity *ent.ContentModel
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert content model failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert content model failed")
	}

	if err = r.batchCreateFields(ctx, tx, entity.ID, req.Data.Fields); err != nil {
		return nil, err
	}

	if err = r.batchCreateModelTranslations(ctx, tx, entity.ID, req.Data.Translations); err != nil {
		return nil, err
	}

	dto := r.mapper.ToDTO(entity)
	return dto, nil
}

func (r *ContentModelRepo) Update(ctx context.Context, req *contentV1.UpdateContentModelRequest) (*contentV1.ContentModel, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	var err error
	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("start transaction failed")
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
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	updateBuilder := tx.ContentModel.UpdateOneID(req.GetId()).
		SetNillableName(req.Data.Name).
		SetNillableDescription(req.Data.Description).
		SetNillableSortOrder(req.Data.SortOrder).
		SetUpdatedAt(time.Now())

	entity, uErr := updateBuilder.Save(ctx)
	if uErr != nil {
		if ent.IsNotFound(uErr) {
			err = contentV1.ErrorNotFound("content model not found")
			return nil, err
		}
		r.log.Errorf("update content model failed: %s", uErr.Error())
		err = contentV1.ErrorInternalServerError("update content model failed")
		return nil, err
	}

	// 字段定义与翻译整体替换（对齐 post translations 惯例）
	if err = r.replaceFields(ctx, tx, req.GetId(), req.Data.Fields); err != nil {
		return nil, err
	}

	if err = r.replaceModelTranslations(ctx, tx, req.GetId(), req.Data.Translations); err != nil {
		return nil, err
	}

	dto := r.mapper.ToDTO(entity)
	return dto, nil
}

func (r *ContentModelRepo) Delete(ctx context.Context, req *contentV1.DeleteContentModelRequest) error {
	if req == nil {
		return contentV1.ErrorBadRequest("invalid parameter")
	}

	var err error
	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("start transaction failed")
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
			err = contentV1.ErrorInternalServerError("transaction commit failed")
		}
	}()

	modelID := req.GetId()

	// 级联清理：字段翻译 → 字段定义 → 模型翻译 → 模型
	fieldIDs, qErr := tx.FieldDefinition.Query().
		Where(fielddefinition.ContentModelIDEQ(modelID)).
		IDs(ctx)
	if qErr != nil {
		r.log.Errorf("query field definition ids failed: %s", qErr.Error())
		err = contentV1.ErrorInternalServerError("query field definition ids failed")
		return err
	}
	if len(fieldIDs) > 0 {
		if _, dErr := tx.FieldDefinitionTranslation.Delete().
			Where(fielddefinitiontranslation.FieldDefinitionIDIn(fieldIDs...)).
			Exec(ctx); dErr != nil {
			r.log.Errorf("delete field definition translations failed: %s", dErr.Error())
			err = contentV1.ErrorInternalServerError("delete field definition translations failed")
			return err
		}
	}
	if _, dErr := tx.FieldDefinition.Delete().
		Where(fielddefinition.ContentModelIDEQ(modelID)).
		Exec(ctx); dErr != nil {
		r.log.Errorf("delete field definitions failed: %s", dErr.Error())
		err = contentV1.ErrorInternalServerError("delete field definitions failed")
		return err
	}
	if _, dErr := tx.ContentModelTranslation.Delete().
		Where(contentmodeltranslation.ContentModelIDEQ(modelID)).
		Exec(ctx); dErr != nil {
		r.log.Errorf("delete content model translations failed: %s", dErr.Error())
		err = contentV1.ErrorInternalServerError("delete content model translations failed")
		return err
	}
	if dErr := tx.ContentModel.DeleteOneID(modelID).Exec(ctx); dErr != nil {
		if ent.IsNotFound(dErr) {
			err = contentV1.ErrorNotFound("content model not found")
			return err
		}
		r.log.Errorf("delete content model failed: %s", dErr.Error())
		err = contentV1.ErrorInternalServerError("delete content model failed")
		return err
	}

	return nil
}

// ListFieldDefinitions 列出模型下的字段定义（含各自翻译），供内容编辑器拉动态表单。
func (r *ContentModelRepo) ListFieldDefinitions(ctx context.Context, req *contentV1.ListFieldDefinitionsRequest) (*contentV1.ListFieldDefinitionsResponse, error) {
	if req == nil || req.GetContentModelId() == 0 {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	entities, err := r.entClient.Client().FieldDefinition.Query().
		Where(fielddefinition.ContentModelIDEQ(req.GetContentModelId())).
		Order(ent.Asc(fielddefinition.FieldSortOrder), ent.Asc(fielddefinition.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query field definitions failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query field definitions failed")
	}

	items := make([]*contentV1.FieldDefinition, 0, len(entities))
	for _, entity := range entities {
		dto := r.fieldDefinitionToDTO(entity)
		translations, tErr := r.listFieldTranslations(ctx, entity.ID)
		if tErr != nil {
			return nil, tErr
		}
		dto.Translations = translations
		items = append(items, dto)
	}

	return &contentV1.ListFieldDefinitionsResponse{
		Items: items,
		Total: uint64(len(items)),
	}, nil
}

// ============ 关联读写辅助 ============

func (r *ContentModelRepo) fillModelAssociations(ctx context.Context, dto *contentV1.ContentModel) error {
	fields, err := r.entClient.Client().FieldDefinition.Query().
		Where(fielddefinition.ContentModelIDEQ(dto.GetId())).
		Order(ent.Asc(fielddefinition.FieldSortOrder), ent.Asc(fielddefinition.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query field definitions failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("query field definitions failed")
	}

	dto.Fields = make([]*contentV1.FieldDefinition, 0, len(fields))
	for _, f := range fields {
		fieldDTO := r.fieldDefinitionToDTO(f)
		translations, tErr := r.listFieldTranslations(ctx, f.ID)
		if tErr != nil {
			return tErr
		}
		fieldDTO.Translations = translations
		dto.Fields = append(dto.Fields, fieldDTO)
	}

	modelTranslations, err := r.entClient.Client().ContentModelTranslation.Query().
		Where(contentmodeltranslation.ContentModelIDEQ(dto.GetId())).
		Order(ent.Asc(contentmodeltranslation.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query content model translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("query content model translations failed")
	}
	for _, t := range modelTranslations {
		dto.Translations = append(dto.Translations, r.modelTrMapper.ToDTO(t))
	}

	return nil
}

func (r *ContentModelRepo) listFieldTranslations(ctx context.Context, fieldDefinitionID uint32) ([]*contentV1.FieldDefinitionTranslation, error) {
	entities, err := r.entClient.Client().FieldDefinitionTranslation.Query().
		Where(fielddefinitiontranslation.FieldDefinitionIDEQ(fieldDefinitionID)).
		Order(ent.Asc(fielddefinitiontranslation.FieldID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query field definition translations failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query field definition translations failed")
	}

	dtos := make([]*contentV1.FieldDefinitionTranslation, 0, len(entities))
	for _, e := range entities {
		dtos = append(dtos, r.fieldTrMapper.ToDTO(e))
	}
	return dtos, nil
}

func (r *ContentModelRepo) fieldDefinitionToDTO(entity *ent.FieldDefinition) *contentV1.FieldDefinition {
	dto := &contentV1.FieldDefinition{
		Id:              trans.Ptr(entity.ID),
		ContentModelId:  entity.ContentModelID,
		Name:            entity.Name,
		Type:            r.typeConverter.ToDTO(entity.Type),
		Label:           entity.Label,
		Description:     entity.Description,
		Placeholder:     entity.Placeholder,
		IsRequired:      entity.IsRequired,
		ValidationRegex: entity.ValidationRegex,
		SortOrder:       entity.SortOrder,
	}
	if entity.Options != nil {
		dto.Options = *entity.Options
	}
	if entity.RelationConfig != nil {
		dto.RelationConfig = r.relationConverter.ToDTO(entity.RelationConfig)
	}
	return dto
}

func (r *ContentModelRepo) newFieldCreateBuilder(fd *ent.FieldDefinitionClient, data *contentV1.FieldDefinition) *ent.FieldDefinitionCreate {
	builder := fd.Create().
		SetNillableContentModelID(data.ContentModelId).
		SetNillableName(data.Name).
		SetNillableType(r.typeConverter.ToEntity(data.Type)).
		SetNillableLabel(data.Label).
		SetNillableDescription(data.Description).
		SetNillablePlaceholder(data.Placeholder).
		SetNillableIsRequired(data.IsRequired).
		SetNillableValidationRegex(data.ValidationRegex).
		SetNillableSortOrder(data.SortOrder).
		SetNillableCreatedBy(data.CreatedBy).
		SetCreatedAt(time.Now())

	if data.Options != nil {
		builder.SetOptions(&data.Options)
	}
	if data.RelationConfig != nil {
		if rc := r.relationConverter.ToEntity(data.RelationConfig); rc != nil {
			builder.SetRelationConfig(rc)
		}
	}

	return builder
}

func (r *ContentModelRepo) batchCreateFields(ctx context.Context, tx *ent.Tx, modelID uint32, fields []*contentV1.FieldDefinition) error {
	for _, f := range fields {
		if f == nil || f.GetName() == "" {
			return contentV1.ErrorBadRequest("field definition name is required")
		}
		f.ContentModelId = trans.Ptr(modelID)
		created, err := r.newFieldCreateBuilder(tx.FieldDefinition, f).Save(ctx)
		if err != nil {
			r.log.Errorf("insert field definition failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("insert field definition failed")
		}
		for _, tr := range f.Translations {
			if tr == nil {
				continue
			}
			tr.FieldDefinitionId = trans.Ptr(created.ID)
			if _, err = tx.FieldDefinitionTranslation.Create().
				SetNillableFieldDefinitionID(tr.FieldDefinitionId).
				SetNillableLanguageCode(tr.LanguageCode).
				SetNillableLabel(tr.Label).
				SetNillableDescription(tr.Description).
				SetNillablePlaceholder(tr.Placeholder).
				SetNillableCreatedBy(tr.CreatedBy).
				SetCreatedAt(time.Now()).
				Save(ctx); err != nil {
				r.log.Errorf("insert field definition translation failed: %s", err.Error())
				return contentV1.ErrorInternalServerError("insert field definition translation failed")
			}
		}
	}
	return nil
}

func (r *ContentModelRepo) replaceFields(ctx context.Context, tx *ent.Tx, modelID uint32, fields []*contentV1.FieldDefinition) error {
	// 清旧（含字段翻译）
	oldIDs, err := tx.FieldDefinition.Query().
		Where(fielddefinition.ContentModelIDEQ(modelID)).
		IDs(ctx)
	if err != nil {
		r.log.Errorf("query old field definition ids failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("query old field definition ids failed")
	}
	if len(oldIDs) > 0 {
		if _, err = tx.FieldDefinitionTranslation.Delete().
			Where(fielddefinitiontranslation.FieldDefinitionIDIn(oldIDs...)).
			Exec(ctx); err != nil {
			r.log.Errorf("delete old field translations failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("delete old field translations failed")
		}
	}
	if _, err = tx.FieldDefinition.Delete().
		Where(fielddefinition.ContentModelIDEQ(modelID)).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old field definitions failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete old field definitions failed")
	}

	return r.batchCreateFields(ctx, tx, modelID, fields)
}

func (r *ContentModelRepo) batchCreateModelTranslations(ctx context.Context, tx *ent.Tx, modelID uint32, translations []*contentV1.ContentModelTranslation) error {
	for _, tr := range translations {
		if tr == nil || tr.GetLanguageCode() == "" {
			continue
		}
		tr.ContentModelId = trans.Ptr(modelID)
		if _, err := tx.ContentModelTranslation.Create().
			SetNillableContentModelID(tr.ContentModelId).
			SetNillableLanguageCode(tr.LanguageCode).
			SetNillableName(tr.Name).
			SetNillableDescription(tr.Description).
			SetNillableCreatedBy(tr.CreatedBy).
			SetCreatedAt(time.Now()).
			Save(ctx); err != nil {
			r.log.Errorf("insert content model translation failed: %s", err.Error())
			return contentV1.ErrorInternalServerError("insert content model translation failed")
		}
	}
	return nil
}

func (r *ContentModelRepo) replaceModelTranslations(ctx context.Context, tx *ent.Tx, modelID uint32, translations []*contentV1.ContentModelTranslation) error {
	if _, err := tx.ContentModelTranslation.Delete().
		Where(contentmodeltranslation.ContentModelIDEQ(modelID)).
		Exec(ctx); err != nil {
		r.log.Errorf("delete old model translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("delete old model translations failed")
	}
	return r.batchCreateModelTranslations(ctx, tx, modelID, translations)
}

// ============ 字段值校验 ============

// ResolveModelIDByCategories 返回分类列表里第一个绑定模型的 ID（0=无绑定）。
// 用于 Post 写入 custom_fields 前解析其所属分类链上的模型。
func (r *ContentModelRepo) ResolveModelIDByCategories(ctx context.Context, categoryIds []uint32) uint32 {
	if len(categoryIds) == 0 {
		return 0
	}

	entities, err := r.entClient.Client().Category.Query().
		Where(category.IDIn(categoryIds...)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query categories for model resolution failed: %s", err.Error())
		return 0
	}

	for _, c := range entities {
		if c.ContentModelID != nil && *c.ContentModelID != 0 {
			return *c.ContentModelID
		}
	}

	return 0
}

// ValidateValues 校验 custom_fields 的值符合绑定模型的字段定义。
// 校验项：必填、text 正则、number 可解析、image/file 引用格式（media_asset:ID）、
// relation 引用格式（entity:ID）+ 实体类型白名单 + 目标实体存在性与租户归属
// （AllowCrossTenant=false 时目标 tenant_id 必须等于当前租户）。
// 未在模型字段定义中的多余键会被拒绝，防止匿名写入任意键。
func (r *ContentModelRepo) ValidateValues(ctx context.Context, contentModelID uint32, customFields map[string]string, tenantID uint32) error {
	if contentModelID == 0 {
		// 未绑定模型：customFields 必须为空，防止绕过模型校验写入任意键
		if len(customFields) > 0 {
			return contentV1.ErrorBadRequest("custom fields are not allowed without a bound content model")
		}
		return nil
	}

	fieldEntities, err := r.entClient.Client().FieldDefinition.Query().
		Where(fielddefinition.ContentModelIDEQ(contentModelID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("query field definitions for validation failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("query field definitions failed")
	}

	// derefStr 解引用 *string（ent Nillable 字段），nil 安全返回空串
	derefStr := func(p *string) string {
		if p == nil {
			return ""
		}
		return *p
	}

	defMap := make(map[string]*ent.FieldDefinition, len(fieldEntities))
	for _, f := range fieldEntities {
		defMap[derefStr(f.Name)] = f
	}

	// 多余键检查
	for key := range customFields {
		if _, ok := defMap[key]; !ok {
			return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q is not defined in the bound content model", key))
		}
	}

	for _, def := range fieldEntities {
		name := derefStr(def.Name)
		value, exists := customFields[name]

		isRequired := def.IsRequired != nil && *def.IsRequired
		if isRequired && (!exists || strings.TrimSpace(value) == "") {
			return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q is required", name))
		}
		if !exists || value == "" {
			continue
		}

		fieldType := fielddefinition.DefaultType
		if def.Type != nil {
			fieldType = *def.Type
		}

		switch fieldType {
		case fielddefinition.TypeFIELD_TYPE_TEXT:
			if def.ValidationRegex != nil && *def.ValidationRegex != "" {
				re, cErr := regexp.Compile(*def.ValidationRegex)
				if cErr != nil {
					return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q has an invalid validation regex", name))
				}
				if !re.MatchString(value) {
					return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q does not match the validation rule", name))
				}
			}

		case fielddefinition.TypeFIELD_TYPE_NUMBER:
			if _, pErr := strconv.ParseFloat(value, 64); pErr != nil {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q must be a number", name))
			}

		case fielddefinition.TypeFIELD_TYPE_IMAGE, fielddefinition.TypeFIELD_TYPE_FILE:
			id, ok := parseStringRef(value, mediaRefPrefix)
			if !ok {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q must be a media reference like %s<id>", name, mediaRefPrefix))
			}
			// 媒体归属校验：存在且属于当前租户（与站点配置读取的字段裁剪同理，防止跨租户引用）
			media, mErr := r.entClient.Client().MediaAsset.Get(ctx, id)
			if mErr != nil {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q references a non-existent media asset", name))
			}
			if media.TenantID != nil && tenantID != 0 && *media.TenantID != tenantID {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q references a media asset from another tenant", name))
			}

		case fielddefinition.TypeFIELD_TYPE_RELATION:
			if def.RelationConfig == nil {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q is a relation without relation config", name))
			}
			entityType, id, ok := parseEntityRef(value)
			if !ok {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q must be an entity reference like post:<id>", name))
			}
			if !relationEntityTypeWhitelist[entityType] {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q references an unsupported entity type %q", name, entityType))
			}
			if def.RelationConfig.TargetEntityType != "" && def.RelationConfig.TargetEntityType != entityType {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q must reference a %s", name, def.RelationConfig.TargetEntityType))
			}
			targetTenant, tErr := r.queryEntityTenant(ctx, entityType, id)
			if tErr != nil {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q references a non-existent %s", name, entityType))
			}
			if !def.RelationConfig.AllowCrossTenant && tenantID != 0 && targetTenant != tenantID {
				return contentV1.ErrorBadRequest(fmt.Sprintf("custom field %q references a %s from another tenant", name, entityType))
			}
		}
	}

	return nil
}

// parseStringRef 解析 "<prefix><id>" 形式的字符串化引用。
func parseStringRef(value, prefix string) (uint32, bool) {
	if !strings.HasPrefix(value, prefix) {
		return 0, false
	}
	id, err := strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, 32)
	if err != nil || id == 0 {
		return 0, false
	}
	return uint32(id), true
}

// parseEntityRef 解析 "entity:id" 形式的关联引用。
func parseEntityRef(value string) (string, uint32, bool) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return "", 0, false
	}
	entityType := strings.ToLower(strings.TrimSpace(parts[0]))
	id, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 32)
	if err != nil || id == 0 {
		return "", 0, false
	}
	return entityType, uint32(id), true
}

// queryEntityTenant 查询目标实体的租户归属（白名单实体类型）。
func (r *ContentModelRepo) queryEntityTenant(ctx context.Context, entityType string, id uint32) (uint32, error) {
	var tenantID *uint32
	var err error
	switch entityType {
	case "post":
		var e *ent.Post
		if e, err = r.entClient.Client().Post.Get(ctx, id); err == nil {
			tenantID = e.TenantID
		}
	case "page":
		var e *ent.Page
		if e, err = r.entClient.Client().Page.Get(ctx, id); err == nil {
			tenantID = e.TenantID
		}
	case "category":
		var e *ent.Category
		if e, err = r.entClient.Client().Category.Get(ctx, id); err == nil {
			tenantID = e.TenantID
		}
	case "tag":
		var e *ent.Tag
		if e, err = r.entClient.Client().Tag.Get(ctx, id); err == nil {
			tenantID = e.TenantID
		}
	default:
		return 0, fmt.Errorf("unsupported entity type: %s", entityType)
	}
	if err != nil {
		return 0, err
	}
	if tenantID == nil {
		return 0, nil
	}
	return *tenantID, nil
}
