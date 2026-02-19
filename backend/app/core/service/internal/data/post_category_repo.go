package data

import (
	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PostCategoryRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.PostCategory, ent.PostCategory]

	repository *entCrud.Repository[
		ent.PostCategoryQuery, ent.PostCategorySelect,
		ent.PostCategoryCreate, ent.PostCategoryCreateBulk,
		ent.PostCategoryUpdate, ent.PostCategoryUpdateOne,
		ent.PostCategoryDelete,
		predicate.PostCategory,
		contentV1.PostCategory, ent.PostCategory,
	]
}

func NewPostCategoryRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *PostCategoryRepo {
	repo := &PostCategoryRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("post-category/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.PostCategory, ent.PostCategory](),
	}

	repo.init()

	return repo
}

func (r *PostCategoryRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.PostCategoryQuery, ent.PostCategorySelect,
		ent.PostCategoryCreate, ent.PostCategoryCreateBulk,
		ent.PostCategoryUpdate, ent.PostCategoryUpdateOne,
		ent.PostCategoryDelete,
		predicate.PostCategory,
		contentV1.PostCategory, ent.PostCategory,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}
