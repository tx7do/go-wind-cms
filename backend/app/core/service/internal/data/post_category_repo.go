package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	entCrud "github.com/tx7do/go-crud/entgo"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/postcategory"
)

type PostCategoryRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewPostCategoryRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *PostCategoryRepo {
	repo := &PostCategoryRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("post-category/repo/core-service"),
	}

	repo.init()

	return repo
}

func (r *PostCategoryRepo) init() {
}

// Create 创建帖子分类关联关系，通常在创建帖子时调用
func (r *PostCategoryRepo) Create(ctx context.Context, postID, categoryID uint32) error {
	if postID == 0 || categoryID == 0 {
		return nil
	}

	_, err := r.entClient.Client().PostCategory.Create().
		SetPostID(postID).
		SetCategoryID(categoryID).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		r.log.Errorf("failed to create post-category: %v", err)
		return err
	}
	return nil
}

// BatchCreate 批量创建帖子分类关联关系，通常在创建或更新帖子时调用
func (r *PostCategoryRepo) BatchCreate(ctx context.Context, tx *ent.Tx, postID uint32, categoryIDs []uint32) error {
	if postID == 0 || len(categoryIDs) == 0 {
		return nil
	}

	builders := make([]*ent.PostCategoryCreate, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		builders = append(builders, tx.PostCategory.Create().
			SetPostID(postID).
			SetCategoryID(id).
			SetCreatedAt(time.Now()),
		)
	}

	err := tx.PostCategory.CreateBulk(builders...).Exec(ctx)
	if err != nil {
		r.log.Errorf("failed to batch create post-category: %v", err)
		return err
	}
	return nil
}

// ListCategoryIDs 列出指定帖子关联的所有分类ID，通常在查询帖子详情时调用
func (r *PostCategoryRepo) ListCategoryIDs(ctx context.Context, postID uint32) ([]uint32, error) {
	if postID == 0 {
		return nil, nil
	}

	categories, err := r.entClient.Client().PostCategory.Query().
		Where(postcategory.PostIDEQ(postID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-category by post ID: %v", err)
		return nil, err
	}

	categoryIDs := make([]uint32, 0, len(categories))
	for _, c := range categories {
		if c == nil || c.CategoryID == nil {
			continue
		}
		categoryIDs = append(categoryIDs, *c.CategoryID)
	}
	return categoryIDs, nil
}

func (r *PostCategoryRepo) ListCategoryIDsByPostIDs(ctx context.Context, postIDs []uint32) ([]uint32, error) {
	if len(postIDs) == 0 {
		return nil, nil
	}

	categories, err := r.entClient.Client().PostCategory.Query().
		Where(postcategory.PostIDIn(postIDs...)).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-category by post IDs: %v", err)
		return nil, err
	}

	categoryIDs := make([]uint32, 0, len(categories))
	for _, c := range categories {
		if c == nil || c.CategoryID == nil {
			continue
		}
		categoryIDs = append(categoryIDs, *c.CategoryID)
	}
	return categoryIDs, nil
}

// ListPostIDs 列出指定分类关联的所有帖子ID，通常在查询分类详情时调用
func (r *PostCategoryRepo) ListPostIDs(ctx context.Context, categoryID uint32) ([]uint32, error) {
	if categoryID == 0 {
		return nil, nil
	}

	posts, err := r.entClient.Client().PostCategory.Query().
		Where(postcategory.CategoryIDEQ(categoryID)).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-category by category ID: %v", err)
		return nil, err
	}

	postIDs := make([]uint32, 0, len(posts))
	for _, p := range posts {
		if p.PostID == nil {
			continue
		}
		postIDs = append(postIDs, *p.PostID)
	}
	return postIDs, nil
}

func (r *PostCategoryRepo) ListPostIDsByCategoryIDs(ctx context.Context, categoryIDs []uint32) ([]uint32, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}

	posts, err := r.entClient.Client().PostCategory.Query().
		Where(postcategory.CategoryIDIn(categoryIDs...)).
		All(ctx)
	if err != nil {
		r.log.Errorf("failed to query post-category by category IDs: %v", err)
		return nil, err
	}

	postIDs := make([]uint32, 0, len(posts))
	for _, p := range posts {
		if p.PostID == nil {
			continue
		}
		postIDs = append(postIDs, *p.PostID)
	}
	return postIDs, nil
}

// CleanCategories 删除指定帖子关联的所有帖子分类关联关系，通常在更新帖子分类时调用
func (r *PostCategoryRepo) CleanCategories(ctx context.Context, tx *ent.Tx, postID uint32) error {
	delBuilder := tx.PostCategory.Delete().Where(postcategory.PostIDEQ(postID))
	// 租户作用域：仅删除本租户关联，避免跨租户删他人数据（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(postcategory.TenantIDEQ(tid))
	}
	_, err := delBuilder.Exec(ctx)
	if err != nil {
		r.log.Errorf("failed to delete post-category by post ID: %v", err)
		return err
	}
	return nil
}

// CleanPosts 删除指定分类关联的所有帖子分类关联关系，通常在删除分类时调用
func (r *PostCategoryRepo) CleanPosts(ctx context.Context, tx *ent.Tx, categoryID uint32) error {
	delBuilder := tx.PostCategory.Delete().Where(postcategory.CategoryIDEQ(categoryID))
	// 租户作用域：仅删除本租户关联，避免跨租户删他人数据（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(postcategory.TenantIDEQ(tid))
	}
	_, err := delBuilder.Exec(ctx)
	if err != nil {
		r.log.Errorf("failed to delete post-category by category ID: %v", err)
		return err
	}
	return nil
}

// IsExist 判断指定帖子和分类的关联关系是否存在，通常在创建或更新帖子分类时调用，避免重复关联
func (r *PostCategoryRepo) IsExist(ctx context.Context, postID, categoryID uint32) (bool, error) {
	exist, err := r.entClient.Client().PostCategory.Query().
		Where(
			postcategory.PostIDEQ(postID),
			postcategory.CategoryIDEQ(categoryID),
		).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query post-category exist failed: %s", err.Error())
		return false, err
	}
	return exist, nil
}

// CountPosts 统计指定分类关联的帖子数量，通常在查询分类详情时调用
func (r *PostCategoryRepo) CountPosts(ctx context.Context, categoryID uint32) (int, error) {
	count, err := r.entClient.Client().PostCategory.Query().
		Where(postcategory.CategoryIDEQ(categoryID)).
		Count(ctx)
	if err != nil {
		r.log.Errorf("query post-category count by category ID failed: %s", err.Error())
		return 0, err
	}
	return count, nil
}

// CountCategories 统计指定帖子关联的分类数量，通常在查询帖子详情时调用
func (r *PostCategoryRepo) CountCategories(ctx context.Context, postID uint32) (int, error) {
	count, err := r.entClient.Client().PostCategory.Query().
		Where(postcategory.PostIDEQ(postID)).
		Count(ctx)
	if err != nil {
		r.log.Errorf("query post-category count by post ID failed: %s", err.Error())
		return 0, err
	}
	return count, nil
}

// Delete 删除指定帖子和分类的关联关系，通常在更新帖子分类时调用，删除旧的关联关系
func (r *PostCategoryRepo) Delete(ctx context.Context, tx *ent.Tx, postID, categoryID uint32) error {
	delBuilder := tx.PostCategory.Delete().Where(
		postcategory.PostIDEQ(postID),
		postcategory.CategoryIDEQ(categoryID),
	)
	// 租户作用域：仅删除本租户关联，避免跨租户删他人数据（按 hasTenant 条件加）
	if tid, hasTenant := maybeTenantFromViewer(ctx); hasTenant {
		delBuilder.Where(postcategory.TenantIDEQ(tid))
	}
	_, err := delBuilder.Exec(ctx)
	if err != nil {
		r.log.Errorf("failed to delete post-category: %v", err)
		return err
	}
	return nil
}
