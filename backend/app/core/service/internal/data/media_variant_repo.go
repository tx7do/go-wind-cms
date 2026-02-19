package data

import (
	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	mediaV1 "go-wind-cms/api/gen/go/media/service/v1"
)

type MediaVariantRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[mediaV1.MediaVariant, ent.MediaVariant]

	repository *entCrud.Repository[
		ent.MediaVariantQuery, ent.MediaVariantSelect,
		ent.MediaVariantCreate, ent.MediaVariantCreateBulk,
		ent.MediaVariantUpdate, ent.MediaVariantUpdateOne,
		ent.MediaVariantDelete,
		predicate.MediaVariant,
		mediaV1.MediaVariant, ent.MediaVariant,
	]
}

func NewMediaVariantRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *MediaVariantRepo {
	repo := &MediaVariantRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("media-variant/repo/core-service"),
		mapper:    mapper.NewCopierMapper[mediaV1.MediaVariant, ent.MediaVariant](),
	}

	repo.init()

	return repo
}

func (r *MediaVariantRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.MediaVariantQuery, ent.MediaVariantSelect,
		ent.MediaVariantCreate, ent.MediaVariantCreateBulk,
		ent.MediaVariantUpdate, ent.MediaVariantUpdateOne,
		ent.MediaVariantDelete,
		predicate.MediaVariant,
		mediaV1.MediaVariant, ent.MediaVariant,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}
