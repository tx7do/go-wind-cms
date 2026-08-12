package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/mapper"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/comment"
	"go-wind-cms/app/core/service/internal/data/ent/predicate"

	commentV1 "go-wind-cms/api/gen/go/comment/service/v1"

	"go-wind-cms/pkg/utils"
)

type CommentRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper

	mapper *mapper.CopierMapper[commentV1.Comment, ent.Comment]

	repository *entCrud.Repository[
		ent.CommentQuery, ent.CommentSelect,
		ent.CommentCreate, ent.CommentCreateBulk,
		ent.CommentUpdate, ent.CommentUpdateOne,
		ent.CommentDelete,
		predicate.Comment,
		commentV1.Comment, ent.Comment,
	]

	statusConverter      *mapper.EnumTypeConverter[commentV1.Comment_Status, comment.Status]
	contentTypeConverter *mapper.EnumTypeConverter[commentV1.Comment_ContentType, comment.ContentType]
	authorTypeConverter  *mapper.EnumTypeConverter[commentV1.Comment_AuthorType, comment.AuthorType]
}

func NewCommentRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *CommentRepo {
	repo := &CommentRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("comment/repo/core-service"),
		mapper:    mapper.NewCopierMapper[commentV1.Comment, ent.Comment](),
		statusConverter: mapper.NewEnumTypeConverter[commentV1.Comment_Status, comment.Status](
			commentV1.Comment_Status_name, commentV1.Comment_Status_value,
		),
		contentTypeConverter: mapper.NewEnumTypeConverter[commentV1.Comment_ContentType, comment.ContentType](
			commentV1.Comment_ContentType_name, commentV1.Comment_ContentType_value,
		),
		authorTypeConverter: mapper.NewEnumTypeConverter[commentV1.Comment_AuthorType, comment.AuthorType](
			commentV1.Comment_AuthorType_name, commentV1.Comment_AuthorType_value,
		),
	}

	repo.init()

	return repo
}

func (r *CommentRepo) init() {
	r.repository = entCrud.NewRepository[
		ent.CommentQuery, ent.CommentSelect,
		ent.CommentCreate, ent.CommentCreateBulk,
		ent.CommentUpdate, ent.CommentUpdateOne,
		ent.CommentDelete,
		predicate.Comment,
		commentV1.Comment, ent.Comment,
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())

	r.mapper.AppendConverters(r.statusConverter.NewConverterPair())
	r.mapper.AppendConverters(r.contentTypeConverter.NewConverterPair())
	r.mapper.AppendConverters(r.authorTypeConverter.NewConverterPair())
}

func (r *CommentRepo) count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.entClient.Client().Comment.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
	}

	return count, err
}

func (r *CommentRepo) Count(ctx context.Context, req *paginationV1.PagingRequest) (int, error) {
	builder := r.entClient.Client().Comment.Query()

	whereSelectors, _, err := r.repository.BuildListSelectorWithPaging(builder, req)
	if len(whereSelectors) != 0 {
		builder.Modify(whereSelectors...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
		return 0, commentV1.ErrorInternalServerError("query count failed")
	}

	return count, nil
}

func (r *CommentRepo) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.entClient.Client().Comment.Query().
		Where(comment.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query comment exist failed: %s", err.Error())
		return false, commentV1.ErrorInternalServerError("query comment exist failed")
	}
	return exist, nil
}

func (r *CommentRepo) buildCommentTree(items []*commentV1.Comment, parentId uint32) []*commentV1.Comment {
	var tree []*commentV1.Comment
	for _, item := range items {
		if item.GetParentId() == parentId {
			// 递归查找子节点
			children := r.buildCommentTree(items, item.GetId())
			item.Children = children
			tree = append(tree, item)
		}
	}
	return tree
}

func (r *CommentRepo) List(ctx context.Context, req *paginationV1.PagingRequest, treeTravel bool) (*commentV1.ListCommentResponse, error) {
	if req == nil {
		return nil, commentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Comment.Query()

	ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return &commentV1.ListCommentResponse{Total: 0, Items: nil}, nil
	}

	if treeTravel {
		roots := r.buildCommentTree(ret.Items, 0) // 假设根节点ParentId为0
		return &commentV1.ListCommentResponse{
			Total: ret.Total,
			Items: roots,
		}, nil
	}

	return &commentV1.ListCommentResponse{
		Total: ret.Total,
		Items: ret.Items,
	}, nil
}

func (r *CommentRepo) Get(ctx context.Context, req *commentV1.GetCommentRequest) (*commentV1.Comment, error) {
	if req == nil {
		return nil, commentV1.ErrorBadRequest("invalid parameter")
	}

	entity, err := r.entClient.Client().Comment.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, commentV1.ErrorFileNotFound("comment not found")
		}

		r.log.Errorf("query comment failed: %s", err.Error())

		return nil, commentV1.ErrorInternalServerError("query comment failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *CommentRepo) Create(ctx context.Context, req *commentV1.CreateCommentRequest) (*commentV1.Comment, error) {
	if req == nil || req.Data == nil {
		return nil, commentV1.ErrorBadRequest("invalid parameter")
	}

	builder := r.entClient.Client().Comment.Create().
		SetNillableContentType(r.contentTypeConverter.ToEntity(req.Data.ContentType)).
		SetNillableObjectID(req.Data.ObjectId).
		SetNillableContent(req.Data.Content).
		SetNillableAuthorID(req.Data.AuthorId).
		SetNillableAuthorName(req.Data.AuthorName).
		SetNillableAuthorEmail(req.Data.AuthorEmail).
		SetNillableAuthorURL(req.Data.AuthorUrl).
		SetNillableAuthorType(r.authorTypeConverter.ToEntity(req.Data.AuthorType)).
		SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
		// 计数列（like_count/dislike_count/reply_count）由 schema Default(0) 兜底，
		// 由 InteractionService 独占递增，Create 路径不再显式设置。
		SetNillableIPAddress(req.Data.IpAddress).
		SetNillableLocation(req.Data.Location).
		SetNillableUserAgent(req.Data.UserAgent).
		SetNillableDetectedLanguage(req.Data.DetectedLanguage).
		SetNillableIsSpam(req.Data.IsSpam).
		SetNillableIsSticky(req.Data.IsSticky).
		SetNillableParentID(req.Data.ParentId).
		SetNillableCreatedBy(req.Data.CreatedBy).
		SetCreatedAt(time.Now())

	var err error
	var entity *ent.Comment
	if entity, err = builder.Save(ctx); err != nil {
		r.log.Errorf("insert comment failed: %s", err.Error())
		return nil, commentV1.ErrorInternalServerError("insert comment failed")
	}

	return r.mapper.ToDTO(entity), nil
}

func (r *CommentRepo) Update(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.Comment, error) {
	if req == nil || req.Data == nil {
		return nil, commentV1.ErrorBadRequest("invalid parameter")
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
			_, err = r.Create(ctx, &commentV1.CreateCommentRequest{Data: req.Data})
			return nil, err
		}
	}

	tid, hasTenant := maybeTenantFromViewer(ctx)
	callerUserID, hasUser := viewerUserIDFromContext(ctx)
	// 计数列（like_count/dislike_count/reply_count）由 InteractionService 独占写入。
	// 经 FilterBlacklist 剥离客户端通过 update_mask 写入计数列的尝试，同时关闭
	// go-crud applyUpdateOneNilFieldMask 在“mask 命中且未传值”时下发 SET NULL 的通道，
	// 避免客户端把计数列 NULL 化造成数据丢失。
	req.GetUpdateMask().Paths = utils.FilterBlacklist(req.GetUpdateMask().Paths, []string{
		"like_count",
		"dislike_count",
		"reply_count",
	})
	builder := r.entClient.Client().Comment.UpdateOneID(req.GetId())
	builder.Where(comment.IDEQ(req.GetId()))
	if hasTenant {
		builder.Where(comment.TenantIDEQ(tid))
	}
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *commentV1.Comment) {
			builder.
				SetNillableContentType(r.contentTypeConverter.ToEntity(req.Data.ContentType)).
				SetNillableObjectID(req.Data.ObjectId).
				SetNillableContent(req.Data.Content).
				SetNillableAuthorID(req.Data.AuthorId).
				SetNillableAuthorName(req.Data.AuthorName).
				SetNillableAuthorEmail(req.Data.AuthorEmail).
				SetNillableAuthorURL(req.Data.AuthorUrl).
				SetNillableAuthorType(r.authorTypeConverter.ToEntity(req.Data.AuthorType)).
				SetNillableStatus(r.statusConverter.ToEntity(req.Data.Status)).
				SetNillableIPAddress(req.Data.IpAddress).
				SetNillableLocation(req.Data.Location).
				SetNillableUserAgent(req.Data.UserAgent).
				SetNillableDetectedLanguage(req.Data.DetectedLanguage).
				SetNillableIsSpam(req.Data.IsSpam).
				SetNillableIsSticky(req.Data.IsSticky).
				SetNillableParentID(req.Data.ParentId).
				SetUpdatedAt(time.Now())

			// updated_by 强制由服务端 viewer context 推导，忽略客户端传入值
			if hasUser {
				builder.SetUpdatedBy(callerUserID)
			}
		},
		func(s *sql.Selector) {
			s.Where(sql.EQ(comment.FieldID, req.GetId()))
		},
	)

	return result, err
}

func (r *CommentRepo) Delete(ctx context.Context, req *commentV1.DeleteCommentRequest) error {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	delBuilder := r.entClient.Client().Comment.Delete()
	delBuilder.Where(comment.IDEQ(req.GetId()))
	if hasTenant {
		delBuilder.Where(comment.TenantIDEQ(tid))
	}
	_, err := delBuilder.Exec(ctx)
	if err != nil {
		r.log.Errorf("delete one data failed: %s", err.Error())
	}

	return err
}
