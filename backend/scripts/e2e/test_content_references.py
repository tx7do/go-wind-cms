"""RELATION 引用解析 E2E：GET /app/v1/content-references/resolve。

链路：
  1. admin 建双语言文章 + submit/approve 发布；另建一篇仅草稿
  2. resolve 已发布文章 → 摘要（标题按 languageCode 命中/回退、slug、urlPath）
  3. 草稿/不存在/category:tag/垃圾形态 → 静默缺失；重复 ref 去重
  4. page:1（种子已发布页）可解析

前置：core/admin(:6600)/app(:6700) 已启动且清代理。
用法：python test_content_references.py
"""
import json
import subprocess
import time

from common import AdminApi, admin_login

APP_BASE = "http://localhost:6700"
passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def app_get(path):
    """app(:6700) GET，绕代理；返回 (status, body dict)。"""
    r = subprocess.run(
        ["curl", "-s", "--noproxy", "*", "--max-time", "15", path],
        capture_output=True, text=True, timeout=30,
    )
    if r.returncode != 0:
        print(f"[app_get] curl rc={r.returncode} stderr={r.stderr[:200]!r}")
    try:
        return 200, json.loads(r.stdout)
    except json.JSONDecodeError:
        return 0, {"_raw": r.stdout[:200]}


hdr = admin_login("admin", "admin")
api = AdminApi(hdr)
TS = int(time.time())

# ===== 1. 造数：双语言已发布文章 + 草稿文章 =====
r = api("/admin/v1/posts", "POST", {"data": {
    "status": "DRAFT",
    "categoryIds": [],
    "translations": [
        {"languageCode": "zh-CN", "title": f"E2E 引用目标中文 {TS}",
         "content": "<p>target</p>", "slug": f"e2e-ref-target-zh-{TS}"},
        {"languageCode": "en-US", "title": f"E2E Ref Target EN {TS}",
         "content": "<p>target</p>", "slug": f"e2e-ref-target-en-{TS}"},
    ],
}})
pub_id = r.get("id")
check("建双语言文章", bool(pub_id), f"id={pub_id}")
api(f"/admin/v1/posts/{pub_id}/review/submit", "POST", {"postId": pub_id})
r = api(f"/admin/v1/posts/{pub_id}/review/approve", "POST", {"postId": pub_id})
check("审核通过→发布", r.get("status") == "POST_STATUS_PUBLISHED",
      f"status={r.get('status')}")

r = api("/admin/v1/posts", "POST", {"data": {
    "status": "DRAFT",
    "categoryIds": [],
    "translations": [{"languageCode": "zh-CN", "title": f"E2E 草稿不可引用 {TS}",
                      "content": "<p>draft</p>"}],
}})
draft_id = r.get("id")
check("建草稿文章", bool(draft_id), f"id={draft_id}")

if not pub_id or not draft_id:
    print(f"\n结果: {passed} passed, {failed} failed（无法继续）")
    raise SystemExit(1)

# ===== 2. 解析已发布文章（语言偏好/回退）=====
st, body = app_get(
    f"{APP_BASE}/app/v1/content-references/resolve"
    f"?refs=post:{pub_id}&languageCode=zh-CN")
items = body.get("items") or []
check("已发布文章解析 200+1 项", st == 200 and len(items) == 1,
      f"status={st} items={len(items)}")
if items:
    it = items[0]
    check("标题命中 zh-CN", it.get("title") == f"E2E 引用目标中文 {TS}",
          f"title={it.get('title')}")
    check("entity/slug 正确",
          it.get("targetEntityType") == "post" and it.get("targetId") == pub_id,
          f"type={it.get('targetEntityType')} id={it.get('targetId')}")
    check("slug 取偏好语言翻译", it.get("slug") == f"e2e-ref-target-zh-{TS}",
          f"slug={it.get('slug')}")
    check("urlPath 拼装", it.get("urlPath") == f"/post/e2e-ref-target-zh-{TS}",
          f"urlPath={it.get('urlPath')}")
    check("ref 回显", it.get("ref") == f"post:{pub_id}")

st, body = app_get(
    f"{APP_BASE}/app/v1/content-references/resolve"
    f"?refs=post:{pub_id}&languageCode=en-US")
items = body.get("items") or []
check("en-US 偏好命中英文标题+slug",
      items and items[0].get("title") == f"E2E Ref Target EN {TS}"
      and items[0].get("slug") == f"e2e-ref-target-en-{TS}",
      f"title={items[0].get('title') if items else None}")

st, body = app_get(f"{APP_BASE}/app/v1/content-references/resolve?refs=post:{pub_id}")
items = body.get("items") or []
check("无语言偏好回退首个非空", items and bool(items[0].get("title")),
      f"title={items[0].get('title') if items else None}")

# ===== 3. 不可见/非法引用静默缺失 =====
st, body = app_get(
    f"{APP_BASE}/app/v1/content-references/resolve"
    f"?refs=post:{draft_id}&refs=post:99999999&refs=category:1&refs=tag:2"
    f"&refs=garbage&refs=post:abc&refs=media_asset:3")
check("全部不可解析 → 200 空列表", st == 200 and body.get("items") == [],
      f"status={st} items={body.get('items')}")

# ===== 4. 去重与混合解析 =====
st, body = app_get(
    f"{APP_BASE}/app/v1/content-references/resolve"
    f"?refs=post:{pub_id}&refs=post:{pub_id}&refs=page:1&languageCode=zh-CN")
items = body.get("items") or []
check("重复 ref 去重 + page 解析", st == 200 and len(items) == 2,
      f"status={st} items={len(items)}")
pages = [it for it in items if it.get("targetEntityType") == "page"]
check("种子已发布页面可解析", pages and pages[0].get("slug")
      and pages[0].get("urlPath", "").startswith("/page/"),
      f"page={pages[0] if pages else None}")

# ===== 5. 清理 =====
api(f"/admin/v1/posts/{pub_id}", "DELETE")
api(f"/admin/v1/posts/{draft_id}", "DELETE")
time.sleep(0.5)
st, body = app_get(f"{APP_BASE}/app/v1/content-references/resolve?refs=post:{pub_id}")
check("删除后引用不再解析", st == 200 and body.get("items") == [],
      f"status={st}")

print(f"\n结果: {passed} passed, {failed} failed")
raise SystemExit(1 if failed else 0)
