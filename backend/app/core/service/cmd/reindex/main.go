// reindex 全量重建 OpenSearch 搜索索引。
//
// 背景搜索索引只由写路径的 asynq "search.reindex" 任务逐条维护（帖子
// 创建/更新时触发），存量数据（如 demo SQL 直灌）永远没有索引。本命令
// 遍历 DB 中全部帖子，复用与写路径完全相同的处理器
//（SearchService.ReindexPost，含 SystemViewer 特权、tenant_id 取自 DB、
// 非 PUBLISHED 清理残留）把已发布文章及其翻译写入 posts 索引。
//
// 用法（在 service 目录下运行，配置取自 ./configs）：
//
//	go run ./cmd/reindex
package main

import (
	"context"
	"fmt"
	"os"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/app/core/service/internal/data/client"
	"go-wind-cms/app/core/service/internal/service"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
	appViewer "go-wind-cms/pkg/entgo/viewer"
	"go-wind-cms/pkg/serviceid"
	"go-wind-cms/pkg/task"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "reindex failed:", err)
		os.Exit(1)
	}
}

func run() error {
	// 与服务进程一致的引导流程：加载 ./configs、初始化日志，再构建 Context
	if err := bConfig.LoadBootstrapConfig("./configs"); err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	cfg := bConfig.GetBootstrapConfig()
	if cfg == nil {
		return fmt.Errorf("bootstrap config is nil")
	}
	appInfo := &conf.AppInfo{
		Project: serviceid.ProjectName,
		AppId:   serviceid.CoreService,
		Version: "1.0.0",
	}
	logger := bLogger.NewLoggerProvider(cfg.Logger, appInfo)

	ctx := bootstrap.NewContextWithParam(context.Background(), appInfo, cfg, logger)
	ctx.RegisterCustomConfig("Authenticator", &authenticationV1.AuthenticatorOptionWrapper{})

	entClient, cleanupEnt, err := client.NewEntClient(ctx)
	if err != nil {
		return fmt.Errorf("init ent client: %w", err)
	}
	defer cleanupEnt()

	esClient, cleanupES, err := client.NewElasticSearchClient(ctx)
	if err != nil {
		return fmt.Errorf("init opensearch client: %w", err)
	}
	defer cleanupES()

	searchRepo := data.NewSearchRepo(ctx, esClient)
	postRepo := data.NewPostRepo(ctx, entClient,
		data.NewPostTranslationRepo(ctx, entClient),
		data.NewPostCategoryRepo(ctx, entClient),
		data.NewPostTagRepo(ctx, entClient),
	)
	searchService := service.NewSearchService(ctx, searchRepo, postRepo, nil)

	// 全量帖子 ID（含草稿/归档：ReindexPost 对非 PUBLISHED 会清理 ES 残留文档）。
	// ent privacy 要求请求带 ViewerContext，此处与 reindex 处理器一致注入 SystemViewer。
	listCtx := appViewer.NewSystemViewerContext(ctx.Context())
	postIDs, err := entClient.Client().Post.Query().IDs(listCtx)
	if err != nil {
		return fmt.Errorf("list post ids: %w", err)
	}

	fmt.Printf("reindexing %d posts...\n", len(postIDs))
	processed := 0
	for _, id := range postIDs {
		if err := searchService.ReindexPost(task.SearchReindexTaskType, &task.SearchReindexPayload{
			Entity: "post",
			ID:     id,
			Op:     "index",
		}); err != nil {
			return fmt.Errorf("reindex post %d: %w", id, err)
		}
		processed++
	}
	fmt.Printf("done: %d posts processed (non-published cleaned up from index)\n", processed)
	return nil
}
