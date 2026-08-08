package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/posttag"
)

type PostTagRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewPostTagRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *PostTagRepo {
	repo := &PostTagRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("post-tag/repo/core-service"),
	}

	repo.init()

	return repo
}

func (r *PostTagRepo) init() {
}

// Create 创建帖子标签关联关系，通常在创建帖子时调用
func (r *PostTagRepo) Create(ctx context.Context, postID, tagID uint32) error {
	if postID == 0 || tagID == 0 {
		return nil
	}

	_, err := r.entClient.Client().PostTag.Create().
		SetPostID(postID).
		SetTagID(tagID).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		r.log.Errorf("failed to create post-tag: %v", err)
		return err
	}
	return nil
}

// BatchCreate 批量创建帖子标签关联关系，通常在创建或更新帖子时调用
func (r *PostTagRepo) BatchCreate(ctx context.Context, tx *ent.Tx, postID uint32, tagIDs []uint32) error {
	if postID == 0 || len(tagIDs) == 0 {
		return nil
	}

	builders := make([]*ent.PostTagCreate, 0, len(tagIDs))
	for _, id := range tagIDs {
		builders = append(builders, tx.PostTag.Create().
			SetPostID(postID).
			SetTagID(id).
			SetCreatedAt(time.Now()),
		)
	}

	_, err := tx.PostTag.CreateBulk(builders...).Save(ctx)
	if err != nil {
		r.log.Errorf("failed to batch create post-tags: %v", err)
		return err
	}
	return nil
}

// CleanTags 删除指定帖子关联的所有帖子标签关联关系，通常在更新帖子标签时调用
func (r *PostTagRepo) CleanTags(ctx context.Context, tx *ent.Tx, postID uint32) error {
	delBuilder := tx.PostTag.Delete().Where(posttag.PostIDEQ(postID))
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(posttag.TenantIDEQ(tid))
	}
	if _, err := delBuilder.Exec(ctx); err != nil {
		r.log.Errorf("delete old post [%d] tags failed: %s", postID, err.Error())
		return err
	}
	return nil
}

// CleanPosts 删除指定标签关联的所有帖子标签关联关系，通常在删除标签时调用
func (r *PostTagRepo) CleanPosts(ctx context.Context, tx *ent.Tx, tagID uint32) error {
	delBuilder := tx.PostTag.Delete().Where(posttag.TagIDEQ(tagID))
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(posttag.TenantIDEQ(tid))
	}
	if _, err := delBuilder.Exec(ctx); err != nil {
		r.log.Errorf("delete old tag [%d] posts failed: %s", tagID, err.Error())
		return err
	}
	return nil
}

// ListTagIDs 列出指定帖子关联的所有标签ID，通常在查询帖子详情时调用
func (r *PostTagRepo) ListTagIDs(ctx context.Context, postID uint32) ([]uint32, error) {
	if postID == 0 {
		return nil, nil
	}

	tags, err := r.entClient.Client().PostTag.Query().
		Where(
			posttag.PostIDEQ(postID),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-tag by post ID: %v", err)
		return nil, err
	}

	tagIDs := make([]uint32, 0, len(tags))
	for _, tag := range tags {
		if tag == nil || tag.TagID == nil {
			continue
		}

		tagIDs = append(tagIDs, *tag.TagID)
	}
	return tagIDs, nil
}

func (r *PostTagRepo) ListTagIDsByPostIDs(ctx context.Context, postIDs []uint32) ([]uint32, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	tags, err := r.entClient.Client().PostTag.Query().
		Where(
			posttag.PostIDIn(postIDs...),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-tag by post IDs: %v", err)
		return nil, err
	}

	tagIDs := make([]uint32, 0, len(tags))
	for _, tag := range tags {
		if tag == nil || tag.TagID == nil {
			continue
		}

		tagIDs = append(tagIDs, *tag.TagID)
	}
	return tagIDs, nil
}

// ListPostIDs 列出指定标签关联的所有帖子ID，通常在查询标签详情时调用
func (r *PostTagRepo) ListPostIDs(ctx context.Context, tagID uint32) ([]uint32, error) {
	if tagID == 0 {
		return nil, nil
	}

	posts, err := r.entClient.Client().PostTag.Query().
		Where(
			posttag.TagIDEQ(tagID),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-tag by tag ID: %v", err)
		return nil, err
	}

	postIDs := make([]uint32, 0, len(posts))
	for _, post := range posts {
		if post.PostID == nil {
			continue
		}

		postIDs = append(postIDs, *post.PostID)
	}
	return postIDs, nil
}

func (r *PostTagRepo) ListPostIDsByTagIDs(ctx context.Context, tagIDs []uint32) ([]uint32, error) {
	if len(tagIDs) == 0 {
		return nil, nil
	}

	posts, err := r.entClient.Client().PostTag.Query().
		Where(
			posttag.TagIDIn(tagIDs...),
		).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-tag by tag IDs: %v", err)
		return nil, err
	}

	postIDs := make([]uint32, 0, len(posts))
	for _, post := range posts {
		if post.PostID == nil {
			continue
		}

		postIDs = append(postIDs, *post.PostID)
	}
	return postIDs, nil
}

// IsExist 判断指定帖子和标签的关联关系是否存在，通常在创建或更新帖子标签时调用，避免重复关联
func (r *PostTagRepo) IsExist(ctx context.Context, postID, tagID uint32) (bool, error) {
	exist, err := r.entClient.Client().PostTag.Query().
		Where(
			posttag.TagIDEQ(tagID),
			posttag.PostIDEQ(postID),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query post-tag exist failed: %s", err.Error())
		return false, err
	}
	return exist, nil
}

// CountPosts 统计指定标签关联的帖子数量，通常在查询标签详情时调用
func (r *PostTagRepo) CountPosts(ctx context.Context, tagID uint32) (int, error) {
	count, err := r.entClient.Client().PostTag.Query().
		Where(posttag.TagIDEQ(tagID)).
		Count(ctx)
	if err != nil {
		r.log.Errorf("failed to count post-tag by tag ID: %v", err)
		return 0, err
	}
	return count, nil
}

// CountTags 统计指定帖子关联的标签数量，通常在查询帖子详情时调用
func (r *PostTagRepo) CountTags(ctx context.Context, postID uint32) (int, error) {
	count, err := r.entClient.Client().PostTag.Query().
		Where(posttag.PostIDEQ(postID)).
		Count(ctx)
	if err != nil {
		r.log.Errorf("failed to count post-tag by post ID: %v", err)
		return 0, err
	}
	return count, nil
}

// Delete 删除指定帖子和标签的关联关系，通常在更新帖子标签时调用，删除旧的关联关系
func (r *PostTagRepo) Delete(ctx context.Context, tx *ent.Tx, postID, tagID uint32) error {
	delBuilder := tx.PostTag.Delete().
		Where(
			posttag.TagIDEQ(tagID),
			posttag.PostIDEQ(postID),
		)
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(posttag.TenantIDEQ(tid))
	}
	if _, err := delBuilder.Exec(ctx); err != nil {
		r.log.Errorf("delete post-tag failed: %s", err.Error())
		return err
	}
	return nil
}
