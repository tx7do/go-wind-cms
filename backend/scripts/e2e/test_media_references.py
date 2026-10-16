"""媒体查重分组 + 引用追踪 E2E。

链路：
  1. Create 两次同 file_hash 的媒体资产 → ListDuplicatedMediaAssets 出一组
  2. 建一篇文章，正文内嵌该资产 URL → reference_count=1 + ListMediaAssetReferences
     命中（post/body），再更新文章删图 → 计数归零
  3. 文章删除 → 引用边清除

前置：core/admin 已启动（无需 OIDC env；密码登录即可）。
用法：python test_media_references.py
"""
import json
import time

from common import APP_BASE, AdminApi, admin_login

hdr = admin_login("admin", "admin")
api = AdminApi(hdr)
passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


TS = int(time.time())
HASH = f"e2e-dup-hash-{TS:016d}" + "0" * 16
URL = f"http://minio.local:9000/gwc/media/e2e-ref-{TS}.png"

# ===== 1. 造两个同哈希资产（走 Create 绕过上传链路查重）=====
ids = []
for i in range(2):
    r = api("/admin/v1/media-assets", "POST", {"data": {
        "filename": f"e2e-dup-{TS}-{i}.png",
        "title": f"E2E dup {i}",
        "type": "IMAGE",
        "mimeType": "image/png",
        "size": 1024,
        "url": URL.replace(".png", f"-{i}.png"),
        "storagePath": f"media/e2e-dup-{TS}-{i}.png",
        "fileHash": HASH,
    }})
    ids.append(r.get("id"))
check("创建两个资产", all(ids), f"ids={ids}")

# ===== 2. 重复组查询 =====
r = api("/admin/v1/media-duplicates", "GET")
groups = r.get("groups") or []
mine = next((g for g in groups if g.get("fileHash") == HASH), None)
check("重复组查询包含新造组", mine is not None, f"groups={len(groups)}")
if mine:
    check("组内条数=2", int(mine.get("count") or 0) == 2, f"count={mine.get('count')}")
    check("组内含两个资产", len(mine.get("assets") or []) == 2)

# ===== 3. 引用抽取：建文章内嵌资产 URL =====
asset_url = URL.replace(".png", "-0.png")  # 与资产 0 的 url 完全一致
post = api("/admin/v1/posts", "POST", {"data": {
    "status": "DRAFT",
    "categoryIds": [],
    "translations": [{
        "languageCode": "zh-CN",
        "title": f"E2E 引用追踪 {TS}",
        "content": f'<p>图一<img src="{asset_url}" alt="a"></p><p>md: ![b]({asset_url})</p>',
    }],
}})
post_id = post.get("id")
check("建文章", bool(post_id), f"id={post_id}")

time.sleep(1)  # 钩子同步执行，保险等待
r = api(f"/admin/v1/media-assets/{ids[0]}", "GET")
check("正文内嵌 → reference_count>=1", int(r.get("referenceCount") or 0) >= 1,
      f"referenceCount={r.get('referenceCount')}")

r = api(f"/admin/v1/media-assets/{ids[0]}/references", "GET")
items = r.get("items") or []
hit = next((x for x in items if x.get("sourceId") == post_id), None)
check("引用关系命中文章", hit is not None, f"items={len(items)}")
if hit:
    check("引用位置=body", hit.get("field") == "body", f"field={hit.get('field')}")
    check("源标题回填", bool(hit.get("sourceTitle")), f"title={hit.get('sourceTitle')}")

# ===== 4. 更新文章去掉图片 → 计数归零 =====
r = api(f"/admin/v1/posts/{post_id}", "PUT", {"data": {
    "translations": [{
        "languageCode": "zh-CN",
        "title": f"E2E 引用追踪 {TS}",
        "content": "<p>没有图了</p>",
    }],
}}, )
check("更新文章去图", r.get("id") == post_id, f"resp={json.dumps(r, ensure_ascii=False)[:80]}")

time.sleep(1)
r = api(f"/admin/v1/media-assets/{ids[0]}", "GET")
check("去图后 reference_count=0", int(r.get("referenceCount") or 0) == 0,
      f"referenceCount={r.get('referenceCount')}")

# ===== 5. 删文章 → 引用边清除 =====
api(f"/admin/v1/posts/{post_id}", "DELETE")
time.sleep(1)
r = api(f"/admin/v1/media-assets/{ids[0]}/references", "GET")
items = r.get("items") or []
check("删文章后引用边清除", not any(x.get("sourceId") == post_id for x in items),
      f"items={len(items)}")

# ===== 10. 页面 sections 内嵌图引用（IMAGE 直链值 + RICH_TEXT 内嵌 img）=====
page = api("/admin/v1/pages", "POST", {"data": {
    "status": "PUBLISHED",
    "translations": [{"languageCode": "zh-CN", "title": f"E2E sections {TS}",
                      "slug": f"e2e-sections-{TS}"}],
    "sections": [
        {"type": "SECTION_TYPE_IMAGE",
         "translations": [{"languageCode": "zh-CN",
                           "content": {"src": asset_url, "alt": "a"}}]},
        {"type": "SECTION_TYPE_RICH_TEXT",
         "translations": [{"languageCode": "zh-CN",
                           "content": {"html": f'<p><img src="{asset_url}"></p>'}}]},
    ],
}})
page_id = page.get("id")
check("建页面（含 IMAGE+RICH_TEXT 区块）", bool(page_id), f"id={page_id}")

time.sleep(1)
r = api(f"/admin/v1/media-assets/{ids[0]}", "GET")
check("sections 内嵌图 → reference_count>=1", int(r.get("referenceCount") or 0) >= 1,
      f"referenceCount={r.get('referenceCount')}")

r = api(f"/admin/v1/media-assets/{ids[0]}/references", "GET")
items = r.get("items") or []
page_hits = [x for x in items if x.get("sourceType") == "page" and x.get("sourceId") == page_id]
check("引用关系命中页面", len(page_hits) >= 1, f"pageHits={len(page_hits)}")
if page_hits:
    check("引用位置=section", any(x.get("field") == "section" for x in page_hits),
          f"fields={[x.get('field') for x in page_hits]}")

api(f"/admin/v1/pages/{page_id}", "DELETE")
time.sleep(1)
r = api(f"/admin/v1/media-assets/{ids[0]}/references", "GET")
items = r.get("items") or []
check("删页面后引用边清除", not any(x.get("sourceId") == page_id and x.get("sourceType") == "page"
                                 for x in items), f"items={len(items)}")

# ===== 11. 内容建模前台：模型绑定分类 → 文章 customFields → 公开模型查询 =====
model = api("/admin/v1/content-models", "POST", {"data": {
    "name": f"E2E 模型 {TS}",
    "code": f"e2e-model-{TS}",
    "fields": [
        {"name": "price", "type": "FIELD_TYPE_NUMBER", "label": "价格"},
        {"name": "spec", "type": "FIELD_TYPE_RICHTEXT", "label": "规格"},
        {"name": "cover", "type": "FIELD_TYPE_IMAGE", "label": "封面"},
    ],
}})
model_id = model.get("id")
check("建内容模型（3 字段）", bool(model_id), f"id={model_id}")

cat = api("/admin/v1/categories", "POST", {"data": {
    "name": f"E2E 模型分类 {TS}",
    "slug": f"e2e-model-cat-{TS}",
    "contentModelId": model_id,
}})
cat_id = cat.get("id")
check("建绑定模型的分类", bool(cat_id), f"id={cat_id}")

asset_url_0 = URL.replace(".png", "-0.png")
cpost = api("/admin/v1/posts", "POST", {"data": {
    "status": "DRAFT",
    "categoryIds": [cat_id],
    "customFields": {
        "price": "99.5",
        "spec": "<b>旗舰</b>规格",
        # IMAGE 字段按校验约定存 media_asset:<id> 引用（非 URL）
        "cover": f"media_asset:{ids[0]}",
    },
    "translations": [{"languageCode": "zh-CN", "title": f"E2E 建模文章 {TS}",
                      "content": "<p>正文</p>"}],
}})
cpost_id = cpost.get("id")
check("建带 customFields 的文章", bool(cpost_id),
      f"id={cpost_id} resp={json.dumps(cpost, ensure_ascii=False)[:160]}")
# 发布走审核链（直改 status 会被审批流拦截）
api(f"/admin/v1/posts/{cpost_id}/review/submit", "POST", {"postId": cpost_id})
api(f"/admin/v1/posts/{cpost_id}/review/approve", "POST", {"postId": cpost_id})

# 公开面：/app/v1/posts/{id}/content-model 返回模型与字段定义
def app_get(path):
    import urllib.request as _u
    req = _u.Request(APP_BASE + path)
    try:
        with _u.urlopen(req, timeout=15) as resp:
            return resp.status, json.loads(resp.read().decode() or "{}")
    except Exception as e:  # noqa: BLE001 —— 断言脚本要打印失败
        print(f"  [app_get {path}] {e}")
        return 0, {}

st, cm = app_get(f"/app/v1/posts/{cpost_id}/content-model")
check("公开模型查询 200", st == 200, f"status={st}")
check("命中绑定模型", cm.get("id") == model_id, f"id={cm.get('id')} vs {model_id}")
fnames = sorted((f.get("name") or "") for f in (cm.get("fields") or []))
check("字段定义全量返回", fnames == ["cover", "price", "spec"], f"fields={fnames}")

# 公开文章详情携带 customFields
st2, pj = app_get(f"/app/v1/posts/{cpost_id}")
check("公开文章详情 200", st2 == 200, f"status={st2}")
cf = pj.get("customFields") or {}
check("详情携带 customFields", cf.get("price") == "99.5" and cf.get("spec", "").startswith("<b>"),
      f"keys={sorted(cf.keys())}")

# 公开媒体端点：IMAGE 引用 media_asset:<id> 换 URL（渲染链）
st3, ma = app_get(f"/app/v1/media-assets/{ids[0]}")
check("公开媒体资产 200", st3 == 200, f"status={st3}")
check("媒体资产带 url", bool(ma.get("url")), f"url={str(ma.get('url'))[:40]}")

api(f"/admin/v1/posts/{cpost_id}", "DELETE")
api(f"/admin/v1/categories/{cat_id}", "DELETE")
api(f"/admin/v1/content-models/{model_id}", "DELETE")

# 清理重复资产
for i in ids:
    if i:
        api(f"/admin/v1/media-assets/{i}", "DELETE")

print(f"\n通过 {passed} / 失败 {failed}")
raise SystemExit(1 if failed else 0)
