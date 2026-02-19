package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/sitesetting"

	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type SiteSettingRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[siteV1.SiteSetting, ent.SiteSetting]

	repository *entCrud.Repository[
		ent.SiteSettingQuery, ent.SiteSettingSelect,
		ent.SiteSettingCreate, ent.SiteSettingCreateBulk,
		ent.SiteSettingUpdate, ent.SiteSettingUpdateOne,
		ent.SiteSettingDelete,
		predicate.SiteSetting,
		siteV1.SiteSetting, ent.SiteSetting,
	]

	settingTypeConverter *mapper.EnumTypeConverter[siteV1.SiteSetting_SettingType, sitesetting.Type]
}

func NewSiteSettingRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SiteSettingRepo {
	repo := &SiteSettingRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("site-setting/repo/core-service"),
		mapper:    mapper.NewCopierMapper[siteV1.SiteSetting, ent.SiteSetting](),
		settingTypeConverter: mapper.NewEnumTypeConverter[siteV1.SiteSetting_SettingType, sitesetting.Type](
			siteV1.SiteSetting_SettingType_name, siteV1.SiteSetting_SettingType_value,
		),
	}

	repo.init()

	return repo
}

func (r *SiteSettingRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.SiteSettingQuery, ent.SiteSettingSelect,
		ent.SiteSettingCreate, ent.SiteSettingCreateBulk,
		ent.SiteSettingUpdate, ent.SiteSettingUpdateOne,
		ent.SiteSettingDelete,
		predicate.SiteSetting,
		siteV1.SiteSetting, ent.SiteSetting,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.settingTypeConverter.NewConverterPair())
}

func (r *SiteSettingRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().SiteSetting.Query().
		Where(sitesetting.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query site setting exist failed: %s", err.Error())
		return false, siteV1.ErrorInternalServerError("query site setting exist failed")
	}
	return exist, nil
}

func (r *SiteSettingRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListSiteSettingResponse, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SiteSetting.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &siteV1.ListSiteSettingResponse{Total: 0, Items: nil}, nil
	}

	return &siteV1.ListSiteSettingResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *SiteSettingRepo) Get(ctx context.Context, req *siteV1.GetSiteSettingRequest) (*siteV1.SiteSetting, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().SiteSetting.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, siteV1.ErrorFileNotFound("sitesetting not found")
		}

		r.log.Errorf("query sitesetting failed: %s", err.Error())

		return nil, siteV1.ErrorInternalServerError("query sitesetting failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *SiteSettingRepo) Create(ctx context.Context, req *siteV1.CreateSiteSettingRequest) (*siteV1.SiteSetting, error) {
	if req == nil || req.Data == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().SiteSetting.Create().
		SetNillableType(r.settingTypeConverter.ToEntity(req.Data.Type)).
		SetNillableSiteID(req.Data.SiteId).
		SetNillableLocale(req.Data.Locale).
		SetNillableGroup(req.Data.Group).
		SetNillableKey(req.Data.Key).
		SetNillableValue(req.Data.Value).
		SetNillableLabel(req.Data.Label).
		SetNillableDescription(req.Data.Description).
		SetNillablePlaceholder(req.Data.Placeholder).
		SetNillableIsRequired(req.Data.IsRequired).
		SetNillableValidationRegex(req.Data.ValidationRegex).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if req.Data.Options != nil {
		builder.SetOptions(trans.Ptr(req.Data.GetOptions()))
	}

	var err error
	var entity *ent.SiteSetting
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert sitesetting failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("insert sitesetting failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *SiteSettingRepo) Update(ctx context.Context, req *siteV1.UpdateSiteSettingRequest) (*siteV1.SiteSetting, error) {
	if req == nil || req.Data == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.IsExist(ctx, req.GetId())
		if err != nil {
			return nil, err
		}
		if !exist {
			req.Data.CreatedBy = req.Data.UpdatedBy
			req.Data.UpdatedBy = nil
			_, err = r.Create(ctx, &siteV1.CreateSiteSettingRequest{Data: req.Data})
			return nil, err
		}
	}

	builder := r.entClient.Client().SiteSetting.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *siteV1.SiteSetting) {
			builder.
				SetNillableType(r.settingTypeConverter.ToEntity(req.Data.Type)).
				SetNillableSiteID(req.Data.SiteId).
				SetNillableLocale(req.Data.Locale).
				SetNillableGroup(req.Data.Group).
				SetNillableKey(req.Data.Key).
				SetNillableValue(req.Data.Value).
				SetNillableLabel(req.Data.Label).
				SetNillableDescription(req.Data.Description).
				SetNillablePlaceholder(req.Data.Placeholder).
				SetNillableIsRequired(req.Data.IsRequired).
				SetNillableValidationRegex(req.Data.ValidationRegex).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if req.Data.Options != nil {
				builder.SetOptions(trans.Ptr(req.Data.GetOptions()))
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(sitesetting.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *SiteSettingRepo) Delete(ctx context.Context, req *siteV1.DeleteSiteSettingRequest) (bool, error) {
	err := r.entClient.Client().SiteSetting.
		DeleteOneID(req.GetId()).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	return err == nil, err
}
