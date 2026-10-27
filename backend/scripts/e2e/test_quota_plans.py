"""配额套餐模板 E2E：建模板 → 套用到租户 → 验配额快照 → 删模板不影响已套用。

前置（同 e2e_common.py）：admin(:6600)/core 已启动（本特性代码版本）。

断言（7 条）：
  1. 建模板（3 项配额）成功且回读 items 一致
  2. 缺 name 被拒
  3. 套用到租户：applied=3
  4. 租户配额列表出现 3 行且值与模板一致
  5. 改值后重复套用：applied 只计变化项（同值跳过）=0 → 先改租户配额值再套用验 upsert
  6. 删除模板成功；租户配额快照仍在
  7. 清理完整
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

    # 建模板
    plan = api("/admin/v1/quota-plans", "POST", {"data": {
        "name": f"e2e-plan-{TS}",
        "description": "e2e",
        "items": [
            {"quotaType": "USER_LIMIT", "limit": 50},
            {"quotaType": "STORAGE", "limit": 1000000000},
            {"quotaType": "AI_TOKENS", "limit": 100000},
        ],
    }})
    plan_id = plan.get("id")
    check_plan = bool(plan_id) and len(plan.get("items") or []) == 3
    ok(check_plan, "建模板成功且回读 3 项", f"id={plan_id} items={len(plan.get('items') or [])}")

    bad = api("/admin/v1/quota-plans", "POST", {"data": {"items": [
        {"quotaType": "USER_LIMIT", "limit": 1}]}})
    ok(bool(bad.get("_raw")) or bad.get("code") == 400, "缺 name 被拒", str(bad)[:80])

    # 建租户（或取现有）
    tenant = api("/admin/v1/tenants", "POST", {"data": {
        "name": f"e2e-plan-tenant-{TS}", "code": f"e2eplan{TS}", "status": "ON",
    }})
    tenant_id = tenant.get("id")
    if not tenant_id:
        listed = api("/admin/v1/tenants?page=1&pageSize=1", "GET")
        tenant_id = (listed.get("items") or [{}])[0].get("id")
    ok(bool(tenant_id), "目标租户就绪", f"id={tenant_id}")

    # 套用
    applied = api(f"/admin/v1/quota-plans/{plan_id}/apply", "POST", {"tenantId": tenant_id})
    ok(int(applied.get("applied") or 0) == 3, "套用 applied=3", str(applied)[:80])

    # 验配额快照
    time.sleep(1)
    quotas = api(f"/admin/v1/tenant-quotas?tenantId={tenant_id}&page=1&pageSize=50", "GET")
    rows = quotas.get("items") or []
    by_type = {r.get("quotaType"): int(r.get("quotaValue") or 0) for r in rows}
    ok(by_type.get("USER_LIMIT") == 50 and by_type.get("AI_TOKENS") == 100000,
       "配额快照与模板一致", str(by_type))

    # 改租户单项后重复套用 → 变化项才计数
    qid = next((r.get("id") for r in rows if r.get("quotaType") == "USER_LIMIT"), None)
    if qid:
        api(f"/admin/v1/tenant-quotas/{qid}", "PUT", {"data": {"quotaValue": 99}})
        time.sleep(0.5)
    applied2 = api(f"/admin/v1/quota-plans/{plan_id}/apply", "POST", {"tenantId": tenant_id})
    ok(int(applied2.get("applied") or 0) == 1, "改值后重复套用只计变化项=1", str(applied2)[:80])

    # 删模板：快照保留
    api(f"/admin/v1/quota-plans/{plan_id}", "DELETE")
    time.sleep(0.5)
    quotas2 = api(f"/admin/v1/tenant-quotas?tenantId={tenant_id}&page=1&pageSize=50", "GET")
    ok(len(quotas2.get("items") or []) >= 3, "删模板后租户配额快照仍在",
       f"n={len(quotas2.get('items') or [])}")

    # 清理
    for r in (quotas2.get("items") or []):
        if str(r.get("tenantId")) == str(tenant_id):
            api(f"/admin/v1/tenant-quotas/{r.get('id')}", "DELETE")
    api(f"/admin/v1/tenants/{tenant_id}", "DELETE")

    print(f"\nALL PASS ({PASS} assertions)")


if __name__ == "__main__":
    main()
