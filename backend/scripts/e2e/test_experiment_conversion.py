"""实验转化归因 E2E：建实验→匿名开页（参与+cookie）→带 cookie 提交表单→统计出转化。

前置：admin/app 服务已启动；demo 数据含 about 页（id=2，FORM 区块 sectionId=3，
字段 name/email/message 均 required）。用法：python test_experiment_conversion.py
"""
import json
import re
import sys

from common import AdminApi, admin_login, app_curl

FORM_PAGE_ID = 2
FORM_SECTION_ID = 3
FORM_DATA = {"name": "E2E", "email": "e2e@test.com", "message": "conversion test"}

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


# 0. 清理目标页的遗留实验（多实验并存取 ID 最小首个，残留会抢走分配）
r = api(f"/admin/v1/experiments?pageId={FORM_PAGE_ID}&pageSize=50")
for it in r.get("items") or []:
    xid = it.get("id")
    if it.get("status") == "EXPERIMENT_STATUS_RUNNING":
        api(f"/admin/v1/experiments/{xid}/status", "PUT", {"id": xid, "status": "EXPERIMENT_STATUS_ENDED"})
    api(f"/admin/v1/experiments/{xid}", "DELETE")
    print(f"  🧹 清理遗留实验 {xid}")

# 1. 建纯转化型实验（变体无 sections 覆盖）并上线
r = api("/admin/v1/experiments", "POST", {"data": {
    "name": "E2E-CONV",
    "pageId": FORM_PAGE_ID,
    "assignmentType": "ASSIGNMENT_TYPE_RANDOM",
    "variants": [{"name": "A", "weight": 100, "sortOrder": 1}],
}})
eid = r.get("id")
check("创建实验", eid is not None, f"id={eid}")
r = api(f"/admin/v1/experiments/{eid}/status", "PUT",
        {"id": eid, "status": "EXPERIMENT_STATUS_RUNNING"})
check("上线", r.get("status") == "EXPERIMENT_STATUS_RUNNING")

# 2. 匿名开页：参与记录 + gwc_vid cookie
payload, set_cookie = app_curl("/app/v1/pages/by-slug/about", user_agent="e2e")
m = re.search(r"gwc_vid=([^;\r\n]+)", set_cookie or "")
vid = m.group(1) if m else None
check("匿名开页种 cookie", vid is not None)

# 3. 带 cookie 提交表单（请求体顶层 pageId/sectionId/data，data 全字符串）
payload, _ = app_curl("/app/v1/form-submissions", "POST",
                      body={"pageId": FORM_PAGE_ID, "sectionId": FORM_SECTION_ID, "data": FORM_DATA},
                      cookie=f"gwc_vid={vid}")
check("表单提交成功", '"id"' in payload, payload[:150])

import time
time.sleep(2)  # 转化归因为后台协程（3s 超时上限），留出执行窗口

# 4. 统计：参与≥1 且 转化≥1
# 归因为后台协程，重试查询直到出现或超时
import time
conv_ok = False
items = []
for _ in range(10):
    r = api(f"/admin/v1/experiments/{eid}/variant-stats")
    items = r.get("items") or []
    norm = [{**it, "participants": int(it.get("participants") or 0),
             "conversions": int(it.get("conversions") or 0)} for it in items]
    if any(it["conversions"] >= 1 and it["participants"] >= 1 for it in norm):
        conv_ok = True
        break
    time.sleep(1)
check("转化归因（参与+转化）", conv_ok, json.dumps(norm, ensure_ascii=False))

# 清理
api(f"/admin/v1/experiments/{eid}/status", "PUT", {"id": eid, "status": "EXPERIMENT_STATUS_ENDED"})
api(f"/admin/v1/experiments/{eid}", "DELETE")
print(f"实验转化归因 E2E：{passed} 通过 / {failed} 失败")
sys.exit(1 if failed else 0)
