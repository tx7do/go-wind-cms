package server

import (
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	gqlTransport "github.com/tx7do/kratos-transport/transport/graphql"

	graphqlV1 "go-wind-cms/api/gen/go/graphql"
	graphqlResolver "go-wind-cms/app/app/service/internal/graphql"
	"go-wind-cms/app/app/service/internal/service"
)

// NewGraphqlServer C 端只读内容 GraphQL 服务器（gqlgen 基座）。
//
// 数据语义与 REST app BFF 对齐：仅已发布内容、翻译按 locale 回退；
// 解析器复用 app BFF 服务（internal/service），不直连 core。
//
// 环境变量：
//   - GRAPHQL_ADDR 监听地址（默认 ":6702"）；显式设为 off/disabled 停用。
//     路径 /graphql（POST 查询、GET ?query= 快速试探、内置 Introspection）。
//
// 契约：backend/api/graphql/schema.graphql；生成产物：api/gen/go/graphql。
func NewGraphqlServer(
	ctx *bootstrap.Context,
	postService *service.PostService,
	pageService *service.PageService,
	categoryService *service.CategoryService,
	tagService *service.TagService,
) *gqlTransport.Server {
	addr := strings.TrimSpace(os.Getenv("GRAPHQL_ADDR"))
	switch strings.ToLower(addr) {
	case "off", "disabled", "false", "0":
		return nil
	case "":
		addr = ":6702"
	}

	exec := graphqlV1.NewExecutableSchema(graphqlV1.Config{
		Resolvers: graphqlResolver.NewResolver(postService, pageService, categoryService, tagService),
	})

	srv := gqlTransport.NewServer(gqlTransport.WithAddress(addr))
	// 上游 Handle() 挂的是 gqlgen 裸 Server（零 transport，一切请求都报
	// "transport not supported"），这里走 HandlePrefix 挂 NewDefaultServer
	//（GET/POST/Introspection/APQ 全量 transport）。
	srv.HandlePrefix("/graphql", handler.NewDefaultServer(exec))

	ctx.NewLoggerHelper("graphql-server/app-service").Infof("graphql server configured at %s (endpoint: /graphql)", addr)
	return srv
}
