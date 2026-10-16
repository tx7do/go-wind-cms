# 内容建模 Components / 可重复字段组 / 动态区 设计

状态：**基础层已落地**（本文件同批提交）；admin 编辑器与 C 端渲染为后续切片（见 §8）。
对标：Strapi Components / Repeatable Fields / Dynamic Zones（原差距清单 P1-1）。

## 1. 概念与目标

- **组件（Component）**：具名、可复用的字段组（如 `seo_meta` = title/description/image；
  `cta` = label/link）。租户内定义一次，多个内容模型引用。
- **可重复字段组（Repeatable）**：一个 COMPONENT 字段的值是组件实例数组
  （相册、步骤、FAQ 列表）。
- **动态区（Dynamic Zone）**：一个 DYNAMIC_ZONE 字段的值是**混合组件**实例的有序数组
  （页头/CTA/段落自由拼装），是页面搭建能力的基座。

非目标（本轮不做）：跨租户组件共享、组件套组件（嵌套 Component 字段——校验器按一层
实现，嵌套留待有真实需求再开）、组件字段级权限。

## 2. 数据模型

### 2.1 新表 `component_definitions`

| 列 | 类型 | 说明 |
|---|---|---|
| id | uint32 PK | |
| code | string | 语言无关引用键，租户内唯一（uix_component_definition_tenant_code） |
| name | string | 显示名 |
| description | string? | |
| — | 通用 mixin | AutoIncrementId / TimeAt / OperatorID / TenantID |

### 2.2 `field_definitions` 复用

组件的字段**不复刻一张表**：`field_definitions` 已有 `content_model_id`（可空），
新增可空 `component_id`——两者恰好其一非空。收益：字段 CRUD/mapper/翻译表
（field_definition_translations）/排序全部复用；新唯一索引
`uix_field_definition_component_name (component_id, name)`
（PG 唯一索引对 NULL 天然放行，不影响模型字段的既有索引）。

### 2.3 FieldDefinition 扩展

- `type` 枚举 += `FIELD_TYPE_COMPONENT`、`FIELD_TYPE_DYNAMIC_ZONE`。
- 新 typed-JSON 列 `component_config`（照 `relation_config` 先例）：

```go
type ComponentConfig struct {
    ComponentCode        string   `json:"component_code,omitempty"`   // COMPONENT 必填
    Repeatable           bool     `json:"repeatable,omitempty"`       // COMPONENT：true=实例数组
    AllowedComponentCodes []string `json:"allowed_component_codes,omitempty"` // DYNAMIC_ZONE 必填，非空
}
```

## 3. 值形状（custom_fields 维持 map[string]string）

组件值以 **JSON 字符串**编码进既有 custom_fields 列——零 schema 迁移、C 端 API
形状不变（此前设计文档设想的 map[string]any 升级对本特性不再是必要条件）：

```jsonc
// COMPONENT 单实例：键的值为对象 JSON
"seo_meta": "{\"title\":\"…\",\"description\":\"…\"}"
// COMPONENT 可重复：对象数组 JSON
"gallery": "[{\"image\":\"media_asset:12\",\"caption\":\"…\"},{\"image\":\"media_asset:13\"}]"
// DYNAMIC_ZONE：有序混合数组，实例用 "_component" 判别键
"sections": "[{\"_component\":\"hero\",\"heading\":\"…\"},{\"_component\":\"cta\",\"label\":\"…\"}]"
```

权衡：组件内部值不可在 SQL 层查询（custom_fields 本就不可查询，无退化）；
C 端按字段类型 JSON.parse 后结构化渲染（§7）。

## 4. 校验（core `ContentModelRepo.ValidateValues` 扩展）

既有挂点不动（PostService.Create/Update 已调用）。对 COMPONENT/DYNAMIC_ZONE
字段值，在既有"未定义键拒绝/必填检查"之后追加：

1. 值必须是合法 JSON；单实例→对象，可重复/动态区→数组（空数组合法，is_required
   用"非空数组"判定）。
2. COMPONENT：实例的键集 ⊆ 组件字段集（未知键拒绝）；组件内 is_required 字段
   必须存在且非空；叶子值复用既有类型校验（number 可数值化、image/file/relation
   为字符串引用等）。
3. DYNAMIC_ZONE：每个元素的 `_component` 必须 ∈ allowed_component_codes；
   元素按其组件定义做与 2 相同的校验。
4. 引用的组件按 code 现查（租户内）；组件不存在/字段定义为空 → 400（模型侧
   及时暴露配置错误，不静默放行）。
5. 单实例 COMPONENT 的 is_required：对象非空即满足。

嵌套（组件字段类型为 COMPONENT）校验器显式拒绝：配置时报错，值校验不递归。

## 5. Proto 面（content.service.v1 + admin 面）

- `ComponentDefinition` 消息：id/code/name/description/createdAt + 统计只读字段
  `field_count`。CRUD 挂在 core `ContentModelService`（List/Get/Create/Update/
  Delete），admin 面新增独立 `ComponentDefinitionService` REST 面
  `/admin/v1/content-components*`（管理页独立于内容模型页）。
- `FieldDefinition`：type 枚举加两值；`component_config` 字段（protojson 整树映射，
  同 RelationConfig 的 proto 消息形态）。
- `ListFieldDefinitionsRequest` 加可选 `component_id`：为 0 时按 content_model_id
  （既有语义不变），非 0 时列组件字段。
- 组件删除保护：被任何模型的 field_definitions.component_config 引用时拒绝删除
  （先解绑再删）；字段定义的 component_id 随组件级联删除。

## 6. admin 交互（后续切片，交互定稿）

- 新页 `/content/component`：列表 + 抽屉（code/name/字段构建器——复用
  ContentModelDrawer 的字段编辑 UX，字段定义按 component_id 存取）。
- ContentModelDrawer 字段类型下拉 += 组件/动态区；选组件：code 选择器 +
  可重复开关；选动态区：组件多选（allowed codes）。
- 文章编辑器（最后做，工作量最大）：COMPONENT 字段渲染为实例列表
  （新增/删除/拖拽排序，实例内按组件字段 schema 渲染子表单）；
  DYNAMIC_ZONE 渲染为"添加组件"选择 + 有序卡片列表。

## 7. C 端渲染契约（后续切片）

`CustomFieldsSection` 对新类型 JSON.parse 后结构化渲染；解析失败回退原样字符串
（与既有容错一致）：
- COMPONENT 单实例 → 键值表；可重复 → 卡片列表。
- DYNAMIC_ZONE → 按 `_component` 分发表渲染（未注册的组件类型跳过该元素）；
  IMAGE 值沿用 media_asset 引用换 URL 链路。四端分两批（react/vue 先行，
  taro/flutter 跟进——与 custom_fields 现状渲染等级对齐）。

## 8. 切片账本

1. ✅ **基础层（本轮）**：ent（component_definitions + field_definitions 扩展）、
   proto（消息/枚举/CRUD/admin 面）、repo（组件 CRUD + ValidateValues 扩展）、
   post 保存路径接线（既有挂点自动生效）、E2E（组件 CRUD + 值校验断言）。
2. 组件管理页 + 模型抽屉新字段类型（admin）。
3. 文章编辑器组件实例编辑（admin，工作量最大）。
4. C 端渲染：react/vue → taro/flutter。
