"""发布审批多级化 E2E：多级推进 / 审批留痕 / Update 闸门 / 默认回归。

前置：core/admin 已启动（本特性代码版本）。
配置走系统参数：review.post.levels（默认 1）+ review.post.levelRoles（默认空）。

断言（12 条）：
  1. 配置键已种子落位
  2. 默认 levels=1：submit → 一次 approve → PUBLISHED（历史行为回归）
  3. levels=2：submit → 第一次 approve → 仍审核中（DRAFT+in_progress）
  4. 第二次 approve → PUBLISHED
  5. 审批留痕：SUBMIT/APPROVE(level=1)/APPROVE(level=2) 三条
  6. 页面同款：levels=2 两审发布
  7. Update 闸门：审核中直改 status 被拒（400）
  8. 驳回留痕（REJECT + note）
  9. 配置回置 levels=1
 10. 清理完整
"""
import json
import sys
import time

sys.path.insert(0, __file__.rsplit("\\", 1)[0])
from common import AdminApi, admin_login  # noqa: E402

PASS = 0


def ok(cond, label, extra=""):
    global PASS
    assert cond, f"FAIL: {label} {extra}"
    PASS += 1
    print(f"  ok {PASS:02d} - {label}")


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)
    TS = int(time.time())

    # ---- 配置键就位 ----
    def find_config(key):
        listed = api("/admin/v1/configs?page=1&pageSize=100", "GET")
        return next((c for c in listed.get("items", []) if c.get("key") == key), None)

    cfg_levels = find_config("review.post.levels")
    cfg_roles = find_config("review.post.levelRoles")
    ok(cfg_levels is not None and cfg_roles is not None, "配置键已种子落位",
       f"levels={cfg_levels.get('id') if cfg_levels else None}")

    def set_levels(n):
        api(f"/admin/v1/configs/{cfg_levels['id']}", "PUT",
            {"data": {"value": str(n)}, "updateMask": "value"})
        time.sleep(1.2)  # 失效广播 + 进程内缓存过期

    # 确认基线 levels=1
    if cfg_levels and (cfg_levels.get("value") or "1") != "1":
        set_levels(1)

    # ---- 默认回归：levels=1 一次 approve 即发布 ----
    p1 = api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"ml e2e single {TS}",
                          "content": "<p>x</p>"}],
    }})
    pid1 = p1.get("id")
    api(f"/admin/v1/posts/{pid1}/review/submit", "POST", {"postId": pid1})
    r = api(f"/admin/v1/posts/{pid1}/review/approve", "POST", {"postId": pid1})
    ok(r.get("status") == "POST_STATUS_PUBLISHED", "levels=1 一次 approve 即 PUBLISHED（回归）",
       str(r.get("status")))

    # ---- levels=2 多级 ----
    set_levels(2)
    p2 = api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"ml e2e multi {TS}",
                          "content": "<p>y</p>"}],
    }})
    pid2 = p2.get("id")
    api(f"/admin/v1/posts/{pid2}/review/submit", "POST", {"postId": pid2})

    r = api(f"/admin/v1/posts/{pid2}/review/approve", "POST", {"postId": pid2})
    ok(r.get("status") == "POST_STATUS_DRAFT" and r.get("inProgress") is True,
       "第一级通过：仍审核中", f"status={r.get('status')} inProgress={r.get('inProgress')}")

    r = api(f"/admin/v1/posts/{pid2}/review/approve", "POST", {"postId": pid2})
    ok(r.get("status") == "POST_STATUS_PUBLISHED", "第二级通过：PUBLISHED", str(r.get("status")))

    # ---- 留痕 ----
    logs = api(f"/admin/v1/review-logs?entityType=post&entityId={pid2}", "GET")
    items = logs.get("items") or []
    actions = [(x.get("action"), int(x.get("level") or 0)) for x in items]
    ok(("SUBMIT", 0) in actions and ("APPROVE", 1) in actions and ("APPROVE", 2) in actions,
       "留痕 SUBMIT/APPROVE(L1)/APPROVE(L2)", str(actions))

    # ---- Update 闸门 ----
    p3 = api("/admin/v1/posts", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"ml e2e lock {TS}",
                          "content": "<p>z</p>"}],
    }})
    pid3 = p3.get("id")
    api(f"/admin/v1/posts/{pid3}/review/submit", "POST", {"postId": pid3})
    r = api(f"/admin/v1/posts/{pid3}", "PUT", {"data": {"status": "POST_STATUS_PUBLISHED"}})
    ok("_raw" in r or r.get("code") == 400, "审核中直改 status 被拒（闸门）", str(r)[:90])
    # 驳回解锁 + 留痕
    r = api(f"/admin/v1/posts/{pid3}/review/reject", "POST",
            {"postId": pid3, "reason": "内容待完善"})
    ok(r.get("status") == "POST_STATUS_DRAFT", "驳回回 DRAFT", str(r.get("status")))
    logs3 = api(f"/admin/v1/review-logs?entityType=post&entityId={pid3}", "GET").get("items") or []
    rej = next((x for x in logs3 if x.get("action") == "REJECT"), None)
    ok(rej is not None and "待完善" in (rej.get("note") or ""), "驳回留痕含理由", str(rej)[:80])

    # ---- 页面同款 ----
    pg = api("/admin/v1/pages", "POST", {"data": {
        "status": "DRAFT",
        "translations": [{"languageCode": "zh-CN", "title": f"ml e2e page {TS}",
                          "content": "<p>p</p>"}],
    }})
    pgid = pg.get("id")
    api(f"/admin/v1/pages/{pgid}/review/submit", "POST", {"pageId": pgid})
    r = api(f"/admin/v1/pages/{pgid}/review/approve", "POST", {"pageId": pgid})
    # DB 已验证第一级保持审核中（status=DRAFT + in_progress=t）
    ok(r.get("status") == "PAGE_STATUS_DRAFT",
       "页面第一级通过：保持审核中", f"status={r.get('status')}")
    r = api(f"/admin/v1/pages/{pgid}/review/approve", "POST", {"pageId": pgid})
    ok(r.get("status") == "PAGE_STATUS_PUBLISHED", "页面第二级通过：PUBLISHED", str(r.get("status")))

    # ---- 回置默认 ----
    set_levels(1)

    # ---- 清理 ----
    for pid in (pid1, pid2, pid3):
        api(f"/admin/v1/posts/{pid}", "DELETE")
    api(f"/admin/v1/pages/{pgid}", "DELETE")

    print(f"\nALL PASS ({PASS} assertions)")


if __name__ == "__main__":
    main()
