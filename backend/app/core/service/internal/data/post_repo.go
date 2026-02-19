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
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/post"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

type PostRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[contentV1.Post, ent.Post]

	repository *entCrud.Repository[
		ent.PostQuery, ent.PostSelect,
		ent.PostCreate, ent.PostCreateBulk,
		ent.PostUpdate, ent.PostUpdateOne,
		ent.PostDelete,
		predicate.Post,
		contentV1.Post, ent.Post,
	]

	statusConverter     *mapper.EnumTypeConverter[contentV1.Post_PostStatus, post.Status]
	editorTypeConverter *mapper.EnumTypeConverter[contentV1.Post_EditorType, post.EditorType]

	postTranslationRepo *PostTranslationRepo
	postCategoryRepo    *PostCategoryRepo
}

func NewPostRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	postTranslationRepo *PostTranslationRepo,
	postCategoryRepo *PostCategoryRepo,
) *PostRepo {
	repo := &PostRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("post/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.Post, ent.Post](),
		statusConverter: mapper.NewEnumTypeConverter[contentV1.Post_PostStatus, post.Status](
			contentV1.Post_PostStatus_name, contentV1.Post_PostStatus_value,
		),
		editorTypeConverter: mapper.NewEnumTypeConverter[contentV1.Post_EditorType, post.EditorType](
			contentV1.Post_EditorType_name, contentV1.Post_EditorType_value,
		),
		postTranslationRepo: postTranslationRepo,
		postCategoryRepo:    postCategoryRepo,
	}

	repo.init()

	return repo
}

func (r *PostRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.PostQuery, ent.PostSelect,
		ent.PostCreate, ent.PostCreateBulk,
		ent.PostUpdate, ent.PostUpdateOne,
		ent.PostDelete,
		predicate.Post,
		contentV1.Post, ent.Post,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
	r.mapper.AppendConverters(r.editorTypeConverter.NewConverterPair())
}

func (r *PostRepo) count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().Post.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
	}

	return count, err
}

func (r *PostRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Post.Query().
		Where(post.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query post exist failed: %s", err.Error())
		return false, contentV1.ErrorInternalServerError("query post exist failed")
	}
	return exist, nil
}

func (r *PostRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Post.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListPostResponse{Total: 0, Items: nil}, nil
	}

	for _, item := range ret.Items {
		translations, err := r.postTranslationRepo.ListTranslations(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query translations failed")
		}
		item.Translations = translations
	}

	return &contentV1.ListPostResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *PostRepo) Get(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.Post, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().Post.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("post not found")
		}

		r.log.Errorf("query post failed: %s", err.Error())

		return nil, contentV1.ErrorInternalServerError("query post failed")
	}

	dto := r.mapper.ToDTO(entity)

	translations, err := r.postTranslationRepo.ListTranslations(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query translations failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query translations failed")
	}
	dto.Translations = translations

	return dto, nil
}

func (r *PostRepo) Create(ctx context.Context, req *contentV1.CreatePostRequest) (dto *contentV1.Post, err error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

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

	builder := tx.Post.Create().
		SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
		SetNillableEditorType(r.editorTypeConverter.ToEntity(req.Data.EditorType)).
		SetNillableSlug(req.Data.Slug).
		SetNillableDisallowComment(req.Data.DisallowComment).
		SetNillableInProgress(req.Data.InProgress).
		SetNillableAutoSummary(req.Data.AutoSummary).
		SetNillableIsFeatured(req.Data.IsFeatured).
		SetNillableSortOrder(req.Data.SortOrder).
		SetNillableVisits(req.Data.Visits).
		SetNillableLikes(req.Data.Likes).
		SetNillableCommentCount(req.Data.CommentCount).
		SetNillableAuthorID(req.Data.AuthorId).
		SetNillableAuthorName(req.Data.AuthorName).
		SetNillablePasswordHash(req.Data.PasswordHash).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	if req.Data.CustomFields != nil {
		builder.SetCustomFields(trans.Ptr(req.Data.GetCustomFields()))
	}

	var entity *ent.Post
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert post failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert post failed")
	}

	if req.Data.Translations != nil {
		if err = r.postTranslationRepo.CleanTranslations(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}
		if err = r.postTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *PostRepo) Update(ctx context.Context, req *contentV1.UpdatePostRequest) (dto *contentV1.Post, err error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
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
			_, err = r.Create(ctx, &contentV1.CreatePostRequest{Data: req.Data})
			return nil, err
		}
	}

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

	if req.Data.Translations != nil {
		if err = r.postTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}
		if err = r.postTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	builder := tx.Post.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *contentV1.Post) {
			builder.
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableEditorType(r.editorTypeConverter.ToEntity(req.Data.EditorType)).
				SetNillableSlug(req.Data.Slug).
				SetNillableDisallowComment(req.Data.DisallowComment).
				SetNillableInProgress(req.Data.InProgress).
				SetNillableAutoSummary(req.Data.AutoSummary).
				SetNillableIsFeatured(req.Data.IsFeatured).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableVisits(req.Data.Visits).
				SetNillableLikes(req.Data.Likes).
				SetNillableCommentCount(req.Data.CommentCount).
				SetNillableAuthorID(req.Data.AuthorId).
				SetNillableAuthorName(req.Data.AuthorName).
				SetNillablePasswordHash(req.Data.PasswordHash).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetUpdatedAt(time.Now())

			if req.Data.CustomFields != nil {
				builder.SetCustomFields(trans.Ptr(req.Data.GetCustomFields()))
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(post.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *PostRepo) Delete(ctx context.Context, req *contentV1.DeletePostRequest) (err error) {
	if req == nil {
		return contentV1.ErrorBadRequest("invalid parameter")
	}

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

	if err = r.entClient.Client().Post.
		DeleteOneID(req.GetId()).
		Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	if err = r.postTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean translations failed")
	}

	return err
}
