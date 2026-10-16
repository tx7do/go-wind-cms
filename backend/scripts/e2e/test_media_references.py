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

from common import AdminApi, admin_login

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

# 清理重复资产
for i in ids:
    if i:
        api(f"/admin/v1/media-assets/{i}", "DELETE")

print(f"\n通过 {passed} / 失败 {failed}")
raise SystemExit(1 if failed else 0)
