//go:build integration
// +build integration

// 集成测试：需要本地 OpenSearch 实例（docker-compose up opensearch）。
// 默认 `go test` 不执行本文件；手动运行：
//   go test -tags integration ./app/core/service/internal/data/ -run TestSearchRepo -v
//
// 测试覆盖：
//   1. 索引/搜索/删除基础流程
//   2. 跨 tenant_id 搜索返回空（租户隔离核心）
//   3. tid==0 / language=="" / status=="" 返回空（不 bypass）
//   4. EnsureIndexTemplate 幂等

package data

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	elasticsearchCrud "github.com/tx7do/go-crud/elasticsearch"
)

// newTestSearchRepo 构造一个连本地 OpenSearch 的 SearchRepo。
// 直接用 go-crud 的底层 option 构造 client，绕过 bootstrap.Context。
func newTestSearchRepo(t *testing.T) *SearchRepo {
	t.Helper()

	esClient, err := elasticsearchCrud.NewElasticsearchClient(
		elasticsearchCrud.WithAddresses("http://localhost:9200"),
		elasticsearchCrud.WithUsername("admin"),
		elasticsearchCrud.WithPassword("@Abcd#123456"),
		elasticsearchCrud.WithLogger(log.DefaultLogger),
	)
	require.NoError(t, err, "连接 OpenSearch 失败，请确认 docker-compose up opensearch 已启动")

	return &SearchRepo{
		esClient: esClient,
		log:      log.NewHelper(log.DefaultLogger),
	}
}

func TestSearchRepo_IndexSearchDelete(t *testing.T) {
	repo := newTestSearchRepo(t)
	ctx := context.Background()

	// 确保索引模板存在
	require.NoError(t, repo.EnsureIndexTemplate(ctx))

	// 索引一篇 tenant=1 的文档
	doc := &PostDocument{
		TenantID: "1",
		PostID:   "99001",
		Language: "zh",
		Status:   "POST_STATUS_PUBLISHED",
		Title:    "集成测试标题集成测试",
		Summary:  "摘要内容",
		Content:  "正文内容集成测试",
	}
	require.NoError(t, repo.IndexPost(ctx, doc))

	// 等待索引刷新（OpenSearch 默认 1s 刷新间隔）
	time.Sleep(2 * time.Second)

	// tenant=1 应搜到
	result, err := repo.SearchPosts(ctx, "集成测试", 1, "zh", "POST_STATUS_PUBLISHED", 0, 10)
	require.NoError(t, err)
	assert.Greater(t, result.Total, 0, "tenant=1 应能搜到自己的文档")

	// 删除
	require.NoError(t, repo.DeletePost(ctx, 99001))
	time.Sleep(2 * time.Second)

	// 删除后应搜不到
	result2, err := repo.SearchPosts(ctx, "集成测试", 1, "zh", "POST_STATUS_PUBLISHED", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, result2.Total, "删除后应搜不到")
}

func TestSearchRepo_TenantIsolation(t *testing.T) {
	repo := newTestSearchRepo(t)
	ctx := context.Background()
	require.NoError(t, repo.EnsureIndexTemplate(ctx))

	// 索引 tenant=1 的文档
	doc := &PostDocument{
		TenantID: "1",
		PostID:   "99002",
		Language: "zh",
		Status:   "POST_STATUS_PUBLISHED",
		Title:    "租户隔离测试租户隔离",
		Content:  "租户隔离内容",
	}
	require.NoError(t, repo.IndexPost(ctx, doc))
	time.Sleep(2 * time.Second)

	// tenant=2 搜 tenant=1 的文档 → 应返回空（核心隔离断言）
	result, err := repo.SearchPosts(ctx, "租户隔离", 2, "zh", "POST_STATUS_PUBLISHED", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, result.Total, "tenant=2 不应搜到 tenant=1 的文档")

	// tenant=1 自己能搜到
	result2, err := repo.SearchPosts(ctx, "租户隔离", 1, "zh", "POST_STATUS_PUBLISHED", 0, 10)
	require.NoError(t, err)
	assert.Greater(t, result2.Total, 0, "tenant=1 应能搜到自己的文档")

	// 清理
	require.NoError(t, repo.DeletePost(ctx, 99002))
}

func TestSearchRepo_NoBypassForZeroTenant(t *testing.T) {
	repo := newTestSearchRepo(t)
	ctx := context.Background()
	require.NoError(t, repo.EnsureIndexTemplate(ctx))

	// 索引 tenant=1 的文档
	doc := &PostDocument{
		TenantID: "1",
		PostID:   "99003",
		Language: "zh",
		Status:   "POST_STATUS_PUBLISHED",
		Title:    "零租户绕过测试",
		Content:  "内容",
	}
	require.NoError(t, repo.IndexPost(ctx, doc))
	time.Sleep(2 * time.Second)

	// tid==0 应返回空（不接受 SystemViewer bypass）
	result, err := repo.SearchPosts(ctx, "绕过", 0, "zh", "POST_STATUS_PUBLISHED", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, result.Total, "tid==0 必须返回空，不接受 bypass")

	// language=="" 应返回空
	result2, err := repo.SearchPosts(ctx, "绕过", 1, "", "POST_STATUS_PUBLISHED", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, result2.Total, "language==空 必须返回空")

	// status=="" 应返回空
	result3, err := repo.SearchPosts(ctx, "绕过", 1, "zh", "", 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, result3.Total, "status==空 必须返回空")

	// 清理
	require.NoError(t, repo.DeletePost(ctx, 99003))
}
