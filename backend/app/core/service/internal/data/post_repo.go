package data

import (
	"context"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/go-crud/pagination"
	paginationFilter "github.com/tx7do/go-crud/pagination/filter"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/post"
	"go-wind-cms/app/core/service/internal/data/ent/postcategory"
	"go-wind-cms/app/core/service/internal/data/ent/posttag"
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
	editorTypeConverter *mapper.EnumTypeConverter[contentV1.EditorType, post.EditorType]

	postTranslationRepo *PostTranslationRepo

	postCategoryRepo *PostCategoryRepo
	postTagRepo      *PostTagRepo
}

func NewPostRepo(
	ctx *bootstrap.Context,
	entClient *entCrud.EntClient[*ent.Client],
	postTranslationRepo *PostTranslationRepo,
	postCategoryRepo *PostCategoryRepo,
	postTagRepo *PostTagRepo,
) *PostRepo {
	repo := &PostRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("post/repo/core-service"),
		mapper:    mapper.NewCopierMapper[contentV1.Post, ent.Post](),
		statusConverter: mapper.NewEnumTypeConverter[contentV1.Post_PostStatus, post.Status](
			contentV1.Post_PostStatus_name, contentV1.Post_PostStatus_value,
		),
		editorTypeConverter: mapper.NewEnumTypeConverter[contentV1.EditorType, post.EditorType](
			contentV1.EditorType_name, contentV1.EditorType_value,
		),
		postTranslationRepo: postTranslationRepo,
		postCategoryRepo:    postCategoryRepo,
		postTagRepo:         postTagRepo,
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

	cnt, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
	}

	return cnt, err
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

func (r *PostRepo) prepareTranslationMaskFields(req *paginationV1.PagingRequest) (
	excludeConditions []*paginationV1.FilterCondition,
	translationMaskFields []string,
	needQueryTranslation bool,
	err error,
) {
	var filterExpr *paginationV1.FilterExpr
	filterExpr, err = paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		r.log.Errorf("convert filter by paging request failed: %s", err.Error())
		return
	}

	excludeFields := []string{"translations", "available_languages"}
	if req.FieldMask != nil && len(req.FieldMask.Paths) > 0 {
		for _, path := range req.FieldMask.Paths {
			path = strings.TrimSpace(path)

			if path == "translations" || strings.HasPrefix(path, "translations.") {
				needQueryTranslation = true
			}
			if strings.HasPrefix(path, "translations.") {
				excludeFields = append(excludeFields, path)
				path = strings.TrimPrefix(path, "translations.")
				if len(path) > 0 {
					translationMaskFields = append(translationMaskFields, path)
				}
			}
		}

		req.FieldMask = FilterViewMask(excludeFields, req.FieldMask)
	} else {
		needQueryTranslation = true
	}

	excludeConditions = pagination.FilterFields(filterExpr, []string{
		"locale",
		"tag_ids",
		"category_ids",
	})
	req.FilteringType = &paginationV1.PagingRequest_FilterExpr{FilterExpr: filterExpr}

	return
}

func (r *PostRepo) List(ctx context.Context, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Debug().Post.Query()

	excludeConditions, translationMaskFields, needQueryTranslation, err := r.prepareTranslationMaskFields(req)
	if err != nil {
		return nil, err
	}

	if len(excludeConditions) > 0 {
		var postIDs []uint32

		var tagIDs []uint32
		var categoryIDs []uint32
		for _, cond := range excludeConditions {
			//r.log.Debugf("exclude condition: %s %v", cond.GetField(), cond.GetValues())

			var val uint64
			switch cond.GetField() {
			case "tag_ids":
				for _, v := range cond.GetValues() {
					if val, err = strconv.ParseUint(v, 10, 64); err == nil {
						tagIDs = append(tagIDs, uint32(val))
					} else {
						r.log.Errorf("parse tag_ids value failed: %s", err.Error())
					}
				}

			case "category_ids":
				for _, v := range cond.GetValues() {
					if val, err = strconv.ParseUint(v, 10, 64); err == nil {
						categoryIDs = append(categoryIDs, uint32(val))
					} else {
						r.log.Errorf("parse category_ids value failed: %s", err.Error())
					}
				}
			}
		}

		if len(tagIDs) > 0 {
			var tempIDs []uint32
			tempIDs, err = r.postTagRepo.ListPostIDsByTagIDs(ctx, tagIDs)
			if err != nil {
				r.log.Errorf("query post IDs by tag IDs failed: %s", err.Error())
				return nil, contentV1.ErrorInternalServerError("query post IDs by tag IDs failed")
			}
			postIDs = append(postIDs, tempIDs...)
		}

		if len(categoryIDs) > 0 {
			var tempIDs []uint32
			tempIDs, err = r.postCategoryRepo.ListPostIDsByCategoryIDs(ctx, categoryIDs)
			if err != nil {
				r.log.Errorf("query post IDs by category IDs failed: %s", err.Error())
				return nil, contentV1.ErrorInternalServerError("query post IDs by category IDs failed")
			}
			postIDs = append(postIDs, tempIDs...)
		}

		builder.Modify(func(s *sql.Selector) {
			if len(tagIDs) > 0 {
				s.Join(
					sql.Table(posttag.Table).
						As("t1"),
				).
					On("\"t1\".post_id", "\"posts\".id")

				var ids = make([]any, 0, len(tagIDs))
				for _, id := range tagIDs {
					ids = append(ids, id)
				}
				s.Where(
					sql.In("\"t1\".tag_id", ids...),
				)
			}

			if len(categoryIDs) > 0 {
				s.Join(
					sql.Table(postcategory.Table).
						As("t1"),
				).
					On("\"t1\".post_id", "\"posts\".id")

				var ids = make([]any, 0, len(categoryIDs))
				for _, id := range categoryIDs {
					ids = append(ids, id)
				}
				s.Where(
					sql.In("\"t1\".category_id", ids...),
				)
			}
		})
	}

	if req.FieldMask != nil && len(req.FieldMask.Paths) > 0 {
		whereSelectors, err := r.repository.BuildSelectorWithTable(post.Table, req.FieldMask.Paths)
		if err != nil {
			r.log.Errorf("build selector with table failed: %s", err.Error())
			return nil, err
		}

		builder.Modify(whereSelectors)

		req.FieldMask = nil
	}

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &contentV1.ListPostResponse{Total: 0, Items: nil}, nil
	}

	if needQueryTranslation {
		viewMask := &fieldmaskpb.FieldMask{
			Paths: translationMaskFields,
		}

		var locale string
		if len(excludeConditions) > 0 {
			for _, cond := range excludeConditions {
				if cond.GetField() == "locale" {
					locale = cond.GetValue()
					break
				}
			}
		}

		for _, item := range ret.Items {
			translations, err := r.postTranslationRepo.ListTranslations(ctx, item.GetId(), locale, viewMask)
			if err != nil {
				r.log.Errorf("query translations failed: %s", err.Error())
				return nil, contentV1.ErrorInternalServerError("query translations failed")
			}
			item.Translations = translations
		}
	}

	for _, item := range ret.Items {
		languages, err := r.postTranslationRepo.ListAvailedLanguages(ctx, item.GetId())
		if err != nil {
			r.log.Errorf("query availed languages failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query availed languages failed")
		}
		item.AvailableLanguages = languages

		if tagIds, err := r.postTagRepo.ListTagIDs(ctx, item.GetId()); err != nil {
			r.log.Errorf("query tag ids failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query tag ids failed")
		} else {
			item.TagIds = tagIds
		}

		if categoryIds, err := r.postCategoryRepo.ListCategoryIDs(ctx, item.GetId()); err != nil {
			r.log.Errorf("query category ids failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query category ids failed")
		} else {
			item.CategoryIds = categoryIds
		}
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

	build := r.entClient.Client().Post.Query()

	switch req.QueryBy.(type) {
	case *contentV1.GetPostRequest_Id:
		build.Where(post.IDEQ(req.GetId()))

	case *contentV1.GetPostRequest_Code:
		build.Where(post.CodeEQ(req.GetCode()))

	default:
		return nil, contentV1.ErrorBadRequest("invalid query field")
	}

	entity, err := build.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, contentV1.ErrorFileNotFound("post not found")
		}

		r.log.Errorf("query post failed: %s", err.Error())

		return nil, contentV1.ErrorInternalServerError("query post failed")
	}

	dto := r.mapper.ToDTO(entity)

	languages, err := r.postTranslationRepo.ListAvailedLanguages(ctx, dto.GetId())
	if err != nil {
		r.log.Errorf("query availed languages failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query availed languages failed")
	}
	dto.AvailableLanguages = languages

	if req.Locale == nil {
		translations, err := r.postTranslationRepo.ListTranslations(ctx, dto.GetId(), "", nil)
		if err != nil {
			r.log.Errorf("query translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query translations failed")
		}
		dto.Translations = translations
	} else {
		translation, err := r.postTranslationRepo.GetTranslation(ctx, dto.GetId(), req.GetLocale())
		if err != nil {
			r.log.Errorf("query translation failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("query translation failed")
		}
		if translation != nil {
			dto.Translations = append(dto.Translations, translation)
		}
	}

	if tagIds, err := r.postTagRepo.ListTagIDs(ctx, dto.GetId()); err != nil {
		r.log.Errorf("query tag ids failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query tag ids failed")
	} else {
		dto.TagIds = tagIds
	}

	if categoryIds, err := r.postCategoryRepo.ListCategoryIDs(ctx, dto.GetId()); err != nil {
		r.log.Errorf("query category ids failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("query category ids failed")
	} else {
		dto.CategoryIds = categoryIds
	}

	return dto, nil
}

func (r *PostRepo) Create(ctx context.Context, req *contentV1.CreatePostRequest) (dto *contentV1.Post, err error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.Translations) == 0 {
		return nil, contentV1.ErrorBadRequest("at least one translation is required")
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
		SetNillableCode(req.Data.Code).
		SetNillableDisallowComment(req.Data.DisallowComment).
		SetNillableInProgress(req.Data.InProgress).
		SetNillableAutoSummary(req.Data.AutoSummary).
		SetNillableIsFeatured(req.Data.IsFeatured).
		SetNillableSortOrder(req.Data.SortOrder).
		// 服务端计数器初始化为 0，不接受客户端值
		SetVisits(0).
		SetLikes(0).
		SetCommentCount(0).
		SetNillableAuthorID(req.Data.AuthorId).
		SetNillableAuthorName(req.Data.AuthorName).
		SetNillablePasswordHash(req.Data.PasswordHash).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetNillablePublishTime(timeutil.TimestamppbToTime(req.Data.PublishTime)).
		SetCreatedAt(time.Now())

	if req.Data.CustomFields != nil {
		builder.SetCustomFields(trans.Ptr(req.Data.GetCustomFields()))
	}

	var entity *ent.Post
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert post failed: %s", err.Error())
		return nil, contentV1.ErrorInternalServerError("insert post failed")
	}

	if len(req.Data.Translations) > 0 {
		if err = r.postTranslationRepo.CleanTranslations(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean translations failed")
		}

		for i := range req.Data.Translations {
			req.Data.Translations[i].PostId = trans.Ptr(entity.ID)
		}

		if err = r.postTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	if len(req.Data.CategoryIds) > 0 {
		if err = r.postCategoryRepo.CleanCategories(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean categories failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean categories failed")
		}

		if err = r.postCategoryRepo.BatchCreate(ctx, tx, entity.ID, req.Data.GetCategoryIds()); err != nil {
			r.log.Errorf("batch insert categories failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert categories failed")
		}
	}

	if len(req.Data.TagIds) > 0 {
		if err = r.postTagRepo.CleanTags(ctx, tx, entity.ID); err != nil {
			r.log.Errorf("clean tags failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean tags failed")
		}

		if err = r.postTagRepo.BatchCreate(ctx, tx, entity.ID, req.Data.GetTagIds()); err != nil {
			r.log.Errorf("batch insert tags failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert tags failed")
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

	if len(req.Data.Translations) > 0 {
		for i := range req.Data.Translations {
			req.Data.Translations[i].PostId = trans.Ptr(req.GetId())
		}

		if err = r.postTranslationRepo.BatchCreate(ctx, tx, req.Data.GetTranslations()); err != nil {
			r.log.Errorf("batch insert translations failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert translations failed")
		}
	}

	if req.Data.CategoryIds != nil {
		if err = r.postCategoryRepo.CleanCategories(ctx, tx, req.GetId()); err != nil {
			r.log.Errorf("clean categories failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean categories failed")
		}

		if err = r.postCategoryRepo.BatchCreate(ctx, tx, req.GetId(), req.Data.GetCategoryIds()); err != nil {
			r.log.Errorf("batch insert categories failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert categories failed")
		}
	}

	if req.Data.TagIds != nil {
		if err = r.postTagRepo.CleanTags(ctx, tx, req.GetId()); err != nil {
			r.log.Errorf("clean tags failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("clean tags failed")
		}

		if err = r.postTagRepo.BatchCreate(ctx, tx, req.GetId(), req.Data.GetTagIds()); err != nil {
			r.log.Errorf("batch insert tags failed: %s", err.Error())
			return nil, contentV1.ErrorInternalServerError("batch insert tags failed")
		}
	}

	builder := tx.Post.UpdateOneID(req.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *contentV1.Post) {
			// 服务端计数器（visits/likes/comment_count）不可由客户端通过 Update 设置，
			// 仅能通过服务端专用的递增接口修改。
			builder.
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableEditorType(r.editorTypeConverter.ToEntity(req.Data.EditorType)).
				SetNillableCode(req.Data.Code).
				SetNillableDisallowComment(req.Data.DisallowComment).
				SetNillableInProgress(req.Data.InProgress).
				SetNillableAutoSummary(req.Data.AutoSummary).
				SetNillableIsFeatured(req.Data.IsFeatured).
				SetNillableSortOrder(req.Data.SortOrder).
				SetNillableAuthorID(req.Data.AuthorId).
				SetNillableAuthorName(req.Data.AuthorName).
				SetNillablePasswordHash(req.Data.PasswordHash).
				SetNillableUpdatedBy(req.Data.UpdatedBy).
				SetNillablePublishTime(timeutil.TimestamppbToTime(req.Data.PublishTime)).
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

	// 删除帖子数据
	if err = r.entClient.Client().Post.
		DeleteOneID(req.GetId()).
		Exec(ctx); err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	// 删除关联数据
	if err = r.postTranslationRepo.CleanTranslations(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean translations failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean translations failed")
	}

	// 删除关联数据
	if err = r.postCategoryRepo.CleanCategories(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean categories failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean categories failed")
	}

	// 删除关联数据
	if err = r.postTagRepo.CleanTags(ctx, tx, req.GetId()); err != nil {
		r.log.Errorf("clean tags failed: %s", err.Error())
		return contentV1.ErrorInternalServerError("clean tags failed")
	}

	return err
}

func (r *PostRepo) TranslationExists(ctx context.Context, postId uint32, languageCode string) (bool, error) {
	return r.postTranslationRepo.TranslationExists(ctx, postId, languageCode)
}

func (r *PostRepo) CreateTranslation(ctx context.Context, req *contentV1.CreatePostTranslationRequest) (*contentV1.PostTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	req.Data.PostId = trans.Ptr(req.GetPostId())

	return r.postTranslationRepo.CreateTranslation(ctx, req.Data)
}

func (r *PostRepo) UpdateTranslation(ctx context.Context, req *contentV1.UpdatePostTranslationRequest) (*contentV1.PostTranslation, error) {
	if req == nil || req.Data == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	if len(req.Data.GetLanguageCode()) == 0 {
		return nil, contentV1.ErrorBadRequest("language code is required")
	}

	if req.Data.GetPostId() == 0 {
		return nil, contentV1.ErrorBadRequest("post id is required")
	}

	if exist, err := r.TranslationExists(ctx, req.Data.GetPostId(), req.Data.GetLanguageCode()); err != nil {
		return nil, err
	} else if !exist {
		if req.GetAllowMissing() {
			return r.CreateTranslation(ctx, &contentV1.CreatePostTranslationRequest{
				Data:   req.Data,
				PostId: req.Data.GetPostId(),
			})
		}

		return nil, contentV1.ErrorFileNotFound("translation not found")
	}

	return r.postTranslationRepo.UpdateTranslation(ctx, req.GetId(), req.Data, req.GetUpdateMask())
}

func (r *PostRepo) GetTranslation(ctx context.Context, req *contentV1.GetPostRequest) (*contentV1.PostTranslation, error) {
	if req == nil {
		return nil, contentV1.ErrorBadRequest("invalid parameter")
	}

	return r.postTranslationRepo.GetTranslation(ctx, req.GetId(), req.GetLocale())
}

func (r *PostRepo) ListTranslations(ctx context.Context, postId uint32) ([]*contentV1.PostTranslation, error) {
	return r.postTranslationRepo.ListTranslations(ctx, postId, "", nil)
}

func (r *PostRepo) DeleteTranslation(ctx context.Context, req *contentV1.DeletePostTranslationRequest) error {
	return r.postTranslationRepo.DeleteTranslation(ctx, req)
}

func (r *PostRepo) CleanTranslations(ctx context.Context, tx *ent.Tx, postID uint32) error {
	return r.postTranslationRepo.CleanTranslations(ctx, tx, postID)
}

func (r *PostRepo) CleanCategories(ctx context.Context, tx *ent.Tx, postID uint32) error {
	return r.postCategoryRepo.CleanCategories(ctx, tx, postID)
}
