# 内容建模（Content Modeling）设计方案

> 状态：**待评审，未实现**
> 对应审计缺口：A1（README 声称"可视化自定义内容模型，支持文本、数字、富文本、图片、文件、关联等字段类型"，全栈缺失）
> 日期：2026-08-22

## 1. 背景与现状

### 1.1 README 声称
- `README.md:99` / `README.en-US.md:99`：可视化自定义内容模型，支持文本、数字、富文本、图片、文件、关联等字段类型
- `README.en-US.md:101`：分类支持绑定内容模型

### 1.2 现状（全栈缺失）
- 无 `field_definition` / `content_model` 实体、proto、服务、UI
- 仅有的半成品：
  - `custom_fields`（`map[string]string` JSON）已存在于 Post/Page/Category 三实体 schema，但 repo 层虽读写（`SetCustomFields`），服务/UI 层从不解释，等同死字段
  - `SiteSetting` 实体（`site_setting.go:57-102`）有一套完整的字段定义列模式（type 枚举/label/description/placeholder/options/is_required/validation_regex），可作为 `field_definition` 的克隆模板
  - `Section` 实体的 `type` 枚举 + `config`（语言无关）+ `content`（按语言）双 JSON map + 翻译表 fallback，是字段定义+翻译的存储模式参照

### 1.3 为什么不改造 Section
- Section 绑定 `page_id`（页面组合/渲染关注点），而内容建模需绑定 Category 并跨 Post/Page/Category
- Section 的 `type` 枚举（GALLERY/CAROUSEL/DIVIDER/SPACIER…）是渲染布局提示，非字段类型
- Section 前端目前忽略 type/config，只用 `content.body`（`section-edit-view.state.ts:161,256`），改造它比新建更大

## 2. 设计目标
- 可视化定义内容模型：每个模型含一组字段定义（名称/类型/校验规则）
- 模型绑定到 Category（README 要求），其下 Post/Page 继承模型字段
- 内容编辑器按模型定义动态渲染字段
- 字段含 6 类：text、number、richtext、image、file、relation
- 字段定义的 label/description/options 支持多语言翻译
- 严格 tenant 隔离与越权校验（relation/file/image 跨实体引用）

## 3. 数据模型

### 3.1 新增实体：`content_model`
内容模型定义，绑定到 Category（一个 Category 可绑定一个 model）。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint32 (mixin AutoIncrementId) | |
| tenant_id | uint32 (mixin TenantID) | 租户隔离 |
| created_at/updated_at/deleted_at | time (mixin TimeAt) | |
| creator_id/updater_id | uint32 (mixin OperatorID) | |
| name | string | 模型名称（语言无关，后端标识） |
| code | string | 模型编码，tenant 内唯一 |
| description | string | 模型描述（语言无关） |

索引：`(tenant_id, code)` unique。

### 3.2 新增实体：`content_model_translation`
模型名称/描述的多语言翻译。参照 `section_translation` 模式（composite unique `(content_model_id, language_code)` + `ListTranslationsWithFallback`）。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint32 (mixin) | |
| content_model_id | uint32 | FK → content_model |
| language_code | string | 语言 |
| name | string | 翻译后名称 |
| description | string | 翻译后描述 |

索引：`(content_model_id, language_code)` unique。

### 3.3 新增实体：`field_definition`
字段定义，归属 content_model。列模式克隆自 `SiteSetting`，type 枚举扩展为 6 类。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint32 (mixin) | |
| content_model_id | uint32 | FK → content_model |
| tenant_id | uint32 (mixin) | |
| created_at/... | time (mixin) | |
| name | string | 字段名（语言无关，键名） |
| type | enum | **text/number/richtext/image/file/relation**（扩展自 SiteSetting 的 9 值枚举） |
| sort_order | uint32 (mixin SortOrder) | 字段排序 |
| is_required | bool | 是否必填 |
| validation_regex | string | 校验正则（仅 text 类适用） |
| options | JSON `map[string]string` | select 类的选项（本方案暂不含 select，预留） |
| relation_config | JSON `*RelationConfig` | **仅 relation 类型**：目标实体类型、是否跨租户、过滤条件 |

`RelationConfig` 结构（参照 `seo.go`/`role_metadata.go` 的 typed-JSON mixin 模式）：
```go
type RelationConfig struct {
    TargetEntityType string   // "post" | "page" | "category" | "tag"（白名单）
    AllowCrossTenant bool     // 是否允许跨租户引用
    FilterCategoryID uint32   // 可选：限定目标 category
}
```

索引：`(content_model_id)`，`(content_model_id, name)` unique。

### 3.4 新增实体：`field_definition_translation`
字段 label/description/placeholder 的多语言翻译。同样克隆 section_translation 模式。

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint32 (mixin) | |
| field_definition_id | uint32 | FK → field_definition |
| language_code | string | |
| label | string | 翻译后标签 |
| description | string | 翻译后描述 |
| placeholder | string | 翻译后占位符 |

索引：`(field_definition_id, language_code)` unique。

### 3.5 Category 绑定模型
Category 表新增字段 `content_model_id uint32`（optional，nullable）。表示该分类下内容使用指定模型。**不建 ent edge**（与现有内容实体无 edge 的惯例一致）。

### 3.6 值存储：升级 `custom_fields`
现有 Post/Page/Category 的 `custom_fields` 是 `map[string]string`。对 text/number 可字符串化，但 **image/file 需存 media_asset 引用（结构化）、relation 需存目标实体引用（结构化）**。

**方案**：将 `custom_fields` 列类型从 `map[string]string` 升级为 `map[string]any`（或 typed-JSON proto，参照 seo）。proto 侧 `CustomFields` 改为 `map[string]string`（仍字符串化引用 ID）或结构化 `map[string]*Fieldvalue`。

权衡：`map[string]string` 保持现状，image/file/relation 的值统一编码为字符串化的引用 ID（如 `"media_asset:123"`、`"post:456"`），服务层解析。这避免 schema 列类型变更，但失去类型安全。**推荐**：升级为 `map[string]any`，值用结构化 proto `FieldValue`（oneof: string_value/number_value/media_ref/relation_ref），参照 `seo` 的 typed-JSON mixin。

### 3.7 为什么不用新建 field_value 实体
独立 field_value 表（field_definition_id + content_entity_type + content_entity_id + value）虽结构化自然，但：
- 内容列表查询需对每行 N 个字段 JOIN，性能差
- custom_fields 已存在且已读写，复用零迁移（除列类型升级）
- 现有 seo/role_metadata 已用 typed-JSON mixin 存结构化行内数据，有先例

## 4. API 设计（proto）

### 4.1 admin proto（`api/protos/admin/service/v1/i_content_model.proto`）
```proto
service ContentModelService {
  rpc List (pagination.PagingRequest) returns (ListContentModelResponse) {}
  rpc Get (GetContentModelRequest) returns (ContentModel) {}
  rpc Create (CreateContentModelRequest) returns (ContentModel) {}
  rpc Update (UpdateContentModelRequest) returns (ContentModel) {}
  rpc Delete (DeleteContentModelRequest) returns (google.protobuf.Empty) {}
  // 列出模型下的字段定义
  rpc ListFieldDefinitions (ListFieldDefinitionRequest) returns (ListFieldDefinitionResponse) {}
  // 翻译相关：参照 i_translator 的 GetTranslation/CreateTranslation/UpdateTranslation 模式
}
```
FieldDefinition 的 CRUD 并入 ContentModelService 或独立 FieldDefinitionService（推荐独立，因字段定义数量多）。

### 4.2 公开读取（app BFF）
**不暴露 content_model / field_definition 给公开端**。公开端读取 Post/Page/Category 时，`custom_fields` 随实体返回（已有），前端按字段定义渲染——但字段定义本身是后台配置，公开端只需值，不需 schema。因此 **app BFF 无需新增端点**，只需确保 `custom_fields` 在公开 List/Get 响应中返回（当前已返回）。

### 4.3 Category 绑定
Category 的 proto/admin 已有 CRUD。新增 `content_model_id` 字段后，admin Category 编辑器加一个"绑定内容模型"下拉。

## 5. 服务层（core）

### 5.1 `content_model_service.go`
- CRUD（参照 site_service.go 模式）
- ListFieldDefinitions：按 content_model_id 列字段
- 翻译：CreateTranslation/UpdateTranslation/GetTranslation，参照 `section_service.go` 的翻译模式 + `ListTranslationsWithFallback`

### 5.2 字段值校验（关键安全点）
Post/Page/Category 的 Create/Update 在写入 `custom_fields` 前，**服务层必须校验值符合绑定模型的字段定义**：
- text/number：按 validation_regex 校验，number 校验数值范围
- image/file：解析 media_ref，校验 media_asset 属于同 tenant、存在、类型匹配（image 类型只接受 image media）
- relation：解析 relation_ref，校验：
  - 目标实体类型在 RelationConfig.TargetEntityType 白名单内
  - 目标实体存在
  - **若 RelationConfig.AllowCrossTenant==false，目标实体 tenant_id 必须等于当前 tenant_id**（防越租户引用）
  - 若 RelationConfig.FilterCategoryID 设定，目标 post 必须属于该 category

这些校验在 `post_service.go`/`page_service.go`/`category_service.go` 的 Create/Update 方法里，调 `fieldDefinitionRepo.ValidateValues(ctx, contentModelID, customFields, tenantID)`。校验失败返回 BadRequest。

### 5.3 字段定义查询缓存
内容编辑器需知道字段定义来渲染表单。admin BFF 的 ContentModelService.ListFieldDefinitions 提供此查询。公开端不需。

## 6. 管理后台 UI

### 6.1 内容模型 builder（`views/app/content/content_model/`）
- 列表页：列出 content_model（name/code/绑定 category 数量）
- 编辑页：
  - 基本信息（name/code/description + 翻译 tab）
  - 字段定义列表（可增删拖拽排序）：
    - 每行：name、type 下拉（6 类）、is_required、validation_regex（type=text 时显示）
    - type=relation 时：RelationConfig 编辑器（目标实体类型下拉、AllowCrossTenant 开关）
    - 翻译 tab：label/description/placeholder 按语言编辑
- 参照 `site_setting` 的表单结构（已有 type/label/description/options/is_required/validation_regex 的渲染）

### 6.2 Category 编辑器加绑定
`views/app/content/category/edit/` 加"内容模型"下拉，绑定 content_model_id。

### 6.3 内容编辑器动态字段
Post/Page 编辑器（`views/app/content/post/edit/`、`page/edit/`）：
- 读取当前 category 绑定的 content_model 的 field_definitions
- 按定义动态渲染字段表单：
  - text → Input
  - number → InputNumber
  - richtext → Editor（复用现有 `adapter/component/Editor`）
  - image → MediaPicker（复用现有媒体选择器）
  - file → MediaPicker（file 类型）
  - relation → EntityPicker（按 RelationConfig.TargetEntityType 渲染对应实体的选择器，受 AllowCrossTenant 约束）
- 值读写 `formData.customFields[name]`

## 7. 翻译层
- content_model_translation、field_definition_translation 克隆 section_translation 的存储与 fallback 模式
- 翻译写入：admin 编辑器的翻译 tab 调 CreateTranslation/UpdateTranslation
- 翻译读取：内容编辑器调 GetTranslation（admin 侧，参照 section 的 GetTranslation 白名单）
- **custom_fields 的值不翻译**——值是内容数据，翻译的是字段定义的 label/description（UI 元素）

## 8. 安全考量

| 风险 | 缓解 |
|------|------|
| relation 越租户引用 | 服务层校验目标实体 tenant_id == 当前 tenant_id（除非 AllowCrossTenant）|
| file/image 越权引用 | 校验 media_asset 属同 tenant、存在、类型匹配 |
| 任意 type 枚举注入 | type 列枚举在 schema 层约束，proto 层枚举校验 |
| 模型逃逸 tenant | content_model/field_definition 都有 TenantID mixin，查询走 tenant 隔离 |
| 公开端泄漏字段定义 | app BFF 不暴露 content_model 服务，只随实体返回 custom_fields 值 |

## 9. 实施清单（分阶段）

### 阶段 1：后端骨架
- [ ] ent schema：content_model、content_model_translation、field_definition、field_definition_translation
- [ ] Category 加 content_model_id 字段
- [ ] custom_fields 列类型升级 map[string]string → map[string]any（+ proto FieldValue typed-JSON）
- [ ] `make ent` 重新生成
- [ ] core service：content_model_service（CRUD + 翻译 + ListFieldDefinitions）
- [ ] core service：字段值校验 `ValidateValues`（text/number/richtext/image/file/relation 全覆盖）
- [ ] proto：admin i_content_model.proto + 字段定义 proto
- [ ] `make api` + 手写装配登记（wiring.go / make register）+ admin BFF 注册
- [ ] 编译验证

### 阶段 2：管理后台 UI
- [ ] content_model builder（列表 + 编辑 + 字段定义 CRUD + 翻译 tab）
- [ ] Category 编辑器加模型绑定下拉
- [ ] Post/Page 编辑器动态字段渲染（6 类控件）
- [ ] 媒体选择器/实体选择器对接

### 阶段 3：公开端
- [ ] 确认 app BFF Post/Page/Category List/Get 返回 custom_fields（已返回）
- [ ] 前台应用按字段渲染（可选，取决于前台是否消费自定义字段）

## 10. 不在本方案范围
- Section 的 type/config/content 不改造（保持页面组合用途）
- 公开端不暴露 content_model 服务
- custom_fields 值不翻译（只翻译字段定义的 UI 元素）
- 不建 ent edge（与现有内容实体无 edge 惯例一致）

## 11. 风险与待决
- **custom_fields 列类型升级**涉及 Post/Page/Category 三实体 schema 变更 + proto 变更，需 `make ent` + `make api`，可能影响已有数据（map[string]string → map[string]any 的迁移）
- **relation 跨实体引用**是最大安全复杂点，需严格白名单 + tenant 校验，建议阶段 1 先实现 text/number/richtext，relation 留阶段 1.5
- **FieldValue proto 的 oneof 设计**需对齐 seo 的 typed-JSON mixin 实现，需确认该 mixin 支持 oneof（或退回字符串化引用）
- Category 绑定模型后，已有 Post/Page 无 custom_fields 值，需定义空值行为
