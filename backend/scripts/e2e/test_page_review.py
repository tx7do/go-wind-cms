"""页面发布审批流 E2E：创建→提审→通过→发布→重复提审拒绝→下架（自动清审核态）→再提审→驳回。

前置：admin/app 服务已启动。用法：python test_page_review.py
"""
import sys

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


# 1. 创建草稿页
r = api("/admin/v1/pages", "POST", {"data": {
    "editorType": "EDITOR_TYPE_RICH_TEXT",
    "translations": [{"languageCode": "zh-CN", "title": "E2E 审批流"}],
}})
pid = r.get("id")
check("创建页面", pid is not None, f"id={pid}")

# 2. 提审
r = api(f"/admin/v1/pages/{pid}/review/submit", "POST", {"pageId": pid})
check("提审 inProgress=true", r.get("inProgress") is True)

# 3. 通过 → 发布
r = api(f"/admin/v1/pages/{pid}/review/approve", "POST", {"pageId": pid})
check("审核通过→发布", r.get("status") == "PAGE_STATUS_PUBLISHED")

# 4. 已发布再提审被拒
r = api(f"/admin/v1/pages/{pid}/review/submit", "POST", {"pageId": pid})
check("已发布再提审被拒(400)", r.get("code") == 400)

# 5. 下架（显式状态变更自动清审核态）
r = api(f"/admin/v1/pages/{pid}", "PUT", {
    "id": pid, "data": {"status": "PAGE_STATUS_DRAFT"}, "updateMask": "status",
})
check("下架", r.get("status") == "PAGE_STATUS_DRAFT")
check("下架自动清审核态", r.get("inProgress") is False)

# 6. 下架后可再提审 → 驳回退草稿
r = api(f"/admin/v1/pages/{pid}/review/submit", "POST", {"pageId": pid})
check("下架后再提审", r.get("inProgress") is True)
r = api(f"/admin/v1/pages/{pid}/review/reject", "POST", {"pageId": pid, "reason": "e2e"})
check("驳回退草稿", r.get("status") == "PAGE_STATUS_DRAFT" and r.get("inProgress") is False)

# 清理
api(f"/admin/v1/pages/{pid}", "DELETE")
print(f"页面审批流 E2E：{passed} 通过 / {failed} 失败")
sys.exit(1 if failed else 0)
