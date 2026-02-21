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
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"
	"go-wind-cms/app/core/service/internal/data/ent/site"

	siteV1 "go-wind-cms/api/gen/go/site/service/v1"
)

type SiteRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[siteV1.Site, ent.Site]

	repository *entCrud.Repository[
		ent.SiteQuery, ent.SiteSelect,
		ent.SiteCreate, ent.SiteCreateBulk,
		ent.SiteUpdate, ent.SiteUpdateOne,
		ent.SiteDelete,
		predicate.Site,
		siteV1.Site, ent.Site,
	]

	statusConverter *mapper.EnumTypeConverter[siteV1.Site_Status, site.Status]
}

func NewSiteRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *SiteRepo {
	repo := &SiteRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("site/repo/core-service"),
		mapper:    mapper.NewCopierMapper[siteV1.Site, ent.Site](),
		statusConverter: mapper.NewEnumTypeConverter[siteV1.Site_Status, site.Status](
			siteV1.Site_Status_name, siteV1.Site_Status_value,
		),
	}

	repo.init()

	return repo
}

func (r *SiteRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.SiteQuery, ent.SiteSelect,
		ent.SiteCreate, ent.SiteCreateBulk,
		ent.SiteUpdate, ent.SiteUpdateOne,
		ent.SiteDelete,
		predicate.Site,
		siteV1.Site, ent.Site,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
}

func (r *SiteRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Site.Query().
		Where(site.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query site exist failed: %s", err.Error())
		return false, siteV1.ErrorInternalServerError("query site exist failed")
	}
	return exist, nil
}

func (r *SiteRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*siteV1.ListSiteResponse, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Site.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &siteV1.ListSiteResponse{Total: 0, Items: nil}, nil
	}

	return &siteV1.ListSiteResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *SiteRepo) Get(ctx context.Context, req *siteV1.GetSiteRequest) (*siteV1.Site, error) {
	if req == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().Site.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, siteV1.ErrorFileNotFound("site not found")
		}

		r.log.Errorf("query site failed: %s", err.Error())

		return nil, siteV1.ErrorInternalServerError("query site failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *SiteRepo) Create(ctx context.Context, req *siteV1.CreateSiteRequest) (*siteV1.Site, error) {
	if req == nil || req.Data == nil {
		return nil, siteV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Site.Create().
		SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
		SetNillableTenantID(req.Data.TenantId).
		SetNillableName(req.Data.Name).
		SetNillableSlug(req.Data.Slug).
		SetNillableDomain(req.Data.Domain).
		SetNillableIsDefault(req.Data.IsDefault).
		SetNillableDefaultLocale(req.Data.DefaultLocale).
		SetNillableTemplate(req.Data.Template).
		SetNillableTheme(req.Data.Theme).
		SetNillableVisitCount(req.Data.VisitCount).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if req.Data.AlternateDomains != nil {
		builder.SetAlternateDomains(req.Data.GetAlternateDomains())
	}

	var err error
	var entity *ent.Site
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert site failed: %s", err.Error())
		return nil, siteV1.ErrorInternalServerError("insert site failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *SiteRepo) Update(ctx context.Context, req *siteV1.UpdateSiteRequest) (*siteV1.Site, error) {
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
			_, err = r.Create(ctx, &siteV1.CreateSiteRequest{Data: req.Data})
			return nil, err
		}
	}

	builder := r.entClient.Client().Site.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *siteV1.Site) {
			builder.
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableName(req.Data.Name).
				SetNillableSlug(req.Data.Slug).
				SetNillableDomain(req.Data.Domain).
				SetNillableIsDefault(req.Data.IsDefault).
				SetNillableDefaultLocale(req.Data.DefaultLocale).
				SetNillableTemplate(req.Data.Template).
				SetNillableTheme(req.Data.Theme).
				SetNillableVisitCount(req.Data.VisitCount).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if req.Data.AlternateDomains != nil {
				builder.SetAlternateDomains(req.Data.GetAlternateDomains())
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(site.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *SiteRepo) Delete(ctx context.Context, req *siteV1.DeleteSiteRequest) (bool, error) {
	err := r.entClient.Client().Site.
		DeleteOneID(req.GetId()).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	return err == nil, err
}
