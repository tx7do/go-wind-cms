# go-wind-cms 自定义 OpenSearch 镜像
#
# 唯一目的：在官方 opensearchproject/opensearch 镜像之上安装 analysis-smartcn
# 中文分词器插件，供前台全文检索（posts 索引，title/summary/content 字段）使用。
#
# 基线版本钉死 2.19.6（2.x 最新稳定线）：
#   - smartcn 自 ES fork 起即存在于该线，安装路径最确定；
#   - 2.12 起 OpenSearch 内置 Elasticsearch 客户端兼容模式，配合 go.mod 中
#     elastic/go-elasticsearch/v9（bootstrap 默认不设 EnableCompatibilityMode，
#     即客户端接受 OpenSearch 的 application/json 响应），核心 index/search API 可用；
#   - 版本兼容性详见集成方案风险项 #2，需在阶段 3 实测 IndexDocument/Search/
#     CreateIndexTemplate 等调用，失败则回退 analysis-icu + ngram 或升级到 3.x。
#
# 升级到 OpenSearch 3.x 留待 v2：需重新验证 smartcn 可装性与 v9 客户端兼容性。
#
# 构建：docker compose build opensearch
FROM docker.io/opensearchproject/opensearch:2.19.6

RUN ./bin/opensearch-plugin install analysis-smartcn
