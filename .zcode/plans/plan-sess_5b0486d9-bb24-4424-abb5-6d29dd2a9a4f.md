# 彻底拆分计数器（Option 2）

## 范围决策（未答复问题的默认）

- **迁移**：`likes`（post）/ `like_count`（comment）→ 独立 `interaction_counter` 表。这俩有 InteractionService ledger 写入路径，是平移。
- **删除**：`visits`(post/page/site)、`dislike_count`/`reply_count`(comment)、`comment_count`(post) 直接从 schema/proto/前端展示全移除。这些列恒为 0、零业务写入，是死代码。后续若要浏览量/点踩作为新功能，counter 表已就位，只新增写入路径即可。
- 这是纯重构，不含新功能开发。若审批时你要求把 visits 也实现，我调整范围补 RecordVisit RPC。

## 关键勘察结论（已确认）

- **OpenSearch 零影响**：`PostDocument`/`PostReindexDocument`/mapping 均不含计数列（`search_repo.go:50-58,94-102`，`post_repo.go:784-792,848-856`），`dynamic:false`。拆列不碰搜索。
- **活跃计数只有 likes/like_count**：InteractionService 已有 ledger 写入（`interaction_repo.go` likePost L188-190 / unlikePost L250-252 / likeComment L302-304 / unlikeComment L360-362），迁 counter 表是平移。
- **读路径缺口**：proto 无批量读计数 RPC，只有 `LikeResponse.like_count`（单条单次）。移除缓存列后列表渲染必须新增 `GetCounts`。
- proto 无 reserved 策略，删字段须配 `reserved` 防 wire 复用。

---

## 一、后端 schema / proto

### 1.1 新增 ent `interaction_counter` 表
新文件 `backend/app/core/service/internal/data/ent/schema/interaction_counter.go`，字段：
- `target_type`（uint8，对应 TargetType 枚举值）
- `target_id`（uint32）
- `metric`（uint8，对应 CounterMetric 枚举值）
- `count`（int64，Default(0)）
- 复合 unique index `(tenant_id, target_type, target_id, metric)`
- mixin：`AutoIncrementId` / `TimeAt` / `TenantID[uint32]`（**完全照抄 post_like.go 的 mixin 写法**）

### 1.2 删除四张表的计数字段
- `post.go`：删 `visits`(L76-80) / `likes`(L82-86) / `comment_count`(L88-92)
- `comment.go`：删 `like_count`(L95-99) / `dislike_count`(L101-105) / `reply_count`(L107-111)
- `page.go`：删 `visits`(L103-107)
- `site.go`：删 `visit_count`(L83-87)

### 1.3 proto 删字段 + reserved
- `post.proto`：删 field 10/11/12，加 `reserved 10, 11, 12; reserved "visits", "likes", "comment_count";`
- `comment.proto`：删 field 21/22/23，加 reserved
- `page.proto`：删 field 20，加 reserved
- `site.proto`：删 field 12，加 reserved

### 1.4 Interaction proto 扩展
`backend/api/protos/interaction/service/v1/interaction.proto`：
- 新增枚举 `CounterMetric { COUNTER_METRIC_UNSPECIFIED=0; COUNTER_METRIC_LIKE=1; }`
- 新增 RPC `GetCounts(GetCountsRequest) returns (GetCountsResponse)`
- 新增 message（避免 enum 作 map key）：
  ```
  message MetricCount { CounterMetric metric=1; int64 count=2; }
  message CountMap { repeated MetricCount counts=1; }
  message GetCountsRequest { TargetType target_type=1; repeated uint32 target_ids=2; repeated CounterMetric metrics=3; }
  message GetCountsResponse { map<uint32, CountMap> counts=1; }  // target_id → metric→count
  ```
- `LikeResponse.like_count` **保留**（mutation 后同事务回读，供乐观更新，无需额外 round-trip）

`backend/api/protos/app/service/v1/i_interaction.proto`：新增 L2 HTTP 注解
```
rpc GetCounts(GetCountsRequest) returns (GetCountsResponse) {
  option (google.api.http) = { post: "/v1/interaction/counts:list" body: "*" };
}
```
**不进白名单**（登录强制），参照现有 Interaction RPC 先例。

---

## 二、后端 repo / service

### 2.1 `interaction_repo.go` 改写
- `likePost`/`unlikePost`/`likeComment`/`unlikeComment`：把 `tx.Post.UpdateOneID(...).AddLikes(±1)` / `AddLikeCount(±1)` 块改为 `interaction_counter` 表 upsert（同 `ent.Tx` 内，原子性不变）：
  - 查 `(tenant_id,target_type,target_id,metric=LIKE)` 行；NotFound → Create count=1；存在 → UpdateOneID.AddCount(±1)，count 归 0 则删行
- `readPostLikeCount`/`readCommentLikeCount`：改读 counter 表
- 新增 `GetCounts(ctx, targetType, targetIDs, metrics)`：批量查 counter 表，组装 `map[target_id]map[metric]count`
- 删除对 `post`/`comment` entity 计数字段的 import 依赖

### 2.2 `interaction_service.go` 新增 `GetCounts` RPC
薄转发，调 `interactionRepo.GetCounts`；`viewerUserIDFromContext(ctx)` 强制身份，未登录 → `interactionV1.ErrorUnauthorized`。

### 2.3 app 网关 `i_interaction` 转发
`backend/app/app/service/internal/service/interaction_service.go`：加 `GetCounts` 转发方法，持 `interactionV1.InteractionServiceClient`（已有）。`rest_server.go` 注册由 `buf generate` 自动覆盖（service 整体注册），**确认不进 AddWhiteList**。

### 2.4 清理失效脚手架
- `post_repo.go`/`comment_repo.go`/`page_repo.go`/`site_repo.go`：删 `SetX(0)` 初始化（字段已无）、删计数相关的 `FilterBlacklist` 调用（字段已无，黑名单空操作）；若 FilterBlacklist 还护其他非计数字段则保留、仅剔计数名。
- `make ent && make proto && wire` 全链重生成（含 OpenAPI 双产物、前端 SDK）。

---

## 三、前端

### 3.1 删 fieldMask 计数条目（7 处）
React `client-page.tsx:292` / `FeaturedPostsSection.tsx:34` / `LatestPostsSection.tsx:27`；Taro 同名 3 处；Vue `[id].vue:242`；Admin `post/index.vue:86`。两套拼写（camelCase + snake_case）都清。

### 3.2 新增 `useGetCounts` hook
React/Taro `api/hooks/interaction.ts`：`useGetCounts(targetType, targetIds[], metrics[])` → 返回 `Map<targetId, Map<metric, number>>`。参照现有 `useInteractionStatus` 的批量模式。Vue/Flutter 镜像。

### 3.3 改读路径（likes → GetCounts）
- **详情页**（React `client-page.tsx` / Taro `detail/index.tsx` / Vue `[id].vue` / Flutter `interaction_bar.dart`）：初始 `likeCount` 从 `useGetCounts(POST,[postId],[LIKE])` 取（替代 `post.likes`）；mutation 后仍用 `LikeResponse.like_count` 乐观更新。
- **列表页**（Featured/Latest/user 收藏 / Flutter carousel / search）：列表组件用当前渲染项的 ID 批量调 `useGetCounts`，计数下发给 PostCard。
- **评论树**（CommentTree / Flutter comment_item）：用所有 comment ID 批量调 `useGetCounts(COMMENT,...,[LIKE])`，`likeCount` 下发给各 CommentItem。

### 3.4 删孤儿列展示
所有 `post.visits` / `page.visits` / `site.visitCount` / `comment.dislikeCount` / `comment.replyCount` / `post.commentCount` 读取点直接删渲染行（约 30 处，四端镜像 + Flutter mock_data.dart）：
- React：`PostMetaBar.tsx`(visits prop) / `PostCard.tsx:104` / `user/page.tsx:325-327` / `CommentTree.tsx:51,136,251`(replyCount)
- Vue/Taro：同名镜像
- Flutter：`interaction_bar.dart:121`(visits) / `comment_item.dart:195`(replyCount→改用 `realChildren.length`) / `my_comments_page.dart` / `featured_carousel.dart:332` / `mock_data.dart`
- PostMetaBar / Card 类组件：删 `visits` prop 及其渲染，仅保留 `likes`（改自 GetCounts）。

### 3.5 Admin
- `post/index.vue` / `page/index.vue` / `site/index.vue` / `comment/index.vue`：删计数 column 定义（9 处）+ fieldMask 计数条目。
- locales：删 zh-CN/en-US `page.json` 各 7 条计数 i18n key（post/page visits/likes/commentCount、comment likeCount/dislikeCount/replyCount、site visitCount）。保留 dashboard 模板的 visits key（无关）。

---

## 四、数据 / 生成

### 4.1 Demo SQL
`backend/sql/postgresql-demo-data.sql`：4 段 INSERT（Comment/Site/Page/Post，L260/689/932/1911）的列清单和 VALUES 删计数列。MySQL 端干净，不动。

### 4.2 代码gen顺序
1. `make ent` — 生成 `interaction_counter` builder + migrate schema（启动自动建表/删列）
2. `make proto` — 生成 GetCounts gRPC/HTTP 代码 + 重生 openapi.yaml（双产物）
3. `wire ./...` — 重生 wire_gen.go（构造签名不变，无需改顶层 wire.Build）
4. 前端 SDK 重生（buf.react/vue/taro typescript gen + flutter dart gen）
5. 手改 fieldMask / hook / 读路径 / admin column / locales / demo SQL / mock

---

## 五、测试

### 5.1 repo 单测（扩展 `interaction_repo_test.go`）
- Like → counter 表行 count=1；Unlike → 递减，归 0 删行；重复 Like 幂等
- GetCounts 批量返回正确 map；跨 tenant 互不可见
- 现有 6 测试适配新 counter 表（断言改读 counter 表而非 `*entity.Likes`）
- 安全回归：update_mask 含 `likes` 且不传值 → 不再可能 NULL（字段已删，验证 Post 更新不报错且无副作用）

### 5.2 鉴权
未登录调 GetCounts → 401；用户 A 的 counter 查询不含 B 的数据（tenant + viewer 隔离）。

### 5.3 端到端（web-gui-tester，登录态）
- 点赞按钮态翻转、计数同步（GetCounts 初始读 + LikeResponse.like_count 乐观更新）
- 列表页点赞数正确批量渲染
- 收藏列表（useWatchedPosts，不变）
- 孤儿列展示已消失（visits/commentCount 等）

---

## 六、实施顺序

1. **P1 后端 schema/proto/repo/service**：1.1–2.4。`make ent/proto/wire` 全链过，core 单测绿。
2. **P2 前端读路径迁移**：3.1–3.4（React→Taro→Vue→Flutter，按端推进，每端可独立 build 验证）。
3. **P3 Admin + 数据清理**：3.5 + 4.1 + Flutter mock。
4. **P4 验证**：四端 build（React next build / Taro build:h5 / Vue build / Flutter analyze 需 Dart CLI，环境内可能不可用则静态检查兜底）+ repo 单测 + 端到端。

---

## 七、风险

- **前端面大**：~50 处读取点四端镜像，机械但量大。按端推进、每端独立验证降低风险。
- **GetCounts 性能**：列表批量查 counter 表需走 unique index（schema 已含），O(target_ids) 行查询；与现有 useInteractionStatus 批量模式同量级，可接受。
- **proto 破坏性**：删字段配 reserved，wire 兼容；旧客户端读不到了的新字段会报错——本次是自研全栈，前后端同步改，无第三方消费者。
- **counter 表 upsert 竞态**：同 (target,metric) 并发 Like → 依赖 unique index + ent.Tx 串行化；与现有 ledger unique 约束同源，已有先例。
- **Flutter 验证缺口**：环境无 Dart CLI 时，flutter 端只能做静态检查（import / widget. 前缀 / 括号平衡），编译验证留待有 SDK 的环境补。