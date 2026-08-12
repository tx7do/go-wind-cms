package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	paginationV1 "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/comment"
	"go-wind-cms/app/core/service/internal/data/ent/commentlike"
	"go-wind-cms/app/core/service/internal/data/ent/post"
	"go-wind-cms/app/core/service/internal/data/ent/postlike"
	"go-wind-cms/app/core/service/internal/data/ent/postwatch"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
	interactionV1 "go-wind-cms/api/gen/go/interaction/service/v1"
)

// InteractionRepo 是点赞/收藏 ledger 与计数缓存的唯一写入方。
//
// 设计要点：
//   - viewer 用户身份由 service 层从鉴权上下文提取后传入，repo 不接受客户端
//     传入的 user_id。
//   - Like/Unlike 在单个 ent.Tx 内同时操作 ledger 表与 post/comment 计数缓存，
//     保证原子性（沿用 post_repo 的跨表事务惯例）。
//   - 重复点赞/收藏走幂等：先 Exist 检查，已存在则 no-op，不依赖约束错误判断。
//   - Watch/Unwatch 仅维护 post_watch ledger，不触碰任何计数列（收藏无缓存列）。
//   - 跨租户隔离由 mixin + 查询带 tenant_id 保证。
type InteractionRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewInteractionRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *InteractionRepo {
	return &InteractionRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("interaction/repo/core-service"),
	}
}

// txn 在事务内执行 fn，沿用 post_repo 的 defer rollback/commit 惯例。
func (r *InteractionRepo) txn(ctx context.Context, fn func(tx *ent.Tx) error) (err error) {
	var tx *ent.Tx
	tx, err = r.entClient.Client().Tx(ctx)
	if err != nil {
		r.log.Errorf("start transaction failed: %s", err.Error())
		return interactionV1.ErrorInternalServerError("start transaction failed")
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
			err = interactionV1.ErrorInternalServerError("transaction commit failed")
		}
	}()
	return fn(tx)
}

// Like 点赞 post 或 comment。幂等：已点赞则 no-op。
// 返回操作后的点赞状态与最新计数。
func (r *InteractionRepo) Like(ctx context.Context, viewerUserID uint32, targetType interactionV1.TargetType, targetID uint32) (bool, int32, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, 0, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	switch targetType {
	case interactionV1.TargetType_TARGET_TYPE_POST:
		likeCount, liked, err := r.likePost(ctx, tid, viewerUserID, targetID)
		if err != nil {
			return false, 0, err
		}
		return liked, likeCount, nil

	case interactionV1.TargetType_TARGET_TYPE_COMMENT:
		likeCount, liked, err := r.likeComment(ctx, tid, viewerUserID, targetID)
		if err != nil {
			return false, 0, err
		}
		return liked, likeCount, nil

	default:
		return false, 0, interactionV1.ErrorBadRequest("invalid target type")
	}
}

// Unlike 取消点赞。幂等：未点赞则 no-op。
func (r *InteractionRepo) Unlike(ctx context.Context, viewerUserID uint32, targetType interactionV1.TargetType, targetID uint32) (bool, int32, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, 0, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	switch targetType {
	case interactionV1.TargetType_TARGET_TYPE_POST:
		likeCount, liked, err := r.unlikePost(ctx, tid, viewerUserID, targetID)
		if err != nil {
			return false, 0, err
		}
		return liked, likeCount, nil

	case interactionV1.TargetType_TARGET_TYPE_COMMENT:
		likeCount, liked, err := r.unlikeComment(ctx, tid, viewerUserID, targetID)
		if err != nil {
			return false, 0, err
		}
		return liked, likeCount, nil

	default:
		return false, 0, interactionV1.ErrorBadRequest("invalid target type")
	}
}

// readPostLikeCount 读取 post.likes 当前值（int32）。
func (r *InteractionRepo) readPostLikeCount(ctx context.Context, postID uint32) (int32, error) {
	entity, err := r.entClient.Client().Post.Query().
		Where(post.IDEQ(postID)).
		Only(ctx)
	if err != nil {
		r.log.Errorf("query post likes failed: %s", err.Error())
		return 0, interactionV1.ErrorInternalServerError("query post likes failed")
	}
	if entity.Likes == nil {
		return 0, nil
	}
	return *entity.Likes, nil
}

// readCommentLikeCount 读取 comment.like_count 当前值（uint32 → int32）。
func (r *InteractionRepo) readCommentLikeCount(ctx context.Context, commentID uint32) (int32, error) {
	entity, err := r.entClient.Client().Comment.Query().
		Where(comment.IDEQ(commentID)).
		Only(ctx)
	if err != nil {
		r.log.Errorf("query comment like_count failed: %s", err.Error())
		return 0, interactionV1.ErrorInternalServerError("query comment like_count failed")
	}
	if entity.LikeCount == nil {
		return 0, nil
	}
	return int32(*entity.LikeCount), nil
}

// likePost 在单个 tx 内：写 post_like ledger + 递增 post.likes。
// 幂等：若 ledger 行已存在则 no-op（不重复递增）。
func (r *InteractionRepo) likePost(ctx context.Context, tid, viewerUserID, postID uint32) (int32, bool, error) {
	exists, err := r.entClient.Client().PostLike.Query().
		Where(
			postlike.TenantIDEQ(tid),
			postlike.UserIDEQ(viewerUserID),
			postlike.PostIDEQ(postID),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query post_like exist failed: %s", err.Error())
		return 0, false, interactionV1.ErrorInternalServerError("query post_like exist failed")
	}

	if exists {
		// 幂等：已点赞，no-op，但仍返回当前状态
		likeCount, qerr := r.readPostLikeCount(ctx, postID)
		if qerr != nil {
			return 0, false, qerr
		}
		return likeCount, true, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		if _, err := tx.PostLike.Create().
			SetTenantID(tid).
			SetUserID(viewerUserID).
			SetPostID(postID).
			Save(ctx); err != nil {
			r.log.Errorf("insert post_like failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("insert post_like failed")
		}

		// 递增 post.likes 缓存
		if _, err := tx.Post.UpdateOneID(postID).
			AddLikes(1).
			Save(ctx); err != nil {
			r.log.Errorf("increment post likes failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("increment post likes failed")
		}

		return nil
	})
	if txErr != nil {
		return 0, false, txErr
	}

	likeCount, err := r.readPostLikeCount(ctx, postID)
	if err != nil {
		return 0, false, err
	}
	return likeCount, true, nil
}

// unlikePost 在单个 tx 内：删 post_like ledger + 递减 post.likes。
// 幂等：若 ledger 行不存在则 no-op。
func (r *InteractionRepo) unlikePost(ctx context.Context, tid, viewerUserID, postID uint32) (int32, bool, error) {
	exists, err := r.entClient.Client().PostLike.Query().
		Where(
			postlike.TenantIDEQ(tid),
			postlike.UserIDEQ(viewerUserID),
			postlike.PostIDEQ(postID),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query post_like exist failed: %s", err.Error())
		return 0, false, interactionV1.ErrorInternalServerError("query post_like exist failed")
	}

	if !exists {
		// 幂等：未点赞，no-op
		likeCount, qerr := r.readPostLikeCount(ctx, postID)
		if qerr != nil {
			return 0, false, qerr
		}
		return likeCount, false, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		affected, err := tx.PostLike.Delete().
			Where(
				postlike.TenantIDEQ(tid),
				postlike.UserIDEQ(viewerUserID),
				postlike.PostIDEQ(postID),
			).
			Exec(ctx)
		if err != nil {
			r.log.Errorf("delete post_like failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("delete post_like failed")
		}
		if affected == 0 {
			// 并发删除，幂等处理
			return nil
		}

		// 递减 post.likes 缓存
		if _, err := tx.Post.UpdateOneID(postID).
			AddLikes(-1).
			Save(ctx); err != nil {
			r.log.Errorf("decrement post likes failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("decrement post likes failed")
		}
		return nil
	})
	if txErr != nil {
		return 0, false, txErr
	}

	likeCount, err := r.readPostLikeCount(ctx, postID)
	if err != nil {
		return 0, false, err
	}
	return likeCount, false, nil
}

// likeComment 在单个 tx 内：写 comment_like ledger + 递增 comment.like_count。
func (r *InteractionRepo) likeComment(ctx context.Context, tid, viewerUserID, commentID uint32) (int32, bool, error) {
	exists, err := r.entClient.Client().CommentLike.Query().
		Where(
			commentlike.TenantIDEQ(tid),
			commentlike.UserIDEQ(viewerUserID),
			commentlike.CommentIDEQ(commentID),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query comment_like exist failed: %s", err.Error())
		return 0, false, interactionV1.ErrorInternalServerError("query comment_like exist failed")
	}

	if exists {
		// 幂等：已点赞，no-op
		likeCount, qerr := r.readCommentLikeCount(ctx, commentID)
		if qerr != nil {
			return 0, false, qerr
		}
		return likeCount, true, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		if _, err := tx.CommentLike.Create().
			SetTenantID(tid).
			SetUserID(viewerUserID).
			SetCommentID(commentID).
			Save(ctx); err != nil {
			r.log.Errorf("insert comment_like failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("insert comment_like failed")
		}

		if _, err := tx.Comment.UpdateOneID(commentID).
			AddLikeCount(1).
			Save(ctx); err != nil {
			r.log.Errorf("increment comment like_count failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("increment comment like_count failed")
		}
		return nil
	})
	if txErr != nil {
		return 0, false, txErr
	}

	likeCount, err := r.readCommentLikeCount(ctx, commentID)
	if err != nil {
		return 0, false, err
	}
	return likeCount, true, nil
}

// unlikeComment 在单个 tx 内：删 comment_like ledger + 递减 comment.like_count。
func (r *InteractionRepo) unlikeComment(ctx context.Context, tid, viewerUserID, commentID uint32) (int32, bool, error) {
	exists, err := r.entClient.Client().CommentLike.Query().
		Where(
			commentlike.TenantIDEQ(tid),
			commentlike.UserIDEQ(viewerUserID),
			commentlike.CommentIDEQ(commentID),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query comment_like exist failed: %s", err.Error())
		return 0, false, interactionV1.ErrorInternalServerError("query comment_like exist failed")
	}

	if !exists {
		// 幂等：未点赞，no-op
		likeCount, qerr := r.readCommentLikeCount(ctx, commentID)
		if qerr != nil {
			return 0, false, qerr
		}
		return likeCount, false, nil
	}

	txErr := r.txn(ctx, func(tx *ent.Tx) error {
		affected, err := tx.CommentLike.Delete().
			Where(
				commentlike.TenantIDEQ(tid),
				commentlike.UserIDEQ(viewerUserID),
				commentlike.CommentIDEQ(commentID),
			).
			Exec(ctx)
		if err != nil {
			r.log.Errorf("delete comment_like failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("delete comment_like failed")
		}
		if affected == 0 {
			return nil
		}

		if _, err := tx.Comment.UpdateOneID(commentID).
			AddLikeCount(-1).
			Save(ctx); err != nil {
			r.log.Errorf("decrement comment like_count failed: %s", err.Error())
			return interactionV1.ErrorInternalServerError("decrement comment like_count failed")
		}
		return nil
	})
	if txErr != nil {
		return 0, false, txErr
	}

	likeCount, err := r.readCommentLikeCount(ctx, commentID)
	if err != nil {
		return 0, false, err
	}
	return likeCount, false, nil
}

// Watch 收藏 post。幂等：已收藏则 no-op。不触碰任何计数列。
func (r *InteractionRepo) Watch(ctx context.Context, viewerUserID uint32, postID uint32) (bool, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	exists, err := r.entClient.Client().PostWatch.Query().
		Where(
			postwatch.TenantIDEQ(tid),
			postwatch.UserIDEQ(viewerUserID),
			postwatch.PostIDEQ(postID),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query post_watch exist failed: %s", err.Error())
		return false, interactionV1.ErrorInternalServerError("query post_watch exist failed")
	}
	if exists {
		// 幂等：已收藏，no-op
		return true, nil
	}

	_, err = r.entClient.Client().PostWatch.Create().
		SetTenantID(tid).
		SetUserID(viewerUserID).
		SetPostID(postID).
		Save(ctx)
	if err != nil {
		r.log.Errorf("insert post_watch failed: %s", err.Error())
		return false, interactionV1.ErrorInternalServerError("insert post_watch failed")
	}
	return true, nil
}

// Unwatch 取消收藏 post。幂等。
func (r *InteractionRepo) Unwatch(ctx context.Context, viewerUserID uint32, postID uint32) (bool, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return false, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	affected, err := r.entClient.Client().PostWatch.Delete().
		Where(
			postwatch.TenantIDEQ(tid),
			postwatch.UserIDEQ(viewerUserID),
			postwatch.PostIDEQ(postID),
		).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("delete post_watch failed: %s", err.Error())
		return false, interactionV1.ErrorInternalServerError("delete post_watch failed")
	}
	// affected==0 视为未收藏的幂等 no-op
	return affected > 0, nil
}

// GetInteractionStatus 批量查询当前 viewer 对指定目标的 {liked, watched} 状态。
// comment 目标的 watched 永远为 false（收藏仅对 post）。
func (r *InteractionRepo) GetInteractionStatus(ctx context.Context, viewerUserID uint32, targetType interactionV1.TargetType, targetIDs []uint32) (map[uint32]*interactionV1.InteractionStatus, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return nil, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	result := make(map[uint32]*interactionV1.InteractionStatus, len(targetIDs))
	for _, id := range targetIDs {
		result[id] = &interactionV1.InteractionStatus{Liked: false, Watched: false}
	}

	switch targetType {
	case interactionV1.TargetType_TARGET_TYPE_POST:
		// 查 post_like
		likedRows, err := r.entClient.Client().PostLike.Query().
			Where(
				postlike.TenantIDEQ(tid),
				postlike.UserIDEQ(viewerUserID),
				postlike.PostIDIn(targetIDs...),
			).
			All(ctx)
		if err != nil {
			r.log.Errorf("query post_like status failed: %s", err.Error())
			return nil, interactionV1.ErrorInternalServerError("query post_like status failed")
		}
		for _, row := range likedRows {
			if row.PostID != nil {
				if _, ok := result[*row.PostID]; ok {
					result[*row.PostID].Liked = true
				}
			}
		}

		// 查 post_watch
		watchedRows, err := r.entClient.Client().PostWatch.Query().
			Where(
				postwatch.TenantIDEQ(tid),
				postwatch.UserIDEQ(viewerUserID),
				postwatch.PostIDIn(targetIDs...),
			).
			All(ctx)
		if err != nil {
			r.log.Errorf("query post_watch status failed: %s", err.Error())
			return nil, interactionV1.ErrorInternalServerError("query post_watch status failed")
		}
		for _, row := range watchedRows {
			if row.PostID != nil {
				if _, ok := result[*row.PostID]; ok {
					result[*row.PostID].Watched = true
				}
			}
		}

	case interactionV1.TargetType_TARGET_TYPE_COMMENT:
		likedRows, err := r.entClient.Client().CommentLike.Query().
			Where(
				commentlike.TenantIDEQ(tid),
				commentlike.UserIDEQ(viewerUserID),
				commentlike.CommentIDIn(targetIDs...),
			).
			All(ctx)
		if err != nil {
			r.log.Errorf("query comment_like status failed: %s", err.Error())
			return nil, interactionV1.ErrorInternalServerError("query comment_like status failed")
		}
		for _, row := range likedRows {
			if row.CommentID != nil {
				if _, ok := result[*row.CommentID]; ok {
					result[*row.CommentID].Liked = true
				}
			}
		}
		// watched 对 comment 永远 false，已在初始化时设好

	default:
		return nil, interactionV1.ErrorBadRequest("invalid target type")
	}

	return result, nil
}

// ListWatchedPosts 列出当前 viewer 收藏的 post。
// 查 post_watch 拿到分页后的 post_ids，再逐个调 PostRepo.Get 复用其完整的
// 附带查询（translations/tags/categories/view_mask）逻辑。
func (r *InteractionRepo) ListWatchedPosts(ctx context.Context, viewerUserID uint32, postRepo *PostRepo, req *paginationV1.PagingRequest) (*contentV1.ListPostResponse, error) {
	tid, hasTenant := maybeTenantFromViewer(ctx)
	if !hasTenant || viewerUserID == 0 {
		return nil, interactionV1.ErrorUnauthorized("viewer identity required")
	}

	// 注：post_watch 的 (tenant_id, user_id, post_id) unique 索引兼作查询索引，
	// 此处按 (tenant, user) 过滤，分页由 PostWatchQuery 的 Limit/Offset 处理。
	q := r.entClient.Client().PostWatch.Query().
		Where(
			postwatch.TenantIDEQ(tid),
			postwatch.UserIDEQ(viewerUserID),
		)

	// 应用分页（Page 为 1-based）
	page := req.GetPage()
	if page == 0 {
		page = 1
	}
	limit := int(req.GetPageSize())
	if limit <= 0 {
		limit = 20
	}
	offset := int(page-1) * limit
	q = q.Limit(limit).Offset(offset)

	watchRows, err := q.All(ctx)
	if err != nil {
		r.log.Errorf("query post_watch list failed: %s", err.Error())
		return nil, interactionV1.ErrorInternalServerError("query post_watch list failed")
	}

	items := make([]*contentV1.Post, 0, len(watchRows))
	for _, row := range watchRows {
		if row.PostID == nil {
			continue
		}
		// 复用 PostRepo.Get 的附带查询。收藏列表请求不携带 locale/view_mask，
		// 故此处 GetPostRequest 不设置二者（默认全字段、默认语言）。
		postEntity, gerr := postRepo.Get(ctx, &contentV1.GetPostRequest{
			QueryBy: &contentV1.GetPostRequest_Id{
				Id: *row.PostID,
			},
		})
		if gerr != nil {
			r.log.Errorf("fetch watched post %d failed: %s", *row.PostID, gerr.Error())
			continue
		}
		items = append(items, postEntity)
	}

	// total 为该 viewer 收藏总数（用于前端分页元数据）
	total, cerr := r.entClient.Client().PostWatch.Query().
		Where(
			postwatch.TenantIDEQ(tid),
			postwatch.UserIDEQ(viewerUserID),
		).
		Count(ctx)
	if cerr != nil {
		r.log.Errorf("count post_watch failed: %s", cerr.Error())
		return nil, interactionV1.ErrorInternalServerError("count post_watch failed")
	}

	return &contentV1.ListPostResponse{
		Items: items,
		Total: uint64(total),
	}, nil
}
