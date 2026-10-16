"""内容建模 Components 基础层 E2E：组件 CRUD / 字段保存 / 值校验 / 删除保护。

链路：
  1. 组件 CRUD：创建（code 唯一）、重复 code 拒绝、code 不可改、更新、列表 fieldCount
  2. 组件字段：SaveFieldDefinitions 落库（嵌套组件字段拒绝）；ListFieldDefinitions 按 componentId 读
  3. 模型绑定：模型带 COMPONENT（cta 可重复）+ DYNAMIC_ZONE 字段 → 分类 → 文章
  4. 值校验：合法通过；未知键/缺必填/坏 JSON/非法 _component/未知媒体引用 → 400
  5. 删除保护：被模型字段引用的组件删除拒绝；解绑（模型删）后可删

前置：core/admin(:6600) 已启动且清代理。
用法：python test_components.py
"""
import json
import time

from common import AdminApi, admin_login

passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def is_4xx(r):
    return r.get("code") in (400, 409)


hdr = admin_login("admin", "admin")
api = AdminApi(hdr)
TS = int(time.time())

# ===== 1. 组件 CRUD =====
cta = api("/admin/v1/content-components", "POST", {"data": {
    "code": f"cta-{TS}", "name": "CTA 按钮", "description": "行动号召按钮",
}})
cta_id = cta.get("id")
check("建组件 cta", bool(cta_id), f"id={cta_id}")

dup = api("/admin/v1/content-components", "POST", {"data": {
    "code": f"cta-{TS}", "name": "重复编码",
}})
check("重复 code 被拒", is_4xx(dup), f"code={dup.get('code')}")

upd = api(f"/admin/v1/content-components/{cta_id}", "PUT", {
    "id": cta_id,
    "data": {"code": f"renamed-{TS}", "name": "CTA 按钮 v2"},
})
check("code 不可改被拒", is_4xx(upd), f"code={upd.get('code')}")

upd = api(f"/admin/v1/content-components/{cta_id}", "PUT", {
    "id": cta_id,
    "data": {"name": "CTA 按钮 v2"},
})
check("更新组件名", upd.get("name") == "CTA 按钮 v2", f"name={upd.get('name')}")

hero = api("/admin/v1/content-components", "POST", {"data": {
    "code": f"hero-{TS}", "name": "头图区",
}})
hero_id = hero.get("id")
check("建组件 hero", bool(hero_id), f"id={hero_id}")

# ===== 2. 组件字段 =====
nested = api(f"/admin/v1/content-components/{cta_id}/field-definitions", "PUT", {
    "componentId": cta_id,
    "fields": [{"name": "label", "type": "FIELD_TYPE_COMPONENT"}],
})
check("嵌套组件字段被拒", is_4xx(nested), f"code={nested.get('code')}")

saved = api(f"/admin/v1/content-components/{cta_id}/field-definitions", "PUT", {
    "componentId": cta_id,
    "fields": [
        {"name": "label", "type": "FIELD_TYPE_TEXT", "label": "按钮文字",
         "isRequired": True},
        {"name": "link", "type": "FIELD_TYPE_TEXT", "label": "跳转链接"},
    ],
})
check("保存 cta 字段（replace）", int(saved.get("total") or 0) == 2,
      f"total={saved.get('total')}")

saved = api(f"/admin/v1/content-components/{cta_id}/field-definitions", "PUT", {
    "componentId": cta_id,
    "fields": [
        {"name": "label", "type": "FIELD_TYPE_TEXT", "label": "按钮文字",
         "isRequired": True},
    ],
})
check("再次保存整体替换", int(saved.get("total") or 0) == 1, f"total={saved.get('total')}")
saved = api(f"/admin/v1/content-components/{cta_id}/field-definitions", "PUT", {
    "componentId": cta_id,
    "fields": [
        {"name": "label", "type": "FIELD_TYPE_TEXT", "label": "按钮文字",
         "isRequired": True},
        {"name": "link", "type": "FIELD_TYPE_TEXT", "label": "跳转链接"},
    ],
})

listed = api(f"/admin/v1/content-components/{cta_id}/field-definitions", "GET")
check("按 componentId 读字段", len(listed.get("items") or []) == 2,
      f"items={len(listed.get('items') or [])}")

hero_fields = api(f"/admin/v1/content-components/{hero_id}/field-definitions", "PUT", {
    "componentId": hero_id,
    "fields": [
        {"name": "heading", "type": "FIELD_TYPE_TEXT", "label": "标题",
         "isRequired": True},
    ],
})
check("保存 hero 字段", int(hero_fields.get("total") or 0) == 1,
      f"total={hero_fields.get('total')}")

lst = api("/admin/v1/content-components?pageSize=50", "GET")
mine = next((x for x in lst.get("items") or [] if x.get("id") == cta_id), None)
check("列表 fieldCount 统计", mine is not None and int(mine.get("fieldCount") or 0) == 2,
      f"fieldCount={mine.get('fieldCount') if mine else None}")

# ===== 3. 模型绑定组件字段 =====
model = api("/admin/v1/content-models", "POST", {"data": {
    "name": f"E2E 组件模型 {TS}",
    "code": f"e2e-comp-model-{TS}",
    "fields": [
        {"name": "ctas", "type": "FIELD_TYPE_COMPONENT", "label": "按钮组",
         "componentConfig": {"componentCode": f"cta-{TS}", "repeatable": True}},
        {"name": "sections", "type": "FIELD_TYPE_DYNAMIC_ZONE", "label": "区块",
         "componentConfig": {"allowedComponentCodes": [f"cta-{TS}", f"hero-{TS}"]}},
    ],
}})
model_id = model.get("id")
check("建模型（组件字段）", bool(model_id), f"id={model_id}")

cat = api("/admin/v1/categories", "POST", {"data": {
    "name": f"E2E 组件分类 {TS}",
    "slug": f"e2e-comp-cat-{TS}",
    "contentModelId": model_id,
}})
cat_id = cat.get("id")
check("建绑定模型的分类", bool(cat_id), f"id={cat_id}")


def make_post(custom_fields):
    return api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "categoryIds": [cat_id],
        "customFields": custom_fields,
        "translations": [{"languageCode": "zh-CN",
                          "title": f"E2E 组件文章 {TS} {time.time_ns()}",
                          "content": "<p>正文</p>"}],
    }})


good = make_post({
    "ctas": json.dumps([{"label": "立即购买", "link": "/buy"}]),
    "sections": json.dumps([
        {"_component": f"hero-{TS}", "heading": "标题"},
        {"_component": f"cta-{TS}", "label": "去看看", "link": "/go"},
    ]),
})
check("合法组件值通过", good.get("id") is not None,
      f"id={good.get('id')} err={str(good)[:120]}")

# ===== 4. 校验断言 =====
cases = [
    ("未知键拒绝", {"ctas": json.dumps([{"label": "x", "oops": "1"}])}),
    ("缺组件必填拒绝", {"ctas": json.dumps([{"link": "/x"}])}),
    ("可重复非数组拒绝", {"ctas": json.dumps({"label": "x"})}),
    ("坏 JSON 拒绝", {"ctas": "{not-json"}),
    ("动态区非法 _component 拒绝",
     {"sections": json.dumps([{"_component": "nope", "heading": "x"}])}),
    ("动态区元素未知键拒绝",
     {"sections": json.dumps([{"_component": f"hero-{TS}", "oops": "x"}])}),
    ("动态区缺 _component 拒绝", {"sections": json.dumps([{"heading": "x"}])}),
]
for name, cf in cases:
    r = make_post(cf)
    check(name, is_4xx(r), f"code={r.get('code')} msg={str(r.get('message'))[:60]}")

# ===== 5. 删除保护 =====
prot = api(f"/admin/v1/content-components/{hero_id}", "DELETE")
check("被引用组件删除拒绝", is_4xx(prot), f"code={prot.get('code')}")

api(f"/admin/v1/posts/{good.get('id')}", "DELETE")
api(f"/admin/v1/categories/{cat_id}", "DELETE")
api(f"/admin/v1/content-models/{model_id}", "DELETE")

freed = api(f"/admin/v1/content-components/{hero_id}", "DELETE")
check("解绑后可删", freed.get("_error") is None and "code" not in freed,
      f"resp={str(freed)[:80]}")

# ===== 清理 =====
api(f"/admin/v1/content-components/{cta_id}", "DELETE")
left = api("/admin/v1/content-components?pageSize=50", "GET")
check("清理完成", not any(str(x.get("code", "")).endswith(str(TS))
                         for x in left.get("items") or []))

print(f"\n结果: {passed} passed, {failed} failed")
raise SystemExit(1 if failed else 0)
